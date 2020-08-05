package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/xlab/treeprint"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes"
)

// Service contains the service object information
type Service struct {
	Client      *kubernetes.Clientset
	Namespace   string
	Name        string
	ServiceType string
	Pods        string
	Ports       string
}

// NewService returns a new service struct
func NewService(client *kubernetes.Clientset, namespace string) *Service {
	return &Service{
		Client:    client,
		Namespace: namespace,
	}
}

// GetInformation gets the service information and the backend pods
func (s *Service) GetInformation(name string) (err error) {

	service, err := s.Client.CoreV1().Services(s.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	s.Name = service.Name
	s.ServiceType = ServiceTypeToString(service.Spec.Type)
	s.Ports = PortsToString(service.Spec.Type, service.Spec.Ports)
	s.Pods, err = PodsToString(s.Client, s.Namespace, service.Spec.Selector)
	if err != nil {
		return err
	}

	return
}

// PrintInformation prints the route information in table format
func (s *Service) PrintInformation() {

	table := &metav1.Table{
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Type", Type: "string"},
			{Name: "Port(s)", Type: "string"},
			{Name: "Pod(s)", Type: "string"},
		},
		Rows: []metav1.TableRow{
			{Cells: []interface{}{s.Name, s.ServiceType, s.Ports, s.Pods}},
		},
	}

	out := bytes.NewBuffer([]byte{})
	printer := printers.NewTablePrinter(printers.PrintOptions{})
	printer.PrintObj(table, out)

	fmt.Print(out.String())
}

// PrintGraph prints the route information in a tree graph format
func (s *Service) PrintGraph() {
	tree := treeprint.New()

	service := tree.AddMetaBranch("Service", s.Name)

	for _, pod := range strings.Split(s.Pods, ",") {
		service.AddMetaNode("Pod", pod)
	}

	fmt.Println(service.String())
}

// PodsToString gets the pods behind the service using the service selectors and
// it concatenates the Pods array into a string of pods separated by semicolons
func PodsToString(client *kubernetes.Clientset, namespace string, selector map[string]string) (podsString string, err error) {

	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(selector).String()})
	if err != nil {
		return "", err
	}

	var podsBuilder strings.Builder
	arrayLength := len(pods.Items)

	for index, pod := range pods.Items {
		podsBuilder.WriteString(pod.Name)

		if (index + 1) != arrayLength {
			podsBuilder.WriteString(",")
		}
	}

	podsString = podsBuilder.String()

	return
}

// PortsToString converts the servicePort array into a string
// of ports separated by semicolons
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

		if serviceType != v1.ServiceTypeClusterIP {
			portsBuilder.WriteString(" ")
			portsBuilder.WriteString(fmt.Sprint(port.NodePort))
		}

		if (index + 1) != arrayLength {
			portsBuilder.WriteString(",")
		}

	}

	return portsBuilder.String()
}

// ServiceTypeToString converts the service type into a string
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
