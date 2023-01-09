package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	models "es/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

func GetESClient() *elasticsearch.Client {
	/* Fetching elastic Search Client */
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	return client

}

func CreateUser(c *gin.Context) {
	client := GetESClient()

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
	c.JSON(http.StatusCreated, gin.H{
		"msg": "user added successfully",
        "id": users.ID,
        "name": users.Name,
        "age":  users.Age,
    })
}

func UpdateUser(c *gin.Context) {
    client := GetESClient()

	index := c.Param("index")
    id := c.Param("id")
    
    var users models.User
    if err := c.BindJSON(&users); err != nil {
        panic(err)
    }
    body := map[string]interface{}{
        "doc": map[string]interface{}{
            "name": users.Name,
            "age": users.Age,
        },
    }

    jsonBody, _ := json.Marshal(body)

    req := esapi.UpdateRequest{
        Index:      index,
        DocumentID: id,
        Body:       bytes.NewReader(jsonBody),
    }
    res, _ := req.Do(context.Background(), client)
    defer res.Body.Close()
    fmt.Println(res.String())
    c.JSON(http.StatusCreated, gin.H{
        "id": id,
        "name": users.Name,
        "age":  users.Age,
    })
}

func DeleteUser(c *gin.Context) {
	client := GetESClient()
	id := c.Param("id")
	index := c.Param("index")

	// Set up the update request
	deleteReq := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	// Perform the update
	res, err := deleteReq.Do(context.Background(), client)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(res)
	c.JSON(http.StatusOK, gin.H{"msg": "user deleted successfully"})
}

func GetUser(c *gin.Context) {
	var doc map[string]interface{}
	client := GetESClient()
	id := c.Param("id")
	index := c.Param("index")

	// Set up the update request
	getReq := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}

	// Perform the update
	res, err := getReq.Do(context.Background(), client)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(res)

	err = json.NewDecoder(res.Body).Decode(&doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response body"})
		return
	}

	// Return the document in the response
	c.JSON(http.StatusOK, doc)
}

func GetAllUser(c *gin.Context) {
	client := GetESClient()

	index := c.Param("index")

	// Set up the update request
	getReq := esapi.SearchRequest{
		Index: []string{index},
	}

	// Perform the update
	res, err := getReq.Do(context.Background(), client)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(res)
	var results struct {
		Hits struct {
			Hits []json.RawMessage `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response body"})
		return
	}

	// Return the search results in the response
	c.JSON(http.StatusOK, results.Hits.Hits)
}
