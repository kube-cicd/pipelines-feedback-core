config
=======

Architecture
------------

- [**Collector:**](./local.go) One or multiple collectors that are reading configuration data from various sources
- **pkgs/controller/ConfigurationController:** Kubernetes CRD controller. Collects and passes to DocumentStore
- **internal/config/DocumentStore:** Keeps all collected document files
- [**Provider**](./provider.go): Provides a resolved, contextual configuration


Configuration structure
-----------------------

`PFConfig` is a Custom Resource Definition that has a `.data` section, just like in a `kind: ConfigMap` or in a `kind: Secret`.

The `.data` section is a key-value flat list **with component prefixes**. For example, when a `jxscm` component wants to retrieve its configuration it will get a map like `{"token-secret-name": "some-secret", "token-secret-key": "token"}`.

> NOTICE: In the code the components are getting key-values without a prefix

```yaml
---
apiVersion: pipelinesfeedback.keskad.pl/v1alpha1
kind: PFConfig
metadata:
    name: keskad-sample-1
spec:
    jobDiscovery: {}
data:
    # {component}.{key}: {value}
    
    # jxscm is a component name
    # In jxscm implementation then we use the ConfigurationProvider in a following way to retrieve all "jxscm" prefixed keys:
    # `jx.sc.Config.FetchContextual("jxscm", pipeline.GetNamespace(), pipeline)`
    
    jxscm.token-secret-name: "some-secret"
    jxscm.token-secret-key: "token"
```

Using ConfigurationProvider in a component
------------------------------------------

```go
cfg := configurationProvider.FetchContextual("jxscm", pipeline.GetNamespace(), pipeline)

// token-secret-key => `jxscm.token-secret-key` from the configuration
// token.txt => a default value in case, when a user would not set anything in the configuration
println(cfg.GetOrDefault("token-secret-key", "token.txt"))
```

Implementing a collector
------------------------

`ConfigurationCollector` interface lets you implement a configuration source that is fetching configuration files at application startup or at runtime, when a Pipeline is reconciled by the controller.

```go
type ConfigurationCollector interface {
	contract.Pluggable
	CollectInitially() ([]*v1alpha1.PFConfig, error)
	CollectOnRequest(info contract.PipelineInfo) ([]*v1alpha1.PFConfig, error)
	SetLogger(logger logging.Logger)
}
```

It is up to you, how to read the data. It can be fetched via HTTP, parsed from YAML/JSON or whatever you implement. Most important thing is
the output format which should be the `v1alpha1.PFConfig` object containing a key-value data section.

**Cases:**
- `.metadata.namespace` not provided: Configuration will be global for Pipelines in all namespaces
- `.spec.jobDiscovery` not provided: Configuration will be for all objects in a namespace or in all namespaces
