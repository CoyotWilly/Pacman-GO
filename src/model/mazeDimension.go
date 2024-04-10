package model

import "Pacman/src/config"

type MazeDimensions struct {
	WidthLines   int
	WidthPixels  int
	HeightLines  int
	HeightPixels int
	Offset       int
}

func ConfigurePixels(dim *MazeDimensions, window config.WindowConfig) MazeDimensions {
	return MazeDimensions{
		WidthLines:   dim.WidthLines,
		WidthPixels:  dim.WidthLines * window.CharSize,
		HeightLines:  dim.HeightLines,
		HeightPixels: dim.HeightLines * window.CharSize,
		Offset:       dim.Offset,
	}
}
