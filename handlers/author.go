package handlers

import (
	"net/http"

	"github.com/awesomeQuotes/database"
	"github.com/awesomeQuotes/operations"
	"github.com/awesomeQuotes/schemas"
	"github.com/gin-gonic/gin"
)

func SearchAuthor(c *gin.Context) {
	db, err := database.Connect()
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
	author, err := operations.SearchAuthor(db, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, author)
}

func CreateAuthor(c *gin.Context) {
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	var input schemas.Author
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	author, err := operations.InsertAuthor(db, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, author)
}
