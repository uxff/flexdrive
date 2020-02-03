package wordcode


{{range .Tables}}
{{$table := .}}
var table_{{.Name}} = `
{{.Name}}
字段名	类型	长度	备注	外键
{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}{{Mapper $col.Name}}	{{Type $col}}	{{$col.Length}}	{{$col.Comment}}
{{end}}
`

{{end}}
