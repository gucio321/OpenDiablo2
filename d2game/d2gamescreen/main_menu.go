// Package d2gamescreen contains the screens
package d2gamescreen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2hero"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2resource"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2gui"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2screen"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2ui"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2client/d2clientconnectiontype"
	"github.com/OpenDiablo2/OpenDiablo2/d2script"
)

type mainMenuScreenMode int

/*const (
	Trademark                   = 0
	MainMenu                    = 1
	SinglePlayerCharacterSelect = 2
	SinglePlayerSelectHeroClass = 3
	MultiplayerMenu             = 4
	MultiplayerTCPIP            = 5
	MultiplayerServerIP         = 6
	MapEngineTest               = 7
	Credits                     = 8
	CinematicsMenu              = 9
)*/

const (
	Trademark                   = 0
	SubMenuMainMenu             = 1
	SinglePlayerCharacterSelect = 2
	SinglePlayerSelectHeroClass = 3
	MultiplayerMenu             = 4
	MultiplayerTCPIP            = 5
	MultiplayerServerIP         = 6
	SubMenuMapEngineTest        = 7
	SubMenuCredits              = 8
	CinematicsMenu              = 9
)

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
	tcpHostBtnX, tcpHostBtnY                 = 264, 280
	tcpJoinBtnX, tcpJoinBtnY                 = 264, 320
	errorLabelX, errorLabelY                 = 400, 250
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

// BuildInfo contains information about the current build
type BuildInfo struct {
	Branch, Commit string
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
	singlePlayerButton  *d2ui.Button
	multiplayerButton   *d2ui.Button
	githubButton        *d2ui.Button
	exitDiabloButton    *d2ui.Button
	creditsButton       *d2ui.Button
	cinematicsButton    *d2ui.Button
	mapTestButton       *d2ui.Button
	networkTCPIPButton  *d2ui.Button
	networkCancelButton *d2ui.Button
	btnTCPIPCancel      *d2ui.Button
	btnTCPIPHostGame    *d2ui.Button
	btnTCPIPJoinGame    *d2ui.Button
	btnServerIPCancel   *d2ui.Button
	btnServerIPOk       *d2ui.Button
	copyrightLabel      *d2ui.Label
	copyrightLabel2     *d2ui.Label
	openDiabloLabel     *d2ui.Label
	versionLabel        *d2ui.Label
	commitLabel         *d2ui.Label
	tcpIPOptionsLabel   *d2ui.Label
	tcpJoinGameLabel    *d2ui.Label
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

	SubMenu int
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

	if len(errorMessageOptional) != 0 {
		mainMenu.errorLabel = ui.NewLabel(d2resource.FontFormal12, d2resource.PaletteUnits)
		mainMenu.errorLabel.SetText(errorMessageOptional[0])
	}

	return mainMenu, nil
}

// OnLoad is called to load the resources for the main menu
func (v *MainMenu) OnLoad(loading d2screen.LoadingState) {
	v.audioProvider.PlayBGM(d2resource.BGMTitle)
	loading.Progress(twentyPercent)

	v.createLabels(loading)
	v.loadBackgroundSprites()
	v.createLogos(loading)
	v.createButtons(loading)

	v.tcpJoinGameEntry = v.uiManager.NewTextbox()
	v.tcpJoinGameEntry.SetPosition(joinGameDialogX, joinGameDialogY)
	v.tcpJoinGameEntry.SetFilter(joinGameCharacterFilter)
	loading.Progress(ninetyPercent)

	if v.screenMode == ScreenModeUnknown {
		v.SetScreenMode(ScreenModeTrademark)
	} else {
		v.SetScreenMode(ScreenModeMainMenu)
	}

	if err := v.inputManager.BindHandler(v); err != nil {
		fmt.Println("failed to add main menu as event handler")
	}
}

func (v *MainMenu) loadBackgroundSprites() {
	var err error

	v.background, err = v.uiManager.NewSprite(d2resource.GameSelectScreen, d2resource.PaletteSky)
	if err != nil {
		log.Print(err)
	}

	v.background.SetPosition(backgroundX, backgroundY)

	v.trademarkBackground, err = v.uiManager.NewSprite(d2resource.TrademarkScreen, d2resource.PaletteSky)
	if err != nil {
		log.Print(err)
	}

	v.trademarkBackground.SetPosition(backgroundX, backgroundY)

	v.tcpIPBackground, err = v.uiManager.NewSprite(d2resource.TCPIPBackground, d2resource.PaletteSky)
	if err != nil {
		log.Print(err)
	}

	v.tcpIPBackground.SetPosition(backgroundX, backgroundY)

	v.serverIPBackground, err = v.uiManager.NewSprite(d2resource.PopUpOkCancel, d2resource.PaletteFechar)
	if err != nil {
		log.Print(err)
	}

	v.serverIPBackground.SetPosition(serverIPbackgroundX, serverIPbackgroundY)
}

func (v *MainMenu) createLabels(loading d2screen.LoadingState) {
	v.versionLabel = v.uiManager.NewLabel(d2resource.FontFormal12, d2resource.PaletteStatic)
	v.versionLabel.Alignment = d2gui.HorizontalAlignRight
	v.versionLabel.SetText("OpenDiablo2 - " + v.buildInfo.Branch)
	v.versionLabel.Color[0] = rgbaColor(white)
	v.versionLabel.SetPosition(versionLabelX, versionLabelY)

	v.commitLabel = v.uiManager.NewLabel(d2resource.FontFormal10, d2resource.PaletteStatic)
	v.commitLabel.Alignment = d2gui.HorizontalAlignLeft
	v.commitLabel.SetText(v.buildInfo.Commit)
	v.commitLabel.Color[0] = rgbaColor(white)
	v.commitLabel.SetPosition(commitLabelX, commitLabelY)

	v.copyrightLabel = v.uiManager.NewLabel(d2resource.FontFormal12, d2resource.PaletteStatic)
	v.copyrightLabel.Alignment = d2gui.HorizontalAlignCenter
	v.copyrightLabel.SetText("Diablo 2 is © Copyright 2000-2016 Blizzard Entertainment")
	v.copyrightLabel.Color[0] = rgbaColor(lightBrown)
	v.copyrightLabel.SetPosition(copyrightX, copyrightY)
	loading.Progress(thirtyPercent)

	v.copyrightLabel2 = v.uiManager.NewLabel(d2resource.FontFormal12, d2resource.PaletteStatic)
	v.copyrightLabel2.Alignment = d2gui.HorizontalAlignCenter
	v.copyrightLabel2.SetText("All Rights Reserved.")
	v.copyrightLabel2.Color[0] = rgbaColor(lightBrown)
	v.copyrightLabel2.SetPosition(copyright2X, copyright2Y)

	v.openDiabloLabel = v.uiManager.NewLabel(d2resource.FontFormal10, d2resource.PaletteStatic)
	v.openDiabloLabel.Alignment = d2gui.HorizontalAlignCenter
	v.openDiabloLabel.SetText("OpenDiablo2 is neither developed by, nor endorsed by Blizzard or its parent company Activision")
	v.openDiabloLabel.Color[0] = rgbaColor(lightYellow)
	v.openDiabloLabel.SetPosition(od2LabelX, od2LabelY)
	loading.Progress(fiftyPercent)

	v.tcpIPOptionsLabel = v.uiManager.NewLabel(d2resource.Font42, d2resource.PaletteUnits)
	v.tcpIPOptionsLabel.SetPosition(tcpOptionsX, tcpOptionsY)
	v.tcpIPOptionsLabel.Alignment = d2gui.HorizontalAlignCenter
	v.tcpIPOptionsLabel.SetText("TCP/IP Options")

	v.tcpJoinGameLabel = v.uiManager.NewLabel(d2resource.Font16, d2resource.PaletteUnits)
	v.tcpJoinGameLabel.Alignment = d2gui.HorizontalAlignCenter

	v.tcpJoinGameLabel.SetText("Enter Host IP Address\nto Join Game")

	v.tcpJoinGameLabel.Color[0] = rgbaColor(gold)

	v.tcpJoinGameLabel.SetPosition(joinGameX, joinGameY)

	if v.errorLabel != nil {
		v.errorLabel.SetPosition(errorLabelX, errorLabelY)
		v.errorLabel.Alignment = d2gui.HorizontalAlignCenter
		v.errorLabel.Color[0] = rgbaColor(red)
	}
}

func (v *MainMenu) createLogos(loading d2screen.LoadingState) {
	var err error

	v.diabloLogoLeft, err = v.uiManager.NewSprite(d2resource.Diablo2LogoFireLeft, d2resource.PaletteUnits)
	if err != nil {
		log.Print(err)
	}

	v.diabloLogoLeft.SetEffect(d2enum.DrawEffectModulate)
	v.diabloLogoLeft.PlayForward()
	v.diabloLogoLeft.SetPosition(diabloLogoX, diabloLogoY)
	loading.Progress(sixtyPercent)

	v.diabloLogoRight, err = v.uiManager.NewSprite(d2resource.Diablo2LogoFireRight, d2resource.PaletteUnits)
	if err != nil {
		log.Print(err)
	}

	v.diabloLogoRight.SetEffect(d2enum.DrawEffectModulate)
	v.diabloLogoRight.PlayForward()
	v.diabloLogoRight.SetPosition(diabloLogoX, diabloLogoY)

	v.diabloLogoLeftBack, err = v.uiManager.NewSprite(d2resource.Diablo2LogoBlackLeft, d2resource.PaletteUnits)
	if err != nil {
		log.Print(err)
	}

	v.diabloLogoLeftBack.SetPosition(diabloLogoX, diabloLogoY)

	v.diabloLogoRightBack, err = v.uiManager.NewSprite(d2resource.Diablo2LogoBlackRight, d2resource.PaletteUnits)
	if err != nil {
		log.Print(err)
	}

	v.diabloLogoRightBack.SetPosition(diabloLogoX, diabloLogoY)
}

func (v *MainMenu) createButtons(loading d2screen.LoadingState) {
	v.exitDiabloButton = v.uiManager.NewButton(d2ui.ButtonTypeWide, "EXIT DIABLO II")
	v.exitDiabloButton.SetPosition(exitDiabloBtnX, exitDiabloBtnY)
	v.exitDiabloButton.OnActivated(func() { v.onExitButtonClicked() })

	v.creditsButton = v.uiManager.NewButton(d2ui.ButtonTypeShort, "CREDITS")
	v.creditsButton.SetPosition(creditBtnX, creditBtnY)
	v.creditsButton.OnActivated(func() { v.onCreditsButtonClicked() })

	v.cinematicsButton = v.uiManager.NewButton(d2ui.ButtonTypeShort, "CINEMATICS")
	v.cinematicsButton.SetPosition(cineBtnX, cineBtnY)
	v.cinematicsButton.OnActivated(func() { v.onCinematicsButtonClicked() })
	loading.Progress(seventyPercent)

	v.singlePlayerButton = v.uiManager.NewButton(d2ui.ButtonTypeWide, "SINGLE PLAYER")
	v.singlePlayerButton.SetPosition(singlePlayerBtnX, singlePlayerBtnY)
	v.singlePlayerButton.OnActivated(func() { v.onSinglePlayerClicked() })

	v.githubButton = v.uiManager.NewButton(d2ui.ButtonTypeWide, "PROJECT WEBSITE")
	v.githubButton.SetPosition(githubBtnX, githubBtnY)
	v.githubButton.OnActivated(func() { v.onGithubButtonClicked() })

	v.mapTestButton = v.uiManager.NewButton(d2ui.ButtonTypeWide, "MAP ENGINE TEST")
	v.mapTestButton.SetPosition(mapTestBtnX, mapTestBtnY)
	v.mapTestButton.OnActivated(func() { v.onMapTestClicked() })

	v.btnTCPIPCancel = v.uiManager.NewButton(d2ui.ButtonTypeMedium,
		v.asset.TranslateString("cancel"))
	v.btnTCPIPCancel.SetPosition(tcpBtnX, tcpBtnY)
	v.btnTCPIPCancel.OnActivated(func() { v.onTCPIPCancelClicked() })

	v.btnServerIPCancel = v.uiManager.NewButton(d2ui.ButtonTypeOkCancel, "CANCEL")
	v.btnServerIPCancel.SetPosition(srvCancelBtnX, srvCancelBtnY)
	v.btnServerIPCancel.OnActivated(func() { v.onBtnTCPIPCancelClicked() })

	v.btnServerIPOk = v.uiManager.NewButton(d2ui.ButtonTypeOkCancel, "OK")
	v.btnServerIPOk.SetPosition(srvOkBtnX, srvOkBtnY)
	v.btnServerIPOk.OnActivated(func() { v.onBtnTCPIPOkClicked() })

	v.createMultiplayerMenuButtons()
	loading.Progress(eightyPercent)
}

func (v *MainMenu) createMultiplayerMenuButtons() {
	v.multiplayerButton = v.uiManager.NewButton(d2ui.ButtonTypeWide, "MULTIPLAYER")
	v.multiplayerButton.SetPosition(multiplayerBtnX, multiplayerBtnY)
	v.multiplayerButton.OnActivated(func() { v.onMultiplayerClicked() })

	v.networkTCPIPButton = v.uiManager.NewButton(d2ui.ButtonTypeWide, "TCP/IP GAME")
	v.networkTCPIPButton.SetPosition(tcpNetBtnX, tcpNetBtnY)
	v.networkTCPIPButton.OnActivated(func() { v.onNetworkTCPIPClicked() })

	v.networkCancelButton = v.uiManager.NewButton(d2ui.ButtonTypeWide,
		v.asset.TranslateString("cancel"))
	v.networkCancelButton.SetPosition(networkCancelBtnX, networkCancelBtnY)
	v.networkCancelButton.OnActivated(func() { v.onNetworkCancelClicked() })

	v.btnTCPIPHostGame = v.uiManager.NewButton(d2ui.ButtonTypeWide, "HOST GAME")
	v.btnTCPIPHostGame.SetPosition(tcpHostBtnX, tcpHostBtnY)
	v.btnTCPIPHostGame.OnActivated(func() { v.onTCPIPHostGameClicked() })

	v.btnTCPIPJoinGame = v.uiManager.NewButton(d2ui.ButtonTypeWide, "JOIN GAME")
	v.btnTCPIPJoinGame.SetPosition(tcpJoinBtnX, tcpJoinBtnY)
	v.btnTCPIPJoinGame.OnActivated(func() { v.onTCPIPJoinGameClicked() })
}

func (v *MainMenu) onMapTestClicked() {
	v.SetSubMenuMode(SubMenuMapEngineTest)
	v.navigator.ToMapEngineTest(0, 1)
}

func (v *MainMenu) onSinglePlayerClicked() {
	if v.heroState.HasGameStates() {
		// Go here only if existing characters are available to select
		v.SetSubMenuMode(SinglePlayerCharacterSelect)
		v.navigator.ToCharacterSelect(d2clientconnectiontype.Local, v.tcpJoinGameEntry.GetText())
	} else {
		v.SetSubMenuMode(SinglePlayerSelectHeroClass)
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
		log.Fatal(err)
	}
}

func (v *MainMenu) onExitButtonClicked() {
	os.Exit(0)
}

func (v *MainMenu) onCreditsButtonClicked() {
	v.SetSubMenuMode(SubMenuCredits)
	v.navigator.ToCredits()
}

func (v *MainMenu) onCinematicsButtonClicked() {
	v.SetSubMenuMode(CinematicsMenu)
	v.navigator.ToCinematics()
}

// Render renders the main menu
func (v *MainMenu) Render(screen d2interface.Surface) {
	v.renderBackgrounds(screen)
	v.renderLogos(screen)
	v.renderLabels(screen)
}

func (v *MainMenu) renderBackgrounds(screen d2interface.Surface) {
	switch v.screenMode {
	case ScreenModeTrademark:
		if err := v.trademarkBackground.RenderSegmented(screen, 4, 3, 0); err != nil {
			return
		}
	case ScreenModeServerIP:
		if err := v.tcpIPBackground.RenderSegmented(screen, 4, 3, 0); err != nil {
			return
		}

		if err := v.serverIPBackground.RenderSegmented(screen, 2, 1, 0); err != nil {
			return
		}
	case ScreenModeTCPIP:
		if err := v.tcpIPBackground.RenderSegmented(screen, 4, 3, 0); err != nil {
			return
		}
	default:
		if err := v.background.RenderSegmented(screen, 4, 3, 0); err != nil {
			return
		}
	}
}

func (v *MainMenu) renderLogos(screen d2interface.Surface) {
	switch v.screenMode {
	case ScreenModeTrademark, ScreenModeMainMenu, ScreenModeMultiplayer:
		v.diabloLogoLeftBack.RenderNoError(screen)
		v.diabloLogoRightBack.RenderNoError(screen)
		v.diabloLogoLeft.RenderNoError(screen)
		v.diabloLogoRight.RenderNoError(screen)
	}
}

func (v *MainMenu) renderLabels(screen d2interface.Surface) {
	switch v.screenMode {
	case ScreenModeServerIP:
		v.tcpIPOptionsLabel.RenderNoError(screen)
		v.tcpJoinGameLabel.RenderNoError(screen)
	case ScreenModeTCPIP:
		v.tcpIPOptionsLabel.RenderNoError(screen)
	case ScreenModeTrademark:
		v.copyrightLabel.RenderNoError(screen)
		v.copyrightLabel2.RenderNoError(screen)

		if v.errorLabel != nil {
			v.errorLabel.RenderNoError(screen)
		}
	case ScreenModeMainMenu:
		v.openDiabloLabel.RenderNoError(screen)
		v.versionLabel.RenderNoError(screen)
		v.commitLabel.RenderNoError(screen)
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
		v.SetSubMenuMode(SubMenuMainMenu)
		v.SetScreenMode(ScreenModeMainMenu)
		return true
	}

	return false
}

func (v *MainMenu) onEscapePressed(event d2interface.KeyEvent) bool {
	if event.Key() == d2enum.KeyEscape {
		return true
	}
	return false
}

// OnKeyUp is called when a key is released
func (v *MainMenu) OnKeyUp(event d2interface.KeyEvent) bool {
	fmt.Println("\n\nOn Key Up called\ninput submenu value is: ", v.SubMenu)
	preventKeyEventPropagation := false

	/*	switch v.SubMenu {
		case Trademark: // escape in trademar, does nothing
			switch event.Key() {
			case d2enum.KeyEscape, d2enum.KeyEnter, d2enum.KeySpace:
				v.SetSubMenuMode(SubMenuMainMenu)
				v.SetScreenMode(ScreenModeMainMenu)
			}
			preventKeyEventPropagation = true
		case SubMenuMainMenu: // for "main" main menu screen pressing escape closes game
			if v.onEscapePressed(event) {
				v.onExitButtonClicked()
			}
		case MultiplayerMenu, SubMenuCredits, CinematicsMenu, SinglePlayerCharacterSelect, SubMenuMapEngineTest: // escape puts from credits, cinematics select menu or single player character select menu onto main menu (trademark)
			if v.onEscapePressed(event) {
				v.navigator.ToMainMenu()
				v.SetSubMenuMode(Trademark)
			}
		case SinglePlayerSelectHeroClass: // escape in select hero class menu puts onto character select menu
			if v.onEscapePressed(event) {
				v.navigator.ToCharacterSelect(d2clientconnectiontype.Local, v.tcpJoinGameEntry.GetText())
				v.SetSubMenuMode(Trademark)
			}
			preventKeyEventPropagation = true
		case MultiplayerTCPIP: // back to previous menu
			if v.onEscapePressed(event) {
				v.onTCPIPCancelClicked()
			}
			preventKeyEventPropagation = true
		case MultiplayerServerIP:
			if v.onEscapePressed(event) {
				v.onBtnTCPIPCancelClicked()
			}
			preventKeyEventPropagation = true
		}*/

	switch v.SubMenu {
	case Trademark: // escape in trademar, does nothing
		switch event.Key() {
		case d2enum.KeyEscape, d2enum.KeyEnter, d2enum.KeySpace:
			v.SetSubMenuMode(SubMenuMainMenu)
			v.SetScreenMode(ScreenModeMainMenu)
		}
		preventKeyEventPropagation = true
	case SubMenuMainMenu: // for "main" main menu screen pressing escape closes game
		if v.onEscapePressed(event) {
			fmt.Println("\n\nEscape button pressed in main menu\nSubMenu is: ", v.SubMenu, " = 1")
			v.onExitButtonClicked()
		}
	case MultiplayerMenu, SubMenuCredits, CinematicsMenu, SinglePlayerCharacterSelect, SubMenuMapEngineTest: // escape puts from credits, cinematics select menu or single player character select menu onto main menu (trademark)
		if v.onEscapePressed(event) {
			v.navigator.ToMainMenu()
			v.SetSubMenuMode(Trademark)
		}
	case SinglePlayerSelectHeroClass: // escape in select hero class menu puts onto character select menu
		if v.onEscapePressed(event) {
			v.navigator.ToCharacterSelect(d2clientconnectiontype.Local, v.tcpJoinGameEntry.GetText())
			v.SetSubMenuMode(Trademark)
		}
		preventKeyEventPropagation = true
	case MultiplayerTCPIP: // back to previous menu
		if v.onEscapePressed(event) {
			v.onTCPIPCancelClicked()
		}
		preventKeyEventPropagation = true
	case MultiplayerServerIP:
		if v.onEscapePressed(event) {
			v.onBtnTCPIPCancelClicked()
		}
		preventKeyEventPropagation = true
	}

	return preventKeyEventPropagation
}

// SetScreenMode sets the screen mode (which sub-menu the screen is on)
func (v *MainMenu) SetScreenMode(screenMode mainMenuScreenMode) {
	v.screenMode = screenMode
	isMainMenu := screenMode == ScreenModeMainMenu
	isMultiplayer := screenMode == ScreenModeMultiplayer
	isTCPIP := screenMode == ScreenModeTCPIP
	isServerIP := screenMode == ScreenModeServerIP

	v.exitDiabloButton.SetVisible(isMainMenu)
	v.creditsButton.SetVisible(isMainMenu)
	v.cinematicsButton.SetVisible(isMainMenu)
	v.singlePlayerButton.SetVisible(isMainMenu)
	v.githubButton.SetVisible(isMainMenu)
	v.mapTestButton.SetVisible(isMainMenu)
	v.multiplayerButton.SetVisible(isMainMenu)
	v.networkTCPIPButton.SetVisible(isMultiplayer)
	v.networkCancelButton.SetVisible(isMultiplayer)
	v.btnTCPIPCancel.SetVisible(isTCPIP)
	v.btnTCPIPHostGame.SetVisible(isTCPIP)
	v.btnTCPIPJoinGame.SetVisible(isTCPIP)
	v.tcpJoinGameEntry.SetVisible(isServerIP)

	if isServerIP {
		v.tcpJoinGameEntry.Activate()
	}

	v.btnServerIPOk.SetVisible(isServerIP)
	v.btnServerIPCancel.SetVisible(isServerIP)
}

func (v *MainMenu) onNetworkCancelClicked() {
	v.SetSubMenuMode(Trademark)
	v.SetScreenMode(ScreenModeMainMenu)
}

func (v *MainMenu) onMultiplayerClicked() {
	v.SetSubMenuMode(MultiplayerMenu)
	v.SetScreenMode(ScreenModeMultiplayer)
}

func (v *MainMenu) onNetworkTCPIPClicked() {
	v.SetSubMenuMode(MultiplayerTCPIP)
	v.SetScreenMode(ScreenModeTCPIP)
}

func (v *MainMenu) onTCPIPCancelClicked() {
	v.SetSubMenuMode(MultiplayerMenu)
	v.SetScreenMode(ScreenModeMultiplayer)
}

func (v *MainMenu) onTCPIPHostGameClicked() {
	v.SetSubMenuMode(SinglePlayerCharacterSelect)
	v.navigator.ToCharacterSelect(d2clientconnectiontype.LANServer, "")
}

func (v *MainMenu) onTCPIPJoinGameClicked() {
	v.SetSubMenuMode(MultiplayerServerIP)
	v.SetScreenMode(ScreenModeServerIP)
}

func (v *MainMenu) onBtnTCPIPCancelClicked() {
	v.SetSubMenuMode(MultiplayerTCPIP)
	v.SetScreenMode(ScreenModeTCPIP)
}

func (v *MainMenu) onBtnTCPIPOkClicked() {
	v.SetSubMenuMode(SinglePlayerCharacterSelect)
	v.navigator.ToCharacterSelect(d2clientconnectiontype.LANClient, v.tcpJoinGameEntry.GetText())
}

func (v *MainMenu) SetSubMenuMode(menu int) {
	v.SubMenu = menu
	fmt.Printf("\n\nSetSubmenuMode was called\nwith input value:", menu)
}
