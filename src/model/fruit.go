package model

import "Pacman/src/enum"

func ConsumeFruit(unit MazeCharacter, maze []string) []string {
	var out []string
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
