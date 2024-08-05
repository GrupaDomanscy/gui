package components

import (
	"fmt"
	"slices"
	"time"

	"domanscy.group/gui/components/atoms"
	"domanscy.group/gui/components/events"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type InputComponent struct {
	position ComponentPosition
	size     rl.Vector2
	eventBus *atoms.EventBus

	fontName              string
	cachedGetFontCallback GetFontCallback
	fontSize              float32

	maxWidth float32

	backgroundColor rl.Color
	textColor       rl.Color

	text              string
	textGlyphHitboxes []rl.Rectangle

	cursorRunePos             int
	cursorXPos                float32
	showCursor                bool
	lastShowCursorStateChange int64

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

		fontName:              fontName,
		cachedGetFontCallback: nil,
		fontSize:              fontSize,

		maxWidth:        maxWidth,
		backgroundColor: backgroundColor,
		textColor:       textColor,

		text:              "",
		textGlyphHitboxes: []rl.Rectangle{},

		cursorRunePos: 0,
		cursorXPos:    0,

		isInitialized: false,

		showCursor:                false,
		lastShowCursorStateChange: time.Now().UnixMilli(),
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

func (component *InputComponent) processLeftClick() {
	component.cursorRunePos--

	if component.cursorRunePos < 0 {
		component.cursorRunePos = 0
	} else {
		component.cursorXPos -= component.textGlyphHitboxes[component.cursorRunePos].Width
	}
}

func (component *InputComponent) processRightClick() {
	textInRunes := []rune(component.text)

	component.cursorRunePos++

	if component.cursorRunePos > len(textInRunes) {
		component.cursorRunePos = len(textInRunes)
	} else {
		component.cursorXPos += component.textGlyphHitboxes[component.cursorRunePos-1].Width
	}

	component.text = string(textInRunes)
}

func (component *InputComponent) processBackspaceClick() {
	textInRunes := []rune(component.text)

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

	component.text = string(textInRunes)
}

func (component *InputComponent) processDeleteClick() {
	textInRunes := []rune(component.text)

	if len(textInRunes) == 0 {
		return
	}

	if component.cursorRunePos >= len(textInRunes) {
		return
	}

	component.textGlyphHitboxes = slices.Delete(component.textGlyphHitboxes, component.cursorRunePos, component.cursorRunePos+1)
	textInRunes = slices.Delete(textInRunes, component.cursorRunePos, component.cursorRunePos+1)

	component.text = string(textInRunes)
}

func (component *InputComponent) processCustomCharPress(pressedChars ...int32) {
	textInRunes := []rune(component.text)

	font, err := component.cachedGetFontCallback(component.fontName)
	if err != nil {
		panic("Font has not been loaded into memory")
	}

	for _, pressedChar := range pressedChars {
		hitbox := convertGlyphInfoToRectangle(rl.GetGlyphInfo(font, pressedChar), font, component.fontSize)
		component.textGlyphHitboxes = append(component.textGlyphHitboxes, hitbox)
		component.cursorXPos += hitbox.Width
		component.cursorRunePos++
	}

	textInRunes = append(textInRunes, pressedChars...)

	component.text = string(textInRunes)
}

func (component *InputComponent) assignEventListeners() {
	component.eventBus.ListenToEvent("gui:keyaction", func(rawArgs ...interface{}) {
		args := rawArgs[0].(events.KeyPressEvent)

		component.showCursor = true
		component.lastShowCursorStateChange = time.Now().UnixMilli()

		if slices.Contains(args.PressedKeys, rl.KeyLeft) {
			component.processLeftClick()
		} else if slices.Contains(args.PressedKeys, rl.KeyRight) {
			component.processRightClick()
		} else if slices.Contains(args.PressedKeys, rl.KeyBackspace) {
			component.processBackspaceClick()
		} else if slices.Contains(args.PressedKeys, rl.KeyDelete) {
			component.processDeleteClick()
		} else {
			component.processCustomCharPress(args.PressedChars...)
		}
	})
}

func (input *InputComponent) CalculateSize(getFont GetFontCallback, maxViewport rl.Vector2) rl.Vector2 {
	if !input.isInitialized {
		input.cachedGetFontCallback = getFont
		input.assignEventListeners()
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

	if time.Now().UnixMilli()-input.lastShowCursorStateChange > 500 {
		input.lastShowCursorStateChange = time.Now().UnixMilli()
		input.showCursor = !input.showCursor
	}

	if input.showCursor {
		rl.DrawRectangle(int32(input.cursorXPos), int32(position.Y), 2, int32(input.fontSize), input.textColor)
	}
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
