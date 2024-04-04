package src

import (
	"Pacman/src/enum"
	"Pacman/src/model"
	"bufio"
	"log"
	"os"
)

type MazeDimensions struct {
	Width  int
	Height int
}

func LoadMaze(
	filePath string, maze *[]string, ghosts *[]*model.Ghost,
	pacman *model.Sprite, dotsCount *int, ghostsCount int, dim *MazeDimensions) error {
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
				*pacman = model.Sprite{X: row, Y: col, XInit: row, YInit: col}
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
