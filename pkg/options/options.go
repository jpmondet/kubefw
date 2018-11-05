package options

import (
	"time"

	"github.com/spf13/pflag"
	//	apiv1 "k8s.io/api/core/v1"
)

//KubeFwConfig stores all the flags
type KubeFwConfig struct {
	CacheSyncTimeout   time.Duration
	CleanupConfig      bool
	EnablePprof        bool
	HealthPort         uint16
	HelpRequested      bool
	HostnameOverride   string
	IPTablesSyncPeriod time.Duration
	Kubeconfig         string
	Master             string
	MetricsEnabled     bool
	MetricsPath        string
	MetricsPort        uint16
	RunFirewall        bool
	Version            bool
	VLevel             string
	StandaloneServer   bool
}

//NewKubeFwConfig creates a new KubeFwConfig
func NewKubeFwConfig() *KubeFwConfig {
	return &KubeFwConfig{
		CacheSyncTimeout:   1 * time.Minute,
		IPTablesSyncPeriod: 5 * time.Minute,
	}
}

//AddFlags Gets flags from parameters (or defaults)
func (s *KubeFwConfig) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&s.HelpRequested, "help", "h", false,
		"Print usage information.")
	fs.BoolVarP(&s.Version, "version", "V", false,
		"Print version information.")
	fs.DurationVar(&s.CacheSyncTimeout, "cache-sync-timeout", s.CacheSyncTimeout,
		"The timeout for cache synchronization (e.g. '5s', '1m'). Must be greater than 0.")
	fs.BoolVar(&s.RunFirewall, "run-firewall", true,
		"Enables Network Policy -- sets up iptables to provide ingress firewall for pods.")
	fs.StringVar(&s.Master, "master", s.Master,
		"The address of the Kubernetes API server (overrides any value in kubeconfig).")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig,
		"Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	fs.BoolVar(&s.CleanupConfig, "cleanup-config", false,
		"Cleanup iptables rules, ipvs, ipset configuration and exit.")
	fs.BoolVar(&s.StandaloneServer, "standalone", false,
		"If you are running this plugin on a server that is external to the cluster. Can be used for tbshooting purposes.")
	fs.DurationVar(&s.IPTablesSyncPeriod, "iptables-sync-period", s.IPTablesSyncPeriod,
		"The delay between iptables rule synchronizations (e.g. '5s', '1m'). Must be greater than 0.")
	fs.StringVar(&s.HostnameOverride, "hostname-override", s.HostnameOverride,
		"Overrides the NodeName of the node. Set this if kubefw is unable to determine your NodeName automatically.")
	fs.BoolVar(&s.EnablePprof, "enable-pprof", false,
		"Enables pprof for debugging performance and memory leak issues.")
	fs.Uint16Var(&s.MetricsPort, "metrics-port", 0, "Prometheus metrics port, (Default 0, Disabled)")
	fs.StringVar(&s.MetricsPath, "metrics-path", "/metrics", "Prometheus metrics path")
	fs.StringVarP(&s.VLevel, "v", "v", "0", "log level for V logs")
	fs.Uint16Var(&s.HealthPort, "health-port", 20244, "Health check port, 0 = Disabled")
}
