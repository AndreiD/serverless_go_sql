package function

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}