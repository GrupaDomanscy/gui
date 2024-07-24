package components

import rl "github.com/gen2brain/raylib-go/raylib"

type ComponentPosition struct {
	Offset   rl.Vector2
	Position rl.Vector2
}

func NewComponentPosition() ComponentPosition {
	return ComponentPosition{
		Offset:   rl.Vector2Zero(),
		Position: rl.Vector2Zero(),
	}
}

func (pos *ComponentPosition) Calculate() rl.Vector2 {
	return rl.Vector2Add(pos.Offset, pos.Position)
}
