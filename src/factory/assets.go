package factory

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"os"
	"strings"
)

const path = "./assets/"

type AssetsFactory struct {
	Pacman   *ebiten.Image
	Fruit    *ebiten.Image
	Dot      *ebiten.Image
	Blinky   *ebiten.Image
	Pinky    *ebiten.Image
	Inky     *ebiten.Image
	Clyde    *ebiten.Image
	Infected *ebiten.Image
	Icon     image.Image
}

func Create(consumer *AssetsFactory) error {
	files, e := os.ReadDir(path)

	if e != nil {
		log.Fatalf("[FACTORY] Could not read assets folder. %s", e)

		return e
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "-") {
			continue
		}

		img, ico, e := ebitenutil.NewImageFromFile(path + file.Name())
		if e != nil {
			log.Fatalf("[FACTORY] Image read failed. %s", e)
		}

		switch strings.Split(file.Name(), ".")[0] {
		case "blinky":
			consumer.Blinky = img
		case "pinky":
			consumer.Pinky = img
		case "inky":
			consumer.Inky = img
		case "clyde":
			consumer.Clyde = img
		case "pacman":
			consumer.Pacman = img
			consumer.Icon = ico
		case "fruit":
			consumer.Fruit = img
		case "dot":
			consumer.Dot = img
		case "infected":
			consumer.Infected = img
		default:
			log.Printf("[FACTORY] Unknown field: %s.", strings.Split(file.Name(), "."))
		}
	}

	return nil
}
