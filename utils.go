package function

import (
	"encoding/json"
	"net/http"
)

func respond(httpStatus int, responseJSON map[string]interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(responseJSON)
}
