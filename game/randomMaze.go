package game

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

var gridWorld [100][100]Tile

func printGridWorld() {
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			fmt.Println(gridWorld[y][x].ToString())
		}
	}
}

func saveGridWorld(levelName string) {
	file, err := os.OpenFile("./game/maps/"+levelName+".map", os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Truncate(0)
	file.Seek(0, 0)

	for y := range gridWorld {
		for x := range gridWorld[y] {
			if gridWorld[y][x] != Blank {
				file.WriteString(gridWorld[y][x].ToString())
			} else {
				file.WriteString(" ")
			}
		}
		file.WriteString("\n")
	}
}

type pair struct {
	x, y int
}

func emptyGridWorld() {
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			gridWorld[y][x] = Blank
		}
	}
}

func createBackground(levelName string) {
	emptyGridWorld()
	var stack []pair
	startX := 1
	startY := 1
	stack = append(stack, pair{startX, startY})

	for len(stack) > 0 {
		index := len(stack) - 1
		currLoc := stack[index]
		stack = stack[:index]

		if currLoc.x == 0 || currLoc.y == 0 || currLoc.x == 99 || currLoc.y == 99 {
			gridWorld[currLoc.y][currLoc.x] = StoneWall
		}
		if gridWorld[currLoc.y][currLoc.x] != Blank {
			continue
		}

		gridWorld[currLoc.y][currLoc.x] = DirtFloor
		var blankLocs []pair
		blankLocs = blankLocs[0:0]
		if currLoc.x+1 < 100 && gridWorld[currLoc.y][currLoc.x+1] == Blank {
			blankLocs = append(blankLocs, pair{currLoc.x + 1, currLoc.y})
		}
		if currLoc.x-1 >= 0 && gridWorld[currLoc.y][currLoc.x-1] == Blank {
			blankLocs = append(blankLocs, pair{currLoc.x - 1, currLoc.y})
		}
		if currLoc.y+1 < 100 && gridWorld[currLoc.y+1][currLoc.x] == Blank {
			blankLocs = append(blankLocs, pair{currLoc.x, currLoc.y + 1})
		}
		if currLoc.y-1 >= 0 && gridWorld[currLoc.y-1][currLoc.x] == Blank {
			blankLocs = append(blankLocs, pair{currLoc.x, currLoc.y - 1})
		}

		if len(blankLocs) == 0 {
			continue
		}

		for _, v := range rand.Perm(len(blankLocs)) {
			newLoc := blankLocs[v]
			isGood := true
			if newLoc.x+1 < 100 && newLoc.x+1 != currLoc.x && gridWorld[newLoc.y][newLoc.x+1] == DirtFloor {
				isGood = false
			}
			if newLoc.x-1 >= 0 && newLoc.x-1 != currLoc.x && gridWorld[newLoc.y][newLoc.x-1] == DirtFloor {
				isGood = false
			}
			if newLoc.y+1 < 100 && newLoc.y+1 != currLoc.y && gridWorld[newLoc.y+1][newLoc.x] == DirtFloor {
				isGood = false
			}
			if newLoc.y-1 >= 0 && newLoc.y-1 != currLoc.y && gridWorld[newLoc.y-1][newLoc.x] == DirtFloor {
				isGood = false
			}

			if isGood {
				stack = append(stack, newLoc)
			} else {
				gridWorld[newLoc.y][newLoc.x] = StoneWall
			}
		}

	}
	saveGridWorld(levelName)
}

func isIntersection(currLoc pair) []pair {
	var dirtLocs []pair
	if gridWorld[currLoc.y][currLoc.x+1] == DirtFloor && !isTravelled[currLoc.y][currLoc.x+1] {
		dirtLocs = append(dirtLocs, pair{currLoc.x + 1, currLoc.y})
	}
	if gridWorld[currLoc.y][currLoc.x-1] == DirtFloor && !isTravelled[currLoc.y][currLoc.x-1] {
		dirtLocs = append(dirtLocs, pair{currLoc.x - 1, currLoc.y})
	}
	if gridWorld[currLoc.y+1][currLoc.x] == DirtFloor && !isTravelled[currLoc.y+1][currLoc.x] {
		dirtLocs = append(dirtLocs, pair{currLoc.x, currLoc.y + 1})
	}
	if gridWorld[currLoc.y-1][currLoc.x] == DirtFloor && !isTravelled[currLoc.y-1][currLoc.x] {
		dirtLocs = append(dirtLocs, pair{currLoc.x, currLoc.y - 1})
	}
	return dirtLocs

}

func depthCount(parentLoc, currLoc pair) (int, []pair) {
	i := 1
	s := 0
	var locSlice []pair
	locSlice = append(locSlice, currLoc)
	var nextLoc pair
	if gridWorld[currLoc.y][currLoc.x+1] == DirtFloor && !(currLoc.y == parentLoc.y && currLoc.x+1 == parentLoc.x) {
		s++
		nextLoc = pair{currLoc.x + 1, currLoc.y}
	}
	if gridWorld[currLoc.y][currLoc.x-1] == DirtFloor && !(currLoc.y == parentLoc.y && currLoc.x-1 == parentLoc.x) {
		s++
		nextLoc = pair{currLoc.x - 1, currLoc.y}
	}
	if gridWorld[currLoc.y+1][currLoc.x] == DirtFloor && !(currLoc.y+1 == parentLoc.y && currLoc.x == parentLoc.x) {
		s++
		nextLoc = pair{currLoc.x, currLoc.y + 1}
	}
	if gridWorld[currLoc.y-1][currLoc.x] == DirtFloor && !(currLoc.y-1 == parentLoc.y && currLoc.x == parentLoc.x) {
		s++
		nextLoc = pair{currLoc.x, currLoc.y - 1}
	}

	if s == 1 {
		newI, newLocSlice := depthCount(currLoc, nextLoc)
		i = i + newI
		locSlice = append(locSlice, newLocSlice...)
	}

	return i, locSlice
}

var isTravelled [100][100]bool

func createEntities(levelName string) {
	file2, err := os.OpenFile("./game/maps/"+levelName+"Entities.map", os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file2.Close()
	file2.Truncate(0)
	file2.Seek(0, 0)

	var stack []pair
	startX := 1
	startY := 1
	file2.WriteString("P " + strconv.Itoa(startX*32) + "," + strconv.Itoa(startY*32) + "\n")
	stack = append(stack, pair{startX, startY})

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			isTravelled[y][x] = false
		}
	}

	for len(stack) > 0 {
		index := len(stack) - 1
		currLoc := stack[index]
		stack = stack[:index]

		isTravelled[currLoc.y][currLoc.x] = true

		dirtLocs := isIntersection(currLoc)
		if len(dirtLocs) > 1 {
			for _, loc := range dirtLocs {
				dC, locSlice := depthCount(currLoc, loc)
				if dC > 2 {
					file2.WriteString("| " + strconv.Itoa(locSlice[0].x*32) + "," + strconv.Itoa(locSlice[0].y*32) + "\n")
					file2.WriteString("m " + strconv.Itoa(locSlice[1].x*32) + "," + strconv.Itoa(locSlice[1].y*32) + "\n")
					/*if rand.Intn(2) == 1 {
						file2.WriteString("C " + strconv.Itoa(locSlice[2].x*32) + "," + strconv.Itoa(locSlice[2].y*32) + "\n")
					} else {
						file2.WriteString("c " + strconv.Itoa(locSlice[2].x*32) + "," + strconv.Itoa(locSlice[2].y*32) + "\n")
					}*/
					file2.WriteString("C " + strconv.Itoa(locSlice[2].x*32) + "," + strconv.Itoa(locSlice[2].y*32) + "\n")
				}
			}
		}
		for _, loc := range dirtLocs {
			stack = append(stack, loc)
		}
	}
}

func CreateRandomMaze(levelName string, ui GameUI) Level {
	createBackground(levelName)
	createEntities(levelName)
	level := Level{}
	level.LevelName = levelName
	level.loadLevelFromFile()
	return level
}
