package export

import (
	"fmt"
	"os"
)

type ExportWriter struct {
	file   *os.File
	indent int
}

func NewExportWriter(file *os.File) *ExportWriter {
	return &ExportWriter{
		file: file,
	}
}

func (w *ExportWriter) Write(s string) {
	for i := 0; i < w.indent; i++ {
		w.file.WriteString("\t")
	}
	w.file.WriteString(s)
	w.file.WriteString("\n")
}

func (w *ExportWriter) Writeln() {
	w.file.WriteString("\n")
}

func (w *ExportWriter) WriteIndent(write func()) {
	w.indent++
	write()
	w.indent--
}

func (w *ExportWriter) Writef(format string, args ...interface{}) {
	w.Write(fmt.Sprintf(format, args...))
}

func (w *ExportWriter) Indent() {
	w.indent++
}

func (w *ExportWriter) Unindent() {
	w.indent--
}
