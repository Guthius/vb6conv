package export

import (
	"fmt"
	"os"
	"strconv"

	"github.com/guthius/vb6conv/vb6"
)

type Control struct {
	Name        string
	TypeName    string
	Resources   map[string]any
	Props       map[string]string
	PropCalls   map[string]string
	Children    []*Control
	MustInit    bool
	SkipAdd     bool // Indicates that the control should not be added to the parent's control collection
	SkipName    bool // Indicates that the control's name property should not be generated
	IsComponent bool
}

type ControlBuilder func(c *vb6.Control) *Control

func applyDefaultProps(c *vb6.Control, props map[string]string) {
	if visible, ok := vb6.GetBool("Visible", c.Properties); ok {
		props["Visible"] = toBool(visible)
	}

	if backColor, ok := vb6.GetColor("BackColor", c.Properties); ok {
		props["BackColor"] = toColor(backColor)
	}

	if font, ok := vb6.GetFont("Font", c.Properties); ok {
		props["Font"] = toFont(font)
	}
}

func applyDefaultPropsForControl(c *vb6.Control, props map[string]string) {
	applyDefaultProps(c, props)

	if x, y, ok := vb6.GetVector2("Left", "Top", c.Properties); ok {
		props["Location"] = toPoint(x, y)
	}

	if w, h, ok := vb6.GetVector2("Width", "Height", c.Properties); ok {
		props["Size"] = toSize(w, h)
	}

	if tabIndex, ok := vb6.GetInt("TabIndex", c.Properties); ok {
		props["TabIndex"] = toInt(tabIndex)
	}
}

func FormBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)

	applyDefaultProps(c, props)

	props["AutoScaleDimensions"] = toSizeF(6, 13)
	props["AutoScaleMode"] = "System.Windows.Forms.AutoScaleMode.None"

	if w, h, ok := vb6.GetVector2("ClientWidth", "ClientHeight", c.Properties); ok {
		props["ClientSize"] = toSize(w, h)
	}

	if caption, ok := vb6.GetProp("Caption", c.Properties); ok {
		props["Text"] = caption
	} else {
		props["Text"] = c.Name
	}

	if backColor, ok := vb6.GetColor("BackColor", c.Properties); ok {
		props["BackColor"] = toColor(backColor)
	}

	if startUpPosition, ok := vb6.GetInt("StartUpPosition", c.Properties); ok {
		switch startUpPosition {
		case 0:
			props["StartPosition"] = "System.Windows.Forms.FormStartPosition.Manual"
		case 1:
			props["StartPosition"] = "System.Windows.Forms.FormStartPosition.CenterParent"
		case 2:
			props["StartPosition"] = "System.Windows.Forms.FormStartPosition.CenterScreen"
		case 3:
			props["StartPosition"] = "System.Windows.Forms.FormStartPosition.WindowsDefaultLocation"
		}
	}

	if controlBox, ok := vb6.GetBool("ControlBox", c.Properties); ok {
		props["ControlBox"] = toBool(controlBox)
	}

	// TODO: WindowState
	// TODO: ShowInTaskbar   =   0   'False

	if minButton, ok := vb6.GetBool("MinButton", c.Properties); ok {
		props["MinimizeBox"] = toBool(minButton)
	}

	if maxButton, ok := vb6.GetBool("MaxButton", c.Properties); ok {
		props["MaximizeBox"] = toBool(maxButton)
	}

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

	applyDefaultPropsForControl(c, props)

	if autoSize, _ := vb6.GetBool("AutoSize", c.Properties); autoSize {
		props["SizeMode"] = "System.Windows.Forms.PictureBoxSizeMode.AutoSize"
	}

	if appearance, ok := vb6.GetInt("Appearance", c.Properties); ok {
		if appearance == 0 {
			props["BorderStyle"] = "System.Windows.Forms.BorderStyle.None"
		} else {
			props["BorderStyle"] = "System.Windows.Forms.BorderStyle.Fixed3D"
		}
	}

	resources := make(map[string]any)
	if locator, ok := vb6.GetProp("Picture", c.Properties); ok {
		bytes, err := vb6.GetBinaryResource(c, locator)
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

	applyDefaultPropsForControl(c, props)

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

	if foreColor, ok := vb6.GetColor("ForeColor", c.Properties); ok {
		props["ForeColor"] = toColor(foreColor)
	}

	if caption, ok := vb6.GetStr("Caption", c.Properties); ok {
		props["Text"] = toStr(caption)
	}

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

func TextBoxBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)

	applyDefaultPropsForControl(c, props)

	borderStyle, ok := vb6.GetInt("BorderStyle", c.Properties)
	if !ok {
		borderStyle = 1
	}

	if borderStyle == 0 {
		props["BorderStyle"] = "System.Windows.Forms.BorderStyle.None"
	} else {
		if appearance, ok := vb6.GetInt("Appearance", c.Properties); ok && appearance == 0 {
			props["BorderStyle"] = "System.Windows.Forms.BorderStyle.FixedSingle"
		}
	}

	if foreColor, ok := vb6.GetColor("ForeColor", c.Properties); ok {
		props["ForeColor"] = toColor(foreColor)
	}

	if maxLength, ok := vb6.GetInt("MaxLength", c.Properties); ok {
		props["MaxLength"] = toInt(maxLength)
	}

	if passwordChar, ok := vb6.GetProp("PasswordChar", c.Properties); ok {
		passwordChar, err := strconv.Unquote(passwordChar)
		if err == nil {
			props["PasswordChar"] = fmt.Sprintf("'%s'", passwordChar)
		}
	}

	if multiLine, ok := vb6.GetBool("MultiLine", c.Properties); ok {
		props["Multiline"] = toBool(multiLine)
	}

	// TODO:  IMEMode         =   3  'DISABLE

	return &Control{
		Name:      c.Name,
		TypeName:  "System.Windows.Forms.TextBox",
		Resources: make(map[string]any),
		Props:     props,
		Children:  buildControlSlice(c.Children),
		MustInit:  false,
	}
}

func FrameBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)

	applyDefaultPropsForControl(c, props)

	if caption, ok := vb6.GetStr("Caption", c.Properties); ok {
		props["Text"] = toStr(caption)
	}

	return &Control{
		Name:      c.Name,
		TypeName:  "System.Windows.Forms.GroupBox",
		Resources: make(map[string]any),
		Props:     props,
		Children:  buildControlSlice(c.Children),
		MustInit:  false,
	}
}

func CommandButtonBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)

	applyDefaultPropsForControl(c, props)

	if caption, ok := vb6.GetStr("Caption", c.Properties); ok {
		props["Text"] = toStr(caption)
	}

	return &Control{
		Name:      c.Name,
		TypeName:  "System.Windows.Forms.Button",
		Resources: make(map[string]any),
		Props:     props,
		Children:  buildControlSlice(c.Children),
		MustInit:  false,
	}
}

func ComboBoxBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)
	propCalls := make(map[string]string)

	applyDefaultPropsForControl(c, props)

	if style, ok := vb6.GetInt("Style", c.Properties); ok {
		switch style {
		case 0:
			props["DropDownStyle"] = "System.Windows.Forms.ComboBoxStyle.DropDown"
		case 1:
			props["DropDownStyle"] = "System.Windows.Forms.ComboBoxStyle.Simple"
		case 2:
			props["DropDownStyle"] = "System.Windows.Forms.ComboBoxStyle.DropDownList"
		}
	}

	if list, ok := vb6.GetProp("List", c.Properties); ok {
		items, err := vb6.GetListResource(c, list)
		if err == nil {
			props["FormattingEnabled"] = toBool(true)
			propCalls["Items"] = fmt.Sprintf("AddRange(%s)", toObjectArray(items))
		}
	}

	return &Control{
		Name:      c.Name,
		TypeName:  "System.Windows.Forms.ComboBox",
		Resources: make(map[string]any),
		Props:     props,
		PropCalls: propCalls,
		Children:  buildControlSlice(c.Children),
		MustInit:  false,
	}
}

func TimerBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)

	if interval, ok := vb6.GetInt("Interval", c.Properties); ok {
		props["Interval"] = toInt(interval)
	}

	return &Control{
		Name:        c.Name,
		TypeName:    "System.Windows.Forms.Timer",
		Resources:   make(map[string]any),
		Props:       props,
		Children:    buildControlSlice(c.Children),
		MustInit:    false,
		SkipAdd:     true,
		SkipName:    true,
		IsComponent: true,
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
	case c.TypeName == "VB.Menu":
		builder = MenuItemBuilder
	case c.TypeName == "VB.PictureBox":
		builder = PictureBoxBuilder
	case c.TypeName == "VB.Label":
		builder = LabelBuilder
	case c.TypeName == "VB.TextBox":
		builder = TextBoxBuilder
	case c.TypeName == "VB.Frame":
		builder = FrameBuilder
	case c.TypeName == "VB.CommandButton":
		builder = CommandButtonBuilder
	case c.TypeName == "VB.ComboBox":
		builder = ComboBoxBuilder
	case c.TypeName == "VB.Timer":
		builder = TimerBuilder
	default:
		return nil
	}

	return builder(c)
}

func MenuItemBuilder(c *vb6.Control) *Control {
	props := make(map[string]string)
	propCalls := make(map[string]string)

	if caption, ok := vb6.GetStr("Caption", c.Properties); ok {
		props["Text"] = toStr(caption)
	}

	children := buildMenuItemSlice(c.Children)
	if len(children) > 0 {
		childNames := make([]string, 0, len(children))
		for i, child := range children {
			child.Props["Index"] = toInt(i)
			childNames = append(childNames, fmt.Sprintf("this.%s", child.Name))
		}
		propCalls["MenuItems"] = fmt.Sprintf("AddRange(%s)", toArrayOfType(childNames, "System.Windows.Forms.MenuItem"))
	}

	return &Control{
		Name:      c.Name,
		TypeName:  "System.Windows.Forms.MenuItem",
		Resources: make(map[string]any),
		Props:     props,
		PropCalls: propCalls,
		Children:  children,
		MustInit:  false,
		SkipAdd:   true,
	}
}

func buildMenuItemSlice(controls []*vb6.Control) []*Control {
	result := make([]*Control, 0, len(controls))
	for _, c := range controls {
		if c.TypeName != "VB.Menu" {
			fmt.Fprintf(os.Stderr, "unexpected control type in menu: %s\n", c.TypeName)
			continue
		}
		control := MenuItemBuilder(c)
		if control != nil {
			result = append(result, control)
		}
	}
	return result
}

func buildMenu(f *Control) {
	menuItems := make([]*Control, 0)
	menuItemNames := make([]string, 0, len(f.Children))

	newChildren := make([]*Control, 0, len(f.Children))
	for _, c := range f.Children {
		if c.TypeName == "System.Windows.Forms.MenuItem" {
			menuItems = append(menuItems, c)
			menuItemNames = append(menuItemNames, fmt.Sprintf("this.%s", c.Name))
		} else {
			newChildren = append(newChildren, c)
		}
	}

	if len(menuItems) == 0 {
		return
	}

	f.Children = newChildren

	propCalls := make(map[string]string)
	propCalls["MenuItems"] = fmt.Sprintf("AddRange(%s)", toArrayOfType(menuItemNames, "System.Windows.Forms.MenuItem"))

	menu := &Control{
		Name:        "mainMenu1",
		TypeName:    "System.Windows.Forms.MainMenu",
		Resources:   make(map[string]any),
		PropCalls:   propCalls,
		Children:    menuItems,
		MustInit:    false,
		SkipAdd:     true,
		IsComponent: true,
	}

	f.Props["Menu"] = "this.mainMenu1"
	f.Children = append(f.Children, menu)
}
