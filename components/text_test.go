package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"testing"
)

type TestFont struct {
}

func (TestFont) GlyphWidth(codepoint rune) float32 {
	return 32
}

func (TestFont) LineHeight() float32 {
	return 32
}

func (TestFont) FontSize() float32 {
	return 32
}

func (TestFont) Spacing() float32 {
	return 1
}

func TestWrapText(t *testing.T) {
	t.Run("Scenario 1", func(t *testing.T) {
		// "Hello world! How are you today? Hello world! How are you today? Hello world! How are you today?"

		testFont := TestFont{}

		expectedText := "Hello\nworld\n,\nhow\nare\nyou\ntoday\n?"
		expectedSize := rl.Vector2{
			X: 32 * 5,
			Y: 32 * 8,
		}

		processedText, size := wrapText(testFont, "Hello world, how are you today?", 32*5) // max 5 characters in one line

		if processedText != expectedText {
			t.Errorf("Processed text expected: \"%s\", received: \"%s\"", expectedText, processedText)
		}

		if !rl.Vector2Equals(size, expectedSize) {
			t.Errorf("Expected size X: %f Y: %f, received: X: %f Y: %f", expectedSize.X, expectedSize.Y, size.X, size.Y)
		}
	})
}
