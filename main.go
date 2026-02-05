package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rigofekete/vhs-club/handlers"
)

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.GET("/tapes", handlers.GetTapes)
	router.GET("/tapes/:id", handlers.GetTapeByID)
	router.POST("/tapes", handlers.PostTape)

	_ = router.Run("localhost:8080")
}
