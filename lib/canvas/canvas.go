package lib

import (
	"errors"
	"image"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/voidwyrm-2/pixelite/lib"
)

type Canvas struct {
	pixels [][]color.RGBA
}

func New(sizeX, sizeY int, canvasColor color.RGBA) Canvas {
	c := [][]color.RGBA{}

	if sizeX <= 0 {
		panic("sizeX cannot be less than 1")
	} else if sizeY <= 0 {
		panic("sizeY cannot be less than 1")
	}

	for y := range sizeY {
		c = append(c, []color.RGBA{})
		for range sizeX {
			c[y] = append(c[y], canvasColor)
		}
	}

	return Canvas{pixels: c}
}

func FromImage(img image.Image) Canvas {
	c := New(img.Bounds().Dy(), img.Bounds().Dx(), rl.Black)

	for x := range img.Bounds().Dy() {
		for y := range img.Bounds().Dx() {
			r, g, b, a := img.At(x, y).RGBA()
			c.pixels[x][y] = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
		}
	}

	return c
}

func (c Canvas) At(x, y int) (color.RGBA, bool) {
	if y >= len(c.pixels) || y < 0 {
		return color.RGBA{}, false
	} else if x >= len(c.pixels[y]) || x < 0 {
		return color.RGBA{}, false
	}
	return c.pixels[y][x], true
}

func (c *Canvas) SetPixel(x, y int, color color.RGBA) {
	if y >= len(c.pixels) || y < 0 {
		return
	} else if x >= len(c.pixels[y]) || x < 0 {
		return
	}
	c.pixels[x][y] = color
}

func colorsEqual(c1, c2 color.RGBA) bool {
	return c1.R == c2.R && c1.G == c2.G && c1.B == c2.B && c1.A == c2.A
}

func (c *Canvas) Fill(x, y int, color color.RGBA) error {
	col, _ := c.At(x, y)
	if colorsEqual(color, col) {
		return nil
	}

	var subfill func(x, y int)
	subfill = func(x, y int) {
		if !colorsEqual(col, lib.Assert(c.At(x, y))) {
			return
		}
		c.SetPixel(x, y, color)

		if _, ok := c.At(x+1, y); ok {
			subfill(x+1, y)
		}
		if _, ok := c.At(x-1, y); ok {
			subfill(x-1, y)
		}
		if _, ok := c.At(x, y+1); ok {
			subfill(x, y+1)
		}
		if _, ok := c.At(x, y-1); ok {
			subfill(x, y-1)
		}
	}

	e := make(chan string, 1)
	func() {
		defer func() {
			s := recover()
			if s == nil {
				e <- ""
			} else {
				e <- s.(string)
			}
		}()

		subfill(x, y)
	}()

	errs := <-e
	if errs != "" {
		return errors.New(errs)
	}
	return nil
}

func (c Canvas) Draw(canvasPosition [2]int, cursorPosition [2]int, pixelSize int) {
	for y := range c.pixels {
		for x := range c.pixels[y] {
			rl.DrawRectangle(int32(canvasPosition[1]+(y*pixelSize)), int32(canvasPosition[0]+(x*pixelSize)), int32(pixelSize), int32(pixelSize), lib.Assert(c.At(x, y)))
			if x == cursorPosition[0] && y == cursorPosition[1] {
				col := lib.Assert(c.At(x, y))
				if (col.R+col.B+col.G)/3 > 127 {
					rl.DrawRectangleLinesEx(rl.NewRectangle(float32(canvasPosition[0]+(x*pixelSize)), float32(canvasPosition[1]+(y*pixelSize)), float32(pixelSize), float32(pixelSize)), 10, rl.Black)
				} else {
					rl.DrawRectangleLinesEx(rl.NewRectangle(float32(canvasPosition[0]+(x*pixelSize)), float32(canvasPosition[1]+(y*pixelSize)), float32(pixelSize), float32(pixelSize)), 10, rl.White)
				}
			}
		}
	}
}

func (c Canvas) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, len(c.pixels), len(c.pixels[0])))

	for x := range len(c.pixels) {
		for y := range len(c.pixels[x]) {
			img.SetRGBA(x, y, c.pixels[x][y])
		}
	}

	return img
}

func (c Canvas) Size() [2]int {
	return [2]int{len(c.pixels[0]), len(c.pixels)}
}
