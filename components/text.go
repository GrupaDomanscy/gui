package components

import (
	"bytes"
	"encoding/gob"
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
	position ComponentPosition
}

func NewTextComponent(text string, loadedFontName string, fontSize float32, spacing float32, color rl.Color) *TextComponent {
	return &TextComponent{
		text:     text,
		fontName: loadedFontName,
		fontSize: fontSize,
		spacing:  spacing,
		color:    color,
		position: NewComponentPosition(),
	}
}

func (comp *TextComponent) SetPosition(pos rl.Vector2) {
	comp.position.Position = pos
}

func (comp *TextComponent) SetPositionOffset(offset rl.Vector2) {
	comp.position.Offset = offset
}

func (comp *TextComponent) CalculateSize(getFont GetFontCallback, maxViewport rl.Vector2) rl.Vector2 {
	font, err := getFont(comp.fontName)
	if err != nil {
		panic(fmt.Sprintf("Provided font (%s) is not loaded into memory", comp.fontName))
	}

	return rl.MeasureTextEx(font, comp.text, comp.fontSize, comp.spacing)
}

func (comp *TextComponent) Render(getFont GetFontCallback) {
	font, err := getFont(comp.fontName)
	if err != nil {
		panic(fmt.Sprintf("Provided font (%s) is not loaded into memory", comp.fontName))
	}

	rl.DrawTextEx(font, comp.text, comp.position.Calculate(), comp.fontSize, comp.spacing, comp.color)
}

func (comp *TextComponent) GetHash() []byte {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(comp)
	return b.Bytes()
}
