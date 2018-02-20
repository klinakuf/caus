package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	elasticityclientset "github.com/klinakuf/caus/pkg/client/clientset/versioned"
)

var (
	kuberconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master      = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
)

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kuberconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}

	elasticityClient, err := elasticityclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %v", err)
	}

	list, err := elasticityClient.CausV1().Elasticities("default").List(metav1.ListOptions{})
	//Databases("default").List(metav1.ListOptions{})
	if err != nil {
		glog.Fatalf("Error listing all databases: %v", err)
	}

	for _, el := range list.Items {
		fmt.Printf("Elasticity rules for deployment %s with threshold %d for the workload %s\n", el.Spec.Deployment.Name, el.Spec.Buffer.Threshold, el.Spec.Workload.Queue)
	}
}
