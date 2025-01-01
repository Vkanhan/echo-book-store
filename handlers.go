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


	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found on the .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the enviroment")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("error loading postgres database")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("error pinging the database")
	}

	fmt.Println("Successfully connected")

	return db
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Fatalf("error decoding the request body: %v", err)
	}

	insertID := insertBook(book)
	res := Response{
		ID:      insertID,
		Message: "Book created successfully",
	}

	json.NewEncoder(w).Encode(res)

}

func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert the string into int: %v", err)
	}

	book, err := getBook(int64(id))
	if err != nil {
		log.Fatalf("unable to get book: %v", err)
	}

	json.NewEncoder(w).Encode(book)
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := getAllBooks()
	if err != nil {
		log.Fatalf("Unable to get all books: %v", err)
	}
	json.NewEncoder(w).Encode(books)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int: %v", err)
	}
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Fatalf("Unable to decode the request body: %v", err)
	}
	updatedRows := updateBook(int64(id), book)
	msg := fmt.Sprintf("Book successfully updated. Total rows affected: %v", updatedRows)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int: %v", err)
	}

	deletedRows := deleteBook(int64(id))
	msg := fmt.Sprintf("Book deleted successfully. Total rows affected: %v", deletedRows)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func insertBook(book Book) int64 {
	db := connectToDB()
	defer db.Close()

	sqlStatement := `INSERT INTO books (title, author, price) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	if err := db.QueryRow(sqlStatement, book.Title, book.Author, book.Price).Scan(&id); err != nil {
		log.Fatalf("unable to execute the query: %v", err)
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
