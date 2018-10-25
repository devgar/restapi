package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

func newBook(isbn string, title string, author Author) *Book {
	lastBookID++
	return &Book{
		ID:     strconv.Itoa(lastBookID),
		Isbn:   isbn,
		Title:  title,
		Author: &author,
	}
}

// Author Struct (Model)
type Author struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func newAuthor(firstname string, lastname string) *Author {
	lastAuthorID++
	return &Author{
		ID:        strconv.Itoa(lastAuthorID),
		Firstname: firstname,
		Lastname:  lastname,
	}
}

var books []Book
var lastBookID = 0
var lastAuthorID = 0

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	lastBookID++
	book.ID = strconv.Itoa(lastBookID)
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index := range books {
		if books[index].ID == params["id"] {
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			if book.Isbn != "" {
				books[index].Isbn = book.Isbn
			}
			if book.Title != "" {
				books[index].Title = book.Title
			}
			_ = json.NewEncoder(w).Encode(books[index])
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	r := mux.NewRouter()

	books = append(books, *newBook("4413743", "Book One",
		*newAuthor("Jowh", "Doe"),
	))

	books = append(books, *newBook("4113843", "Book Two",
		*newAuthor("Steve", "Smith"),
	))

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
