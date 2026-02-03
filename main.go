package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tape struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Director string  `json:"director"`
	Genre    string  `json:"genre"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

var tapes = []tape{
	{ID: "1", Title: "2001: A Space Odyssey", Director: "Stanley Kubrick", Genre: "Sci-Fi", Quantity: 1, Price: 5999.99},
	{ID: "2", Title: "Blade Runner", Director: "Ridley Scott", Genre: "Sci-Fi", Quantity: 1, Price: 3999.99},
	{ID: "3", Title: "Dune", Director: "David Lynch", Genre: "Sci-Fi", Quantity: 1, Price: 2999.99},
	{ID: "4", Title: "Stalker", Director: "Andrei Tarkovsky", Genre: "Sci-Fi", Quantity: 1, Price: 5999.99},
	{ID: "5", Title: "Amarcord", Director: "Federico Fellini", Genre: "Drama", Quantity: 1, Price: 5999.99},
}

func main() {
	router := gin.Default()

	// TODO: Add health handler. Read the MVC link that Péter sent me

	router.GET("/tapes", getTapes)
	router.GET("/tapes/:id", getTapeByID)
	router.POST("/tapes", postTape)

	_ = router.Run("localhost:8080")
}

// Custom Handlers and Helpers ///
/////////////////////////////////

func getTapes(c *gin.Context) {
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
	c.IndentedJSON(http.StatusOK, tapes)
}

// Helper for genre query parameter
func filterGenres(genre string) []tape {
	var filtered []tape
	for _, tape := range tapes {
		if tape.Genre == genre {
			filtered = append(filtered, tape)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return filtered
}

func getTapeByID(c *gin.Context) {
	id := c.Param("id")

	for _, tape := range tapes {
		if tape.ID == id {
			c.IndentedJSON(http.StatusOK, tape)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "tape not found"})
}

// Admin-only handler
func postTape(c *gin.Context) {
	var newTape tape
	if err := c.BindJSON(&newTape); err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "error binding object"})
		return
	}

	// TODO: Needs to be refactored once DB is added
	// TODO: Should we update resources inside a handler with a POST method?
	// For now, just increment tape.Quantity and return a message, if tape exists, indicating the new quantity
	if exists, quantity := tapeExists(newTape); exists {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(
				"The title, \"%s\", is already in the catalog. We now have %d units in stock.", newTape.Title, quantity),
		})
		return
	}

	tapes = append(tapes, newTape)
	c.IndentedJSON(http.StatusCreated, newTape)
}

// Helper for the postTape function
func tapeExists(newTape tape) (bool, int) {
	for i := range tapes {
		if newTape.Title == tapes[i].Title &&
			newTape.Director == tapes[i].Director {
			// TODO: increment quantity in the DB! For now, increment directly and return to respond with message.
			tapes[i].Quantity += 1
			return true, tapes[i].Quantity
		}
	}
	return false, 0
}
