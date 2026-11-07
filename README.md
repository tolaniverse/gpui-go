# gpui-go

A **native Go**, GPUI-inspired desktop UI framework experiment. This is a new
implementation, not Rust bindings or a feature-complete port of GPUI.

The first vertical slice has Go-owned entity state, declarative views, a small
row/column element layout, paint commands, pointer hit-testing, and a desktop
counter example. Ebitengine is the Go desktop backend for
initial window/input/GPU/text backend. Framework state and layout do not depend
on Ebitengine, so other backends can be added later. No Rust is used. The
packages build and test with `CGO_ENABLED=0` on macOS; the desktop window has
not yet been interactively verified.

```sh
go test ./...
go run ./examples/counter
```

In the counter, click a button or press Tab/Shift-Tab to move focus and
Enter/Space to activate it. Give each keyboard-activatable element a stable,
unique `FocusID` and an `OnClick` callback. Focus follows clicks and survives
view re-renders while that ID remains in the tree.

Current limitations: one window and root view; no text measurement, clipping,
editable text, IME, accessibility, scrolling, or GPUI-compatible API.
The simple layout distributes unallocated space equally among auto-sized children.
The included font is Go Regular; system fonts are not yet integrated.

Next milestone: editable text and IME, then measured layout and clipping.

Reference architecture: [GPUI README](https://github.com/zed-industries/zed/blob/main/crates/gpui/README.md).
