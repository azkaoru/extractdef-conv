package main

import (
	"encoding/xml"
	"flag"
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

//grouping patterndef.Type
func typeGrouping(printdataList []PrintData) map[string][]PrintData {

	printDataMap := make(map[string][]PrintData)

	for _, printData := range printdataList {

		var newList []PrintData
		if list, ok := printDataMap[printData.Ptype]; ok {
			newList = append(list, printData)
		} else {
			newList = append([]PrintData{}, printData)
		}
		printDataMap[printData.Ptype] = newList

	}

	return printDataMap
}

func newPrinter(printType string) Printer {
	var printer Printer
	switch printType {
	case DEFAULT_CSV:
		printer = &CsvPrinter{w: os.Stdout, headPrint: false}
	case TUBAME_CSV:
		printer = &TubameCsvPrinter{
			w:                os.Stdout,
			knowledgeMap:     make(map[string]string),
			counter:          1,
			checkItemCounter: 1,
		}
	}
	return printer
}

const DEFAULT_CSV = "csv"
const TUBAME_CSV = "tubame"

func main() {

	var printType string
	flag.StringVar(&printType, "print", DEFAULT_CSV, "if you want to convert for tubame knowledge converter,use tubame")
	flag.Parse()
	xdef, err := decode("extractdef.xml")
	if err != nil {
		fmt.Println("Error:", err)
	}
	//fmt.Println("extractdef.xml ,PATTERNDEF size = ", len(xdef.Patterndefs))

	var data []PrintData

	for _, patterndef := range xdef.Patterndefs {
		if len(patterndef.Keywords) == 0 {
			printdata := PrintData{
				patterndef.Type,
				patterndef.Id,
				patterndef.Msg.Id,
				patterndef.Msg.Level,
				patterndef.Name,
				patterndef.Pattern,
				strings.Trim(strings.TrimSpace(patterndef.Msg.Value), "\t"),
				``,
				``}
			data = append(data, printdata)
		} else {
			for _, keyword := range patterndef.Keywords {
				printdata := PrintData{
					patterndef.Type,
					patterndef.Id,
					keyword.Msg.Id,
					keyword.Msg.Level,
					patterndef.Name,
					patterndef.Pattern,
					strings.Trim(strings.TrimSpace(keyword.Msg.Value), "\t"),
					keyword.Pos + `:` + keyword.Pattern,
					patterndef.Name + ` ` + keyword.Name}
				data = append(data, printdata)
			}
		}
	}

	printDataMap := typeGrouping(data)
	csvPrinter := newPrinter(printType)
	for _, printdataList := range printDataMap {
		csvPrinter.Print(printdataList)
	}

}
