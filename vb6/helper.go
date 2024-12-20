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

	v, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}

	return v, true
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

func parseString(str string) string {
	// TODO: Should implement this properly
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return ""
	}
	if str[0] == '"' {
		str = str[1:]
	}
	if str[len(str)-1] == '"' {
		str = str[:len(str)-1]
	}
	return str
}

type frxHeader struct {
	U1   uint32
	U2   uint32
	Size uint32
}

func GetResource(c *Control, path string) ([]byte, error) {
	colon := strings.Index(path, ":")
	if colon == -1 {
		return nil, nil
	}

	filename := parseString(path[:colon])
	if len(filename) == 0 {
		return nil, nil
	}

	filePath := filepath.Join(c.Form.Folder, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	offset, err := strconv.ParseInt(path[colon+1:], 16, 32)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	var header frxHeader
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
