package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type endMenuObj struct {
	gameEnd          *sdl.Texture
	isTerminated     bool
	isRestarted      bool
	returnToMainMenu button
	exit             button
}

func endMenuInit(ui *UI2d) {
	ui.endMenu = endMenuObj{}
	ui.endMenu.isTerminated = false
	ui.endMenu.isRestarted = false

	ui.endMenu.gameEnd = getTextTexture("Game Over", sdl.Color{255, 255, 255, 255})
	ui.endMenu.exit = button{pos: pos{}, isClicked: false}
	ui.endMenu.exit.srcRect = append(ui.endMenu.exit.srcRect, &sdl.Rect{362, 349, 25, 32})
	ui.endMenu.exit.dstRect = append(ui.endMenu.exit.dstRect, &sdl.Rect{winWidth * .4, winHeight * .65, 25 * 2, 32 * 2})
	ui.endMenu.exit.srcRect = append(ui.endMenu.exit.srcRect, &sdl.Rect{313, 381, 70, 25})
	ui.endMenu.exit.dstRect = append(ui.endMenu.exit.dstRect, &sdl.Rect{winWidth*.4 + 24*2, winHeight*.65 + 3*2, 46 * 2, 25 * 2})
	ui.endMenu.exit.str = getTextTexture("Exit", sdl.Color{255, 255, 255, 0})
	ui.endMenu.exit.dstRect = append(ui.endMenu.exit.dstRect, &sdl.Rect{winWidth*.4 + 28*2, winHeight*.65 + 5*2, 36 * 2, 20 * 2})

	ui.endMenu.returnToMainMenu = button{pos: pos{}, isClicked: false}
	ui.endMenu.returnToMainMenu.srcRect = append(ui.endMenu.returnToMainMenu.srcRect, &sdl.Rect{310, 349, 25, 32})
	ui.endMenu.returnToMainMenu.dstRect = append(ui.endMenu.returnToMainMenu.dstRect, &sdl.Rect{winWidth*.4 - 50, winHeight*.65 - 80, 25 * 2, 32 * 2})
	ui.endMenu.returnToMainMenu.srcRect = append(ui.endMenu.returnToMainMenu.srcRect, &sdl.Rect{313, 381, 70, 25})
	ui.endMenu.returnToMainMenu.dstRect = append(ui.endMenu.returnToMainMenu.dstRect, &sdl.Rect{winWidth*.4 + 24*2 - 50, winHeight*.65 - 80 + 3*2, 91 * 2, 25 * 2})
	ui.endMenu.returnToMainMenu.str = getTextTexture("Main Menu", sdl.Color{255, 255, 255, 0})
	ui.endMenu.returnToMainMenu.dstRect = append(ui.endMenu.returnToMainMenu.dstRect, &sdl.Rect{winWidth*.4 + 28*2 - 50, winHeight*.65 - 80 + 5*2, 81 * 2, 20 * 2})

}

func drawEndMenuButtons(ui *UI2d) {
	renderer.Copy(ui.endMenu.gameEnd, nil, &sdl.Rect{175, 100, 450, 100})
	renderer.Copy(uiAtlas, ui.endMenu.exit.srcRect[0], ui.endMenu.exit.dstRect[0])
	renderer.Copy(uiAtlas, ui.endMenu.exit.srcRect[1], ui.endMenu.exit.dstRect[1])
	renderer.Copy(ui.endMenu.exit.str, nil, ui.endMenu.exit.dstRect[2])

	renderer.Copy(uiAtlas, ui.endMenu.returnToMainMenu.srcRect[0], ui.endMenu.returnToMainMenu.dstRect[0])
	renderer.Copy(uiAtlas, ui.endMenu.returnToMainMenu.srcRect[1], ui.endMenu.returnToMainMenu.dstRect[1])
	renderer.Copy(ui.endMenu.returnToMainMenu.str, nil, ui.endMenu.returnToMainMenu.dstRect[2])

}

func endMenu(ui *UI2d) stateFunc {
	renderer.Copy(mainMenuBackground, nil, nil)

	drawEndMenuButtons(ui)

	if ui.input.leftButton && !ui.input.prevLeftButton {
		inp := &sdl.Rect{int32(ui.input.x), int32(ui.input.y), 1, 1}
		if inp.HasIntersection(ui.endMenu.exit.dstRect[0]) || inp.HasIntersection(ui.endMenu.exit.dstRect[1]) {
			ui.endMenu.isTerminated = true

		}
		if inp.HasIntersection(ui.endMenu.returnToMainMenu.dstRect[0]) || inp.HasIntersection(ui.endMenu.returnToMainMenu.dstRect[1]) {
			ui.endMenu.isRestarted = true
		}
	}
	return determineToken
}
