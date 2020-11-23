package d2player

import (
	//"strconv"

	//"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2resource"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
	//"github.com/OpenDiablo2/OpenDiablo2/d2core/d2hero"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2ui"
)

const ( // for the dc6 frames
	moveGoldPanelTopLeft = iota
	moveGoldPanelTopRight
	moveGoldPanelBottomLeft
	moveGoldPanelBottomRight
)

const (
	moveGoldPanelOffsetX, moveGoldPanelOffsetY = 80, 64
)

const (
	moveGoldActionLabelX, moveGoldActionLabelY = 110, 100
)

const (
	moveGoldCloseButtonX, moveGoldCloseButtonY = 208, 453
	moveGoldOkButtonX, moveGoldOkButtonY       = 208, 453
)

// NewHeroStatsPanel creates a new hero status panel
func NewMoveGoldPanel(
	asset *d2asset.AssetManager,
	ui *d2ui.UIManager,
	l d2util.LogLevel) *MoveGoldPanel {
	originX := 0
	originY := 0

	mg := &MoveGoldPanel{
		asset:     asset,
		uiManager: ui,
		originX:   originX,
		originY:   originY,
	}

	mg.Logger = d2util.NewLogger()
	mg.Logger.SetLevel(l)
	mg.Logger.SetPrefix(logPrefix)

	return mg
}

// HeroStatsPanel represents the hero status panel
type MoveGoldPanel struct {
	asset               *d2asset.AssetManager
	uiManager           *d2ui.UIManager
	closeButton         *d2ui.Button
	okButton            *d2ui.Button
	panel               *d2ui.Sprite
	moveGoldActionLabel *d2ui.Label
	onCloseCb           func()
	panelGroup          *d2ui.WidgetGroup

	originX  int
	originY  int
	isOpen   bool
	isChest  bool
	wantDrop bool

	*d2util.Logger
}

// Load the data for the hero status panel
func (s *MoveGoldPanel) Load() {
	var err error

	s.panelGroup = s.uiManager.NewWidgetGroup(d2ui.RenderPriorityHeroStatsPanel)

	frame := d2ui.NewUIFrame(s.asset, s.uiManager, d2ui.FrameLeft)
	s.panelGroup.AddWidget(frame)

	s.panel, err = s.uiManager.NewSprite(d2resource.PopUpOkCancel, d2resource.PaletteSky)
	if err != nil {
		s.Error(err.Error())
	}

	w, h := frame.GetSize()
	staticPanel := s.uiManager.NewCustomWidgetCached(s.renderStaticMenu, w, h)
	s.panelGroup.AddWidget(staticPanel)

	closeButton := s.uiManager.NewButton(d2ui.ButtonTypeSquareClose, "")
	closeButton.SetVisible(false)
	closeButton.SetPosition(moveGoldCloseButtonX, moveGoldCloseButtonY)
	closeButton.OnActivated(func() { s.Close() })
	s.panelGroup.AddWidget(closeButton)

	okButton := s.uiManager.NewButton(d2ui.ButtonTypeSquareClose, "")
	okButton.SetVisible(false)
	okButton.SetPosition(moveGoldOkButtonX, moveGoldOkButtonY)
	okButton.OnActivated(func() { s.Close() })
	s.panelGroup.AddWidget(okButton)

	s.moveGoldActionLabel = s.uiManager.NewLabel(d2resource.Font16, d2resource.PaletteStatic)
	s.moveGoldActionLabel.Alignment = d2ui.HorizontalAlignCenter
	s.moveGoldActionLabel.SetText("lalala")
	s.moveGoldActionLabel.SetPosition(moveGoldActionLabelX, moveGoldActionLabelY)
	s.panelGroup.AddWidget(s.moveGoldActionLabel)

	s.panelGroup.SetVisible(true)
}

// IsOpen returns true if the hero status panel is open
func (s *MoveGoldPanel) IsOpen() bool {
	return s.isOpen
}

// Toggle toggles the visibility of the hero status panel
func (s MoveGoldPanel) Toggle() {
	if s.isOpen {
		s.Close()
	} else {
		s.Open()
	}
}

// Open opens the hero status panel
func (s *MoveGoldPanel) Open() {
	s.isOpen = true
	s.panelGroup.SetVisible(true)
}

// Close closed the hero status panel
func (s *MoveGoldPanel) Close() {
	s.isOpen = false
	s.panelGroup.SetVisible(false)
	s.onCloseCb()
}

// SetOnCloseCb the callback run on closing the HeroStatsPanel
func (s *MoveGoldPanel) SetOnCloseCb(cb func()) {
	s.onCloseCb = cb
}

// Advance updates labels on the panel
func (s *MoveGoldPanel) Advance(elapsed float64) {
	if !s.isOpen {
		return
	}
}

func (s *MoveGoldPanel) renderStaticMenu(target d2interface.Surface) {
	s.renderStaticPanelFrames(target)
	s.moveGoldActionLabel.Render(target)
}

func (s *MoveGoldPanel) renderStaticPanelFrames(target d2interface.Surface) {
	frames := []int{
		statsPanelTopLeft,
		statsPanelTopRight,
		statsPanelBottomRight,
		statsPanelBottomLeft,
	}

	currentX := s.originX + statsPanelOffsetX
	currentY := s.originY + statsPanelOffsetY

	for _, frameIndex := range frames {
		if err := s.panel.SetCurrentFrame(frameIndex); err != nil {
			s.Error(err.Error())
		}

		w, h := s.panel.GetCurrentFrameSize()

		switch frameIndex {
		case statsPanelTopLeft:
			s.panel.SetPosition(currentX, currentY+h)
			currentX += w
		case statsPanelTopRight:
			s.panel.SetPosition(currentX, currentY+h)
			currentY += h
		case statsPanelBottomRight:
			s.panel.SetPosition(currentX, currentY+h)
		case statsPanelBottomLeft:
			s.panel.SetPosition(currentX-w, currentY+h)
		}

		s.panel.Render(target)
	}
}

/*
func (s *HeroStatsPanel) initStatValueLabels() {
	valueLabelConfigs := []struct {
		assignTo **d2ui.Label
		value    int
		x, y     int
	}{
		{&s.labels.Level, s.heroState.Level, 112, 110},
		{&s.labels.Experience, s.heroState.Experience, 200, 110},
		{&s.labels.NextLevelExp, s.heroState.NextLevelExp, 330, 110},
		{&s.labels.Strength, s.heroState.Strength, 175, 147},
		{&s.labels.Dexterity, s.heroState.Dexterity, 175, 207},
		{&s.labels.Vitality, s.heroState.Vitality, 175, 295},
		{&s.labels.Energy, s.heroState.Energy, 175, 355},
		{&s.labels.MaxStamina, s.heroState.MaxStamina, 330, 295},
		{&s.labels.Stamina, int(s.heroState.Stamina), 370, 295},
		{&s.labels.MaxHealth, s.heroState.MaxHealth, 330, 320},
		{&s.labels.Health, s.heroState.Health, 370, 320},
		{&s.labels.MaxMana, s.heroState.MaxMana, 330, 355},
		{&s.labels.Mana, s.heroState.Mana, 370, 355},
	}

	for _, cfg := range valueLabelConfigs {
		*cfg.assignTo = s.createStatValueLabel(cfg.value, cfg.x, cfg.y)
	}
}

func (s *HeroStatsPanel) setStatValues() {
	s.labels.Level.SetText(strconv.Itoa(s.heroState.Level))
	s.labels.Experience.SetText(strconv.Itoa(s.heroState.Experience))
	s.labels.NextLevelExp.SetText(strconv.Itoa(s.heroState.NextLevelExp))

	s.labels.Strength.SetText(strconv.Itoa(s.heroState.Strength))
	s.labels.Dexterity.SetText(strconv.Itoa(s.heroState.Dexterity))
	s.labels.Vitality.SetText(strconv.Itoa(s.heroState.Vitality))
	s.labels.Energy.SetText(strconv.Itoa(s.heroState.Energy))

	s.labels.MaxHealth.SetText(strconv.Itoa(s.heroState.MaxHealth))
	s.labels.Health.SetText(strconv.Itoa(s.heroState.Health))

	s.labels.MaxStamina.SetText(strconv.Itoa(s.heroState.MaxStamina))
	s.labels.Stamina.SetText(strconv.Itoa(int(s.heroState.Stamina)))

	s.labels.MaxMana.SetText(strconv.Itoa(s.heroState.MaxMana))
	s.labels.Mana.SetText(strconv.Itoa(s.heroState.Mana))
}

func (s *HeroStatsPanel) createStatValueLabel(stat, x, y int) *d2ui.Label {
	text := strconv.Itoa(stat)
	return s.createTextLabel(PanelText{X: x, Y: y, Text: text, Font: d2resource.Font16, AlignCenter: true})
}

func (s *HeroStatsPanel) createTextLabel(element PanelText) *d2ui.Label {
	label := s.uiManager.NewLabel(element.Font, d2resource.PaletteStatic)
	if element.AlignCenter {
		label.Alignment = d2ui.HorizontalAlignCenter
	}

	label.SetText(element.Text)
	label.SetPosition(element.X, element.Y)
	s.panelGroup.AddWidget(label)

	return label
}*/
