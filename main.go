package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

type User struct {
	Name   string
	ID     int
	NUM    float64
	Sports map[string]int
	Skills []string
	Meta   UserMeta
}

type UserMeta struct {
	Visits int
}

func Logger(next http.Handler) http.Handler {
	// return middleware.Logger(next) // just logger with default info

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Incoming request: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	tpl, err := template.ParseFiles(filepath) // we cannot just put path there on windows
	if err != nil {
		log.Printf("parsing error: %v", err)
		http.Error(w, "Parsing error", http.StatusInternalServerError)
		return
	}

	user := User{
		Name: "Dauren",
		ID:   001,
		Sports: map[string]int{
			"Tennis":   3,
			"Football": 2,
			"Swimming": 4,
		},
		Skills: []string{
			"Go",
			"JS",
			"Python",
		},
		Meta: UserMeta{
			Visits: 999,
		},
	}

	err = tpl.Execute(w, user)
	if err != nil {
		log.Printf("Execution error: %v", err)
		http.Error(w, "Execution error", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath)
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "name")

	if userID == "" {
		userID = "Guest"
	}

	w.Write([]byte(fmt.Sprintf(`<h1>USES PAGE</h1>
	<p>The userID is %v</p>`, userID)))
}

func main() {
	r := chi.NewRouter()

	r.Use(Logger)

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)

	// exercise
	r.Get("/users", urlHandler)
	r.Get("/users/{name}", urlHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Server on :3000...")
	http.ListenAndServe(":3000", r)
}
