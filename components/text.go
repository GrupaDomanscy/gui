package components

import (
	"fmt"
	"strings"

	"domanscy.group/gui/components/atoms"
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

	eventBus *atoms.EventBus
}

func NewTextComponent(eventBus *atoms.EventBus, text string, loadedFontName string, fontSize float32, spacing float32, color rl.Color) *TextComponent {
	return &TextComponent{
		text:          text,
		processedText: "",
		wrapText:      false,
		fontName:      loadedFontName,
		fontSize:      fontSize,
		spacing:       spacing,
		color:         color,
		position:      NewComponentPosition(),
		eventBus:      eventBus,
	}
}

func (comp *TextComponent) SetWrapText(wrap bool) {
	comp.wrapText = wrap
}

func wrapText(font atoms.Font, text string, maxWidth float32) (processedText string, calculatedSize rl.Vector2) {
	type IterationState struct {
		wrappedText      strings.Builder
		currentLineText  strings.Builder
		currentLineWidth float32
	}

	processEndOfInput := func(state *IterationState) {
		state.wrappedText.WriteString(state.currentLineText.String())
		calculatedSize.Y += font.LineHeight()

		if state.currentLineWidth > calculatedSize.X {
			calculatedSize.X = state.currentLineWidth
		}
	}

	processEndOfLine := func(state *IterationState) {
		if state.currentLineWidth > calculatedSize.X {
			calculatedSize.X = state.currentLineWidth
		}

		state.wrappedText.WriteString(strings.Trim(state.currentLineText.String(), " "))
		state.currentLineText.Reset()

		state.currentLineWidth = 0
		state.wrappedText.WriteRune('\n')

		calculatedSize.Y += font.LineHeight()

		// This is the line spacing, for more context see original raylib function:
		// https://github.com/raysan5/raylib/blob/9e39788e077f1d35c5fe54600f2143423a80bb3d/src/rtext.c#L1164
		// Be aware that text line spacing is a global variable in raylib
		calculatedSize.Y += 2
	}

	lastWhitespaceCharacterIndex := -1
	lastWhitespaceCharacterState := IterationState{}
	previousCharacterState := IterationState{}
	state := IterationState{}

	runes := []rune(text)

	for characterIndex := 0; characterIndex < len(runes); characterIndex++ {
		character := runes[characterIndex]

		previousCharacterState = state

		if character == '\n' {
			processEndOfLine(&state)
			continue
		}

		glyphWidth := font.GlyphWidth(character)

		if state.currentLineWidth+glyphWidth > maxWidth {
			if state.currentLineText.Len() == 1 {
				panic("Length of state.currentLineText is equal to 1, unhandled variant.")
			}

			if lastWhitespaceCharacterIndex != -1 {
				state = lastWhitespaceCharacterState
				characterIndex = lastWhitespaceCharacterIndex
			} else {
				state = previousCharacterState
				characterIndex -= 1
			}

			lastWhitespaceCharacterIndex = -1
			processEndOfLine(&state)

			continue
		}

		if character == ' ' && state.currentLineText.Len() == 0 {
			continue
		}

		state.currentLineText.WriteRune(character)
		state.currentLineWidth += glyphWidth

		if character == ' ' {
			lastWhitespaceCharacterIndex = characterIndex
			lastWhitespaceCharacterState = state
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
		comp.processedText, calculatedSize = wrapText(atoms.NewRaylibFont(font, comp.fontSize, comp.spacing), comp.text, maxViewport.X)

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

func (comp *TextComponent) GetPosition() rl.Vector2 {
	return comp.position.Calculate()
}

func (comp *TextComponent) GetEventBus() *atoms.EventBus {
	return comp.eventBus
}
