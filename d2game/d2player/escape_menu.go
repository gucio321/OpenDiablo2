package d2player

import (
	"time"

	//"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2resource"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2gui"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2ui"
)

type (
	layoutID int
	optionID int
)

const (
	// UI
	labelGutter    = 10
	sidePanelsSize = 80
	pentSize       = 54
	menuSize       = 500
	spacerWidth    = 10

	// layouts
	noLayoutID layoutID = iota - 2
	saveLayoutID
	mainLayoutID
	optionsLayoutID
	soundOptionsLayoutID
	videoOptionsLayoutID
	automapOptionsLayoutID
	configureControlsLayoutID

	// audio
	optAudioSoundVolume optionID = iota
	optAudioMusicVolume
	optAudio3dSound
	optAudioHardwareAcceleration
	optAudioEnvEffects
	optAudioNpcSpeech
	// video
	optVideoResolution
	optVideoLightingQuality
	optVideoBlendedShadows
	optVideoPerspective
	optVideoGamma
	optVideoContrast
	// automap
	optAutomapSize
	optAutomapFade
	optAutomapCenterWhenCleared
	optAutomapShowParty
	optAutomapShowNames
)

const (
	singleFrame = time.Millisecond * 16
)

// NewEscapeMenu creates a new escape menu
func NewEscapeMenu(navigator d2interface.Navigator,
	renderer d2interface.Renderer,
	audioProvider d2interface.AudioProvider,
	uiManager *d2ui.UIManager,
	guiManager *d2gui.GuiManager,
	assetManager *d2asset.AssetManager,
	l d2util.LogLevel,
	keyMap *KeyMap,
) *EscapeMenu {
	m := &EscapeMenu{
		audioProvider: audioProvider,
		renderer:      renderer,
		navigator:     navigator,
		guiManager:    guiManager,
		uiManager:     uiManager,
		assetManager:  assetManager,
		keyMap:        keyMap,
	}

	keyBindingMenu := NewKeyBindingMenu(assetManager, renderer, uiManager, guiManager, keyMap, l, m)
	m.keyBindingMenu = keyBindingMenu

	m.layouts = make(map[layoutID]*d2ui.WidgetGroup)
	//m.layouts[mainLayoutID] = m.newMainLayout()
	//m.layouts[optionsLayoutID] = m.newOptionsLayout()
	//m.layouts[soundOptionsLayoutID] = m.newSoundOptionsLayout()
	//m.layouts[videoOptionsLayoutID] = m.newVideoOptionsLayout()
	//m.layouts[automapOptionsLayoutID] = m.newAutomapOptionsLayout()

	m.Logger = d2util.NewLogger()
	m.Logger.SetLevel(l)
	m.Logger.SetPrefix(logPrefix)

	return m
}

// EscapeMenu represents the in-game menu that shows up when the esc key is pressed
type EscapeMenu struct {
	isOpen        bool
	selectSound   d2interface.SoundEffect
	currentLayout layoutID

	// leftPent and rightPent are generated once and shared between the layouts
	leftPent  *d2ui.Sprite
	rightPent *d2ui.Sprite
	layouts   map[layoutID]*d2ui.WidgetGroup

	renderer       d2interface.Renderer
	audioProvider  d2interface.AudioProvider
	navigator      d2interface.Navigator
	guiManager     *d2gui.GuiManager
	uiManager      *d2ui.UIManager
	assetManager   *d2asset.AssetManager
	keyMap         *KeyMap
	keyBindingMenu *KeyBindingMenu

	onCloseCb func()

	*d2util.Logger
}

func (e *EscapeMenu) OnLoad() {
	for i := noLayoutID; i < configureControlsLayoutID; i++ {
		e.layouts[i] = e.uiManager.NewWidgetGroup(d2ui.RenderPriorityEscapeMenu)
	}

	l, _ := e.uiManager.NewSprite(d2resource.EscapeMenuSaveGame, d2resource.PaletteSky)
	l.SetPosition(200, 200)
	e.layouts[noLayoutID].AddWidget(l)
	e.showLayout(noLayoutID)
}

func (e *EscapeMenu) showLayout(id layoutID) {
	for i := noLayoutID; i < configureControlsLayoutID; i++ {
		e.layouts[i].SetVisible(i == id)
	}

}

func (e *EscapeMenu) SetOnCloseCb(cb func()) {
	e.onCloseCb = cb
}

func (e *EscapeMenu) IsOpen() bool {
	return e.isOpen
}

func (e *EscapeMenu) Render(target d2interface.Surface) {
	// no op
}

func (e *EscapeMenu) Advance(_ float64) error {
	return nil
}

func (e *EscapeMenu) open() {
	e.isOpen = true
}

func (e *EscapeMenu) OnEscKey() {
	// noop
}
