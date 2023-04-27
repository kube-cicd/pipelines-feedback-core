package app

import (
	pipelinesfeedbackv1alpha1scheme "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/client/clientset/versioned/scheme"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/controller"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/feedback"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/feedback/jxscm"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/logging"
	"github.com/pkg/errors"
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
)

type PipelinesFeedbackApp struct {
	// can read configuration from various sources
	ConfigCollector config.ConfigurationCollector

	JobController          *controller.GenericController
	ConfigController       *controller.ConfigurationController
	Debug                  bool
	MetricsBindAddress     string
	HealthProbeBindAddress string
	LeaderElect            bool

	CustomFeedbackReceiver string
	CustomConfigCollector  string

	// Feedback receivers available to choose by the user. Falls back to default, embedded list if not specified
	AvailableFeedbackReceivers []feedback.Receiver

	// Config providers available to choose by the user. Falls back to default, embedded list if not specified
	AvailableConfigCollectors []config.ConfigurationCollector

	Logger logging.Logger
}

func (app *PipelinesFeedbackApp) Run() error {
	app.Logger = logging.CreateLogger(app.Debug)

	if err := app.populateFeedbackReceiver(); err != nil {
		return err
	}
	if err := app.populateConfigCollector(); err != nil {
		return err
	}

	// add a standard scheme and Pipelines Feedback Core CRDs
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(pipelinesfeedbackv1alpha1scheme.AddToScheme(scheme))
	// todo: add custom scheme

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                        scheme,
		MetricsBindAddress:            app.MetricsBindAddress,
		Port:                          9443,
		HealthProbeBindAddress:        app.HealthProbeBindAddress,
		LeaderElection:                app.LeaderElect,
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
	if err := app.ConfigController.Initialize(kubeconfig, app.ConfigCollector, app.Logger, app.JobController.Store); err != nil {
		return errors.Wrap(err, "cannot push dependencies to ConfigurationController")
	}
	if err := app.JobController.InjectDependencies(recorder, kubeconfig, app.Logger, app.ConfigController.Provider); err != nil {
		return errors.Wrap(err, "cannot inject dependencies to GenericController")
	}

	// register controllers
	if err = app.JobController.SetupWithManager(mgr); err != nil {
		app.Logger.Error(err, "unable to setup job controller", "controller")
		return err
	}
	if err = app.ConfigController.SetupWithManager(mgr); err != nil {
		app.Logger.Error(err, "unable to setup configuration controller", "config")
		return err
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		app.Logger.Error(err, "unable to set up healthz")
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		app.Logger.Error(err, "unable to set up readyz")
		return err
	}

	app.Logger.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		app.Logger.Error(err, "cannot start manager")
		return err
	}
	return nil
}

func (app *PipelinesFeedbackApp) populateFeedbackReceiver() error {
	//
	// The mechanism allows to register multiple options and let the user to chose one option
	//
	if app.CustomFeedbackReceiver == "" {
		return nil
	}
	if app.AvailableFeedbackReceivers == nil {
		app.AvailableFeedbackReceivers = []feedback.Receiver{
			&jxscm.Receiver{},
		}
	}
	for _, pluggable := range app.AvailableFeedbackReceivers {
		if pluggable.CanHandle(app.CustomFeedbackReceiver) {
			app.JobController.FeedbackReceiver = pluggable
			return nil
		}
	}
	return errors.New("unrecognized FeedbackProvider")
}

func (app *PipelinesFeedbackApp) populateConfigCollector() error {
	// if the user did not select anything
	if app.CustomConfigCollector == "" {
		app.ConfigCollector = &config.LocalFileConfigurationCollector{}
		app.ConfigCollector.SetLogger(app.Logger)
		return nil
	}
	// if there are no available collectors
	if app.AvailableConfigCollectors == nil {
		app.AvailableConfigCollectors = []config.ConfigurationCollector{
			&config.LocalFileConfigurationCollector{},
		}
	}
	collectors := make([]config.ConfigurationCollector, 0)
	for _, pluggable := range app.AvailableConfigCollectors {
		if pluggable.CanHandle(app.CustomConfigCollector) {
			pluggable.SetLogger(app.Logger)
			collectors = append(collectors, pluggable)
		}
	}
	if len(collectors) == 0 {
		return errors.New("unrecognized ConfigProviders: " + app.CustomConfigCollector)
	}
	app.ConfigCollector = config.CreateMultipleCollector(collectors)
	app.ConfigCollector.SetLogger(app.Logger)
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
