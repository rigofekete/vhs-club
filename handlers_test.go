package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/tapes", getTapes)
	router.GET("/tapes/:id", getTapeByID)
	router.POST("/tapes", postTape)

	return router
}

func TestGetTapes(t *testing.T) {
	router := setupRouter()

	// Test basic GetTapes to return all tapes
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tapes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []tape
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(response))
	assert.Equal(t, "2001: A Space Odyssey", response[0].Title)

	// Test GetTapes with genre query parameter, returning only 1 tape.
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/tapes?genre=Drama", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var response2 []tape
	err2 := json.Unmarshal(w2.Body.Bytes(), &response2)
	assert.NoError(t, err2)
	assert.Equal(t, 1, len(response2))
	assert.Equal(t, "Drama", response2[0].Genre)

	// Test GetTapes with genre query parameter, returning multiple tapes
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/tapes?genre=Sci-Fi", nil)
	router.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Code)

	var response3 []tape
	err3 := json.Unmarshal(w3.Body.Bytes(), &response3)
	assert.NoError(t, err3)
	assert.Equal(t, 4, len(response3))

	for _, tape := range response3 {
		assert.Equal(t, "Sci-Fi", tape.Genre)
	}

	// Test unavalilable genre
	w4 := httptest.NewRecorder()
	genre := "Romantic"
	req4, _ := http.NewRequest("GET", fmt.Sprintf("/tapes?genre=%s", genre), nil)
	router.ServeHTTP(w4, req4)

	assert.Equal(t, http.StatusNotFound, w4.Code)

	// NOTE: gin.H type is a shortcut to map[sting]any, according to the gin-gonic docs
	var response4 gin.H
	err4 := json.Unmarshal(w4.Body.Bytes(), &response4)
	assert.NoError(t, err4)
	assert.Equal(t, fmt.Sprintf("we currently have no %s movies available", genre), response4["message"])
}

func TestGetTapeByID(t *testing.T) {
	router := setupRouter()

	// Test existing ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tapes/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response tape
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "2", response.ID)
	assert.Equal(t, "Blade Runner", response.Title)

	// Test getting non existing ID
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/tapes/10", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusNotFound, w2.Code)

	var response2 map[string]any
	err2 := json.Unmarshal(w2.Body.Bytes(), &response2)
	assert.NoError(t, err2)
	assert.Equal(t, "tape not found", response2["message"])
}

func TestPostTape(t *testing.T) {
	router := setupRouter()

	// Test adding a brand new tape
	w := httptest.NewRecorder()
	request := tape{
		ID:       "6",
		Title:    "A torinói ló",
		Director: "Tarr Béla",
		Genre:    "Drama",
		Quantity: 1,
		Price:    5999.99,
	}

	tapeJSON, err := json.Marshal(request)
	assert.NoError(t, err)
	req, _ := http.NewRequest("POST", "/tapes", strings.NewReader(string(tapeJSON)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response tape
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(tapes))
	assert.Equal(t, 1, tapes[5].Quantity)
	assert.Equal(t, "A torinói ló", tapes[5].Title)

	// Test adding a tape which already exists in the tapes slice
	w2 := httptest.NewRecorder()
	request2 := tape{
		ID:       "6",
		Title:    "A torinói ló",
		Director: "Tarr Béla",
		Genre:    "Drama",
		Quantity: 1,
		Price:    5999.99,
	}

	tapeJSON2, err2 := json.Marshal(request2)
	assert.NoError(t, err2)
	req2, _ := http.NewRequest("POST", "/tapes", strings.NewReader(string(tapeJSON2)))
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response2 gin.H
	err2 = json.Unmarshal(w2.Body.Bytes(), &response2)
	assert.NoError(t, err2)
	assert.Equal(t, 6, len(tapes))
	assert.Equal(t, 2, tapes[5].Quantity)
	assert.Equal(t, fmt.Sprintf("The title, \"%s\", is already in the catalog. We now have %d units in stock.", tapes[5].Title, tapes[5].Quantity), response2["message"])
	// fmt.Println("Message: ", response2["message"])
}
