Store
=====

Stores state, controller needs a cache. There are so many events being sent to the controller, but not every one deserves processing. That's why Pipelines Feedback is using strong caching.

Choosing a store
----------------

```bash
# use a commandline switch
-s, --store string                          Sets a Store adapter (default "redis")
```

```yaml
# helm values
controller:
    adapters:
        store: redis
    deployment:
        env:
            REDIS_HOST: "redis:6379"
```

Memory
------

Stores configuration in Pod's memory. Whole cache is wiped on controller restart. There is no configuration needed.

Redis
-----

Connects to a Redis instance for a persistent cache. Configurable using environment variables.

```bash
# use commandline switch to activate 
pipelines-feedback-tekton --store redis
```

| Environment variable name | Default value  | Description       |
|---------------------------|----------------|-------------------|
| REDIS_HOST                | localhost:6379 | Host + port       |
| REDIS_DB                  | 0              | Database number   |
| REDIS_PASSWORD            |                | Optional password |
