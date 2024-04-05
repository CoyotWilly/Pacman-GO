package model

import (
	"Pacman/src/config"
	"Pacman/src/enum"
	"Pacman/src/factory"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
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
		options.GeoM.Translate(float64(pacman.X),
			float64(pacman.Y))
	} else {
		options.GeoM.Translate(float64(unit.Col*windowConfig.CharSize),
			float64(unit.Row*windowConfig.CharSize))
	}

	screen.DrawImage(rect, options)
}

func drawDirection() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}

	return move[dir]
}
