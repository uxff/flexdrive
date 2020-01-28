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

func (t *{{Mapper .Name}}) UpdateById(cols []string) error {
	_, err := base.UpdateByCol("id", t.Id, t, cols)
	return err
}

func Get{{Mapper .Name}}ById(id int) (*{{Mapper .Name}}, error) {
	e := &{{Mapper .Name}}{}
	exist, err := base.GetByCol("id", id, e)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	return e, err
}

{{end}}
