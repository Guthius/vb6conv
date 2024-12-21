package vb6

import (
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// These are a bunch of helper functions that are used to extract properties from the VB6 controls.

func GetStr(key string, props PropertyMap) (string, bool) {
	prop, ok := props[key]
	if !ok {
		return "", false
	}

	str, err := strconv.Unquote(prop.Value)
	if err != nil {
		return "", false
	}

	return str, true
}

func GetProp(key string, props PropertyMap) (string, bool) {
	prop, ok := props[key]
	if !ok {
		return "", false
	}

	str := prop.Value

	q := strings.Index(str, "'")
	if q != -1 {
		str = str[:q]
	}

	return str, true
}

func TwipsToPixels(twips int) int {
	inches := float64(twips) / 1440.0
	return int(inches * 96.0)
}

func GetInt(key string, props PropertyMap) (int, bool) {
	str, ok := GetProp(key, props)
	if !ok {
		return 0, false
	}

	str = strings.TrimSpace(str)
	v, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, false
	}

	return int(v), true
}

func GetFloat32(key string, props PropertyMap) (float32, bool) {
	str, ok := GetProp(key, props)
	if !ok {
		return 0, false
	}

	str = strings.TrimSpace(str)
	v, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, false
	}

	return float32(v), true
}

func GetBool(key string, props PropertyMap) (bool, bool) {
	v, ok := GetInt(key, props)
	if !ok {
		return false, false
	}

	return v == -1, true
}

func GetTwips(key string, props PropertyMap) (int, bool) {
	str, ok := GetProp(key, props)
	if !ok {
		return 0, false
	}

	v, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}

	return TwipsToPixels(v), true
}

func GetVector2(x string, y string, props PropertyMap) (int, int, bool) {
	left, ok := GetTwips(x, props)
	if !ok {
		return 0, 0, false
	}

	top, ok := GetTwips(y, props)
	if !ok {
		return 0, 0, false
	}

	return left, top, true
}

func GetColor(key string, props PropertyMap) (uint32, bool) {
	str, ok := GetProp(key, props)
	if !ok {
		return 0, false
	}

	str = strings.TrimPrefix(str, "&H")
	str = strings.TrimSuffix(str, "&")

	v, err := strconv.ParseInt(str, 16, 32)
	if err != nil {
		return 0, false
	}

	return uint32(v), true
}

type Font struct {
	Family        string
	Size          float32
	Charset       int
	Weight        int
	Underline     bool
	Italic        bool
	Strikethrough bool
}

func GetFont(key string, props PropertyMap) (*Font, bool) {
	prop, ok := props[key]
	if !ok {
		return nil, false
	}

	font := &Font{}
	font.Family, _ = GetStr("Name", prop.Properties)
	font.Size, _ = GetFloat32("Size", prop.Properties)
	font.Charset, _ = GetInt("Charset", prop.Properties)
	font.Weight, _ = GetInt("Weight", prop.Properties)
	font.Underline, _ = GetBool("Underline", prop.Properties)
	font.Italic, _ = GetBool("Italic", prop.Properties)
	font.Strikethrough, _ = GetBool("Strikethrough", prop.Properties)

	return font, true
}

type resourceInfo struct {
	filename string
	offset   int64
}

type frxBinaryHeader struct {
	U1   uint32
	U2   uint32
	Size uint32
}

type frxListHeader struct {
	Size   uint16
	MaxLen uint16
}

func parseResourceLocator(root string, locator string) (*resourceInfo, error) {
	colon := strings.Index(locator, ":")
	if colon == -1 {
		return nil, nil
	}

	filename, err := strconv.Unquote(locator[:colon])
	if err != nil {
		return nil, err
	}

	if len(filename) == 0 {
		return nil, nil
	}

	offset, err := strconv.ParseInt(locator[colon+1:], 16, 32)
	if err != nil {
		return nil, err
	}

	return &resourceInfo{
		filename: filepath.Join(root, filename),
		offset:   int64(offset),
	}, nil
}

func GetBinaryResource(c *Control, path string) ([]byte, error) {
	res, err := parseResourceLocator(c.Form.Folder, path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(res.filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	_, err = file.Seek(res.offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var header frxBinaryHeader
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	if header.Size == 0 {
		return []byte{}, nil
	}

	data := make([]byte, header.Size)

	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetListResource(c *Control, path string) ([]string, error) {
	res, err := parseResourceLocator(c.Form.Folder, path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(res.filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	_, err = file.Seek(res.offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var header frxListHeader
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	items := make([]string, 0, header.Size)
	itembuf := make([]byte, header.MaxLen)

	var len uint16
	for i := 0; i < int(header.Size); i++ {
		err = binary.Read(file, binary.LittleEndian, &len)
		if err != nil {
			return nil, err
		}
		err = binary.Read(file, binary.LittleEndian, itembuf[:len])
		if err != nil {
			return nil, err
		}
		items = append(items, string(itembuf[:len]))
	}

	return items, nil
}
