package main

import (
    "fmt"
    // "path/filepath"
    "flag"
    // "os"
    "context"
    "k8s.io/client-go/rest"
    // "k8s.io/client-go/tools/clientcmd"
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

func main() {
    fmt.Println("podcleaner running")

    // var kubeconfig *string

    //TODO ensure kubeconfig can be retrieved in container

    // if home := os.Getenv("HOME"); home != "" {
    //     kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    // } else {
    //     kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    // }
    flag.Parse()

    config, err := rest.InClusterConfig(); if err != nil {
        fmt.Println(err)
    }

    // config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    // if err != nil {
    //     panic(err)
    // }

    clientset, err := kubernetes.NewForConfig(config)

    pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }

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
