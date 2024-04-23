package src

import (
	"Pacman/src/config"
	"Pacman/src/enum"
	"Pacman/src/factory"
	"Pacman/src/model"
	"Pacman/src/pathfinder"
	"bufio"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
	"os"
	"slices"
	"unicode"
)

const (
	HorizontalOffset = 1.5
)

func LoadMaze(
	filePath string, maze *[]string, ghosts *[]*model.Ghost, pacman *model.Sprite, ghostsCount int,
	factory *factory.AssetsFactory, dotsCount *int, dim *model.MazeDimensions, window *config.WindowConfig) error {
	file, e := os.Open(filePath)
	if e != nil {
		return e
	}

	defer func(file *os.File) {
		e = file.Close()
		if e != nil {
			log.Fatal("[MAZE] Could not close file.")
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		*maze = append(*maze, line)
	}
	ghostsImg := map[int32]*ebiten.Image{'p': factory.Pinky, 'b': factory.Blinky, 'i': factory.Inky, 'c': factory.Clyde}
	for row, line := range *maze {
		for col, char := range line {
			dim.WidthLines = col
			switch char {
			case enum.PACMAN:
				x := int((float64(row) + HorizontalOffset) * float64(window.CharSize))
				y := col * window.CharSize
				*pacman = model.Sprite{X: x, Y: y, XInit: x, YInit: y}
			case enum.POINT:
				*dotsCount++
			}

			if len(*ghosts) < ghostsCount &&
				(char == enum.BLINKY || char == enum.INKY || char == enum.PINKY || char == enum.CLYDE) {
				*ghosts = append(*ghosts,
					&model.Ghost{
						PositionLines: model.Sprite{X: col, Y: row, XInit: col, YInit: row},
						PositionPixels: model.Sprite{X: col * window.CharSize, Y: row * window.CharSize,
							XInit: col * window.CharSize, YInit: row * window.CharSize},
						Shape:  ghostsImg[char],
						Status: enum.Normal,
						Name:   enum.Rune2GhostName(char),
						Movement: model.Movement{
							DirectionCounter: enum.UNDEFINED,
							DirectionLock:    enum.UNDEFINED,
						},
					})
			}

			dim.HeightLines = row
		}
	}
	dim.WidthLines++
	dim.HeightLines++

	return nil
}

func DrawMaze(screen *ebiten.Image, unit *model.MazeCharacter,
	windowConfig *config.WindowConfig, factory *factory.AssetsFactory) {
	rect := ebiten.NewImage(windowConfig.CharSize, windowConfig.CharSize)
	options := &ebiten.DrawImageOptions{}
	offset := 0.0

	switch unit.Char {
	case enum.WALL:
		rect.Fill(color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff})
	case enum.POINT:
		rect = factory.Dot
		options.GeoM.Scale(windowConfig.ScaleFactor, windowConfig.ScaleFactor)
		offset = float64(windowConfig.CharSize / 2)
	case enum.FRUIT:
		rect = factory.Fruit
		options.GeoM.Scale(windowConfig.ScaleFactor, windowConfig.ScaleFactor)
		offset = (float64(windowConfig.CharSize) * windowConfig.ScaleFactor) / 2
	default:
		return
	}
	options.GeoM.Translate(float64(unit.Col*windowConfig.CharSize)+offset,
		float64(unit.Row*windowConfig.CharSize)+offset)

	screen.DrawImage(rect, options)
}

func LoadGhostsMaze(maze []string, name enum.GhostsName) []string {
	var ghostsMaze []string
	char := pathfinder.Name2Rune(name)
	chars := []int32{pathfinder.Name2Rune(enum.Pinky), unicode.ToUpper(pathfinder.Name2Rune(enum.Pinky)),
		pathfinder.Name2Rune(enum.Inky), unicode.ToUpper(pathfinder.Name2Rune(enum.Inky)),
		pathfinder.Name2Rune(enum.Blinky), unicode.ToUpper(pathfinder.Name2Rune(enum.Blinky)),
		pathfinder.Name2Rune(enum.Clyde), unicode.ToUpper(pathfinder.Name2Rune(enum.Clyde))}

	for _, row := range maze {
		var rowString string
		for _, c := range row {
			if c == enum.POINT || (slices.Contains(chars, c) && (unicode.ToUpper(c) != unicode.ToUpper(char))) {
				rowString += string(enum.EMPTY)
			} else {
				rowString += string(c)
			}
		}

		ghostsMaze = append(ghostsMaze, rowString)
	}

	return ghostsMaze
}
