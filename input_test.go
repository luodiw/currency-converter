package main

import (
	"testing"
)

type TestData struct {
	Inputs []string
	Output DataInput
}

func TestEvaluateInput(t *testing.T) {
	cacheData := readFromCache()
	rates := castRateFromLatest(cacheData)

	allTestData := []TestData{
		{
			Inputs: []string{
				"USD to EUR", "usd to eur", "USD TO EUR",
				"Convert USD to EUR", "convert usd to eur", "CONVERT USD TO EUR",
				"1 USD in EUR", "1 usd In Eur", "1 USD IN EUR",
				"exchange USD to EUR", "Exchange usd to eur", "EXCHANGE USD TO EUR",
				"What is the Equivalent of 1 USD in EUR?", "What is the equivalent of 1 USD in EUR?", "What is THE EQUIVALANT of 1 USD in EUR?",
				"USD to EUR CONVERSION", "USD to eur conv", "usd to Eur conversion",
				"How much is 1 usd in EUR", "How MUHC is 1 usd in EUR", "How MUCH is 1 usd in EUR",
				"usd to eur", "usd in eur", "USD in EUR",
			},
			Output: DataInput{Value: 1, CurrencyFrom: "USD", CurrencyTo: "EUR"},
		},
		{
			Inputs: []string{
				"100JPY in USD", "100 JPY to USD", "100 jpy TO usd",
				"Convert 100 JPY to USD", "convert 100 jpy to usd", "CONVERT 100 JPY to usd",
				"100 JPY in USD", "100 jpy In Usd", "100 JPY IN USd",
				"exchange 100 JPY to USD", "Exchange 100 jpy to USD", "EXCHANGE 100 JPY to Usd",
				"What is the Equivalent of 100 JPY in USD?", "What is the equivalent of 100 JPY in USD?", "What is THE EQUIVALANT of 100 JPY in USD?",
				"100 JPY to USD CONVERSION", "100 JPY to usd conv", "100 jpy to USd conversion",
				"How much is 100 jpy in USD", "How MUHC is 100 jpy in USD", "How MUCH is 100 jpy in USD",
				"100jpy to usd", "100jpy in usd", "100JPY in USD",
			},
			Output: DataInput{Value: 100, CurrencyFrom: "JPY", CurrencyTo: "USD"},
		},
		{
			Inputs: []string{
				"CAD AUD", "cad to aud", "CAD TO AUD",
				"Convert CAD to AUD", "convert cad to aud", "CONVERT CAD TO AUD",
				"1 CAD in AUD", "1 cad In Aud", "1 CAD IN AUD",
				"exchange CAD to AUD", "Exchange cad to aud", "EXCHANGE CAD TO AUD",
				"What is the Equivalent of 1 CAD in AUD?", "What is the equivalent of 1 CAD in AUD?", "What is THE EQUIVALANT of 1 CAD in AUD?",
				"CAD to AUD CONVERSION", "CAD to aud conv", "cad to Aud conversion",
				"How much is 1 cad in AUD", "How MUHC is 1 cad in AUD", "How MUCH is 1 cad in AUD",
				"cad to aud", "cad in aud", "CAD in AUD",
			},
			Output: DataInput{Value: 1, CurrencyFrom: "CAD", CurrencyTo: "AUD"},
		},
		{
			Inputs: []string{
				"50.5 EUR to CHF", "50.5 eur to chf", "50.5 EUR TO CHF",
				"Convert 50.5 EUR to CHF", "convert 50.5 eur to chf", "CONVERT 50.5 EUR TO CHF",
				"50.5 EUR in CHF", "50.5 eur In Chf", "50.5 EUR IN CHF",
				"exchange 50.5 EUR to CHF", "Exchange 50.5 eur to chf", "EXCHANGE 50.5 EUR TO CHF",
				"What is the Equivalent of 50.5 EUR in CHF?", "What is the equivalent of 50.5 EUR in CHF?", "What is THE EQUIVALANT of 50.5 EUR in CHF?",
				"50.5 EUR to CHF CONVERSION", "50.5 EUR to chf conv", "50.5 eur to Chf conversion",
				"How much is 50.5 eur in CHF", "How MUHC is 50.5 eur in CHF", "How MUCH is 50.5 eur in CHF",
				"50.5eur to chf", "50.5eur in chf", "50.5EUR in CHF",
			},
			Output: DataInput{Value: 50.5, CurrencyFrom: "EUR", CurrencyTo: "CHF"},
		},
	}

	for _, testData := range allTestData {
		for _, input := range testData.Inputs {
			evalOutput := inputWrapper(rates)
			if testData.Output != evalOutput {
				t.Errorf(`FAILED: evaluateInput(rates, %v) Expected %v, got %v.`, input, testData.Output, evalOutput)
			} else {
				t.Logf(`PASSED: evaluateInput(rates, %v) Returned %v`, input, testData.Output)
			}
		}
	}
}
