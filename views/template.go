package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/SuWh1/WebDevGo/context"
	"github.com/SuWh1/WebDevGo/models"
	"github.com/gorilla/csrf"
)

// it specifically sends error string to public users
type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(path.Base(patterns[0])) // when we use galleries/new.gohtml for example, it should take whole path not just new one -> so path.Base used
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser not implemented")
			},
			"errors": func() []string {
				return nil
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsfsing template %w", err)
	}

	return Template{
		htmlTpl: tpl,
	}, nil
}

type Template struct {
	htmlTpl *template.Template
}

type Data struct {
	Yield interface{}
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Execution error: %v", err)
		http.Error(w, "Execution error", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "SMTH went wrong.")
		}
	}
	return msgs
}
