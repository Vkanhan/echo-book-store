package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var db *sql.DB

func connectToDB() (*sql.DB, error) {

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

	return db, err
}

func insertBook(book Book) int64 {
	sqlStatement := `INSERT INTO books (title, author, price) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	if err := db.QueryRow(sqlStatement, book.Title, book.Author, book.Price).Scan(&id); err != nil {
		log.Printf("Error inserting the book: %v", err)
		return 0
	}
	log.Printf("Inserted a single book with ID: %v", id)
	return id
}

func getBook(id int64) (Book, error) {
	var book Book
	sqlStatement := `SELECT *FROM books WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Println("No rows were returned")
			return book, nil
		default:
			log.Printf("Unable to scan the row: %v", err)
		}
	}
	return book, nil
}

func getAllBooks() ([]Book, error) {
	var books []Book

	sqlStatement := `SELECT *FROM books`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Printf("Unable to execute the query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
			log.Printf("Unable to scan the row: %v", err)
		}
		books = append(books, book)
	}
	return books, err
}

func updateBook(id int64, book Book) int64 {
	sqlStatement := `UPDATE books SET title=$2, author=$3, price=$4 WHERE id=$1`
	res, err := db.Exec(sqlStatement, id, book.Title, book.Author, book.Price)
	if err != nil {
		log.Printf("Unable to execute the query: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("error while checking the affected rows: %v", err)
	}
	log.Printf("Total rows affected: %v", rowsAffected)
	return rowsAffected
}

func deleteBook(id int64) int64 {
	sqlStatement := `DELETE FROM books WHERE id=$1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Printf("Unable to execute the query: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error while checking the affected rows: %v", rowsAffected)
	}
	return rowsAffected
}
