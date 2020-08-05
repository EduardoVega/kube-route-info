package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/xlab/treeprint"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes"
)

// Ingress contains the ingress object information
type Ingress struct {
	Client    *kubernetes.Clientset
	Namespace string
	Name      string
	Rules     []RulesInfo
}

// RulesInfo contains the rules configured in the ingress
type RulesInfo struct {
	Host         string
	IngressRules []IngressRuleInfo
}

// IngressRuleInfo contains the information of each rule configured in the ingress
type IngressRuleInfo struct {
	Path        string
	ServiceName string
	ServicePort string
	ServiceInfo *Service
}

// NewIngress returns a new ingress struct
func NewIngress(client *kubernetes.Clientset, namespace string) *Ingress {
	return &Ingress{
		Client:    client,
		Namespace: namespace,
	}
}

// GetInformation gets the ingress information and its backend services
func (i *Ingress) GetInformation(name string) (err error) {

	ingress, err := i.Client.NetworkingV1beta1().Ingresses(i.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	i.Name = ingress.Name
	i.Rules = []RulesInfo{}

	for _, rule := range ingress.Spec.Rules {
		rulesInfo := RulesInfo{}

		rulesInfo.Host = rule.Host
		rulesInfo.IngressRules = []IngressRuleInfo{}

		for _, path := range rule.IngressRuleValue.HTTP.Paths {
			ingressRuleInfo := IngressRuleInfo{}
			ingressRuleInfo.Path = path.Path
			ingressRuleInfo.ServiceName = path.Backend.ServiceName
			ingressRuleInfo.ServiceInfo = GetServiceInfo(i.Client, i.Namespace, path.Backend.ServiceName)

			switch path.Backend.ServicePort.Type {
			case 0:
				ingressRuleInfo.ServicePort = fmt.Sprint(path.Backend.ServicePort.IntVal)
			case 1:
				ingressRuleInfo.ServicePort = path.Backend.ServicePort.StrVal
			}

			rulesInfo.IngressRules = append(rulesInfo.IngressRules, ingressRuleInfo)
		}

		i.Rules = append(i.Rules, rulesInfo)
	}

	return
}

// PrintInformation prints the route information in table format
func (i *Ingress) PrintInformation() {

	rows := []metav1.TableRow{}

	for _, rule := range i.Rules {
		for _, ingressRule := range rule.IngressRules {
			rows = append(rows, metav1.TableRow{
				Cells: []interface{}{i.Name, rule.Host, ingressRule.Path, ingressRule.ServicePort, ingressRule.ServiceInfo.Name, ingressRule.ServiceInfo.ServiceType, ingressRule.ServiceInfo.Ports, ingressRule.ServiceInfo.Pods},
			})
		}
	}

	table := &metav1.Table{
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Host", Type: "string"},
			{Name: "Path", Type: "string"},
			{Name: "Port", Type: "string"},
			{Name: "Service", Type: "string"},
			{Name: "Type", Type: "string"},
			{Name: "Port(s)", Type: "string"},
			{Name: "Pod(s)", Type: "string"},
		},
		Rows: rows,
	}

	out := bytes.NewBuffer([]byte{})
	printer := printers.NewTablePrinter(printers.PrintOptions{})
	printer.PrintObj(table, out)

	fmt.Print(out.String())
}

// PrintGraph prints the route information in a tree graph format
func (i *Ingress) PrintGraph() {
	tree := treeprint.New()

	ingress := tree.AddMetaBranch("Ingress", i.Name)

	for _, rule := range i.Rules {

		hostBranch := ingress.AddBranch(rule.Host)

		for _, ingressRule := range rule.IngressRules {
			ruleBranch := hostBranch.AddBranch(ingressRule.Path)
			serviceBranch := ruleBranch.AddMetaBranch("Service", ingressRule.ServiceInfo.Name)

			for _, pod := range strings.Split(ingressRule.ServiceInfo.Pods, ",") {
				serviceBranch.AddMetaNode("Pod", pod)
			}
		}
	}

	fmt.Println(ingress.String())
}

// GetServiceInfo gets the information of the services configured in the ingress
func GetServiceInfo(client *kubernetes.Clientset, namespace, serviceName string) (service *Service) {
	service = NewService(client, namespace)

	var err error
	err = service.GetInformation(serviceName)
	if err != nil {
		service.Name = "<none>"
		service.ServiceType = "<none>"
		service.Pods = "<none>"
		service.Ports = "<none>"
	}

	return
}
