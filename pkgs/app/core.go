package app

import (
	"os"
	"strings"

	pipelinesfeedbackv1alpha1scheme "github.com/kube-cicd/pipelines-feedback-core/pkgs/client/clientset/versioned/scheme"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/controller"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/feedback"
	debugFeedback "github.com/kube-cicd/pipelines-feedback-core/pkgs/feedback/debug"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/feedback/jxscm"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/store"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	scheme = runtime.NewScheme()
)

type SchemeSetter func(*runtime.Scheme) error

type PipelinesFeedbackApp struct {
	// can read configuration from various sources
	ConfigCollector config.ConfigurationCollector

	JobController          *controller.GenericController
	ConfigController       *controller.ConfigurationController
	Debug                  bool
	DisableCRD             bool
	ControllerName         string
	MetricsBindAddress     string
	HealthProbeBindAddress string
	LeaderElect            bool
	LeaderElectId          string
	RestrictNamespaces     string

	CustomFeedbackReceiver string
	CustomStore            string
	CustomConfigCollector  string

	// error handling
	DelayAfterErrorNum          int
	RequeueDelaySecs            int
	StopProcessingAfterErrorNum int

	// Feedback receivers available to choose by the user. Falls back to default, embedded list if not specified
	AvailableFeedbackReceivers []feedback.Receiver

	// Stores available to pick by the user
	AvailableStores []store.Store

	// Config providers available to choose by the user. Falls back to default, embedded list if not specified
	AvailableConfigCollectors []config.ConfigurationCollector

	// Allows to register custom CRD schema for the controller
	KubernetesSchemeSetters []SchemeSetter

	Logger *logging.InternalLogger
	schema *config.SchemaValidator
}

func (app *PipelinesFeedbackApp) Run() error {
	app.Logger = logging.CreateLogger(app.Debug)

	if err := app.populateFeedbackReceiver(); err != nil {
		return err
	}
	if err := app.populateConfigCollector(); err != nil {
		return err
	}
	if err := app.populateStoreAdapter(); err != nil {
		return err
	}

	// add a standard scheme and Pipelines Feedback Core CRDs
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	if !app.DisableCRD {
		utilruntime.Must(pipelinesfeedbackv1alpha1scheme.AddToScheme(scheme))
	}
	// custom CRD schemes
	for _, schemeSetter := range app.KubernetesSchemeSetters {
		utilruntime.Must(schemeSetter(scheme))
	}

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	managerOpts := ctrl.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: app.MetricsBindAddress,
		},
		HealthProbeBindAddress:        app.HealthProbeBindAddress,
		LeaderElection:                app.LeaderElect,
		LeaderElectionID:              app.LeaderElectId + app.ControllerName + ".keskad.pl",
		LeaderElectionReleaseOnCancel: true,
	}

	// restrict controller to specific namespaces only
	if app.RestrictNamespaces != "" {
		nsList := strings.Split(app.RestrictNamespaces, ",")
		nsMap := make(map[string]cache.Config, len(nsList))
		leaderElectionNamespace := ""
		for i, ns := range nsList {
			ns = strings.TrimSpace(ns)
			if i == 0 {
				leaderElectionNamespace = ns
			}
			// nil will have the namespace use the default settings
			nsMap[ns] = cache.Config{}
		}
		managerOpts.Cache.DefaultNamespaces = nsMap
		managerOpts.LeaderElectionNamespace = leaderElectionNamespace
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), managerOpts)
	if err != nil {
		return err
	}

	recorder := mgr.GetEventRecorderFor(app.ControllerName)
	kubeconfig, err := createKubeConfiguration(os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err.Error())
	}

	// PFConfig schema registration: There all dynamically loaded components are able to register their schema
	// for configuration keys validation
	app.schema = &config.SchemaValidator{Debug: app.Debug}
	app.schema.Add(config.Schema{
		Name: "global",
		AllowedFields: []string{
			"dashboard-url",
			"logs-enabled",
			"logs-max-line-length",
			"logs-max-full-length-lines-count",
			"logs-split-separator",
		},
	})

	// dependencies
	if err := app.ConfigController.Initialize(kubeconfig, app.ConfigCollector, app.Logger, app.JobController.Store, app.schema); err != nil {
		return errors.Wrap(err, "cannot push dependencies to ConfigurationController")
	}
	if err := app.JobController.InjectDependencies(recorder, kubeconfig, app.Logger,
		app.ConfigController.Provider, app.schema); err != nil {

		return errors.Wrap(err, "cannot inject dependencies to GenericController")
	}

	// collect configuration initially right after all components are injected (and registered in ConfigurationProvider)
	if err := app.ConfigController.CollectInitially(app.ConfigCollector); err != nil {
		return errors.Wrap(err, "cannot initially collect configuration")
	}

	// register controllers
	if err = app.JobController.SetupWithManager(mgr); err != nil {
		app.Logger.Error(err, "unable to setup job controller", "controller")
		return err
	}
	if !app.DisableCRD {
		if err = app.ConfigController.SetupWithManager(mgr); err != nil {
			app.Logger.Error(err, "unable to setup configuration controller", "config")
			return err
		}
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
			&debugFeedback.Receiver{},
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
		app.ConfigCollector = config.NewLocalFileConfigurationCollector(app.Logger, "")
		return nil
	}
	// if there are no available collectors
	if app.AvailableConfigCollectors == nil {
		app.AvailableConfigCollectors = []config.ConfigurationCollector{
			config.NewLocalFileConfigurationCollector(app.Logger, ""),
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
	app.ConfigCollector = config.CreateMultipleCollector(collectors, app.Logger)
	app.ConfigCollector.SetLogger(app.Logger)
	return nil
}

func (app *PipelinesFeedbackApp) populateStoreAdapter() error {
	if app.CustomStore == "" {
		return nil
	}
	if app.AvailableStores == nil {
		app.AvailableStores = []store.Store{
			store.NewMemory(),
			store.NewRedis(),
		}
	}
	for _, pluggable := range app.AvailableStores {
		if pluggable.CanHandle(app.CustomStore) {
			app.JobController.Store = store.Operator{Store: pluggable}
			app.Logger.Infof("Initializing store adapter '%s'", pluggable.GetImplementationName())
			return app.JobController.Store.Initialize()
		}
	}
	return errors.New("unrecognized store type")
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
