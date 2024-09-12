package main

import (
	"fmt"
	"net/http"

	"github.com/SuWh1/WebDevGo/controllers"
	"github.com/SuWh1/WebDevGo/models"
	"github.com/SuWh1/WebDevGo/templates"
	"github.com/SuWh1/WebDevGo/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService, // we have access to users controllers
	} // user-related actions: browsing diff pages

	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	)) // signup new user
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	)) // gives signin templates inside user controller

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Server on :3000...")

	csrfKey := "qwertyuioplkjhgfdsazxcvbnmkijhg"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false), // because it is local
	)
	http.ListenAndServe(":3000", csrfMw(r)) // listens for incoming HTTP requests
}
