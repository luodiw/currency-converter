// package main

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"time"
// )

// func useApi(date string) (Latest, error) {
// 	fileName := getEnvVar("FILE_NAME")
// 	var apiEndPoint string
// 	var appID string = getEnvVar("APP_ID")
// 	var result Latest

// 	if date == "" {
// 		apiEndPoint = "https://openexchangerates.org/api/latest.json?app_id=" + appID
// 	} else {
// 		apiEndPoint = fmt.Sprintf("https://openexchangerates.org/api/historical/%s.json?app_id=%s", date, appID)
// 	}

// 	response, err := http.Get(apiEndPoint)

// 	if err != nil {
// 		return result, errors.New("get request failed")
// 	}

// 	responseData, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	writeToCache(responseData, fileName)

// 	json.Unmarshal([]byte(responseData), &result)

// 	return result, nil
// }

// func getHistoricalRate(date string) (Latest, error) {
// 	_, err := time.Parse("2006-01-02", date)
// 	if err != nil {
// 		return Latest{}, errors.New("invalid date format; please use YYYY-MM-DD")
// 	}

// 	return useApi(date)
// }

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func useApi(date string) (Latest, error) {
	var apiEndPoint string
	var appID string = getEnvVar("APP_ID")
	var result Latest

	cacheKey := fmt.Sprintf("exchange_rates:%s", date)

	cachedData, err := getFromCache(cacheKey)
	if err == nil {
		log.Printf("Data found in Redis cache for date: %s\n", date)
		return cachedData, nil
	} else {
		log.Printf("Data not found in Redis cache for date: %s. Error: %v\n", date, err)
	}

	if date == "" {
		apiEndPoint = "https://openexchangerates.org/api/latest.json?app_id=" + appID
		log.Println("Fetching latest exchange rates")
	} else {
		apiEndPoint = fmt.Sprintf("https://openexchangerates.org/api/historical/%s.json?app_id=%s", date, appID)
		log.Printf("Fetching historical exchange rates for date: %s\n", date)
	}

	response, err := http.Get(apiEndPoint)
	if err != nil {
		log.Printf("API request failed: %v\n", err)
		return result, fmt.Errorf("get request failed: %w", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(responseData, &result)
	if err != nil {
		log.Printf("Failed to parse JSON response: %v\n", err)
		return result, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	err = setToCache(cacheKey, result, 24*time.Hour)
	if err != nil {
		log.Printf("Failed to cache data in Redis: %v\n", err)
	} else {
		log.Printf("Successfully cached data in Redis for date: %s\n", date)
	}

	log.Println("Successfully fetched and processed exchange rates")
	return result, nil
}

func getHistoricalRate(date string) (Latest, error) {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Printf("Invalid date format provided: %s\n", date)
		return Latest{}, errors.New("invalid date format. please use YYYY-MM-DD")
	}

	log.Printf("Fetching historical rates for date: %s\n", date)
	return useApi(date)
}
