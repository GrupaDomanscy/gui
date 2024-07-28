package components

import (
	"domanscy.group/gui/components/atoms"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GetFontCallback = func(fontName string) (rl.Font, error)

type Component interface {
	Render(GetFontCallback)
	CalculateSize(getFont GetFontCallback, maxViewport rl.Vector2) rl.Vector2
	SetPosition(rl.Vector2)
	SetPositionOffset(rl.Vector2)
	GetPosition() rl.Vector2

	GetEventBus() *atoms.EventStore
	PropagateEvent(eventType string, args ...interface{})
}
