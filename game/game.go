package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type GameUI interface {
	Init()
	Draw(*Level, bool) bool
	SelectLevel() (*Level, bool)
	AddPreview(Level)
}

type Tile rune

const (
	StoneWall     Tile = '#'
	DirtFloor     Tile = '.'
	DoorC         Tile = '|'
	DoorO         Tile = '/'
	MainCharacter Tile = 'P'
	Spider        Tile = 'S'
	Rat           Tile = 'R'
	ChestC        Tile = 'C'
	ChestO        Tile = 'c'
	Monster       Tile = 'm' // will be used as specific monster like rat or snake etc
	Blank         Tile = 0
	Pending       Tile = -1
)

type Grid struct {
	Layers     []Tile
	Background Tile
}

type Row struct {
	x, y  int
	Grids []Grid
}

type Pos struct {
	X, Y int
}

type GridWorld struct {
	Rows []Row
}

type Entity struct {
	Pos
	Tile Tile
}

type Attackable interface {
	GetActionPoints() float64
	SetActionPoints(float64)
	GetHitpoints() int
	SetHitpoints(int)
	GetAttackPower() int
}

//Character is an attackable by implementing the attackable functions
type Character struct {
	Entity
	Hitpoints     int
	FullHitpoints int
	Name          string
	Strength      int
	Speed         float64
	ActionPoints  float64
}

func (c *Character) GetActionPoints() float64 {
	return c.ActionPoints
}

func (c *Character) SetActionPoints(ap float64) {
	c.ActionPoints = ap
}

func (c *Character) GetHitpoints() int {
	return c.Hitpoints
}

func (c *Character) SetHitpoints(hp int) {
	c.Hitpoints = hp
}

func (c *Character) GetAttackPower() int {
	return c.Strength
}

// Attack -> two attackables like player and a monster are taken as args
func Attack(a1, a2 Attackable) {
	// fmt.Println(a1, "  attacks  ", a2)
	a1.SetActionPoints(a1.GetActionPoints() - 1)
	a2.SetHitpoints(a2.GetHitpoints() - a1.GetAttackPower())
	if a2.GetHitpoints() > 0 {
		a2.SetActionPoints(a2.GetActionPoints() - 1)
		a1.SetHitpoints(a1.GetHitpoints() - a2.GetAttackPower())
	}
}

type Player struct {
	Character
}

func (player *Player) Move(to Pos, level *Level2) {
	monster, exists := level.Enemies[to]
	pod, existsPod := level.HealthPots[to]

	if existsPod {
		fmt.Println(player.Hitpoints)
		player.Hitpoints += pod.Hitpoint
		delete(level.HealthPots, pod.Pos)
		fmt.Println(player.Hitpoints)
	}

	if !exists {
		player.Pos = to
	} else {
		Attack(player, monster)

		fmt.Println("Player attacked Monster")
		fmt.Println(level.Player.Hitpoints, monster.Hitpoints)
		if monster.Hitpoints <= 0 {
			pos := monster.Pos
			delete(level.Enemies, monster.Pos)
			newPot := &HealthPot{10, pos}
			level.HealthPots[pos] = newPot
			level.EnemiesForHealthBars = RemoveEnemyFromHealthArray(level.EnemiesForHealthBars, monster)
		}
		if level.Player.Hitpoints <= 0 {
			panic("You Died...")
		}
		if level.Player.Hitpoints <= 0 {
			fmt.Println("YOU DIED!!!")
			panic("YOU DIED!!!")
		}
	}
}

type Level struct {
	GridWorld GridWorld
	LevelName string
	Entities  []Entity
}

type Level2 struct {
	Map                  [][]Tile
	Player               *Player
	Enemies              map[Pos]*Enemy
	EnemiesForHealthBars []*Enemy
	HealthPots           map[Pos]*HealthPot
	FileName             string
	//Debug    map[Pos]bool
}

func LoadLevelFromFile2(fileName string) *Level2 {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	levelLines := make([]string, 0)
	longestRow := 0
	index := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if longestRow < len(levelLines[index]) {
			longestRow = len(levelLines[index])
		}
		index++
	}
	defer file.Close()

	level := &Level2{}
	level.Player = &Player{}

	// TODO where should we initialize the player?
	level.Player.ActionPoints = 0
	level.Player.Strength = 30
	level.Player.Hitpoints = 300
	level.Player.FullHitpoints = 300
	level.Player.Tile = 'R'
	level.Player.Speed = 1.0
	level.Player.Name = "PurpleWIZARD"

	level.Map = make([][]Tile, len(levelLines))
	level.EnemiesForHealthBars = make([]*Enemy, 0)
	level.HealthPots = make(map[Pos]*HealthPot, 0)
	level.Enemies = make(map[Pos]*Enemy)

	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]
		for x, c := range line {

			var t Tile
			switch c {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case '#':
				t = StoneWall
			case '.':
				t = DirtFloor
			case '|':
				t = DoorC
			case '-':
				t = DoorO
			case 'P':
				level.Player.X = x
				level.Player.Y = y
				t = Pending
			case 'R':
				level.Enemies[Pos{x, y}] = NewRat(Pos{x, y})
				t = Pending
			case 'S':
				level.Enemies[Pos{x, y}] = NewSpider(Pos{x, y})
				t = Pending
			default:
				fmt.Printf("%c\n", c)
				panic("Invalid character in the map file")
			}
			level.Map[y][x] = t

		}
	}

	// If tile pending, it assigns the background of that tile. E.g. when player is encountered on the map, it puts a dirt floor under that tile
	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
				level.Map[y][x] = level.bfsFloor(Pos{x, y})
			}
		}
	}
	level.FileName = fileName

	return level
}

func (level *Level) ToString() {
	gw := level.GridWorld
	for y := range gw.Rows {
		for x := range gw.Rows[y].Grids {
			fmt.Print(gw.Rows[y].Grids[x].Background.ToString())
		}
		fmt.Println("")
	}
}

//TODO check this func to be able to print the map i think... ?
func (tile Tile) ToString() string {
	switch tile {
	case StoneWall:
		return "#"
	case DirtFloor:
		return "."
	case DoorC:
		return "|"
	case DoorO:
		return "-"
	case Blank:
		return " "
	case MainCharacter:
		return "P"
	case ChestC:
		return "C"
	case ChestO:
		return "c"
	case Monster:
		return "m"
	case Rat:
		return "R"
	case Spider:
		return "S"
	default:
		panic("unknown toString tile")
	}
}

func createPreviews(ui GameUI) {
	for levelindex := 1; levelindex <= 4; levelindex++ {
		level := Level{}
		level.LevelName = "level" + strconv.Itoa(levelindex)
		level.loadLevelFromFile()
		ui.AddPreview(level)
	}
}

func Run(ui GameUI) {
	rand.Seed(time.Now().UnixNano())
	isReplayed := true
	for isReplayed {
		ui.Init()
		createPreviews(ui)
		var editBeforeStart bool
		level, editBeforeStart := ui.SelectLevel()
		if level == nil {
			return
		}
		level.loadLevelFromFile()
		isReplayed = ui.Draw(level, editBeforeStart)
	}
}
