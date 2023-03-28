package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	http.Handle("/", router)
	fmt.Println("Connected to port 8492")
	log.Println("Connected to port 8492")
	log.Fatal(http.ListenAndServe(":8492", router))
}
