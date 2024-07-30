package gui

import rl "github.com/gen2brain/raylib-go/raylib"

type WindowResizedEventArgs struct {
	oldWindowSize rl.Vector2
	newWindowSize rl.Vector2
}
