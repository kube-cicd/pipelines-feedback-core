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
