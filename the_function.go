package function

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	_ "github.com/go-sql-driver/mysql" // needed by sqlx
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

// DB is the database reference.
var DB *sqlx.DB

var (
	dbHost        = os.Getenv("DB_HOST")
	dbUser        = os.Getenv("DB_USER")
	dbPassword    = os.Getenv("DB_PASS")
	dbName        = os.Getenv("DB_DATABASE")
	connectionURI = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
)

func init() {
	var err error
	DB, err = sqlx.Open("mysql", connectionURI)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connection to the database OK")
	// Only allow 1 connection to the database to avoid overloading it.
	DB.SetMaxIdleConns(1)
	DB.SetMaxOpenConns(1)
}

// TheFunction is our main function.
func TheFunction(w http.ResponseWriter, r *http.Request) {
	// Generate some fake data...
	// SomeStruct ...
	type SomeStruct struct {
		Email     string `faker:"email"`
		FirstName string `faker:"first_name"`
		LastName  string `faker:"last_name"`
	}

	a := SomeStruct{}
	err := faker.FakeData(&a)
	if err != nil {
		log.Fatal(err)
	}

	query := "INSERT INTO users(id, name, email) VALUES (?,?,?)"
	newID := uuid.New().String()
	DB.MustExec(query, newID, a.FirstName+" "+a.LastName, a.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "created", "id": newID})
}
