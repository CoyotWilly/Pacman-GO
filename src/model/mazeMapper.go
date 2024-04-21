package model

import (
	"Pacman/src/enum"
	"slices"
	"strings"
)

func Maze2MazeString(maze []string) string {
	var out string
	for i, line := range maze {
		if i == len(maze)-1 {
			out += line

			break
		}
		out += line + "\n"
	}

	return out
}

func MazeWithPath2Directions(path string, position MazeCharacter, movesCount int) []int {
	pathArray := strings.Split(strings.ReplaceAll(path, "\r\n", "\n"), "\n")
	var directions []int

	for i := 0; i < movesCount-1; i++ {
		var filtered []int // TODO fix moves generator
		moves, newPosition := CheckPossiblePaths(position, pathArray)

		if i > 0 && len(moves) > 1 {
			for _, move := range moves {
				if (move == enum.DOWN && directions[i-1] != enum.UP) ||
					(move == enum.UP && directions[i-1] != enum.DOWN) ||
					(move == enum.RIGHT && directions[i-1] != enum.LEFT) ||
					(move == enum.LEFT && directions[i-1] != enum.RIGHT) {
					filtered = append(filtered, move)
				}
			}
		} else {
			filtered = moves
		}

		if filtered[0] < 2 {
			newPosition.Col = position.Col
		} else if filtered[0] > 1 {
			newPosition.Row = position.Row
		}

		switch filtered[0] {
		case enum.UP:
			if slices.Contains(moves, enum.DOWN) {
				newPosition.Row--
			}
			break
		case enum.DOWN:
			if slices.Contains(moves, enum.UP) {
				newPosition.Row++
			}
			break
		case enum.LEFT:
			if slices.Contains(moves, enum.RIGHT) {
				newPosition.Col--
			}
			break
		case enum.RIGHT:
			if slices.Contains(moves, enum.LEFT) {
				newPosition.Col++
			}
			break
		}
		position = newPosition
		directions = append(directions, filtered[0])
	}

	return directions
}
