Helm Chart
==========

Releases
--------

We do releases on Quay.io in order to be more compatible with RedHat stack and also to have cool download stats. Helm Charts are published as OCI images in a separate repository in the same organization.

- [Check Quay.io releases page](https://quay.io/organization/pipelines-feedback)

Installation
------------

### ArgoCD

```yaml
---
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: pipelines-feedback-batchv1
  namespace: argocd
spec:
  destination:
    namespace: default
    server: https://kubernetes.default.svc
  project: default
  source:
    chart: batchv1-chart
    helm:
      #values: |
      #  rbac:
      #      resourceNames: ["my-secret-name-in-every-namespace"]
    repoURL: quay.io/pipelines-feedback
    targetRevision: v0.1
  syncPolicy: {}
```

### Plain Helm from CLI

```bash
helm install pfc oci://quay.io/pipelines-feedback/batchv1-chart --version 0.0.1-latest-main
```

Security
--------

### Limit `kind: Secret` by name

`kind: PFConfig` can reference to `kind: Secret` containing Gitlab/GitHub/etc secrets. In multi-tenant environment, where each team has it's own `kind: Namespace` and Gitlab/GitHub/etc token the `kind: Secret` may have
the same name in each namespace, so the controller permissions could be easily limited with RBAC rule.

```yaml
rbac:
    secretResourceNames: ["my-gitlab-token-secret-name"]
```

### Watch Pipelines only in selected namespaces

You may want to explicitly set the list of allowed namespaces controller has access to. Proper RBAC rules would be generated for you.

There is no possibility to use e.g. labelled namespaces, only fixed namespace names are allowed due to RBAC nature in Kubernetes.

```yaml
rbac:
    bindToNamespaces: ["team-1", "team-2"]
```

### Set resources and actions

This section should be configured automatically by controller like _Tekton Pipelines Feedback_.

```yaml
rbac:
    jobRules: 
        - apiGroups: ["tekton.dev"]
          resources: ["pipelineruns"]
          verbs: ["list", "get", "watch"]
```
