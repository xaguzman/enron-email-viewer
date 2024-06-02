package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func indexExists(zincURL, indexName string) bool {
	url := fmt.Sprintf("%s/api/index/%s", zincURL, indexName)

	req, err := authenticatedRequest("GET", url, "")

	if err != nil {
		fmt.Printf("error making request tp get index from zinc server: %s\n", err)
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("sending request to check for index failed: %v", err)
		return false
	}

	defer resp.Body.Close()

	// Check if the status code is 200 OK, which means the index exists
	if resp.StatusCode == http.StatusOK {
		return true
	} else if resp.StatusCode == http.StatusNotFound {
		return false
	}

	fmt.Printf("unexpected status code: %d", resp.StatusCode)
	return false
}

/*
*
Returns true if index was succesfully created, false otherwise
*
*/
func createZincIndex(zincSearchUrl string) (bool, error) {

	url := zincSearchUrl + "/api/index"

	exists := indexExists(zincSearchUrl, indexName)

	if exists {
		fmt.Println("Index already exists, skipping index creation and document upload.")
		return false, nil
	}

	fmt.Println("Index doesn't exist; creating...")

	index := `{
		"name": "` + indexName + `",
		"storage_type": "disk",
		"mappings": {
			"properties": {
				"_id": {
					"type":          "keyword",
					"index":         true,
					"store":         false,
					"sortable":      true,
					"aggregatable":  true,
					"highlightable": false
				},
				"@timestamp": {
					"type":          "date",
					"index":         true,
					"store":         false,
					"sortable":      true,
					"aggregatable":  true,
					"highlightable": false
				},
				"Bcc": {
					"type":          "text",
					"index":         true,
					"store":         false,
					"sortable":      false,
					"aggregatable":  false,
					"highlightable": false
				},
				"Body": {
					"type":          "text",
					"index":         true,
					"store":         false,
					"sortable":      false,
					"aggregatable":  false,
					"highlightable": true
				},
				"Cc": {
					"type":          "text",
					"index":         true,
					"store":         false,
					"sortable":      false,
					"aggregatable":  false,
					"highlightable": true
				},
				"From": {
					"type":          "text",
					"index":         true,
					"store":         false,
					"sortable":      false,
					"aggregatable":  false,
					"highlightable": false
				},
				"Subject": {
					"type":          "text",
					"index":         true,
					"store":         false,
					"sortable":      false,
					"aggregatable":  false,
					"highlightable": true
				},
				"To": {
					"type":          "text",
					"index":         true,
					"store":         false,
					"sortable":      false,
					"aggregatable":  false,
					"highlightable": true
				}
			}
		}
	}`

	req, err := authenticatedRequest("POST", url, index)

	if err != nil {
		log.Println("creating authenticated request for /api/index: ", err)
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending /api/index request to ZincSearch: ", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return false, nil
	}

	log.Println("Created index")
	return true, nil
}

func postToZinc(jsonData string, zincSearchUrl string) {
	url := zincSearchUrl + "/api/_bulkv2"

	req, err := authenticatedRequest("POST", url, jsonData)
	if err != nil {
		log.Println("Error creating /api/_bulkv2 request: ", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending /api/bulkv2 request to ZincSearch: ", err)
		return
	}
	defer resp.Body.Close()
}

func authenticatedRequest(method string, url string, body string) (*http.Request, error) {

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Println("Error creating request: ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")

	return req, nil
}
