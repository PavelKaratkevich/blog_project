package postgresRepository

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sqlx.DB, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file")
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
		log.Printf("Error while opening DB: %v", err.Error())
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		log.Printf("Error while connection ping: %s", err.Error())
		return nil, err
	}
	log.Println("Database is running")

	// Reading file with SQL instructions
	res, err := ioutil.ReadFile("instructions.sql")
	if err != nil {
		log.Printf("Error while reading file with instructions: %v", err.Error())
		return nil, err
	}
	var schema = string(res)
	client.MustExec(schema)

	return client, nil
}
