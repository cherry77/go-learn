package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

type ApiResponse struct {
	Status bool `json:"status"`
	Data   struct {
		Data []map[string]interface{} `json:"data"`
	} `json:"data"`
}

func main() {
	// 1. Prepare the API request URL with authentication
	url := prepareApiUrl()
	fmt.Println("Request URL:", url)

	// 2. Make the API request
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("API request failed: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Reading response failed: %v", err))
	}

	// 3. Parse the JSON response
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		panic(fmt.Sprintf("JSON parsing failed: %v", err))
	}

	if !apiResponse.Status {
		panic("API returned non-success status")
	}

	// 4. Export to CSV
	err = exportToCSVWithSortedFields(apiResponse.Data.Data, "output.csv")
	if err != nil {
		panic(fmt.Sprintf("CSV export failed: %v", err))
	}

	fmt.Println("Data successfully exported to output.csv")
}

func prepareApiUrl() string {
	clientKey := "13669"
	clientSecretKey := "M7T1EE6E6JTQXL07OUDM"
	baseURL := "https://open.3s.mobvista.com/channel/iaa/v1"

	params := map[string]string{
		"time":              fmt.Sprintf("%d", time.Now().Unix()),
		"client_key":        clientKey,
		"client_secret_key": clientSecretKey,
		"start_date":        "2025-07-04",
		"end_date":          "2025-07-11",
		"page":              "1",
		"per_page":          "200",
	}

	// Sort the keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string for token generation
	values := url.Values{}
	for _, k := range keys {
		values.Add(k, params[k])
	}
	queryString := values.Encode()

	// Generate token
	token := fmt.Sprintf("%x", sha256.Sum256([]byte(queryString)))

	// Remove client_secret_key and add token
	delete(params, "client_secret_key")
	params["token"] = token

	// Build final URL
	finalValues := url.Values{}
	for k, v := range params {
		finalValues.Add(k, v)
	}
	finalURL := fmt.Sprintf("%s?%s", baseURL, finalValues.Encode())

	return finalURL
}

func exportToCSVWithSortedFields(data []map[string]interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if len(data) == 0 {
		return nil
	}

	// Get and sort all unique field names
	allFields := make(map[string]bool)
	for _, record := range data {
		for field := range record {
			allFields[field] = true
		}
	}

	// Convert to slice and sort
	sortedFields := make([]string, 0, len(allFields))
	for field := range allFields {
		sortedFields = append(sortedFields, field)
	}
	sort.Strings(sortedFields)

	// Write sorted headers
	if err := writer.Write(sortedFields); err != nil {
		return err
	}

	// Write data rows with fields in sorted order
	for _, record := range data {
		row := make([]string, len(sortedFields))
		for i, field := range sortedFields {
			val := record[field]
			if val == nil {
				row[i] = ""
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
