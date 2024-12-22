package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/SuWh1/WebDevGo/controllers"
	"github.com/SuWh1/WebDevGo/migrations"
	"github.com/SuWh1/WebDevGo/models"
	"github.com/SuWh1/WebDevGo/templates"
	"github.com/SuWh1/WebDevGo/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig

	CSRF struct {
		Key    string
		Secure bool
	}

	Server struct {
		Address string // for domain: localhost:3000
	}
}

func loadEnvConfig() (config, error) {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// TODO: Read te PSQL val from an ENV var
	cfg.PSQL = models.DefaultPostgresConfig()

	// TODO: SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}

	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// TODO: Read the CSRF val from an ENV var
	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	// TODO: Read the server val from an ENV var
	cfg.Server.Address = ":3000"

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Database setup
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Service setups
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	galleryService := &models.GalleryService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	// Middleware setup
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// Controllers setup
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS,
		"check-your-email.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS,
		"reset-pw.gohtml", "tailwind.gohtml",
	))
	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleriesC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"galleries/new.gohtml", "tailwind.gohtml",
	))
	galleriesC.Templates.Edit = views.Must(views.ParseFS(
		templates.FS,
		"galleries/edit.gohtml", "tailwind.gohtml",
	))

	// Routers and routes setup
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Route("/galleries", func(r chi.Router) {
		r.Group(func(r chi.Router) { // we group it so that we can use any route that needed signin in there; and outside it can use routes like galleries by id that is public
			r.Use(umw.RequireUser)
			r.Get("/", galleriesC.Index)
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
		})
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Printf("Server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}

// func TimerMiddleware(r *http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		h(w, r)
// 		fmt.Println("Request time:", time.Since(start))
// 	}
// }
