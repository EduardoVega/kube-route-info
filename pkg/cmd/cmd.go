package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // combined authprovider import
)

var cmdExample = `
	# View the route information of the service my-service
	%[1]s route-info service my-service

	# View the route information of the ingress my-ingress in namespace my-namespace
	%[1]s route-info ingress my-ingress --namespace my-namespace
`

// Resource provides the information required to get
// the route configuration from ingress and service objects
type Resource struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
	resourceInterface ResourceInterface
	printGraph        bool
	resourceType      string
	resourceName      string
}

// ResourceInterface defines the methods the must be
// implemented in the ingress and service structs
type ResourceInterface interface {
	PrintTable(string, io.Writer) error
	PrintGraph(string, io.Writer) error
}

// NewResource creates a new Resource struct with the required information
// to get the route configuration from ingress and service objects
func NewResource(streams genericclioptions.IOStreams) *Resource {
	return &Resource{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

// NewCmd returns the new route-info cobra command
func NewCmd(streams genericclioptions.IOStreams) *cobra.Command {
	r := NewResource(streams)

	cmd := &cobra.Command{
		Use:          "route-info [TYPE] [NAME] [flags]",
		Short:        "View route information from ingresses or services to pods",
		Example:      fmt.Sprintf(cmdExample, "kubectl"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := r.Validate(args); err != nil {
				return err
			}
			if err := r.Complete(c, args); err != nil {
				return err
			}
			if err := r.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&r.printGraph, "graph", r.printGraph, "if true, print the route information in a tree graph format")
	r.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Validate ensures that all required arguments and flags values are provided
func (r *Resource) Validate(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("requires 2 arguments. Run: kubectl route-info -h")
	}

	if args[0] != "ingress" && args[0] != "service" {
		return fmt.Errorf("only ingress and service types are supported. Run: kubectl route-info -h")
	}

	return nil
}

// Complete sets all information required for the command
func (r *Resource) Complete(cmd *cobra.Command, args []string) error {

	var err error

	r.resourceType = args[0]
	r.resourceName = args[1]

	config, err := r.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	// TODO: Test this with the kubectl plugin ns
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}

	if namespace == "" {
		namespace = "default"
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	client := NewClient(clientset, namespace)

	switch r.resourceType {
	case "service":
		r.resourceInterface = NewService(client, namespace)

	case "ingress":
		r.resourceInterface = NewIngress(client, namespace)
	}

	return nil
}

// Run executes the command of printing the route information
func (r *Resource) Run() (err error) {

	if r.printGraph {
		err := r.resourceInterface.PrintGraph(r.resourceName, os.Stdout)
		if err != nil {
			return err
		}
	} else {
		err := r.resourceInterface.PrintTable(r.resourceName, os.Stdout)
		if err != nil {
			return err
		}
	}

	return
}
