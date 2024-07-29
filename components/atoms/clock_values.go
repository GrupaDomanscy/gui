package atoms

type ClockValues struct {
	top    float32
	bottom float32
	left   float32
	right  float32
}

func NewClockValues() ClockValues {
	return ClockValues{0, 0, 0, 0}
}

func (c ClockValues) Top() float32 {
	return c.top
}

func (c ClockValues) Bottom() float32 {
	return c.bottom
}

func (c ClockValues) Left() float32 {
	return c.left
}

func (c ClockValues) Right() float32 {
	return c.right
}

func (c ClockValues) SetTop(value float32) {
	c.top = value
}

func (c ClockValues) SetBottom(value float32) {
	c.bottom = value
}

func (c ClockValues) SetLeft(value float32) {
	c.left = value
}

func (c ClockValues) SetRight(value float32) {
	c.right = value
}

func (c ClockValues) HorizontalSum() float32 {
	return c.left + c.right
}

func (c ClockValues) VerticalSum() float32 {
	return c.top + c.bottom
}
