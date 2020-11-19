package d2player

import (
	"fmt"
	"log"

	//"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2resource"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
	//"github.com/OpenDiablo2/OpenDiablo2/d2core/d2hero"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2ui"
)

const (
	msgLogX, msgLogY                 = 401, 222 //toset
	msgTitleX, msgTitleY             = 100, 100 //toset
	msgStartX, msgStartY             = 50, 120
	msgCloseButtonX, msgCloseButtonY = 416, 22 //toset
)

type MessageLog struct {
	title       *d2ui.Label
	closeButton *d2ui.Button
	messages    []string
	panel       *d2ui.Sprite
	isOpen      bool
	onCloseCb   func()

	asset      *d2asset.AssetManager
	uiManager  *d2ui.UIManager
	panelGroup *d2ui.WidgetGroup
}

func newMessageLog(
	messages []string,
	asset *d2asset.AssetManager,
	ui *d2ui.UIManager,
) *MessageLog {
	ml := &MessageLog{
		messages:  messages,
		asset:     asset,
		uiManager: ui,
		isOpen:    false,
	}

	return ml
}

func (m *MessageLog) load() {
	var err error
	frame := d2ui.NewUIFrame(m.asset, m.uiManager, d2ui.FrameLeft)
	m.panelGroup.AddWidget(frame)

	m.panel, err = m.uiManager.NewSprite(d2resource.BoxPieces, d2resource.PaletteSky)
	if err != nil {
		log.Print(err)
	}

	panel := m.uiManager.NewCustomWidget(m.Render, 400, 600)
	m.panelGroup.AddWidget(panel)

	m.closeButton = m.uiManager.NewButton(d2ui.ButtonTypeSquareClose, "")
	m.closeButton.SetVisible(false)
	m.closeButton.OnActivated(func() { m.Close() })
	m.panelGroup.AddWidget(m.closeButton)

	m.panelGroup.SetVisible(false)
}

func (s *MessageLog) IsOpen() bool {
	return s.isOpen
}

// Toggle the skill tree visibility
func (m *MessageLog) Toggle() {
	fmt.Println("MessageLog Toggled")

	if m.isOpen == true {
		m.Close()
	} else {
		m.Open()
	}
}

// Open the skill tree
func (m *MessageLog) Open() {
	m.isOpen = true

	m.panelGroup.SetVisible(true)
}

// Close the skill tree
func (m *MessageLog) Close() {
	m.isOpen = false

	m.panelGroup.SetVisible(false)

	m.onCloseCb()
}

// Set the callback run on closing the skilltree
func (s *MessageLog) SetOnCloseCb(cb func()) {
	s.onCloseCb = cb
}

func (m *MessageLog) Advance() {
	if !m.isOpen {
		return
	}
}

func (m *MessageLog) renderStaticMenu(target d2interface.Surface) {
	m.renderStaticPanelFrames(target)
}

func (m *MessageLog) renderStaticPanelFrames(target d2interface.Surface) {
	frames := []int{
		msgLogX,
		msgLogY,
	}
	for _, frameIndex := range frames {
		if err := m.panel.SetCurrentFrame(frameIndex); err != nil {
			log.Printf("%e", err)
		}
		w, h := m.panel.GetCurrentFrameSize()
		switch frameIndex {
		case statsPanelTopLeft:
			m.panel.SetPosition(msgLogX, msgLogY+h)
			msgLogX += w
		case statsPanelTopRight:
			m.panel.SetPosition(msgLogX, msgLogY+h)
			msgLogY += h
		case statsPanelBottomRight:
			m.panel.SetPosition(msgLogX, msgLogY+h)
		case statsPanelBottomLeft:
			m.panel.SetPosition(msgLogX-w, msgLogY+h)
		}
		m.panel.Render(target)
	}
}

// Render the skill tree panel
func (m *MessageLog) Render(target d2interface.Surface) {
	//m.renderTabCommon(target)
	//m.renderTab(target, m.selectedTab)
}

func (s *MessageLog) renderPanelSegment(
	target d2interface.Surface,
	frame int) {
	if err := s.resources.skillPanel.SetCurrentFrame(frame); err != nil {
		log.Printf("%e", err)
		return
	}

	s.resources.skillPanel.Render(target)
}

/*
func (s *skillTree) loadForHeroType() {
	sp, err := s.uiManager.NewSprite(s.resources.skillPanelPath, d2resource.PaletteSky)
	if err != nil {
		log.Print(err)
	}

	s.resources.skillPanel = sp

	si, err := s.uiManager.NewSprite(s.resources.skillIconPath, d2resource.PaletteSky)
	if err != nil {
		log.Print(err)
	}

	s.resources.skillSprite = si

	s.tab[firstTab].createButton(s.uiManager, tabButtonX, tabButton0Y)
	s.tab[firstTab].button.OnActivated(func() { s.setTab(firstTab) })
	s.panelGroup.AddWidget(s.tab[firstTab].button)

	s.tab[secondTab].createButton(s.uiManager, tabButtonX, tabButton1Y)
	s.tab[secondTab].button.OnActivated(func() { s.setTab(secondTab) })
	s.panelGroup.AddWidget(s.tab[secondTab].button)

	s.tab[thirdTab].createButton(s.uiManager, tabButtonX, tabButton2Y)
	s.tab[thirdTab].button.OnActivated(func() { s.setTab(thirdTab) })
	s.panelGroup.AddWidget(s.tab[thirdTab].button)

	s.availSPLabel = s.uiManager.NewLabel(d2resource.Font16, d2resource.PaletteSky)
	s.availSPLabel.SetPosition(availSPLabelX, availSPLabelY)
	s.availSPLabel.Alignment = d2ui.HorizontalAlignCenter
	s.availSPLabel.SetText(s.makeTabString("StrSklTree1", "StrSklTree2", "StrSklTree3"))
	s.panelGroup.AddWidget(s.availSPLabel)
}

type heroTabData struct {
	resources        *skillTreeHeroTypeResources
	str1, str2, str3 string
	closeButtonPos   [numTabs]int
}

func (s *skillTree) makeTabString(keys ...interface{}) string {
	translations := make([]interface{}, len(keys))

	token := "%s"
	format := token

	for idx, key := range keys {
		if idx > 0 {
			format += "\n" + token
		}

		translations[idx] = s.asset.TranslateString(key.(string))
	}

	return fmt.Sprintf(format, translations...)
}

func makeCloseButtonPos(close1, close2, close3 int) [numTabs]int {
	return [numTabs]int{close1, close2, close3}
}

func (s *skillTree) getTab(class d2enum.Hero) *heroTabData {
	tabMap := map[d2enum.Hero]*heroTabData{
		d2enum.HeroBarbarian: {
			&skillTreeHeroTypeResources{
				skillPanelPath: d2resource.SkillsPanelBarbarian,
				skillIconPath:  d2resource.BarbarianSkills,
			},
			s.makeTabString("StrSklTree21", "StrSklTree4"),
			s.makeTabString("StrSklTree21", "StrSklTree22"),
			s.makeTabString("StrSklTree20"),
			makeCloseButtonPos(
				skillCloseButtonXRight,
				skillCloseButtonXLeft,
				skillCloseButtonXRight),
		},
		d2enum.HeroNecromancer: {
			&skillTreeHeroTypeResources{
				skillPanelPath: d2resource.SkillsPanelNecromancer,
				skillIconPath:  d2resource.NecromancerSkills,
			},
			s.makeTabString("StrSklTree19"),
			s.makeTabString("StrSklTree17", "StrSklTree18", "StrSklTree5"),
			s.makeTabString("StrSklTree16", "StrSklTree5"),
			makeCloseButtonPos(
				skillCloseButtonXLeft,
				skillCloseButtonXRight,
				skillCloseButtonXLeft),
		},
		d2enum.HeroPaladin: {
			&skillTreeHeroTypeResources{
				skillPanelPath: d2resource.SkillsPanelPaladin,
				skillIconPath:  d2resource.PaladinSkills,
			},
			s.makeTabString("StrSklTree15", "StrSklTree4"),
			s.makeTabString("StrSklTree14", "StrSklTree13"),
			s.makeTabString("StrSklTree12", "StrSklTree13"),
			makeCloseButtonPos(
				skillCloseButtonXLeft,
				skillCloseButtonXMiddle,
				skillCloseButtonXLeft),
		},
		d2enum.HeroAssassin: {
			&skillTreeHeroTypeResources{
				skillPanelPath: d2resource.SkillsPanelAssassin,
				skillIconPath:  d2resource.AssassinSkills,
			},
			s.makeTabString("StrSklTree30"),
			s.makeTabString("StrSklTree31", "StrSklTree32"),
			s.makeTabString("StrSklTree33", "StrSklTree34"),
			makeCloseButtonPos(
				skillCloseButtonXMiddle,
				skillCloseButtonXRight,
				skillCloseButtonXLeft),
		},
		d2enum.HeroSorceress: {
			&skillTreeHeroTypeResources{
				skillPanelPath: d2resource.SkillsPanelSorcerer,
				skillIconPath:  d2resource.SorcererSkills,
			},
			s.makeTabString("StrSklTree25", "StrSklTree5"),
			s.makeTabString("StrSklTree24", "StrSklTree5"),
			s.makeTabString("StrSklTree23", "StrSklTree5"),
			makeCloseButtonPos(
				skillCloseButtonXLeft,
				skillCloseButtonXLeft,
				skillCloseButtonXRight),
		},
		d2enum.HeroAmazon: {
			&skillTreeHeroTypeResources{
				skillPanelPath: d2resource.SkillsPanelAmazon,
				skillIconPath:  d2resource.AmazonSkills,
			},
			s.makeTabString("StrSklTree10", "StrSklTree11", "StrSklTree4"),
			s.makeTabString("StrSklTree8", "StrSklTree9", "StrSklTree4"),
			s.makeTabString("StrSklTree6", "StrSklTree7", "StrSklTree4"),
			makeCloseButtonPos(
				skillCloseButtonXRight,
				skillCloseButtonXMiddle,
				skillCloseButtonXLeft),
		},
		d2enum.HeroDruid: {
			&skillTreeHeroTypeResources{
				skillPanelPath: d2resource.SkillsPanelDruid,
				skillIconPath:  d2resource.DruidSkills,
			},
			s.makeTabString("StrSklTree26"),
			s.makeTabString("StrSklTree27", "StrSklTree28"),
			s.makeTabString("StrSklTree29"),
			makeCloseButtonPos(
				skillCloseButtonXRight,
				skillCloseButtonXRight,
				skillCloseButtonXRight),
		},
	}

	return tabMap[class]
}

func (s *skillTree) setHeroTypeResourcePath() {
	entry := s.getTab(s.heroClass)
	if entry == nil {
		log.Fatal("Unknown Hero Type")
	}

	s.resources = entry.resources
	s.tab[firstTab].buttonText = entry.str1
	s.tab[secondTab].buttonText = entry.str2
	s.tab[thirdTab].buttonText = entry.str3

	for i := 0; i < numTabs; i++ {
		s.tab[i].closeButtonPosX = entry.closeButtonPos[i]
	}
}






func (s *skillTree) setTab(tab int) {
	s.selectedTab = tab
	s.closeButton.SetPosition(s.tab[tab].closeButtonPosX, skillCloseButtonY)

	for _, si := range s.skillIcons {
		si.SetVisible(si.skill.SkillPage == tab+1)
	}
}


func (s *skillTree) renderTabCommon(target d2interface.Surface) {
	skillPanel := s.resources.skillPanel
	x, y := s.originX, s.originY

	// top
	w, h, err := skillPanel.GetFrameSize(frameCommonTabTopLeft)
	if err != nil {
		log.Printf("%e", err)
		return
	}

	y += h

	skillPanel.SetPosition(x, y)
	s.renderPanelSegment(target, frameCommonTabTopLeft)

	skillPanel.SetPosition(x+w, y)
	s.renderPanelSegment(target, frameCommonTabTopRight)

	// bottom
	_, h, err = skillPanel.GetFrameSize(frameCommonTabBottomLeft)
	if err != nil {
		log.Printf("%e", err)
		return
	}

	y += h

	skillPanel.SetPosition(x, y)
	s.renderPanelSegment(target, frameCommonTabBottomLeft)

	skillPanel.SetPosition(x+w, y)
	s.renderPanelSegment(target, frameCommonTabBottomRight)
}

func (s *skillTree) renderTab(target d2interface.Surface, tab int) {
	topFrame := frameOffsetTop + (tabIndexOffset * tab)
	bottomFrame := frameOffsetBottom + (tabIndexOffset * tab)

	skillPanel := s.resources.skillPanel
	x, y := s.originX, s.originY

	// top
	_, h0, err := skillPanel.GetFrameSize(topFrame)
	if err != nil {
		log.Printf("%e", err)
		return
	}

	y += h0

	skillPanel.SetPosition(x, y)
	s.renderPanelSegment(target, topFrame)

	// bottom
	w, h1, err := skillPanel.GetFrameSize(bottomFrame)
	if err != nil {
		log.Printf("%e", err)
		return
	}

	skillPanel.SetPosition(x, y+h1)

	s.renderPanelSegment(target, bottomFrame)

	// tab button highlighted
	switch tab {
	case firstTab:
		skillPanel.SetPosition(x+w, y+h1)
		s.renderPanelSegment(target, frameSelectedTab1Full)
	case secondTab:
		x += w
		skillPanel.SetPosition(x, s.originY+h0)
		s.renderPanelSegment(target, frameSelectedTab2Top)

		skillPanel.SetPosition(x, y+h1)
		s.renderPanelSegment(target, frameSelectedTab2Bottom)
	case thirdTab:
		skillPanel.SetPosition(x+w, y)
		s.renderPanelSegment(target, frameSelectedTab3Full)
	}
}
*/
