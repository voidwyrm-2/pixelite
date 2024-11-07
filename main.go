package main

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/akamensky/argparse"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/voidwyrm-2/goconf"
	"github.com/voidwyrm-2/pixelite/lib"
	canvas "github.com/voidwyrm-2/pixelite/lib/canvas"
	"github.com/voidwyrm-2/pixelite/misc"
	"github.com/voidwyrm-2/pixelite/palettes"
)

func sExit(s string) {
	fmt.Println(s)
	os.Exit(1)
}

func eExit(e error) {
	sExit(e.Error())
}

var (
	config     map[string]any
	configPath = path.Join(rl.HomeDir(), ".pixelite")
)

//go:embed version.txt
var version string

type pixeliteState uint8

const (
	EXIT pixeliteState = iota
	CANVAS
	PALETTES
)

var (
	mouseMode      = true
	layers         []canvas.Canvas
	layerIndex     = 0
	hiddenLayers   = []int{}
	cursorPosition = [2]int{0, 0}
	canvasOffset   = [2]int{0, 0}
	pixelSize      = 10

	state                = CANVAS
	previousState        = CANVAS
	showingPaletteColors = false
)

func switchState(newState pixeliteState) {
	previousState = state
	state = newState
}

var shiftKey = false

var (
	loadedPalettes = []palettes.Palette{misc.DefaultPalette}
	paletteIndex   = 0
	colorIndex     = 0
)

func applyConfig() error {
	if mmode, ok := config["mouseMode"]; !ok {
		mouseMode = true
	} else if mm, ok := mmode.(bool); !ok {
		return fmt.Errorf("expected config option 'mouseMode' to be 'bool' not '%s'", reflect.TypeOf(mmode).Name())
	} else {
		mouseMode = mm
	}

	return nil
}

func main() {
	parser := argparse.NewParser("pixelite", "A pixel art editor")

	showDebug := parser.Flag("d", "debug", &argparse.Options{Required: false, Help: "Shows extra debug messages"})
	showVersion := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "Shows the current Pixelite version"})
	filename := parser.String("f", "file", &argparse.Options{Required: false, Help: "The file to edit or create"})
	savedAs := parser.String("s", "saveas", &argparse.Options{Required: false, Help: "The path to save the edited file at instead of saving to the file opened"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if *showVersion {
		fmt.Println(version)
	}

	debug := lib.DebugPrintln(showDebug)

	debug("beginning pre-window init")

	debug("checking verison...")
	if vr, err := lib.GetPixeliteVersion(); err != nil {
		eExit(err)
	} else {
		if verRemote, err := lib.NewVersionFromVersionString(vr); err != nil {
			eExit(err)
		} else if verLocal, err := lib.NewVersionFromVersionString(version); err != nil {
			eExit(err)
		} else if verLocal.Compare(verRemote) == -1 {
			fmt.Printf("A new version of Pixelite is available!(%s -> %s)\nrun `go get github.com/voidwyrm-2/pixelite` to install it\n", verLocal.Fmt(), verRemote.Fmt())
			return
		}
	}
	debug("version checked")

	if *showVersion {
		return
	}

	debug("loading config...")
	if cnf, err := goconf.Load(configPath); err != nil {
		if lib.IsFileNotFoundError(err) {
			err = goconf.Save(configPath, misc.DefaultConfig)
			if err != nil {
				eExit(err)
			}
			config = misc.DefaultConfig
		} else {
			eExit(err)
		}
	} else {
		config = cnf
	}
	debug("loaded config")

	if err := applyConfig(); err != nil {
		eExit(err)
	}

	rl.InitWindow(800, 800, "Pixelite")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// layers = append(layers, canvas.New(20, 20, rl.White))

	if strings.TrimSpace(*filename) == "" {
		sExit("filename cannot be empty")
	}

	if strings.TrimSpace(*savedAs) == "" {
		*savedAs = *filename
	}

	img, err := lib.LoadImage(*filename)
	if err != nil {
		eExit(err)
	}
	layers = append(layers, canvas.FromImage(img))

	for {
		if state == EXIT {
			if rl.IsKeyPressed(rl.KeyY) {
				break
			} else if rl.IsKeyPressed(rl.KeyN) {
				switchState(previousState)
			}
		} else if state == CANVAS {
			if rl.WindowShouldClose() {
				switchState(EXIT)
				continue
			}

			if rl.IsKeyDown(rl.KeyQ) {
				pixelSize += 1
			}
			if rl.IsKeyDown(rl.KeyE) && pixelSize > 1 {
				pixelSize -= 1
			}

			/*
				if rl.IsKeyPressed(int32(rl.MouseButtonLeft)) {
					mDelta := rl.GetMouseDelta()
					canvasOffset[0] += int(mDelta.X)
					canvasOffset[1] += int(mDelta.Y)
				}
			*/

			if rl.IsKeyPressed(rl.KeyP) {
				showingPaletteColors = !showingPaletteColors
			}

			if rl.IsKeyDown(rl.KeyUp) {
				if showingPaletteColors {
					colorIndex = lib.Clamp(colorIndex-10, 0, loadedPalettes[paletteIndex].Length()-1)
				} else {
					cursorPosition[1] = lib.Clamp(cursorPosition[1]-1, 0, layers[0].Size()[1]-1)
				}
			}
			if rl.IsKeyDown(rl.KeyDown) {
				if showingPaletteColors {
					colorIndex = lib.Clamp(colorIndex+10, 0, loadedPalettes[paletteIndex].Length()-1)
				} else {
					cursorPosition[1] = lib.Clamp(cursorPosition[1]+1, 0, layers[0].Size()[1]-1)
				}
			}
			if rl.IsKeyDown(rl.KeyRight) {
				if showingPaletteColors {
					colorIndex = lib.Clamp(colorIndex-1, 0, loadedPalettes[paletteIndex].Length()-1)
				} else {
					cursorPosition[0] = lib.Clamp(cursorPosition[0]-1, 0, layers[0].Size()[0]-1)
				}
			}
			if rl.IsKeyDown(rl.KeyRight) {
				if showingPaletteColors {
					colorIndex = lib.Clamp(colorIndex+1, 0, loadedPalettes[paletteIndex].Length()-1)
				} else {
					cursorPosition[0] = lib.Clamp(cursorPosition[0]+1, 0, layers[0].Size()[0]-1)
				}
			}

			if rl.IsKeyDown(rl.KeyW) {
				canvasOffset[1] -= 5
			}
			if rl.IsKeyDown(rl.KeyS) && !shiftKey {
				canvasOffset[1] += 5
			}
			if rl.IsKeyDown(rl.KeyA) {
				canvasOffset[0] -= 5
			}
			if rl.IsKeyDown(rl.KeyD) {
				canvasOffset[0] += 5
			}

			if rl.IsKeyDown(rl.KeySlash) {
				layers[layerIndex].SetPixel(cursorPosition[0], cursorPosition[1], loadedPalettes[paletteIndex].GetColor(colorIndex))
			}

			if mouseMode {
				mouseWM := rl.GetMouseWheelMove()
				pixelSize += int(mouseWM)

				mousePos := rl.GetMousePosition()
				cursorPosition[0] = lib.Clamp((int(mousePos.X)-canvasOffset[0])/pixelSize, 0, layers[0].Size()[0]-1)
				cursorPosition[1] = lib.Clamp((int(mousePos.Y)-canvasOffset[1])/pixelSize, 0, layers[0].Size()[1]-1)

				if rl.IsKeyDown(int32(rl.MouseLeftButton)) {
					layers[layerIndex].SetPixel(cursorPosition[0], cursorPosition[1], loadedPalettes[paletteIndex].GetColor(colorIndex))
				}
			}

			shiftKey = rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

			if shiftKey {
				if rl.IsKeyPressed(rl.KeyS) {
					if err := lib.SaveImage(*savedAs, layers[layerIndex].ToImage()); err != nil {
						fmt.Println("error while saving to file '"+*savedAs+"'", err.Error())
					}
				}
			}
		} else {
			panic(fmt.Sprintf("invalid Pixelite state %v", state))
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		if state == EXIT {
			lib.DrawTextLines([]string{"Are you sure you want to exit?", "Any unsaved changes will be lost", "(Y/N)"}, 10, 10, 20, rl.RayWhite)
		} else if state == CANVAS {
			for _, layer := range layers {
				layer.Draw(canvasOffset, cursorPosition, pixelSize)
			}
		}

		rl.EndDrawing()
	}
}
