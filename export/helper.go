package export

import "fmt"

// These are a bunch of helper methods to convert values to C# initializer statements

func toInt(v int) string {
	return fmt.Sprintf("%v", v)
}

func toColor(c uint32) string {
	r := (c & 0xFF0000) >> 16
	g := (c & 0x00FF00) >> 8
	b := (c & 0x0000FF)

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
