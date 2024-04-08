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
	mazeDimensions model.MazeDimensions

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

type Game struct {
	row int
	col int
}

func (g *Game) Update() error {
	moves := model.CheckPossibleMoves(pacman, windowConfig, mazeDimensions,
		model.MazeCharacter{Row: g.row, Col: g.col}, maze)
	model.ProcessInput(&pacman, moves)
	maze = model.ProcessPoint(maze, model.MazeCharacter{Row: g.row, Col: g.col}, &score)

	g.row = pacman.Y / windowConfig.CharSize
	g.col = pacman.X / windowConfig.CharSize

	if g.row > mazeDimensions.Width-1 {
		g.row = mazeDimensions.Width - 1
	} else if g.row < 1 {
		g.row = 0
	}

	if g.col > mazeDimensions.Height-1 {
		g.col = mazeDimensions.Height - 1
	} else if g.col < 1 {
		g.col = 0
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

	if e := ebiten.RunGame(&Game{
		row: pacman.X / windowConfig.CharSize,
		col: pacman.Y / windowConfig.CharSize,
	}); e != nil {
		log.Fatalf("[GAME] Startup failed. %s", e)
	}
}
