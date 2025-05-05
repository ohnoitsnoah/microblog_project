package main

import(
	"html/template"
	"log"
	"net/http"
	"sync"
)

type Post struct {
	Username string
	Content string
}

var(
	posts []Post
	mu    sync.Mutex
	tmpl  *template.Template
)

func main() {
	tmpl = template.Must(template.Parse.Glob("templates/*.html"))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post" postHandler)

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	content := r.FormValue("content")

	mu.Lock()
	posts = append([]Post{{Username: username, Content: content}}, posts...)
	// Allows for new posts to show up on top
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
