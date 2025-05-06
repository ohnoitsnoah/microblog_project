package main

import (
	"database/sql"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

var authTmpl = template.Must(template.ParseGlob("templates/*.html"))

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hash)
		if err != nil {
			http.Error(w, "Username may already exist", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	authTmpl.ExecuteTemplate(w, "register.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var hash string
		err := db.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&hash)
		if err != nil {
			http.Error(w, "Invalid login", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: username,
			Path:  "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	authTmpl.ExecuteTemplate(w, "login.html", nil)
}

func getUsername(r *http.Request) string {
	cookie, err := r.Cookie("username")
	if err != nil {
		return ""
	}
	return cookie.Value
}
