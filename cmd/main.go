package main

import (
	"fmt"
	"log"
	"net/http"
	function "github.com/AndreiD/serverless_go_sql"
)

func main() {
	http.HandleFunc("/", function.TheFunction)
	fmt.Println("Listening on http://localhost:3000/")
	log.Fatal(http.ListenAndServe(":3000", nil))
}