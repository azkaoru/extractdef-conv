# extractdef-conv

extractdef-convは、下記のファイルをconvertする。

https://github.com/db-syntax-diff/db_syntax_diff/blob/master/config/extractdef.xml


{1:20}(extractdef-decoder/extractdef.xml)



``` xml (extractdef-decoder/extractdef.xml)
 <PATTERNDEF patid="SQL-102" type="SQL">
    <NAME>CREATE\s+(?:UNIQUE\s+|BITMAP\s+)?INDEX</NAME>
    <PATTERN>CREATE(\s+.*\s+|\s)INDEX ($2) (ON .*)</PATTERN>
    <KEYWORDDEF pos="$1 ">
      <NAME>BITMAP</NAME>
      <PATTERN>BITMAP</PATTERN>
      <MESSAGE id="SQL-102-001" level="LOW1">%keyword%句は未サポートです。</MESSAGE>
    </KEYWORDDEF>
    <KEYWORDDEF pos="$2 ">
      <NAME>schema.index</NAME>
      <PATTERN>[^.(), ]+\s*\.\s*[^.(), ]+</PATTERN>
      <MESSAGE id="SQL-102-002" level="LOW2">インデックス名にスキーマ名は指定できません。</MESSAGE>
    </KEYWORDDEF>
    <KEYWORDDEF pos="$3 ">
      <NAME>CLUSTER</NAME>
      <PATTERN>CLUSTER</PATTERN>
      <MESSAGE id="SQL-102-003" level="ERROR LV2">%keyword%句は未サポートです。</MESSAGE>
    </KEYWORDDEF>
 </PATTERNDEF>
```


 - 上記の検索条件があるか対応可能か？
   確認する必要がある。

