package frx

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type locator struct {
	filename string
	offset   int64
}

type binaryHeader struct {
	U1   uint32
	U2   uint32
	Size uint32
}

type listHeader struct {
	Size   uint16
	MaxLen uint16
}

var (
	errMalformedLocatorMissingColon    = errors.New("malformed locator: missing colon")
	errMalformedLocatorMissingFilename = errors.New("malformed locator: missing filename")
)

func parseLocator(searchPath string, s string) (*locator, error) {
	colon := strings.Index(s, ":")
	if colon == -1 {
		return nil, errMalformedLocatorMissingColon
	}

	filename, err := strconv.Unquote(s[:colon])
	if err != nil {
		return nil, err
	}

	if len(filename) == 0 {
		return nil, errMalformedLocatorMissingFilename
	}

	offset, err := strconv.ParseInt(s[colon+1:], 16, 32)
	if err != nil {
		return nil, err
	}

	return &locator{
		filename: filepath.Join(searchPath, filename),
		offset:   int64(offset),
	}, nil
}

// LoadBinary loads a binary resource from a FRX file.
func LoadBinary(searchPath string, path string) ([]byte, error) {
	res, err := parseLocator(searchPath, path)
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

	var header binaryHeader
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

// LoadList loads a list of strings from a FRX file.
func LoadList(searchPath string, path string) ([]string, error) {
	res, err := parseLocator(searchPath, path)
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

	var header listHeader
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
