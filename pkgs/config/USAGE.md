Configuring Pipelines Feedback
==============================

Pipelines Feedback has configurable [Configuration Collectors](./interface.go). `pipeline-feedback-core` provides few core collector types, and your controller of choice can implement a new collector not described there.

```bash
# configurable via CLI
-c, --config-provider string                Sets a ConfigCollector - possible to set multiple, comma separated, without spaces) (default "local")
```

```yaml
# helm values
controller:
    adapters:
        config: local
```

local
-----

Loads a configuration in JSON format from file specified by environment variable `CONFIG_PATH`.

> Note: This is a global configuration, it will be inherited into all Pipelines


```bash
export CONFIG_PATH=/etc/keskad/pipelines-feedback.json
```

A JSON file should consist only a key-value dictionary.

**Example configuration file:**

```json
{
  "jxscm.git-kind": "gitlab",
  "jxscm.git-server": "example.org",
  "jxscm.token": "glpat-hello",
  "jxscm.git-user": "__token__",
  "dashboard-url": "https://console-openshift-console.apps.my-cluster.org/k8s/ns/{{ .namespace }}/tekton.dev~v1beta1~PipelineRun/{{ .name }}"
}
```

kubernetes (always turned on)
-----------------------------

```yaml
---
apiVersion: pipelinesfeedback.keskad.pl/v1alpha1
kind: PFConfig
metadata:
    name: keskad-sample-1
    namespace: team-1
spec:
    jobDiscovery:
        # https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#resources-that-support-set-based-requirements
        labelSelector:
            - matchLabels:
                  some-label: value
            - matchExpressions:
                  - key: "some-other-label-name"
                    operator: In
                    values: ["a", "b"]
data:
    "jxscm.git-kind": "gitlab",
    "jxscm.git-server": "example.org",
    "jxscm.token": "glpat-hello",
    "jxscm.git-user": "__token__",
    "dashboard-url": "https://console-openshift-console.apps.my-cluster.org/k8s/ns/{{ .namespace }}/tekton.dev~v1beta1~PipelineRun/{{ .name }}"
```
