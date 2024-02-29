package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}

	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("error creating dynamic client: %v\n", err)
		os.Exit(1)
	}

	gvr := schema.GroupVersionResource{
		//kb api-versions
		Group:   "monitoring.coreos.com",
		Version: "v1",
		//kb api-resources
		Resource: "servicemonitors",
	}

	plans, err := dynamicClient.Resource(gvr).Namespace("kube-system").List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Printf("error getting %s: %v\n", gvr.Resource, err)
		os.Exit(1)
	}

	for _, pod := range plans.Items {
		fmt.Printf(
			"Name: %s\n",
			pod.Object["metadata"].(map[string]interface{})["name"],
		)
	}
}
