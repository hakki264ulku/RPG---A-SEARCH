package ui

import (
	"github.com/BiesGo/sdlWorkSpace/rpg/game"
	"github.com/veandco/go-sdl2/sdl"
)

var chestSrcStorage []*sdl.Rect

func initEntityInterface(ui *UI2d) {
	chestSrcStorage = append(chestSrcStorage, &textureIndex[game.ChestC][0])
	chestSrcStorage = append(chestSrcStorage, &textureIndex[game.ChestO][0])
}

type entityInterface interface {
	getX() int
	getY() int
	getRect() *sdl.Rect
}

func getX(intf interface{}) int {
	switch t := intf.(type) {
	case Door:
		return t.x
	case mainCharacter:
		return t.x
	}
	panic("error")
}

type entity struct {
	x, y    int
	srcRect *sdl.Rect
}

type Door struct {
	entity
	is_open    bool
	srcStorage []*sdl.Rect
}

func newDoor(obj game.Entity) *Door {
	if obj.Tile == game.DoorC {
		return &Door{entity{obj.X, obj.Y, &textureIndex[obj.Tile][0]}, false, nil} // , make([]*sdl.Rect, 2)
	} else if obj.Tile == game.DoorO {
		return &Door{entity{obj.X, obj.Y, &textureIndex[obj.Tile][0]}, true, nil}
	}
	panic("error")
}

func (d Door) getX() int {
	return d.x
}
func (d Door) getY() int {
	return d.y
}
func (d Door) getRect() *sdl.Rect {
	return d.srcRect
}

type mainCharacter struct {
	entity
}

func createMainCharacter(obj game.Entity) *mainCharacter {
	return &mainCharacter{entity{obj.X, obj.Y, &textureIndex[obj.Tile][0]}}
}

func (mc mainCharacter) getX() int {
	return mc.x
}
func (mc mainCharacter) getY() int {
	return mc.y
}
func (mc mainCharacter) getRect() *sdl.Rect {
	return mc.srcRect
}

type chest struct {
	entity
	isOpen     bool
	srcStorage []*sdl.Rect
}

func createChest(obj game.Entity) chest {
	if obj.Tile == game.ChestC {
		return chest{entity{obj.X, obj.Y, &textureIndex[obj.Tile][0]}, false, chestSrcStorage}
	} else if obj.Tile == game.ChestO {
		return chest{entity{obj.X, obj.Y, &textureIndex[obj.Tile][0]}, true, chestSrcStorage}
	}
	panic("panic")
}

func (c chest) getX() int {
	return c.x
}
func (c chest) getY() int {
	return c.y
}
func (c chest) getRect() *sdl.Rect {
	return c.srcRect
}

type monster struct {
	entity
}

type enemy struct {
	entity
}

func createMonster(obj game.Entity) monster {
	return monster{entity{obj.X, obj.Y, &textureIndex[obj.Tile][0]}}
}

func (m monster) getX() int {
	return m.x
}
func (m monster) getY() int {
	return m.y
}
func (m monster) getRect() *sdl.Rect {
	return m.srcRect
}

func createEnemy(obj game.Entity) enemy {
	return enemy{entity{obj.X, obj.Y, &textureIndex[obj.Tile][0]}}
}
func (e enemy) getX() int {
	return e.x
}
func (e enemy) getY() int {
	return e.y
}
func (e enemy) getRect() *sdl.Rect {
	return e.srcRect
}
