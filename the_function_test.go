package function

import (
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStuff(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TheFunction)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Error(err)
	}

	userID, _, _, err := jsonparser.Get(body, "id")
	if err != nil {
		t.Error(err)
	}

	// cleanup the user
	err = deleteUserByID(string(userID))
	if err != nil {
		t.Error(err)
	}
}
