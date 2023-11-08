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
