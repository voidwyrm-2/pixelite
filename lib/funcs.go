package lib

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func IsFileNotFoundError(e error) bool {
	return strings.HasSuffix(e.Error(), ": no such file or directory")
}

func GetPixeliteVersion() (string, error) {
	res, err := http.Get("https://raw.githubusercontent.com/voidwyrm-2/pixelite/refs/heads/main/version.txt")
	if err != nil {
		return "", err
	}

	version, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	} else if string(version) == "404: Not Found" {
		return "", err
	}

	return string(version), nil
}

func DebugPrintln(isDebug *bool) func(a ...any) {
	return func(a ...any) {
		if *isDebug {
			fmt.Println(a...)
		}
	}
}

func Assert[A any, B any](a A, _ B) A {
	return a
}

func LoadImage(fpath string) (image.Image, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	image, _, err := image.Decode(f)
	return image, err
}

func SaveImage(fpath string, img image.Image) error {
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	switch path.Ext(fpath) {
	case ".jpg", "jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	case ".png":
		return png.Encode(f, img)
	case "":
		fpath += ".png"
		return png.Encode(f, img)
	}

	return fmt.Errorf("invalid image format '%s'", path.Ext(fpath))
}

func Clamp[N float32 | float64 | int | int8 | int16 | int32 | int64](n, min, max N) N {
	if n < min {
		return min
	} else if n > max {
		return max
	}
	return n
}

func DrawTextLines(text []string, startX, startY, fontSize int32, color color.RGBA) {
	for ln, line := range text {
		rl.DrawText(line, startX, startY+(fontSize*int32(ln)), fontSize, color)
	}
}
