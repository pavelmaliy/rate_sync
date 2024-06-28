package main

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"rateSync/resources"
	"time"
)

func main() {
	url := "https://v6.exchangerate-api.com/v6/9d90a307165975a4b2958f79/pair/EUR/ILS"

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch the URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Failed to parse JSON response: %v\n", err)
		return
	}

	fmt.Println(result["conversion_rate"])
	persist(result["conversion_rate"].(float64))
}

func persist(rate float64) {
	ctx := context.Background()
	bytez, err := resources.EmbeddedFS.ReadFile("firestore/annual-report-52e9f-firebase.json")
	if err != nil {
		log.Fatal(err)
	}
	opt := option.WithCredentialsJSON(bytez)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	now := time.Now()
	localDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	_, _, err = client.Collection("eur_ils").Add(ctx, map[string]interface{}{
		"rate": rate,
		"date": localDate,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully updated rates!!")
}
