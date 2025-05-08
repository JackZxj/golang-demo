package main

import (
	"log"
	"os"
	"strings"
	"text/template"
)

type Obj struct {
	Name string `json:"name"`
}

func main() {
	// First we create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"title": strings.Title,
		"test": func(v string) map[string]string {
			return map[string]string{"test": v}
		},
		"test2": func(v string) Obj {
			return Obj{v}
		},
		"test3": func(v string) map[string]interface{} {
			return map[string]interface{}{"test": v}
		},
	}

	// A simple template definition to test our function.
	// We print the input text several ways:
	// - the original
	// - title-cased
	// - title-cased and then printed with %q
	// - printed with %q and then title-cased.
	const templateText = `
Input: {{printf "%q" .}}
Output 0: {{title .}}
Output 1: {{title . | printf "%q"}}
Output 2: {{printf "%q" . | title}}

Output 3: {{ $v := test .}} get {{printf "%q" $v }} v: {{ $v.test }}
Output 3.1: {{- range $k, $v := test .}} get {{printf "%q:%q" $k $v }} {{- end}}
Output 3.2: {{ $v := test3 .}} get {{printf "%q" $v }} v: {{ $v.test }}
Output 4: {{ $v := test2 .}} get {{printf "%q" $v }} v: {{ $v.Name }}
`

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, "the go programming language")
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

}
