config
=======

Architecture
------------

- [**Collector:**](./local.go) One or multiple collectors that are reading configuration data from various sources
- **pkgs/controller/ConfigurationController:** Kubernetes CRD controller. Collects and passes to DocumentStore
- **internal/config/DocumentStore:** Keeps all collected document files
- [**Provider**](./provider.go): Provides a resolved, contextual configuration
