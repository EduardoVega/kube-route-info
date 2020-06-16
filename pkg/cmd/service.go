package cmd

import (
	"context"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	"fmt"
	"os"
	"strings"
	"github.com/liggitt/tabwriter"
	"k8s.io/apimachinery/pkg/labels"
)

type Service struct {
	Client *kubernetes.Clientset
	Namespace string
	Name string 
	ServiceType string
	Pods  string
	Ports string
}

// type ServiceInterface interface {
// 	GetInformation(name string)
// 	PrintInformation()
// }

func NewService(client *kubernetes.Clientset, namespace string) *Service {
	return &Service {
		Client: client,
		Namespace: namespace,
	}
}

func (s *Service) GetInformation(name string) (err error) {

	service, err := s.Client.CoreV1().Services(s.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if (err != nil){
		return err
	}

	s.Name = service.Name
	s.ServiceType = ServiceTypeToString(service.Spec.Type)
	s.Ports = PortsToString(service.Spec.Type, service.Spec.Ports)
	s.Pods, err = PodsToString(s.Client, s.Namespace, service.Spec.Selector)
	if (err != nil){
		return err
	}

	return
}

func (s *Service) PrintInformation() {
	// initialize tabwriter
	w := new(tabwriter.Writer)
	
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	
	defer w.Flush()
	
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t", "NAME", "TYPE", "PORT(S)", "POD(S)")
	
	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t", s.Name, s.ServiceType, s.Ports, s.Pods)

	fmt.Fprintf(w, "\n")

}

func PodsToString(client *kubernetes.Clientset, namespace string, selector map[string]string) (podsString string, err error){

	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(selector).String() })
	if (err != nil){
		return "", err
	}

	var podsBuilder strings.Builder
	arrayLength := len(pods.Items)

	for index, pod := range pods.Items {
		podsBuilder.WriteString(pod.Name)	

		if ((index + 1) != arrayLength){
			podsBuilder.WriteString(",")
		}
	}

	podsString = podsBuilder.String()

	return
}

func PortsToString(serviceType v1.ServiceType, ports []v1.ServicePort) string {

	var portsBuilder strings.Builder
	arrayLength := len(ports)

	for index, port := range ports {
		portsBuilder.WriteString(fmt.Sprint(port.Port))
		portsBuilder.WriteString(" ")

		switch port.TargetPort.Type {
		case 0:
			portsBuilder.WriteString(fmt.Sprint(port.TargetPort.IntVal))
		case 1:
			portsBuilder.WriteString(port.TargetPort.StrVal)
		} 

		if (serviceType != v1.ServiceTypeClusterIP ){
			portsBuilder.WriteString(" ")
			portsBuilder.WriteString(fmt.Sprint(port.NodePort))
		}

		if ((index + 1) != arrayLength){
			portsBuilder.WriteString(",")
		}

	}

	return portsBuilder.String()
}

func ServiceTypeToString(serviceType v1.ServiceType) string {
	switch serviceType {
		case v1.ServiceTypeLoadBalancer:
			return "LoadBalancer"
		case v1.ServiceTypeNodePort:
			return "NodePort" 
		case v1.ServiceTypeClusterIP:
			return "ClusterIP"
	}

	return "Unknown"
}