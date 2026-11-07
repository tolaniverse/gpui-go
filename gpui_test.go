package gpui

import (
	"image/color"
	"testing"
)

type testView struct{ count *Entity[int] }

func (v *testView) Render() *Element {
	count := 0
	v.count.Read(func(value *int) { count = *value })
	child := Text(string(rune('0' + count)))
	child.Style = Style{Height: 40, Background: color.RGBA{R: 42, A: 255}}
	child.OnClick = func() { v.count.Update(func(value *int) { *value++ }) }
	root := Column(child, Text("footer"))
	root.Style = Style{Padding: 10, Gap: 5}
	return root
}

func TestEntityClickAndLayout(t *testing.T) {
	view := &testView{}
	app := New(view)
	view.count = NewEntity(app, 0)
	scene := app.Frame(200, 100)
	if len(scene.Commands) != 3 || scene.Commands[0].Bounds != (Rect{X: 10, Y: 10, W: 180, H: 40}) {
		t.Fatalf("unexpected scene: %+v", scene.Commands)
	}
	if app.Click(5, 5) || !app.Click(20, 20) {
		t.Fatal("click dispatch did not match element bounds")
	}
	if got := app.Frame(200, 100).Commands[1].Text; got != "1" {
		t.Fatalf("after click, got label %q, want 1", got)
	}
	if app.Click(20, 70) {
		t.Fatal("footer should not be clickable")
	}
}

func TestRowFixedAndAutoSpace(t *testing.T) {
	first, second, third := Text("a"), Text("b"), Text("c")
	first.Style.Width = 50
	root := Row(first, second, third)
	root.Style.Gap = 5
	scene := buildScene(root, 200, 30)
	if got := scene.Commands[2].Bounds; got != (Rect{X: 130, Y: 0, W: 70, H: 30}) {
		t.Fatalf("last child bounds = %+v", got)
	}
}
