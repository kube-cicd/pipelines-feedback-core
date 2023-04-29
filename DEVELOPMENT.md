## Concept

See documentation of each component by visiting its directory in this repository.

**Components:**
- [Receivers](./pkgs/feedback): Connector to a party that receives the feedback. Default implementation is `jxscm` which handles `Gitea`, `Gitlab`, `Github`, etc.
- [Providers](./pkgs/provider): Pipeline data providers. Default implementation is collecting labelled `kind: Job` from the cluster and parsing their status. Feel free to implement your **kinds** to support e.g. `Tekton`, `Argo Workflows` or `Jenkins X`
- [ConfigurationCollector](./pkgs/config): Provides settings & secrets to access the **Receiver** (e.g. credentials to log-in into Gitlab to post a PR update)
- [Store](./pkgs/store): State storage (key-value) that stores values not available in Kubernetes manifests. Default backends: `memory`, `redis`
