package {{.Models}}

{{$ilen := len .Imports}}
{{if gt $ilen 0}}
import (
	{{range .Imports}}"{{.}}"{{end}}
        "github.com/uxff/flexdrive/pkg/dao/base"
)
{{end}}

{{range .Tables}}
type {{Mapper .Name}} struct {
{{$table := .}}
 {{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{Mapper1up $col.Name}}	{{Type $col}} {{Tag $table $col}}
{{end}}
}

func (t {{Mapper .Name}}) TableName() string {
	return "{{.Name}}"
}

func (t *{{Mapper .Name}}) GetById(int id) error {
	_, err := base.GetByCol("id", id, t)
	return err
}

{{end}}
