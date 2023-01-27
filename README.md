Pipelines Feedback Core
=======================

Generic Kubernetes controller watching Jobs on your cluster and notifying external systems, mainly Github, Gitlab and everything what is supported by JX [go-scm](https://github.com/jenkins-x/go-scm).

> NOTICE: THIS IS A WORK IN PROGRESS

Framework
---------

This repository acts as a core library. You can easily create a Kubernetes controller for your CI/CD of choice by just implementing a simple interface we provide.
The idea of this project is to create a unified, core library that could be used with Tekton, Argo Workflows, plain Kubernetes Jobs and other possible CI/CD stacks.

Implementations
---------------

- [Kubernetes batch/v1 Jobs: Reference implementation](./pkgs/implementation)
