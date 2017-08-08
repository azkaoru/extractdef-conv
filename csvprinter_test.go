package main

import (
	"os"
	"testing"
)

func TestCsvPrinterWithHeader(t *testing.T) {
	data := []PrintData{
		PrintData{`FUNC`, `FNC-001-001`, `FNC-001-001`, `LOW1`, `ADD_MONTHS`, `ADD_MONTHS\s*\(`, `%pattern%関数は未サポートです。演算子を利用することで対処可能ですが、結果の型に注意が必要です。`, ``, ``},
		PrintData{`SQL`, `SQL-101`, `SQL-101-002`, `ERROR LV2`, `COMMIT`, `COMMIT ($1)`, `%keyword%句は未サポートです。`, `$1:^(?:WORK\s)?\s*FORCE`, `COMMIT FORCE`}}
	printer := &CsvPrinter{w: os.Stdout, headPrint: true}
	printer.Print(data)
	//TODO: assert
}

func TestCsvPrinterWithoutHeader(t *testing.T) {
	data := []PrintData{
		PrintData{`FUNC`, `FNC-001-001`, `FNC-001-001`, `LOW1`, `ADD_MONTHS`, `ADD_MONTHS\s*\(`, `%pattern%関数は未サポートです。演算子を利用することで対処可能ですが、結果の型に注意が必要です。`, ``, ``},
		PrintData{`SQL`, `SQL-101`, `SQL-101-002`, `ERROR LV2`, `COMMIT`, `COMMIT ($1)`, `%keyword%句は未サポートです。`, `$1:^(?:WORK\s)?\s*FORCE`, `COMMIT FORCE`}}
	printer := &CsvPrinter{w: os.Stdout}
	printer.Print(data)
	//TODO: assert
}
