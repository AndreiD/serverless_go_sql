package function

import (
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
	respond(http.StatusOK, map[string]interface{}{"get user...": "info..."}, w)
}

// createUser a user
func createUser(r *http.Request, w http.ResponseWriter) {
	r.RequestURI
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

	respond(http.StatusCreated, map[string]interface{}{"status": "created", "id": newID}, w)
}

// updateUser a user
func updateUser(r *http.Request, w http.ResponseWriter) {
	respond(http.StatusOK, map[string]interface{}{"update...": "info..."}, w)
}

// deleteUser a user
func deleteUser(r *http.Request, w http.ResponseWriter) {
	respond(http.StatusOK, map[string]interface{}{"deleted...": "info..."}, w)
}

// deleteUserByID delete a user by id
func deleteUserByID(userID string) error {
	query := "DELETE FROM users WHERE id = ? LIMIT 1"
	DB.MustExec(query, userID)
	return nil
}
