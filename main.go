package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Person struct {
	Id   int
	Name string
	Age  string
}

var Pasha Person

func main() {
	log.Println("My SQL in docker running")

	client, err := ConnectDB()
	if err != nil {
		log.Fatalf("Error with DB Connection: %v", err.Error())
	}

	get, err := client.Query("Select Id, Name, Age from blog_project")
	if err != nil {
		log.Fatalf("Error while inseting data into database: %v", err.Error())
	}

	for get.Next() {
		var user Person
		err = get.Scan(&user.Id, &user.Name, &user.Age)

		if err != nil {
			log.Fatalf("Error while getting data into database: %v", err.Error())
		}
		log.Printf("User name is %v, id - %v, age - %v", user.Id, user.Name, user.Age)
	} 
	// get.Close()
}

func ConnectDB() (*sqlx.DB, error) {
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
