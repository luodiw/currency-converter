package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	initRedis()

	if len(os.Args) > 1 && os.Args[1] == "api" {
		fmt.Println("Starting API server...")
		startAPIServer()
	} else {
		runCLI()
	}

	if len(os.Args) > 1 && os.Args[1] == "api" {
		fmt.Println("Starting API server...")
		startAPIServer()
	} else {
		runCLI()
	}
}

func runCLI() {
	displayWelcomeScreen()

	if checkForForce() {
		fmt.Println("Conversion is forced to update from API")
	}

	for {
		rawInput := getInput()

		if shouldTerminate(rawInput) {
			fmt.Println("Thank you for using the Currency Converter CLI. Goodbye!")
			os.Exit(0)
		}

		input, date := extractDate(rawInput)
		formattedInput := formatInput(input)

		var rates map[string]float64
		var err error

		if date != "" {
			historicalData, err := getHistoricalRate(date)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			rates = castRateFromLatest(historicalData)
		} else {
			now := time.Now().Unix()
			rates, err = caller(now)
			if err != nil {
				fmt.Println("Error: No currency exchange data found")
				continue
			}
		}

		inputData := processInput(formattedInput, rates)

		if inputData.Value == 0 && inputData.CurrencyFrom == "" && inputData.CurrencyTo == "" {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		result := convert(inputData, rates)
		if date != "" {
			fmt.Printf("On %s: %.2f %s = %.2f %s\n", date, inputData.Value, inputData.CurrencyFrom, result, inputData.CurrencyTo)
		} else {
			fmt.Printf("%.2f %s = %.2f %s\n", inputData.Value, inputData.CurrencyFrom, result, inputData.CurrencyTo)
		}
	}
}

func displayWelcomeScreen() {
	fmt.Println("====================================")
	fmt.Println("Welcome to the Currency Converter CLI")
	fmt.Println("====================================")
	fmt.Println("Enter your conversion query, e.g.:")
	fmt.Println("  '5 USD to EUR'")
	fmt.Println("  'How much is 100 JPY in GBP?'")
	fmt.Println("  '10 USD to EUR on 2022-01-01'")
	fmt.Println()
	fmt.Println("To exit, type 'exit', 'quit', 'end', or 'thank you'")
	fmt.Println("====================================")
	fmt.Println()
}

func shouldTerminate(input string) bool {
	terminationPhrases := []string{"exit", "quit", "end", "thank you", "goodbye", "bye", "finished", "done"}
	lowercaseInput := strings.ToLower(strings.TrimSpace(input))

	for _, phrase := range terminationPhrases {
		if strings.Contains(lowercaseInput, phrase) {
			return true
		}
	}

	return false
}

func extractDate(input string) (string, string) {
	re := regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}\b`)
	date := re.FindString(input)
	if date != "" {
		return strings.TrimSpace(strings.Replace(input, date, "", -1)), date
	}
	return input, ""
}

func processInput(formattedInput []string, rates map[string]float64) DataInput {
	var input DataInput

	for i := 0; i < len(formattedInput); i++ {
		if input.Value == 0 {
			val, err := strconv.ParseFloat(formattedInput[i], 64)
			if err == nil {
				input.Value = val
			}
		}
		if checkValidCurrency(rates, formattedInput[i]) {
			currentWord := formattedInput[i]
			if input.CurrencyFrom == "" {
				input.CurrencyFrom = currentWord
			} else if input.CurrencyTo == "" {
				input.CurrencyTo = currentWord
			}
		}
	}

	if input.Value == 0 {
		input.Value = 1
	}

	return input
}
