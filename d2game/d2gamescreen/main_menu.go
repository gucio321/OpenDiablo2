// Package d2gamescreen contains the screens
package d2gamescreen

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2hero"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2resource"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2screen"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2ui"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2client/d2clientconnectiontype"
	"github.com/OpenDiablo2/OpenDiablo2/d2script"
)

type mainMenuScreenMode int

// mainMenuScreenMode types
const (
	ScreenModeUnknown mainMenuScreenMode = iota
	ScreenModeTrademark
	ScreenModeMainMenu
	ScreenModeMultiplayer
	ScreenModeTCPIP
	ScreenModeServerIP
)

const (
	joinGameDialogX, joinGameDialogY         = 318, 245
	serverIPbackgroundX, serverIPbackgroundY = 270, 175
	backgroundX, backgroundY                 = 0, 0
	versionLabelX, versionLabelY             = 795, -10
	commitLabelX, commitLabelY               = 2, 2
	copyrightX, copyrightY                   = 400, 500
	copyright2X, copyright2Y                 = 400, 525
	od2LabelX, od2LabelY                     = 400, 580
	tcpOptionsX, tcpOptionsY                 = 400, 23
	joinGameX, joinGameY                     = 400, 190
	diabloLogoX, diabloLogoY                 = 400, 120
	exitDiabloBtnX, exitDiabloBtnY           = 264, 535
	creditBtnX, creditBtnY                   = 264, 505
	cineBtnX, cineBtnY                       = 401, 505
	singlePlayerBtnX, singlePlayerBtnY       = 264, 290
	githubBtnX, githubBtnY                   = 264, 400
	mapTestBtnX, mapTestBtnY                 = 264, 440
	tcpBtnX, tcpBtnY                         = 33, 543
	srvCancelBtnX, srvCancelBtnY             = 285, 305
	srvOkBtnX, srvOkBtnY                     = 420, 305
	multiplayerBtnX, multiplayerBtnY         = 264, 330
	tcpNetBtnX, tcpNetBtnY                   = 264, 280
	networkCancelBtnX, networkCancelBtnY     = 264, 540
	tcpHostBtnX, tcpHostBtnY                 = 264, 200
	tcpJoinBtnX, tcpJoinBtnY                 = 264, 240
	errorLabelX, errorLabelY                 = 400, 250
	machineIPX, machineIPY                   = 400, 90
)

const (
	white       = 0xffffffff
	lightYellow = 0xffff8cff
	gold        = 0xd8c480ff
	red         = 0xff0000ff
)

const (
	joinGameCharacterFilter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890._:"
)

const (
	logPrefix = "Game Screen"
)

// BuildInfo contains information about the current build
type BuildInfo struct {
	Branch, Commit string
}

// CreateMainMenu creates an instance of MainMenu
func CreateMainMenu(
	navigator d2interface.Navigator,
	asset *d2asset.AssetManager,
	renderer d2interface.Renderer,
	inputManager d2interface.InputManager,
	audioProvider d2interface.AudioProvider,
	ui *d2ui.UIManager,
	buildInfo BuildInfo,
	l d2util.LogLevel,
	errorMessageOptional ...string,
) (*MainMenu, error) {
	heroStateFactory, err := d2hero.NewHeroStateFactory(asset)
	if err != nil {
		return nil, err
	}

	mainMenu := &MainMenu{
		asset:          asset,
		screenMode:     ScreenModeUnknown,
		leftButtonHeld: true,
		renderer:       renderer,
		inputManager:   inputManager,
		audioProvider:  audioProvider,
		navigator:      navigator,
		buildInfo:      buildInfo,
		uiManager:      ui,
		heroState:      heroStateFactory,
	}

	mainMenu.Logger = d2util.NewLogger()
	mainMenu.Logger.SetPrefix(logPrefix)
	mainMenu.Logger.SetLevel(l)

	if len(errorMessageOptional) != 0 {
		mainMenu.errorLabel = ui.NewLabel(d2resource.FontFormal12, d2resource.PaletteUnits)
		mainMenu.errorLabel.SetText(errorMessageOptional[0])
	}

	return mainMenu, nil
}

// MainMenu represents the main menu
type MainMenu struct {
	tcpIPBackground     *d2ui.Sprite
	trademarkBackground *d2ui.Sprite
	background          *d2ui.Sprite
	diabloLogoLeft      *d2ui.Sprite
	diabloLogoRight     *d2ui.Sprite
	diabloLogoLeftBack  *d2ui.Sprite
	diabloLogoRightBack *d2ui.Sprite
	serverIPBackground  *d2ui.Sprite
	mainLayout          *d2ui.WidgetGroup
	networkLayout       *d2ui.WidgetGroup
	tcpIPLayout         *d2ui.WidgetGroup
	serverIPLayout      *d2ui.WidgetGroup
	trademarkLayout     *d2ui.WidgetGroup
	errorLabel          *d2ui.Label
	tcpJoinGameEntry    *d2ui.TextBox
	screenMode          mainMenuScreenMode
	leftButtonHeld      bool

	asset         *d2asset.AssetManager
	inputManager  d2interface.InputManager
	renderer      d2interface.Renderer
	audioProvider d2interface.AudioProvider
	scriptEngine  *d2script.ScriptEngine // nolint:structcheck,unused // it will be used...
	navigator     d2interface.Navigator
	uiManager     *d2ui.UIManager
	heroState     *d2hero.HeroStateFactory

	buildInfo BuildInfo

	*d2util.Logger
}

// OnLoad is called to load the resources for the main menu
func (v *MainMenu) OnLoad(loading d2screen.LoadingState) {
	v.mainLayout = v.uiManager.NewWidgetGroup(d2ui.RenderPriorityBackground)
	v.networkLayout = v.uiManager.NewWidgetGroup(d2ui.RenderPriorityBackground)
	v.tcpIPLayout = v.uiManager.NewWidgetGroup(d2ui.RenderPriorityBackground)
	v.serverIPLayout = v.uiManager.NewWidgetGroup(d2ui.RenderPriorityBackground)
	v.trademarkLayout = v.uiManager.NewWidgetGroup(d2ui.RenderPriorityBackground)

	v.audioProvider.PlayBGM(d2resource.BGMTitle)
	loading.Progress(twentyPercent)

	v.createLabels(loading)
	v.loadBackgroundSprites()
	v.createLogos(loading)
	v.createButtons(loading)

	v.tcpJoinGameEntry = v.uiManager.NewTextbox()
	v.tcpJoinGameEntry.SetPosition(joinGameDialogX, joinGameDialogY)
	v.tcpJoinGameEntry.SetFilter(joinGameCharacterFilter)
	v.serverIPLayout.AddWidget(v.tcpJoinGameEntry)
	loading.Progress(ninetyPercent)

	if v.screenMode == ScreenModeUnknown {
		v.SetScreenMode(ScreenModeTrademark)
	} else {
		v.SetScreenMode(ScreenModeMainMenu)
	}

	if err := v.inputManager.BindHandler(v); err != nil {
		v.Error("failed to add main menu as event handler")
	}
}

func (v *MainMenu) loadBackgroundSprites() {
	var err error

	v.background, err = v.uiManager.NewSprite(d2resource.GameSelectScreen, d2resource.PaletteSky)
	if err != nil {
		v.Error(err.Error())
	}

	v.background.SetPosition(backgroundX, backgroundY)

	v.trademarkBackground, err = v.uiManager.NewSprite(d2resource.TrademarkScreen, d2resource.PaletteSky)
	if err != nil {
		v.Error(err.Error())
	}

	v.trademarkBackground.SetPosition(backgroundX, backgroundY)

	v.tcpIPBackground, err = v.uiManager.NewSprite(d2resource.TCPIPBackground, d2resource.PaletteSky)
	if err != nil {
		v.Error(err.Error())
	}

	v.tcpIPBackground.SetPosition(backgroundX, backgroundY)

	v.serverIPBackground, err = v.uiManager.NewSprite(d2resource.PopUpOkCancel, d2resource.PaletteFechar)
	if err != nil {
		v.Error(err.Error())
	}

	v.serverIPBackground.SetPosition(serverIPbackgroundX, serverIPbackgroundY)
}

func (v *MainMenu) createLabels(loading d2screen.LoadingState) {
	versionLabel := v.uiManager.NewLabel(d2resource.FontFormal12, d2resource.PaletteStatic)
	versionLabel.Alignment = d2ui.HorizontalAlignRight
	versionLabel.SetText("OpenDiablo2 - " + v.buildInfo.Branch)
	versionLabel.Color[0] = d2util.Color(white)
	versionLabel.SetPosition(versionLabelX, versionLabelY)
	v.mainLayout.AddWidget(versionLabel)

	commitLabel := v.uiManager.NewLabel(d2resource.FontFormal10, d2resource.PaletteStatic)
	commitLabel.Alignment = d2ui.HorizontalAlignLeft
	commitLabel.SetText(v.buildInfo.Commit)
	commitLabel.Color[0] = d2util.Color(white)
	commitLabel.SetPosition(commitLabelX, commitLabelY)
	v.mainLayout.AddWidget(commitLabel)

	copyrightLabel := v.uiManager.NewLabel(d2resource.FontFormal12, d2resource.PaletteStatic)
	copyrightLabel.Alignment = d2ui.HorizontalAlignCenter
	copyrightLabel.SetText(v.asset.TranslateLabel(d2enum.CopyrightLabel))
	copyrightLabel.Color[0] = d2util.Color(lightBrown)
	copyrightLabel.SetPosition(copyrightX, copyrightY)
	v.trademarkLayout.AddWidget(copyrightLabel)
	loading.Progress(thirtyPercent)

	copyrightLabel2 := v.uiManager.NewLabel(d2resource.FontFormal12, d2resource.PaletteStatic)
	copyrightLabel2.Alignment = d2ui.HorizontalAlignCenter
	copyrightLabel2.SetText(v.asset.TranslateLabel(d2enum.AllRightsReservedLabel))
	copyrightLabel2.Color[0] = d2util.Color(lightBrown)
	copyrightLabel2.SetPosition(copyright2X, copyright2Y)
	v.trademarkLayout.AddWidget(copyrightLabel2)

	openDiabloLabel := v.uiManager.NewLabel(d2resource.FontFormal10, d2resource.PaletteStatic)
	openDiabloLabel.Alignment = d2ui.HorizontalAlignCenter
	openDiabloLabel.SetText("OpenDiablo2 is neither developed by, nor endorsed by Blizzard or its parent company Activision")
	openDiabloLabel.Color[0] = d2util.Color(lightYellow)
	openDiabloLabel.SetPosition(od2LabelX, od2LabelY)
	v.mainLayout.AddWidget(openDiabloLabel)
	loading.Progress(fiftyPercent)

	tcpIPOptionsLabel := v.uiManager.NewLabel(d2resource.Font42, d2resource.PaletteUnits)
	tcpIPOptionsLabel.SetPosition(tcpOptionsX, tcpOptionsY)
	tcpIPOptionsLabel.Alignment = d2ui.HorizontalAlignCenter
	tcpIPOptionsLabel.SetText(v.asset.TranslateLabel(d2enum.TCPIPOptionsLabel))
	v.tcpIPLayout.AddWidget(tcpIPOptionsLabel)

	tcpIPOptionsLabel2 := v.uiManager.NewLabel(d2resource.Font42, d2resource.PaletteUnits)
	tcpIPOptionsLabel2.SetPosition(tcpOptionsX, tcpOptionsY)
	tcpIPOptionsLabel2.Alignment = d2ui.HorizontalAlignCenter
	tcpIPOptionsLabel2.SetText(v.asset.TranslateLabel(d2enum.TCPIPOptionsLabel))
	v.serverIPLayout.AddWidget(tcpIPOptionsLabel2)

	tcpJoinGameLabel := v.uiManager.NewLabel(d2resource.Font16, d2resource.PaletteUnits)
	tcpJoinGameLabel.Alignment = d2ui.HorizontalAlignCenter
	tcpJoinGameLabel.SetText(strings.Join(d2util.SplitIntoLinesWithMaxWidth(v.asset.TranslateLabel(d2enum.TCPIPEnterHostIPLabel), 27), "\n"))
	tcpJoinGameLabel.Color[0] = d2util.Color(gold)
	tcpJoinGameLabel.SetPosition(joinGameX, joinGameY)
	v.serverIPLayout.AddWidget(tcpJoinGameLabel)

	machineIP := v.uiManager.NewLabel(d2resource.Font24, d2resource.PaletteUnits)
	machineIP.Alignment = d2ui.HorizontalAlignCenter
	machineIP.SetText(v.asset.TranslateLabel(d2enum.TCPIPYourIPLabel) + "\n" + v.getLocalIP())
	machineIP.Color[0] = d2util.Color(lightYellow)
	machineIP.SetPosition(machineIPX, machineIPY)
	v.serverIPLayout.AddWidget(machineIP)

	machineIP2 := v.uiManager.NewLabel(d2resource.Font24, d2resource.PaletteUnits)
	machineIP2.Alignment = d2ui.HorizontalAlignCenter
	machineIP2.SetText(v.asset.TranslateLabel(d2enum.TCPIPYourIPLabel) + "\n" + v.getLocalIP())
	machineIP2.Color[0] = d2util.Color(lightYellow)
	machineIP2.SetPosition(machineIPX, machineIPY)
	v.serverIPLayout.AddWidget(machineIP2)

	if v.errorLabel != nil {
		v.errorLabel.SetPosition(errorLabelX, errorLabelY)
		v.errorLabel.Alignment = d2ui.HorizontalAlignCenter
		v.errorLabel.Color[0] = d2util.Color(red)
		v.mainLayout.AddWidget(v.errorLabel)
	}
}

func (v *MainMenu) createLogos(loading d2screen.LoadingState) {
	var err error

	v.diabloLogoLeft, err = v.uiManager.NewSprite(d2resource.Diablo2LogoFireLeft, d2resource.PaletteUnits)
	if err != nil {
		v.Error(err.Error())
	}

	v.diabloLogoLeft.SetEffect(d2enum.DrawEffectModulate)
	v.diabloLogoLeft.PlayForward()
	v.diabloLogoLeft.SetPosition(diabloLogoX, diabloLogoY)
	loading.Progress(sixtyPercent)

	v.diabloLogoRight, err = v.uiManager.NewSprite(d2resource.Diablo2LogoFireRight, d2resource.PaletteUnits)
	if err != nil {
		v.Error(err.Error())
	}

	v.diabloLogoRight.SetEffect(d2enum.DrawEffectModulate)
	v.diabloLogoRight.PlayForward()
	v.diabloLogoRight.SetPosition(diabloLogoX, diabloLogoY)

	v.diabloLogoLeftBack, err = v.uiManager.NewSprite(d2resource.Diablo2LogoBlackLeft, d2resource.PaletteUnits)
	if err != nil {
		v.Error(err.Error())
	}

	v.diabloLogoLeftBack.SetPosition(diabloLogoX, diabloLogoY)

	v.diabloLogoRightBack, err = v.uiManager.NewSprite(d2resource.Diablo2LogoBlackRight, d2resource.PaletteUnits)
	if err != nil {
		v.Error(err.Error())
	}

	v.diabloLogoRightBack.SetPosition(diabloLogoX, diabloLogoY)
}

func (v *MainMenu) createButtons(loading d2screen.LoadingState) {
	exitDiabloButton := v.uiManager.NewButton(d2ui.ButtonTypeWide, v.asset.TranslateLabel(d2enum.ExitGameLabel))
	exitDiabloButton.SetPosition(exitDiabloBtnX, exitDiabloBtnY)
	exitDiabloButton.OnActivated(func() { v.onExitButtonClicked() })
	v.mainLayout.AddWidget(exitDiabloButton)

	creditsButton := v.uiManager.NewButton(d2ui.ButtonTypeShort, v.asset.TranslateLabel(d2enum.CreditsLabel))
	creditsButton.SetPosition(creditBtnX, creditBtnY)
	creditsButton.OnActivated(func() { v.onCreditsButtonClicked() })
	v.mainLayout.AddWidget(creditsButton)

	cinematicsButton := v.uiManager.NewButton(d2ui.ButtonTypeShort, v.asset.TranslateLabel(d2enum.CinematicsLabel))
	cinematicsButton.SetPosition(cineBtnX, cineBtnY)
	cinematicsButton.OnActivated(func() { v.onCinematicsButtonClicked() })
	v.mainLayout.AddWidget(cinematicsButton)
	loading.Progress(seventyPercent)

	singlePlayerButton := v.uiManager.NewButton(d2ui.ButtonTypeWide, v.asset.TranslateLabel(d2enum.SinglePlayerLabel))
	singlePlayerButton.SetPosition(singlePlayerBtnX, singlePlayerBtnY)
	singlePlayerButton.OnActivated(func() { v.onSinglePlayerClicked() })
	v.mainLayout.AddWidget(singlePlayerButton)

	githubButton := v.uiManager.NewButton(d2ui.ButtonTypeWide, "PROJECT WEBSITE")
	githubButton.SetPosition(githubBtnX, githubBtnY)
	githubButton.OnActivated(func() { v.onGithubButtonClicked() })
	v.mainLayout.AddWidget(githubButton)

	mapTestButton := v.uiManager.NewButton(d2ui.ButtonTypeWide, "MAP ENGINE TEST")
	mapTestButton.SetPosition(mapTestBtnX, mapTestBtnY)
	mapTestButton.OnActivated(func() { v.onMapTestClicked() })
	v.mainLayout.AddWidget(mapTestButton)

	btnTCPIPCancel := v.uiManager.NewButton(d2ui.ButtonTypeMedium,
		v.asset.TranslateLabel(d2enum.CancelLabel))
	btnTCPIPCancel.SetPosition(tcpBtnX, tcpBtnY)
	btnTCPIPCancel.OnActivated(func() { v.onTCPIPCancelClicked() })
	v.tcpIPLayout.AddWidget(btnTCPIPCancel)

	btnServerIPCancel := v.uiManager.NewButton(d2ui.ButtonTypeOkCancel, v.asset.TranslateLabel(d2enum.CancelLabel))
	btnServerIPCancel.SetPosition(srvCancelBtnX, srvCancelBtnY)
	btnServerIPCancel.OnActivated(func() { v.onBtnTCPIPCancelClicked() })
	v.serverIPLayout.AddWidget(btnServerIPCancel)

	btnServerIPOk := v.uiManager.NewButton(d2ui.ButtonTypeOkCancel, v.asset.TranslateString(d2enum.OKLabel))
	btnServerIPOk.SetPosition(srvOkBtnX, srvOkBtnY)
	btnServerIPOk.OnActivated(func() { v.onBtnTCPIPOkClicked() })
	v.serverIPLayout.AddWidget(btnServerIPOk)

	v.createMultiplayerMenuButtons()
	loading.Progress(eightyPercent)
}

func (v *MainMenu) createMultiplayerMenuButtons() {
	multiplayerButton := v.uiManager.NewButton(d2ui.ButtonTypeWide,
		v.asset.TranslateLabel(d2enum.OtherMultiplayerLabel))
	multiplayerButton.SetPosition(multiplayerBtnX, multiplayerBtnY)
	multiplayerButton.OnActivated(func() { v.onMultiplayerClicked() })
	v.mainLayout.AddWidget(multiplayerButton)

	networkTCPIPButton := v.uiManager.NewButton(d2ui.ButtonTypeWide, v.asset.TranslateLabel(d2enum.TCPIPGameLabel))
	networkTCPIPButton.SetPosition(tcpNetBtnX, tcpNetBtnY)
	networkTCPIPButton.OnActivated(func() { v.onNetworkTCPIPClicked() })
	v.networkLayout.AddWidget(networkTCPIPButton)

	networkCancelButton := v.uiManager.NewButton(d2ui.ButtonTypeWide,
		v.asset.TranslateLabel(d2enum.CancelLabel))
	networkCancelButton.SetPosition(networkCancelBtnX, networkCancelBtnY)
	networkCancelButton.OnActivated(func() { v.onNetworkCancelClicked() })
	v.networkLayout.AddWidget(networkCancelButton)

	btnTCPIPHostGame := v.uiManager.NewButton(d2ui.ButtonTypeWide, v.asset.TranslateLabel(d2enum.TCPIPHostGameLabel))
	btnTCPIPHostGame.SetPosition(tcpHostBtnX, tcpHostBtnY)
	btnTCPIPHostGame.OnActivated(func() { v.onTCPIPHostGameClicked() })
	v.tcpIPLayout.AddWidget(btnTCPIPHostGame)

	btnTCPIPJoinGame := v.uiManager.NewButton(d2ui.ButtonTypeWide, v.asset.TranslateLabel(d2enum.TCPIPJoinGameLabel))
	btnTCPIPJoinGame.SetPosition(tcpJoinBtnX, tcpJoinBtnY)
	btnTCPIPJoinGame.OnActivated(func() { v.onTCPIPJoinGameClicked() })
	v.tcpIPLayout.AddWidget(btnTCPIPJoinGame)
}

func (v *MainMenu) onMapTestClicked() {
	v.navigator.ToMapEngineTest(0, 1)
}

func (v *MainMenu) onSinglePlayerClicked() {
	v.SetScreenMode(ScreenModeUnknown)

	if v.heroState.HasGameStates() {
		// Go here only if existing characters are available to select
		v.navigator.ToCharacterSelect(d2clientconnectiontype.Local, v.tcpJoinGameEntry.GetText())
	} else {
		v.navigator.ToSelectHero(d2clientconnectiontype.Local, v.tcpJoinGameEntry.GetText())
	}
}

func (v *MainMenu) onGithubButtonClicked() {
	url := "https://www.github.com/OpenDiablo2/OpenDiablo2"

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		v.Error(err.Error())
	}
}

func (v *MainMenu) onExitButtonClicked() {
	os.Exit(0)
}

func (v *MainMenu) onCreditsButtonClicked() {
	v.navigator.ToCredits()
}

func (v *MainMenu) onCinematicsButtonClicked() {
	v.navigator.ToCinematics()
}

// Render renders the main menu
func (v *MainMenu) Render(screen d2interface.Surface) {
	v.renderBackgrounds(screen)
	v.renderLogos(screen)
}

func (v *MainMenu) renderBackgrounds(screen d2interface.Surface) {
	switch v.screenMode {
	case ScreenModeTrademark:
		v.trademarkBackground.RenderSegmented(screen, 4, 3, 0)
	case ScreenModeServerIP:
		v.tcpIPBackground.RenderSegmented(screen, 4, 3, 0)
		v.serverIPBackground.RenderSegmented(screen, 2, 1, 0)
	case ScreenModeTCPIP:
		v.tcpIPBackground.RenderSegmented(screen, 4, 3, 0)
	default:
		v.background.RenderSegmented(screen, 4, 3, 0)
	}
}

func (v *MainMenu) renderLogos(screen d2interface.Surface) {
	switch v.screenMode {
	case ScreenModeTrademark, ScreenModeMainMenu, ScreenModeMultiplayer:
		v.diabloLogoLeftBack.Render(screen)
		v.diabloLogoRightBack.Render(screen)
		v.diabloLogoLeft.Render(screen)
		v.diabloLogoRight.Render(screen)
	}
}

// Advance runs the update logic on the main menu
func (v *MainMenu) Advance(tickTime float64) error {
	switch v.screenMode {
	case ScreenModeMainMenu, ScreenModeTrademark, ScreenModeMultiplayer:
		if err := v.diabloLogoLeftBack.Advance(tickTime); err != nil {
			return err
		}

		if err := v.diabloLogoRightBack.Advance(tickTime); err != nil {
			return err
		}

		if err := v.diabloLogoLeft.Advance(tickTime); err != nil {
			return err
		}

		if err := v.diabloLogoRight.Advance(tickTime); err != nil {
			return err
		}
	}

	return nil
}

// OnMouseButtonDown is called when a mouse button is clicked
func (v *MainMenu) OnMouseButtonDown(event d2interface.MouseEvent) bool {
	if v.screenMode == ScreenModeTrademark && event.Button() == d2enum.MouseButtonLeft {
		v.SetScreenMode(ScreenModeMainMenu)
		return true
	}

	return false
}

func (v *MainMenu) onEscapePressed(event d2interface.KeyEvent, mode mainMenuScreenMode) {
	if event.Key() == d2enum.KeyEscape {
		v.SetScreenMode(mode)
	}
}

// OnKeyUp is called when a key is released
func (v *MainMenu) OnKeyUp(event d2interface.KeyEvent) bool {
	preventKeyEventPropagation := false

	switch v.screenMode {
	case ScreenModeTrademark: // On retail version of D2, some specific key events (Escape, Space and Enter) puts you onto the main menu.
		switch event.Key() {
		case d2enum.KeyEscape, d2enum.KeyEnter, d2enum.KeySpace:
			v.SetScreenMode(ScreenModeMainMenu)
		}

		preventKeyEventPropagation = true
	case ScreenModeMainMenu: // pressing escape in Main Menu close the game
		if event.Key() == d2enum.KeyEscape {
			v.onExitButtonClicked()
		}
	case ScreenModeMultiplayer: // back to previous menu
		v.onEscapePressed(event, ScreenModeMainMenu)

		preventKeyEventPropagation = true
	case ScreenModeTCPIP: // back to previous menu
		v.onEscapePressed(event, ScreenModeMultiplayer)

		preventKeyEventPropagation = true
	case ScreenModeServerIP: // back to previous menu
		v.onEscapePressed(event, ScreenModeTCPIP)

		preventKeyEventPropagation = true
	}

	return preventKeyEventPropagation
}

// SetScreenMode sets the screen mode (which sub-menu the screen is on)
func (v *MainMenu) SetScreenMode(screenMode mainMenuScreenMode) {
	v.screenMode = screenMode
	isTrademark := screenMode == ScreenModeTrademark
	isMainMenu := screenMode == ScreenModeMainMenu
	isMultiplayer := screenMode == ScreenModeMultiplayer
	isTCPIP := screenMode == ScreenModeTCPIP
	isServerIP := screenMode == ScreenModeServerIP

	v.trademarkLayout.SetVisible(isTrademark)
	v.mainLayout.SetVisible(isMainMenu)
	v.networkLayout.SetVisible(isMultiplayer)
	v.tcpIPLayout.SetVisible(isTCPIP)

	if isServerIP {
		v.tcpJoinGameEntry.Activate()
	}

	v.serverIPLayout.SetVisible(isServerIP)
}

func (v *MainMenu) onNetworkCancelClicked() {
	v.SetScreenMode(ScreenModeMainMenu)
}

func (v *MainMenu) onMultiplayerClicked() {
	v.SetScreenMode(ScreenModeMultiplayer)
}

func (v *MainMenu) onNetworkTCPIPClicked() {
	v.SetScreenMode(ScreenModeTCPIP)
}

func (v *MainMenu) onTCPIPCancelClicked() {
	v.SetScreenMode(ScreenModeMultiplayer)
}

func (v *MainMenu) onTCPIPHostGameClicked() {
	v.navigator.ToCharacterSelect(d2clientconnectiontype.LANServer, "")
}

func (v *MainMenu) onTCPIPJoinGameClicked() {
	v.SetScreenMode(ScreenModeServerIP)
}

func (v *MainMenu) onBtnTCPIPCancelClicked() {
	v.SetScreenMode(ScreenModeTCPIP)
}

func (v *MainMenu) onBtnTCPIPOkClicked() {
	v.navigator.ToCharacterSelect(d2clientconnectiontype.LANClient, v.tcpJoinGameEntry.GetText())
}

// getLocalIP returns local machine IP address
func (v *MainMenu) getLocalIP() string {
	// https://stackoverflow.com/a/28862477
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}

	v.Warning("no IPv4 Address could be found")

	return v.asset.TranslateLabel(d2enum.IPNotFoundLabel)
}
