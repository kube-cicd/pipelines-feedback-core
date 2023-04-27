Feedback Receiver
=================

Lets you implement **an integration point with external system** in order to notify about changes coming from reconciliation loop of Pipeline resources on the cluster.

Everytime the Kubernetes controller finds e.g. `kind: Job` or `kind: PipelineRun` or other configured kind it will call your `Receiver` implementation in order to notify e.g. Github, Gitlab, JIRA, Matrix, Slack or whatever integration
is implemented and enabled.

`jxscm` is our core, reference implementation that out-of-the-box provides a pleasant experience with many GIT provides (all supported by Jenkins X GO-SCM library).

Use case: Progress update (e.g. progress bar)
---------------------------------------------

Use `UpdateProgress()` method to get known about each state change.

> NOTICE: Same state change may be triggered multiple times due to how Kubernetes handles events. Consider using `store.Store` or higher level interface `store.Operator` to keep the information about already processed events.

Use case: Alerting & Notifications
----------------------------------

`WhenCreated()`, `WhenStarted()` and `WhenFinished()` are fired exactly once, so those methods are perfect candidates to implement notifications or alerting.

Interface
---------

```go
package feedback

type Receiver interface {
	contract.Pluggable

	// UpdateProgress is called each time a status is changed
	UpdateProgress(ctx context.Context, status contract.PipelineInfo) error

	// WhenCreated is an event, when a Pipeline was created and is in Pending or already in Running state
	WhenCreated(ctx context.Context, status contract.PipelineInfo) error

	// WhenStarted is an event, when a Pipeline is started
	WhenStarted(ctx context.Context, status contract.PipelineInfo) error

	// WhenFinished is an event, when a Pipeline is finished - Failed, Errored, Aborted or Succeeded
	WhenFinished(ctx context.Context, status contract.PipelineInfo) error
}
```

contract.Pluggable interface
----------------------------

Mandatory interface. Allows to self-identify as a plugin in CLI, so the user could choose an implementation.


contract.wiring.WithInitialization interface
--------------------------------------------

Optional interface. Method `InitializeWithContext(sc *ServiceContext) error` will inject standard services to your implementation, those services includes e.g. a logger, kube config and a configuration provider.
