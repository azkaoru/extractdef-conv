package main

import (
	"fmt"
	"testing"
)

func TestPatterndefUnmarshal(t *testing.T) {
	data := `
	  <PATTERNDEF patid="SQL-034" type="SQL" filter="Oracle8">
    	<NAME>ALTER SNAPSHOT</NAME>
    	<PATTERN>ALTER SNAPSHOT ($1)</PATTERN>
    	<KEYWORDDEF pos="$1 ">
      		<NAME>ALTER SNAPSHOT LOG</NAME>
      		<PATTERN>LOG</PATTERN>
      		<MESSAGE id="SQL-035-001" level="ERROR LV4">%pattern%文は未サポートです。</MESSAGE>
    	</KEYWORDDEF>
    	<KEYWORDDEF pos="$1 ">
      		<NAME>ALTER SNAPSHOT</NAME>
      		<PATTERN>!LOG</PATTERN>
      		<MESSAGE id="SQL-034-001" level="ERROR LV4">%pattern%文は未サポートです。</MESSAGE>
    	</KEYWORDDEF>
 	 </PATTERNDEF>
	`

	pdef, err := patterndefUnmarshal(data)
	if err != nil {
		t.Error(err)
	}

	if pdef.Name != "ALTER SNAPSHOT" {
		t.Error("Wrong name ,was expecting 'ALTER SNAPSHOT' but got ", pdef.Name)
	}

	if pdef.Pattern != "ALTER SNAPSHOT ($1)" {
		t.Error("Wrong name ,was expecting 'ALTER SNAPSHOT ($1)' but got ", pdef.Pattern)
	}

	if len(pdef.Keywords) != 2 {
		t.Error("Wrong Keywords ,was expecting 'Keywords size 2' but got ", len(pdef.Keywords))
	}

	for i, keyword := range pdef.Keywords {
		fmt.Println(i, keyword.Msg.Id)
		fmt.Println(i, keyword.Msg.Level)
		fmt.Println(i, keyword.Msg.Value)
	}
}

func TestExtractdefUnmarshal(t *testing.T) {

	data := `
    <?xml version="1.0" encoding="UTF-8" ?>
        <!-- Copyright (C) 2007-2010 NTT -->
        <DEFINITION>
            <COMMON>
                <MACRO>
                    <NAME>%todb%</NAME>
                    <VALUE>PostgreSQL</VALUE>
                </MACRO>
                <MACRO>
                    <NAME>%fromdb%</NAME>
                    <VALUE>Oracle</VALUE>
                </MACRO>
            </COMMON>
            <PATTERNDEF patid="SQL-034" type="SQL" filter="Oracle8">
                <NAME>ALTER SNAPSHOT</NAME>
                <PATTERN>ALTER SNAPSHOT ($1)</PATTERN>
                <KEYWORDDEF pos="$1 ">
                    <NAME>ALTER SNAPSHOT LOG</NAME>
                    <PATTERN>LOG</PATTERN>
                    <MESSAGE id="SQL-035-001" level="ERROR LV4">%pattern%文は未サポートです。</MESSAGE>
                </KEYWORDDEF>
                <KEYWORDDEF pos="$1 ">
                    <NAME>ALTER SNAPSHOT</NAME>
                    <PATTERN>!LOG</PATTERN>
                    <MESSAGE id="SQL-034-001" level="ERROR LV4">%pattern%文は未サポートです。</MESSAGE>
                </KEYWORDDEF>
            </PATTERNDEF>
        </DEFINITION>
	`
	xdef, err := extrafdefUnmarshal(data)
	if err != nil {
		t.Error(err)
	}
	if len(xdef.Patterndefs) != 1 {
		t.Error("Wrong Patterndef ,was expecting 'Patterndef size 1' but got ", len(xdef.Patterndefs))
	}

}
