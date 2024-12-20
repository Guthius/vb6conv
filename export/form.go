package export

import (
	"os"
	"path/filepath"
)

func exportForm(p *ProjectInfo, f *Control) error {
	filename := filepath.Join(p.Output, f.Name+".cs")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := NewExportWriter(file)
	writer.Write("using System;")
	writer.Write("using System.Collections.Generic;")
	writer.Write("using System.ComponentModel;")
	writer.Write("using System.Data;")
	writer.Write("using System.Drawing;")
	writer.Write("using System.Linq;")
	writer.Write("using System.Text;")
	writer.Write("using System.Threading.Tasks;")
	writer.Write("using System.Windows.Forms;")
	writer.Writeln()
	writer.Writef("namespace %s;", p.Namespace)
	writer.Writeln()
	writer.Writef("public partial class %s : Form", f.Name)
	writer.Write("{")
	writer.WriteIndent(func() {
		writer.Writef("public %s()", f.Name)
		writer.Write("{")
		writer.WriteIndent(func() {
			writer.Write("InitializeComponent();")
		})
		writer.Write("}")
	})
	writer.Write("}")

	return nil
}
