package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type crtbResponse struct {
	ApiVersion string
	Items      []crtb
}

type crtb struct {
	Metadata metav1.ObjectMeta
	UserName string
}

func main() {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d cluster roles in the cluster\n", len(clusterRoles.Items))

	c := crtbResponse{}
	crtbs, err := clientset.RESTClient().Get().AbsPath("/apis/management.cattle.io/v3").Resource("clusterroletemplatebindings").DoRaw(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(crtbs, &c)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(c)
}
