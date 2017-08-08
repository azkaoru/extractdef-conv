package main

import (
	"os"
	"testing"
)

func TestTubameCsvPrinter(t *testing.T) {

	data := []PrintData{
		PrintData{`EMBSQL`, `EMB-056-001`, `EMB-056-001`, `ERROR LV4`, `OBJECT DELETE`, `OBJECT\s+DELETE`, `%pattern%文は未サポートです。`, ``, ``},
		PrintData{`EMBSQL`, `EMB-019`, `EMB-019-001`, `ERROR LV2`, `DEALLOCATE DESCRIPTOR`, `DEALLOCATE\s+DESCRIPTOR ($1)`, `%keyword%句は未サポートです。`, `$1 :GLOBAL`, `DEALLOCATE DESCRIPTOR GLOBAL`},
		PrintData{`EMBSQL`, `EMB-019`, `EMB-019-002`, `ERROR LV2`, `DEALLOCATE DESCRIPTOR`, `DEALLOCATE\s+DESCRIPTOR ($1)`, `%keyword%句は未サポートです。`, `$1 :LOCAL`, `DEALLOCATE DESCRIPTOR LOCAL`},
		PrintData{`EMBSQL`, `EMB-003-001`, `EMB-003-001`, `ERROR LV4`, `CACHE FREE ALL`, `CACHE\s+FREE\s+ALL`, `%pattern%文は未サポートです。`, ``, ``}}

	printer := &TubameCsvPrinter{
		w:                   os.Stdout,
		knowledgeMap:        make(map[string]string),
		counter:             1,
		checkItemCounter:    1,
		tubamePrintDataList: []TubamePrintData{}}

	printer.Print(data)

}
