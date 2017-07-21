package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Printer interface {
	Print(data interface{})
}

type PrintData struct {
	Ptype    string
	Pid      string
	Msgid    string
	Level    string
	Pname    string
	Ppattern string
	Msg      string
	Kpattern string
	Kname    string
}

type CsvPrinter struct {
	w         io.Writer
	headPrint bool
}

//http://sambaiz.net/article/37/
// Convert interface{} to []interface{}
func (printer *CsvPrinter) toSlice(src interface{}) []interface{} {

	ret := []interface{}{}
	if v := reflect.ValueOf(src); v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			ret = append(ret, v.Index(i).Interface())
		}

	} else {
		ret = append(ret, v.Interface())
	}

	return ret
}

// Generate csv rows including header from interface{} slice or object
func (printer *CsvPrinter) getRows(src interface{}) [][]string {
	s1 := printer.toSlice(src)
	rows := make([][]string, 1)
	for i, d := range s1 {

		if i != 0 || printer.headPrint {
			rows = append(rows, []string{})
		}
		v := reflect.ValueOf(d)

		for j := 0; j < v.NumField(); j++ {

			if printer.headPrint && i == 0 {
				//Header
				colName := strings.ToLower(v.Type().Field(j).Name)
				rows[0] = append(rows[0], colName)
			}
			rows[len(rows)-1] = append(rows[len(rows)-1], fmt.Sprint(v.Field(j).Interface()))

		}
	}
	return rows
}

func (printer *CsvPrinter) Print(data interface{}) {
	rows := printer.getRows(data)
	for _, row := range rows {
		csvRecord := strings.Join(row, ",")
		fmt.Println(csvRecord)
	}
}
