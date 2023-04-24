Pipelines Feedback Core
=======================

> NOTICE: THIS IS A WORK IN PROGRESS, currently more a PoC. I try to make it working, design it well, then stabilize and release.

Generic Kubernetes controller watching Jobs on your cluster and notifying external systems, mainly Github, Gitlab, but not only.
There are **Feedback receivers**, **Feedback providers** and **Configuration providers**.

**Bundled Feedback Receivers:**
- [go-scm](https://github.com/jenkins-x/go-scm) (strongest Gitlab supported)

**Bundled Configuration Providers:**
- localfile (read configuration from local YAML)
- annotation (read from Kubernetes annotations)

Roadmap
-------

**First beta release:**
- [x] Reference implementation implementing `kind: Job` support
- [ ] Modular architecture (pluggable: `config`, `receiver`, `provider`, `store`)
- [ ] Split on `pkgs` and `internal` to hide internally used methods
- [ ] Configuration as CRD and as local file, inherited and merged

**First stable release:**
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

Development
-----------

[If you want to join the development or create your own controller, take a look there.](./DEVELOPMENT.md)
