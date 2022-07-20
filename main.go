package main

import (
    "fmt"
    "path/filepath"
    "flag"
    "os"
    "context"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/kubernetes"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
    fmt.Println("podcleaner running")

    var kubeconfig *string

    //TODO ensure kubeconfig can be retrieved in container

    if home := os.Getenv("HOME"); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }
    flag.Parse()

    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err)
    }

    clientset, err := kubernetes.NewForConfig(config)

    pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }

    for i := range pods.Items {
        podPhase := pods.Items[i].Status.Phase
        fmt.Println(pods.Items[i].ObjectMeta.Name)
        if podPhase == corev1.PodSucceeded {
            fmt.Printf("Deleting pod %s which finished with status %s\n", pods.Items[i].ObjectMeta.Name, podPhase)
            clientset.CoreV1().Pods("default").Delete(context.TODO(), pods.Items[i].ObjectMeta.Name, metav1.DeleteOptions{})
        }
    }

}
