package controller

import (
	"html/template"
	"net/http"
	"os"

	"github.com/kokweikhong/go-playground/chi-middleware-jwt-auth/middleware"
)

// create base controller interface
type BaseController interface {
	Index(w http.ResponseWriter, r *http.Request)
	UserForm(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

// create base controller struct
type baseController struct{}

// create base controller constructor
func NewBaseController() BaseController {
	return &baseController{}
}

// create base controller index method
func (c *baseController) Index(w http.ResponseWriter, r *http.Request) {
	// get user from cookie
	cookie, err := r.Cookie("user")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get current working dir with error handling
	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse template with error handling
	tmpl, err := template.ParseFiles(cwd + "/chi-middleware-jwt-auth/frontend/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute template with error handling
	err = tmpl.Execute(w, cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// create base controller login method
func (c *baseController) UserForm(w http.ResponseWriter, r *http.Request) {

	// check if user is already logged in and token is valid
	cookie, err := r.Cookie("token")
	if err == nil {
		tokenString := cookie.Value
		isValid, err := middleware.ValidateJWTToken(tokenString)
		if err == nil && isValid {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// get current working dir with error handling
	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse template with error handling
	tmpl, err := template.ParseFiles(cwd + "/chi-middleware-jwt-auth/frontend/templates/user-form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute template with error handling
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *baseController) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// get username and password from form
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	tokenString, err := middleware.NewJWTClaims(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
	})

	// set user cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "user",
		Value: username,
		Path:  "/",
	})

	// redirect to index
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Logout method
func (c *baseController) Logout(w http.ResponseWriter, r *http.Request) {
	// delete token cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// delete user cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "user",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// redirect to index
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
