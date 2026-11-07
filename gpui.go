// Package gpui provides a Go-native, GPUI-inspired retained-state and
// immediate-view UI core. Render builds a fresh element tree when state changes.
// App, views, and entities must be used from the UI goroutine.
package gpui

import (
	"image/color"
)

// Rect uses logical pixels, independent of the display's scale factor.
type Rect struct{ X, Y, W, H float64 }

// Style describes the small layout and painting vocabulary implemented so far.
// A zero Width or Height means that the parent determines that dimension.
type Style struct {
	Width, Height float64
	Padding, Gap  float64
	Background    color.RGBA
	Foreground    color.RGBA
	FontSize      float64
}

type Kind uint8

const (
	ColumnKind Kind = iota
	RowKind
	TextKind
)

// Element is a declarative description, rebuilt by View.Render after updates.
type Element struct {
	Kind     Kind
	Style    Style
	Text     string
	Children []*Element
	OnClick  func()
	// FocusID is a stable, unique ID for an interactive element. It persists
	// across Render calls and enables Tab traversal and keyboard activation.
	FocusID string
}

func Column(children ...*Element) *Element { return &Element{Kind: ColumnKind, Children: children} }
func Row(children ...*Element) *Element    { return &Element{Kind: RowKind, Children: children} }
func Text(value string) *Element           { return &Element{Kind: TextKind, Text: value} }

// View owns state and describes its appearance. Render must not mutate state.
type View interface{ Render() *Element }

// App hosts one root view. The platform backend owns the event loop.
type App struct {
	root    View
	dirty   bool
	scene   Scene
	focusID string
}

func New(root View) *App { return &App{root: root, dirty: true} }

// Notify schedules a new view tree on the next Frame call.
func (a *App) Notify() { a.dirty = true }

// Entity holds application state. Update invalidates its owning app.
type Entity[T any] struct {
	app   *App
	value T
}

func NewEntity[T any](app *App, value T) *Entity[T] { return &Entity[T]{app: app, value: value} }

func (e *Entity[T]) Read(read func(*T)) { read(&e.value) }
func (e *Entity[T]) Update(update func(*T)) {
	update(&e.value)
	e.app.Notify()
}

// Frame computes a scene only when the view or viewport has changed.
func (a *App) Frame(width, height float64) Scene {
	if a.dirty || a.scene.Width != width || a.scene.Height != height {
		a.scene = buildScene(a.root.Render(), width, height)
		if a.focusID != "" && !a.scene.hasFocus(a.focusID) {
			a.focusID = ""
		}
		for _, target := range a.scene.focuses {
			if target.id == a.focusID {
				a.scene.Commands = append(a.scene.Commands, Command{Kind: FocusRect, Bounds: target.rect,
					Color: color.RGBA{R: 175, G: 215, B: 255, A: 255}})
				break
			}
		}
		a.dirty = false
	}
	return a.scene
}

// Click dispatches to the topmost hit-tested element, then invalidates the view.
// The caller should first call Frame to establish the current layout.
func (a *App) Click(x, y float64) bool {
	for i := len(a.scene.hits) - 1; i >= 0; i-- {
		hit := a.scene.hits[i]
		if x >= hit.rect.X && x < hit.rect.X+hit.rect.W && y >= hit.rect.Y && y < hit.rect.Y+hit.rect.H {
			a.focusID = hit.focusID
			hit.fn()
			a.Notify()
			return true
		}
	}
	if a.focusID != "" {
		a.focusID = ""
		a.Notify()
	}
	return false
}

// FocusID returns the currently focused element's stable ID, or an empty string.
func (a *App) FocusID() string { return a.focusID }

type KeyCode uint8

const (
	KeyTab KeyCode = iota
	KeyEnter
	KeySpace
)

type KeyEvent struct {
	Code  KeyCode
	Shift bool
}

// HandleKey moves focus with Tab/Shift-Tab or activates the focused element
// with Enter/Space. Call Frame before dispatching input to establish layout.
func (a *App) HandleKey(event KeyEvent) bool {
	switch event.Code {
	case KeyTab:
		if len(a.scene.focuses) == 0 {
			return false
		}
		index := -1
		for i, target := range a.scene.focuses {
			if target.id == a.focusID {
				index = i
				break
			}
		}
		if event.Shift {
			if index < 0 {
				index = 0
			}
			index = (index - 1 + len(a.scene.focuses)) % len(a.scene.focuses)
		} else {
			index = (index + 1) % len(a.scene.focuses)
		}
		a.focusID = a.scene.focuses[index].id
		a.Notify()
		return true
	case KeyEnter, KeySpace:
		for _, target := range a.scene.focuses {
			if target.id == a.focusID {
				target.fn()
				a.Notify()
				return true
			}
		}
	}
	return false
}
