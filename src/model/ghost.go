package model

import (
	"Pacman/src/config"
	"Pacman/src/enum"
	"Pacman/src/factory"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Ghost struct {
	Position Sprite
	Shape    ebiten.Image
	Status   enum.GhostsStatus
}

func UpdateGhosts() {

}

func MoveGhosts() {

}

func DrawGhosts(screen *ebiten.Image, unit *MazeCharacter, windowConfig *config.WindowConfig,
	factory *factory.AssetsFactory, pacman *Sprite, ghosts []*Ghost, ghostsCount int, dotsCount *int) {
	rect := ebiten.NewImage(windowConfig.CharSize, windowConfig.CharSize)
	ghostsImg := []*ebiten.Image{factory.Pinky, factory.Blinky, factory.Inky, factory.Clyde}
	options := &ebiten.DrawImageOptions{}

	switch unit.Char {
	case 'P':
		rect = factory.Pacman
		options.GeoM.Scale(windowConfig.ScaleFactor+0.1, windowConfig.ScaleFactor+0.1)
	case 'G':
		if len(ghosts) < ghostsCount {
			ghosts = append(ghosts, &Ghost{
				Position: Sprite{X: unit.Row, Y: unit.Col, XInit: unit.Row, YInit: unit.Col},
				Status:   enum.Normal})
		} else {
			ghosts[unit.Col%ghostsCount].Position = Sprite{X: unit.Row, Y: unit.Col, XInit: unit.Row, YInit: unit.Col}
		}
		rect = ghostsImg[unit.Col%ghostsCount]
		options.GeoM.Scale(windowConfig.ScaleFactor+0.1, windowConfig.ScaleFactor+0.1)
	case '.':
		*dotsCount++
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

func CheckMovePossibility() {

}
