package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type GameConfig struct {
	Player        string        `json:"player"`
	Ghost         string        `json:"ghost"`
	InfectedGhost string        `json:"eatableGhost"`
	Wall          string        `json:"wall"`
	Point         string        `json:"point"`
	Fruit         string        `json:"fruit"`
	Death         string        `json:"death"`
	Space         string        `json:"space"`
	ImgSize       int           `json:"imgSize"`
	UseEmoji      bool          `json:"useEmoji"`
	FruitDuration time.Duration `json:"fruitDuration"`
}

func LoadGameConfiguration(filePath string, configuration *GameConfig) error {
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
