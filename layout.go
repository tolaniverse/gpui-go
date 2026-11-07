package gpui

import "image/color"

type CommandKind uint8

const (
	FillRect CommandKind = iota
	DrawText
)

type Command struct {
	Kind     CommandKind
	Bounds   Rect
	Color    color.RGBA
	Text     string
	FontSize float64
}

type hitTarget struct {
	rect    Rect
	focusID string
	fn      func()
}

type focusTarget struct {
	id   string
	rect Rect
	fn   func()
}

// Scene is an ordered list of paint commands in logical pixels.
type Scene struct {
	Width, Height float64
	Commands      []Command
	hits          []hitTarget
	focuses       []focusTarget
}

func (s Scene) hasFocus(id string) bool {
	for _, target := range s.focuses {
		if target.id == id {
			return true
		}
	}
	return false
}

func buildScene(root *Element, width, height float64) Scene {
	scene := Scene{Width: width, Height: height}
	layout(&scene, root, Rect{W: width, H: height})
	return scene
}

func layout(scene *Scene, node *Element, bounds Rect) {
	if node == nil || bounds.W <= 0 || bounds.H <= 0 {
		return
	}
	style := node.Style
	if style.Background.A != 0 {
		scene.Commands = append(scene.Commands, Command{Kind: FillRect, Bounds: bounds, Color: style.Background})
	}
	if node.OnClick != nil {
		scene.hits = append(scene.hits, hitTarget{rect: bounds, focusID: node.FocusID, fn: node.OnClick})
		if node.FocusID != "" {
			scene.focuses = append(scene.focuses, focusTarget{id: node.FocusID, rect: bounds, fn: node.OnClick})
		}
	}
	inner := Rect{X: bounds.X + style.Padding, Y: bounds.Y + style.Padding,
		W: max(0, bounds.W-2*style.Padding), H: max(0, bounds.H-2*style.Padding)}
	if node.Kind == TextKind {
		foreground := style.Foreground
		if foreground.A == 0 {
			foreground = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		fontSize := style.FontSize
		if fontSize <= 0 {
			fontSize = 16
		}
		scene.Commands = append(scene.Commands, Command{Kind: DrawText, Bounds: inner,
			Color: foreground, Text: node.Text, FontSize: fontSize})
		return
	}
	// Fixed-size children claim space first. Remaining space is divided equally
	// among auto-sized children. Children are clipped to the parent's bounds by
	// the backend; overflow and intrinsic text measurement are not implemented.
	visible := make([]*Element, 0, len(node.Children))
	for _, child := range node.Children {
		if child != nil {
			visible = append(visible, child)
		}
	}
	if len(visible) == 0 {
		return
	}
	row := node.Kind == RowKind
	available := inner.H
	if row {
		available = inner.W
	}
	available = max(0, available-style.Gap*float64(len(visible)-1))
	fixed, auto := 0.0, 0
	for _, child := range visible {
		size := child.Style.Height
		if row {
			size = child.Style.Width
		}
		if size > 0 {
			fixed += size
		} else {
			auto++
		}
	}
	share := 0.0
	if auto > 0 {
		share = max(0, available-fixed) / float64(auto)
	}
	pos := inner.Y
	if row {
		pos = inner.X
	}
	for _, child := range visible {
		size := child.Style.Height
		if row {
			size = child.Style.Width
		}
		if size <= 0 {
			size = share
		}
		cross := inner.W
		if row {
			cross = inner.H
			if child.Style.Height > 0 {
				cross = min(cross, child.Style.Height)
			}
		} else if child.Style.Width > 0 {
			cross = min(cross, child.Style.Width)
		}
		if row {
			layout(scene, child, Rect{X: pos, Y: inner.Y, W: min(size, max(0, inner.X+inner.W-pos)), H: cross})
		} else {
			layout(scene, child, Rect{X: inner.X, Y: pos, W: cross, H: min(size, max(0, inner.Y+inner.H-pos))})
		}
		pos += size + style.Gap
	}
}
