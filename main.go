package main

import (
	"Pacman/src"
	"Pacman/src/config"
	"Pacman/src/factory"
	"Pacman/src/model"
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for row, line := range maze {
		for col, character := range line {
			rect := ebiten.NewImage(windowConfig.CharSize, windowConfig.CharSize)
			offset := float64(0)
			scale := false

			switch character {
			case '#':
				rect.Fill(color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff})
			case '.':
				rect = imgFactory.Dot
				scale = true
				offset = float64(windowConfig.CharSize / 2)
			case 'X':
				rect = imgFactory.Fruit
				scale = true
				offset = (float64(windowConfig.CharSize) * windowConfig.ScaleFactor) / 2
			default:
				continue
			}
			options := &ebiten.DrawImageOptions{}
			if scale {
				options.GeoM.Scale(windowConfig.ScaleFactor, windowConfig.ScaleFactor)
			}

			options.GeoM.Translate(float64(col*windowConfig.CharSize)+offset,
				float64(row*windowConfig.CharSize)+offset)
			screen.DrawImage(rect, options)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
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

	e = src.LoadMaze(*mazeLayoutFile, &maze, &ghosts, &pacman, &dotsCount, ghostsCount, &mazeDimensions)
	if e != nil {
		log.Panicf("[GAME] Maze load failed")
	}

	e = factory.Create(&imgFactory, gameConfig.ImgSize, gameConfig.ImgSize)
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
