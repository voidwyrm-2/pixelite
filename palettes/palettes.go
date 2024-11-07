package palettes

import "image/color"

type Palette struct {
	name   string
	colors []PaletteColor
}

func New(name string, colors []PaletteColor) Palette {
	return Palette{name: name, colors: colors}
}

func (p Palette) Length() int {
	return len(p.colors)
}

func (p Palette) GetColor(index int) color.RGBA {
	return p.colors[index].color
}

type PaletteColor struct {
	name  string
	color color.RGBA
}

func NewPaletteColor(name string, color color.RGBA) PaletteColor {
	return PaletteColor{name: name, color: color}
}
