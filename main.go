package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	router.Get("/healthz", handlerReadiness)
	router.Post("/books", CreateBook)
	router.Get("/books/{id}", GetBook)
	router.Get("/books", GetAllBooks)
	router.Put("/books/{id}", UpdateBook)
	router.Delete("/books/{id}", DeleteBook)

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the env")
	}

	server := &http.Server{
		Addr: ":"+ portString,
		Handler: router,
	}

	log.Println("Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
