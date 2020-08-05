package main

import (
	"os"

	"kube-route-info/pkg/cmd"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {

	root := cmd.NewCmd(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
