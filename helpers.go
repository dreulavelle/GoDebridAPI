package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// GetApiKey loads the API key from environment variables
func GetApiKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading keys")
	}

	apiKey := os.Getenv("rd_api")
	return apiKey
}

// HttpClient initializes the HTTP client
func HttpClient(api_key string) *Client {
	client := &http.Client{}
	return &Client{
		c:      client,
		apiKey: api_key,
	}
}

// commonGetReq is a helper function for making GET requests
func (c *Client) CommonGetReq(endpoint string, target interface{}) error {
	var wg sync.WaitGroup
	wg.Add(1)

	var resBody []byte
	var err error
	go func() {
		defer wg.Done()
		resBody, err = c.GetReq(endpoint)
	}()
	wg.Wait()

	if err != nil {
		return err
	}

	if err := json.Unmarshal(resBody, target); err != nil {
		return fmt.Errorf("Decode failed")
	}

	return nil
}

func (c *Client) GetReq(endpoint string) ([]byte, error) {
	reqUrl := api_url.ResolveReference(&url.URL{Path: "1.0" + endpoint})
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	resp, err := c.c.Do(req)
	if err != nil {
		fmt.Println("Error when sending get request")
		return nil, err
	}

	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read body")
	}

	return resBody, nil
}

func (c *Client) PostReq(endpoint string, data url.Values) ([]byte, error) {
	reqUrl := api_url.ResolveReference(&url.URL{Path: "1.0" + endpoint})
	req, err := http.NewRequest("POST", reqUrl.String(), strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.c.Do(req)
	if err != nil {
		fmt.Println("Error when sending post request")
		return nil, err
	}

	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read body")
	}

	return resBody, nil
}

func GetFileIdsFromTorrent(val rdTorrentInfoSchema) string {
	allowedFileTypes := []string{"mkv", "srt"}
	var fileIds []string

	for _, val := range val.Files {
		for _, id := range allowedFileTypes {
			if strings.Contains(val.Path, id) {
				fileIds = append(fileIds, strconv.Itoa(val.Id))
			}
		}
	}
	downloadFiles := strings.Join(fileIds, ",")

	return downloadFiles
}
