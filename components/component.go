package components

import rl "github.com/gen2brain/raylib-go/raylib"

type GetFontCallback = func(fontName string) (rl.Font, error)

type Component interface {
	Render(GetFontCallback)
	CalculateSize(maxViewport rl.Vector2) rl.Vector2
	SetPosition(rl.Vector2)
	SetPositionOffset(rl.Vector2)
}
