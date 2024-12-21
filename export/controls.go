package export

import (
	"fmt"

	"github.com/guthius/vb6conv/vb6"
)

type Control struct {
	Name      string
	TypeName  string
	Resources map[string]any
	Props     map[string]string
	Children  []*Control
	MustInit  bool
}

type ControlBuilder func(c *vb6.Control) *Control

func applyDefaultProps(c *vb6.Control, props map[string]string) {
	if v, ok := vb6.GetBool("Visible", c.Properties); ok {
		props["Visible"] = toBool(v)
	}

	if v, ok := vb6.GetColor("BackColor", c.Properties); ok {
		props["BackColor"] = toColor(v)
	}

	if font, ok := vb6.GetFont("Font", c.Properties); ok {
		props["Font"] = toFont(font)
	}

	if v, ok := vb6.GetInt("TabIndex", c.Properties); ok {
		props["TabIndex"] = toInt(v)
	}
}

func FormBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)
	props["AutoScaleDimensions"] = toSizeF(6, 13)
	props["AutoScaleMode"] = "System.Windows.Forms.AutoScaleMode.None"

	if w, h, ok := vb6.GetVector2("ClientWidth", "ClientHeight", c.Properties); ok {
		props["ClientSize"] = toSize(w, h)
	}

	if v, ok := vb6.GetProp("Caption", c.Properties); ok {
		props["Text"] = v
	} else {
		props["Text"] = c.Name
	}

	if v, ok := vb6.GetColor("BackColor", c.Properties); ok {
		props["BackColor"] = toColor(v)
	}

	applyDefaultProps(c, props)

	// TODO: ClientLeft, ClientTop
	// TODO: ControlBox
	// TODO: StartUpPosition, WindowState
	// TODO: MaxButton, MinButton
	// TODO: KeyPreview

	return &Control{
		Name:      c.Name,
		TypeName:  "System.Windows.Forms.Form",
		Resources: make(map[string]any),
		Props:     props,
		Children:  buildControlSlice(c.Children),
		MustInit:  false,
	}
}

func PictureBoxBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)

	if x, y, ok := vb6.GetVector2("Left", "Top", c.Properties); ok {
		props["Location"] = toPoint(x, y)
	}
	if w, h, ok := vb6.GetVector2("Width", "Height", c.Properties); ok {
		props["Size"] = toSize(w, h)
	}

	if autoSize, _ := vb6.GetBool("AutoSize", c.Properties); autoSize {
		props["SizeMode"] = "System.Windows.Forms.PictureBoxSizeMode.AutoSize"
	}

	if v, ok := vb6.GetInt("Appearance", c.Properties); ok {
		if v == 0 {
			props["BorderStyle"] = "System.Windows.Forms.BorderStyle.None"
		} else {
			props["BorderStyle"] = "System.Windows.Forms.BorderStyle.Fixed3D"
		}
	}

	applyDefaultProps(c, props)

	resources := make(map[string]any)
	if locator, ok := vb6.GetProp("Picture", c.Properties); ok {
		bytes, err := vb6.GetResource(c, locator)
		if err == nil {
			resource := fmt.Sprintf("%s.Image", c.Name)
			resources[resource] = bytes
			props["Image"] = fmt.Sprintf("((System.Drawing.Image)(resources.GetObject(\"%s\")))", resource)
		} else {
			fmt.Printf("unable to load resource: %s (%v)\n", locator, err)
		}
	}

	props["TabStop"] = "false"

	return &Control{
		Name:      c.Name,
		TypeName:  "System.Windows.Forms.PictureBox",
		Resources: resources,
		Props:     props,
		Children:  buildControlSlice(c.Children),
		MustInit:  true,
	}
}

func LabelBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)

	if alignment, ok := vb6.GetInt("Alignment", c.Properties); ok {
		switch alignment {
		case 0:
			props["TextAlign"] = "System.Drawing.ContentAlignment.TopLeft"
		case 1:
			props["TextAlign"] = "System.Drawing.ContentAlignment.TopRight"
		case 2:
			props["TextAlign"] = "System.Drawing.ContentAlignment.TopCenter"
		}
	}

	if x, y, ok := vb6.GetVector2("Left", "Top", c.Properties); ok {
		props["Location"] = toPoint(x, y)
	}

	if w, h, ok := vb6.GetVector2("Width", "Height", c.Properties); ok {
		props["Size"] = toSize(w, h)
	}

	if foreColor, ok := vb6.GetColor("ForeColor", c.Properties); ok {
		props["ForeColor"] = toColor(foreColor)
	}

	if caption, ok := vb6.GetProp("Caption", c.Properties); ok {
		props["Text"] = caption
	}

	applyDefaultProps(c, props)

	if backStyle, ok := vb6.GetInt("BackStyle", c.Properties); ok {
		if backStyle == 0 { // Transparent
			props["BackColor"] = "System.Drawing.Color.Transparent"
		}
	}

	return &Control{
		Name:      c.Name,
		TypeName:  "System.Windows.Forms.Label",
		Resources: make(map[string]any),
		Props:     props,
		Children:  buildControlSlice(c.Children),
		MustInit:  false,
	}
}

func buildControlSlice(controls []*vb6.Control) []*Control {
	result := make([]*Control, 0, len(controls))
	for _, c := range controls {
		if control := buildControl(c); control != nil {
			result = append(result, control)
		}
	}
	return result
}

func buildControl(c *vb6.Control) *Control {
	var builder ControlBuilder

	switch {
	case c.TypeName == "VB.Form":
		builder = FormBuilder
	case c.TypeName == "VB.PictureBox":
		builder = PictureBoxBuilder
	case c.TypeName == "VB.Label":
		builder = LabelBuilder
	default:
		return nil
	}

	return builder(c)
}
