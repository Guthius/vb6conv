package resx

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

type Resx interface {
	Add(key string, value any)
	Count() int
	Save(filename string) error
}

func NewResx() Resx {
	return &resxImpl{
		entries: make(map[string]resxElem),
	}
}

type resxElem struct {
	dataType string
	mimeType string
	value    any
}

type resxImpl struct {
	entries map[string]resxElem
}

func (res *resxImpl) Add(key string, value any) {
	switch (value).(type) {
	case []byte:
		res.entries[key] = resxElem{
			dataType: "System.Drawing.Bitmap, System.Drawing",
			mimeType: "application/x-microsoft.net.object.bytearray.base64",
			value:    value,
		}
	}
}

func (res *resxImpl) Count() int {
	return len(res.entries)
}

func (res *resxImpl) Save(name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(`<?xml version="1.0" encoding="utf-8"?>
<root>
  <xsd:schema id="root" xmlns="" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:msdata="urn:schemas-microsoft-com:xml-msdata">
    <xsd:import namespace="http://www.w3.org/XML/1998/namespace" />
    <xsd:element name="root" msdata:IsDataSet="true">
      <xsd:complexType>
        <xsd:choice maxOccurs="unbounded">
          <xsd:element name="metadata">
            <xsd:complexType>
              <xsd:sequence>
                <xsd:element name="value" type="xsd:string" minOccurs="0" />
              </xsd:sequence>
              <xsd:attribute name="name" use="required" type="xsd:string" />
              <xsd:attribute name="type" type="xsd:string" />
              <xsd:attribute name="mimetype" type="xsd:string" />
              <xsd:attribute ref="xml:space" />
            </xsd:complexType>
          </xsd:element>
          <xsd:element name="assembly">
            <xsd:complexType>
              <xsd:attribute name="alias" type="xsd:string" />
              <xsd:attribute name="name" type="xsd:string" />
            </xsd:complexType>
          </xsd:element>
          <xsd:element name="data">
            <xsd:complexType>
              <xsd:sequence>
                <xsd:element name="value" type="xsd:string" minOccurs="0" msdata:Ordinal="1" />
                <xsd:element name="comment" type="xsd:string" minOccurs="0" msdata:Ordinal="2" />
              </xsd:sequence>
              <xsd:attribute name="name" type="xsd:string" use="required" msdata:Ordinal="1" />
              <xsd:attribute name="type" type="xsd:string" msdata:Ordinal="3" />
              <xsd:attribute name="mimetype" type="xsd:string" msdata:Ordinal="4" />
              <xsd:attribute ref="xml:space" />
            </xsd:complexType>
          </xsd:element>
          <xsd:element name="resheader">
            <xsd:complexType>
              <xsd:sequence>
                <xsd:element name="value" type="xsd:string" minOccurs="0" msdata:Ordinal="1" />
              </xsd:sequence>
              <xsd:attribute name="name" type="xsd:string" use="required" />
            </xsd:complexType>
          </xsd:element>
        </xsd:choice>
      </xsd:complexType>
    </xsd:element>
  </xsd:schema>
  <resheader name="resmimetype">
    <value>text/microsoft-resx</value>
  </resheader>
  <resheader name="version">
    <value>2.0</value>
  </resheader>
  <resheader name="reader">
    <value>System.Resources.ResXResourceReader, System.Windows.Forms, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089</value>
  </resheader>
  <resheader name="writer">
    <value>System.Resources.ResXResourceWriter, System.Windows.Forms, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089</value>
  </resheader>
  <assembly alias="System.Drawing" name="System.Drawing, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" />
`)
	if err != nil {
		return err
	}
	for key, elem := range res.entries {
		_, err = file.WriteString(fmt.Sprintf("\t<data name=\"%s\" type=\"%s\" mimetype=\"%s\">\n", key, elem.dataType, elem.mimeType))
		if err != nil {
			return err
		}
		err = writeValue(file, elem.value)
		if err != nil {
			return err
		}
		_, err = file.WriteString("\t</data>\n")
		if err != nil {
			return err
		}
	}
	_, err = file.WriteString("</root>")
	if err != nil {
		return err
	}
	return nil
}

func encodeBase64(data []byte) string {
	base64 := base64.StdEncoding.EncodeToString(data)
	sb := strings.Builder{}
	for i := 0; i < len(base64); i += 80 {
		end := i + 80
		if end > len(base64) {
			end = len(base64)
		}
		sb.WriteString("        ")
		sb.WriteString(base64[i:end])
		sb.WriteString("\n")
	}
	return sb.String()
}

func writeValue(f *os.File, value any) error {
	_, err := f.WriteString("\t\t<value>\n")
	if err != nil {
		return err
	}
	switch v := value.(type) {
	case []byte:
		_, err = f.WriteString(encodeBase64(v))
		if err != nil {
			return err
		}
	}
	_, err = f.WriteString("\t\t</value>\n")
	return err
}
