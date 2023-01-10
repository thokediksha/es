package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	models "es/models"
	"fmt"
	"math"
	"net/http"
	"strconv"
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

	// Create a new user
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
		DocumentID: strconv.Itoa(users.ID),
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
		"msg":  "user added successfully",
		"id":   users.ID,
		"name": users.Name,
		"age":  users.Age,
	})
}

func CreateUserBatch(c *gin.Context) {
	client := GetESClient()

	// Create a new user
	var bulkRequest bytes.Buffer
	var users []models.User

	if err := c.BindJSON(&users); err != nil {
		panic(err)
	}

	for _, user := range users {
		// Convert the user to JSON
		userJSON, err := json.Marshal(&user)
		if err != nil {
			panic(err)
		}

		// Add the user to the Bulk request
		bulkRequest.Write(userJSON)
		bulkRequest.Write([]byte("\n"))

		req := esapi.IndexRequest{
			Index:      "my_index",
			DocumentID: strconv.Itoa(user.ID),
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
	c.JSON(http.StatusCreated, gin.H{
		"msg": "user added successfully"})

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
			"age":  users.Age,
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
		"id":   id,
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
	pageNumber, err := strconv.Atoi(c.Query("page"))
    if err != nil {
        pageNumber = 1
    }
    pageSize, err := strconv.Atoi(c.Query("size"))
    if err != nil {
        pageSize = 10
    }

	from := (pageNumber-1)*pageSize
    size := pageSize

	// Set up the update request
	getReq := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(`{"from": ` + strconv.Itoa(from) + `, "size": ` + strconv.Itoa(size) + `}`),
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

	totalPages := math.Ceil(float64(100000) / float64(size))
	c.JSON(http.StatusOK, gin.H{"results": results.Hits.Hits, "currentPage": pageNumber, 
	"totalPages": totalPages})

	// Return the search results in the response
	// c.JSON(http.StatusOK, results.Hits.Hits)
}

func SearchData(c *gin.Context) {
	client := GetESClient()
	// value := c.Query("id")
	value := c.Query("name")
	// value := c.Query("age")
	
	// Parse the page number and page size from the request query parameters
    pageNumber, err := strconv.Atoi(c.Query("page"))
    if err != nil {
        pageNumber = 1
    }
    pageSize, err := strconv.Atoi(c.Query("size"))
    if err != nil {
        pageSize = 10
    }
	
	// Create a search request
	req := esapi.SearchRequest{
		Index: []string{"my_index"},
		Body:  strings.NewReader(`{"query":{"query_string":{"query": "` +value+ `"}},"from": ` + strconv.Itoa((pageNumber-1)*pageSize) + `,
		"size": ` + strconv.Itoa(pageSize) + `}`),
	}

	// Execute the search
	res, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Printf("Error executing search: %s", err)
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