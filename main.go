package main

import (
	"github.com/guthius/vb6conv/export"
	"github.com/guthius/vb6conv/vb6"
)

func main() {
	form, err := vb6.Load("D:/Projects/VB6/ms3_0_3/client/frmCredits.frm")
	if err != nil {
		panic(err)
	}

	project := export.ProjectInfo{
		Name:      "ms3_0_3",
		Namespace: "Mirage",
	}

	export.Export(&project, form)
}
