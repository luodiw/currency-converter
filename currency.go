package main

import "math"

func convert(dataInput DataInput, rates map[string]float64) float64 {

	var conversion float64
	precision := int(getIntEnvVar("PRECISION"))

	var convertFromUsd = func(value float64, targetRate float64) float64 {
		var convertedValue float64 = value * targetRate
		return convertedValue
	}

	if dataInput.CurrencyFrom == "USD" {
		conversion = convertFromUsd(dataInput.Value, rates[dataInput.CurrencyTo])
	} else {
		rateFrom := rates[dataInput.CurrencyFrom]
		var value float64 = dataInput.Value / rateFrom
		conversion = convertFromUsd(value, rates[dataInput.CurrencyTo])
	}

	return math.Round(conversion*(math.Pow10(precision))) / math.Pow10(precision)
}

func checkValidCurrency(rates map[string]float64, currency string) bool {
	for key := range rates {
		if currency == key {
			return true
		}
	}
	return false
}
