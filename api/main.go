// main.go
// author: hankbao

package main

import (
	"log"
	"net/http"
	"os"

	"main/handlers"
	"main/models"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Print("Server started")

	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_username := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := "reader"

	var handler handlers.Handler
	db, err := models.ConnectDatabase(db_host, db_port, db_username, db_password, db_name)
	if err != nil {
		log.Fatal(err)
	} else {
		handler = *handlers.NewHandler(db)
	}

	router := gin.Default()

	router.POST("/subscribe", handler.Subscribe)
	router.GET("/feeds/:id", handler.GetFeedById)
	router.GET("/articles/:id", handler.GetArticleById)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "<h1>Hello, World!</h1>")
	})

	log.Fatal(router.Run(":8070"))
}
