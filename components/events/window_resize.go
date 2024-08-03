package events

import rl "github.com/gen2brain/raylib-go/raylib"

type WindowResizedEvent struct {
	OldWindowSize rl.Vector2
	NewWindowSize rl.Vector2
}
