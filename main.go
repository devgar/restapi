package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	routesBooks(r.Group("/api/books"))
	routesBlogPosts(r.Group("/api/blogposts"))
	routesStream(r.Group("/api/stream"))
	return r
}

func main() {

	if err := initDB(); err != nil {
		panic("failed to connect database")
	}

	log.Fatal(setupRouter().Run(":8000"))

	db.Close()
}
