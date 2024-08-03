package gui

import (
	"domanscy.group/gui/components/atoms"
	"domanscy.group/gui/components/events"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var keysToCheck []int32 = []int32{
	32,
	256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269,
	280, 281, 282, 283, 284,
	290, 291, 292, 293, 294, 295, 296, 297, 298, 299, 300, 301,
	340, 341, 342, 343, 344, 345, 346, 347, 348,
	91, 92, 93,
	96,
	320, 321, 322, 323, 324, 325, 326, 327, 328, 329, 330, 331, 332, 333, 334, 335, 336,
	39,
	44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64,
	65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
	80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
}

func processKeyEvents(eventBus *atoms.EventBus) {
	// im really really sorry.
	// afaik there is no other way. GetCharPressed and GetKeyPressed are functions which remove
	// keys from internal array, but at one frame key must be checked with both of these functions.
	// This makes it impossible to program the elegant and efficient way.
	//
	// ~110 iterations to the moon i guess

	pressedChars := make([]rune, 0)

	var pressedChar rune = rl.GetCharPressed()

	for pressedChar != 0 {
		pressedChars = append(pressedChars, pressedChar)

		pressedChar = rl.GetCharPressed()
	}

	pressedKeys := make([]int32, 0)
	downKeys := make([]int32, 0)
	releasedKeys := make([]int32, 0)

	for _, key := range keysToCheck {
		keyPressed := rl.IsKeyPressed(key)
		keyDown := rl.IsKeyDown(key)
		keyReleased := rl.IsKeyReleased(key)

		if keyDown {
			downKeys = append(downKeys, key)
		}

		if keyPressed {
			pressedKeys = append(pressedKeys, key)
		}

		if keyReleased {
			releasedKeys = append(releasedKeys, key)
		}
	}

	if len(pressedChars) > 0 || len(downKeys) > 0 || len(pressedKeys) > 0 || len(releasedKeys) > 0 {
		eventBus.DispatchEvent("gui:keyaction", events.KeyPressEvent{
			PressedChars: pressedChars,
			DownKeys:     downKeys,
			PressedKeys:  pressedKeys,
			ReleasedKeys: releasedKeys,
		})
	}
}
