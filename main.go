package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
)

type Income struct {
	ReceivedTotal     float64 `json:"received_total"`
	ServicesPerformed float64 `json:"services_performed"`
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

func round(number float64) float64 {
	return float64(int(number*100)) / 100
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
				ReceivedTotal:     round(rand.Float64() * 1000),
				ServicesPerformed: round(rand.Float64() * 1000),
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

func delete_item(c *gin.Context, items *[]Item) {
	var requestBody struct {
		ID uint64 `json:"id"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON format"})
		return
	}
	id := requestBody.ID

}
func init_db() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./db/example.db")
	db, err = create_db(db)
	return db, err
}
func create_db(db *sql.DB) (*sql.DB, error) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    date DATETIME,
    statement_id BIGINT UNSIGNED,
    statement_description VARCHAR(255),
    posting_description VARCHAR(255),
    received_total DECIMAL(10, 2),
    services_performed DECIMAL(10, 2),
    comment TEXT);`)
	return db, err
}
func main() {
	db, err := init_db()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
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
