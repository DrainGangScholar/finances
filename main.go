package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

type Income struct {
	ReceivedTotal     int64 `json:"received_total"`
	ServicesPerformed int64 `json:"services_performed"`
}

type Item struct {
	ID                   uint64    `json:"id"`
	Date                 time.Time `json:"date"`
	StatementID          uint64    `json:"statement_id"`
	StatementDescription string    `json:"statement_description"`
	PostingDescription   string    `json:"posting_description"`
	Income               Income    `json:"income"`
	Comment              string    `json:"comment"`
}

func generate_test_data(numItems int) []Item {
	var testData []Item
	for i := 0; i < numItems; i++ {
		testData = append(testData, Item{
			ID:                   uint64(i + 1),
			Date:                 time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour),
			StatementID:          uint64(rand.Intn(1000) + 1),
			StatementDescription: fmt.Sprintf("Statement %d", rand.Intn(10)+1),
			PostingDescription:   fmt.Sprintf("Posting %d", rand.Intn(10)+1),
			Income: Income{
				ReceivedTotal:     rand.Int63n(10000) + 1,
				ServicesPerformed: rand.Int63n(1000) + 1,
			},
			Comment: fmt.Sprintf("Comment %d", rand.Intn(10)+1),
		})
	}
	return testData
}

func get_items(c *gin.Context, items []Item) {
	c.IndentedJSON(http.StatusOK, items)
}

func get_item(c *gin.Context, items *[]Item) {
	var requestBody struct {
		ID uint64 `json:"id"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON format"})
		return
	}
	id := requestBody.ID
	for _, item := range *items {
		if item.ID == id {
			c.IndentedJSON(http.StatusOK, item)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
}

func create_item(c *gin.Context, items *[]Item) {
	var new_item Item

	if err := c.BindJSON(&new_item); err != nil {
		return
	}
	*items = append(*items, new_item)
	c.IndentedJSON(http.StatusCreated, new_item)
}

func main() {
	items := generate_test_data(10)
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/items", func(c *gin.Context) {
		get_items(c, items)
	})
	router.GET("/item", func(c *gin.Context) {
		get_item(c, &items)
	})
	router.POST("/item", func(c *gin.Context) {
		create_item(c, &items)
	})
	router.Run("localhost:8080")
}
