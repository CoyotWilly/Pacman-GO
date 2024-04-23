package model

import (
	"Pacman/src/config"
	"Pacman/src/enum"
	"Pacman/src/factory"
	"Pacman/src/pathfinder"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
	"unicode"
)

type Ghost struct {
	PositionLines  Sprite
	PositionPixels Sprite
	Name           enum.GhostsName
	Shape          *ebiten.Image
	Status         enum.GhostsStatus
	Movement       Movement
	Maze           []string
}

var (
	ghostUpdateMtx sync.RWMutex
)

func DrawDirection(towards ebiten.Key) int {
	move := map[ebiten.Key]int{
		ebiten.KeyW: 0,
		ebiten.KeyS: 1,
		ebiten.KeyD: 2,
		ebiten.KeyA: 3,
	}

	return move[towards]
}

func DrawGhosts(screen *ebiten.Image, unit *MazeCharacter, windowConfig *config.WindowConfig,
	factory *factory.AssetsFactory, pacman *Sprite, ghosts []*Ghost, dotsCount *int) {
	rect := ebiten.NewImage(windowConfig.CharSize, windowConfig.CharSize)
	options := &ebiten.DrawImageOptions{}
	ghostsMap := make(map[int32]*Ghost)

	for _, ghost := range ghosts {
		ghostsMap[pathfinder.Name2Rune(ghost.Name)] = ghost
	}

	switch unit.Char {
	case enum.PACMAN:
		rect = factory.Pacman
		options.GeoM.Scale(windowConfig.ScaleFactor+0.1, windowConfig.ScaleFactor+0.1)
		break
	case enum.BLINKY:
		break
	case enum.INKY:
		break
	case enum.PINKY:
		break
	case enum.CLYDE:
		break
	case enum.POINT:
		*dotsCount++
		break
	default:
		return
	}

	if unit.Char == enum.BLINKY || unit.Char == enum.INKY || unit.Char == enum.PINKY || unit.Char == enum.CLYDE {
		rect = ghostsMap[unit.Char].Shape
		options.GeoM.Scale(windowConfig.ScaleFactor+0.1, windowConfig.ScaleFactor+0.1)
		options.GeoM.Translate(float64(ghostsMap[unit.Char].PositionPixels.X),
			float64(ghostsMap[unit.Char].PositionPixels.Y))
	}

	if rect == factory.Pacman {
		switch pacman.Direction {
		case enum.UP:
			if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.DOWN {
				pacman.X -= windowConfig.CharSize
			} else if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.LEFT {
				pacman.X -= windowConfig.CharSize
				pacman.Y += windowConfig.CharSize
			} else if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.RIGHT {
				pacman.Y += windowConfig.CharSize
			}
			options.GeoM.Rotate(-math.Pi / 2)
			break
		case enum.DOWN:
			if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.UP {
				pacman.X += windowConfig.CharSize
				pacman.Y -= windowConfig.CharSize
			} else if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.RIGHT {
				pacman.X += windowConfig.CharSize
			}
			options.GeoM.Rotate(math.Pi / 2)
			break
		case enum.LEFT:
			if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.UP {
				pacman.X += windowConfig.CharSize
				pacman.Y -= windowConfig.CharSize
			} else if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.RIGHT {
				pacman.X += windowConfig.CharSize
			}
			options.GeoM.Rotate(math.Pi)
			options.GeoM.Scale(1, -1)
			break
		case enum.RIGHT:
			if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.UP {
				pacman.Y -= windowConfig.CharSize
			} else if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.DOWN {
				pacman.X -= windowConfig.CharSize
			} else if pacman.Direction != pacman.PrevDirection && pacman.PrevDirection == enum.LEFT {
				pacman.X -= windowConfig.CharSize
			}
		}

		options.GeoM.Translate(float64(pacman.X), float64(pacman.Y))
		pacman.PrevDirection = pacman.Direction
	}

	screen.DrawImage(rect, options)
}

func UpdateGhosts(ghosts *[]*Ghost, unit MazeCharacter, maze *[]string, fruitTimer *time.Timer, fruitMtx *sync.Mutex,
	factory factory.AssetsFactory, gameConfig config.GameConfig) {
	if unit.Char == enum.FRUIT {
		go changeStatus(ghosts, []*ebiten.Image{factory.Infected}, &ghostUpdateMtx)
		*maze = ConsumeFruit(unit, *maze)

		fruitMtx.Lock()
		defer fruitMtx.Unlock()
		fruitTimer = time.NewTimer(gameConfig.FruitDuration * time.Millisecond)
		<-fruitTimer.C
		changeStatus(ghosts, []*ebiten.Image{factory.Pinky, factory.Blinky, factory.Inky, factory.Clyde}, &ghostUpdateMtx)
	}
}

func MoveGhosts(ghosts *[]*Ghost, windowConfig config.WindowConfig, mazeDim MazeDimensions) {
	var updated []*Ghost
	for _, ghost := range *ghosts {
		if ghost.Movement.DirectionCounter == len(ghost.Movement.Directions)-1 && ghost.Movement.Directions != nil {
			ghost.Maze = changeGhostMarker(&ghost.Maze, ghost.Name, mazeDim)
			ghost.Movement.Directions = nil
		}

		if ghost.Movement.Directions == nil {
			directions := generatePath(ghost)
			ghost.Movement = Movement{
				DirectionCounter: 0,
				Directions:       directions,
				DirectionLock:    directions[0],
			}
		}

		if ghost.Movement.DirectionLock < windowConfig.CharSize {
			go pixelMove(ghost, ghost.Movement.Directions[ghost.Movement.DirectionCounter])
		} else if ghost.Movement.DirectionCounter < len(ghost.Movement.Directions) &&
			len(ghost.Movement.Directions) != 0 {
			ghost.Movement.DirectionLock = 0
			ghost.Movement.DirectionCounter++
			go pixelMove(ghost, ghost.Movement.Directions[ghost.Movement.DirectionCounter])
		}

		ghost.PositionLines.X = int(math.Round(float64(ghost.PositionPixels.X) / float64(windowConfig.CharSize)))
		ghost.PositionLines.Y = int(math.Round(float64(ghost.PositionPixels.Y) / float64(windowConfig.CharSize)))
		updated = append(updated, ghost)
	}
	ghosts = &updated
}

func changeStatus(ghosts *[]*Ghost, assets []*ebiten.Image, mtx *sync.RWMutex) {
	mtx.Lock()
	defer mtx.Unlock()
	for i, ghost := range *ghosts {
		if len(assets) == 1 {
			ghost.Shape = assets[0]
			ghost.Status = enum.Infected

			continue
		}

		ghost.Shape = assets[i]
		ghost.Status = enum.Normal
	}
}

func pixelMove(ghost *Ghost, direction int) {
	ghost.Movement.DirectionMtx.Lock()
	defer ghost.Movement.DirectionMtx.Unlock()

	if direction > 1 {
		if direction == enum.RIGHT {
			ghost.PositionPixels.X++
		} else {
			ghost.PositionPixels.X--
		}
	} else {
		if direction == enum.DOWN {
			ghost.PositionPixels.Y++
		} else {
			ghost.PositionPixels.Y--
		}
	}

	ghost.Movement.DirectionLock++
}

func changeGhostMarker(maze *[]string, ghostName enum.GhostsName, dim MazeDimensions) []string {
	nameChar := pathfinder.Name2Rune(ghostName)
	newX := enum.UNDEFINED
	newY := enum.UNDEFINED
	var newMaze []string
	var newMazeTemp []string

	// SEARCHING FOR NEW DESTINATION
	for {
		x, y := rand.Int()%dim.WidthLines, rand.Int()%dim.HeightLines
		var tempMaze []string

		for i, row := range *maze {
			var tempRow string
			for j, c := range row {
				if i == x && j == y && c == enum.EMPTY {
					tempRow += string(pathfinder.Name2Rune(enum.NoName))
				} else {
					tempRow += string(c)
				}
			}
			tempMaze = append(tempMaze, tempRow)
		}

		for i, row := range tempMaze {
			if i != x {
				continue
			}

			for j, c := range row {
				if i == x && j == y && c == pathfinder.Name2Rune(enum.NoName) {
					world := pathfinder.ParseWorld(Maze2MazeString(*maze))
					_, _, found := pathfinder.Path(world.From(ghostName), world.To(enum.NoName))
					if found {
						newX = x
						newY = y
						break
					}
				}
			}
		}

		if newX != enum.UNDEFINED && newY != enum.UNDEFINED {
			newMazeTemp = tempMaze
			break
		}
	}
	// New maze config creation
	for i, row := range newMazeTemp {
		var newRow string
		for j, c := range row {
			stringChar := c
			if c == nameChar {
				stringChar = enum.EMPTY
			} else if c == unicode.ToUpper(nameChar) {
				stringChar = nameChar
			} else if i == newX && j == newY {
				stringChar = unicode.ToUpper(nameChar)
			}
			newRow += string(stringChar)
		}

		newMaze = append(newMaze, newRow)
	}

	return newMaze
}

func generatePath(ghost *Ghost) []int {
	var fixedMaze []string
	for i, row := range ghost.Maze {
		var fixedRow string
		for j, c := range row {
			if i == ghost.PositionLines.Y && j == ghost.PositionLines.X {
				fixedRow += string(pathfinder.Name2Rune(ghost.Name))
			} else if c == pathfinder.Name2Rune(ghost.Name) {
				fixedRow += string(enum.EMPTY)
			} else {
				fixedRow += string(c)
			}
		}
		fixedMaze = append(fixedMaze, fixedRow)
	}
	ghost.Maze = fixedMaze

	world := pathfinder.ParseWorld(Maze2MazeString(ghost.Maze))
	p, _, found := pathfinder.Path(world.From(ghost.Name), world.To(ghost.Name))
	if !found {
		log.Panic("Could not find a path")
	} else {
		sMaze := world.RenderPath(p)

		return MazeWithPath2Directions(sMaze, MazeCharacter{
			Row: ghost.PositionLines.Y,
			Col: ghost.PositionLines.X,
		}, pathfinder.EstimateDistance(sMaze))
	}

	return nil
}
