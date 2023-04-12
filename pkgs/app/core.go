package app

import (
	pipelinesfeedbackv1alpha1scheme "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/client/clientset/versioned/scheme"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/controller"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/feedback"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

type PipelinesFeedbackApp struct {
	JobController          *controller.GenericController
	ConfigController       *controller.ConfigurationController
	Debug                  bool
	MetricsBindAddress     string
	HealthProbeBindAddress string
	LeaderElect            bool

	CustomFeedbackReceiver string
	CustomConfigProvider   string

	// Feedback receivers available to choose by the user. Falls back to default, embedded list if not specified
	AvailableFeedbackReceivers []feedback.Receiver

	// Config providers available to choose by the user. Falls back to default, embedded list if not specified
	AvailableConfigProviders []config.ConfigurationProvider
}

func (app *PipelinesFeedbackApp) Run() error {
	logrus.SetLevel(logrus.InfoLevel)
	if app.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := app.populateFeedbackReceiver(); err != nil {
		return err
	}
	if err := app.populateConfigProvider(); err != nil {
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
	if err := app.JobController.InjectDependencies(recorder, kubeconfig); err != nil {
		return errors.Wrap(err, "cannot inject dependencies")
	}

	if err = app.JobController.SetupWithManager(mgr); err != nil {
		appLog.Error(err, "unable to setup job controller", "controller")
		return err
	}

	if err = app.ConfigController.SetupWithManager(mgr); err != nil {
		appLog.Error(err, "unable to setup configuration controller", "config")
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

func (app *PipelinesFeedbackApp) populateFeedbackReceiver() error {
	//
	// The mechanism allows to register multiple options and let the user to chose one option
	//

	if app.CustomFeedbackReceiver == "" {
		return nil
	}
	if app.AvailableFeedbackReceivers == nil {
		app.AvailableFeedbackReceivers = []feedback.Receiver{
			&feedback.JXSCMReceiver{},
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

func (app *PipelinesFeedbackApp) populateConfigProvider() error {
	if app.CustomConfigProvider == "" {
		return nil
	}
	if app.AvailableConfigProviders == nil {
		app.AvailableConfigProviders = []config.ConfigurationProvider{
			&config.LocalFileConfigurationProvider{},
		}
	}
	for _, pluggable := range app.AvailableConfigProviders {
		if pluggable.CanHandle(app.CustomConfigProvider) {
			app.JobController.ConfigProvider = pluggable
			return nil
		}
	}
	return errors.New("unrecognized ConfigProvider: " + app.CustomConfigProvider)
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
