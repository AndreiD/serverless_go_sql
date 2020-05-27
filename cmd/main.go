package main

import (
	"fmt"
	function "github.com/AndreiD/serverless_go_sql"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", function.TheFunction)
	fmt.Println("Listening on http://localhost:3000/")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
