package ui

import (
	"bufio"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BiesGo/sdlWorkSpace/rpg/game"
	"github.com/veandco/go-sdl2/sdl"
)

func (result *inputState) updateMouseState() {
	result.prevLeftButton = result.leftButton
	result.prevRightButton = result.rightButton
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	result.x = int(mouseX)
	result.y = int(mouseY)
	result.leftButton = leftButton != 0
	result.rightButton = rightButton != 0

}
func (result *inputState) updateKeyboardState() {
	for i := range result.currKeyState {
		result.prevKeyState[i] = result.currKeyState[i]
	}
}

func loadTextureIndex() {
	textureIndex = make(map[game.Tile][]sdl.Rect)
	infile, err := os.Open("ui/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}
	defer infile.Close()
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := game.Tile(line[0])
		xy := line[2:]
		splitXYC := strings.Split(xy, ",")
		x, err := strconv.ParseInt(splitXYC[0], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(splitXYC[1], 10, 64)
		if err != nil {
			panic(err)
		}
		variationCount, err := strconv.ParseInt(splitXYC[2], 10, 64)
		if err != nil {
			panic(err)
		}

		var rects []sdl.Rect
		for i := 0; i < int(variationCount); i++ {
			rects = append(rects, sdl.Rect{X: int32(x * 32), Y: int32(y * 32), W: 32, H: 32})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		textureIndex[tileRune] = rects

	}
}

func imgFileToTexture(filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)

	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	return tex
}

func getEntity(obj game.Entity) entityInterface {
	switch obj.Tile {
	case '|':
		return newDoor(obj)
	case 'P':
		return createMainCharacter(obj)
	case 'C', 'c':
		return createChest(obj)
	case 'm':
		return createMonster(obj)
	case 'R', 'S':
		return createEnemy(obj)
	}

	panic("error")
}

func createLayers(level *game.Level, ui *UI2d) {
	for y := range ui.background.srcRect {
		for x := range ui.background.srcRect[y] {
			ui.background.srcRect[y][x] = nil
			ui.background.dstRect[y][x] = nil
		}
	}

	gridWorld := level.GridWorld
	for y, row := range gridWorld.Rows {
		for x, grid := range row.Grids {

			layer := grid.Background
			if layer != game.Blank {
				srcRects := textureIndex[layer]
				ui.background.srcRect[y][x] = &srcRects[rand.Intn(len(srcRects))]
				ui.background.dstRect[y][x] = &sdl.Rect{X: int32(x) * 32, Y: int32(y) * 32, W: 32, H: 32}

				renderer.Copy(textureAtlas, ui.background.srcRect[y][x], ui.background.dstRect[y][x])
			} else {
				ui.background.srcRect[y][x] = nil
				ui.background.dstRect[y][x] = nil

			}

		}

	}
	ui.background.entities = ui.background.entities[0:0]
	for _, obj := range level.Entities {
		ui.background.entities = append(ui.background.entities, getEntity(obj))
	}
	ui.mc = ui.background.entities[0].(*mainCharacter)
}

func (ui *UI2d) AddPreview(level game.Level) {
	ui.levelPreviews = append(ui.levelPreviews, layer{})
	index := len(ui.levelPreviews) - 1
	for y := range ui.levelPreviews[index].srcRect {
		for x := range ui.levelPreviews[index].srcRect[y] {
			ui.levelPreviews[index].srcRect[y][x] = nil
			ui.levelPreviews[index].dstRect[y][x] = nil
		}
	}

	gridWorld := level.GridWorld
	for y, row := range gridWorld.Rows {
		for x, grid := range row.Grids {

			layer := grid.Background
			if layer != game.Blank {
				srcRects := textureIndex[layer]
				ui.levelPreviews[index].srcRect[y][x] = &srcRects[rand.Intn(len(srcRects))]
				ui.levelPreviews[index].dstRect[y][x] = &sdl.Rect{X: 150 + int32(x)*32, Y: int32(y) * 32, W: 32, H: 32}

				renderer.Copy(textureAtlas, ui.levelPreviews[index].srcRect[y][x], ui.levelPreviews[index].dstRect[y][x])
			} else {
				ui.levelPreviews[index].srcRect[y][x] = nil
				ui.levelPreviews[index].dstRect[y][x] = nil
			}

		}

	}
	for _, obj := range level.Entities {
		ui.levelPreviews[index].entities = append(ui.levelPreviews[index].entities, getEntity(obj))
	}
}

func (ui *UI2d) ReCreatePreview(levelName string, index int) {
	level := game.CreateRandomMaze(levelName, ui)
	for y := range ui.levelPreviews[index].srcRect {
		for x := range ui.levelPreviews[index].srcRect[y] {
			ui.levelPreviews[index].srcRect[y][x] = nil
			ui.levelPreviews[index].dstRect[y][x] = nil
		}
	}

	gridWorld := level.GridWorld
	for y, row := range gridWorld.Rows {
		for x, grid := range row.Grids {

			layer := grid.Background
			if layer != game.Blank {
				srcRects := textureIndex[layer]
				ui.levelPreviews[index].srcRect[y][x] = &srcRects[rand.Intn(len(srcRects))]
				ui.levelPreviews[index].dstRect[y][x] = &sdl.Rect{X: 150 + int32(x)*32, Y: int32(y) * 32, W: 32, H: 32}

				renderer.Copy(textureAtlas, ui.levelPreviews[index].srcRect[y][x], ui.levelPreviews[index].dstRect[y][x])
			} else {
				ui.levelPreviews[index].srcRect[y][x] = nil
				ui.levelPreviews[index].dstRect[y][x] = nil
			}

		}

	}
	ui.levelPreviews[index].entities = ui.levelPreviews[index].entities[0:0]
	for _, obj := range level.Entities {
		ui.levelPreviews[index].entities = append(ui.levelPreviews[index].entities, getEntity(obj))
	}
}

func determineToken(ui *UI2d) stateFunc {
	switch currentState {
	case mainScreen:
		return mainMenu(ui)
	case editLevel:
		return editMenu(ui)
	case selectScreen:
		return selectMenu(ui)
	case endScreen:
		return endMenu(ui)
	case playScreen:
		return playMenu(ui)
	default:
		return nil
	}
}

func createOnePixel(r, g, b, a byte) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	tex.GetBlendMode()
	if err != nil {
		panic(err)
	}
	pixels := make([]byte, 4)
	pixels[0] = r
	pixels[1] = g
	pixels[2] = b
	pixels[3] = a
	tex.Update(nil, pixels, 4)

	return tex
}

func getTextTexture(str string, color sdl.Color) *sdl.Texture {
	textSurface, _ := font.RenderUTF8Solid(str, color)
	textTexture, _ := renderer.CreateTextureFromSurface(textSurface)
	return textTexture
}

func (ui *UI2d) Draw(level *game.Level, startingState bool) bool {
	if startingState {
		currentState = editLevel
	} else {
		currentState = endScreen
	}
	globalLevel = level
	GlobalLevel2 = game.LoadLevelFromFile2("./game/maps/" + level.LevelName + ".map")
	ui.background = layer{}

	createLayers(level, ui)

	for {
		currTime := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) { // theEvent := event.(type) //remember this
			case *sdl.QuitEvent:
				return false
			}
		}

		//FUNCTION-SELECTOR MECHANISM
		determineToken(ui)

		if ui.endMenu.isTerminated {
			return false // with this return value, the infinite loop in the game.Run becomes broken so that game screen ends
		}
		if ui.endMenu.isRestarted {
			return true
		}
		renderer.Present()

		// INPUT UPDATES
		ui.input.updateKeyboardState()
		ui.input.updateMouseState()

		// FRAME LOCKS
		elapsedTime := time.Since(currTime).Milliseconds()
		if elapsedTime > 10 {
			fmt.Println("elapsed time->", elapsedTime)
			//sdl.Delay(uint32(16 - elapsedTime))
		}
		if elapsedTime < 16 {
			sdl.Delay(uint32(16 - elapsedTime))
		}
	}
}

var editBeforeStart bool

func (ui *UI2d) SelectLevel() (*game.Level, bool) {
	currentState = mainScreen
	globalLevel = nil

	//TODO ELAPSED TIME IS COMMENTED OUT FROM THE CODE becasue when player moves there was happenening a bad delay on tiles ????

	//start := time.Now()
	for {
		currTime := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) { // theEvent := event.(type) //remember this
			case *sdl.QuitEvent:
				return nil, false
			}
		}
		determineToken(ui)

		if globalLevel != nil {
			return globalLevel, editBeforeStart
		}

		if ui.endMenu.isTerminated {
			return globalLevel, false // with this return value, the infinite loop in the game.Run becomes broken so that game screen ends
		}
		if ui.endMenu.isRestarted {
			return globalLevel, true
		}

		renderer.Present()

		ui.input.updateKeyboardState()
		ui.input.updateMouseState()
		elapsedTime := time.Since(currTime).Milliseconds()
		if elapsedTime > 10 {
			fmt.Println(elapsedTime)
		}
		if elapsedTime < 16 {
			sdl.Delay(uint32(16 - elapsedTime))
		}
		sdl.Delay(1)
	}
}
