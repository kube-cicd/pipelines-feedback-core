Installation
============

With ArgoCD
-----------

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: pipelines-feedback-batchv1
  namespace: argocd
spec:
  destination:
    namespace: my-namespace
    server: https://kubernetes.default.svc
  project: default
  source:
    chart: batchv1-chart
    helm:
      values: |
        rbac:
            resourceNames: ["my-secret-name-in-every-namespace"]
    repoURL: quay.io/pipelines-feedback

    # please pay attention to the version
    # check available options: https://quay.io/repository/pipelines-feedback/batchv1-chart?tab=tags
    targetRevision: 0.1

  syncPolicy: {}
```

Manually using Helm
-------------------

```bash
# please pay attention to the version
# check available options: https://quay.io/repository/pipelines-feedback/batchv1-chart?tab=tags
helm install pfb oci://quay.io/pipelines-feedback/batchv1-chart --version 0.1
```
