package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AndreiD/serverless_go_sql/serverless"
)

func main() {
	http.HandleFunc("/", serverless.Register)
	fmt.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}