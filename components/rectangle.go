package components

import (
	"domanscy.group/gui/components/atoms"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RectangleComponent struct {
	eventBus *atoms.EventStore
	child    Component

	position ComponentPosition

	size rl.Vector2

	backgroundColor rl.Color
}

func NewRectangleComponent(child Component, backgroundColor rl.Color) *RectangleComponent {
	return &RectangleComponent{
		eventBus: atoms.NewEventStore(),
		child:    child,
		position: NewComponentPosition(),

		size:            rl.Vector2Zero(),
		backgroundColor: backgroundColor,
	}
}

func (rec *RectangleComponent) SetChild(child Component) {
	rec.child = child
}

func (rec *RectangleComponent) SetPosition(pos rl.Vector2) {
	rec.position.Position = pos

	rec.child.SetPositionOffset(rec.child.GetPosition())
}

func (rec *RectangleComponent) SetPositionOffset(offset rl.Vector2) {
	rec.position.Offset = offset

	rec.child.SetPositionOffset(rec.child.GetPosition())
}

func (rec *RectangleComponent) Render(getFont GetFontCallback) {
	position := rec.GetPosition()

	rl.DrawRectangle(int32(position.X), int32(position.Y), int32(rec.size.X), int32(rec.size.Y), rec.backgroundColor)
}

func (rec *RectangleComponent) CalculateSize(getFont GetFontCallback, maxViewport rl.Vector2) rl.Vector2 {
	childSize := rec.child.CalculateSize(getFont, maxViewport)

	rec.size = childSize

	return rec.size
}

func (rec *RectangleComponent) GetPosition() rl.Vector2 {
	return rec.position.Calculate()
}

func (rec *RectangleComponent) GetEventBus() *atoms.EventStore {
	return rec.eventBus
}

func (rec *RectangleComponent) PropagateEvent(eventType string, args ...interface{}) {
	rec.eventBus.DispatchEvent(eventType, args...)

	rec.child.PropagateEvent(eventType, args)
}
