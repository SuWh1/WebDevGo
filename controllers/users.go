package controllers // responsible for handling user interactions,
// such as HTTP requests, and orchestrating
// the flow of data between the views and the models.

import (
	"fmt"
	"net/http"

	"github.com/SuWh1/WebDevGo/models"
)

type Users struct { // users can access to various pages that is why we use templates inside the users, that is where we store data
	Templates struct { // for rendering "new" signup page
		New    Template
		SignIn Template // it should assign the signin.gohtml
	}

	UserService *models.UserService // access to UserService to work with users
}

func (u Users) New(w http.ResponseWriter, r *http.Request) { // signup page
	var data struct {
		Email string // place where data of email will be stored
	} // capture data to pass it next

	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data) // just rendering
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email") // in signup.gohtml we use input name, then we use it there
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password) // creating and storing user credentials in database
	if err != nil {
		fmt.Println(err)
		http.Error(w, "SMTH went wrong", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User created: %+v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) { // it is almost the same as sign in (it is obvious)
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data) // just rendering
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) { // processing sign in form
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "SMTH went wrong", http.StatusInternalServerError)
		return
	}

	// cookie
	cookie := http.Cookie{
		Name:     "email",
		Value:    user.Email,
		Path:     "/",  // what path have access to cookie
		HttpOnly: true, // saying dont allow cookies allow work to JS
	}
	http.SetCookie(w, &cookie) // it is silent (does not return error)

	//cookie created
	fmt.Fprintf(w, "User auth: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) { // handle cookie
	email, err := r.Cookie("email")
	if err != nil {
		fmt.Fprint(w, "Email cookie connot be read.")
		return
	}
	fmt.Fprintf(w, "Email cookie: %s\n", email.Value)
	fmt.Fprintf(w, "Headers: %+v\n", r.Header)
}
