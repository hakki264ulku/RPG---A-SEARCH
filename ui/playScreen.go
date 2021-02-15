package ui

import (
	"math/rand"

	"github.com/BiesGo/sdlWorkSpace/rpg/game"
	"github.com/veandco/go-sdl2/sdl"
)

// drawHealthBars uses the gloabal 'healthBarTextures' array to draw the healthbars
func drawHealthBars(c *game.Player) {

	charSrcRect := &textureIndex['P'][0]
	renderer.Copy(textureAtlas, charSrcRect, &sdl.Rect{0, 0, 16, 16})

	unitValue := (1 * c.FullHitpoints) / 15
	var index int
	index = -1
	if unitValue*15 == c.Hitpoints || unitValue*14 < c.Hitpoints { // if health is full
		index = 14
	} else {
		for i := 14; i >= 0; i-- {
			if unitValue*i >= c.Hitpoints && unitValue*(i-1) < c.Hitpoints {
				index = i
				break
			}
		}
		if index == -1 {
			index = 14
		}
	}

	texture := healthBarTextures[index]
	renderer.Copy(texture, &sdl.Rect{5, 25, 380, 75}, &sdl.Rect{16, 0, 380 / 4, 75 / 4})
}

func drawPods(level *game.Level2, offsetX, offsetY int32) {
	for pos := range level.HealthPots {
		renderer.Copy(potTexture, &sdl.Rect{415, 5, 185, 295}, &sdl.Rect{int32(pos.X)*32 + offsetX + 16, int32(pos.Y)*32 + offsetY + 8, 16, 24})
	}
}

func drawHealthBarsEnemy() {
	i := 0
	for _, c := range GlobalLevel2.EnemiesForHealthBars {
		charSrcRect := &textureIndex[c.Tile][0]
		renderer.Copy(textureAtlas, charSrcRect, &sdl.Rect{winWidth - 100, int32(i * 16), 16, 16})

		unitValue := (1 * c.FullHitpoints) / 15
		var index int
		index = -1
		if unitValue*15 == c.Hitpoints || unitValue*14 < c.Hitpoints { // if health is full
			index = 14
		} else {
			for i := 14; i >= 0; i-- {
				if unitValue*i >= c.Hitpoints && unitValue*(i-1) < c.Hitpoints {
					index = i
					break
				}
			}
			if index == -1 {
				index = 14
			}
		}

		texture := healthBarTextures[index]
		renderer.Copy(texture, &sdl.Rect{5, 25, 380, 75}, &sdl.Rect{winWidth - 100 + 16, int32(i * 16), 380 / 4, 75 / 4})
		i++
	}
}

func playMenu(ui *UI2d) stateFunc {
	renderer.Clear() // Clears everything on the screen 

	newLevel := GlobalLevel2

	r := rand.New(rand.NewSource(1))

	if centerX == -1 && centerY == -1 {
		centerX = newLevel.Player.X
		centerY = newLevel.Player.Y
	}

	limit := 5
	if newLevel.Player.X > centerX+limit {
		centerX++
	} else if newLevel.Player.X < centerX-limit {
		centerX--
	} else if newLevel.Player.Y > centerY+limit {
		centerY++
	} else if newLevel.Player.Y < centerY-limit {
		centerY--
	}
	var offsetX int32
	var offsetY int32
	lastOffsetX := offsetX
	lastOffsetY := offsetY

	offsetX = int32((winWidth / 2) - centerX*32)
	offsetY = int32((winHeight / 2) - centerY*32)

	for y, row := range newLevel.Map {
		for x, tile := range row {

			if tile != game.Blank {
				srcRects := textureIndex[tile]
				srcRect := srcRects[r.Intn(len(srcRects))] // get a random tile from a specific group of rects,
				//this makes difference if the variaton count of the current tile is greater than 1
				dstRect := sdl.Rect{int32(x*32) + offsetX, int32(y*32) + offsetY, 32, 32}
				// for seeing how breadth-first search works, can be said that this is going to be for debugging purposes
				renderer.Copy(textureAtlas, &srcRect, &dstRect) // WRITES THE SRC RECT TO THE DST RECT
			}
		}
	}

	playerSrcRect := textureIndex['P'][0]

	if lastOffsetX != offsetX || lastOffsetY != offsetY {
		renderer.Copy(textureAtlas, &playerSrcRect, &sdl.Rect{int32(newLevel.Player.X)*32 + offsetX, int32(newLevel.Player.Y)*32 + offsetY, 32, 32})
	}
	isInputTaken := game.HandleInput(ui.input.currKeyState, ui.input.prevKeyState, newLevel)

	// for the sake of TURN BASED Playing ability
	if isInputTaken {
		for _, monster := range GlobalLevel2.Enemies {
			monster.Update(GlobalLevel2)
		}
		for _, monster := range newLevel.Enemies {
			monsterSrcrect := textureIndex[game.Tile(monster.Tile)][0]
			renderer.Copy(textureAtlas, &monsterSrcrect, &sdl.Rect{int32(monster.X)*32 + offsetX, int32(monster.Y)*32 + offsetY, 32, 32})
		}
	}

	for _, monster := range newLevel.Enemies {
		monsterSrcrect := textureIndex[game.Tile(monster.Tile)][0]
		renderer.Copy(textureAtlas, &monsterSrcrect, &sdl.Rect{int32(monster.X)*32 + offsetX, int32(monster.Y)*32 + offsetY, 32, 32})
	}

	drawHealthBarsEnemy()
	drawHealthBars(newLevel.Player)
	drawPods(newLevel, offsetX, offsetY)
	if ui.input.currKeyState[sdl.SCANCODE_T] != 0 && ui.input.prevKeyState[sdl.SCANCODE_T] == 0 {
		// game.Save2(GlobalLevel2)
		// GlobalLevel2 = game.LoadLevelFromFile2(newLevel.FileName)
		currentState = endScreen
	}

	return determineToken
}
