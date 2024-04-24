package model

import (
	"Pacman/src/enum"
	"github.com/hajimehoshi/ebiten/v2"
	"slices"
	"sync"
)

const (
	GhostEliminationReward = 10
)

func ProcessInput(pacman *Sprite, moves []int, mtx *sync.RWMutex) {
	if ebiten.IsKeyPressed(ebiten.KeyW) && slices.Contains(moves, enum.UP) {
		mtx.Lock()
		pacman.Y--
		pacman.Direction = DrawDirection(ebiten.KeyW)
		mtx.Unlock()
	} else if ebiten.IsKeyPressed(ebiten.KeyS) && slices.Contains(moves, enum.DOWN) {
		mtx.Lock()
		pacman.Y++
		pacman.Direction = DrawDirection(ebiten.KeyS)
		mtx.Unlock()
	} else if ebiten.IsKeyPressed(ebiten.KeyD) && slices.Contains(moves, enum.RIGHT) {
		mtx.Lock()
		pacman.X++
		pacman.Direction = DrawDirection(ebiten.KeyD)
		mtx.Unlock()
	} else if ebiten.IsKeyPressed(ebiten.KeyA) && slices.Contains(moves, enum.LEFT) {
		mtx.Lock()
		pacman.X--
		pacman.Direction = DrawDirection(ebiten.KeyA)
		mtx.Unlock()
	}
}

func ProcessTeleport(pacman *Sprite, mazeDim MazeDimensions, mtx *sync.RWMutex) {
	if pacman.X > mazeDim.WidthPixels {
		mtx.Lock()
		pacman.X = 0
		mtx.Unlock()
	} else if pacman.X < 1 {
		mtx.Lock()
		pacman.X = mazeDim.WidthPixels
		mtx.Unlock()
	}

}

func ProcessGhostElimination(unit MazeCharacter, ghosts *[]*Ghost, score *int,
	lives *int, isOver *bool) {
	for _, ghost := range *ghosts {
		if ghost.PositionLines.X == unit.Col && ghost.PositionLines.Y == unit.Row {
			if ghost.Status == enum.Infected {
				ghost.PositionLines = Sprite{
					X:     ghost.PositionLines.XInit,
					Y:     ghost.PositionLines.YInit,
					XInit: ghost.PositionLines.XInit,
					YInit: ghost.PositionLines.YInit,
				}
				*score += GhostEliminationReward

			} else {
				*lives--
				if *lives < 1 {
					*isOver = true
				}
			}
		}
	}
}

func CheckPossibleMoves(position MazeCharacter, maze []string) []int {
	var moves []int
	if maze[position.Row-1][position.Col] != enum.WALL {
		moves = append(moves, enum.UP)
	}
	if maze[position.Row+1][position.Col] != enum.WALL {
		moves = append(moves, enum.DOWN)
	}
	if maze[position.Row][position.Col+1] != enum.WALL {
		moves = append(moves, enum.RIGHT)
	}
	if maze[position.Row][position.Col-1] != enum.WALL {
		moves = append(moves, enum.LEFT)
	}

	return moves
}

func CheckPossiblePaths(position MazeCharacter, maze []string) ([]int, MazeCharacter) {
	var moves []int
	newPosition := position
	if maze[position.Row-1][position.Col] == enum.PATH {
		moves = append(moves, enum.UP)
		newPosition.Row--
	}
	if maze[position.Row+1][position.Col] == enum.PATH {
		moves = append(moves, enum.DOWN)
		newPosition.Row++
	}

	if maze[position.Row][position.Col+1] == enum.PATH {
		moves = append(moves, enum.RIGHT)
		newPosition.Col++
	}
	if maze[position.Row][position.Col-1] == enum.PATH {
		moves = append(moves, enum.LEFT)
		newPosition.Col--
	}

	return moves, newPosition
}
