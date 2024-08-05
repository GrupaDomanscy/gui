package components

import (
	"fmt"
	"slices"

	"domanscy.group/gui/components/atoms"
	"domanscy.group/gui/components/events"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type InputComponent struct {
	position ComponentPosition
	size     rl.Vector2
	eventBus *atoms.EventBus

	fontName string
	fontSize float32

	maxWidth float32

	backgroundColor rl.Color
	textColor       rl.Color

	text              string
	textGlyphHitboxes []rl.Rectangle

	cursorRunePos int
	cursorXPos    float32

	isInitialized bool
}

func NewInputComponent(eventBus *atoms.EventBus, fontName string, fontSize float32, maxWidth float32, backgroundColor rl.Color, textColor rl.Color) *InputComponent {
	if eventBus == nil {
		panic("event bus can't be nil")
	}

	component := &InputComponent{
		position: NewComponentPosition(),
		size:     rl.Vector2Zero(),
		eventBus: eventBus,

		fontName: fontName,
		fontSize: fontSize,

		maxWidth:        maxWidth,
		backgroundColor: backgroundColor,
		textColor:       textColor,

		text:              "",
		textGlyphHitboxes: []rl.Rectangle{},

		cursorRunePos: 0,
		cursorXPos:    0,

		isInitialized: false,
	}

	return component
}

func convertGlyphInfoToRectangle(info rl.GlyphInfo, font rl.Font, fontSize float32) rl.Rectangle {
	return rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(info.AdvanceX) * float32(font.BaseSize) / fontSize,
		Height: fontSize,
	}
}

func (component *InputComponent) assignEventListeners(getFont GetFontCallback) {
	component.eventBus.ListenToEvent("gui:keyaction", func(rawArgs ...interface{}) {
		args := rawArgs[0].(events.KeyPressEvent)

		textInRunes := []rune(component.text)

		if slices.Contains(args.PressedKeys, rl.KeyLeft) {
			component.cursorRunePos--

			if component.cursorRunePos < 0 {
				component.cursorRunePos = 0
			} else {
				component.cursorXPos -= component.textGlyphHitboxes[component.cursorRunePos].Width
			}
		} else if slices.Contains(args.PressedKeys, rl.KeyRight) {
			component.cursorRunePos++

			if component.cursorRunePos > len(textInRunes) {
				component.cursorRunePos = len(textInRunes)
			} else {
				component.cursorXPos += component.textGlyphHitboxes[component.cursorRunePos-1].Width
			}
		} else if slices.Contains(args.PressedKeys, rl.KeyBackspace) {
			if len(textInRunes) == 0 {
				return
			}

			if component.cursorRunePos == 0 {
				return
			}

			component.cursorXPos -= component.textGlyphHitboxes[component.cursorRunePos-1].Width
			component.textGlyphHitboxes = slices.Delete(component.textGlyphHitboxes, component.cursorRunePos-1, component.cursorRunePos)
			textInRunes = slices.Delete(textInRunes, component.cursorRunePos-1, component.cursorRunePos)
			component.cursorRunePos--
		} else {
			font, err := getFont(component.fontName)
			if err != nil {
				panic("Font has not been loaded into memory")
			}

			for _, pressedChar := range args.PressedChars {
				hitbox := convertGlyphInfoToRectangle(rl.GetGlyphInfo(font, pressedChar), font, component.fontSize)
				component.textGlyphHitboxes = append(component.textGlyphHitboxes, hitbox)
				component.cursorXPos += hitbox.Width
				component.cursorRunePos++
			}

			textInRunes = append(textInRunes, args.PressedChars...)
		}

		component.text = string(textInRunes)
	})
}

func (input *InputComponent) CalculateSize(getFont GetFontCallback, maxViewport rl.Vector2) rl.Vector2 {
	if !input.isInitialized {
		input.assignEventListeners(getFont)
		input.isInitialized = true
	}

	var width float32 = 0

	if maxViewport.X < input.maxWidth {
		width = maxViewport.X
	} else {
		width = input.maxWidth
	}

	height := input.fontSize

	input.size = rl.Vector2{
		X: width,
		Y: height,
	}

	return input.size
}

func (input *InputComponent) Render(getFont GetFontCallback) {
	font, err := getFont(input.fontName)
	if err != nil {
		panic(fmt.Sprintf("font %s has not been loaded into memory", input.fontName))
	}

	position := input.position.Calculate()
	size := input.size

	rl.DrawRectangle(int32(position.X), int32(position.Y), int32(size.X), int32(size.Y), input.backgroundColor)

	rl.DrawTextEx(font, input.text, position, input.fontSize, 0, input.textColor)

	rl.DrawRectangle(int32(input.cursorXPos), int32(position.Y), 2, int32(input.fontSize), input.textColor)
}

func (input *InputComponent) SetPosition(pos rl.Vector2) {
	input.position.Position = pos
}

func (input *InputComponent) SetPositionOffset(offset rl.Vector2) {
	input.position.Offset = offset
}

func (input *InputComponent) GetPosition() rl.Vector2 {
	return input.position.Calculate()
}

func (input *InputComponent) GetEventBus() *atoms.EventBus {
	return input.eventBus
}
