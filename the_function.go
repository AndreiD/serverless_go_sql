package function

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"context"
	"database/sql"
	"fmt"
	"os"
	uuid "github.com/google/uuid"
	// Import the MySQL SQL driver.
	_ "github.com/go-sql-driver/mysql"

)

var (
	db *sql.DB

	connectionName = os.Getenv("DB_HOST")
	dbUser         = os.Getenv("DB_USER")
	dbPassword     = os.Getenv("DB_PASS")
	dbName         = os.Getenv("DB_DATABASE")
	dsn            = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", dbUser, dbPassword, connectionName, dbName)
)

func init() {
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("could not open db: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("could not ping db: %v", err)
	}
	// Only allow 1 connection to the database to avoid overloading it.
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
}


// TheFunction ...
func TheFunction(w http.ResponseWriter, r *http.Request) {
	randomUserClient := http.Client{
		Timeout: time.Second * 3,
	}

	req, err := http.NewRequest(http.MethodGet, "https://randomuser.me/api/", nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	res, err := randomUserClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var o map[string]interface{}
	err = json.Unmarshal(body, &o)
	if err != nil {
		log.Fatal(err)
		return
	}

	results := o["results"].([]interface{})
	result := results[0].(map[string]interface{})

	result["generator"] = "google-cloud-function"

	stmt, err := db.Prepare("INSERT INTO users(id, name, email) VALUES (?,?,?)")
	if err != nil {
		log.Printf("db.Query: %v", err)
		return
	}
	_, err = stmt.Exec(uuid.New().String(), "Jimmy", "jim@me.com")
	if err != nil {
		log.Printf("db.Query: %v", err)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
