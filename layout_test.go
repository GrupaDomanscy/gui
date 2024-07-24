package gui

import (
	"fmt"
	"testing"

	. "domanscy.group/gui/components"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestLayoutComponent(t *testing.T) {
	testCases := []struct {
		direction          int
		mainAxisAlignment  int
		crossAxisAlignment int
		text1ExpectedPos   rl.Vector2
		text2ExpectedPos   rl.Vector2
		text3ExpectedPos   rl.Vector2
	}{
		{DirectionColumn, AlignStart, AlignStart, rl.Vector2{X: 0.0, Y: 0.0}, rl.Vector2{X: 0.0, Y: 32.0}, rl.Vector2{X: 0.0, Y: 64.0}},
		{DirectionColumn, AlignCenter, AlignStart, rl.Vector2{X: 0.0, Y: 300.0 - (96.0 / 2)}, rl.Vector2{X: 0.0, Y: 300.0 - (96.0 / 2) + 32.0}, rl.Vector2{X: 0.0, Y: 300.0 - (96.0 / 2) + 64.0}},
		{DirectionColumn, AlignEnd, AlignStart, rl.Vector2{X: 0.0, Y: 600.0 - 96.0}, rl.Vector2{X: 0.0, Y: 600.0 - 64.0}, rl.Vector2{X: 0.0, Y: 600.0 - 32.0}},
		{DirectionColumn, AlignStart, AlignCenter, rl.Vector2{X: 334.5, Y: 0}, rl.Vector2{X: 315.0, Y: 32.0}, rl.Vector2{X: 184.0, Y: 64.0}},
		{DirectionColumn, AlignStart, AlignEnd, rl.Vector2{X: 669.0, Y: 0}, rl.Vector2{X: 630.0, Y: 32.0}, rl.Vector2{X: 368.0, Y: 64.0}},

		{DirectionRow, AlignStart, AlignStart, rl.Vector2{X: 0.0, Y: 0.0}, rl.Vector2{X: 131.0, Y: 0.0}, rl.Vector2{X: 301.0, Y: 0.0}},
		{DirectionRow, AlignCenter, AlignStart, rl.Vector2{X: 33.5, Y: 0}, rl.Vector2{X: 164.5, Y: 0.0}, rl.Vector2{X: 334.5, Y: 0.0}},
		{DirectionRow, AlignEnd, AlignStart, rl.Vector2{X: 67.0, Y: 0}, rl.Vector2{X: 198.0, Y: 0.0}, rl.Vector2{X: 368.0, Y: 0}},
		{DirectionRow, AlignStart, AlignCenter, rl.Vector2{X: 0, Y: 284.0}, rl.Vector2{X: 131.0, Y: 284.0}, rl.Vector2{X: 301.0, Y: 284.0}},
		{DirectionRow, AlignStart, AlignEnd, rl.Vector2{X: 0.0, Y: 568.0}, rl.Vector2{X: 131.0, Y: 568.0}, rl.Vector2{X: 301.0, Y: 568.0}},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Subtest %d", i+1), func(t *testing.T) {
			layout := NewLayoutComponent(testCase.direction, testCase.mainAxisAlignment, testCase.crossAxisAlignment)
			text1 := NewTextComponent("Hello world", "Roboto", 32, 0, WhiteColor)
			text2 := NewTextComponent("Mumbo jambo", "Roboto", 32, 0, WhiteColor)
			text3 := NewTextComponent("DSAOIJDSAOIJSAKJNDSAKJNBDSA", "Roboto", 32, 0, WhiteColor)

			layout.AddChild(text1)
			layout.AddChild(text2)
			layout.AddChild(text3)

			BuildApp().
				WithTitle("Hello world").
				WithInitialSize(800, 600).
				WithFont("Roboto", "assets-for-testing\\Roboto-Regular.ttf").
				WithRootElement(layout).
				WithAppRoutine(func() (stop bool) {
					stop = true
					return
				}).
				Run()

			text1Pos := text1.GetPosition()
			text2Pos := text2.GetPosition()
			text3Pos := text3.GetPosition()

			if text1Pos.X != testCase.text1ExpectedPos.X {
				t.Errorf("text1Pos.X was expected to be %f, %f returned", testCase.text1ExpectedPos.X, text1Pos.X)
			}

			if text1Pos.Y != testCase.text1ExpectedPos.Y {
				t.Errorf("text1Pos.Y was expected to be %f, %f returned", testCase.text1ExpectedPos.Y, text1Pos.Y)
			}

			if text2Pos.X != testCase.text2ExpectedPos.X {
				t.Errorf("text2Pos.X was expected to be %f, %f returned", testCase.text2ExpectedPos.X, text2Pos.X)
			}

			if text2Pos.Y != testCase.text2ExpectedPos.Y {
				t.Errorf("text2Pos.Y was expected to be %f, %f returned", testCase.text2ExpectedPos.Y, text2Pos.Y)
			}

			if text3Pos.X != testCase.text3ExpectedPos.X {
				t.Errorf("text3Pos.X was expected to be %f, %f returned", testCase.text3ExpectedPos.X, text3Pos.X)
			}

			if text3Pos.Y != testCase.text3ExpectedPos.Y {
				t.Errorf("text3Pos.Y was expected to be %f, %f returned", testCase.text3ExpectedPos.Y, text3Pos.Y)
			}
		})
	}
}
