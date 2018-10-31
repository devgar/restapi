package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func routeBooksGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []Book
	db.Find(&books)
	json.NewEncoder(w).Encode(books)
}

func routeBooksGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book Book
	db.Preload("Author").First(&book, params["id"])
	if book.ID == 0 {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode("NOT FOUND")
	} else {
		book.AuthorID = 0
		json.NewEncoder(w).Encode(book)
	}
}

func routeBooksPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	db.Create(&book)
	json.NewEncoder(w).Encode(book)
}

func routeBooksPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updates, book Book
	_ = json.NewDecoder(r.Body).Decode(&updates)
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	db.First(&book, id)
	if book.ID == 0 {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode("NOT FOUND")
	} else {
		db.Model(&book).Updates(&updates)
		json.NewEncoder(w).Encode(&book)
	}
}

func routeBooksDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	db.First(&book, id)
	if book.ID == 0 {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode("NOT FOUND")
	} else {
		db.Delete(&book)
		json.NewEncoder(w).Encode(book)
	}
}
