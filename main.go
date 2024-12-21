package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guthius/vb6conv/export"
	"github.com/guthius/vb6conv/vb6"
	"github.com/guthius/vb6conv/vb6/vbp"
	"github.com/spf13/pflag"
)

var (
	project   string
	namespace string
	output    string
)

func main() {
	pflag.StringVarP(&project, "project", "p", "", "Path to the project file (required)")
	pflag.StringVarP(&output, "output", "o", "", "Output directory (required)")
	pflag.StringVarP(&namespace, "namespace", "n", "", "Namespace for the project (optional)")
	pflag.Parse()

	if len(project) == 0 {
		fmt.Fprintln(os.Stderr, "Error:  missing required argument 'project'")
		pflag.Usage()
		os.Exit(1)
	}

	if len(output) == 0 {
		fmt.Fprintln(os.Stderr, "Error:  missing required argument 'output'")
		pflag.Usage()
		os.Exit(1)
	}

	output, err := filepath.Abs(output)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		err := os.MkdirAll(output, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	vbproj, err := vbp.Open(project)
	if err != nil {
		panic(err)
	}

	if len(namespace) == 0 {
		namespace = vbproj.Name
	}

	project := export.ProjectInfo{
		Name:      vbproj.Name,
		Namespace: namespace,
		Output:    output,
	}

	for _, form := range vbproj.Forms {
		f, err := vb6.Load(form)
		if err != nil {
			panic(err)
		}

		export.Export(&project, f)
	}
}
