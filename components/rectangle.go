package components

import (
	"domanscy.group/gui/components/atoms"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RectangleComponent struct {
	eventBus *atoms.EventBus
	child    Component

	position ComponentPosition

	size rl.Vector2

	backgroundColor rl.Color

	padding   atoms.ClockValues
	roundness float32
}

func NewRectangleComponent(eventBus *atoms.EventBus, child Component, backgroundColor rl.Color, roundness float32) *RectangleComponent {
	if roundness < 0 {
		panic("Roundness can't be less than 0.")
	}

	return &RectangleComponent{
		eventBus: eventBus,
		child:    child,
		position: NewComponentPosition(),

		size:            rl.Vector2Zero(),
		backgroundColor: backgroundColor,
		padding:         atoms.NewClockValues(),
		roundness:       roundness,
	}
}

func (rec *RectangleComponent) SetChild(child Component) {
	rec.child = child
}

func (rec *RectangleComponent) getChildPositionOffset() rl.Vector2 {
	return rl.Vector2Add(rec.GetPosition(), rl.Vector2{
		X: rec.GetPaddingLeft(),
		Y: rec.GetPaddingTop(),
	})
}

func (rec *RectangleComponent) SetPosition(pos rl.Vector2) {
	rec.position.Position = pos

	rec.child.SetPositionOffset(rec.getChildPositionOffset())
}

func (rec *RectangleComponent) SetPositionOffset(offset rl.Vector2) {
	rec.position.Offset = offset

	rec.child.SetPositionOffset(rec.getChildPositionOffset())
}

func (rec *RectangleComponent) Render(getFont GetFontCallback) {
	position := rec.GetPosition()

	if rec.roundness != 0 {
		rectangleBoundaries := rl.Rectangle{
			X:      position.X,
			Y:      position.Y,
			Width:  rec.size.X,
			Height: rec.size.Y,
		}

		rl.DrawRectangleRounded(rectangleBoundaries, rec.roundness, 0, rec.backgroundColor)
	} else {
		rl.DrawRectangle(int32(position.X), int32(position.Y), int32(rec.size.X), int32(rec.size.Y), rec.backgroundColor)
	}

	rec.child.Render(getFont)
}

func (rec *RectangleComponent) CalculateSize(getFont GetFontCallback, maxViewport rl.Vector2) rl.Vector2 {
	childSize := rec.child.CalculateSize(
		getFont,
		rl.Vector2Add(
			maxViewport,
			rl.Vector2{X: -rec.padding.HorizontalSum(), Y: -rec.padding.VerticalSum()},
		),
	)

	rec.size = rl.Vector2Add(childSize, rl.Vector2{X: rec.padding.HorizontalSum(), Y: rec.padding.VerticalSum()})

	rec.child.SetPositionOffset(rec.getChildPositionOffset())

	return rec.size
}

func (rec *RectangleComponent) GetPosition() rl.Vector2 {
	return rec.position.Calculate()
}

func (rec *RectangleComponent) GetEventBus() *atoms.EventBus {
	return rec.eventBus
}

func (rec *RectangleComponent) GetPaddingTop() float32 {
	return rec.padding.Top()
}

func (rec *RectangleComponent) GetPaddingLeft() float32 {
	return rec.padding.Left()
}

func (rec *RectangleComponent) GetPaddingRight() float32 {
	return rec.padding.Right()
}

func (rec *RectangleComponent) GetPaddingBottom() float32 {
	return rec.padding.Bottom()
}

func (rec *RectangleComponent) SetPaddingTop(value float32) {
	rec.padding.SetTop(value)
	rec.eventBus.DispatchEvent("gui:schedule-recalculation", nil)
}

func (rec *RectangleComponent) SetPaddingLeft(value float32) {
	rec.padding.SetLeft(value)
	rec.eventBus.DispatchEvent("gui:schedule-recalculation", nil)
}

func (rec *RectangleComponent) SetPaddingRight(value float32) {
	rec.padding.SetRight(value)
	rec.eventBus.DispatchEvent("gui:schedule-recalculation", nil)
}

func (rec *RectangleComponent) SetPaddingBottom(value float32) {
	rec.padding.SetBottom(value)
	rec.eventBus.DispatchEvent("gui:schedule-recalculation", nil)
}
