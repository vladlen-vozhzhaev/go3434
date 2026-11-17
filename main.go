package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
var posts []Post

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello World</h1>")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getAllPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Posts []Post
	}{
		Posts: posts,
	}
	templates.ExecuteTemplate(w, "index.html", data)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "create.html", nil)
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		title := r.FormValue("title")
		content := r.FormValue("content")
		if title == "" || content == "" {
			http.Error(w, "title or content is empty", http.StatusBadRequest)
			return
		}
		if err := createPost(r.Context(), title, content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusInternalServerError)
	}

	if r.Method == http.MethodGet {

	}

	if r.Method == http.MethodPost {

	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	if err := initDB(); err != nil {
		fmt.Println("init db failed, err:%v\n", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	fmt.Println("Listening on port http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
