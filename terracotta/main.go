package main

import (
    "html/template"
    "log"
    "net/http"
    "sync"
)

type Post struct {
    Username string
    Content  string
}

var (
    posts []Post
    mu    sync.Mutex
    tmpl  *template.Template
)

func main() {
    tmpl = template.Must(template.ParseGlob("templates/*.html"))

    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/post", postHandler)

    log.Println("Starting server on :8081...")
    log.Fatal(http.ListenAndServe(":8081", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    tmpl.ExecuteTemplate(w, "layout.html", posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    username := r.FormValue("username")
    content := r.FormValue("content")

    mu.Lock()
    posts = append([]Post{{Username: username, Content: content}}, posts...) // newest first
    mu.Unlock()

    http.Redirect(w, r, "/", http.StatusSeeOther)
}
