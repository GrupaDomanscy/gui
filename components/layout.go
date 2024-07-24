package components

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const DirectionColumn = 0
const DirectionRow = 1

const AlignStart = 0
const AlignCenter = 1
const AlignEnd = 2

type LayoutComponent struct {
	children           []Component
	direction          int
	mainAxisAlignment  int
	crossAxisAlignment int

	position ComponentPosition
}

func NewLayoutComponent(direction int, mainAxisAlignment int, crossAxisAlignment int) *LayoutComponent {
	if direction != DirectionColumn && direction != DirectionRow {
		panic(fmt.Sprintf("Unknown value for direction property: %d", direction))
	}

	if mainAxisAlignment != AlignCenter && mainAxisAlignment != AlignStart && mainAxisAlignment != AlignEnd {
		panic(fmt.Sprintf("Unknwon value for mainAxisAlignment property: %d", mainAxisAlignment))
	}

	if crossAxisAlignment != AlignCenter && crossAxisAlignment != AlignStart && crossAxisAlignment != AlignEnd {
		panic(fmt.Sprintf("Unknwon value for crossAxisAlignment property: %d", crossAxisAlignment))
	}

	return &LayoutComponent{
		children:           make([]Component, 0),
		direction:          direction,
		mainAxisAlignment:  mainAxisAlignment,
		crossAxisAlignment: crossAxisAlignment,
		position:           NewComponentPosition(),
	}
}

func (layout *LayoutComponent) calculateSizesOfChildren(getFont GetFontCallback, maxViewport rl.Vector2) []rl.Vector2 {
	sizes := make([]rl.Vector2, len(layout.children))

	currentMaxViewport := maxViewport

	if layout.direction == DirectionColumn {
		for i, child := range layout.children {
			childSize := child.CalculateSize(getFont, currentMaxViewport)

			sizes[i] = childSize

			currentMaxViewport.Y -= childSize.Y
		}
	} else {
		for i, child := range layout.children {
			childSize := child.CalculateSize(getFont, currentMaxViewport)

			sizes[i] = childSize

			currentMaxViewport.X -= childSize.X
		}
	}

	return sizes
}

func getArrayOfXAxisFromVector2(arr []rl.Vector2) []float32 {
	xAxis := make([]float32, len(arr))

	for i, item := range arr {
		xAxis[i] = item.X
	}

	return xAxis
}

func getArrayOfYAxisFromVector2(arr []rl.Vector2) []float32 {
	yAxis := make([]float32, len(arr))

	for i, item := range arr {
		yAxis[i] = item.Y
	}

	return yAxis
}

func sum[T float32](values []T) T {
	var result T

	for i := 0; i < len(values); i++ {
		result += values[i]
	}

	return result
}

func (layout *LayoutComponent) calculateChildPositionsAndParentSizeForMainAxis(sizes []float32, maxViewport float32) (positions []float32, parentSize float32) {
	positions = make([]float32, len(sizes))

	var currentPos float32

	switch layout.mainAxisAlignment {
	case AlignStart:
		currentPos = 0

		for i, size := range sizes {
			positions[i] = currentPos
			currentPos += size
			parentSize += size
		}

		break
	case AlignCenter:
		sizeSum := sum(sizes)
		currentPos = (maxViewport / 2) - (sizeSum / 2)

		for i, size := range sizes {
			positions[i] = currentPos
			currentPos += size
		}

		if sizeSum > maxViewport {
			parentSize = sizeSum
		} else {
			parentSize = maxViewport
		}
		break
	case AlignEnd:
		sizeSum := sum(sizes)
		currentPos = maxViewport - sizeSum

		for i, size := range sizes {
			positions[i] = currentPos
			currentPos += size
		}

		if sizeSum > maxViewport {
			parentSize = sizeSum
		} else {
			parentSize = maxViewport
		}
		break
	default:
		panic(fmt.Sprintf("Unhandled main axis alignment value: %d", layout.mainAxisAlignment))
	}

	return
}

func (layout *LayoutComponent) calculateChildPositionsAndParentSizeForCrossAxis(sizes []float32, maxViewport float32) (positions []float32, parentSize float32) {
	positions = make([]float32, len(sizes))

	switch layout.crossAxisAlignment {
	case AlignStart:
		var maxChildSize float32 = 0

		for i := 0; i < len(sizes); i++ {
			positions[i] = 0

			if sizes[i] > maxChildSize {
				maxChildSize = sizes[i]
			}
		}

		parentSize = maxChildSize
		break
	case AlignCenter:
		var maxChildSize float32 = 0

		center := maxViewport / 2

		for i := 0; i < len(sizes); i++ {
			positions[i] = center - (sizes[i] / 2)

			if sizes[i] > maxChildSize {
				maxChildSize = sizes[i]
			}
		}

		if maxChildSize > maxViewport {
			parentSize = maxChildSize
		} else {
			parentSize = maxViewport
		}
		break
	case AlignEnd:
		var maxChildSize float32 = 0

		for i := 0; i < len(sizes); i++ {
			positions[i] = maxViewport - sizes[i]

			if sizes[i] > maxChildSize {
				maxChildSize = sizes[i]
			}
		}

		if maxChildSize > maxViewport {
			parentSize = maxChildSize
		} else {
			parentSize = maxViewport
		}
		break
	default:
		panic(fmt.Sprintf("Unhandled cross axis alignment value: %d", layout.crossAxisAlignment))
	}

	return
}

func joinFloatArraysToVector2Array(xArr []float32, yArr []float32) []rl.Vector2 {
	if len(xArr) != len(yArr) {
		panic(fmt.Sprintf("first and second array are not equal in length! %d != %d", len(xArr), len(yArr)))
	}

	joinedArr := make([]rl.Vector2, len(xArr))

	for i := 0; i < len(xArr); i++ {
		joinedArr[i].X = xArr[i]
		joinedArr[i].Y = yArr[i]
	}

	return joinedArr
}

func (layout *LayoutComponent) CalculateSize(getFont GetFontCallback, maxViewport rl.Vector2) rl.Vector2 {
	childrenSizes := layout.calculateSizesOfChildren(getFont, maxViewport)

	xAxisSizes := getArrayOfXAxisFromVector2(childrenSizes)
	yAxisSizes := getArrayOfYAxisFromVector2(childrenSizes)

	var xAxisPositions []float32
	var yAxisPositions []float32

	var xAxisParentSize float32
	var yAxisParentSize float32

	switch layout.direction {
	case DirectionColumn:
		yAxisPositions, yAxisParentSize = layout.calculateChildPositionsAndParentSizeForMainAxis(yAxisSizes, maxViewport.Y)
		xAxisPositions, xAxisParentSize = layout.calculateChildPositionsAndParentSizeForCrossAxis(xAxisSizes, maxViewport.X)
		break
	case DirectionRow:
		xAxisPositions, xAxisParentSize = layout.calculateChildPositionsAndParentSizeForMainAxis(xAxisSizes, maxViewport.X)
		yAxisPositions, yAxisParentSize = layout.calculateChildPositionsAndParentSizeForCrossAxis(yAxisSizes, maxViewport.Y)
		break
	default:
		panic(fmt.Sprintf("Unhandled direction parameter value in layout component: %d", layout.direction))
	}

	positions := joinFloatArraysToVector2Array(xAxisPositions, yAxisPositions)

	for i, child := range layout.children {
		child.SetPosition(positions[i])
	}

	return rl.Vector2{
		X: xAxisParentSize,
		Y: yAxisParentSize,
	}
}

func (layout *LayoutComponent) Render(getFont GetFontCallback) {
	for _, child := range layout.children {
		child.Render(getFont)
	}
}

func (layout *LayoutComponent) SetPosition(pos rl.Vector2) {
	layout.position.Position = pos

	newLayoutPosition := layout.position.Calculate()

	for _, child := range layout.children {
		child.SetPositionOffset(newLayoutPosition)
	}
}

func (layout *LayoutComponent) SetPositionOffset(offset rl.Vector2) {
	layout.position.Offset = offset

	newLayoutPosition := layout.position.Calculate()

	for _, child := range layout.children {
		child.SetPositionOffset(newLayoutPosition)
	}
}

func (layout *LayoutComponent) AddChild(child Component) {
	layout.children = append(layout.children, child)
}

func (layout *LayoutComponent) GetPosition() rl.Vector2 {
	return layout.position.Calculate()
}
