package main

import (
	"encoding/xml"
	"fmt"
	"os"
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
	Msg     Message
}

type Patterndef struct {
	XMLName  xml.Name     `xml:"PATTERNDEF"`
	Name     string       `xml:"NAME"`
	Pattern  string       `xml:"PATTERN"`
	Keywords []Keyworddef `xml:"KEYWORDDEF"`
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
}
