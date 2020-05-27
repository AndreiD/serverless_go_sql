package serverless

import (
	"fmt"
	"net/http"
)

// Register ...
func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! v1.1")
}