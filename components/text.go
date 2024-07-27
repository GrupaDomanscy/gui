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

func (component *TextComponent) SetWrapText(wrap bool) {
	component.wrapText = wrap
}

func getGlyphWidth(font *rl.Font, glyph int32) float32 {
	glyphInfo := rl.GetGlyphInfo(*font, glyph)

	if glyphInfo.AdvanceX != 0 {
		return float32(glyphInfo.AdvanceX)
	} else {
		return (*font).Recs.Width + float32(glyphInfo.OffsetX)
	}
}

func wrapText(font rl.Font, text string, fontSize float32, spacing float32, maxWidth float32) (processedText string, calculatedSize rl.Vector2) {
	type IterationState struct {
		wrappedText                  strings.Builder
		currentLineText              strings.Builder
		currentLineWidth             float32
		defaultLineHeight            float32 // this is bullshit, we have for ex. arabic characters which break this rule
		lastWhitespaceCharacterIndex int
	}

	fontScalingFactor := float32(font.BaseSize) / fontSize

	processEndOfInput := func(state *IterationState) {
		state.wrappedText.WriteString(state.currentLineText.String())
		calculatedSize.Y += state.defaultLineHeight

		if state.currentLineWidth > calculatedSize.X {
			calculatedSize.X = state.currentLineWidth * fontScalingFactor
		}
	}

	processEndOfLine := func(state *IterationState) {
		if state.currentLineWidth > calculatedSize.X {
			calculatedSize.X = state.currentLineWidth * fontScalingFactor
		}

		state.wrappedText.WriteString(state.currentLineText.String())
		state.currentLineText.Reset()

		state.currentLineWidth = 0
		state.currentLineText.WriteRune('\n')

		calculatedSize.Y += state.defaultLineHeight
	}

	previousCharacterState := IterationState{
		lastWhitespaceCharacterIndex: -1,
		currentLineWidth:             0,
		defaultLineHeight:            fontSize,
	}

	state := IterationState{
		lastWhitespaceCharacterIndex: -1,
		currentLineWidth:             0,
		defaultLineHeight:            fontSize,
	}

	runes := []rune(text)

	for characterIndex := 0; characterIndex < len(runes); characterIndex++ {
		character := runes[characterIndex]

		previousCharacterState = state

		if character == ' ' {
			state.lastWhitespaceCharacterIndex = characterIndex
		}

		state.currentLineText.WriteRune(character)

		if character == '\n' {
			processEndOfLine(&state)
		}

		glyphWidth := getGlyphWidth(&font, character)
		state.currentLineWidth += glyphWidth

		if state.currentLineWidth > maxWidth {
			if state.currentLineText.Len() == 1 {
				panic("Length of state.currentLineText is equal to 1, unhandled variant.")
			}

			state = previousCharacterState
			characterIndex -= 1

			processEndOfLine(&state)
		}
	}

	processEndOfInput(&state)

	processedText = state.wrappedText.String()

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

	var calculatedSize rl.Vector2

	if !comp.wrapText {
		calculatedSize = rl.MeasureTextEx(font, comp.text, comp.fontSize, comp.spacing)
		comp.processedText = comp.text

		return calculatedSize
	} else {
		comp.processedText, calculatedSize = wrapText(font, comp.text, comp.fontSize, comp.spacing, maxViewport.X)

		return calculatedSize
	}
}

func (comp *TextComponent) Render(getFont GetFontCallback) {
	font, err := getFont(comp.fontName)
	if err != nil {
		panic(fmt.Sprintf("Provided font (%s) is not loaded into memory", comp.fontName))
	}

	rl.DrawTextEx(font, comp.processedText, comp.position.Calculate(), comp.fontSize, comp.spacing, comp.color)
}

func (text *TextComponent) GetPosition() rl.Vector2 {
	return text.position.Calculate()
}
