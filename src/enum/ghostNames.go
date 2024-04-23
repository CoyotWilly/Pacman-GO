package enum

type GhostsName string

const (
	NoName GhostsName = "Unnamed"
	Pinky  GhostsName = "Pinky"
	Inky   GhostsName = "Inky"
	Blinky GhostsName = "Blinky"
	Clyde  GhostsName = "Clyde"
)

func Rune2GhostName(char int32) GhostsName {
	mapper := map[int32]GhostsName{
		'p': Pinky,
		'i': Inky,
		'b': Blinky,
		'c': Clyde,
		'g': NoName,
	}

	return mapper[char]
}
