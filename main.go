package main

import (
	"Pacman/src"
	"Pacman/src/config"
	"Pacman/src/factory"
	model "Pacman/src/model"
	"flag"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"log"
	"time"
)

var (
	gameConfigFile   = flag.String("game", "./config/game.json", "path to game configuration file")
	windowConfigFile = flag.String("window", "./config/window.json", "path to window configuration file")
	mazeLayoutFile   = flag.String("maze", "./config/maze.txt", "path to maze layout file")

	gameConfig     config.GameConfig
	windowConfig   config.WindowConfig
	mazeDimensions src.MazeDimensions

	imgFactory factory.AssetsFactory

	lives       = 3
	ghostsCount = 4
	isOver      = false

	pacman     model.Sprite
	ghosts     []*model.Ghost
	maze       []string
	score      int
	dotsCount  int
	fruitTimer *time.Timer
)

type Game struct{}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		pacman.Y--
		pacman.Direction = model.DrawDirection(ebiten.KeyW)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		pacman.Y++
		pacman.Direction = model.DrawDirection(ebiten.KeyS)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		pacman.X++
		pacman.Direction = model.DrawDirection(ebiten.KeyD)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		pacman.X--
		pacman.Direction = model.DrawDirection(ebiten.KeyA)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for row, line := range maze {
		for col, char := range line {
			read := model.MazeCharacter{Row: row, Col: col, Line: line, Char: char}
			src.DrawMaze(screen, &read, &windowConfig, &imgFactory)
			model.DrawGhosts(screen, &read, &windowConfig, &imgFactory, &pacman, ghosts, ghostsCount, &dotsCount)
		}
	}
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	if windowConfig.Height == 0 || windowConfig.Width == 0 {
		return windowConfig.CharSize * mazeDimensions.Width, windowConfig.CharSize * mazeDimensions.Height
	}
	return windowConfig.Width, windowConfig.Height
}

func main() {
	e := src.InitializeConfiguration(*gameConfigFile, *windowConfigFile, &gameConfig, &windowConfig)
	if e != nil {
		log.Fatal(e)

		return
	}

	e = src.LoadMaze(*mazeLayoutFile, &maze, &ghosts, &pacman, &dotsCount, ghostsCount, &mazeDimensions, &windowConfig)
	if e != nil {
		log.Panicf("[GAME] Maze load failed")
	}

	e = factory.Create(&imgFactory)
	if e != nil {
		log.Panicf("[GAME] Assest load failed")
	}

	ebiten.SetWindowTitle(windowConfig.Name)
	ebiten.SetWindowIcon([]image.Image{imgFactory.Icon})
	if windowConfig.Height == 0 || windowConfig.Width == 0 {
		ebiten.SetWindowSize(windowConfig.CharSize*mazeDimensions.Width, windowConfig.CharSize*mazeDimensions.Height)
	} else {
		ebiten.SetWindowSize(windowConfig.Width, windowConfig.Height)
	}

	if e := ebiten.RunGame(&Game{}); e != nil {
		log.Fatalf("[GAME] Startup failed. %s", e)
	}
}
