package model

import "Pacman/src/enum"

func ProcessPoint(maze []string, unit MazeCharacter, score *int) []string {
	var out []string
	if maze[unit.Row][unit.Col] == enum.POINT {
		*score++
		for row, line := range maze {
			if row != unit.Row {
				out = append(out, line)

				continue
			}
			var newLine []uint8
			for col, _ := range line {
				if unit.Col == col {
					newLine = append(newLine, enum.EMPTY)

					continue
				}
				newLine = append(newLine, maze[row][col])
			}
			out = append(out, string(newLine))
		}

		return out
	}

	return maze
}
