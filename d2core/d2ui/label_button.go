package d2ui

import (
	"image/color"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
)

const (
	white = 0xffffffff
)

// static checks to ensure LabelButton implemented Widget and ClickableWidget
var _ Widget = &LabelButton{}
var _ ClickableWidget = &LabelButton{}

// LabelButton represents LabelButton
type LabelButton struct {
	*BaseWidget
	alignment  HorizontalAlign
	stdColor   color.Color
	hoverColor color.Color
	onClick    func()
	label      *Label

	*d2util.Logger
}

func (ui *UIManager) NewLabelButton(font, palette string) *LabelButton {
	base := NewBaseWidget(ui)
	base.SetVisible(true)

	result := &LabelButton{
		BaseWidget: base,
		stdColor:   d2util.Color(white),
	}

	result.label = ui.NewLabel(font, palette)
	result.label.Alignment = HorizontalAlignCenter
	result.label.Color[0] = result.stdColor

	ui.addWidget(result)

	return result
}

func (b *LabelButton) SetText(text string) {
	b.label.SetText(text)
	b.width, b.height = b.label.GetSize()
}

func (b *LabelButton) SetColors(normColor, hoverColor color.Color) {
	b.stdColor = normColor
	b.hoverColor = hoverColor
}

func (b *LabelButton) GetSize() (x, y int) {
	return b.label.GetSize()
}

// GetPosition returns real position (including offset for the alignment)
func (b *LabelButton) GetPosition() (x, y int) {
	return b.x - b.label.getAlignOffset(b.width), b.y
}

// OnActivated defines the callback handler for the activate event
func (b *LabelButton) OnActivated(cb func()) {
	b.onClick = cb
}

// Activate calls the on activated callback handler, if any
func (b *LabelButton) Activate() {
	if b.onClick == nil {
		return
	}

	b.onClick()
}

// SetEnabled sets the enabled state
func (b *LabelButton) SetEnabled(_ bool) {
	// noop
}

// GetEnabled returns the enabled state
func (b *LabelButton) GetEnabled() bool {
	return true
}

func (b *LabelButton) SetPressed(_ bool) {
	// noop
}

func (b *LabelButton) GetPressed() bool {
	return false
}

func (b *LabelButton) Advance(_ float64) error {
	return nil
}

func (b *LabelButton) Render(target d2interface.Surface) {
	target.PushTranslation(b.GetPosition())
	defer target.Pop()

	b.label.Render(target)

	if b.isHovered() {
		b.label.Color[0] = b.hoverColor
	} else {
		b.label.Color[0] = b.stdColor
	}
}
