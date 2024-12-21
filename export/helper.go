package export

import (
	"fmt"
	"strconv"

	"github.com/guthius/vb6conv/vb6"
)

// These are a bunch of helper methods to convert values to C# initializer statements

func toInt(v int) string {
	return fmt.Sprintf("%v", v)
}

func toBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func toStr(s string) string {
	return strconv.Quote(s)
}

func toColor(c uint32) string {
	r := (c & 0x0000FF)
	g := (c & 0x00FF00) >> 8
	b := (c & 0xFF0000) >> 16

	return fmt.Sprintf("System.Drawing.Color.FromArgb(%d, %d, %d)", r, g, b)
}

func toPoint(x int, y int) string {
	return fmt.Sprintf("new System.Drawing.Point(%v, %v)", x, y)
}

func toSizeF(w float32, h float32) string {
	return fmt.Sprintf("new System.Drawing.SizeF(%vF, %vF)", w, h)
}

func toSize(w int, h int) string {
	return fmt.Sprintf("new System.Drawing.Size(%v, %v)", w, h)
}

func toFont(f *vb6.Font) string {
	var fontStyle string
	switch f.Weight {
	case 700:
		fontStyle = "System.Drawing.FontStyle.Bold"
	default:
		fontStyle = "System.Drawing.FontStyle.Regular"
	}

	if f.Italic {
		fontStyle += " | System.Drawing.FontStyle.Italic"
	}
	if f.Underline {
		fontStyle += " | System.Drawing.FontStyle.Underline"
	}
	if f.Strikethrough {
		fontStyle += " | System.Drawing.FontStyle.Strikeout"
	}

	graphicsUnit := "System.Drawing.GraphicsUnit.Point"

	return fmt.Sprintf("new System.Drawing.Font(\"%v\", %vF, %v, %v, ((byte)(%v)))", f.Family, f.Size, fontStyle, graphicsUnit, f.Charset)
}
