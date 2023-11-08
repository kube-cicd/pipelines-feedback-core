Feedback Receiver
=================

Configures a connection to external system which will act as a Feedback Receiver.

**Examples:**
- GIT: Update commit status
- JIRA: Update release info

```bash
# select an adapter via commandline switch
-f, --feedback-receiver string              Sets a FeedbackReceiver (default "jxscm")
```

jxscm
-----

Jenkins X go-scm library provides a support for multiple GIT servers like Gitlab.com, Gitlab Self-hosted, GitHub, Gitea self-hosted, Bitbucket.
JX-SCM is configured via `kind: PFConfig` or using a JSON loaded locally at controller startup.

| Name                         | Example value                        | Description                                                                                                 |
|------------------------------|--------------------------------------|-------------------------------------------------------------------------------------------------------------|
| jxscm.token                  | glpat-blablabla                      | Plaintext access token. Avoid using this field. Use `token-secret-name` and `token-secret-key` pair instead |
| jxscm.token-secret-name      | my-secret-name                       | `kind: Secret` name placed in same namespace as `kind: PFConfig` and Pipeline is                            |
| jxscm.token-secret-key       | token                                | Name of the key in `.data` section of the `kind: Secret`                                                    |
| jxscm.git-repo-url           | https://username:password@gitlab.com | Full SCM url                                                                                                |
| jxscm.git-kind               | gitlab                               | Jenkins X go-scm git-kind parameter                                                                         |
| jxscm.git-token              | glpat-blablabla                      | Same as `token`                                                                                             |
| jxscm.git-user               | __token__                            |                                                                                                             |
| jxscm.bb-oauth-client-id     |                                      |                                                                                                             |
| jxscm.bb-oauth-client-secret |                                      |                                                                                                             |
| jxscm.progress-comment       |                                      | Go template formatted PR progress comment                                                                   |
| jxscm.finished-comment       |                                      | Go template formatted PR summary comment                                                                    |


**Example configuration:**

```yaml
---
apiVersion: pipelinesfeedback.keskad.pl/v1alpha1
kind: PFConfig
metadata:
    name: keskad-sample-1
    namespace: team-1
spec:
    jobDiscovery: {}
data:
    jxscm.token-secret-name: "some-secret"
    jxscm.token-secret-key: "token"
```
