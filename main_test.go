package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEmptyData(t *testing.T) {
	if err := initDB(); err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/books", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[]", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/books/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	body, _ := json.Marshal(&Book{
		Title: "New title",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/books/1", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/books/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestCreateBookWithNonExistingAuthor(t *testing.T) {
	if err := initDB(); err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	err := db.Create(&Book{
		Title:    "Test to fail",
		Isbn:     "123456789-ABCD",
		AuthorID: 99,
	}).Error

	assert.NotNil(t, err)
}

func TestPingRoute(t *testing.T) {

	if err := initDB(); err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	router := setupRouter()

	body, _ := json.Marshal(&Book{
		Title: "The first book",
		Isbn:  "981234155123-12312",
		Author: &Author{
			Firstname: "Josh",
			Lastname:  "Weleer",
		},
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/books", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "[]", w.Body.String())

	req, _ = http.NewRequest("GET", "/api/books", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "[]", w.Body.String())

	req, _ = http.NewRequest("GET", "/api/books/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
