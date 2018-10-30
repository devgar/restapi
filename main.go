package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []Book
	db.Find(&books)
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book Book
	db.Preload("Author").First(&book, params["id"])
	if book.ID == 0 {
		json.NewEncoder(w).Encode("NOT FOUND")
		return
	}
	book.AuthorID = 0
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	db.Create(&book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updates, base Book
	_ = json.NewDecoder(r.Body).Decode(&updates)
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	db.First(&base, id)
	db.Model(&base).Updates(&updates)
	json.NewEncoder(w).Encode(&base)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&[]Book{})
}

func main() {
	r := mux.NewRouter()

	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Book{})
	db.AutoMigrate(&Author{})

	db.Create(&Book{
		Isbn:  "4142312",
		Title: "Book One",
		Author: &Author{
			Firstname: "John",
			Lastname:  "Doe",
		},
	})
	db.Create(&Author{
		Firstname: "Arturo",
		Lastname:  "Perez Reverte",
	})
	db.Create(&Book{
		Isbn:  "4413743",
		Title: "Book Two",
	})

	// ! Need to create authors inside books

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
