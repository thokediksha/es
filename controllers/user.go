package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	models "es/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
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
	c.JSON(http.StatusOK, gin.H{"msg": "user added successfully"})
}

// UpdateDocument updates the specified document in Elasticsearch.
func UpdateDocument(client *elasticsearch.Client, index, typ, id string, doc interface{}) (*esapi.Response, error) {
	// Convert the updated document to JSON
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	// Create the update request
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(docJSON),
	}

	// Execute the update request
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}

	// Check the response status
	if res.IsError() {
		return nil, fmt.Errorf("error updating document: %s", res.String())
	}

	return res, nil
}

func UpdateUser(c *gin.Context) {
	client := GetESClient()
	// Get the index, type, and ID from the URL parameters
	index := c.Param("index")
	id := c.Param("id")

	// Bind the request data to the UpdateDocumentRequest struct
	var req models.User

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create a map with the updated document data
	doc := map[string]interface{}{
		"name": req.Name,
		"age":  req.Age,
	}

	// Update the document in Elasticsearch
	_, err = UpdateDocument(client, index, "_doc", id, doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a successful response
	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
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
