package model

type Sprite struct {
	X             int
	Y             int
	XInit         int
	YInit         int
	Direction     int
	PrevDirection int
}

type Offset struct {
	X int
	Y int
}

func NewSprite(sprite Sprite) *Sprite {
	return &Sprite{
		X:     sprite.X,
		Y:     sprite.Y,
		XInit: sprite.XInit,
		YInit: sprite.YInit,
	}
}
