# extractdef-conv

extractdef-convは、db-syntax-diffにあるextractdef.xmlをTUBAME(https://github.com/TUBAME/migration-tool)ナレッジ用のXMLにconvertするプロジェクトです。

 * db-syntax-diff
  https://github.com/db-syntax-diff/db_syntax_diff/blob/master/config/extractdef.xml

extractdef-conv内に同梱しているextractdef.xmlを対象に処理を実施しています。
extractdef-convは処理の対象のextractdef.xmlを外部指定できないです。

## build & run

```
go build
extractdef-conv.exe -print tubame > oracleToPostgres.csv
```

## convert tubame knowledge

 * oracleToPostgres.csvをunicodeに変換してください。
 
 * unicodeに変換したoracleToPostgres.csvをconvツール(https://github.com/tak7iji/conv)を利用して、TUBAMEナレッジのXMLに変換する。
 
```
conv -f oracleToPostgres.csv
Convert oracleToPostgres.csv to TUBAME Knowledge XML.
```

作成されたoracleToPostgres.xmlがTUBAMEナレッジXMLです。

  




