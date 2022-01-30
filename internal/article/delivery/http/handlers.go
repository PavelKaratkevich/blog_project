package handlers

import (
	"fmt"
	"log"
	"net/http"
	"new/internal/domain"
	"text/template"
	"new/internal/article/repository/postgresDB"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var posts []domain.Article
var showPost domain.Article

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	client, err := postgresRepository.ConnectDB()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err.Error())
	}
	defer client.Close()

	res, err := client.Query("SELECT * FROM articles")
	if err != nil {
		log.Printf("Error whle selecting all articles: %v", err.Error())
	}

	posts = []domain.Article{}
	for res.Next() {
		var article domain.Article
		err := res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)
		if err != nil {
			log.Printf("Error whle scanning all articles: %v", err.Error())
		}
		posts = append(posts, article)
		// fmt.Println(fmt.Sprintf("Article: %v with id %v", article.Title, article.Id))
	}

	t.ExecuteTemplate(w, "index", posts)
}

func Create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}

func Save_article(w http.ResponseWriter, r *http.Request) {
	var client *sqlx.DB
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	}

	client, err := postgresRepository.ConnectDB()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err.Error())
	}
	defer client.Close()

	insert, err := client.Query("Insert into articles (title, anons, full_text) values ($1, $2, $3)", title, anons, full_text)
	if err != nil {
		log.Printf("Error while inserting data into database: %v", err.Error())
	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		log.Printf("Error while parsing HTML files: %v", err.Error())
	}

	client, err := postgresRepository.ConnectDB()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err.Error())
	}
	defer client.Close()

	res, err := client.Query("SELECT * FROM articles WHERE id = $1", vars["id"])
	if err != nil {
		log.Printf("Error whle selecting all articles: %v", err.Error())
	}

	showPost = domain.Article{}
	for res.Next() {
		var article domain.Article
		err := res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)
		if err != nil {
			log.Printf("Error while scanning all articles: %v", err.Error())
		}
		showPost = article
	}

	t.ExecuteTemplate(w, "show", showPost)
}