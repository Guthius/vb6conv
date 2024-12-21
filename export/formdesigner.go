package export

import (
	"fmt"
	"os"
	"path/filepath"
)

func exportFormDesigner(p *ProjectInfo, f *Control, hasRes bool) error {
	filename := filepath.Join(p.Output, f.Name+".Designer.cs")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := NewExportWriter(file)
	writer.Write("using System;")
	writer.Writeln()
	writer.Writef("namespace %s;", p.Namespace)
	writer.Writeln()
	writer.Writef("partial class %s", f.Name)
	writer.Write("{")
	writer.WriteIndent(func() {
		writer.Write("/// <summary>")
		writer.Write("/// Required designer variable.")
		writer.Write("/// </summary>")
		writer.Write("private System.ComponentModel.IContainer components = null;")
		writer.Writeln()
		writer.Write("/// <summary>")
		writer.Write("/// Clean up any resources being used.")
		writer.Write("/// </summary>")
		writer.Write("/// <param name=\"disposing\">true if managed resources should be disposed; otherwise, false.</param>")
		writer.Write("protected override void Dispose(bool disposing)")
		writer.Write("{")
		writer.WriteIndent(func() {
			writer.Write("if (disposing && (components != null))")
			writer.Write("{")
			writer.WriteIndent(func() {
				writer.Write("components.Dispose();")
			})
			writer.Write("}")
			writer.Write("base.Dispose(disposing);")
		})
		writer.Write("}")
		writer.Writeln()
		writeInitializeComponent(p, f, writer, hasRes)
		writer.Writeln()
		for _, c := range f.Children {
			writeControlDefinitions(p, c, writer)
		}
	})

	writer.Write("}")

	return nil
}

func writeControlDefinitions(p *ProjectInfo, f *Control, w *ExportWriter) {
	for _, c := range f.Children {
		writeControlDefinitions(p, c, w)
	}
	w.Writef("private %s %s;", f.TypeName, f.Name)
}

func writeControlInitializers(p *ProjectInfo, f *Control, w *ExportWriter) {
	for _, c := range f.Children {
		writeControlInitializers(p, c, w)
	}
	if f.IsComponent {
		w.Writef("%s = new %s(this.components);", f.Name, f.TypeName)
	} else {
		w.Writef("%s = new %s();", f.Name, f.TypeName)
	}
}

func writeControlProperties(p *ProjectInfo, f *Control, w *ExportWriter, root bool) {
	for _, c := range f.Children {
		writeControlProperties(p, c, w, false)
	}
	name := "this"
	if !root {
		name = fmt.Sprintf("this.%s", f.Name)
	}
	w.Write("//")
	w.Writef("// %s", f.Name)
	w.Write("//")
	if !f.SkipName {
		w.Writef("%s.Name = \"%s\";", name, f.Name)
	}
	for k, v := range f.Props {
		w.Writef("%s.%s = %s;", name, k, v)
	}
	for k, v := range f.PropCalls {
		w.Writef("%s.%s.%s;", name, k, v)
	}
	for _, c := range f.Children {
		if !c.SkipAdd {
			w.Writef("%s.Controls.Add(this.%s);", name, c.Name)
		}
	}
}

func getControlsToInit(f *Control, current []*Control) []*Control {
	if f.MustInit {
		current = append(current, f)
	}
	for _, c := range f.Children {
		current = getControlsToInit(c, current)
	}
	return current
}

func writeBeginInit(w *ExportWriter, init []*Control) {
	if len(init) == 0 {
		return
	}
	for _, c := range init {
		w.Writef("((System.ComponentModel.ISupportInitialize)(this.%s)).BeginInit();", c.Name)
	}
	w.Write("this.SuspendLayout();")
}

func writeEndInit(w *ExportWriter, init []*Control) {
	if len(init) == 0 {
		return
	}
	for _, c := range init {
		w.Writef("((System.ComponentModel.ISupportInitialize)(this.%s)).EndInit();", c.Name)
	}
	w.Write("this.ResumeLayout(false);")
	w.Write("this.PerformLayout();")
}

func writeInitializeComponent(p *ProjectInfo, f *Control, w *ExportWriter, hasRes bool) {
	init := getControlsToInit(f, []*Control{})

	w.Write("#region Windows Form Designer generated code")
	w.Writeln()
	w.Write("/// <summary>")
	w.Write("/// Required method for Designer support - do not modify")
	w.Write("/// the contents of this method with the code editor.")
	w.Write("/// </summary>")
	w.Write("private void InitializeComponent()")
	w.Write("{")
	w.WriteIndent(func() {
		w.Write("this.components = new System.ComponentModel.Container();")
		if hasRes {
			w.Writef("System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(%s));", f.Name)
		}
		for _, c := range f.Children {
			writeControlInitializers(p, c, w)
		}
		writeBeginInit(w, init)
		writeControlProperties(p, f, w, true)
		writeEndInit(w, init)
	})

	w.Write("}")
	w.Writeln()
	w.Write("#endregion")
}
