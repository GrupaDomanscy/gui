package components

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextComponent struct {
	Component

	text          string
	processedText string
	wrapText      bool
	fontName      string
	fontSize      float32
	spacing       float32
	color         rl.Color
	position      ComponentPosition
}

func NewTextComponent(text string, loadedFontName string, fontSize float32, spacing float32, color rl.Color) *TextComponent {
	return &TextComponent{
		text:          text,
		processedText: "",
		wrapText:      false,
		fontName:      loadedFontName,
		fontSize:      fontSize,
		spacing:       spacing,
		color:         color,
		position:      NewComponentPosition(),
	}
}

func getGlyphWidth(font *rl.Font, glyph int32, fontSize float32, cachedFontScalingFactor float32) float32 {
	var fontScalingFactor float32

	if cachedFontScalingFactor == nil {
		fontScalingFactor = float32(font.BaseSize) / fontSize
	} else {
		fontScalingFactor := *cachedFontScalingFactor
	}

	glyphInfo := rl.GetGlyphInfo(*font, glyph)

	if glyphInfo.AdvanceX != 0 {
		return float32(glyphInfo.AdvanceX)
	} else {
		return float32(glyphInfo.AdvanceX)
	}
}

func wrapText(font rl.Font, text string, fontSize float32, spacing float32, maxWidth float32) (processedText string, calculatedSize rl.Vector2) {
	fontScalingFactor := float32(font.BaseSize) / fontSize

	lastWhitespaceCharacterIndex := -1
	var currentLineText strings.Builder
	var currentLineWidth float32 = 0
	var defaultLineHeight float32 = fontSize // this is bullshit, we have for ex. arabic characters which break this rule

	for characterIndex, character := range text {
		if character == ' ' {
			lastWhitespaceCharacterIndex = characterIndex
		}

		currentLineText.WriteRune(character)

		currentLineWidth += glyphWidth
	}

	return
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

	if !comp.wrapText {
		return rl.MeasureTextEx(font, comp.text, comp.fontSize, comp.spacing)
	} else {
		var calculatedSize rl.Vector2

		comp.processedText, calculatedSize = wrapText(font, comp.text, comp.fontSize, comp.spacing, maxViewport.X)

		return calculatedSize
	}
}

func (comp *TextComponent) Render(getFont GetFontCallback) {
	font, err := getFont(comp.fontName)
	if err != nil {
		panic(fmt.Sprintf("Provided font (%s) is not loaded into memory", comp.fontName))
	}

	rl.DrawTextEx(font, comp.text, comp.position.Calculate(), comp.fontSize, comp.spacing, comp.color)
}

func (text *TextComponent) GetPosition() rl.Vector2 {
	return text.position.Calculate()
}
