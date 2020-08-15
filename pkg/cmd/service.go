package cmd

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xlab/treeprint"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// Service defines Service attributes
type Service struct {
	Client    ClientInterface
	Namespace string
}

// NewService returns a new Service struct
func NewService(client ClientInterface, namespace string) *Service {
	return &Service{
		Client:    client,
		Namespace: namespace,
	}
}

// PrintGraph prints service route information in a tree graph format
func (s *Service) PrintGraph(name string, w io.Writer) (err error) {

	service, err := s.Client.GetServiceByName(name)
	if err != nil {
		return err
	}

	pods, err := s.Client.GetPodsByLabels(service.Spec.Selector)
	if err != nil {
		return err
	}

	tree := treeprint.New()

	serviceBranch := tree.AddMetaBranch("Service", service.Name)

	if ServiceTypeToString(service.Spec.Type) == "ExternalName" {
		serviceBranch.AddMetaNode("Hostname", service.Spec.ExternalName)

	} else {
		for _, pod := range pods.Items {
			serviceBranch.AddMetaNode("Pod", pod.Name)
		}
	}

	fmt.Fprint(w, serviceBranch.String())

	return
}

// PrintTable prints service route information in table format
func (s *Service) PrintTable(name string, w io.Writer) (err error) {

	service, err := s.Client.GetServiceByName(name)
	if err != nil {
		return err
	}

	pods, err := s.Client.GetPodsByLabels(service.Spec.Selector)
	if err != nil {
		return err
	}

	columnName := ""
	cellValue := ""

	if ServiceTypeToString(service.Spec.Type) == "ExternalName" {
		columnName = "Hostname"
		cellValue = service.Spec.ExternalName

	} else {
		columnName = "Pod(s)"
		cellValue = PodsToString(pods)
	}

	table := &metav1.Table{
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Type", Type: "string"},
			{Name: "Port(s)", Type: "string"},
			{Name: columnName, Type: "string"},
		},
		Rows: []metav1.TableRow{
			{
				Cells: []interface{}{
					service.Name,
					ServiceTypeToString(service.Spec.Type),
					PortsToString(service.Spec.Ports),
					cellValue,
				},
			},
		},
	}

	out := bytes.NewBuffer([]byte{})
	printer := printers.NewTablePrinter(printers.PrintOptions{})
	printer.PrintObj(table, out)

	fmt.Fprint(w, out.String())

	return
}
