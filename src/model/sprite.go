package model

import "sync"

type Sprite struct {
	X             int
	Y             int
	XInit         int
	YInit         int
	Direction     int
	PrevDirection int
}

type Movement struct {
	DirectionCounter int
	Directions       []int
	DirectionLock    int
	DirectionMtx     sync.RWMutex
}
