package main

import (
    "fmt"
    "os"
    "context"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/kubernetes"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    policy "k8s.io/api/policy/v1beta1"
)

// Testing purposes
func evictNginx(client *kubernetes.Clientset) error {
    return client.PolicyV1beta1().Evictions("default").Evict(context.TODO(), &policy.Eviction{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "nginx",
            Namespace: "default",
        },
    })
}

func buildFromKubeConfig() *rest.Config {
    home := os.Getenv("HOME")
    config, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", home)); if err != nil {
        fmt.Println(err)
    }
    return config
}

func main() {
    fmt.Println("podcleaner running")

    // Fetch kube config
    config, err := rest.InClusterConfig(); if err != nil {
        if err.Error() == "unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined" {
            config = buildFromKubeConfig()
        } else {
            fmt.Println(err)
        }
    }

    // Create client
    clientset, err := kubernetes.NewForConfig(config)

    // List pods
    pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }
    // Iterate through pods
    for i := range pods.Items {
        podReason := pods.Items[i].Status.Reason
        if podReason == "Evicted" {
            fmt.Printf("Deleting pod %s which finished with status: %s\n", pods.Items[i].ObjectMeta.Name, podReason)
            err = clientset.CoreV1().Pods("default").Delete(context.TODO(), pods.Items[i].ObjectMeta.Name, metav1.DeleteOptions{}); if err != nil {
                fmt.Println(err)
            }
        }
    }
}
