package pathfinder

import (
	"Pacman/src/enum"
	"fmt"
	"strings"
	"unicode"
)

// Kind* constants refer to tile kinds for input and output.
const (
	// KindPlain (.) is a plain tile with a movement cost of 1.
	KindPlain = iota
	// KindRiver (~) is a river tile with a movement cost of 2.
	KindRiver
	// KindMountain (M) is a mountain tile with a movement cost of 3.
	KindMountain
	// KindBlocker (X) is a tile which blocks movement.
	KindBlocker
	KindPlayer
	// KindFrom (F) is a tile which marks where the path should be calculated
	// from.
	KindFrom
	KindFrom1
	KindFrom2
	KindFrom3
	KindFrom4
	// KindTo (T) is a tile which marks the goal of the path.
	KindTo
	KindTo1
	KindTo2
	KindTo3
	KindTo4
	// KindPath (●) is a tile to represent where the path is in the output.
	KindPath
)

// KindRunes map tile kinds to output runes.
var KindRunes = map[int]rune{
	KindPlain:    ' ',
	KindRiver:    '.',
	KindMountain: 'X',
	KindBlocker:  '#',
	KindPlayer:   'M',
	KindFrom:     'g',
	KindFrom1:    'p',
	KindFrom2:    'i',
	KindFrom3:    'b',
	KindFrom4:    'c',
	KindTo:       'G',
	KindTo1:      'P',
	KindTo2:      'I',
	KindTo3:      'B',
	KindTo4:      'C',
	KindPath:     '*',
}

var ghostName2Rune = map[enum.GhostsName]rune{
	enum.Pinky:  'p',
	enum.Inky:   'i',
	enum.Blinky: 'b',
	enum.Clyde:  'c',
	enum.NoName: 'g',
}

// RuneKinds map input runes to tile kinds.
var RuneKinds = map[rune]int{
	' ': KindPlain,
	'.': KindRiver,
	'X': KindMountain,
	'#': KindBlocker,
	'M': KindPlayer,
	'p': KindFrom,
	'i': KindFrom,
	'b': KindFrom,
	'c': KindFrom,
	'g': KindFrom,
	'G': KindTo,
	'B': KindTo,
	'I': KindTo,
	'P': KindTo,
	'C': KindTo,
}

// KindCosts map tile kinds to movement costs.
var KindCosts = map[int]float64{
	KindPlayer:   1.0,
	KindPlain:    1.0,
	KindFrom:     1.0,
	KindTo:       1.0,
	KindRiver:    2.0,
	KindMountain: 3.0,
}

// A Tile is a tile in a grid which implements Pattern.
type Tile struct {
	// Kind is the kind of tile, potentially affecting movement.
	Kind int
	// X and Y are the coordinates of the tile.
	X, Y int
	// W is a reference to the World that the tile is a part of.
	W World
}

func Name2Rune(name enum.GhostsName) rune {
	return ghostName2Rune[name]
}

// PathNeighbors returns the neighbors of the tile, excluding blockers and
// tiles off the edge of the board.
func (t *Tile) PathNeighbors() []Pattern {
	var neighbors []Pattern
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if n := t.W.Tile(t.X+offset[0], t.Y+offset[1]); n != nil &&
			n.Kind != KindBlocker {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

// PathNeighborCost returns the movement cost of the directly neighboring tile.
func (t *Tile) PathNeighborCost(to Pattern) float64 {
	toT := to.(*Tile)
	return KindCosts[toT.Kind]
}

// PathEstimatedCost uses Manhattan distance to estimate orthogonal distance
// between non-adjacent nodes.
func (t *Tile) PathEstimatedCost(to Pattern) float64 {
	toT := to.(*Tile)
	absX := toT.X - t.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - t.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}

// World is a two dimensional map of Tiles.
type World map[int]map[int]*Tile

// Tile gets the tile at the given coordinates in the world.
func (w World) Tile(x, y int) *Tile {
	if w[x] == nil {
		return nil
	}
	return w[x][y]
}

// SetTile sets a tile at the given coordinates in the world.
func (w World) SetTile(t *Tile, x, y int) {
	if w[x] == nil {
		w[x] = map[int]*Tile{}
	}
	w[x][y] = t
	t.X = x
	t.Y = y
	t.W = w
}

// FirstOfKind gets the first tile on the board of a kind, used to get the from
// and to tiles as there should only be one of each.
func (w World) FirstOfKind(kind int) *Tile {
	for _, row := range w {
		for _, t := range row {
			if t.Kind == kind {
				return t
			}
		}
	}
	return nil
}

// From gets the from tile from the world.
func (w World) From(ghostName enum.GhostsName) *Tile {
	return w.FirstOfKind(RuneKinds[ghostName2Rune[ghostName]])
}

// To gets the to tile from the world.
func (w World) To(ghostName enum.GhostsName) *Tile {
	return w.FirstOfKind(RuneKinds[unicode.ToUpper(ghostName2Rune[ghostName])])
}

// RenderPath renders a path on top of a world.
func (w World) RenderPath(path []Pattern) string {
	width := len(w)
	if width == 0 {
		return ""
	}
	height := len(w[0])
	pathLocs := map[string]bool{}
	for _, p := range path {
		pT := p.(*Tile)
		pathLocs[fmt.Sprintf("%d,%d", pT.X, pT.Y)] = true
	}
	rows := make([]string, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			t := w.Tile(x, y)
			r := ' '
			if pathLocs[fmt.Sprintf("%d,%d", x, y)] {
				r = KindRunes[KindPath]
			} else if t != nil {
				r = KindRunes[t.Kind]
			}
			rows[y] += string(r)
		}
	}
	return strings.Join(rows, "\n")
}

// ParseWorld parses a textual representation of a world into a world map.
func ParseWorld(input string) World {
	w := World{}
	for y, row := range strings.Split(strings.TrimSpace(input), "\n") {
		for x, raw := range row {
			kind, ok := RuneKinds[raw]
			if !ok {
				kind = KindBlocker
			}
			w.SetTile(&Tile{
				Kind: kind,
			}, x, y)
		}
	}
	return w
}
