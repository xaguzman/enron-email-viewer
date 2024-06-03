package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const defaultZincUrl = "http://localhost:4080"
const defaultZincUser = "admin"
const defaultZincPwd = "admin"

type ZincSearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []Hit   `json:"hits"`
	} `json:"hits"`
}

type Hit struct {
	Index     string  `json:"_index"`
	Type      string  `json:"_type"`
	ID        string  `json:"_id"`
	Score     float64 `json:"_score"`
	Timestamp string  `json:"@timestamp"`
	Source    struct {
		Timestamp string   `json:"@timestamp"`
		Bcc       []string `json:"Bcc"`
		Body      string   `json:"Body"`
		Cc        []string `json:"Cc"`
		Date      string   `json:"Date"`
		From      string   `json:"From"`
		Subject   string   `json:"Subject"`
		To        []string `json:"To"`
	} `json:"_source"`
}

func queryZinc(queryTerm string, indexName string) (*ZincSearchResponse, error) {
	zincSearchUrl := os.Getenv("ZINCSEARCH_URL")

	if zincSearchUrl == "" {
		zincSearchUrl = defaultZincUrl
	}

	if indexName == "" {
		indexName = "enron-emails"
	}

	jsonData := `
    {
        "search_type": "matchphrase",
        "query": {
            "term": "` + queryTerm + `",
            "field": "_all"
        },
        "sort_fields": ["Date"],
        "from": 0,
        "max_results": 20
    }`

	url := fmt.Sprintf("%s/api/%s/_search", zincSearchUrl, indexName)

	req, err := authenticatedRequest("POST", url, jsonData)

	if err != nil {
		log.Printf("error creating authenticated request for /api/%s/index: %s", indexName, err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error sending authenticated request for /api/%s/index: %s", indexName, err)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var zincResponse ZincSearchResponse

	err = json.Unmarshal(bodyBytes, &zincResponse)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// bodyString := string(bodyBytes)

	if resp.StatusCode == 200 {
		return &zincResponse, nil
	} else {
		return nil, errors.New("something went wrong")
	}

}

func authenticatedRequest(method string, url string, body string) (*http.Request, error) {

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Println("Error creating request: ", err)
		return nil, err
	}

	zincSearchUser := os.Getenv("ZINCSEARCH_USER")

	if zincSearchUser == "" {
		zincSearchUser = defaultZincUser
	}

	zincSearchPwd := os.Getenv("ZINCSEARCH_PWD")

	if zincSearchPwd == "" {
		zincSearchPwd = defaultZincPwd
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(zincSearchUser, zincSearchPwd)

	return req, nil
}
