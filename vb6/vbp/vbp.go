package vbp

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Reference struct {
	Guid    string
	Version string
	LCID    int
	Path    string
	Name    string
}

type Object struct {
	Guid    string
	Version string
	LCID    int
	Path    string
}

type Module struct {
	Name     string
	Filename string
}

type Version struct {
	Major    int
	Minor    int
	Revision int
}

type Project struct {
	Name       string
	Folder     string
	References []*Reference
	Objects    []*Object
	Modules    []*Module
	Forms      []string
	Startup    string
	Title      string
	ExeName32  string
	Version    Version
}

func Open(name string) (*Project, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	project := &Project{
		Name:       name,
		Folder:     filepath.Dir(name),
		References: make([]*Reference, 0),
		Objects:    make([]*Object, 0),
		Modules:    make([]*Module, 0),
		Forms:      make([]string, 0),
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		eq := strings.Index(line, "=")
		if eq == -1 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		value := strings.TrimSpace(line[eq+1:])
		switch key {
		case "Type":
			if value != "Exe" {
				return nil, errors.New("only Exe projects are supported")
			}
		case "Reference":
			ref, err := parseReference(value)
			if err != nil {
				return nil, err
			}
			project.References = append(project.References, ref)
		case "Object":
			obj, err := parseObject(value)
			if err != nil {
				return nil, err
			}
			project.Objects = append(project.Objects, obj)
		case "Module":
			mod, err := parseModule(value)
			if err != nil {
				return nil, err
			}
			project.Modules = append(project.Modules, mod)
		case "Form":
			project.Forms = append(project.Forms, value)
		case "Name":
			value = strings.TrimSpace(value)
			if len(value) > 0 {
				s, err := strconv.Unquote(value)
				if err != nil {
					return nil, err
				}
				project.Name = s
			}
		case "Startup":
			s, err := strconv.Unquote(value)
			if err != nil {
				return nil, err
			}
			project.Startup = s
		case "Title":
			s, err := strconv.Unquote(value)
			if err != nil {
				return nil, err
			}
			project.Title = s
		case "ExeName32":
			s, err := strconv.Unquote(value)
			if err != nil {
				return nil, err
			}
			project.ExeName32 = s
		case "MajorVer":
			project.Version.Major, _ = strconv.Atoi(value)
		case "MinorVer":
			project.Version.Minor, _ = strconv.Atoi(value)
		case "RevisionVer":
			project.Version.Revision, _ = strconv.Atoi(value)
		}
	}
	return project, nil
}

func parseGuid(s string) string {
	s = strings.TrimPrefix(s, "*\\G{")
	s = strings.TrimSuffix(s, "}")
	return s
}

func parseReference(s string) (*Reference, error) {
	tok := strings.Split(s, "#")
	if len(tok) != 5 {
		return nil, errors.New("invalid reference")
	}
	lcid, err := strconv.Atoi(tok[2])
	if err != nil {
		return nil, err
	}
	ref := &Reference{
		Guid:    parseGuid(tok[0]),
		Version: tok[1],
		LCID:    lcid,
		Path:    tok[3],
		Name:    tok[4],
	}
	return ref, nil
}

func parseObject(s string) (*Object, error) {
	tok := strings.Split(s, "#")
	if len(tok) != 3 {
		return nil, errors.New("invalid object")
	}
	tok2 := strings.Split(tok[2], ";")
	if len(tok2) != 2 {
		return nil, errors.New("invalid object")
	}
	guid := strings.TrimPrefix(tok[0], "{")
	lcid, err := strconv.Atoi(tok2[0])
	if err != nil {
		return nil, err
	}
	obj := &Object{
		Guid:    strings.TrimSuffix(guid, "}"),
		Version: tok[1],
		LCID:    lcid,
		Path:    strings.TrimSpace(tok2[1]),
	}
	return obj, nil
}

func parseModule(s string) (*Module, error) {
	tok := strings.Split(s, ";")
	if len(tok) != 2 {
		return nil, errors.New("invalid module")
	}
	return &Module{
		Name:     tok[0],
		Filename: strings.TrimSpace(tok[1]),
	}, nil
}
