package cli

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/app"
	"github.com/spf13/cobra"
	"strings"
)

func NewRootCommand(app *app.PipelinesFeedbackApp) *cobra.Command {
	command := &cobra.Command{
		Use:   "pipelines-feedback",
		Short: "Runs a Kubernetes controller that observes Pipeline status and reports the status to the external system",
		RunE: func(command *cobra.Command, args []string) error {
			return app.Run()
		},
	}

	app.CustomFeedbackReceiver = ""
	app.CustomConfigCollector = ""

	//
	// FeedbackReceiver and ConfigCollector can be enforced by the controller.
	// When it is not enforced, then the user can select an implementation
	//
	if app.JobController.FeedbackReceiver == nil {
		available := ""
		if app.AvailableFeedbackReceivers != nil {
			for _, option := range app.AvailableFeedbackReceivers {
				available += option.GetImplementationName() + ", "
			}
		}
		command.Flags().StringVarP(&app.CustomFeedbackReceiver, "feedback-receiver", "f", "jxscm", "Sets a FeedbackReceiver (possible options: "+strings.TrimRight(available, ", ")+")")
	}
	if app.ConfigCollector == nil {
		available := ""
		if app.AvailableConfigCollectors != nil {
			for _, option := range app.AvailableConfigCollectors {
				available += option.GetImplementationName() + ", "
			}
		}
		command.Flags().StringVarP(&app.CustomConfigCollector, "config-provider", "c", "local", "Sets a ConfigCollector (possible options: "+strings.TrimRight(available, ", ")+" - possible to set multiple, comma separated, without spaces)")
	}

	command.Flags().BoolVarP(&app.Debug, "debug", "v", false, "Increase verbosity to the debug level")
	command.Flags().StringVarP(&app.MetricsBindAddress, "metrics-bind-address", "m", ":8080", "Metrics bind address")
	command.Flags().StringVarP(&app.HealthProbeBindAddress, "health-probe-bind-address", "p", ":8081", "Health probe bind address")
	command.Flags().BoolVarP(&app.LeaderElect, "leader-elect", "l", false, "Enable leader election")
	command.Flags().StringVarP(&app.LeaderElectId, "instance-id", "", "aSaMKO0", "Leader election ID (if running multiple controller instances with different configuration)")

	// error handling
	command.Flags().IntVarP(&app.DelayAfterErrorNum, "requeue-delay-after-error-count", "", 100, "Delay reconciliation of this resource, after it failed X times")
	command.Flags().IntVarP(&app.RequeueDelaySecs, "requeue-delay-secs", "", 15, "After (--requeue-delay-after-error-count) failed retries every reconciliation of this resource should be delayed by X seconds")
	command.Flags().IntVarP(&app.StopProcessingAfterErrorNum, "requeue-stop-after-error-count", "", 150, "Stop processing resource after X failed retries")

	return command
}
