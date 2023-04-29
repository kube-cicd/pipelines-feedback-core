App
===

Common bootstrap logic setting up controllers, wiring dependencies and running main loop.

Usage
-----

`PipelinesFeedbackApp` is using a [Dependency Injection (Constructor Injection) pattern](https://en.wikipedia.org/wiki/Dependency_injection) to let your controller be configured with implementations of choice.
Replace any piece by simply injecting it in the object construction.

```go
pfcApp := app.PipelinesFeedbackApp{
    JobController:    batchjob.CreateJobController(),
    ConfigController: &controller.ConfigurationController{},
}
```

```go
// Interface definition
type PipelinesFeedbackApp struct {
	// can read configuration from various sources
	ConfigCollector config.ConfigurationCollector

	JobController          *controller.GenericController
	ConfigController       *controller.ConfigurationController

	CustomFeedbackReceiver string
	CustomConfigCollector  string

	// Feedback receivers available to choose by the user. Falls back to default, embedded list if not specified
	AvailableFeedbackReceivers []feedback.Receiver

	// Config providers available to choose by the user. Falls back to default, embedded list if not specified
	AvailableConfigCollectors []config.ConfigurationCollector
}
```

Building your own controller application
----------------------------------------

Copy [main.go](../../main.go) to your project and adjust to your needs, then use `go build` to build a customized controller.
