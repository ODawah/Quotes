package main

import (
	"log"

	"github.com/awesomeQuotes/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	var router = gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"Health": "Alive"})
	})
	router.POST("/create_author", handlers.CreateAuthor)
	router.GET("/find_author/:name", handlers.SearchAuthor)
	router.POST("/create_quote", handlers.CreateQuote)
	router.GET("/find_quote", handlers.SearchQuote)
	router.GET("/find_Author_quotes/:name", handlers.SearchAuthorQuotes)

	log.Fatalln(router.Run(":8080"))
}
