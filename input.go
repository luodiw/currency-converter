package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type DataInput struct {
	Value        float64
	CurrencyFrom string
	CurrencyTo   string
}

func getInput() string {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	rawInput, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading input:", err)
		return getInput()
	}

	checkForQuitSequence(rawInput)

	return rawInput
}

func checkForQuitSequence(input string) {
	if input == ":q\n" {
		os.Exit(0)
	}
}

func formatInput(input string) []string {
	var filterAlphanumericAndDash = func(input string) string {
		var result []rune

		for _, char := range input {
			if unicode.IsLetter(char) || unicode.IsDigit(char) || unicode.IsSpace(char) || char == '-' {
				result = append(result, char)
			}
		}

		return string(result)
	}

	var separateIntegerFromNumber = func(input string) []string {
		var i int

		for i < len(input) && unicode.IsDigit(rune(input[i])) {
			i++
		}

		if i == 0 || i == len(input) {
			return []string{input}
		}

		return []string{input[:i], input[i:]}
	}

	filteredInput := filterAlphanumericAndDash(input)
	capitalizedInput := strings.ToUpper(filteredInput)
	splitInput := strings.Fields(capitalizedInput)

	var formattedInput []string
	for _, word := range splitInput {
		slice := separateIntegerFromNumber(word)
		formattedInput = append(formattedInput, slice...)
	}

	return formattedInput
}

func inputWrapper(rates map[string]float64) DataInput {
	var input DataInput
	rawInput := getInput()
	formattedInput := formatInput(rawInput)

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

	if input.CurrencyFrom == "" || input.CurrencyTo == "" {
		return inputWrapper(rates)
	}
	return input
}
