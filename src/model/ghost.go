package model

import (
	"Pacman/src/config"
	"Pacman/src/enum"
	"Pacman/src/factory"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math"
	"sync"
	"time"
)

type Ghost struct {
	Position Sprite
	Shape    *ebiten.Image
	Status   enum.GhostsStatus
}

var (
	ghostUpdateMtx sync.RWMutex
)

func UpdateGhosts(ghosts *[]*Ghost, unit MazeCharacter, maze *[]string, fruitTimer *time.Timer, fruitMtx *sync.Mutex,
	factory factory.AssetsFactory, gameConfig config.GameConfig) {
	if unit.Char == 'X' {
		go changeStatus(ghosts, []*ebiten.Image{factory.Infected}, &ghostUpdateMtx)
		*maze = ConsumeFruit(unit, *maze)

		log.Printf("Countdown start: %s", time.Now().String())
		fruitMtx.Lock()
		defer fruitMtx.Unlock()
		fruitTimer = time.NewTimer(gameConfig.FruitDuration * time.Millisecond)
		<-fruitTimer.C
		changeStatus(ghosts, []*ebiten.Image{factory.Pinky, factory.Blinky, factory.Inky, factory.Clyde}, &ghostUpdateMtx)
		log.Printf("Countdown stop: %s", time.Now().String())
	}
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
		ghosts[unit.Col%ghostsCount].Position = Sprite{X: unit.Row, Y: unit.Col, XInit: unit.Row, YInit: unit.Col}
		rect = ghosts[unit.Col%ghostsCount].Shape
		options.GeoM.Scale(windowConfig.ScaleFactor+0.1, windowConfig.ScaleFactor+0.1)
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
	} else {
		options.GeoM.Translate(float64((unit.Col)*windowConfig.CharSize),
			float64(unit.Row*windowConfig.CharSize))
	}

	screen.DrawImage(rect, options)
}

func DrawDirection(towards ebiten.Key) int {
	move := map[ebiten.Key]int{
		ebiten.KeyW: 0,
		ebiten.KeyS: 1,
		ebiten.KeyD: 2,
		ebiten.KeyA: 3,
	}

	return move[towards]
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
