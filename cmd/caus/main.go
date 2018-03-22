package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	elasticityclientset "github.com/klinakuf/caus/pkg/client/clientset/versioned"
	"github.com/klinakuf/caus/pkg/client/informers/externalversions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kuberconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master      = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	//TimeFrequency is the time frequency of processing elasticity definitions
	TimeFrequency = flag.String("timefreq", "10", "Frequency in seconds for processiong elasticity definitions e.g. if 10 every 10 sec")
	//ScaleDown is the duration between the last decision and the next scaling in decision
	ScaleDown = flag.String("scaledown", "3", "Duration in muinutes for the coming scale down decision")
)

//MyMonitor represents the data of the monitoring object
type MyMonitor struct {
	name string
}

//GetRate fetches the rate of the processing queue
func (m *MyMonitor) GetRate() (float64, error) {
	return float64(8), nil
}

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kuberconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}
	//elasticity client
	elasticityClient, err := elasticityclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building elasticity client: %v", err)
	}

	//kubernetes client
	kubernetesClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes client: %v", err)
	}
	myMonitorInstance := MyMonitor{}
	elSharedInformerFactory := externalversions.NewSharedInformerFactory(elasticityClient, 0)
	elInformer := elSharedInformerFactory.Caus().V1().Elasticities().Informer()
	elLister := elSharedInformerFactory.Caus().V1().Elasticities().Lister()
	controller := NewController(elInformer, elLister, kubernetesClient, &myMonitorInstance, elasticityClient)

	stop := make(chan struct{})
	defer close(stop)
	fmt.Println("Controller will run")
	elSharedInformerFactory.Start(stop)
	controller.Run(1, stop)
}
