package main

import (
	"Pacman/src"
	"Pacman/src/config"
	"Pacman/src/factory"
	"Pacman/src/model"
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	fontSize = 72
	Dpi      = 72
)

var (
	gameConfigFile   = flag.String("game", "./config/game.json", "path to game configuration file")
	windowConfigFile = flag.String("window", "./config/window.json", "path to window configuration file")
	mazeLayoutFile   = flag.String("maze", "./config/maze.txt", "path to maze layout file")
	messagesFile     = flag.String("messages", "./config/messages.json", "path to file with communication messages")

	gameConfig     config.GameConfig
	windowConfig   config.WindowConfig
	messages       config.Message
	mazeDimensions model.MazeDimensions

	imgFactory factory.AssetsFactory

	lives       = 3
	ghostsCount = 4
	isOver      = false

	fontFamily font.Face

	pacman     model.Sprite
	ghosts     []*model.Ghost
	maze       []string
	score      int
	dotsCount  int
	fruitTimer *time.Timer

	pacmanMtx sync.RWMutex
	fruitMtx  sync.Mutex
)

type Game struct {
	row int
	col int
}

func (g *Game) Update() error {
	unit := model.MazeCharacter{Row: g.row, Col: g.col}
	moves := model.CheckPossibleMoves(unit, maze)
	maze = model.ProcessPoint(maze, unit, &score)
	go model.ProcessTeleport(&pacman, mazeDimensions, &pacmanMtx)
	go model.ProcessInput(&pacman, moves, &pacmanMtx)
	go model.UpdateGhosts(&ghosts, model.MazeCharacter{Row: g.row, Col: g.col, Char: int32(maze[g.row][g.col])},
		&maze, fruitTimer, &fruitMtx, imgFactory, gameConfig)
	go model.ProcessGhostElimination(unit, &ghosts, &score, &lives, &isOver)
	go model.MoveGhosts(&ghosts, windowConfig, mazeDimensions)

	g.row = pacman.Y / windowConfig.CharSize
	g.col = pacman.X / windowConfig.CharSize

	if g.row > mazeDimensions.HeightLines-2 {
		g.row = mazeDimensions.HeightLines - 2
	} else if g.row < 1 {
		g.row = 1
	}

	if g.col > mazeDimensions.WidthLines-2 {
		g.col = mazeDimensions.WidthLines - 2
	} else if g.col < 1 {
		g.col = 1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if isOver {
		var msg string
		var fontColor color.Color
		if lives > 0 {
			msg = messages.Win
			fontColor = colornames.Gold
		} else {
			msg = messages.Lose
			fontColor = colornames.Red
		}

		for i := 0; i < len(msg); i++ {
			display := fmt.Sprintf(strings.ToUpper(string([]rune(msg)[i])))
			xMargin := ((windowConfig.CharSize * mazeDimensions.WidthLines) - (len(msg) * fontSize)) / 2
			text.Draw(screen, display, fontFamily, i*fontSize+xMargin,
				windowConfig.CharSize*mazeDimensions.HeightLines/2, fontColor)
		}

		return
	}

	for row, line := range maze {
		for col, char := range line {
			read := model.MazeCharacter{Row: row, Col: col, Line: line, Char: char}
			src.DrawMaze(screen, &read, &windowConfig, &imgFactory)
			model.DrawGhosts(screen, &read, &windowConfig, &imgFactory, &pacman, ghosts, &dotsCount)
		}
	}
	model.DrawLives(screen, imgFactory.Live, lives, model.MazeCharacter{
		Row: windowConfig.CharSize * mazeDimensions.WidthLines, Col: 0,
	})
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	if windowConfig.Height == 0 || windowConfig.Width == 0 {
		return windowConfig.CharSize * mazeDimensions.WidthLines, windowConfig.CharSize * mazeDimensions.HeightLines
	}
	return windowConfig.Width, windowConfig.Height
}

func main() {
	e := src.InitializeConfiguration(*gameConfigFile, *windowConfigFile, &gameConfig, &windowConfig)
	if e != nil {
		log.Fatal(e)

		return
	}

	e = factory.Create(&imgFactory)
	if e != nil {
		log.Panicf("[GAME] Assest load failed")
	}

	e = src.LoadMaze(*mazeLayoutFile, &maze, &ghosts, &pacman, ghostsCount,
		&imgFactory, &dotsCount, &mazeDimensions, &windowConfig)
	if e != nil {
		log.Panicf("[GAME] Maze load failed")
	}

	for _, ghost := range ghosts {
		ghost.Maze = src.LoadGhostsMaze(maze, ghost.Name)
	}

	e = loadFont()
	if e != nil {
		log.Printf("[GAME] Font load failed. %s", e)
	}

	e = config.LoadMessages(*messagesFile, &messages)
	if e != nil {
		log.Printf("[GAME] UI messages load failed. %s", e)
	}

	ebiten.SetWindowTitle(windowConfig.Name)
	ebiten.SetWindowIcon([]image.Image{imgFactory.Icon})
	if windowConfig.Height == 0 || windowConfig.Width == 0 {
		mazeDimensions = model.ConfigurePixels(&mazeDimensions, windowConfig)
		ebiten.SetWindowSize(windowConfig.CharSize*mazeDimensions.WidthLines, windowConfig.CharSize*mazeDimensions.HeightLines)
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

func loadFont() error {
	ttf, e := opentype.Parse(fonts.MPlus1pRegular_ttf)

	if e != nil {
		log.Fatal(e)

		return e
	}

	fontFamily, e = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     Dpi,
		Hinting: font.HintingFull,
	})

	if e != nil {
		log.Fatal(e)

		return e
	}

	return nil
}
