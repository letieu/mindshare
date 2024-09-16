package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Thought struct {
	ID   int
	Text string
}

type Comment struct {
	ID        int
	Text      string
	ThoughtID int
}

var db *sql.DB
var indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
var detailTmpl = template.Must(template.ParseFiles("templates/detail.html"))

func main() {
    fmt.Println("Initializing the database")
	var err error
	db, err = sql.Open("sqlite3", "./thoughts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable()
    fmt.Println("Database initialized")

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/submit", submitHandler)
	r.HandleFunc("/thought/{id}", detailHandler)
	r.HandleFunc("/thought/{id}/comment", commentHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

    fmt.Println("Server started on port 8080")
	log.Fatal(srv.ListenAndServe())
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS thoughts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

    query = `
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT NOT NULL,
        thoughtID INTEGER NOT NULL,
        FOREIGN KEY (thoughtID) REFERENCES thoughts(id)
    );`

    _, err = db.Exec(query)
    if err != nil {
        log.Fatal(err)
    }
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, text FROM thoughts ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var thoughts []Thought
	for rows.Next() {
		var thought Thought
		if err := rows.Scan(&thought.ID, &thought.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		thoughts = append(thoughts, thought)
	}

	indexTmpl.Execute(w, thoughts)
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	thoughtID := vars["id"]

	var thought Thought
	err := db.QueryRow("SELECT id, text FROM thoughts WHERE id = ?", thoughtID).Scan(&thought.ID, &thought.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, text, thoughtID FROM comments WHERE thoughtID = ? ORDER BY id DESC", thoughtID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Text, &comment.ThoughtID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	detailTmpl.Execute(w, struct {
		Thought
		Comments []Comment
	}{
		thought,
		comments,
	})
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

    if r.Method != http.MethodPost {
        http.Redirect(w, r, fmt.Sprintf("/thought/%s", vars["id"]), http.StatusSeeOther)
        return
    }

    text := r.FormValue("text")
    text = strings.TrimSpace(text)
    if text == "" {
        http.Redirect(w, r, fmt.Sprintf("/thought/%s", vars["id"]), http.StatusSeeOther)
        return
    }

    _, err := db.Exec("INSERT INTO comments (text, thoughtID) VALUES (?, ?)", text, vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/thought/%s", vars["id"]), http.StatusSeeOther)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	text := r.FormValue("text")
    text = strings.TrimSpace(text)
	if text == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err := db.Exec("INSERT INTO thoughts (text) VALUES (?)", text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
