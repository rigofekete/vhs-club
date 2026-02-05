package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostTape(c *gin.Context) {
	var newTape Tape
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

	Tapes = append(Tapes, newTape)
	c.IndentedJSON(http.StatusCreated, newTape)
}

// Helper for the postTape function
func tapeExists(newTape Tape) (bool, int) {
	for i := range Tapes {
		if newTape.Title == Tapes[i].Title &&
			newTape.Director == Tapes[i].Director {
			// TODO: increment quantity in the DB! For now, increment directly and return to respond with message.
			Tapes[i].Quantity += 1
			return true, Tapes[i].Quantity
		}
	}
	return false, 0
}
