package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*AccordionContainer)(nil)

// AccordionContainer displays a list of AccordionItems.
// Each item is represented by a button that reveals a detailed view when tapped.
type AccordionContainer struct {
	BaseWidget
	Items     []*AccordionItem
	MultiOpen bool
}

// NewAccordionContainer creates a new accordion widget.
func NewAccordionContainer(items ...*AccordionItem) *AccordionContainer {
	a := &AccordionContainer{
		Items: items,
	}
	a.ExtendBaseWidget(a)
	return a
}

// MinSize returns the size that this widget should not shrink below.
func (a *AccordionContainer) MinSize() fyne.Size {
	a.ExtendBaseWidget(a)
	return a.BaseWidget.MinSize()
}

// Append adds the given item to this AccordionContainer.
func (a *AccordionContainer) Append(item *AccordionItem) {
	a.Items = append(a.Items, item)
	a.Refresh()
}

// Remove deletes the given item from this AccordionContainer.
func (a *AccordionContainer) Remove(item *AccordionItem) {
	for i, ai := range a.Items {
		if ai == item {
			a.RemoveIndex(i)
			break
		}
	}
}

// Remove deletes the item at the given index from this AccordionContainer.
func (a *AccordionContainer) RemoveIndex(index int) {
	a.Items = append(a.Items[:index], a.Items[index+1:]...)
	a.Refresh()
}

// Open expands the item at the given index.
func (a *AccordionContainer) Open(index int) {
	if index < 0 || index >= len(a.Items) {
		return
	}
	for i, ai := range a.Items {
		if i == index {
			ai.Open = true
		} else if !a.MultiOpen {
			ai.Open = false
		}
	}
	a.Refresh()
}

// OpenAll expands all items.
func (a *AccordionContainer) OpenAll() {
	if !a.MultiOpen {
		return
	}
	for _, i := range a.Items {
		i.Open = true
	}
	a.Refresh()
}

// Close collapses the item at the given index.
func (a *AccordionContainer) Close(index int) {
	if index < 0 || index >= len(a.Items) {
		return
	}
	a.Items[index].Open = false
	a.Refresh()
}

// CloseAll collapses all items.
func (a *AccordionContainer) CloseAll() {
	for _, i := range a.Items {
		i.Open = false
	}
	a.Refresh()
}

func (a *AccordionContainer) toggleForIndex(index int) func() {
	return func() {
		if a.Items[index].Open {
			a.Close(index)
		} else {
			a.Open(index)
		}
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (a *AccordionContainer) CreateRenderer() fyne.WidgetRenderer {
	a.ExtendBaseWidget(a)
	r := &accordionContainerRenderer{
		container: a,
	}
	r.updateObjects()
	return r
}

type accordionContainerRenderer struct {
	container *AccordionContainer
	headers   []*Button
}

func (r *accordionContainerRenderer) MinSize() fyne.Size {
	width := 0
	height := 0
	for i, ai := range r.container.Items {
		if i != 0 {
			height += theme.Padding()
		}
		min := r.headers[i].MinSize()
		width = fyne.Max(width, min.Width)
		height += min.Height
		if ai.Open {
			height += theme.Padding()
			min := ai.Detail.MinSize()
			width = fyne.Max(width, min.Width)
			height += min.Height
		}
	}
	return fyne.NewSize(width, height)
}

func (r *accordionContainerRenderer) Layout(size fyne.Size) {
	x := 0
	y := 0
	for i, ai := range r.container.Items {
		if i != 0 {
			y += theme.Padding()
		}
		h := r.headers[i]
		h.Move(fyne.NewPos(x, y))
		min := h.MinSize().Height
		h.Resize(fyne.NewSize(size.Width, min))
		y += min
		if ai.Open {
			y += theme.Padding()
			d := ai.Detail
			d.Move(fyne.NewPos(x, y))
			min := d.MinSize().Height
			d.Resize(fyne.NewSize(size.Width, min))
			y += min
		}
	}
}

func (r *accordionContainerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *accordionContainerRenderer) Objects() (objects []fyne.CanvasObject) {
	for _, h := range r.headers {
		objects = append(objects, h)
	}
	for _, i := range r.container.Items {
		objects = append(objects, i.Detail)
	}
	return
}

func (r *accordionContainerRenderer) Refresh() {
	r.updateObjects()
	r.Layout(r.container.Size())
	canvas.Refresh(r.container)
}

func (r *accordionContainerRenderer) updateObjects() {
	is := len(r.container.Items)
	hs := len(r.headers)
	i := 0
	for ; i < is; i++ {
		ai := r.container.Items[i]
		var h *Button
		if i < hs {
			h = r.headers[i]
		} else {
			h = &Button{}
			r.headers = append(r.headers, h)
		}
		h.Hidden = false
		h.Text = ai.Title
		h.OnTapped = r.container.toggleForIndex(i)
		if ai.Open {
			h.Icon = theme.MoveUpIcon()
			ai.Detail.Show()
		} else {
			h.Icon = theme.MoveDownIcon()
			ai.Detail.Hide()
		}
		h.Refresh()
	}
	// Hide extras
	for ; i < hs; i++ {
		r.headers[i].Hide()
	}
}

func (r *accordionContainerRenderer) Destroy() {
}

type AccordionItem struct {
	Title  string
	Detail fyne.CanvasObject
	Open   bool
}

func NewAccordionItem(title string, detail fyne.CanvasObject) *AccordionItem {
	return &AccordionItem{
		Title:  title,
		Detail: detail,
	}
}
