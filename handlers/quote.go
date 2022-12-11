package handlers

import (
	"net/http"

	"github.com/awesomeQuotes/database"
	"github.com/awesomeQuotes/operations"
	"github.com/awesomeQuotes/schemas"
	"github.com/gin-gonic/gin"
)

func SearchQuote(c *gin.Context) {
	var input *schemas.Quote
	db, _ := database.Connect()
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't Bind the json body"})
		return
	}
	quote, err := operations.SearchQuote(db, input.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quote)

}

func CreateQuote(c *gin.Context) {
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input schemas.Quote
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := operations.InsertQuote(db, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, author)

}

func SearchAuthorQuotes(c *gin.Context) {
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No name inserted"})
		return
	}
	Quotes, err := operations.AuthorQuotes(db, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, Quotes)
}
