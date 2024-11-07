package misc

import (
	"image/color"

	"github.com/voidwyrm-2/pixelite/palettes"
)

var DefaultConfig = map[string]any{
	"mouseMode": true,
	"moveUp":    "up",
	"moveDown":  "down",
	"moveLeft":  "left",
	"moveRight": "right",
}

func newColor(r, g, b uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

var DefaultPalette = palettes.New("pixite", []palettes.PaletteColor{
	// palettes.NewPaletteColor("Black", newColor(0, 0, 0)),
	// palettes.NewPaletteColor("White", newColor(255, 255, 255)),
	palettes.NewPaletteColor("Red", newColor(255, 0, 0)),
	palettes.NewPaletteColor("Green", newColor(0, 255, 0)),
	palettes.NewPaletteColor("Blue", newColor(0, 0, 255)),
	palettes.NewPaletteColor("Yellow", newColor(255, 255, 0)),
	palettes.NewPaletteColor("Cyan", newColor(0, 255, 255)),
	palettes.NewPaletteColor("Purple", newColor(255, 0, 255)),
})
