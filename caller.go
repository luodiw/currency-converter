package main

import (
	"fmt"
	"os"
)

type Latest struct {
	Timestamp int64          `json:"timestamp"`
	Rates     map[string]any `json:"rates"`
}

func checkForForce() bool {
	if len(os.Args) > 1 && (os.Args[1] == "--force" || os.Args[1] == "-f") {
		return true
	}
	return false
}

func caller(now int64) (map[string]float64, error) {
	rates := make(map[string]float64)

	cache, cacheErr := useCache(now, checkForForce())
	if cacheErr == nil {
		rates = castRateFromLatest(cache)
		return rates, nil
	}

	api, apiErr := useApi("")
	if apiErr == nil {
		rates = castRateFromLatest(api)
		return rates, nil
	}

	forcedCache, forcedCacheErr := useCache(now, true)
	if forcedCacheErr == nil {
		rates = castRateFromLatest(forcedCache)
		return rates, nil
	}

	return rates, forcedCacheErr
}

func castRateFromLatest(latestData Latest) map[string]float64 {
	rates := make(map[string]float64)

	for key, value := range latestData.Rates {
		rate, ok := value.(float64)
		if !ok {
			fmt.Println("Error: Unable to convert rate to float64")
			os.Exit(1)
		}

		rates[key] = rate
	}

	return rates
}
