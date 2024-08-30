package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlejoFM/Laravel-Routes-Exporter-PostmanReady/internal/postman"
)

func main() {
	rootDir := "./"

	items := postman.GeneratePostmanItems(rootDir)

	collection := postman.PostmanCollection{
		Info: postman.Info{
			PostmanID: "uuid.New().String()",
			Name:      "Laravel API Requests",
			Schema:    "https://schema.getpostman.com/json/collection/v2.0.0/collection.json",
		},
		Items: items,
	}

	jsonData, err := postman.FormatJSON(collection)
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	outputFile := filepath.Join(rootDir, "postman_collection.json")
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing the JSON file:", err)
		return
	}

	fmt.Printf("JSON File generated: %s\n", outputFile)
}
