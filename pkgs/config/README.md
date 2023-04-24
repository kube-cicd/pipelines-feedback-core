config
=======

Architecture
------------

- **Collector:** One or multiple collectors that are reading configuration data from various sources
- **pkgs/controller/ConfigurationController:** Kubernetes CRD controller. Collects and passes to DocumentStore
- **internal/config/DocumentStore:** Keeps all collected document files
- **Provider**: Provides a resolved, contextual configuration
