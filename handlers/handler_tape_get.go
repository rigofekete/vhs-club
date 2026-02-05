package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Tape struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Director string  `json:"director"`
	Genre    string  `json:"genre"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

var Tapes = []Tape{
	{ID: "1", Title: "2001: A Space Odyssey", Director: "Stanley Kubrick", Genre: "Sci-Fi", Quantity: 1, Price: 5999.99},
	{ID: "2", Title: "Blade Runner", Director: "Ridley Scott", Genre: "Sci-Fi", Quantity: 1, Price: 3999.99},
	{ID: "3", Title: "Dune", Director: "David Lynch", Genre: "Sci-Fi", Quantity: 1, Price: 2999.99},
	{ID: "4", Title: "Stalker", Director: "Andrei Tarkovsky", Genre: "Sci-Fi", Quantity: 1, Price: 5999.99},
	{ID: "5", Title: "Amarcord", Director: "Federico Fellini", Genre: "Drama", Quantity: 1, Price: 5999.99},
}

func GetTapes(c *gin.Context) {
	genre := c.Query("genre")
	if genre != "" {
		filtered := filterGenres(genre)
		if filtered == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("we currently have no %s movies available", genre),
			})
			return
		}
		c.IndentedJSON(http.StatusOK, filtered)
		return
	}
	c.IndentedJSON(http.StatusOK, Tapes)
}

// Helper for genre query parameter
func filterGenres(genre string) []Tape {
	var filtered []Tape
	for _, tape := range Tapes {
		if tape.Genre == genre {
			filtered = append(filtered, tape)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return filtered
}

func GetTapeByID(c *gin.Context) {
	id := c.Param("id")

	for _, tape := range Tapes {
		if tape.ID == id {
			c.IndentedJSON(http.StatusOK, tape)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "tape not found"})
}
