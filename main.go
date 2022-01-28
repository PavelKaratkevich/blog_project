package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "index", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	client, err := ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err.Error())
	}
	defer client.Close()

	insert, err := client.Query("Insert into articles (title, anons, full_text) values ($1, $2, $3)", title, anons, full_text)
	if err != nil {
		log.Printf("Error while inserting data into database: %v", err.Error())
	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func HandleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", save_article)
	http.ListenAndServe(":8080", nil)
}

func main() {
	HandleFunc()
}

func ConnectDB() (*sqlx.DB, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		return nil, err
	}

	// get environment variables
	db_user := os.Getenv("POSTGRES_USER")
	db_pswd := os.Getenv("POSTGRES_PASSWORD")
	db_address := os.Getenv("DB_ADDRESS")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("POSTGRES_DB")

	dataSource := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", db_address, db_port, db_user, db_name, db_pswd)

	client, err := sqlx.Open("postgres", dataSource)
	if err != nil || client == nil {
		log.Fatal("Error while opening DB: ", err.Error())
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		log.Fatalf("Error while connection ping: %s", err.Error())
		return nil, err
	}
	log.Println("Database is running")

	// Reading file with SQL instructions
	res, err := ioutil.ReadFile("instructions.sql")
	if err != nil {
		log.Fatalf("Error while reading file with instructions: %v", err.Error())
		return nil, err
	}
	var schema = string(res)
	client.MustExec(schema)

	return client, nil
}
