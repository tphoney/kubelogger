package main

import (
	"context"
	"fmt"
	"io"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go kubeConfigPath namespace podName containerName")
		return
	}

	kubeConfigPath := os.Args[1]
	namespace := os.Args[2]
	podName := os.Args[3]
	containerName := os.Args[4]

	fmt.Printf("kube config path: %q, namespace: %q, pod: %q, container %q\n", kubeConfigPath, namespace, podName, containerName)

	// path-to-kubeconfig -- for example, /root/.kube/config
	config, _ := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	// creates the clientset
	clientset, _ := kubernetes.NewForConfig(config)
	// access the API to list pods
	pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	// print the pod names
	for _, pod := range pods.Items {
		fmt.Printf("Pod name %s\n", pod.GetName())
		// print the container names
		for _, container := range pod.Spec.Containers {
			fmt.Printf("Container name %s\n", container.Name)
		}
	}

	opts := &v1.PodLogOptions{
		Follow:    true,
		Container: containerName,
	}

	ctx := context.Background()
	req := clientset.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(podName).
		Resource("pods").
		SubResource("log").
		VersionedParams(opts, scheme.ParameterCodec)

	readCloser, err := req.Stream(ctx)

	if err != nil {
		fmt.Println("failed to stream logs", err)
	}
	// print the logs to stdout
	io.Copy(os.Stdout, readCloser)
	defer readCloser.Close()
}
