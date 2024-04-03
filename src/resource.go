package src

import (
	"Pacman/src/config"
	"Pacman/src/model"
	"log"
)

func InitializeConfiguration(
	gameConfigFile string,
	windowConfigFile string,
	gameConfig *config.GameConfig,
	windowConfig *config.WindowConfig) error {
	e := config.LoadWindowConfiguration(windowConfigFile, windowConfig)

	if e != nil {
		log.Fatalf("[RESOURCE] Window configuration load failed. %s", e)

		return e
	}

	e = config.LoadGameConfiguration(gameConfigFile, gameConfig)

	if e != nil {
		log.Fatalf("[RESOURCE] Game configuration load failed. %s", e)

		return e
	}

	return nil
}

func InitializeGame(
	mazeFile string,
	maze *[]string,
	ghosts *[]*model.Ghost,
	pacman *model.Sprite,
	dotsCount *int,
	ghostsCount int,
	dimensions *MazeDimensions) error {
	e := LoadMaze(mazeFile, maze, ghosts, pacman, dotsCount, ghostsCount, dimensions)

	if e != nil {
		log.Panicf("[RESOURCE] Maze layout load failed. %s", e)

		return e
	}

	return nil
}
