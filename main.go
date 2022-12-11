package main

import (
	"log"

	"github.com/awesomeQuotes/Handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	var router = gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"Health": "Alive"})
	})
	router.POST("/create_author", Handlers.CreateAuthor)
	router.GET("/find_author/:name", Handlers.SearchAuthor)
	router.POST("/create_quote", Handlers.CreateQuote)
	router.GET("/find_quote", Handlers.SearchQuote)
	router.GET("/find_Author_quotes/:name", Handlers.SearchAuthorQuotes)

	log.Fatalln(router.Run(":8080"))
}
