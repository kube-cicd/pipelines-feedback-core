Getting started guide
=====================

**Pipelines Feedback** can integrate any CI/CD with any external system to report Pipelines status into.

Feedback Provider (CI/CD system - source of Pipelines execution knowledge)
-----------------

At first please choose an implementation you wish to use.

**Known implementations:**
- [For Tekton Pipelines](https://github.com/kube-cicd/pipelines-feedback-tekton)
- [For Kubernetes `kind: Job`](https://github.com/kube-cicd/pipelines-feedback-core/tree/main/pkgs/implementation/batchjob)

Common configuration
====================

Regardless of selected controller if you are using Tekton, or plain `kind: Job` or other CI/CD there are things that are
common for all controllers based on `pipelines-feedback-core`.

[Choosing Store type](./pkgs/store/README.md)
-------------------

[Configuring Feedback Receiver (external system link)](./pkgs/feedback/USAGE.md)
---------------------

[Configuring controller Globally, per Namespace and per Pipeline](./pkgs/config/USAGE.md)
---------------------------------------------------------------

Global configuration reference
------------------------------

Pipelines Feedback Core has a core set of settings which are not associated with any _Feedback Receiver_, _Controller_, _Store_ or _Feedback Provider_. Those configuration options are called `global`.

| Name                             | Default value | Description                                                                                                                                                                                                                             |
|----------------------------------|---------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| dashboard-url                    |               | Go-template formatted URL to any dashboard e.g. OpenShift Pipelines, Tekton Dashboard, other. <br/> Example: `https://console-openshift-console.apps.my-cluster.org/k8s/ns/{{ .namespace }}/tekton.dev~v1beta1~PipelineRun/{{ .name }}` |
| logs-enabled                     | true          | Fetch logs from builds [true/false]                                                                                                                                                                                                     |
| logs-max-line-length             | 64            | How many characters a single log line could have                                                                                                                                                                                        |
| logs-max-full-length-lines-count | 10            | How many log lines should be returned                                                                                                                                                                                                   |
| logs-split-separator             | (...)         | A string that replaces ending in truncated logs                                                                                                                                                                                         | 
