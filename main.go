package main

import (
	"os"

	"github.com/devgar/restapi/api"
	"github.com/devgar/restapi/database"
	"github.com/devgar/restapi/lib/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	// load .env environment variables
	// err := godotenv.Load()
	// if err != nil {
	// 	panic(err)
	// }

	// initializes database
	db, _ := database.Initialize()

	port := os.Getenv("PORT")
	app := gin.Default() // create gin app
	app.Use(database.Inject(db))
	app.Use(middlewares.JWTMiddleware())
	api.ApplyRoutes(app) // apply api router
	app.Run(":" + port)  // listen to given port
}
