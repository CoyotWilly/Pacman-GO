package model

import (
	"Pacman/src/enum"
	"math/rand"
)

type Ghost struct {
	Position Sprite
	Status   enum.GhostsStatus
}

func UpdateGhosts() {

}

func MoveGhosts() {

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
