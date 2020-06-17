package cmd

import (
	"context"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"os"
	"github.com/liggitt/tabwriter"
)

type Ingress struct {
	Client *kubernetes.Clientset
	Namespace string
	Name string 
	Rules []RulesInfo
}

type RulesInfo struct {
	Host string
	IngressRules []IngressRuleInfo
}

type IngressRuleInfo struct {
	Path string
	ServiceName string
	ServicePort string
	ServiceInfo *Service
}

// type IngressInterface interface {
// 	GetInformation(name string)
// 	PrintInformation()
// }

func NewIngress(client *kubernetes.Clientset, namespace string) *Ingress {
	return &Ingress {
		Client: client,
		Namespace: namespace,
	}
}

func (i *Ingress) GetInformation(name string) (err error) {

	ingress, err := i.Client.NetworkingV1beta1().Ingresses(i.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if (err != nil){
		return err
	}
	
	i.Name = ingress.Name
	i.Rules = []RulesInfo{}

	for _, rule := range ingress.Spec.Rules{
		rulesInfo := RulesInfo {}

		rulesInfo.Host = rule.Host
		rulesInfo.IngressRules = []IngressRuleInfo{}

		for _, path := range rule.IngressRuleValue.HTTP.Paths{
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

func (i *Ingress) PrintInformation() {
	// initialize tabwriter
	w := new(tabwriter.Writer)
	
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	
	defer w.Flush()
	
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t", "NAME", "HOST", "PATH", "PORT", "SERVICE", "TYPE", "PORT(S)", "POD(S)")
	
	for _, rule := range i.Rules {
		for _, ingressRule := range rule.IngressRules {
			fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t", i.Name, rule.Host, ingressRule.Path, ingressRule.ServicePort, ingressRule.ServiceInfo.Name, ingressRule.ServiceInfo.ServiceType, ingressRule.ServiceInfo.Ports, ingressRule.ServiceInfo.Pods)
		}
	}

	fmt.Fprintf(w, "\n")
}

func GetServiceInfo(client *kubernetes.Clientset, namespace, serviceName string) (service *Service) {
	service = NewService(client, namespace)

	var err error
	err = service.GetInformation(serviceName)
	if (err != nil){
		service.Name = "<none>"
		service.ServiceType = "<none>"
		service.Pods = "<none>"
		service.Ports = "<none>"
	}

	return
}

