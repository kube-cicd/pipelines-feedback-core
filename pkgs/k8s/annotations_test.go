package k8s_test

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/k8s"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

type TestCase struct {
	name                     string
	inputAnnotations         map[string]string
	expectedError            string
	expectedPrId             string
	expectedCommit           string
	expectedReference        string
	expectedRepoHttpsUrl     string
	expectedOrganizationName string
	expectedRepositoryName   string
	expectedIsTechnicalJob   bool
	expectedIsValid          bool
}

func TestCreateJobContextFromKubernetesAnnotations(t *testing.T) {
	testCases := []TestCase{
		{
			name: "Commit + PR id + Ref + URL is present, a full context is present. Is not a technical job",
			inputAnnotations: map[string]string{
				"pipelinesfeedback.keskad.pl/ref":            "refs/heads/main",
				"pipelinesfeedback.keskad.pl/pr-id":          "1",
				"pipelinesfeedback.keskad.pl/commit":         "2d6cc283fb5be9f963f2b70c504e4fedc6c025b8",
				"pipelinesfeedback.keskad.pl/https-repo-url": "https://github.com/kube-cicd/pipelines-feedback-core.git",
			},
			expectedError:            "",
			expectedPrId:             "1",
			expectedCommit:           "2d6cc283fb5be9f963f2b70c504e4fedc6c025b8",
			expectedReference:        "refs/heads/main",
			expectedRepoHttpsUrl:     "https://github.com/kube-cicd/pipelines-feedback-core.git",
			expectedOrganizationName: "kube-cicd",
			expectedRepositoryName:   "pipelines-feedback-core",
			expectedIsTechnicalJob:   false,
			expectedIsValid:          true,
		},
		{
			name: "Its a technical job, without SCM context",
			inputAnnotations: map[string]string{
				"pipelinesfeedback.keskad.pl/technical-job": "true",
			},
			expectedError:            "",
			expectedPrId:             "",
			expectedCommit:           "",
			expectedReference:        "",
			expectedRepoHttpsUrl:     "",
			expectedOrganizationName: "",
			expectedRepositoryName:   "",
			expectedIsTechnicalJob:   true,
			expectedIsValid:          true,
		},
		{
			name:                     "Everything is empty, not a technical job",
			inputAnnotations:         map[string]string{},
			expectedError:            "",
			expectedPrId:             "",
			expectedCommit:           "",
			expectedReference:        "",
			expectedRepoHttpsUrl:     "",
			expectedOrganizationName: "",
			expectedRepositoryName:   "",
			expectedIsTechnicalJob:   false,
			expectedIsValid:          false,
		},
		{
			name: "Invalid url",
			inputAnnotations: map[string]string{
				"pipelinesfeedback.keskad.pl/https-repo-url": "ftp://hehe",
			},
			expectedError:            "cannot create JobContext: repository url does not contain valid organization and repository names",
			expectedPrId:             "",
			expectedCommit:           "",
			expectedReference:        "",
			expectedRepoHttpsUrl:     "",
			expectedOrganizationName: "",
			expectedRepositoryName:   "",
			expectedIsTechnicalJob:   false,
			expectedIsValid:          false,
		},
		{
			name: "Invalid organization name",
			inputAnnotations: map[string]string{
				"pipelinesfeedback.keskad.pl/https-repo-url": "https://github.com",
			},
			expectedError:            "cannot create JobContext: repository url does not contain valid organization and repository names",
			expectedPrId:             "",
			expectedCommit:           "",
			expectedReference:        "",
			expectedRepoHttpsUrl:     "",
			expectedOrganizationName: "",
			expectedRepositoryName:   "",
			expectedIsTechnicalJob:   false,
			expectedIsValid:          false,
		},
	}

	for _, testCase := range testCases {
		objectMeta := metav1.ObjectMeta{
			Name:        "bakunin-1",
			Annotations: testCase.inputAnnotations,
		}
		k8sContext, err := k8s.CreateJobContextFromKubernetesAnnotations(objectMeta)
		hasUsableAnnotations, _ := k8s.HasUsableAnnotations(objectMeta)

		if err == nil {
			assert.Equal(t, testCase.expectedCommit, k8sContext.Commit)
			assert.Equal(t, testCase.expectedPrId, k8sContext.PrId)
			assert.Equal(t, testCase.expectedReference, k8sContext.Reference)
			assert.Equal(t, testCase.expectedRepoHttpsUrl, k8sContext.RepoHttpsUrl)
			assert.Equal(t, testCase.expectedOrganizationName, k8sContext.OrganizationName)
			assert.Equal(t, testCase.expectedRepositoryName, k8sContext.RepositoryName)
			assert.Equal(t, testCase.expectedIsTechnicalJob, k8sContext.IsTechnicalJob())
		} else {
			assert.Equal(t, testCase.expectedError, err.Error())
		}
		assert.Equal(t, testCase.expectedIsValid, hasUsableAnnotations)
	}
}
