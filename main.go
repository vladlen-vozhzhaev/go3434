package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
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

func apiPostHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getAllPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func apiCreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	/*userID, ok := userIdFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}*/
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if input.Title == "" || input.Content == "" {
		http.Error(w, "title and content must be non-empty", http.StatusBadRequest)
	}
	if err := createPost(r.Context(), input.Title, input.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated) // 201
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

/*func editHandler(w http.ResponseWriter, r *http.Request) {
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

}*/

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "register.html", nil)
		return
	}
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, "email or password is empty", http.StatusBadRequest)
		return
	}
	if len(password) < 8 {
		http.Error(w, "password is too short", http.StatusBadRequest)
		return
	}
	hash, err := hashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(r.Context(),
		"INSERT INTO users (email, password_hash) VALUES ($1, $2)",
		email, hash,
	)
	if err != nil {
		http.Error(w, "Пользователь с таким Email уже существует", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "login.html", nil)
		return
	}
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	var userID int
	var hash string
	err := db.QueryRow(r.Context(),
		"SELECT id, password_hash FROM users WHERE email = $1",
		email,
	).Scan(&userID, &hash)
	if err != nil || !checkPassword(hash, password) {
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}
	token, err := createSession(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24, // сутки
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		deleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		MaxAge:   -1, // сутки
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	/*	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		todayRUB, err := getRub("")
		if err != nil {
			panic(err)
		}
		yesterdayRUB, err := getRub(yesterday)
		if err != nil {
			panic(err)
		}
		fmt.Println("Курс на сегодня: %v", todayRUB)
		fmt.Println("Курс на вчера: %v", yesterdayRUB)
		if todayRUB > yesterdayRUB {
			url, err := getGifURL("rich")
			if err != nil {
				panic(err)
			}
			fmt.Println(url)
		} else {
			url, err := getGifURL("broke")
			if err != nil {
				panic(err)
			}
			fmt.Println(url)
		}

	*/
	if err := initDB(); err != nil {
		fmt.Println("init db failed, err:%v\n", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", requireAuth(createHandler))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/api/posts", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			apiPostHandler(writer, request)
			return
		}
		if request.Method == http.MethodPost {
			requireAuth(apiCreatePostHandler)(writer, request)
			return
		}
		http.Error(writer, "Метод не разрешен", http.StatusMethodNotAllowed) // 405
	})
	fmt.Println("Listening on port http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
