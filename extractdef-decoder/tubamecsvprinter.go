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

func (printer *TubameCsvPrinter) createKnowledge(data interface{}, cate string) string {
	printdata := data.(PrintData)
	if cateName, ok := printer.knowledgeMap[printdata.Pid]; ok {
		return cateName
	} else {
		t := TubamePrintData{}
		t.No = strconv.Itoa(printer.counter)
		printer.counter++
		t.ChapterName = ""
		t.CategoryName = printdata.Pname
		t.KnowledgeTitle = "knowledge" + t.No
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

func (printer *TubameCsvPrinter) createKnowledgeAndCheckItem(data interface{}, cate string) {
	printdata := data.(PrintData)
	t := TubamePrintData{}
	t.No = strconv.Itoa(printer.counter)
	printer.counter++
	t.ChapterName = ""
	t.CategoryName = printdata.Pname + " (" + printdata.Msgid + ")"
	t.KnowledgeTitle = "knowledge" + t.No

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
	t.CheckItemName = "item" + strconv.Itoa(printer.checkItemCounter)
	printer.checkItemCounter++
	//検索手順
	t.CheckItemSeachProcedure = "TODO:"
	//検索実施
	t.CheckItemIsSearch = "True"
	//移植要因
	t.CheckItemPortingFactor = "Java バージョンアップによる変更"
	//難易度
	t.ChekcItemPortingLevel = printdata.Level
	//難易度詳細
	t.CheckItemDifficultyDetail = "TODO:"
	//目視確認
	t.ChekcItemManualCheckFlg = "check"
	//ヒアリング確認
	t.ChekcItemHearingCheckFlg = ""

	//検索条件
	t.SearchTarget = "TODO:"
	t.SearchKey1 = printdata.Ppattern
	if printdata.Pid == printdata.Msgid {
		t.SearchKey2 = ""
	} else {
		t.SearchKey2 = printdata.Kpattern
	}
	t.SearchModule = "TODO:"

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
	t.Survey = ""

	printer.tubamePrintDataList = append(printer.tubamePrintDataList, t)

}

//map
func (printer *TubameCsvPrinter) print(data interface{}) {
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
	printer.CsvPrinter.print(printer.tubamePrintDataList)

}
