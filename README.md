Pipelines Feedback Core
=======================

![[image](https://quay.io/repository/pipelines-feedback/batchv1?tab=tags)](https://img.shields.io/badge/container-quay.io-green.svg)
![[chart](https://quay.io/repository/pipelines-feedback/batchv1-chart?tab=tags)](https://img.shields.io/badge/chart-quay.io-green.svg)

> NOTICE: THIS IS A WORK IN PROGRESS, currently more a PoC. I try to make it working, design it well, then stabilize and release.

Generic Kubernetes controller watching Jobs on your cluster and notifying external systems, mainly Github, Gitlab, but not only.
There are **Feedback receivers**, **Feedback providers** and **Configuration providers**.

**Bundled Feedback Receivers:**
- [jxscm](https://github.com/jenkins-x/go-scm) (Github, Gitea, Gitlab, Bitbucket, etc.)

**Bundled Configuration Providers:**
- local (read configuration from local JSON file)
- crd (read from Kubernetes CRD - `kind: PFConfig`)

**Releases:**

We do releases on Quay.io in order to be more compatible with RedHat stack and also to have cool download stats. Helm Charts are published as OCI images in a separate repository in the same organization.

- [Check Quay.io releases page](https://quay.io/organization/pipelines-feedback)

Roadmap
-------

**First alpha release - 0.1:**
- [x] Reference implementation implementing `kind: Job` support
- [x] Modular architecture (pluggable: `config`, `receiver`, `provider`, `store`)
- [x] Split on `pkgs` and `internal` to hide internally used methods
- [x] Configuration as CRD and as local file, inherited and merged
- [x] Support for administrative jobs (jobs without SCM context e.g. backup jobs, identified by group-id)
- [x] Add support for logs fetching
- [x] Configuration schema support

**First beta release - 0.2:**
- [ ] Support for optional arguments in API for easier future interface extension (https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
- [x] Support Redis as a cache store
- [ ] Support for Matrix (Federated, Secure, Slack-like Messenger)

**First stable release - 1.0:**
- [ ] Freeze the API (in the code as well as in CRD)
- [ ] Document the API

**Next:**
- [ ] WebAssembly support to write Feedback receivers in language of choice

Framework
---------

This repository acts as a core library. You can easily create a Kubernetes controller for your CI/CD of choice by just implementing a simple interface we provide.
The idea of this project is to create a unified, core library that could be used with Tekton, Argo Workflows, Jenkins X, plain Kubernetes Jobs and other possible CI/CD stacks.

Implementations
---------------

- [Kubernetes batch/v1 Jobs: Reference implementation](./pkgs/implementation)

batch/v1 Jobs (jobs-feedback)
-----------------------------

This repository contains a exemplary and fully functional implementation for basic Kubernetes jobs.

[Development](./DEVELOPMENT.md)
-----------

[See development API for developing integrations or whole custom controllers](./DEVELOPMENT.md)

[Usage](./USAGE.md)
-------

For end users - `pipelines-feedback-core` is an opinionated framework, so the usage could be different depending on the controller you are going to use.

Check [generic usage](./USAGE.md) tips for common parts of every controller.
