package game

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type pos struct {
	x, y int
}

type Direction int

const (
	UP    Direction = 0
	RIGHT Direction = 1
	DOWN  Direction = 2
	LEFT  Direction = 3
	NONE  Direction = 4
)

type Cell struct {
	pos
	dir      Direction
	isFilled bool
}

type Path struct {
	Cells []*Cell
}

// EmptyMaze fills the maze with Blank tiles
func EmptyMaze(maze [][]Cell) {
	for y, row := range maze {
		for x := range row {
			row[x] = Cell{pos{x, y}, NONE, false}
		}
	}
}

// getEmptyCell returns the first empty cell's pointer in the maze and a bool value which is true if there is an empty cell in the maze
func getEmptyCellFromMaze(maze [][]Cell) (*Cell, bool) {
	for _, row := range maze {
		for x := range row {
			if !row[x].isFilled {
				emptyCell := Cell{pos{row[x].x, row[x].y}, row[x].dir, row[x].isFilled}
				return &emptyCell, true
			}
		}
	}

	return &Cell{}, false // means that all cells are consumed so algorithm should stop
}

//doesPathIncludes checks whether the given path includes the given cell or not
func doesPathIncludes(path []*Cell, cellToCheck Cell) bool {
	for _, cell := range path {
		if cell.x == cellToCheck.x && cell.y == cellToCheck.y {
			return true
		}
	}
	return false
}

func printMaze(maze [][]Cell) {
	for _, row := range maze {
		for i := range row {
			fmt.Print(row[i], " ")
		}
		fmt.Println()
	}
}

// randomWalkOnMaze runs on the maze and returns success:true if walk can be done
func randomWalkOnMaze(maze [][]Cell, width, height int) bool {
	startCell, success := getEmptyCellFromMaze(maze)
	if !success {
		return false // means that there is no empty cell in the maze
	}

	path := Path{}
	path.Cells = make([]*Cell, 0)
	path.Cells = append(path.Cells, startCell)

	rand.Seed(time.Now().UnixNano())

	// a loop for random walk
	for {
		// pick a random direction to walk
		dirNum := rand.Intn(4) // 0,1,2,3 randomly since 4(NONE) is not a 'dir', we don't choose it
		nextCell := &Cell{pos{startCell.x, startCell.y}, NONE, false}

		switch dirNum {
		case 0:
			if startCell.y != 0 {
				nextCell.y--
				startCell.dir = UP
			} else {
				continue
			}
		case 1:
			if startCell.x != width-1 {
				nextCell.x++
				startCell.dir = RIGHT
			} else {
				continue
			}
		case 2:
			if startCell.y != height-1 {
				nextCell.y++
				startCell.dir = DOWN
			} else {
				continue
			}
		case 3:
			if startCell.x != 0 {
				nextCell.x--
				startCell.dir = LEFT
			} else {
				continue
			}
		}

		//means that it intersects with the path so erase the loop
		if doesPathIncludes(path.Cells, *nextCell) {

			for {

				if path.Cells[len(path.Cells)-1].x == nextCell.x && path.Cells[len(path.Cells)-1].y == nextCell.y {
					break
				}

				//path.Cells[len(path.Cells)-1].isFilled = false
				//path.Cells[len(path.Cells)-1].dir = NONE
				path.Cells = path.Cells[:len(path.Cells)-1] // delete the last element from the array
				//fmt.Println(len(path.Cells))

				if len(path.Cells) == 0 {
					break
				}
			}
			startCell = path.Cells[len(path.Cells)-1]
			continue
		}

		// if the path intersects with a cell which is in the maze, use the path to fill the maze
		if maze[nextCell.y][nextCell.x].isFilled {
			for _, cell := range path.Cells {
				maze[cell.y][cell.x].isFilled = true
				maze[cell.y][cell.x].dir = cell.dir
			}

			maze[nextCell.y][nextCell.x].isFilled = true

			startCell, success = getEmptyCellFromMaze(maze)
			if !success {
				return false // means that there is no empty cell in the maze
			}

			path = Path{}
			path.Cells = make([]*Cell, 0)
			path.Cells = append(path.Cells, startCell)

			continue
		}

		path.Cells = append(path.Cells, &Cell{nextCell.pos, nextCell.dir, nextCell.isFilled})

		// startCell.x = nextCell.x
		// startCell.y = nextCell.y
		startCell = nextCell

		//startCell.isFilled = nextCell.isFilled
		//startCell.dir = nextCell.dir

	}

}

//GenerateRandomMaze will produce a random maze with given width and height by using Cell struct
func GenerateRandomMaze(width, height int) {

	maze := make([][]Cell, height)
	for i := range maze {
		maze[i] = make([]Cell, width)
	}

	EmptyMaze(maze)

	//set initial cell as pos(0,0)
	maze[0][0].isFilled = true

	randomWalkOnMaze(maze, width, height)

	//printMaze(maze)

	generateMapFromMaze(width, height, maze)

}

func generateMapFromMaze(width, height int, maze [][]Cell) [][]Tile {
	var Map [][]Tile
	Map = make([][]Tile, 2*height+1)
	for i := range Map {
		Map[i] = make([]Tile, 2*width+1)
	}

	for y, row := range maze {
		for x, cell := range row {
			Map[2*y+1][2*x+1] = DirtFloor

			switch cell.dir {
			case UP:
				Map[2*y+1-1][2*x+1] = DirtFloor
			case RIGHT:
				Map[2*y+1][2*x+1+1] = DirtFloor
			case DOWN:
				Map[2*y+1+1][2*x+1] = DirtFloor
			case LEFT:
				Map[2*y+1][2*x+1-1] = DirtFloor
			case NONE:

				// in left but directing to right
				if x != 0 {
					switch row[x-1].dir {
					case RIGHT:
						Map[2*y+1][2*x+1+1] = DirtFloor
						Map[2*y+1][2*x+1-1] = DirtFloor
						Map[2*y+1-1][2*x+1] = DirtFloor
						Map[2*y+1+1][2*x+1] = DirtFloor
					}
				}

				if x != len(row)-1 {
					// in right but directing to left
					switch row[x+1].dir {
					case LEFT:
						Map[2*y+1][2*x+1+1] = DirtFloor
						Map[2*y+1][2*x+1-1] = DirtFloor
						Map[2*y+1-1][2*x+1] = DirtFloor
						Map[2*y+1+1][2*x+1] = DirtFloor
					}
				}

				if y != 0 {
					switch maze[y-1][x].dir {
					case DOWN:
						Map[2*y+1][2*x+1+1] = DirtFloor
						Map[2*y+1][2*x+1-1] = DirtFloor
						Map[2*y+1-1][2*x+1] = DirtFloor
						Map[2*y+1+1][2*x+1] = DirtFloor
					}
				}

				if y != len(maze)-1 {
					switch maze[y+1][x].dir {
					case UP:
						Map[2*y+1][2*x+1+1] = DirtFloor
						Map[2*y+1][2*x+1-1] = DirtFloor
						Map[2*y+1-1][2*x+1] = DirtFloor
						Map[2*y+1+1][2*x+1] = DirtFloor
					}
				}

			}
		}
	}

	for y, row := range Map {
		for x := range row {
			if x == 0 || y == 0 || x == 2*width || y == 2*height || Map[y][x] == 0 {
				Map[y][x] = '#'
			}
		}
	}

OUTER:
	for y, row := range Map {
		for x := range row {
			if Map[y][x] == '.' {
				Map[y][x] = 'P'
				break OUTER
			}
		}
	}

	saveMap(Map, "new")

	return Map
}

func walkFromPlayerToACornerAndPutADoor(mp [][]Tile, winWidth, winHeight int) {

}

func saveMap(Map [][]Tile, levelName string) {
	file, err := os.OpenFile("./game/maps/"+levelName+".map", os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Truncate(0)
	file.Seek(0, 0)
	for y := range Map {
		for x := range Map[y] {
			if Map[y][x] != 0 {
				file.WriteString(Map[y][x].ToString())
			} else {
				file.WriteString(" ")
			}
		}
		file.WriteString("\n")
	}

}

func returnDirtOrStone() Tile {
	rand.Seed(time.Now().UnixNano())
	dirNum := rand.Intn(2) // 0 or 1

	if dirNum == 0 {
		return StoneWall
	}
	return DirtFloor
}
