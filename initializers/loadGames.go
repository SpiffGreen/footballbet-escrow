package initializers

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var GameData []map[string]interface{}

func LoadGames() {
	jsonFile, err := os.Open("games.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %s", err)
	}
	defer jsonFile.Close()

	// Read the JSON file into a byte array
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %s", err)
	}

	// Unmarshal the JSON data into the slice of maps
	err = json.Unmarshal(byteValue, &GameData)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
}
