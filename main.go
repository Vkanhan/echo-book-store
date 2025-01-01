package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db := connectToDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/books", CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", GetBook).Methods("GET")
	r.HandleFunc("/books", GetAllBooks).Methods("GET")
	r.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
