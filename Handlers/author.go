package Handlers

import (
	"net/http"

	"github.com/awesomeQuotes/Database"
	"github.com/awesomeQuotes/Operations"
	"github.com/awesomeQuotes/Schemas"
	"github.com/gin-gonic/gin"
)

func SearchAuthor(c *gin.Context) {
	db, err := Database.Connect()
	defer db.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No name inserted"})
		return
	}
	author, err := Operations.SearchAuthor(db, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, author)
}

func CreateAuthor(c *gin.Context) {
	db, err := Database.Connect()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	var input Schemas.Author
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	author, err := Operations.InsertAuthor(db, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, author)
}
