package d2ui

import (
	"image/color"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
)

const (
	white = 0xffffffff
)

var _ Widget = &LabelButton{}

// LabelButton represents LabelButton
type LabelButton struct {
	*BaseWidget
	text       string
	alignment  HorizontalAlign
	font       *d2asset.Font
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
	result.label.Alignment = HorizontalAlignLeft
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

func (b *LabelButton) OnActivated(cb func()) {
	b.onClick = cb
}

func (b *LabelButton) Activate() {
	if b.onClick == nil {
		return
	}

	b.onClick()
}

func (b *LabelButton) SetEnabled(_ bool) {
	// noop
}

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
	if b.IsHovered() {
		b.label.Color[0] = b.hoverColor
	} else {
		b.label.Color[0] = b.stdColor
	}

}
