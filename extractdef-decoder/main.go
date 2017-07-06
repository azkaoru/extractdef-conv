package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Message struct {
	XMLName xml.Name `xml:"MESSAGE"`
	Id      string   `xml:"id,attr"`
	Level   string   `xml:"level,attr"`
	Value   string   `xml:",chardata"`
}

type Keyworddef struct {
	Name    string `xml:"NAME"`
	Pattern string `xml:"PATTERN"`
	Pos     string `xml:"pos,attr"`
	Msg     Message
}

type Patterndef struct {
	XMLName  xml.Name     `xml:"PATTERNDEF"`
	Name     string       `xml:"NAME"`
	Pattern  string       `xml:"PATTERN"`
	Keywords []Keyworddef `xml:"KEYWORDDEF"`
	Id       string       `xml:"patid,attr"`
	Type     string       `xml:"type,attr"`
	Msg      Message
}

type Extractdef struct {
	Name        string       `xml:"DEFINITION"`
	Patterndefs []Patterndef `xml:"PATTERNDEF"`
}

func patterndefUnmarshal(body string) (patterndef Patterndef, err error) {
	err = xml.Unmarshal([]byte(body), &patterndef)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	return
}

func extrafdefUnmarshal(body string) (extractdef Extractdef, err error) {
	err = xml.Unmarshal([]byte(body), &extractdef)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	return
}

func decode(filename string) (extractdef Extractdef, err error) {
	xmlfile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening extractdef.xml:", err)
		return
	}
	defer xmlfile.Close()
	decoder := xml.NewDecoder(xmlfile)
	err = decoder.Decode(&extractdef)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	return
}

func main() {
	xdef, err := decode("extractdef.xml")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("extractdef.xml ,PATTERNDEF size = ", len(xdef.Patterndefs))

	var data []PrintData

	for _, patterndef := range xdef.Patterndefs {
		if len(patterndef.Keywords) == 0 {
			printdata := PrintData{
				patterndef.Id,
				patterndef.Type,
				patterndef.Msg.Id,
				patterndef.Msg.Level,
				patterndef.Name,
				patterndef.Pattern,
				strings.Trim(patterndef.Msg.Value, "\n "),
				``,
				``}
			data = append(data, printdata)
		} else {
			for _, keyword := range patterndef.Keywords {
				printdata := PrintData{
					patterndef.Id,
					patterndef.Type,
					keyword.Msg.Id,
					keyword.Msg.Level,
					patterndef.Name,
					patterndef.Pattern,
					strings.Trim(keyword.Msg.Value, "\n "),
					keyword.Pos + `:` + keyword.Pattern,
					patterndef.Name + ` ` + keyword.Name}
				data = append(data, printdata)
			}
		}
	}

	printer := &CsvPrinter{w: os.Stdout, headPrint: true}
	printer.print(data)

	fmt.Println("extractdef.xml ,check rule size = ", len(data))
}
