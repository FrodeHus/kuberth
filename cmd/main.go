package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	kuberth "github.com/frodehus/kuberth/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	dnsClient, err := kuberth.NewForConfig(config)
	kuberthAPI := dnsClient.KuberthV1alpha1()
	listOptions := metav1.ListOptions{}
	entries, err := kuberthAPI.DnsEntries("default").List(listOptions)

	for _, entry := range entries.Items {
		fmt.Printf("%s\n", entry.Name)
		template := "%-32s%-8s%-8s\n"
		for _, spec := range entry.Spec {
			fmt.Printf(template, spec.Name, spec.Type, spec.Value)
		}
	}
}
