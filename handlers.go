package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func connectToDB() *sql.DB {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Print("Database connection error")
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	return db
}
	
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return 
	}

	if book.Title == "" || book.Author == "" || book.Price <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	insertID := insertBook(book)
	res := Response{
		ID:      insertID,
		Message: "Book created successfully",
	}

	respondWithJSON(w, http.StatusCreated, res)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return 
	}

	book, err := getBook(int64(id))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch the book")
		return 
	}

	respondWithJSON(w, http.StatusOK, book)
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := getAllBooks()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get all books")
		return 
	}

	respondWithJSON(w, http.StatusOK, books)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return 
	}
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to respond your request")
		return 
	}

	updatedRows := updateBook(int64(id), book)
	msg := fmt.Sprintf("Book successfully updated: %v", updatedRows)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	respondWithJSON(w, http.StatusOK, res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return 
	}

	deletedRows := deleteBook(int64(id))
	msg := fmt.Sprintf("Book deleted successfully: %v", deletedRows)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	respondWithJSON(w, http.StatusOK, res)
}

func insertBook(book Book) int64 {
	db := connectToDB()
	defer db.Close()

	sqlStatement := `INSERT INTO books (title, author, price) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	if err := db.QueryRow(sqlStatement, book.Title, book.Author, book.Price).Scan(&id); err != nil {
		log.Fatalf("Error inserting the book: %v", err)
	}
	fmt.Printf("Inserted a single book: %v", id)
	return id
}

func getBook(id int64) (Book, error) {
	db := connectToDB()
	defer db.Close()

	var book Book
	sqlStatement := `SELECT *FROM books WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned")
			return book, nil
		default:
			log.Fatalf("Unable to scan the row: %v", err)
		}
	}
	return book, nil
}

func getAllBooks() ([]Book, error) {
	db := connectToDB()
	defer db.Close()

	var books []Book

	sqlStatement := `SELECT *FROM books`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
			log.Fatalf("Unable to scan the row: %v", err)
		}
		books = append(books, book)
	}
	return books, err
}

func updateBook(id int64, book Book) int64 {
	db := connectToDB()
	defer db.Close()

	sqlStatement := `UPDATE books SET title=$2, author=$3, price=$4 WHERE id=$1`
	res, err := db.Exec(sqlStatement, id, book.Title, book.Author, book.Price)
	if err != nil {
		log.Fatalf("Unable to execute the query: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking the affected rows: %v", err)
	}
	fmt.Printf("Total rows affected: %v", rowsAffected)
	return rowsAffected
}

func deleteBook(id int64) int64 {
	db := connectToDB()
	defer db.Close()

	sqlStatement := `DELETE FROM books WHERE id=$1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows: %v", rowsAffected)
	}
	return rowsAffected
}
