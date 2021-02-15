package ui

import (
	"io/ioutil"
	"strings"

	"github.com/BiesGo/sdlWorkSpace/rpg/game"
	"github.com/veandco/go-sdl2/sdl"
)

type levelRelativity struct {
	startX, starY, endX, endY, relativeX, relativeY, scale int
}

var buttonSelected *sdl.Texture

type levelButton struct {
	levelName string
	isClicked bool
	texture   *sdl.Texture
	rect      *sdl.Rect
}

type selectMenuObj struct {
	levels []levelButton
	start  button
	edit   button
	rand   button
	levelRelativity
	preview button
}

func createSelectMenu(ui *UI2d) {
	ui.selectMenu = selectMenuObj{}
	ui.selectMenu.start = button{pos: pos{x: winWidth * .4, y: winHeight * .4}, isClicked: false}
	ui.selectMenu.start.srcRect = append(ui.selectMenu.start.srcRect, &sdl.Rect{310, 349, 25, 32})
	ui.selectMenu.start.srcRect = append(ui.selectMenu.start.srcRect, &sdl.Rect{313, 381, 70, 25})
	ui.selectMenu.start.dstRect = append(ui.selectMenu.start.dstRect, &sdl.Rect{0, 480 - 50, 25*1.5 - 0.5, 32 * 1.5})
	ui.selectMenu.start.dstRect = append(ui.selectMenu.start.dstRect, &sdl.Rect{0 + 24*1.5, 480 + 3*1.5 - 50 - 0.5, 55*1.5 - .5, 25 * 1.6})
	ui.selectMenu.start.str = getTextTexture("Start", sdl.Color{255, 255, 255, 0})
	ui.selectMenu.start.dstRect = append(ui.selectMenu.start.dstRect, &sdl.Rect{0 + 28*1.5, 480 + 5*1.5 - .5 - 50, 45*1.5 - 0.5, 20 * 1.5})

	ui.selectMenu.edit = button{pos: pos{x: winWidth * 0.4, y: winHeight * .4}, isClicked: false}
	ui.selectMenu.edit.srcRect = append(ui.selectMenu.edit.srcRect, &sdl.Rect{336, 349, 25, 32})
	ui.selectMenu.edit.srcRect = append(ui.selectMenu.edit.srcRect, &sdl.Rect{313, 381, 70, 25})
	ui.selectMenu.edit.dstRect = append(ui.selectMenu.edit.dstRect, &sdl.Rect{0, 480, 25*1.5 - 0.5, 32 * 1.5})
	ui.selectMenu.edit.dstRect = append(ui.selectMenu.edit.dstRect, &sdl.Rect{0 + 24*1.5, 480 + 3*1.5 - 0.5, 46 * 1.5, 25 * 1.6})
	ui.selectMenu.edit.str = getTextTexture("Edit", sdl.Color{255, 255, 255, 0})
	ui.selectMenu.edit.dstRect = append(ui.selectMenu.edit.dstRect, &sdl.Rect{0 + 28*1.5, 480 + 5*1.5 - .5, 36 * 1.5, 20 * 1.5})

	ui.selectMenu.rand = button{pos: pos{x: winWidth * .4, y: winHeight * .4}, isClicked: false}
	ui.selectMenu.rand.srcRect = append(ui.selectMenu.rand.srcRect, &sdl.Rect{362, 349, 25, 32})
	ui.selectMenu.rand.srcRect = append(ui.selectMenu.rand.srcRect, &sdl.Rect{313, 381, 70, 25})
	ui.selectMenu.rand.dstRect = append(ui.selectMenu.rand.dstRect, &sdl.Rect{0, 480 + 50, 25*1.5 - 0.5, 32 * 1.5})
	ui.selectMenu.rand.dstRect = append(ui.selectMenu.rand.dstRect, &sdl.Rect{0 + 24*1.5, 480 + 3*1.5 + 50 - 0.5, 82 * 1.5, 25 * 1.6})
	ui.selectMenu.rand.str = getTextTexture("ReCreate", sdl.Color{255, 255, 255, 0})
	ui.selectMenu.rand.dstRect = append(ui.selectMenu.rand.dstRect, &sdl.Rect{0 + 28*1.5, 480 + 5*1.5 - .5 + 50, 72 * 1.5, 20 * 1.5})

	ui.selectMenu.preview = button{pos: pos{}, isClicked: false}
	ui.selectMenu.preview.srcRect = append(ui.selectMenu.preview.srcRect, &sdl.Rect{311, 143, 20, 20})
	ui.selectMenu.preview.srcRect = append(ui.selectMenu.preview.srcRect, &sdl.Rect{332, 143, 20, 20})
	ui.selectMenu.preview.dstRect = append(ui.selectMenu.preview.dstRect, &sdl.Rect{0, 400, 20, 20})
	ui.selectMenu.preview.str = getTextTexture("Show Preview", sdl.Color{255, 255, 255, 0})
	ui.selectMenu.preview.dstRect = append(ui.selectMenu.preview.dstRect, &sdl.Rect{25, 400, 108, 20})

	files, err := ioutil.ReadDir("./game/maps/")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		name := f.Name()
		if strings.Contains(name, "level") && !strings.Contains(name, "Entities") && strings.HasSuffix(name, ".map") {
			ui.selectMenu.levels = append(ui.selectMenu.levels, levelButton{strings.TrimSuffix(name, ".map"), false, nil, nil})
		}
	}
	for i := range ui.selectMenu.levels {
		ui.selectMenu.levels[i].texture = getTextTexture(ui.selectMenu.levels[i].levelName, sdl.Color{255, 255, 255, 0})
		ui.selectMenu.levels[i].rect = &sdl.Rect{5, int32(40 * i), 90, 40}
	}

	ui.selectMenu.levelRelativity = levelRelativity{0, 0, 21, 19, 0, 0, 54}
}

func updateSelections(ui *UI2d) {
	if ui.input.leftButton && !ui.input.prevLeftButton {
		clickRect := &sdl.Rect{int32(ui.input.x), int32(ui.input.y), 1, 1}
		// if ui.selectMenu.edit.dstRect[0].HasIntersection(clickRect) || ui.selectMenu.edit.dstRect[1].HasIntersection(clickRect) {
		// 	for _, level := range ui.selectMenu.levels {
		// 		if level.isClicked {
		// 			globalLevel = &game.Level{}
		// 			globalLevel.LevelName = level.levelName
		// 			editBeforeStart = true
		// 			break
		// 		}
		// 	}
		// }
		if ui.selectMenu.start.dstRect[0].HasIntersection(clickRect) || ui.selectMenu.start.dstRect[1].HasIntersection(clickRect) {
			for _, level := range ui.selectMenu.levels {
				if level.isClicked {
					GlobalLevel2 = game.LoadLevelFromFile2("game/maps/" + level.levelName + ".map")
					currentState = playScreen
					// globalLevel = &game.Level{}
					// globalLevel.LevelName = level.levelName
					// editBeforeStart = false
					break
				}
			}
		}
		// if ui.selectMenu.rand.dstRect[0].HasIntersection(clickRect) || ui.selectMenu.rand.dstRect[1].HasIntersection(clickRect) {
		// 	for i, level := range ui.selectMenu.levels {
		// 		if level.isClicked {
		// 			go ui.ReCreatePreview(level.levelName, i)
		// 		}
		// 	}
		// }
		for i, level := range ui.selectMenu.levels {
			if level.rect.HasIntersection(clickRect) {
				for i := range ui.selectMenu.levels {
					ui.selectMenu.levels[i].isClicked = false
				}
				ui.selectMenu.levels[i].isClicked = true
				break
			}
		}
		if clickRect.HasIntersection(ui.selectMenu.preview.dstRect[0]) {
			ui.selectMenu.preview.isClicked = !ui.selectMenu.preview.isClicked
		}
	}
}

func updatePreviewRelativity(ui *UI2d) {
	if ui.input.currKeyState[sdl.SCANCODE_RIGHT] != 0 && ui.input.prevKeyState[sdl.SCANCODE_RIGHT] == 0 {
		if ui.selectMenu.endX+21 <= 105 {
			ui.selectMenu.startX += 21
			ui.selectMenu.endX += 21
			ui.selectMenu.relativeX -= 672
		}
	}
	if ui.input.currKeyState[sdl.SCANCODE_LEFT] != 0 && ui.input.prevKeyState[sdl.SCANCODE_LEFT] == 0 {
		if ui.selectMenu.startX-21 >= 0 {
			ui.selectMenu.startX -= 21
			ui.selectMenu.endX -= 21
			ui.selectMenu.relativeX += 672
		}
	}
	if ui.input.currKeyState[sdl.SCANCODE_UP] != 0 && ui.input.prevKeyState[sdl.SCANCODE_UP] == 0 {
		if ui.selectMenu.starY-19 >= 0 {
			ui.selectMenu.starY -= 19
			ui.selectMenu.endY -= 19
			ui.selectMenu.relativeY += 608
		}
	}
	if ui.input.currKeyState[sdl.SCANCODE_DOWN] != 0 && ui.input.prevKeyState[sdl.SCANCODE_DOWN] == 0 {
		if ui.selectMenu.endY+19 <= 114 {
			ui.selectMenu.starY += 19
			ui.selectMenu.endY += 19
			ui.selectMenu.relativeY -= 608
		}
	}
}

func showPreview(level *layer, ui *UI2d) {
	startX := ui.selectMenu.startX
	starY := ui.selectMenu.starY
	endX := ui.selectMenu.endX
	endY := ui.selectMenu.endY

	for y := starY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			if x > -1 && y > -1 && x < 100 && y < 100 && level.dstRect[y][x] != nil {
				renderer.Copy(textureAtlas, level.srcRect[y][x], &sdl.Rect{150 + (int32(x%21) * 32), (int32(y%19) * 32), 32, 32})
			}
		}
	}
	// relativeX := ui.selectMenu.relativeX
	// relativeY := ui.selectMenu.relativeY
	// for _, intf := range level.entities {
	// 	obj := intf.(entityInterface)
	// 	if obj.getX()+relativeX+150 > 150 {
	// 		renderer.Copy(textureAtlas, obj.getRect(), &sdl.Rect{150 + int32(obj.getX()+relativeX), int32(obj.getY() + relativeY), 32, 32})
	// 	}
	// }
}

func selectMenuMiniMap(level *layer, ui *UI2d) {
	scale := ui.selectMenu.scale
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if level.srcRect[y][x] != nil {
				renderer.Copy(textureAtlas, level.srcRect[y][x], &sdl.Rect{150 + 600/int32(scale)*int32(x), 600 / int32(scale) * int32(y), int32((600 / scale)), int32((600 / scale))})
			}
		}
	}
	for _, intf := range level.entities {
		obj := intf.(entityInterface)
		renderer.Copy(textureAtlas, obj.getRect(), &sdl.Rect{150 + 600/int32(scale)*int32(obj.getX()/32), 600 / int32(scale) * int32(obj.getY()/32), int32((600 / scale)), int32((600 / scale))})
	}
}

func updateZoomScale(ui *UI2d) {
	if ui.input.currKeyState[sdl.SCANCODE_PAGEDOWN] != 0 && ui.input.prevKeyState[sdl.SCANCODE_PAGEDOWN] == 0 {
		if ui.selectMenu.scale+9 < 100 {
			ui.selectMenu.scale += 9
		}
	}
	if ui.input.currKeyState[sdl.SCANCODE_PAGEUP] != 0 && ui.input.prevKeyState[sdl.SCANCODE_PAGEUP] == 0 {
		if ui.selectMenu.scale-9 > 0 {
			ui.selectMenu.scale -= 9
		}
	}
}

func selectMenu(ui *UI2d) stateFunc {
	renderer.Copy(mainMenuBackground, nil, nil)
	renderer.Copy(ui.mainMenu.infoTab, nil, &sdl.Rect{0, 0, 150, winHeight})
	updateSelections(ui)

	for i, level := range ui.selectMenu.levels {
		if level.isClicked {
			GlobalLevel2 = game.LoadLevelFromFile2("game/maps/" + level.levelName + ".map")

			px := createOnePixel(255, 255, 255, 200)
			renderer.Copy(px, nil, &sdl.Rect{0, int32(40 * i), 110, 40})
			renderer.Copy(ui.mainMenu.infoTab, nil, &sdl.Rect{5, int32(i*40) + 5, 100, 30})
			if ui.selectMenu.preview.isClicked {
				if ui.input.currKeyState[sdl.SCANCODE_LSHIFT] != 0 {
					updateZoomScale(ui)
					selectMenuMiniMap(&ui.levelPreviews[i], ui)
				} else {
					updatePreviewRelativity(ui)
					showPreview(&ui.levelPreviews[i], ui)
				}
				if ui.input.currKeyState[sdl.SCANCODE_LSHIFT] == 0 && ui.input.prevKeyState[sdl.SCANCODE_LSHIFT] == 0 {
					ui.selectMenu.scale = 50
				}
			}
		}
		renderer.Copy(level.texture, nil, level.rect)
	}

	for i := 0; i < 2; i++ {
		renderer.Copy(uiAtlas, ui.selectMenu.start.srcRect[i], ui.selectMenu.start.dstRect[i])
		//renderer.Copy(uiAtlas, ui.selectMenu.edit.srcRect[i], ui.selectMenu.edit.dstRect[i])
		//renderer.Copy(uiAtlas, ui.selectMenu.rand.srcRect[i], ui.selectMenu.rand.dstRect[i])
	}
	renderer.Copy(ui.selectMenu.start.str, nil, ui.selectMenu.start.dstRect[2])
	// renderer.Copy(ui.selectMenu.edit.str, nil, ui.selectMenu.edit.dstRect[2])
	// renderer.Copy(ui.selectMenu.rand.str, nil, ui.selectMenu.rand.dstRect[2])

	if ui.selectMenu.preview.isClicked {
		renderer.Copy(uiAtlas, ui.selectMenu.preview.srcRect[1], ui.selectMenu.preview.dstRect[0])
	} else {
		renderer.Copy(uiAtlas, ui.selectMenu.preview.srcRect[0], ui.selectMenu.preview.dstRect[0])
	}
	renderer.Copy(ui.selectMenu.preview.str, nil, ui.selectMenu.preview.dstRect[1])

	if ui.input.currKeyState[sdl.SCANCODE_ESCAPE] != 0 && ui.input.prevKeyState[sdl.SCANCODE_ESCAPE] == 0 {
		currentState = mainScreen
	}

	// WHEN YOU PRESS 'P' the game will start
	if ui.input.currKeyState[sdl.SCANCODE_P] != 0 && ui.input.prevKeyState[sdl.SCANCODE_P] == 0 {
		currentState = playScreen
	}

	return determineToken
}
