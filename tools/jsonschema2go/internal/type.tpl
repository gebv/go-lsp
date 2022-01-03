// {{.Description}}
type {{.Name}} struct {
    {{range $el := .Fields}}
    {{$el.Name}} {{$el.Type}}
    {{end}}
}
