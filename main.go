package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

const (
	serverAddr = "localhost:8080"
)

type Book struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Published bool   `json:"published"`
}

var books = []*Book{
	{Title: "Percy Jackson and the Lightning Thief", Author: "Rick Riordan", Published: true},
	{Title: "Percy Jackson and the Sea of Monsters", Author: "Rick Riordan", Published: true},
	{Title: "A Feast for Crows", Author: "George R.R. Martin", Published: true},
	{Title: "A Dance with Dragons", Author: "George R.R. Martin", Published: true},
	{Title: "The Winds of Winter", Author: "George R.R. Martin", Published: false},
	{Title: "A Dream of Spring", Author: "George R.R. Martin", Published: false},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getBookFromTitle(c *gin.Context) {
	searchTitle := strings.ToLower(c.Param("title"))
	var bookTitles []string
	for _, book := range books {
		bookTitles = append(bookTitles, strings.ToLower(book.Title))
	}

	matches := fuzzy.Find(searchTitle, bookTitles)
	if len(matches) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	for i, title := range bookTitles {
		if title == matches[0] {
			c.IndentedJSON(http.StatusOK, books[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:title", getBookFromTitle)

	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
