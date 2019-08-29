package main

import (
	"flag"
	"os"
	"time"

	"k8s.io/client-go/util/homedir"

	"github.com/homedepot/k8s-global-objects/runner"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig  string
	runInterval time.Duration
	runOnce     bool
	debug       bool
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", homedir.HomeDir()+"/.kube/config", "KUBECONFIG location")
	flag.DurationVar(&runInterval, "runinterval", time.Second*60, "interval to kick off sync")
	flag.BoolVar(&runOnce, "runonce", false, "Run App once")
	flag.BoolVar(&debug, "debug", false, "Debug")
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	// for fluentd format
	//log.SetFormatter(&joonix.FluentdFormatter{})

	if debug {
		log.SetLevel(log.DebugLevel)
		//log.SetReportCaller(true)
	}

	log.Debugf("Flag kubeconfig: %v", kubeconfig)
	log.Debugf("Flag runinterval: %v", runInterval)
	log.Debugf("Flag runOnce: %v", runOnce)
	log.Debugf("Flag debug: %v", debug)
}

func main() {
	// attempting to see if running in kubernetes
	config, err := rest.InClusterConfig()
	if err != nil {
		log.WithError(err)
	}

	// if not running in container will use kubeconfig flag value
	if config == nil {
		log.Debug("No InClusterConfig found, using kubeconfig flag instead")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err)
		}
	}

	// creating clientset
	client := &runner.K8S{}
	client.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.WithError(err)
	}

	// start runner
	var run *runner.Runner
	{
		runnerConfig := &runner.Config{
			Client:      client,
			RunInterval: runInterval,
			Debug:       debug,
			Once:        runOnce,
		}

		log.Info("Starting K8S Global Objects Runner")
		run = runner.NewRunner(runnerConfig)

		err = run.Init()
		if err != nil {
			log.Fatal(err)
		}

		err := run.Start()
		if err != nil {
			log.Fatal(err)
		}
	}

	defer run.Close()
}
