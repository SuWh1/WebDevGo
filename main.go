package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	fmt.Fprint(w, "<h1>Hllo home</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch email me at <a href=\"https://github.com/SuWh1/WebDevGo\">GitHub</a>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	html := `<h1>FAQ</h1>
	<p><b>Q: Is there a free version?</b></p>
	<p>A: Yes! We offer a free trial for 30 days on any paid plans.</p>
	<p><b>Q: What are your support hours?</b></p>
	<p>A: We have support staff answering emails 24/7, though response times may be a bit slower on weekends.</p>
	<p><b>Q: How do I contact support?</b></p> 
	<p>A: Email us - <a href="https://github.com/SuWh1/WebDevGo">GitHub</a></p>`

	fmt.Fprint(w, html)
}

type Router struct {
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.Error(w, "Page not found!", http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println("Server on :3000...")
	http.ListenAndServe(":3000", router)
}
