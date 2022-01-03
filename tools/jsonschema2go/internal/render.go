package internal

import (
	"embed"
	"html/template"
	"io"

	js "github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed type.tpl

var statis embed.FS
var structTpl *template.Template

func init() {
	typeTpl, _ := statis.ReadFile("type.tpl")
	structTpl = template.Must(template.New("").Parse(string(typeTpl)))
}

type Render_Type struct {
	Description string
	Name        string
	Fields      []Render_TypeField
}

type Render_TypeField struct {
	Name        string
	Description string
	Type        string
}

func RenderType(out io.Writer, in *js.Schema) error {

	// default value
	model := &Render_Type{
		Description: in.Description,
		Fields:      []Render_TypeField{},
	}

	// TODO

	return structTpl.Execute(out, model)
}
