package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	formatPtr := flag.String("format", "html", "Format to output (html, md)")
	outPtr := flag.String("out", "er.html", "Output file name")
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		path = "./ent/schema"
	}

	graph, err := entc.LoadGraph(path, &gen.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	var b bytes.Buffer

	tmpl := htmlTmpl
	out := "er.html"
	if *formatPtr == "md" {
		tmpl = mdTmpl
		out = "er.md"
	}
	if *outPtr != "" {
		out = *outPtr
	}

	if err := tmpl.Execute(&b, graph); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(out, b.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}

var mdTmpl = template.Must(template.New("er").
	Funcs(template.FuncMap{
		"fmtType": func(s string) string {
			return strings.NewReplacer(
				".", "DOT",
				"*", "STAR",
				"[", "LBRACK",
				"]", "RBRACK",
			).Replace(s)
		},
	}).
	Parse(`{{- with $.Nodes }}
` + "```" + `mermaid
erDiagram
{{- range $n := . }}
    {{ $n.Name }} {
	{{- if $n.HasOneFieldID }}
        {{ fmtType $n.ID.Type.String }} {{ $n.ID.Name }}
	{{- end }}
	{{- range $f := $n.Fields }}
        {{ fmtType $f.Type.String }} {{ $f.Name }}
	{{- end }}
    }
{{- end }}
{{- range $n := . }}
    {{- range $e := $n.Edges }}
	{{- if not $e.IsInverse }}
		{{- $rt := "|o--o|" }}{{ if $e.O2M }}{{ $rt = "|o--o{" }}{{ else if $e.M2O }}{{ $rt = "}o--o|" }}{{ else if $e.M2M }}{{ $rt = "}o--o{" }}{{ end }}
    	{{ $n.Name }} {{ $rt }} {{ $e.Type.Name }} : "{{ $e.Name }}{{ with $e.Ref }}/{{ .Name }}{{ end }}"
	{{- end }}
	{{- end }}
{{- end }}
` + "```" + `
{{- end }}`))

var htmlTmpl = template.Must(template.New("er").
	Funcs(template.FuncMap{
		"fmtType": func(s string) string {
			return strings.NewReplacer(
				".", "DOT",
				"*", "STAR",
				"[", "LBRACK",
				"]", "RBRACK",
			).Replace(s)
		},
	}).
	Parse(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
</head>
<body>
	{{- with $.Nodes }}
		<div class="mermaid" id="er-diagram">
erDiagram
{{- range $n := . }}
    {{ $n.Name }} {
	{{- if $n.HasOneFieldID }}
        {{ fmtType $n.ID.Type.String }} {{ $n.ID.Name }}
	{{- end }}
	{{- range $f := $n.Fields }}
        {{ fmtType $f.Type.String }} {{ $f.Name }}
	{{- end }}
    }
{{- end }}
{{- range $n := . }}
    {{- range $e := $n.Edges }}
	{{- if not $e.IsInverse }}
		{{- $rt := "|o--o|" }}{{ if $e.O2M }}{{ $rt = "|o--o{" }}{{ else if $e.M2O }}{{ $rt = "}o--o|" }}{{ else if $e.M2M }}{{ $rt = "}o--o{" }}{{ end }}
    	{{ $n.Name }} {{ $rt }} {{ $e.Type.Name }} : "{{ $e.Name }}{{ with $e.Ref }}/{{ .Name }}{{ end }}"
	{{- end }}
	{{- end }}
{{- end }}
		</div>
	{{- end }}
	<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
	<script src="https://unpkg.com/panzoom@9.4.3/dist/panzoom.min.js"></script>
	<script>
		mermaid.mermaidAPI.initialize({
			startOnLoad: true,
		});
		var observer = new MutationObserver((event) => {
			document.querySelectorAll('text[id^=text-entity]').forEach(text => {
				text.textContent = text.textContent.replace('DOT', '.');
				text.textContent = text.textContent.replace('STAR', '*');
				text.textContent = text.textContent.replace('LBRACK', '[');
				text.textContent = text.textContent.replace('RBRACK', ']');
			});
			observer.disconnect();
			panzoom(document.getElementById('er-diagram'));
		});
		observer.observe(document.getElementById('er-diagram'), { attributes: true, childList: true });
	</script>
</body>
</html>
`))
