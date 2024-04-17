package config

import (
	"encoding/json"
	"log"
	"os"
)

type Message struct {
	Win  string `json:"win"`
	Lose string `json:"lose"`
}

func LoadMessages(filePath string, configuration *WindowConfig) error {
	file, e := os.Open(filePath)

	if e != nil {
		log.Fatalf("[CONFIG] Couldn't find/open config file. %s", e)

		return e
	}

	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			log.Panicf("[CONFIG] Couldn't close config file. %s", e)
		}
	}(file)

	parser := json.NewDecoder(file)
	e = parser.Decode(&configuration)

	if e != nil {
		log.Fatalf("[CONFIG] Couldn't parse passed config file. %s", e)

		return e
	}

	return nil
}
