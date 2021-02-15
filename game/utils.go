package game

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func inRange(level *Level2, pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func checkDoor(level *Level2, pos Pos) {
	t := level.Map[pos.Y][pos.X]

	if t == DoorC {
		level.Map[pos.Y][pos.X] = DoorO
	}
}

// HandleInput returns true if input is taken in that specific loop cycle, and returns false if no input taken from the user
// for the sake of TURN BASES ABILITY
func HandleInput(currKeyState, prevKeyState []uint8, level *Level2) bool {
	p := level.Player

	if currKeyState[sdl.SCANCODE_UP] != 0 && prevKeyState[sdl.SCANCODE_UP] == 0 {
		newPos := Pos{p.X, p.Y - 1}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}
		return true
	}

	if currKeyState[sdl.SCANCODE_DOWN] != 0 && prevKeyState[sdl.SCANCODE_DOWN] == 0 {
		newPos := Pos{p.X, p.Y + 1}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}
		return true
	}

	if currKeyState[sdl.SCANCODE_LEFT] != 0 && prevKeyState[sdl.SCANCODE_LEFT] == 0 {
		newPos := Pos{p.X - 1, p.Y}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}
		return true
	}

	if currKeyState[sdl.SCANCODE_RIGHT] != 0 && prevKeyState[sdl.SCANCODE_RIGHT] == 0 {
		newPos := Pos{p.X + 1, p.Y}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}
		return true
	}

	return false
}

func canWalk(level *Level2, pos Pos) bool {

	if inRange(level, pos) {
		t := level.Map[pos.Y][pos.X]
		switch t {
		case StoneWall, DoorC, Blank:
			return false
		default:
			return true
		}
	}
	return false
}

func getNeighbours(level *Level2, pos Pos) []Pos {
	neighbours := make([]Pos, 0, 4)
	left := Pos{pos.X - 1, pos.Y}
	right := Pos{pos.X + 1, pos.Y}
	up := Pos{pos.X, pos.Y - 1}
	down := Pos{pos.X, pos.Y + 1}

	if canWalk(level, right) {
		neighbours = append(neighbours, right)
	}

	if canWalk(level, left) {
		neighbours = append(neighbours, left)
	}

	if canWalk(level, up) {
		neighbours = append(neighbours, up)
	}

	if canWalk(level, down) {
		neighbours = append(neighbours, down)
	}

	return neighbours
}

func (level *Level2) bfsFloor(start Pos) Tile {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true
	//level.Debug = visited

	for len(frontier) > 0 {
		current := frontier[0]

		currentTile := level.Map[current.Y][current.X]
		switch currentTile {
		case DirtFloor:
			return DirtFloor
		default:

		}

		frontier = frontier[1:] // pops the first appended element
		for _, next := range getNeighbours(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}

	}
	return DirtFloor
}

func (level *Level2) astar(start Pos, goal Pos) []Pos {
	frontier := make(pqueue, 0, 8)
	frontier = frontier.push(start, 1)
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	var current Pos
	for len(frontier) > 0 {
		frontier, current = frontier.pop()

		if current == goal {
			path := make([]Pos, 0)
			p := current
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}
			path = append(path, p)

			// reverse path to be able to follow it
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			return path // this path is from the start position to the goal position
		}

		for _, next := range getNeighbours(level, current) {
			newCost := costSoFar[current] + 1 // always one for now

			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = frontier.push(next, priority)
				cameFrom[next] = current
			}
		}
	}

	return nil
}

// ContainsEnemy checks whether the given []*Enemy array contains the enemy
func ContainsEnemy(arr []*Enemy, e *Enemy) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}

// FindIndexOfEnemy returns the index of the enemy from the health bar slice, returns -1 if not found any equality
func FindIndexOfEnemy(healthBarArr []*Enemy, e *Enemy) int {
	for i, en := range healthBarArr {
		if en == e {
			fmt.Println(i)
			return i
		}
	}
	return -1
}

// REMOVE THE DEAD MONSTER FROM THE HEALTH BAR ARRAY

func RemoveEnemyFromHealthArray(arr []*Enemy, e *Enemy) []*Enemy {
	// var returnArr []*Enemy
	// for _, m := range arr {
	// 	if m.Hitpoints > 0 {
	// 		returnArr = append(returnArr, m)
	// 	}
	// }
	i := FindIndexOfEnemy(arr, e)
	fmt.Println(i)
	arr[i] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}
