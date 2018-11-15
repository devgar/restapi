package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	routesBooks(r.Group("/api/books"))
	return r
}

func main() {

	if err := initDB(); err != nil {
		panic("failed to connect database")
	}

	// db.Create(&Book{
	// 	Isbn:  "4142312",
	// 	Title: "Book One",
	// 	Author: &Author{
	// 		Firstname: "John",
	// 		Lastname:  "Doe",
	// 	},
	// })
	// db.Create(&Book{
	// 	Isbn:   "4413743",
	// 	Title:  "Book Two",
	//	Author: &Author{
	//		Firstname: "Arturo",
	//		Lastname:  "Perez",
	//	},
	// })

	log.Fatal(setupRouter().Run(":8000"))

	db.Close()
}
