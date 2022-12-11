package Handlers

import (
	"net/http"

	"github.com/awesomeQuotes/Database"
	"github.com/awesomeQuotes/Operations"
	"github.com/awesomeQuotes/Schemas"
	"github.com/gin-gonic/gin"
)

func SearchQuote(c *gin.Context) {
	var input *Schemas.Quote
	db, _ := Database.Connect()
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't Bind the json body"})
		return
	}
	quote, err := Operations.SearchQuote(db, input.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quote)

}

func CreateQuote(c *gin.Context) {
	db, err := Database.Connect()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input Schemas.Quote
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := Operations.InsertQuote(db, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, author)

}

func SearchAuthorQuotes(c *gin.Context) {
	db, err := Database.Connect()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No name inserted"})
		return
	}
	Quotes, err := Operations.AuthorQuotes(db, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, Quotes)
}
