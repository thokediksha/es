package controllers

import (
	"context"
	"encoding/json"
	models "es/models"
	"fmt"
	// "net/http"
	// "strconv"

	// "strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/gin-gonic/gin"
	// uuid "github.com/satori/go.uuid"
)

var client *elasticsearch.Client

func CreateUser(c *gin.Context) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	// Create a new user.
	var users models.User

	if err := c.BindJSON(&users); err != nil {
		panic(err)
	}

	fmt.Println(users.ID)
	// Convert the user to JSON.
	userJSON, err := json.Marshal(&users)
	if err != nil {
		panic(err)
	}

	// id := uuid.NewV4()

	// Create a new Index request
	req := esapi.IndexRequest{
		Index:      "my_index",
		DocumentID: users.ID,
		Body:       strings.NewReader(string(userJSON)),
		Refresh:    "true",
	}

	// Send the Index request
	res, err := req.Do(context.Background(), client)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Print the response
	fmt.Println(res)
}

func UpdateUser(c *gin.Context) {

	// Get the ID of the document to update
	id := c.Param("id")

	// Create a new user.
	var users models.User

	if err := c.BindJSON(&users); err != nil {
		panic(err)
	}
  
	fmt.Println(id)
	// Convert the user to JSON.
	userJSON, err := json.Marshal(&users)
	if err != nil {
		panic(err)
	}

	// id := uuid.NewV4()

	// Create a new Index request
	req := esapi.UpdateRequest{
		DocumentID: id,
		Body:       strings.NewReader(string(userJSON)),
		Refresh:    "true",
	}

	// Send the Index request
	res, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// Print the response
	fmt.Println(res)
}

func GetUser(c *gin.Context) {

}
