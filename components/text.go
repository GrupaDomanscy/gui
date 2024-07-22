package components

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextComponent struct {
	Component

	text     string
	fontName string
	fontSize float32
	spacing  float32
	color    rl.Color
	position rl.Vector2
}

func NewTextComponent(text string, loadedFontName string, fontSize float32, spacing float32, color rl.Color) *TextComponent {
	return &TextComponent{
		text:     text,
		fontName: loadedFontName,
		fontSize: fontSize,
		spacing:  spacing,
		color:    color,
		position: rl.Vector2Zero(),
	}
}

func (comp *TextComponent) Render(getFont GetFontCallback) {
	font, err := getFont(comp.fontName)
	if err != nil {
		panic(fmt.Sprintf("Provided font (%s) is not loaded into memory", comp.fontName, err))
	}

	rl.DrawTextEx(font, comp.text, comp.position, comp.fontSize, comp.spacing, comp.color)
}
