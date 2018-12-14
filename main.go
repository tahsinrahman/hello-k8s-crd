package main

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	myclient "github.com/tahsinrahman/hello-k8s-crd/pkg/client/clientset/versioned"
)

func main() {
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatal(err)
	}

	myClientset, err := myclient.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	k8sClientSet, err := kubernetes.NewConfig(config)
	if err != nil {
		log.Fatal(err)
	}
}
