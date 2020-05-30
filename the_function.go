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
	"strings"
	"time"
)

// DB is the database reference.
var DB *sqlx.DB

var (
	dbHost        = os.Getenv("DB_HOST")
	dbUser        = os.Getenv("DB_USER")
	dbPassword    = os.Getenv("DB_PASS")
	dbName        = os.Getenv("DB_DATABASE")
	connectionURI = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)
)

// SomeStruct ...
type SomeStruct struct {
	ID         string    `db:"id" json:"id"`
	Email      string    `faker:"email" db:"email" json:"email"`
	Name       string    `faker:"first_name" db:"name" json:"name"`
	LastUpdate time.Time `db:"last_update" json:"last_update"`
}

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
	log.Println("initialized OK! connection to the database OK!")
	// Only allow 1 connection to the database to avoid overloading it.
	DB.SetMaxIdleConns(1)
	DB.SetMaxOpenConns(1)
}

// TheFunction is our main function.
func TheFunction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUser(r, w)
	case http.MethodPost:
		createUser(r, w)
	case http.MethodPut:
		updateUser(r, w)
	case http.MethodDelete:
		deleteUser(r, w)
	default:
		respond(http.StatusBadRequest, map[string]interface{}{"error": "unsupported http verb"}, w)
	}
}

// getUser a user
func getUser(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respond(http.StatusBadRequest, map[string]interface{}{"error": "URL Param 'id' is missing"}, w)
		return
	}

	s := SomeStruct{}
	err := DB.Get(&s, "SELECT * FROM users WHERE id = ? LIMIT 1", id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respond(http.StatusNotFound, map[string]interface{}{"error": "no user with this ID was found in the database"}, w)
			return
		}
		respond(http.StatusBadRequest, map[string]interface{}{"error": err.Error()}, w)
		return
	}
	respond(http.StatusOK, s, w)
}

// createUser a user
func createUser(r *http.Request, w http.ResponseWriter) {
	//you can read the json like this
	//decoder := json.NewDecoder(req.Body)
	//var t test_struct
	//err := decoder.Decode(&t)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println(t.Test)
	log.Printf("got a new request with user agent %s", r.UserAgent())

	a := SomeStruct{}
	err := faker.FakeData(&a)
	if err != nil {
		log.Fatal(err)
	}

	newID := uuid.New().String()
	_, err = DB.NamedExec(`INSERT INTO users (id, name, email) VALUES (:id, :name, :email)`,
		map[string]interface{}{
			"id":    newID,
			"name":  a.Name,
			"email": a.Email,
		})

	if err != nil {
		respond(http.StatusBadRequest, map[string]interface{}{"error": err.Error()}, w)
		return
	}
	respond(http.StatusCreated, map[string]interface{}{"status": "created", "id": newID}, w)
}

// updateUser a user
func updateUser(r *http.Request, w http.ResponseWriter) {
	//you can read the json like this
	decoder := json.NewDecoder(r.Body)
	var payload struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	err := decoder.Decode(&payload)
	if err != nil {
		respond(http.StatusBadRequest, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	_, err = DB.NamedExec(`UPDATE users SET name=:name WHERE id = :id`,
		map[string]interface{}{
			"id":   payload.ID,
			"name": payload.Name,
		})
	if err != nil {
		respond(http.StatusBadRequest, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	respondStatus(http.StatusOK, w)
}

// deleteUser a user
func deleteUser(r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respond(http.StatusBadRequest, map[string]interface{}{"error": "URL Param 'id' is missing"}, w)
		return
	}
	query := "DELETE FROM users WHERE id = ? LIMIT 1"
	DB.MustExec(query, id)
	respondStatus(http.StatusOK, w)
}
