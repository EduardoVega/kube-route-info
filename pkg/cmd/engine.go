package cmd

import (

	"fmt"
	"github.com/spf13/cobra"
	// "k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	// "encoding/json"
	
)

var cmdExample = `
	# work in progress
	%[1]s work in progress
`

type Engine struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
	namespace string
	config *rest.Config
	args []string
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
		Use:          "path [TYPE] [NAME] [flags]",
		Short:        "View the path that a request takes until it reaches pods",
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

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (e *Engine) Validate() error {
	if len(e.args) != 2 {
		return fmt.Errorf("requires 2 arguments. Run: kubectl path -h")
	}

	return nil
}

func (e *Engine) Run() error{
	
	// config, err := e.configFlags.ToRESTConfig()
	// if err != nil {
	// 	panic(err.Error())
	// }

	clientset, err := kubernetes.NewForConfig(e.config)
	if err != nil {
		return err
	}

	// TODO: Use interfaces
	
	// switch e.args[0] {
	// 	case "service":
	// 		service := Service{
	// 			clientset,
	// 			e.namespace,
	// 		}

	// 		data, err := service.GetInformation(e.args[1])
	// 		if err != nil {
	// 			return err
	// 		}

	// 		// jsonData, err := json.Marshal(data)
	// 		// if err != nil {
	// 		// 	return err
	// 		// }

	// 		// fmt.Println(string(jsonData))
	// 		service.PrintInformation(data)
	// 	default:
	// 		return fmt.Errorf("resource type is not valid. Run: kubectl path -h")
	// }

	ingress := NewIngress(clientset, e.namespace)

	err = ingress.GetInformation(e.args[1])
	if err != nil {
		return err
	}

	return nil
}