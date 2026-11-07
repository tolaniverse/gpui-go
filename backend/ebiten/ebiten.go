// Package ebiten adapts the gpui-go scene and input model to a native desktop
// window. Ebitengine supplies platform windows, GPU drawing, and font shaping.
package ebiten

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/tolaniverse/gpui-go"

	engine "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

type game struct {
	app    *gpui.App
	width  int
	height int
	font   *text.GoTextFaceSource
}

// Run opens a window and runs the UI event loop. Call it from main.
func Run(app *gpui.App, title string, width, height int) error {
	if app == nil || width <= 0 || height <= 0 {
		return fmt.Errorf("ebiten backend: invalid app or window dimensions")
	}
	font, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return fmt.Errorf("ebiten backend: load font: %w", err)
	}
	engine.SetWindowTitle(title)
	engine.SetWindowSize(width, height)
	engine.SetWindowResizingMode(engine.WindowResizingModeEnabled)
	return engine.RunGame(&game{app: app, width: width, height: height, font: font})
}

func (g *game) Update() error {
	g.app.Frame(float64(g.width), float64(g.height))
	if inpututil.IsKeyJustPressed(engine.KeyTab) {
		g.app.HandleKey(gpui.KeyEvent{Code: gpui.KeyTab, Shift: engine.IsKeyPressed(engine.KeyShift)})
	}
	if inpututil.IsKeyJustPressed(engine.KeyEnter) || inpututil.IsKeyJustPressed(engine.KeyNumpadEnter) {
		g.app.HandleKey(gpui.KeyEvent{Code: gpui.KeyEnter})
	}
	if inpututil.IsKeyJustPressed(engine.KeySpace) {
		g.app.HandleKey(gpui.KeyEvent{Code: gpui.KeySpace})
	}
	if inpututil.IsMouseButtonJustPressed(engine.MouseButtonLeft) {
		x, y := engine.CursorPosition()
		g.app.Click(float64(x), float64(y))
	}
	return nil
}

func (g *game) Draw(screen *engine.Image) {
	screen.Fill(color.RGBA{R: 18, G: 21, B: 28, A: 255})
	for _, command := range g.app.Frame(float64(g.width), float64(g.height)).Commands {
		switch command.Kind {
		case gpui.FillRect:
			vector.FillRect(screen, float32(command.Bounds.X), float32(command.Bounds.Y),
				float32(command.Bounds.W), float32(command.Bounds.H), command.Color, false)
		case gpui.DrawText:
			face := &text.GoTextFace{Source: g.font, Size: command.FontSize}
			options := &text.DrawOptions{}
			options.GeoM.Translate(command.Bounds.X, command.Bounds.Y)
			options.ColorScale.ScaleWithColor(command.Color)
			text.Draw(screen, command.Text, face, options)
		case gpui.FocusRect:
			bounds := command.Bounds
			if bounds.W > 4 && bounds.H > 4 {
				vector.StrokeRect(screen, float32(bounds.X+2), float32(bounds.Y+2),
					float32(bounds.W-4), float32(bounds.H-4), 2, command.Color, true)
			}
		}
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width, g.height = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}
