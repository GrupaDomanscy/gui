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

	text string
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

		text: "",
	}

	eventBus.ListenToEvent("gui:keyaction", func(rawArgs ...interface{}) {
		args := rawArgs[0].(events.KeyPressEvent)

		textInRunes := []rune(component.text)

		if slices.Contains(args.PressedKeys, rl.KeyBackspace) {
			if len(textInRunes) == 0 {
				return
			}

			textInRunes = textInRunes[:len(textInRunes)-1]
		} else {
			textInRunes = append(textInRunes, args.PressedChars...)
		}

		component.text = string(textInRunes)
	})

	return component
}

func (input *InputComponent) CalculateSize(getFont GetFontCallback, maxViewport rl.Vector2) rl.Vector2 {
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

	rl.DrawTextEx(font, input.text, input.position.Calculate(), input.fontSize, 0, input.textColor)
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
