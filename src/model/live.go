package model

import "github.com/hajimehoshi/ebiten/v2"

const (
	margin       = 8
	paddingRight = 162.5
)

func DrawLives(screen *ebiten.Image, icon *ebiten.Image, lives int, unit MazeCharacter) {

	for i := 0; i < lives; i++ {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(
			float64(unit.Row-(i*icon.Bounds().Dx())-((i+1)*margin+icon.Bounds().Dx()))-paddingRight,
			float64(unit.Col)+margin)
		screen.DrawImage(icon, options)
	}
}
