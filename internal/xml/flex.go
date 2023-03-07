package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"golang.org/x/net/html/charset"
)

/*
<T3FlexForms>
    <data>
        <sheet index="sDEF">
            <language index="lDEF">
                <field index="field_content">
                    <value index="vDEF">705</value>
                </field>
                <field index="field_errors">
                    <value index="vDEF"></value>
                </field>
            </language>
        </sheet>
    </data>
</T3FlexForms>
*/

type T3FlexFormsTag struct {
	XMLName xml.Name `xml:"T3FlexForms"`
	Meta    MetaTag  `xml:"meta"`
	Data    DataTag  `xml:"data"`
}

type MetaTag struct {
	XMLName        xml.Name `xml:"meta"`
	CurrentSheetId string   `xml:"current_sheet_id"`
}

type DataTag struct {
	XMLName xml.Name `xml:"data"`
	Sheet   SheetTag `xml:"sheet"`
}

type SheetTag struct {
	XMLName  xml.Name    `xml:"sheet"`
	Language LanguageTag `xml:"language"`
}

type LanguageTag struct {
	XMLName xml.Name   `xml:"language"`
	Field   []FieldTag `xml:"field"`
}

type FieldTag struct {
	XMLName xml.Name `xml:"field"`
	Index   string   `xml:"index,attr"`
	Value   string   `xml:"value"`
}

// ParseFlexForm parses TV flex form data into a struct
func ParseFlexForm(flexFormData []byte) (*T3FlexFormsTag, error) {
	reader := bytes.NewReader(flexFormData)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	flexForm := &T3FlexFormsTag{}
	err := decoder.Decode(&flexForm)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal flex form XML: %w", err)
	}
	return flexForm, nil
}
