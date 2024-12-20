package vb6

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Property struct {
	Name       string
	Value      string
	Properties PropertyMap
}

type PropertyMap map[string]Property

type Control struct {
	Form       *Form
	TypeName   string
	Name       string
	Children   []*Control
	Properties PropertyMap
}

type Attribute struct {
	Name  string
	Value string
}

type Form struct {
	Filename   string
	Folder     string
	Root       *Control
	Attributes []Attribute
	Script     string
}

// Errors returned by Load
var (
	ErrFileNotExist          = errors.New("file not found")
	ErrFileEmpty             = errors.New("file is empty")
	ErrUnexpectedEOF         = errors.New("unexpected end of file")
	ErrBadVersion            = errors.New("version is missing or version is not supported")
	ErrMalformed             = errors.New("malformed statement")
	ErrExpectedBegin         = errors.New("expected Begin keyword")
	ErrExpectedBeginProperty = errors.New("expected BeginProperty keyword")
)

// readLines reads all lines from a file and returns them as a slice of strings.
func readLines(file *os.File) []string {
	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// readVersion reads the version from the first line of the file.
func readVersion(line string) (string, error) {
	if !strings.HasPrefix(line, "VERSION ") {
		return "", ErrBadVersion
	}
	return line[8:], nil
}

// readProperty reads a property from a line and adds it to the properties map.
func readProperty(line string, properties PropertyMap) {
	equal := strings.Index(line, "=")
	if equal != -1 {
		name := strings.TrimSpace(line[:equal])
		properties[name] = Property{
			Name:       name,
			Value:      strings.TrimSpace(line[equal+1:]),
			Properties: make(PropertyMap),
		}
	}
}

func readComplexProperty(lines []string, properties PropertyMap) ([]string, error) {
	line := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line, "BeginProperty ") {
		return lines, ErrExpectedBeginProperty
	}

	name := strings.TrimSpace(line[14:])
	results := make(PropertyMap)
	lines = lines[1:]
	for len(lines) > 0 {
		line := strings.TrimSpace(lines[0])
		if line == "EndProperty" {
			break
		}
		readProperty(line, results)
		lines = lines[1:]
	}
	lines = lines[1:]
	properties[name] = Property{
		Name:       name,
		Value:      "",
		Properties: results,
	}
	return lines, nil
}

// readControl reads a control from the lines and returns the remaining lines and the control.
func readControl(lines []string, form *Form) ([]string, *Control, error) {
	if len(lines) == 0 {
		return lines, nil, ErrUnexpectedEOF
	}

	line := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line, "Begin ") {
		return lines, nil, ErrExpectedBegin
	}

	line = line[6:]

	space := strings.Index(line, " ")
	if space == -1 {
		return lines, nil, ErrMalformed
	}

	control := &Control{
		Form:       form,
		TypeName:   strings.TrimSpace(line[:space]),
		Name:       strings.TrimSpace(line[space+1:]),
		Children:   make([]*Control, 0),
		Properties: make(PropertyMap),
	}

	lines = lines[1:]
	for len(lines) > 0 {
		line := strings.TrimSpace(lines[0])

		// Stop reading once we reach the end
		if line == "End" {
			lines = lines[1:]
			break
		}

		// Read the details of nested children
		if strings.HasPrefix(line, "Begin ") {
			newLines, child, err := readControl(lines, form)
			if err != nil {
				return newLines, nil, err
			}
			lines = newLines

			control.Children = append(control.Children, child)
			continue
		}

		// Read a complex property
		if strings.HasPrefix(line, "BeginProperty ") {
			newLines, err := readComplexProperty(lines, control.Properties)
			if err != nil {
				return newLines, nil, err
			}
			lines = newLines
			continue
		}

		// Read a simple property
		readProperty(line, control.Properties)

		// Advance to the next line
		lines = lines[1:]
	}

	return lines, control, nil
}

// readAttributes reads the attributes from the lines and returns the remaining lines and the attributes.
func readAttributes(lines []string) ([]string, []Attribute) {
	attributes := make([]Attribute, 0)

	for len(lines) > 0 {
		line := strings.TrimSpace(lines[0])
		if !strings.HasPrefix(line, "Attribute ") {
			break
		}

		line = line[10:]

		equal := strings.Index(line, "=")
		if equal != -1 {
			attributes = append(attributes, Attribute{
				Name:  strings.TrimSpace(line[:equal]),
				Value: strings.TrimSpace(line[equal+1:]),
			})
		}

		lines = lines[1:]
	}

	return lines, attributes
}

func Load(path string) (*Form, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, ErrFileNotExist
	}

	defer file.Close()

	lines := readLines(file)
	if len(lines) == 0 {
		return nil, ErrFileEmpty
	}

	v, err := readVersion(lines[0])
	if err != nil {
		return nil, err
	}

	if v != "5.00" {
		return nil, ErrBadVersion
	}

	form := &Form{
		Filename: path,
		Folder:   filepath.Dir(path),
	}

	lines, root, err := readControl(lines[1:], form)
	if err != nil {
		return nil, err
	}

	lines, attr := readAttributes(lines)

	form.Root = root
	form.Attributes = attr
	form.Script = strings.Join(lines, "\n")

	return form, nil
}
