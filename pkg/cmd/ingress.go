package cmd

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xlab/treeprint"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// Ingress defines Ingress atributes
type Ingress struct {
	Client    ClientInterface
	Namespace string
}

// NewIngress returns a new Ingress struct
func NewIngress(client ClientInterface, namespace string) *Ingress {
	return &Ingress{
		Client:    client,
		Namespace: namespace,
	}
}

// PrintGraph prints ingress route information in a tree graph format
func (i *Ingress) PrintGraph(name string, w io.Writer) (err error) {

	ingress, err := i.Client.GetIngressByName(name)
	if err != nil {
		return err
	}

	tree := treeprint.New()

	ingressBranch := tree.AddMetaBranch("Ingress", ingress.Name)

	for _, rule := range ingress.Spec.Rules {

		hostBranch := ingressBranch.AddBranch(rule.Host)

		for _, ingressRule := range rule.IngressRuleValue.HTTP.Paths {
			ruleBranch := hostBranch.AddBranch(ingressRule.Path)

			// Get service resource
			service, err := i.Client.GetServiceByName(ingressRule.Backend.ServiceName)
			if err != nil {
				ruleBranch.AddMetaBranch("Service", ingressRule.Backend.ServiceName+" *Not found*")

			} else {
				serviceBranch := ruleBranch.AddMetaBranch("Service", service.Name)

				// Check service type
				if ServiceTypeToString(service.Spec.Type) == "ExternalName" {
					serviceBranch.AddMetaNode("Hostname", service.Spec.ExternalName)

				} else {

					// Get pod resources
					pods, err := i.Client.GetPodsByLabels(service.Spec.Selector)
					if err != nil {
						return err
					}

					for _, pod := range pods.Items {
						serviceBranch.AddMetaNode("Pod", pod.Name)
					}
				}
			}
		}
	}

	fmt.Fprint(w, ingressBranch.String())

	return
}

// PrintTable prints ingress route information in table format
func (i *Ingress) PrintTable(name string, w io.Writer) (err error) {

	ingress, err := i.Client.GetIngressByName(name)
	if err != nil {
		return nil
	}

	rows := []metav1.TableRow{}

	// Default column name for pods
	podColumnName := "Pod(s)"

	for _, rule := range ingress.Spec.Rules {
		for _, ingressRule := range rule.IngressRuleValue.HTTP.Paths {

			// If service does not exist
			serviceName := ingressRule.Backend.ServiceName + " *Not found*"
			serviceType := ""
			servicePorts := ""
			servicePodsHostname := ""

			// Get service resource
			service, err := i.Client.GetServiceByName(ingressRule.Backend.ServiceName)

			// If service does exist
			if err == nil {
				serviceName = service.Name
				serviceType = ServiceTypeToString(service.Spec.Type)
				servicePorts = PortsToString(service.Spec.Ports)

				// Check service type
				if ServiceTypeToString(service.Spec.Type) == "ExternalName" {
					servicePodsHostname = service.Spec.ExternalName
					podColumnName = "Pod(s)/Hostname"

				} else {
					// Get pod resources
					pods, err := i.Client.GetPodsByLabels(service.Spec.Selector)
					if err != nil {
						return nil
					}

					servicePodsHostname = PodsToString(pods)
				}
			}

			rows = append(rows, metav1.TableRow{
				Cells: []interface{}{
					ingress.Name,
					rule.Host,
					ingressRule.Path,
					PortToString(
						ingressRule.Backend.ServicePort.Type,
						ingressRule.Backend.ServicePort.StrVal,
						ingressRule.Backend.ServicePort.IntVal,
					),
					serviceName,
					serviceType,
					servicePorts,
					servicePodsHostname,
				},
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
			{Name: "Service Port(s)", Type: "string"},
			{Name: podColumnName, Type: "string"},
		},
		Rows: rows,
	}

	out := bytes.NewBuffer([]byte{})
	printer := printers.NewTablePrinter(printers.PrintOptions{})
	printer.PrintObj(table, out)

	fmt.Fprint(w, out.String())

	return
}
