package main

import (
	"io"
	"strconv"
	"strings"
)

type TubamePrintData struct {

	//for knowledge
	No string
	//チャプター名
	ChapterName string
	//カテゴリ名
	CategoryName string
	//親カテゴリ名
	ParentCategoryName string
	//ナレッジタイトル
	KnowledgeTitle string
	//ナレッジ概要
	KnowledgeContent string

	//for checkitem
	//チェックアイテム名
	CheckItemName string
	//検索手順
	CheckItemSeachProcedure string
	//検索実施
	CheckItemIsSearch string
	//移植要因
	CheckItemPortingFactor string
	//難易度
	ChekcItemPortingLevel string
	//難易度詳細
	CheckItemDifficultyDetail string
	//目視確認
	ChekcItemManualCheckFlg string
	//ヒアリング確認
	ChekcItemHearingCheckFlg string

	//check checkCondtion
	SearchTarget string
	SearchKey1   string
	SearchKey2   string
	SearchModule string

	//Calc
	//ライン数算出
	CalcEnable string
	//非算出理由
	CalcNotCountStepDesc string
	//不明/TODO
	CalcResultTodo string
	//ライン数
	CalcStepCount string
	//ライン数根拠
	CalcStepBasis string
	//調査内容
	Survey string
}

type TubameCsvPrinter struct {
	w         io.Writer
	headPrint bool
	CsvPrinter
	chapterCategoryName string
	knowledgeMap        map[string]string
	counter             int
	checkItemCounter    int
	tubamePrintDataList []TubamePrintData
}

//検索対象と、検索対象モジュールを返す
func (printer *TubameCsvPrinter) getSearchTargetAndSearchModule(param1 string) (string, string) {
	var searchTarget, searchModule string
	switch param1 {
	case "EMBSQL":
		searchTarget = "*.p?c|*.cpp|*.h"
		searchModule = "ext_search_sql_parser.py"
	case "FUNC", "SQL", "TYPE":
		searchTarget = "*.sql|*.java|*.ddl|*.p?c|*.cpp|*.h"
		searchModule = "ext_search_sql_parser.py"
	case "ORACA", "SQLCA", "SQLDA":
		searchTarget = "*.p?c|*.h|*.cpp"
		searchModule = ""
	}
	return searchTarget, searchModule
}

//調査内容を返す
func (printer *TubameCsvPrinter) getInvestigation(param1 string, param2 string) string {
	//get $1 from $1 :GLOBAL
	replaceTarget := strings.TrimSpace(strings.Split(param2, ":")[0])

	newval := strings.Replace(param1, replaceTarget, param2, -1)
	return newval
}

func (printer *TubameCsvPrinter) getSearchKey1(param1 string, param2 string) string {
	//get ($1) from $1 :GLOBAL
	replaceTarget := "(" + strings.TrimSpace(strings.Split(param2, ":")[0]) + ")"

	// get Global from $1 :GLOBAL
	replaceValue := strings.TrimSpace(strings.Split(param2, ":")[1])

	newval := strings.Replace(param1, replaceTarget, replaceValue, -1)
	return newval
}

func (printer *TubameCsvPrinter) getPortabilityDegree(level string) (string, string) {
	var val string
	var val2 string
	switch level {
	case "CHECK_LOW2", "LOW1", "LOW2":
		val = "Low"
		if level == "LOW1" {
			val2 = "低1"
		} else if level == "LOW2" || level == "CHECK_LOW2" {
			val2 = "低2"
		}
	case "CHECK_MIDDLE", "MIDDLE":
		val = "Middle"
		val2 = "中"
	case "High":
		val = "High"
		val2 = "高"
	case "ERROR LV1", "ERROR LV2", "ERROR LV3", "ERROR LV4", "ERROR LV5":
		val = "Unknown"
		val2 = "不明1"
	case "WARNING":
		val = "Unknown"
		val2 = "不明2"
	}
	return val, val2
}
func (printer *TubameCsvPrinter) createKnowledge(data interface{}, cate string) string {
	printdata := data.(PrintData)
	if cateName, ok := printer.knowledgeMap[printdata.Pid]; ok {
		return cateName
	} else {
		t := TubamePrintData{}
		t.No = strconv.Itoa(printer.counter)

		if len(printer.tubamePrintDataList) == 0 {
			t.CategoryName = printdata.Ptype + " MIGRATION"
		} else {
			t.CategoryName = printdata.Pname
		}
		printer.counter++
		t.ChapterName = ""

		t.KnowledgeTitle = "knowledge-" + t.No
		t.KnowledgeContent = "<ns3:para>" + t.ChapterName + "</ns3:para>"
		if cate != "" {
			t.ParentCategoryName = cate
			t.ChapterName = ""
		} else {
			t.ChapterName = printdata.Ptype
		}
		printer.tubamePrintDataList = append(printer.tubamePrintDataList, t)
		printer.knowledgeMap[printdata.Pid] = t.CategoryName
		return printer.knowledgeMap[printdata.Pid]
	}
}

func (printer *TubameCsvPrinter) createCheckItemForXml(printData TubamePrintData) TubamePrintData {
	printData.CheckItemName = "item-" + strconv.Itoa(printer.checkItemCounter)
	printData.SearchTarget = "*.xml"
	printData.SearchKey1 = `"//*[contains(translate(text(), 'abcdefghijklmnopqrstuvwxyz','ABCDEFGHIJKLMNOPQRSTUVWXYZ'),'KEY1') and contains(translate(text(), 'abcdefghijklmnopqrstuvwxyz','ABCDEFGHIJKLMNOPQRSTUVWXYZ'),'KEY2')]"`
	printData.KnowledgeContent = ""
	printData.SearchModule = "ext_search_xml_parser.py"
	printer.checkItemCounter++
	return printData
}

func (printer *TubameCsvPrinter) createCheckItemForProperties(printData TubamePrintData) TubamePrintData {
	//チェックアイテム名
	printData.CheckItemName = "item-" + strconv.Itoa(printer.checkItemCounter)
	printData.SearchTarget = "*.properties"
	printData.SearchModule = ""
	printData.KnowledgeContent = ""
	printer.checkItemCounter++
	return printData
}

func (printer *TubameCsvPrinter) createKnowledgeAndCheckItem(data interface{}, cate string) {
	printdata := data.(PrintData)

	if printdata.Level == "FATAL" {
		return
	}
	t := TubamePrintData{}
	t.No = strconv.Itoa(printer.counter)
	printer.counter++
	t.ChapterName = ""
	t.CategoryName = printdata.Pname + " (" + printdata.Msgid + ")"
	t.KnowledgeTitle = "knowledge-" + t.No

	if printdata.Pid == printdata.Msgid {
		t.KnowledgeContent = "<ns3:para>" + printdata.Ppattern + " " + strings.Replace(printdata.Msg, "%", "", -1) + "</ns3:para>"
	} else {
		replaceTarget := strings.TrimSpace(strings.Split(printdata.Kpattern, ":")[0])
		repaceValue := strings.Split(printdata.Kpattern, ":")[1]
		newval := strings.Replace(printdata.Ppattern, replaceTarget, repaceValue, -1)
		t.KnowledgeContent = "<ns3:para>" + newval + " " + strings.Replace(printdata.Msg, "%", "", -1) + "</ns3:para>"
	}

	t.ParentCategoryName = cate

	//チェックアイテム名
	t.CheckItemName = "item-" + strconv.Itoa(printer.checkItemCounter)
	printer.checkItemCounter++
	//検索手順
	t.CheckItemSeachProcedure = "TODO:"
	//検索実施
	t.CheckItemIsSearch = "True"
	//移植要因
	t.CheckItemPortingFactor = "DBMS の変更"

	level, levelDesc := printer.getPortabilityDegree(printdata.Level)
	//難易度
	t.ChekcItemPortingLevel = level
	//難易度詳細
	t.CheckItemDifficultyDetail = levelDesc
	//目視確認
	t.ChekcItemManualCheckFlg = "check"
	//ヒアリング確認
	t.ChekcItemHearingCheckFlg = ""

	searchTarget, searchModule := printer.getSearchTargetAndSearchModule(printdata.Ptype)
	//検索条件
	t.SearchTarget = searchTarget

	if printdata.Pid == printdata.Msgid {
		t.SearchKey1 = printdata.Ppattern
		t.SearchKey2 = ""
		t.Survey = printdata.Ppattern
		t.CheckItemSeachProcedure = printdata.Ppattern
	} else {
		t.SearchKey1 = printer.getSearchKey1(printdata.Ppattern, printdata.Kpattern)
		t.SearchKey2 = ""
		t.Survey = printer.getInvestigation(printdata.Ppattern, printdata.Kpattern)
		t.CheckItemSeachProcedure = t.Survey
	}
	t.SearchModule = searchModule

	//ライン数算出
	t.CalcEnable = "TRUE"
	//非算出理由
	t.CalcNotCountStepDesc = ""
	//不明/TODO
	t.CalcResultTodo = "TODO:SE 手動算出"
	//ライン数
	t.CalcStepCount = ""
	//ライン数根拠
	t.CalcStepBasis = ""
	//調査内容
	//t.Survey = ""

	printer.tubamePrintDataList = append(printer.tubamePrintDataList, t)

	switch printdata.Ptype {

	case "FUNC", "SQL", "TYPE":
		//add check item for xml , properties
		printer.tubamePrintDataList = append(printer.tubamePrintDataList, printer.createCheckItemForXml(t))
		printer.tubamePrintDataList = append(printer.tubamePrintDataList, printer.createCheckItemForProperties(t))
	}

}

//map
func (printer *TubameCsvPrinter) Print(data interface{}) {
	printer.tubamePrintDataList = []TubamePrintData{}
	slice := printer.CsvPrinter.toSlice(data)
	for i, row := range slice {
		if i == 0 {
			printer.chapterCategoryName = printer.createKnowledge(row, "")
			printer.createKnowledgeAndCheckItem(row, printer.chapterCategoryName)
		} else {
			printdata := row.(PrintData)
			if cate, ok := printer.knowledgeMap[printdata.Pid]; ok {
				printer.createKnowledgeAndCheckItem(printdata, cate)
			} else {
				parentCate := printer.createKnowledge(row, printer.chapterCategoryName)
				printer.createKnowledgeAndCheckItem(printdata, parentCate)
			}
		}
	}

	printer.CsvPrinter.Print(printer.tubamePrintDataList)

}
