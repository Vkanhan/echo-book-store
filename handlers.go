package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
