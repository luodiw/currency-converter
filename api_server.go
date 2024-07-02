package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ConversionRequest struct {
	From   string  `json:"from" query:"from"`
	To     string  `json:"to" query:"to"`
	Amount float64 `json:"amount" query:"amount"`
	Date   string  `json:"date" query:"date"`
}

type ConversionResponse struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Amount    float64   `json:"amount"`
	Result    float64   `json:"result"`
	Date      string    `json:"date"`
	Timestamp time.Time `json:"timestamp"`
}

// func startAPIServer() {
// 	e := echo.New()

// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())

// 	e.GET("/convert", handleConversion)
// 	e.GET("/rates", handleGetRates)

// 	e.Logger.Fatal(e.Start(":8080"))
// }

func startAPIServer() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handleRoot)
	e.GET("/convert", handleConversion)
	e.GET("/rates", handleGetRates)

	// Handle 404 Not Found
	e.Any("*", handle404)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func handleRoot(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Welcome to the Currency Converter API",
		"version": "1.0",
		"endpoints": `
			GET /convert?from=USD&to=EUR&amount=100&date=2023-06-30
			GET /rates?date=2023-06-30
		`,
	})
}

func handle404(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]string{
		"error":   "Route not found",
		"message": "Please check the API documentation for valid endpoints",
	})
}

func handleConversion(c echo.Context) error {
	req := new(ConversionRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	var rates map[string]float64
	var err error

	if req.Date != "" {
		historicalData, err := getHistoricalRate(req.Date)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		rates = castRateFromLatest(historicalData)
	} else {
		now := time.Now().Unix()
		rates, err = caller(now)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get exchange rates"})
		}
	}

	inputData := DataInput{
		Value:        req.Amount,
		CurrencyFrom: req.From,
		CurrencyTo:   req.To,
	}

	result := convert(inputData, rates)

	response := ConversionResponse{
		From:      req.From,
		To:        req.To,
		Amount:    req.Amount,
		Result:    result,
		Date:      req.Date,
		Timestamp: time.Now(),
	}

	return c.JSON(http.StatusOK, response)
}

func handleGetRates(c echo.Context) error {
	date := c.QueryParam("date")

	var rates map[string]float64
	var err error

	if date != "" {
		historicalData, err := getHistoricalRate(date)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		rates = castRateFromLatest(historicalData)
	} else {
		now := time.Now().Unix()
		rates, err = caller(now)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get exchange rates"})
		}
	}

	return c.JSON(http.StatusOK, rates)
}
