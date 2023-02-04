package cli

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/controller"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/feedback"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

var (
	scheme = runtime.NewScheme()
	appLog = ctrl.Log.WithName("controller")
)

func NewRootCommand(app *PipelinesFeedbackApp) *cobra.Command {
	command := &cobra.Command{
		Use:   "pipelines-feedback",
		Short: "Runs a Kubernetes controller that observes Pipeline status and reports the status to the external system",
		RunE: func(command *cobra.Command, args []string) error {
			return app.Run()
		},
	}

	app.customFeedbackReceiver = ""
	app.customConfigProvider = ""

	// FeedbackReceiver and ConfigProvider could be enforced by the controller. When it is not enforced then the user can select an implementation
	if app.Controller.FeedbackReceiver == nil {
		command.Flags().StringVarP(&app.customFeedbackReceiver, "feedback-receiver", "f", "jxscm", "Sets a FeedbackReceiver (possible options: jxscm)")
	}
	if app.Controller.ConfigProvider == nil {
		command.Flags().StringVarP(&app.customConfigProvider, "config-provider", "c", "local", "Sets a ConfigProvider (possible options: local)")
	}

	command.Flags().BoolVarP(&app.debug, "debug", "v", true, "Increase verbosity to the debug level")
	command.Flags().StringVarP(&app.metricsBindAddress, "metrics-bind-address", "m", ":8080", "Metrics bind address")
	command.Flags().StringVarP(&app.healthProbeBindAddress, "health-probe-bind-address", "p", ":8081", "Health probe bind address")
	command.Flags().BoolVarP(&app.leaderElect, "leader-elect", "l", false, "Enable leader election")

	return command
}

type PipelinesFeedbackApp struct {
	Controller             *controller.GenericController
	debug                  bool
	metricsBindAddress     string
	healthProbeBindAddress string
	leaderElect            bool

	customFeedbackReceiver string
	customConfigProvider   string
}

func (app *PipelinesFeedbackApp) populateFeedbackReceiver() error {
	if app.customFeedbackReceiver == "" {
		return nil
	}
	if app.customFeedbackReceiver == "jxscm" {
		app.Controller.FeedbackReceiver = &feedback.JXSCMReceiver{}
		return nil
	}
	return errors.New("unrecognized FeedbackProvider")
}

func (app *PipelinesFeedbackApp) populateConfigProvider() error {
	if app.customConfigProvider == "" {
		return nil
	}
	if app.customConfigProvider == "local" {
		app.Controller.ConfigProvider = &config.LocalFileConfigurationProvider{}
		return nil
	}
	return errors.New("unrecognized ConfigProvider")
}

func (app *PipelinesFeedbackApp) Run() error {
	if app.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := app.populateFeedbackReceiver(); err != nil {
		return err
	}
	if err := app.populateConfigProvider(); err != nil {
		return err
	}

	// add a standard scheme
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	// todo: add custom scheme

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                        scheme,
		MetricsBindAddress:            app.metricsBindAddress,
		Port:                          9443,
		HealthProbeBindAddress:        app.healthProbeBindAddress,
		LeaderElection:                app.leaderElect,
		LeaderElectionID:              "aSaMKO0.keskad.pl",
		LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		return err
	}

	recorder := mgr.GetEventRecorderFor("pipelines-feedback")
	kubeconfig, err := createKubeConfiguration(os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err.Error())
	}

	// dependencies
	if err := app.Controller.InjectDependencies(recorder, kubeconfig); err != nil {
		return errors.Wrap(err, "cannot inject dependencies")
	}

	if err = app.Controller.SetupWithManager(mgr); err != nil {
		appLog.Error(err, "unable to setup controller", "controller")
		return err
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		appLog.Error(err, "unable to set up healthz")
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		appLog.Error(err, "unable to set up readyz")
		return err
	}

	appLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		appLog.Error(err, "cannot start manager")
		return err
	}
	return nil
}

func createKubeConfiguration(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		fromFlags, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return fromFlags, nil
	}
	inCluster, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return inCluster, nil
}
