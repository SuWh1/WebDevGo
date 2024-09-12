package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

// ParseFS accesses files from virtual file system that implements the fs.FS interface
// retrives the binary template files from templates/fs.go
func ParseFS(fs fs.FS, patterns ...string) (Template, error) { // parses one or more templates from a file system, and returns parsed templates
	tpl := template.New(patterns[0]) // it is not obvious, but we need it to make parsefs see the csrf protection
	tpl = tpl.Funcs(                 // placeholder func
		template.FuncMap{
			"csrfField": func() (template.HTML, error) { // it gives access to csrf func in html templates
				return "", fmt.Errorf("csrfField not implemented")
			},
		},
	)
	// template.ParseFs creates new template to parse everything, but now it is not, now it is using existing templates with csrf
	tpl, err := tpl.ParseFS(fs, patterns...) // there it parsing our tamlates and see the csrf
	if err != nil {
		return Template{}, fmt.Errorf("parsfsing template %w", err)
	}

	return Template{ // return parsed template
		htmlTpl: tpl,
	}, nil
}

type Template struct { // wraps the parsed template
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := t.htmlTpl.Clone() // we need clone to every user have unique instance of web request
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(
		template.FuncMap{ // where we need request specific information, we update them || id we do not need them we just define them in ParseFs
			"csrfField": func() template.HTML { // it gives access to csrf func in html templates
				return csrf.TemplateField(r)
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	var buf bytes.Buffer // is takes quite memory in huge web sitess

	err = tpl.Execute(&buf, data) // we should execute to responce writer, because when csrfField error occur it will make
	// the templates render only till the csrf which is not we want
	// so we execute it ot kind of buffer
	if err != nil {
		log.Printf("Execution error: %v", err)
		http.Error(w, "Execution error", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf) // of all goes well we copy from buffer to responce writer
}
