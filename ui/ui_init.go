package ui

import (
	"strconv"

	"github.com/BiesGo/sdlWorkSpace/rpg/game"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const winWidth, winHeight = 800, 600 //1280, 720 //800, 600

type stateFunc func(*UI2d) stateFunc

type gameState int

const (
	mainScreen   gameState = 0
	selectScreen gameState = 1
	inGame       gameState = 2
	editLevel    gameState = 3
	endScreen    gameState = 4
	playScreen   gameState = 5
)

var currentState gameState = mainScreen

var renderer *sdl.Renderer
var textureAtlas *sdl.Texture
var textureIndex map[game.Tile][]sdl.Rect
var blackPixel *sdl.Texture
var font *ttf.Font

var healthBarTextures []*sdl.Texture
var potTexture *sdl.Texture

//var enemiesForHealthBars []*game.Enemy

// TODO -> TO PLAY THE GAME FOR NOW, but should be changed with a smarter way
var GlobalLevel2 *game.Level2
var centerX int
var centerY int

type inputState struct {
	leftButton      bool
	prevLeftButton  bool
	rightButton     bool
	prevRightButton bool
	x, y            int
	currKeyState    []uint8
	prevKeyState    []uint8
}

type layer struct {
	srcRect  [100][100]*sdl.Rect
	dstRect  [100][100]*sdl.Rect
	entities []interface{}
}

type UI2d struct {
	levelPreviews []layer
	background    layer
	mc            *mainCharacter
	input         *inputState
	mainMenu      mainMenuObj
	selectMenu    selectMenuObj
	editMenu      editMenuObj
	endMenu       endMenuObj
}

func (ui *UI2d) Init() {
	createMainMenu(ui)
	var input inputState
	input.updateMouseState()
	input.currKeyState = sdl.GetKeyboardState()
	input.prevKeyState = make([]uint8, len(input.currKeyState))
	input.updateKeyboardState()
	ui.input = &input
	createSelectMenu(ui)
	createEditMenu(ui)
	endMenuInit(ui)
}

func init() {
	healthBarTextures = make([]*sdl.Texture, 15)
	ttf.Init()
	font, _ = ttf.OpenFont("ui/assets/OpenSans-Regular.ttf", 64)

	//sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	window, err := sdl.CreateWindow("RPG", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	/*explosionBytes, audioSpec := sdl.LoadWAV("29301__junggle__btn121.wav")
	audioID, err := sdl.OpenAudioDevice("", false, audioSpec, nil, 0)
	if err != nil {
		panic(err)
	}
	defer sdl.FreeWAV(explosionBytes)*/
	// rand.Seed(time.Now().UTC().UnixNano())
	blackPixel = createOnePixel(0, 0, 0, 0)
	textureAtlas = imgFileToTexture("ui/assets/tiles.png")
	loadTextureIndex()

	GlobalLevel2 = game.LoadLevelFromFile2("game/maps/new.map")
	centerX = -1
	centerY = -1
	potTexture = imgFileToTexture("ui/assets/healthBars/potion.png")

	for i := 0; i < len(healthBarTextures); i++ {
		healthBarTextures[i] = imgFileToTexture("ui/assets/healthBars/bar" + strconv.Itoa(i+1) + ".png")
	}

}
