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

// Fixed field order as specified
var fixedFieldOrder = []string{
	"date",
	"channel_id",
	"channel_uuid",
	"offer_id",
	"offer_uuid",
	"package",
	"install",
	"impressions",
	"rr_d0",
	"rr_d1",
	"rr_d3",
	"rr_d7",
	"rr_d14",
	"rr_d30",
	"d0_roas",
	"d1_roas",
	"d3_roas",
	"d7_roas",
	"d14_roas",
	"d30_roas",
	"revenue_d0",
	"revenue_d1",
	"revenue_d3",
	"revenue_d7",
	"revenue_d14",
	"revenue_d30",
}

var endDate = time.Now().Format("2006-01-02")
var startDate = time.Now().AddDate(0, 0, -6).Format("2006-01-02") // 6天前+今天=7天

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

	// 4. Export to CSV with fixed field order
	err = exportToCSVWithFixedOrder(apiResponse.Data.Data, fmt.Sprintf("Mobvista_IAA_%s_%s.csv", startDate, endDate))
	if err != nil {
		panic(fmt.Sprintf("CSV export failed: %v", err))
	}

	fmt.Println("Data successfully exported to output.csv with fixed field order")
}

func prepareApiUrl() string {
	clientKey := "13669"
	clientSecretKey := "M7T1EE6E6JTQXL07OUDM"
	baseURL := "https://open.3s.mobvista.com/channel/iaa/v1"

	params := map[string]string{
		"time":              fmt.Sprintf("%d", time.Now().Unix()),
		"client_key":        clientKey,
		"client_secret_key": clientSecretKey,
		"start_date":        startDate,
		"end_date":          endDate,
		"page":              "1",
		"per_page":          "300",
	}

	// Sort keys and build query string for token
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := url.Values{}
	for _, k := range keys {
		values.Add(k, params[k])
	}
	queryString := values.Encode()

	// Generate token
	token := fmt.Sprintf("%x", sha256.Sum256([]byte(queryString)))

	// Build final URL
	delete(params, "client_secret_key")
	params["token"] = token

	finalValues := url.Values{}
	for k, v := range params {
		finalValues.Add(k, v)
	}

	return fmt.Sprintf("%s?%s", baseURL, finalValues.Encode())
}

func exportToCSVWithFixedOrder(data []map[string]interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers in fixed order
	if err := writer.Write(fixedFieldOrder); err != nil {
		return err
	}

	// Write data rows with fields in fixed order
	for _, record := range data {
		row := make([]string, len(fixedFieldOrder))
		for i, field := range fixedFieldOrder {
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
