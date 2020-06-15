package cmd

import (
	"context"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	//"gopkg.in/oleiade/reflections.v1"
	"fmt"
	"os"
	"strings"
	"github.com/liggitt/tabwriter"
	// "encoding/json"
)

type Service struct {
	Client *kubernetes.Clientset
	Namespace string
	Name string 
	ServiceType v1.ServiceType
	Selector map[string]string 
	Pods []PodObj
	Ports []v1.ServicePort
}

type ServiceInterface interface {
	GetInformation(name string)
	PrintInformation()
}

func NewService(client *kubernetes.Clientset, namespace string) *Service {
	return &Service {
		Client: client,
		Namespace: namespace,
	}
}

// type ServiceInformation struct {
// 	Name string `json:"name"`
// 	ServiceType v1.ServiceType `json:"type"`
// 	Selector map[string]string `json:"selector"`
// 	Pods []PodObj `json:"pods"`
// 	Ports []v1.ServicePort `json:"ports"`
// }

func (s *Service) GetInformation(name string) (err error) {

	service, err := s.Client.CoreV1().Services(s.Namespace).Get(context.TODO(), name, metav1.GetOptions{})

	s.Name = service.Name
	s.Selector = service.Spec.Selector
	s.ServiceType = service.Spec.Type
	s.Ports = service.Spec.Ports

	pod := Pod{
		s.Client,
		s.Namespace,
	}

	s.Pods, err = pod.GetPods(s.Selector)

	return
}

func (s *Service) PrintInformation() {
	// initialize tabwriter
	w := new(tabwriter.Writer)
	
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	
	defer w.Flush()
	
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t", "NAMESPACE", "SERVICE", "TYPE", "PORTS Port TargetPort NodePort", "PODS Name Status")
	
	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t%s\t", s.Namespace, s.Name, s.ServiceType, PortsToString(s.ServiceType, s.Ports), PodsToString(s.Pods))
	//fmt.Fprintf(w, "\n%s\t%s\t%v\t%s\t%s\t", "", "", "", "└──", "")

	// for _, pod := range serviceInformation.Pods {
	// 	fmt.Fprintf(w, "\n%s\t%s\t%s\t %s\t%s\t", "", "", "", pod.Name, pod.Status)
	// }

	fmt.Fprintf(w, "\n")

}

func PodsToString(pods []PodObj) string{
	var podsString strings.Builder
	arrayLength := len(pods)

	for index, pod := range pods {
		// podsString.WriteString("{")
		podsString.WriteString(pod.Name)
		podsString.WriteString(" ")
		podsString.WriteString(fmt.Sprint(pod.Status))
		// podsString.WriteString("}")

		if ((index + 1) != arrayLength){
			podsString.WriteString(", ")
		}
	}

	return podsString.String()
}

func PortsToString(serviceType v1.ServiceType, ports []v1.ServicePort) string {

	var portsString strings.Builder
	arrayLength := len(ports)

	for index, item := range ports {

		// portsString.WriteString("{")
		portsString.WriteString(fmt.Sprint(item.Port))
		portsString.WriteString(" ")
		
		switch item.TargetPort.Type {
			case 0:
				portsString.WriteString(fmt.Sprint(item.TargetPort.IntVal))
			case 1:
				portsString.WriteString(item.TargetPort.StrVal)
		}

		if (serviceType != "ClusterIP" ){
			portsString.WriteString(" ")
			portsString.WriteString(fmt.Sprint(item.NodePort))
		}
		// portsString.WriteString("}")

		if ((index + 1) != arrayLength){
			portsString.WriteString(", ")
		}
	}

	return portsString.String()
}

// func CheckServiceType(serviceType v1.ServiceType) string {
// 	switch serviceType {
// 		case v1.ServiceTypeLoadBalancer:
// 			return "LoadBalancer"
// 		case v1.ServiceTypeNodePort:
// 			return "NodePort" 
// 		case v1.ServiceTypeClusterIP:
// 			return "ClusterIP"
// 	}

// 	return "Unknown"
// }