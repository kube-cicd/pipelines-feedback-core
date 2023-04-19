package cli

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/app"
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
	app.CustomConfigProvider = ""

	//
	// FeedbackReceiver and ConfigProvider can be enforced by the controller.
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
	if app.JobController.ConfigProvider == nil {
		available := ""
		if app.AvailableConfigProviders != nil {
			for _, option := range app.AvailableConfigProviders {
				available += option.GetImplementationName() + ", "
			}
		}
		command.Flags().StringVarP(&app.CustomConfigProvider, "config-provider", "c", "local", "Sets a ConfigProvider (possible options: "+strings.TrimRight(available, ", ")+" - possible to set multiple, comma separated, without spaces)")
	}

	command.Flags().BoolVarP(&app.Debug, "debug", "v", false, "Increase verbosity to the debug level")
	command.Flags().StringVarP(&app.MetricsBindAddress, "metrics-bind-address", "m", ":8080", "Metrics bind address")
	command.Flags().StringVarP(&app.HealthProbeBindAddress, "health-probe-bind-address", "p", ":8081", "Health probe bind address")
	command.Flags().BoolVarP(&app.LeaderElect, "leader-elect", "l", false, "Enable leader election")

	return command
}
