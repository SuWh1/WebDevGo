package controllers

import ( // we have text/template but when we work with html we should use that
	"html/template" // to parse and execute HTML safely
	"net/http"
)

// HandlerFunc -> adapter to allow the use of ordinary functions as HTTP handlers
func StaticHandler(tpl Template) http.HandlerFunc { // handle static pages: home contact and so on
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil) // when executed -> redered in client's page
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML // to make links work -> for safe only
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! We offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails 24/7, though response times may be a bit slower on weekends.",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="https://github.com/SuWh1/WebDevGo">GitHub</a>`,
		},
		{
			Question: "Do you have office?",
			Answer:   "Our team is remote.",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
