package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Username string
	Content  string
	TimeAgo  string
}

var db *sql.DB
var templates = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post", postHandler)

	log.Println("Starting server on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT username, content, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var username, content string
		var createdAt time.Time
		err := rows.Scan(&username, &content, &createdAt)
		if err != nil {
			log.Println(err)
			continue
		}

		posts = append(posts, Post{
			Username: username,
			Content:  content,
			TimeAgo:  timeAgo(createdAt),
		})
	}

	templates.ExecuteTemplate(w, "index.html", posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	content := r.FormValue("content")

	if username == "" || content == "" {
		http.Error(w, "Missing fields", 400)
		return
	}

	_, err := db.Exec("INSERT INTO posts (username, content) VALUES (?, ?)", username, content)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func timeAgo(t time.Time) string {
	diff := time.Since(t)
	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	default:
		return t.Format("Jan 2")
	}
}
