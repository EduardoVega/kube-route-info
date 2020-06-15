package main

import (
	"os"
	//"fmt"
	//"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"kube-route-info/pkg/cmd"
)

func main() {

	// flags := pflag.NewFlagSet("kubectl-graph", pflag.ExitOnError)

	// pflag.CommandLine = flags

	root := cmd.NewCmd(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}