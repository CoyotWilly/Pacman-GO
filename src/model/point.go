package model

func ProcessPoint(maze []string, unit MazeCharacter, score *int) []string {
	var out []string
	if maze[unit.Row][unit.Col] == '.' {
		*score++
		for row, line := range maze {
			if row != unit.Row {
				out = append(out, line)

				continue
			}
			var newLine []uint8
			for col, _ := range line {
				if unit.Col == col {
					newLine = append(newLine, ' ')

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
