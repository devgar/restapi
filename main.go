package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	err := initDB()
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

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
		Isbn:     "4413743",
		Title:    "Book Two",
		AuthorID: 2,
	})

	r.HandleFunc("/api/books", routeBooksGet).Methods("GET")
	r.HandleFunc("/api/books/{id}", routeBooksGetOne).Methods("GET")
	r.HandleFunc("/api/books", routeBooksPost).Methods("POST")
	r.HandleFunc("/api/books/{id}", routeBooksPut).Methods("PUT")
	r.HandleFunc("/api/books/{id}", routeBooksDelete).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
