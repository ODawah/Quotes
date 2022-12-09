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
	router.GET("/find_author/:name", Handlers.SearchAuthor)

	log.Fatalln(router.Run(":8080"))
}
