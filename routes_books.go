package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func routeBooksGet(c *gin.Context) {
	var books []Book
	db.Find(&books)
	c.JSON(200, books)
}

func routeBooksGetOne(c *gin.Context) {
	id := c.Param("id")
	var book Book
	db.Preload("Author").First(&book, id)
	if book.ID == 0 {
		c.JSON(404, "Book not found")
	} else {
		book.AuthorID = 0
		c.JSON(200, book)
	}
}

func routeBooksPost(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, book)
}

func routeBooksPut(c *gin.Context) {
	id := c.Param("id")
	var changes, book Book
	if err := c.ShouldBindJSON(&changes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.First(&book, id)
	if book.ID == 0 {
		c.JSON(404, "Book not found")
		return
	}
	if err := db.Model(&book).Updates(&changes).Error; err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, book)

}

func routeBooksDelete(c *gin.Context) {
	id := c.Param("id")
	var book Book
	db.First(&book, id)
	if book.ID == 0 {
		c.JSON(404, "Book not found")
	} else {
		db.Delete(&book)
		c.JSON(200, book)
	}
}

func routesBooks(r *gin.RouterGroup) {
	r.GET("", routeBooksGet)
	r.GET("/:id", routeBooksGetOne)
	r.POST("", routeBooksPost)
	r.PUT("/:id", routeBooksPut)
	r.DELETE("/:id", routeBooksDelete)
}
