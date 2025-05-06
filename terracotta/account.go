package main

import (
	"database/sql"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "register.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		http.Error(w, "Username may already be taken", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "login.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var passwordHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&passwordHash)
	if err != nil {
		http.Error(w, "Invalid login", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid login", http.StatusUnauthorized)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getUsername(r *http.Request) string {
	cookie, err := r.Cookie("username")
	if err != nil {
		return ""
	}
	return cookie.Value
}
