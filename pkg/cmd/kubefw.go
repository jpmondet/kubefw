package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/golang/glog"
	"github.com/jpmondet/kubefw/pkg/controllers/netpol"
	"github.com/jpmondet/kubefw/pkg/healthcheck"
	"github.com/jpmondet/kubefw/pkg/metrics"
	"github.com/jpmondet/kubefw/pkg/options"
	"github.com/jpmondet/kubefw/pkg/utils"

	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// These get set at build time via -ldflags magic
var version string
var buildDate string

// KubeFw holds the information needed to run server
type KubeFw struct {
	Client kubernetes.Interface
	Config *options.KubeFwConfig
}

// NewKubeFwDefault returns a KubeFw object
func NewKubeFwDefault(config *options.KubeFwConfig) (*KubeFw, error) {

	var clientconfig *rest.Config
	var err error
	PrintVersion(true)
	// Use out of cluster config if the URL or kubeconfig have been specified. Otherwise use incluster config.
	if len(config.Master) != 0 || len(config.Kubeconfig) != 0 {
		clientconfig, err = clientcmd.BuildConfigFromFlags(config.Master, config.Kubeconfig)
		if err != nil {
			return nil, errors.New("Failed to build configuration from CLI: " + err.Error())
		}
	} else {
		clientconfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, errors.New("unable to initialize inclusterconfig: " + err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(clientconfig)
	if err != nil {
		return nil, errors.New("Failed to create Kubernetes client: " + err.Error())
	}

	return &KubeFw{Client: clientset, Config: config}, nil
}

// CleanupConfigAndExit performs Cleanup on all three controllers
func CleanupConfigAndExit() {
	npc := netpol.NetworkPolicyController{}
	npc.Cleanup()
}

// Run starts the controllers and waits forever till we get SIGINT or SIGTERM
func (kfw *KubeFw) Run() error {
	var err error
	var wg sync.WaitGroup
	healthChan := make(chan *healthcheck.ControllerHeartbeat, 10)
	defer close(healthChan)
	stopCh := make(chan struct{})

	if !(kfw.Config.RunFirewall) {
		glog.Info("Firewall functionality must be specified. Exiting!")
		os.Exit(0)
	}

	hc, err := healthcheck.NewHealthController(kfw.Config)
	if err != nil {
		return errors.New("Failed to create health controller: " + err.Error())
	}
	wg.Add(1)
	go hc.RunServer(stopCh, &wg)

	informerFactory := informers.NewSharedInformerFactory(kfw.Client, 0)
	podInformer := informerFactory.Core().V1().Pods().Informer()
	nsInformer := informerFactory.Core().V1().Namespaces().Informer()
	npInformer := informerFactory.Networking().V1().NetworkPolicies().Informer()
	informerFactory.Start(stopCh)

	err = kfw.CacheSyncOrTimeout(informerFactory, stopCh)
	if err != nil {
		return errors.New("Failed to synchronize cache: " + err.Error())
	}

	hc.SetAlive()
	wg.Add(1)
	go hc.RunCheck(healthChan, stopCh, &wg)

	if (kfw.Config.MetricsPort > 0) && (kfw.Config.MetricsPort <= 65535) {
		kfw.Config.MetricsEnabled = true
		mc, err := metrics.NewMetricsController(kfw.Client, kfw.Config)
		if err != nil {
			return errors.New("Failed to create metrics controller: " + err.Error())
		}
		wg.Add(1)
		go mc.Run(healthChan, stopCh, &wg)

	} else if kfw.Config.MetricsPort > 65535 {
		glog.Errorf("Metrics port must be over 0 and under 65535, given port: %d", kfw.Config.MetricsPort)
		kfw.Config.MetricsEnabled = false
	} else {
		kfw.Config.MetricsEnabled = false
	}

	if kfw.Config.RunFirewall {
		if kfw.Config.StandaloneServer {
			nodes, err := utils.AllNodesObjects(kfw.Client)
			if err != nil {
				return errors.New("Failed to retrieve cluster nodes objects")
			}

			for _, node := range nodes.Items {
				kfw.Config.HostnameOverride = node.Name
				npc, err := netpol.NewNetworkPolicyController(kfw.Client,
					kfw.Config, podInformer, npInformer, nsInformer)
				if err != nil {
					return errors.New("Failed to create network policy controller: " + err.Error())
				}

				podInformer.AddEventHandler(npc.PodEventHandler)
				nsInformer.AddEventHandler(npc.NamespaceEventHandler)
				npInformer.AddEventHandler(npc.NetworkPolicyEventHandler)

				wg.Add(1)
				go npc.Run(healthChan, stopCh, &wg)
			}
		} else {
			npc, err := netpol.NewNetworkPolicyController(kfw.Client,
				kfw.Config, podInformer, npInformer, nsInformer)
			if err != nil {
				return errors.New("Failed to create network policy controller: " + err.Error())
			}

			podInformer.AddEventHandler(npc.PodEventHandler)
			nsInformer.AddEventHandler(npc.NamespaceEventHandler)
			npInformer.AddEventHandler(npc.NetworkPolicyEventHandler)

			wg.Add(1)
			go npc.Run(healthChan, stopCh, &wg)
		}
	}

	// Handle SIGINT and SIGTERM
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	glog.Infof("Shutting down the controllers")
	close(stopCh)

	wg.Wait()
	return nil
}

// CacheSyncOrTimeout performs cache synchronization under timeout limit
func (kfw *KubeFw) CacheSyncOrTimeout(informerFactory informers.SharedInformerFactory, stopCh <-chan struct{}) error {
	syncOverCh := make(chan struct{})
	go func() {
		informerFactory.WaitForCacheSync(stopCh)
		close(syncOverCh)
	}()

	select {
	case <-time.After(kfw.Config.CacheSyncTimeout):
		return errors.New(kfw.Config.CacheSyncTimeout.String() + " timeout")
	case <-syncOverCh:
		return nil
	}
}

// PrintVersion print version and exit
func PrintVersion(logOutput bool) {
	output := fmt.Sprintf("Running %v version %s, built on %s, %s\n", os.Args[0], version, buildDate, runtime.Version())

	if !logOutput {
		fmt.Fprintf(os.Stderr, output)
	} else {
		glog.Info(output)
	}
}
