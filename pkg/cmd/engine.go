package cmd

import (

	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	
)

var cmdExample = `
	# View the route information of the service my-service
	%[1]s route-info service my-service

	# View the route information of the ingress my-ingress in namespace my-namespace
	%[1]s route-info ingress my-ingress --namespace my-namespace
`

type Engine struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
	namespace string
	config *rest.Config
	args []string
	resourceInterface ResourceInterface
}

type ResourceInterface interface {
	GetInformation(name string) error
	PrintInformation()
}

func NewEngine(streams genericclioptions.IOStreams) *Engine {
	return &Engine{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams: streams,
	}
}

func NewCmd(streams genericclioptions.IOStreams) *cobra.Command {
	e := NewEngine(streams)

	cmd := &cobra.Command{
		Use:          "route-info [TYPE] [NAME] [flags]",
		Short:        "View route information from ingresses or services to pods",
		Example:      fmt.Sprintf(cmdExample, "kubectl"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := e.Complete(c, args); err != nil {
				return err
			}
			if err := e.Validate(); err != nil {
				return err
			}
			if err := e.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	//cmd.Flags().BoolVar(&g.listNamespaces, "list", g.listNamespaces, "if true, print the list of all namespaces in the current KUBECONFIG")
	e.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete sets all information required for the command
func (e *Engine) Complete(cmd *cobra.Command, args []string) error {
	
	e.args = args

	var err error
	e.config, err = e.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	e.namespace, err = cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}

	// TODO: Get namespace from context in kubeconfig

	if ( e.namespace == "" ){
		e.namespace = "default"
	}

	clientset, err := kubernetes.NewForConfig(e.config)
	if err != nil {
		return err
	}

	switch e.args[0] {
		case "service":
			e.resourceInterface = NewService(clientset, e.namespace)

		case "ingress":
			e.resourceInterface = NewIngress(clientset, e.namespace)
	}

	return nil
}

// Validate ensures that all required arguments and flags values are provided
func (e *Engine) Validate() error {
	if len(e.args) != 2 {
		return fmt.Errorf("requires 2 arguments. Run: kubectl route-info -h")
	}

	if e.resourceInterface == nil {
		return fmt.Errorf("resource type is not valid. Run: kubectl route-info -h")
	}

	return nil
}

// Run passes the information and executes the command 
func (e *Engine) Run() (err error) {
	
	err = e.resourceInterface.GetInformation(e.args[1])
	if err != nil {
		return err
	}
	e.resourceInterface.PrintInformation()

	return nil
}