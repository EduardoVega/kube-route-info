package cmd

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apilabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Client defines Client atributes
type Client struct {
	Clientset *kubernetes.Clientset
	Namespace string
}

// NewClient returns a new Client struct
func NewClient(clientset *kubernetes.Clientset, namespace string) *Client {
	return &Client{
		Clientset: clientset,
		Namespace: namespace,
	}
}

// ClientInterface defines the Client functions
type ClientInterface interface {
	GetPodsByLabels(map[string]string) (*v1.PodList, error)
	GetServiceByName(string) (*v1.Service, error)
	GetIngressByName(name string) (*v1beta1.Ingress, error)
}

// GetPodsByLabels returns a list of pods that match the given labels
func (c *Client) GetPodsByLabels(labels map[string]string) (pods *v1.PodList, err error) {
	pods, err = c.Clientset.CoreV1().Pods(c.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: apilabels.Set(labels).String()})
	return
}

// GetServiceByName returns a service that matches a given name
func (c *Client) GetServiceByName(name string) (service *v1.Service, err error) {
	service, err = c.Clientset.CoreV1().Services(c.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return
}

// GetIngressByName returns an ingress that matches a given name
func (c *Client) GetIngressByName(name string) (ingress *v1beta1.Ingress, err error) {
	ingress, err = c.Clientset.NetworkingV1beta1().Ingresses(c.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return
}
