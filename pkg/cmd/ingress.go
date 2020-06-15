package cmd

import (
	"context"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/api/core/v1"
	//"gopkg.in/oleiade/reflections.v1"
	"fmt"
	// "os"
	// "strings"
	// "github.com/liggitt/tabwriter"
	// "encoding/json"
)

type Ingress struct {
	Client *kubernetes.Clientset
	Namespace string
	Name string 
}

type IngressInterface interface {
	GetInformation(name string)
	PrintInformation()
}

func NewIngress(client *kubernetes.Clientset, namespace string) *Ingress {
	return &Ingress {
		Client: client,
		Namespace: namespace,
	}
}

func (i *Ingress) GetInformation(name string) (err error) {

	ingress, err := i.Client.NetworkingV1beta1().Ingresses(i.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }

	fmt.Println(ingress.Name)

	return
}

func (i *Ingress) PrintInformation() {

}

