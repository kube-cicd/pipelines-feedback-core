package jxscm

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/azure"
	"github.com/jenkins-x/go-scm/scm/driver/bitbucket"
	"github.com/jenkins-x/go-scm/scm/driver/fake"
	"github.com/jenkins-x/go-scm/scm/driver/gitea"
	"github.com/jenkins-x/go-scm/scm/driver/github"
	"github.com/jenkins-x/go-scm/scm/driver/gitlab"
	"github.com/jenkins-x/go-scm/scm/driver/gogs"
	"github.com/jenkins-x/go-scm/scm/driver/stash"
	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/jenkins-x/go-scm/scm/transport"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// todo: Contribute to JX's GO SCM library a factory method - "NewFromConfig()" and delete this duplicate code there

// ClientOptionFunc is a function taking a client as its argument
type ClientOptionFunc func(*scm.Client)

// ErrMissingGitServerURL the error returned if you use a git driver that needs a git server URL
var ErrMissingGitServerURL = fmt.Errorf("No git serverURL was specified")

type authOptions struct {
	oauthToken   string
	clientID     string
	clientSecret string
}

// setUsername allows the username to be set
func setUsername(username string) ClientOptionFunc {
	return func(client *scm.Client) {
		if username != "" {
			client.Username = username
		}
	}
}

// ensureBBCEndpoint lets ensure we have the /api/v3 suffix on the URL
func ensureBBCEndpoint(u string) string {
	if strings.HasPrefix(u, "https://bitbucket.org") || strings.HasPrefix(u, "http://bitbucket.org") {
		return "https://api.bitbucket.org"
	}
	return u
}

// ensureGHEEndpoint lets ensure we have the /api/v3 suffix on the URL
func ensureGHEEndpoint(u string) string {
	if strings.HasPrefix(u, "https://github.com") || strings.HasPrefix(u, "http://github.com") {
		return "https://api.github.com"
	}
	// lets ensure we use the API endpoint to login
	if !strings.Contains(u, "/api/") {
		u = scm.URLJoin(u, "/api/v3")
	}
	return u
}

func NewClientFromConfig(data config.Data, gitToken string) (*scm.Client, error) {
	if repoURL := data.GetOrDefault("git-repo-url", ""); repoURL != "" {
		return factory.FromRepoURL(repoURL)
	}
	driver := data.GetOrDefault("git-kind", "")
	serverURL := strings.TrimSpace(data.GetOrDefault("git-server", ""))

	// make sure the URL is valid / issue #45
	if serverURL != "" {
		// fallback to https:// in case, when URL is defined, but does not contain the scheme
		if !strings.Contains(serverURL, "://") {
			serverURL = "https://" + serverURL
		}

		// make sure the url is valid
		if _, urlErr := url.ParseRequestURI(serverURL); urlErr != nil {
			return nil, fmt.Errorf("invalid git-server URL: %q. valid values are empty or a correctly formatted URL address", serverURL)
		}
	}

	oauthToken := data.GetOrDefault("git-token", gitToken)
	username := data.GetOrDefault("git-user", "")
	if oauthToken == "" {
		return nil, errors.New("No Git OAuth token specified")
	}

	authOptions := &authOptions{
		oauthToken: oauthToken,
	}

	clientID := data.GetOrDefault("bb-oauth-client-id", "")
	clientSecret := data.GetOrDefault("bb-oauth-client-secret", "")
	authOptions.clientID = clientID
	authOptions.clientSecret = clientSecret

	client, err := newClient(driver, serverURL, authOptions, setUsername(username))
	if driver == "" {
		driver = client.Driver.String()
	}
	return client, err
}

func newClient(driver, serverURL string, authOptions *authOptions, opts ...ClientOptionFunc) (*scm.Client, error) {
	oauthToken := authOptions.oauthToken
	if driver == "" {
		driver = "github"
	}
	var client *scm.Client
	var err error

	switch driver {
	case "azure":
		client = azure.NewDefault()
	case "bitbucket", "bitbucketcloud":
		if serverURL != "" {
			client, err = bitbucket.New(ensureBBCEndpoint(serverURL))
		} else {
			client = bitbucket.NewDefault()
		}
	case "fake", "fakegit":
		client, _ = fake.NewDefault()
	case "gitea":
		if serverURL == "" {
			return nil, ErrMissingGitServerURL
		}
		client, err = gitea.NewWithToken(serverURL, oauthToken)
	case "github":
		if serverURL != "" {
			client, err = github.New(ensureGHEEndpoint(serverURL))
		} else {
			client = github.NewDefault()
		}
	case "gitlab":
		if serverURL != "" {
			client, err = gitlab.New(serverURL)
		} else {
			client = gitlab.NewDefault()
		}
	case "gogs":
		if serverURL == "" {
			return nil, ErrMissingGitServerURL
		}
		client, err = gogs.New(serverURL)
	case "stash", "bitbucketserver":
		if serverURL == "" {
			return nil, ErrMissingGitServerURL
		}
		client, err = stash.New(serverURL)
	default:
		return nil, fmt.Errorf("Unsupported $GIT_KIND value: %s", driver)
	}
	if err != nil {
		return client, err
	}
	if oauthToken != "" {
		switch driver {
		case "azure":
			client.Client = &http.Client{
				Transport: &transport.Custom{
					Before: func(r *http.Request) {
						encoded := base64.StdEncoding.EncodeToString([]byte(":" + oauthToken))
						r.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))
					},
				},
			}
		case "gitea":
			client.Client = &http.Client{
				Transport: &transport.Authorization{
					Scheme:      "token",
					Credentials: oauthToken,
				},
			}
		case "gitlab":
			client.Client = &http.Client{
				Transport: &transport.PrivateToken{
					Token: oauthToken,
				},
			}
		case "bitbucketcloud":
			// lets process any options now so that we can populate the username
			for _, o := range opts {
				o(client)
			}
			if client.Username == "" {
				return nil, errors.Errorf("no username supplied")
			}
			if authOptions.clientID != "" && authOptions.clientSecret != "" {
				config := clientcredentials.Config{
					ClientID:     authOptions.clientID,
					ClientSecret: authOptions.clientSecret,
					TokenURL:     "https://bitbucket.org/site/oauth2/access_token",
				}
				client.Client = config.Client(context.Background())
				return client, nil
			}
			// BB App Password / PAT
			client.Client = &http.Client{
				Transport: &transport.BasicAuth{
					Username: client.Username,
					Password: oauthToken,
				},
			}
			return client, nil
		default:
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: oauthToken},
			)
			client.Client = oauth2.NewClient(context.Background(), ts)
		}
	}
	for _, o := range opts {
		o(client)
	}
	return client, err
}
