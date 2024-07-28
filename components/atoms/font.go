package atoms

import rl "github.com/gen2brain/raylib-go/raylib"

type Font interface {
	GlyphWidth(codepoint rune) float32
	LineHeight() float32
	FontSize() float32
	Spacing() float32
}

type RaylibFont struct {
	font     rl.Font
	fontSize float32
	spacing  float32
}

func NewRaylibFont(font rl.Font, fontSize float32, spacing float32) *RaylibFont {
	return &RaylibFont{
		font,
		fontSize,
		spacing,
	}
}

func (font RaylibFont) GlyphWidth(codepoint rune) float32 {
	glyphInfo := rl.GetGlyphInfo(font.font, codepoint)

	fontScalingFactor := float32(font.font.BaseSize) / font.fontSize

	if glyphInfo.AdvanceX != 0 {
		return float32(glyphInfo.AdvanceX) * fontScalingFactor
	} else {
		return (font.font.Recs.Width + float32(glyphInfo.OffsetX)) * fontScalingFactor
	}
}

func (font RaylibFont) LineHeight() float32 {
	return font.fontSize
}

func (font RaylibFont) Spacing() float32 {
	return font.spacing
}

func (font RaylibFont) FontSize() float32 {
	return font.fontSize
}
