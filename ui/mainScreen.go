package ui

import (
	"bufio"
	"os"
	"unicode/utf8"

	"github.com/veandco/go-sdl2/sdl"
)

var mainMenuBackground *sdl.Texture
var uiAtlas *sdl.Texture

type pos struct {
	x, y int
}

type button struct {
	pos
	srcRect   []*sdl.Rect
	dstRect   []*sdl.Rect
	str       *sdl.Texture
	isClicked bool
}

type mainMenuObj struct {
	play       button
	info       button
	infoTab    *sdl.Texture
	infoStr    []*sdl.Texture
	infoStrLen []int32
}

func createMainMenu(ui *UI2d) {
	mainMenuBackground = imgFileToTexture("ui/assets/main_menu_background.png")
	uiAtlas = imgFileToTexture("ui/assets/ui_split.png")
	ui.mainMenu = mainMenuObj{}

	ui.mainMenu.play = button{pos: pos{x: winWidth * 0.4, y: winHeight * .4}, isClicked: false}
	ui.mainMenu.play.srcRect = append(ui.mainMenu.play.srcRect, &sdl.Rect{310, 349, 25, 32})
	ui.mainMenu.play.srcRect = append(ui.mainMenu.play.srcRect, &sdl.Rect{313, 381, 70, 25})
	ui.mainMenu.play.dstRect = append(ui.mainMenu.play.dstRect, &sdl.Rect{winWidth * .4, winHeight * .65, 25 * 2, 32 * 2})
	ui.mainMenu.play.dstRect = append(ui.mainMenu.play.dstRect, &sdl.Rect{(winWidth * .4) + 24*2, (winHeight * .65) + 3*2, 70 * 2, 25 * 2})
	ui.mainMenu.play.str = getTextTexture("Play", sdl.Color{255, 255, 255, 0})
	ui.mainMenu.play.dstRect = append(ui.mainMenu.play.dstRect, &sdl.Rect{(winWidth * .4) + 28*2, (winHeight * .65) + 5*2, 45 * 2, 20 * 2})

	ui.mainMenu.info = button{pos: pos{x: winWidth * 0.4, y: winHeight * .5}, isClicked: false}
	ui.mainMenu.info.srcRect = append(ui.mainMenu.info.srcRect, &sdl.Rect{336, 349, 25, 32})
	ui.mainMenu.info.srcRect = append(ui.mainMenu.info.srcRect, &sdl.Rect{313, 381, 70, 25})
	ui.mainMenu.info.dstRect = append(ui.mainMenu.info.dstRect, &sdl.Rect{winWidth * .4, winHeight * .8, 25 * 2, 32 * 2})
	ui.mainMenu.info.dstRect = append(ui.mainMenu.info.dstRect, &sdl.Rect{(winWidth * .4) + 24*2, (winHeight * .8) + 3*2, 70 * 2, 25 * 2})
	ui.mainMenu.info.str = getTextTexture("Info", sdl.Color{255, 255, 255, 0})
	ui.mainMenu.info.dstRect = append(ui.mainMenu.info.dstRect, &sdl.Rect{(winWidth * .4) + 28*2, (winHeight * .8) + 5*2, 45 * 2, 20 * 2})

	ui.mainMenu.infoTab = createOnePixel(0, 0, 0, 180)
	ui.mainMenu.infoTab.SetBlendMode(sdl.BLENDMODE_BLEND)
	ui.mainMenu.infoStr = append(ui.mainMenu.infoStr, getTextTexture("Artificial Intelligence Final Project", sdl.Color{255, 255, 255, 255}))
	ui.mainMenu.infoStrLen = append(ui.mainMenu.infoStrLen, 34)
	//ui.mainMenu.infoStr[0].SetBlendMode(sdl.BLENDMODE_BLEND)

	file, err := os.Open("ui/assets/mainMenuInfo.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ui.mainMenu.infoStrLen = append(ui.mainMenu.infoStrLen, int32(utf8.RuneCountInString(scanner.Text())))
		if ui.mainMenu.infoStrLen[len(ui.mainMenu.infoStrLen)-1] != 0 {
			ui.mainMenu.infoStr = append(ui.mainMenu.infoStr, getTextTexture(scanner.Text(), sdl.Color{255, 255, 255, 255}))
		} else {
			ui.mainMenu.infoStr = append(ui.mainMenu.infoStr, nil)
		}
	}
}

func drawMenuButtons(ui *UI2d) {
	for i := 0; i < 2; i++ {
		renderer.Copy(uiAtlas, ui.mainMenu.play.srcRect[i], ui.mainMenu.play.dstRect[i])
		renderer.Copy(uiAtlas, ui.mainMenu.info.srcRect[i], ui.mainMenu.info.dstRect[i])
	}

	renderer.Copy(ui.mainMenu.play.str, nil, ui.mainMenu.play.dstRect[2])
	renderer.Copy(ui.mainMenu.info.str, nil, ui.mainMenu.info.dstRect[2])
}

func updateMenu(ui *UI2d) {
	if !ui.mainMenu.play.isClicked && !ui.mainMenu.info.isClicked {
		drawMenuButtons(ui)
	} else if ui.mainMenu.play.isClicked {
		currentState = selectScreen
		ui.mainMenu.play.isClicked = false
	} else if ui.mainMenu.info.isClicked {
		renderer.Copy(ui.mainMenu.infoTab, nil, &sdl.Rect{winWidth * .1, winHeight * .1, winWidth * .8, winHeight * .8})
		for i, infoStr := range ui.mainMenu.infoStr {
			if ui.mainMenu.infoStrLen[i] != 0 {
				renderer.Copy(infoStr, nil, &sdl.Rect{winWidth*.1 + 10, int32(winHeight*.1 + 10 + i*20), ui.mainMenu.infoStrLen[i] * 10, 20})
			}
		}
		if ui.input.currKeyState[sdl.SCANCODE_ESCAPE] != 0 && ui.input.prevKeyState[sdl.SCANCODE_ESCAPE] == 0 {
			ui.mainMenu.info.isClicked = false
		}
	}
}

func mainMenu(ui *UI2d) stateFunc {
	renderer.Copy(mainMenuBackground, nil, nil)
	updateMenu(ui)
	ui.input.updateMouseState()
	if ui.input.leftButton && !ui.input.prevLeftButton {
		for i := 0; i < 2; i++ {
			if ui.mainMenu.play.dstRect[i].HasIntersection(&sdl.Rect{int32(ui.input.x), int32(ui.input.y), 1, 1}) {
				ui.mainMenu.play.isClicked = true
			}
			if ui.mainMenu.info.dstRect[i].HasIntersection(&sdl.Rect{int32(ui.input.x), int32(ui.input.y), 1, 1}) {
				ui.mainMenu.info.isClicked = true
			}
		}
	}
	return determineToken
}
