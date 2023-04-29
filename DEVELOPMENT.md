## Developing a controller with Pipelines Feedback Core

Pipelines Feedback Core is a framework for developing your own controller, we together maintain universal and powerful tools for integrating CI/CD Pipelines with external systems.
Using this framework you can create an Open/Closed Source integration like e.g. Slack notifications or Tekton support.

See documentation of each component by visiting its directory in this repository.

**Components:**
- [Receivers](./pkgs/feedback): Connector to a party that receives the feedback. Default implementation is `jxscm` which handles `Gitea`, `Gitlab`, `Github`, etc.
- [Providers](./pkgs/provider): Pipeline data providers. Default implementation is collecting labelled `kind: Job` from the cluster and parsing their status. Feel free to implement your **kinds** to support e.g. `Tekton`, `Argo Workflows` or `Jenkins X`
- [ConfigurationCollector](./pkgs/config): Provides settings & secrets to access the **Receiver** (e.g. credentials to log-in into Gitlab to post a PR update)
- [Store](./pkgs/store): State storage (key-value) that stores values not available in Kubernetes manifests. Default backends: `memory`, `redis`

**API:**
- [Bootstrapping your own controller](./pkgs/app/README.md)
- [Contract of objects used inside Pipelines Feedback - the API](./pkgs/contract)
- [Logger interface available in your component](./pkgs/logging/interface.go)
- [Kubernetes Annotations parser - use it when parsing your custom kind/crd](./pkgs/k8s)
- [ConfigurationProvider - reads configuration data and provides it to your component](./pkgs/config/README.md) (see also the [ConfigurationProvider code](./pkgs/config/provider.go))
