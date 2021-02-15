package ui

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/BiesGo/sdlWorkSpace/rpg/game"
	"github.com/veandco/go-sdl2/sdl"
)

var editingTile game.Tile = game.DirtFloor

// controlling the tile change menu with the numbers
func checkEditingTileChange(ui *UI2d) {
	if ui.input.currKeyState[sdl.SCANCODE_1] != 0 && ui.input.prevKeyState[sdl.SCANCODE_1] == 0 {
		editingTile = editingTileSlice[0]
	}
	if ui.input.currKeyState[sdl.SCANCODE_2] != 0 && ui.input.prevKeyState[sdl.SCANCODE_2] == 0 {
		editingTile = editingTileSlice[1]
	}
	if ui.input.currKeyState[sdl.SCANCODE_3] != 0 && ui.input.prevKeyState[sdl.SCANCODE_3] == 0 {
		editingTile = editingTileSlice[2]
	}
	if ui.input.currKeyState[sdl.SCANCODE_4] != 0 && ui.input.prevKeyState[sdl.SCANCODE_4] == 0 {
		editingTile = editingTileSlice[3]
	}
}

func getTileType() int {
	switch editingTile {
	case game.DirtFloor:
		return 0
	case game.StoneWall:
		return 0
	case game.MainCharacter:
		return 0
	case game.DoorC, game.DoorO:
		return 0
	case game.Spider, game.Rat:
		return 0
	default:
		panic("unknown tile in getTileLayer")
	}
}

var editingTileSlice = []game.Tile{
	game.DirtFloor,
	game.StoneWall,
	game.DoorC,
	game.MainCharacter,
	game.Spider,
	game.Rat,
}
var globalLevel *game.Level

func addToGridWorld(x, y int, tile game.Tile) {
	gridY := len(globalLevel.GridWorld.Rows)
	for gridY < y+1 {
		globalLevel.GridWorld.Rows = append(globalLevel.GridWorld.Rows, game.Row{})
		gridY++
	}
	gridX := len(globalLevel.GridWorld.Rows[y].Grids)
	for gridX < x+1 {
		globalLevel.GridWorld.Rows[y].Grids = append(globalLevel.GridWorld.Rows[y].Grids, game.Grid{Layers: []game.Tile{}})
		gridX++
	}
	if tile == game.Blank || tile == game.DirtFloor || tile == game.StoneWall {
		globalLevel.GridWorld.Rows[y].Grids[x].Background = tile
	}
}

//addTileToTheMap adds the current editing tile to the current GlobalLevel2
func addTileToTheMap(x, y int, editingTile game.Tile) {
	fmt.Println(x, y, editingTile)
	for Y, row := range GlobalLevel2.Map {
		for X := range row {
			if x == X && y == Y {
				GlobalLevel2.Map[y][x] = editingTile
				fmt.Println("Editing tile is added the place")
			}
		}
	}

	file, err := os.OpenFile(GlobalLevel2.FileName, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	file.Truncate(0)
	file.Seek(0, 0)
	for _, row := range GlobalLevel2.Map {
		for _, t := range row {
			if t != game.Blank {
				file.WriteString(t.ToString())
			} else {
				file.WriteString(" ")
			}
		}
		file.WriteString("\n")
	}
	file.Close()

}

func editTile(ui *UI2d) {
	if ui.input.leftButton { // && !ui.input.prevLeftButton
		x := int(math.Floor(float64(ui.input.x)/32)) + ui.editMenu.startX
		y := int(math.Floor(float64(ui.input.y)/32)) + ui.editMenu.starY
		l := getTileType()
		if ui.background.dstRect[y][x] == nil || l == 0 {
			ui.background.srcRect[y][x] = &textureIndex[editingTile][rand.Intn(len(textureIndex[editingTile]))]
			ui.background.dstRect[y][x] = &sdl.Rect{X: int32(x) * 32, Y: int32(y) * 32, W: 32, H: 32}
			// add editing tile to the map
			GlobalLevel2.Map[y][x] = editingTile
			go addTileToTheMap(x, y, editingTile)

			addToGridWorld(x, y, editingTile)
		} else if l == 1 && !ui.input.prevLeftButton {
			globalLevel.Entities = append(globalLevel.Entities, game.Entity{game.Pos{x * 32, y * 32}, editingTile})
			ui.background.entities = append(ui.background.entities, getEntity(globalLevel.Entities[len(globalLevel.Entities)-1]))
		}

	}

	if ui.input.rightButton && !ui.input.prevRightButton {
		isDeleted := false
		x := int(math.Floor(float64(ui.input.x)/32)) + ui.editMenu.startX
		y := int(math.Floor(float64(ui.input.y)/32)) + ui.editMenu.starY
		if len(ui.background.entities) > 0 {
			for i, intf := range ui.background.entities {
				obj := intf.(entityInterface)
				if ui.input.x < obj.getX()+32 && ui.input.x >= obj.getX() && ui.input.y < obj.getY()+32 && ui.input.y >= obj.getY() {
					if len(ui.background.entities) > 1 {
						ui.background.entities = append(ui.background.entities[0:i], ui.background.entities[i+1:len(ui.background.entities)]...)
						globalLevel.Entities = append(globalLevel.Entities[0:i], globalLevel.Entities[i+1:len(globalLevel.Entities)]...)
					} else {
						ui.background.entities = ui.background.entities[0:0]
						globalLevel.Entities = globalLevel.Entities[0:0]
					}
					isDeleted = true
				}
			}
		}
		if !isDeleted && ui.background.dstRect[y][x] != nil { //!isDeleted &&
			ui.background.dstRect[y][x] = nil
			ui.background.srcRect[y][x] = nil
			globalLevel.GridWorld.Rows[y].Grids[x].Background = game.Blank
		}
	}

	if ui.input.currKeyState[sdl.SCANCODE_BACKSPACE] != 0 && ui.input.prevKeyState[sdl.SCANCODE_BACKSPACE] == 0 {
		fmt.Println("Level Reloaded")
		globalLevel = globalLevel.ReLoadTheLevel()
		createLayers(globalLevel, ui)

	}
}

// currTileChangeMenu runs when you press LShift
func currTileChangeMenu(ui *UI2d) {

	var w int = 240
	var h int = int(math.Ceil(float64(len(editingTileSlice))/3.0))*76 + 48
	var x int32 = int32(800-w) / 2
	var y int32 = int32(600-h) / 2
	var tileTabDst []*sdl.Rect
	renderer.Copy(ui.mainMenu.infoTab, &sdl.Rect{0, 0, 1, 1}, &sdl.Rect{x, y, int32(w), int32(h)})
	j := 0
	for i, tile := range editingTileSlice {
		if 76*(j+1) > w {
			y = y + 76
			j = 0
		}
		j++
		if editingTile == tile {
			px := createOnePixel(255, 255, 255, 200)
			renderer.Copy(px, nil, &sdl.Rect{x + 6 + (int32(76*i))%int32(w-12), y, 76, 76})
			renderer.Copy(ui.mainMenu.infoTab, &sdl.Rect{0, 0, 1, 1}, &sdl.Rect{x + 12 + (int32(76*i))%int32(w-12), y + 6, 64, 64})
		}
		tileTabDst = append(tileTabDst, &sdl.Rect{x + 12 + (int32(76*i))%int32(w-12), y + 6, 64, 64})
		renderer.Copy(textureAtlas, &textureIndex[tile][0], tileTabDst[i])
	}

	x = 360
	y += 76
	renderer.Copy(uiAtlas, ui.selectMenu.start.srcRect[0], &sdl.Rect{x, y, 25*1.5 - 0.5, 32 * 1.5})
	renderer.Copy(uiAtlas, ui.selectMenu.start.srcRect[1], &sdl.Rect{x + 24*1.5, y + 4, 55*1.5 - .5, 25 * 1.6})
	renderer.Copy(ui.selectMenu.start.str, nil, &sdl.Rect{x + 28*1.5, y + 7, 45*1.5 - 0.5, 20 * 1.5})

	if ui.input.leftButton && !ui.input.prevLeftButton {
		clickRect := &sdl.Rect{int32(ui.input.x), int32(ui.input.y), 1, 1}
		for i, rect := range tileTabDst {
			if clickRect.HasIntersection(rect) {
				fmt.Println("Editing tile is set")
				editingTile = editingTileSlice[i]
				break
			}
		}
		if clickRect.HasIntersection(&sdl.Rect{x, y, 25*1.5 - 0.5, 32 * 1.5}) || clickRect.HasIntersection(&sdl.Rect{x + 24*1.5, y + 4, 55*1.5 - .5, 25 * 1.6}) {
			game.Save2(GlobalLevel2)
			GlobalLevel2 = game.LoadLevelFromFile2(GlobalLevel2.FileName)
			currentState = playScreen
		}
	}
}

type editMenuObj struct {
	levelRelativity
}

func createEditMenu(ui *UI2d) {
	ui.editMenu = editMenuObj{}
	ui.editMenu.levelRelativity = levelRelativity{0, 0, 25, 19, 0, 0, 54}
}

func updateEditRelativity(ui *UI2d) {
	if ui.input.currKeyState[sdl.SCANCODE_RIGHT] != 0 && ui.input.prevKeyState[sdl.SCANCODE_RIGHT] == 0 {
		if ui.editMenu.endX+25 <= 100 {
			ui.editMenu.startX += 25
			ui.editMenu.endX += 25
			ui.editMenu.relativeX -= 800
		}
	}
	if ui.input.currKeyState[sdl.SCANCODE_LEFT] != 0 && ui.input.prevKeyState[sdl.SCANCODE_LEFT] == 0 {
		if ui.editMenu.startX-25 >= 0 {
			ui.editMenu.startX -= 25
			ui.editMenu.endX = ui.editMenu.startX + 25
			ui.editMenu.relativeX += 800
		}
	}
	if ui.input.currKeyState[sdl.SCANCODE_UP] != 0 && ui.input.prevKeyState[sdl.SCANCODE_UP] == 0 {
		if ui.editMenu.starY-19 >= 0 {
			ui.editMenu.starY -= 19
			ui.editMenu.endY -= 19
			ui.editMenu.relativeY += 608
		}
	}
	if ui.input.currKeyState[sdl.SCANCODE_DOWN] != 0 && ui.input.prevKeyState[sdl.SCANCODE_DOWN] == 0 {
		if ui.editMenu.endY+19 <= 114 {
			ui.editMenu.starY += 19
			ui.editMenu.endY += 19
			ui.editMenu.relativeY -= 608
		}
	}
}

func showEditLevel(ui *UI2d) {
	startX := ui.editMenu.startX
	starY := ui.editMenu.starY
	endX := ui.editMenu.endX
	endY := ui.editMenu.endY

	for y := starY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			if x > -1 && y > -1 && x < 100 && y < 100 && ui.background.dstRect[y][x] != nil {
				renderer.Copy(textureAtlas, ui.background.srcRect[y][x], &sdl.Rect{(int32(x%25) * 32), (int32(y%19) * 32), 32, 32})
			}
		}
	}
	// relativeX := ui.editMenu.relativeX
	// relativeY := ui.editMenu.relativeY
	// for _, intf := range ui.background.entities {
	// 	obj := intf.(entityInterface)
	// 	renderer.Copy(textureAtlas, obj.getRect(), &sdl.Rect{int32(obj.getX() + relativeX), int32(obj.getY() + relativeY), 32, 32})
	// }
}

func editMenuMiniMap(ui *UI2d) {
	scale := ui.editMenu.scale
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if ui.background.srcRect[y][x] != nil {
				renderer.Copy(textureAtlas, ui.background.srcRect[y][x], &sdl.Rect{600 / int32(scale) * int32(x), 600 / int32(scale) * int32(y), int32((600 / scale)), int32((600 / scale))})
			}
		}
	}
	for _, intf := range ui.background.entities {
		obj := intf.(entityInterface)
		renderer.Copy(textureAtlas, obj.getRect(), &sdl.Rect{600 / int32(scale) * int32(obj.getX()/32), 600 / int32(scale) * int32(obj.getY()/32), int32((600 / scale)), int32((600 / scale))})
	}
}
func updateEditScale(ui *UI2d) {
	if ui.input.currKeyState[sdl.SCANCODE_PAGEDOWN] != 0 && ui.input.prevKeyState[sdl.SCANCODE_PAGEDOWN] == 0 {
		if ui.editMenu.scale+9 < 100 {
			ui.editMenu.scale += 9
		}
	}
	if ui.input.currKeyState[sdl.SCANCODE_PAGEUP] != 0 && ui.input.prevKeyState[sdl.SCANCODE_PAGEUP] == 0 {
		if ui.editMenu.scale-9 > 0 {
			ui.editMenu.scale -= 9
		}
	}
}

func editMenu(ui *UI2d) stateFunc {
	if ui.input.currKeyState[sdl.SCANCODE_S] == 0 && ui.input.prevKeyState[sdl.SCANCODE_S] != 0 {
		game.Save2(GlobalLevel2)
		GlobalLevel2 = game.LoadLevelFromFile2(GlobalLevel2.FileName)
		fmt.Println("saving done")
	}
	checkEditingTileChange(ui)

	renderer.Copy(mainMenuBackground, nil, nil)
	//renderer.Copy(blackPixel, nil, &sdl.Rect{0, 0, winWidth, winHeight})

	if ui.input.currKeyState[sdl.SCANCODE_TAB] != 0 {
		updateEditScale(ui)
		editMenuMiniMap(ui)
	} else {
		showEditLevel(ui)
	}
	if ui.input.currKeyState[sdl.SCANCODE_LSHIFT] != 0 {
		currTileChangeMenu(ui)
	} else if ui.input.currKeyState[sdl.SCANCODE_TAB] == 0 {
		updateEditRelativity(ui)
		editTile(ui)
	}

	if ui.input.currKeyState[sdl.SCANCODE_P] != 0 && ui.input.prevKeyState[sdl.SCANCODE_P] == 0 {
		GlobalLevel2 = game.LoadLevelFromFile2(GlobalLevel2.FileName)
		currentState = playScreen
	}

	return determineToken
}
