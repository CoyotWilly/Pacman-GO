package src

import (
	"Pacman/src/config"
	"Pacman/src/enum"
	"Pacman/src/factory"
	"Pacman/src/model"
	"bufio"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
	"os"
)

const (
	VerticalOffset   = 1
	HorizontalOffset = 1.5
)

type MazeDimensions struct {
	Width  int
	Height int
}

func LoadMaze(
	filePath string, maze *[]string, ghosts *[]*model.Ghost,
	pacman *model.Sprite, dotsCount *int, ghostsCount int, dim *MazeDimensions, window *config.WindowConfig) error {
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

	for row, line := range *maze {
		for col, char := range line {
			dim.Width = col
			switch char {
			case 'P':
				x := int((float64(row) + HorizontalOffset) * float64(window.CharSize))
				y := (col - 1) * window.CharSize
				*pacman = model.Sprite{X: x, Y: y, XInit: x, YInit: y}
			case 'G':
				if len(*ghosts) < ghostsCount {
					*ghosts = append(*ghosts,
						&model.Ghost{
							Position: model.Sprite{X: row, Y: col, XInit: row, YInit: col},
							Status:   enum.Normal,
						})
				}
			case '.':
				*dotsCount++
			}
			dim.Height = row
		}
	}
	dim.Width++
	dim.Height++

	return nil
}

func DrawMaze(screen *ebiten.Image, unit *model.MazeCharacter,
	windowConfig *config.WindowConfig, factory *factory.AssetsFactory) {
	rect := ebiten.NewImage(windowConfig.CharSize, windowConfig.CharSize)
	options := &ebiten.DrawImageOptions{}
	offset := float64(0)

	switch unit.Char {
	case '#':
		rect.Fill(color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff})
	case '.':
		rect = factory.Dot
		options.GeoM.Scale(windowConfig.ScaleFactor, windowConfig.ScaleFactor)
		offset = float64(windowConfig.CharSize / 2)
	case 'X':
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
