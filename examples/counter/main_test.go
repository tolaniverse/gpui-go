package main

import (
	"testing"

	"github.com/tolaniverse/gpui-go"
)

func TestCounterPointerAndKeyboard(t *testing.T) {
	view := &counter{}
	app := gpui.New(view)
	view.count = gpui.NewEntity(app, 0)
	app.Frame(460, 230)
	if !app.HandleKey(gpui.KeyEvent{Code: gpui.KeyTab}) || app.FocusID() != "increment" {
		t.Fatal("Tab did not focus increment")
	}
	if !app.HandleKey(gpui.KeyEvent{Code: gpui.KeySpace}) {
		t.Fatal("Space did not increment")
	}
	app.Frame(460, 230)
	if !app.HandleKey(gpui.KeyEvent{Code: gpui.KeyTab}) || app.FocusID() != "decrement" {
		t.Fatal("Tab did not focus decrement")
	}
	if !app.HandleKey(gpui.KeyEvent{Code: gpui.KeyEnter}) {
		t.Fatal("Enter did not decrement")
	}
	app.Frame(460, 230)
	if !app.Click(50, 110) {
		t.Fatal("pointer did not activate increment")
	}
	value := 0
	view.count.Read(func(n *int) { value = *n })
	if value != 1 || app.FocusID() != "increment" {
		t.Fatalf("count = %d, focus = %q", value, app.FocusID())
	}
}
