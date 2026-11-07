package gpui

import "testing"

type focusView struct {
	ids       []string
	activated []string
}

func (v *focusView) Render() *Element {
	root := Column()
	for _, id := range v.ids {
		id := id
		child := Text(id)
		child.Style.Height = 30
		child.FocusID = id
		child.OnClick = func() { v.activated = append(v.activated, id) }
		root.Children = append(root.Children, child)
	}
	return root
}

func TestFocusTraversalAndActivation(t *testing.T) {
	view := &focusView{ids: []string{"first", "second"}}
	app := New(view)
	app.Frame(100, 60)
	if app.HandleKey(KeyEvent{Code: KeyEnter}) {
		t.Fatal("nothing should activate before focus")
	}
	for _, step := range []struct {
		key  KeyEvent
		want string
	}{
		{KeyEvent{Code: KeyTab}, "first"},
		{KeyEvent{Code: KeyTab}, "second"},
		{KeyEvent{Code: KeyTab}, "first"},
		{KeyEvent{Code: KeyTab, Shift: true}, "second"},
	} {
		if !app.HandleKey(step.key) || app.FocusID() != step.want {
			t.Fatalf("key %+v: focus = %q, want %q", step.key, app.FocusID(), step.want)
		}
		app.Frame(100, 60)
		commands := app.Frame(100, 60).Commands
		if commands[len(commands)-1].Kind != FocusRect {
			t.Fatal("focused element has no focus indicator")
		}
	}
	if !app.HandleKey(KeyEvent{Code: KeySpace}) || len(view.activated) != 1 || view.activated[0] != "second" {
		t.Fatalf("space activated %+v", view.activated)
	}
	app.Frame(100, 60)
	if !app.HandleKey(KeyEvent{Code: KeyEnter}) || len(view.activated) != 2 {
		t.Fatalf("enter activated %+v", view.activated)
	}
}

func TestFocusSurvivesRenderAndClearsWhenRemoved(t *testing.T) {
	view := &focusView{ids: []string{"first", "second"}}
	app := New(view)
	app.Frame(100, 60)
	if !app.Click(5, 35) || app.FocusID() != "second" {
		t.Fatalf("click focus = %q", app.FocusID())
	}
	app.Frame(100, 60)
	if app.FocusID() != "second" {
		t.Fatal("focus lost across re-render")
	}
	view.ids = []string{"first"}
	app.Notify()
	app.Frame(100, 60)
	if app.FocusID() != "" || app.HandleKey(KeyEvent{Code: KeySpace}) {
		t.Fatal("removed element retained focus")
	}
	if !app.HandleKey(KeyEvent{Code: KeyTab, Shift: true}) || app.FocusID() != "first" {
		t.Fatal("reverse traversal did not select last focusable element")
	}
	if app.Click(90, 50) || app.FocusID() != "" {
		t.Fatal("clicking empty space should blur")
	}
}
