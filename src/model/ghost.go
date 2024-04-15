package model

import (
	"Pacman/src/config"
	"Pacman/src/enum"
	"Pacman/src/factory"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math"
	"slices"
	"sync"
	"time"
)

type Ghost struct {
	PositionLines  Sprite
	PositionPixels Sprite
	Shape          *ebiten.Image
	Status         enum.GhostsStatus
	Movement       Movement
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
	factory *factory.AssetsFactory, pacman *Sprite, ghosts []*Ghost, ghostsCount int, dotsCount *int) {
	rect := ebiten.NewImage(windowConfig.CharSize, windowConfig.CharSize)
	options := &ebiten.DrawImageOptions{}

	switch unit.Char {
	case enum.PACMAN:
		rect = factory.Pacman
		options.GeoM.Scale(windowConfig.ScaleFactor+0.1, windowConfig.ScaleFactor+0.1)
		break
	case enum.GHOST:
		rect = ghosts[unit.Col%ghostsCount].Shape
		options.GeoM.Scale(windowConfig.ScaleFactor+0.1, windowConfig.ScaleFactor+0.1)
		options.GeoM.Translate(float64(ghosts[unit.Col%ghostsCount].PositionPixels.X),
			float64(ghosts[unit.Col%ghostsCount].PositionPixels.Y))
		break
	case enum.POINT:
		*dotsCount++
		break
	default:
		return
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

func MoveGhosts(ghosts *[]*Ghost, maze []string, windowConfig config.WindowConfig, mazeDim MazeDimensions) {
	var updated []*Ghost
	for _, ghost := range *ghosts {
		if ghost.Movement.DirectionLock < windowConfig.CharSize && ghost.Movement.Direction != enum.UNDEFINED {
			go headTowardsGivenPoint(MazeCharacter{Row: 6, Col: 14}, ghost,
				[]int{ghost.Movement.Direction}, windowConfig, mazeDim)
		} else {
			ghost.Movement.DirectionLock = 0
			log.Printf("Lock: %d, Direction: %d", ghost.Movement.DirectionLock, ghost.Movement.Direction)
			moves := CheckPossibleMoves(MazeCharacter{Row: ghost.PositionLines.Y, Col: ghost.PositionLines.X}, maze)
			go headTowardsGivenPoint(MazeCharacter{Row: 6, Col: 14}, ghost, moves, windowConfig, mazeDim)
		}

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

func headTowardsGivenPoint(point MazeCharacter, ghost *Ghost, moves []int,
	conf config.WindowConfig, mazeDim MazeDimensions) {
	moveX := false
	if ghost.PositionLines.X > point.Row && slices.Contains(moves, enum.UP) {
		pixelMove(ghost, conf.CharSize*point.Row, enum.UP)
		moveX = true
	} else if ghost.PositionLines.X < point.Row && slices.Contains(moves, enum.DOWN) {
		pixelMove(ghost, conf.CharSize*point.Row, enum.DOWN)
		moveX = true
	}
	ghost.PositionLines.X = int(math.Round(float64(ghost.PositionPixels.X) / float64(conf.CharSize)))

	if moveX {
		ghost.PositionLines.Y = int(math.Round(float64(ghost.PositionPixels.Y) / float64(conf.CharSize)))

		return
	}

	if slices.Contains(moves, enum.LEFT) && ghost.PositionPixels.Y > conf.CharSize &&
		ghost.PositionLines.Y > point.Col {
		pixelMove(ghost, conf.CharSize*point.Col, enum.LEFT)
	} else if slices.Contains(moves, enum.RIGHT) && ghost.PositionPixels.Y < mazeDim.HeightPixels &&
		ghost.PositionLines.Y < point.Col {
		pixelMove(ghost, conf.CharSize*point.Col, enum.RIGHT)
	}
	ghost.PositionLines.Y = int(math.Round(float64(ghost.PositionPixels.Y) / float64(conf.CharSize)))
}

func pixelMove(ghost *Ghost, position int, direction int) {
	ghost.Movement.DirectionMtx.Lock()
	defer ghost.Movement.DirectionMtx.Unlock()

	if ghost.PositionPixels.X != position && direction > 1 {
		if direction == enum.RIGHT {
			ghost.PositionPixels.X++
		} else {
			ghost.PositionPixels.X--
		}
	} else if ghost.PositionPixels.Y != position && direction < 2 {
		if direction == enum.DOWN {
			ghost.PositionPixels.Y++
		} else {
			ghost.PositionPixels.Y--
		}
	}

	ghost.Movement.Direction = direction
	ghost.Movement.DirectionLock++
}
