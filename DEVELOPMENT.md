

## Concept

- [Receivers](./pkgs/feedback): Party that receives the feedback. Default implementation is JX SCM which handles Gitea, Gitlab, Github, etc.
- [Providers](./pkgs/provider): Pipeline data providers. Default implementation is collecting labelled `kind: Job` from the cluster and parsing their status
- [ConfigurationProvider](./pkgs/config): Configuration source that provides secrets to access the **Receiver** (e.g. credentials to log-in into Gitlab to post a PR update)
