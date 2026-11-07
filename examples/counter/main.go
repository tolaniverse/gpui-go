package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/tolaniverse/gpui-go"
	"github.com/tolaniverse/gpui-go/backend/ebiten"
)

type counter struct{ count *gpui.Entity[int] }

func (c *counter) Render() *gpui.Element {
	value := 0
	c.count.Read(func(n *int) { value = *n })
	increment := gpui.Text("Increment")
	increment.FocusID = "increment"
	increment.Style = gpui.Style{
		Height: 56, Padding: 16, FontSize: 19,
		Background: color.RGBA{R: 55, G: 100, B: 220, A: 255},
	}
	increment.OnClick = func() { c.count.Update(func(n *int) { *n++ }) }
	decrement := gpui.Text("Decrement")
	decrement.FocusID = "decrement"
	decrement.Style = gpui.Style{
		Height: 56, Padding: 16, FontSize: 19,
		Background: color.RGBA{R: 67, G: 75, B: 96, A: 255},
	}
	decrement.OnClick = func() { c.count.Update(func(n *int) { *n-- }) }
	label := gpui.Text(fmt.Sprintf("Count: %d", value))
	label.Style = gpui.Style{Height: 56, FontSize: 28}
	buttons := gpui.Row(increment, decrement)
	buttons.Style.Gap = 12
	root := gpui.Column(label, buttons)
	root.Style = gpui.Style{Padding: 32, Gap: 12}
	return root
}

func main() {
	view := &counter{}
	app := gpui.New(view)
	view.count = gpui.NewEntity(app, 0)
	if err := ebiten.Run(app, "gpui-go counter", 460, 230); err != nil {
		log.Fatal(err)
	}
}
