package model

import (
	"Pacman/src/config"
	"Pacman/src/enum"
	"github.com/hajimehoshi/ebiten/v2"
	"slices"
)

func ProcessInput(pacman *Sprite, moves []int) {
	if ebiten.IsKeyPressed(ebiten.KeyW) && slices.Contains(moves, enum.UP) {
		pacman.Y--
		pacman.Direction = DrawDirection(ebiten.KeyW)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) && slices.Contains(moves, enum.DOWN) {
		pacman.Y++
		pacman.Direction = DrawDirection(ebiten.KeyS)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) && slices.Contains(moves, enum.RIGHT) {
		pacman.X++
		pacman.Direction = DrawDirection(ebiten.KeyD)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) && slices.Contains(moves, enum.LEFT) {
		pacman.X--
		pacman.Direction = DrawDirection(ebiten.KeyA)
	}
}

func CheckPossibleMoves(pacman Sprite, windowConfig config.WindowConfig, mazeDimensions MazeDimensions,
	position MazeCharacter, maze []string) []int {
	var moves []int
	if pacman.Y >= windowConfig.CharSize && maze[position.Row-1][position.Col] != enum.WALL {
		moves = append(moves, enum.UP)
	}
	if pacman.Y <= windowConfig.CharSize*mazeDimensions.Height && maze[position.Row+1][position.Col] != enum.WALL {
		moves = append(moves, enum.DOWN)
	}
	if pacman.X >= windowConfig.CharSize && maze[position.Row][position.Col+1] != enum.WALL {
		moves = append(moves, enum.RIGHT)
	}
	if pacman.X <= windowConfig.CharSize*mazeDimensions.Width && maze[position.Row][position.Col-1] != enum.WALL {
		moves = append(moves, enum.LEFT)
	}

	return moves
}
