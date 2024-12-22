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

type ref struct {
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

// parseRef parses a reference string into a [ref] struct.
//
// The reference string is in the format "filename:offset".
//   - The filename is relative to the search path.
//   - The offset is a hexadecimal number.
//
// Example reference: "frmDeleteAccount.frx":0C81
func parseRef(searchPath string, str string) (*ref, error) {
	colon := strings.Index(str, ":")
	if colon == -1 {
		return nil, errMalformedLocatorMissingColon
	}

	filename, err := strconv.Unquote(str[:colon])
	if err != nil {
		return nil, err
	}

	if len(filename) == 0 {
		return nil, errMalformedLocatorMissingFilename
	}

	offset, err := strconv.ParseInt(str[colon+1:], 16, 32)
	if err != nil {
		return nil, err
	}

	return &ref{
		filename: filepath.Join(searchPath, filename),
		offset:   int64(offset),
	}, nil
}

// LoadBinary loads a binary resource from a FRX file.
func LoadBinary(searchPath string, refStr string) ([]byte, error) {
	ref, err := parseRef(searchPath, refStr)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(ref.filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	_, err = file.Seek(ref.offset, io.SeekStart)
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
func LoadList(searchPath string, refStr string) ([]string, error) {
	ref, err := parseRef(searchPath, refStr)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(ref.filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	_, err = file.Seek(ref.offset, io.SeekStart)
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
