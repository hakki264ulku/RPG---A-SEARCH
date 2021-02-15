package game

import (
	"fmt"
)

type Enemy struct {
	Character
}

func NewRat(p Pos) *Enemy {
	monster := Enemy{}
	monster.Pos = p
	monster.Tile = 'R'
	monster.Name = "Rat"
	monster.Hitpoints = 60
	monster.FullHitpoints = 60
	monster.Strength = 10
	monster.Speed = 1.0
	monster.ActionPoints = 0.0
	return &monster
}

func NewSpider(p Pos) *Enemy {
	monster := Enemy{}
	monster.Pos = p
	monster.Tile = 'S'
	monster.Name = "Spider"
	monster.Hitpoints = 90
	monster.FullHitpoints = 90
	monster.Strength = 15
	monster.Speed = 1.0
	monster.ActionPoints = 0.0
	return &monster
}

// HealthPot is left by dead monsters
type HealthPot struct {
	Hitpoint int
	Pos
}

func (m *Enemy) Update(level *Level2) {
	m.ActionPoints = m.Speed // they can move only for the amount of their speed
	playerPos := level.Player.Pos

	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)

	if positions != nil {
		if !ContainsEnemy(level.EnemiesForHealthBars, m) {
			level.EnemiesForHealthBars = append(level.EnemiesForHealthBars, m)
		}
	}

	// Must be >1 because the 1st position is the monsters current
	moveIndex := 1

	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Enemy) Move(to Pos, level *Level2) {
	_, exists := level.Enemies[to]

	if !exists && to != level.Player.Pos {
		delete(level.Enemies, m.Pos)
		level.Enemies[to] = m
		m.Pos = to
	} else {
		if to.X == level.Player.X && to.Y == level.Player.Y {
			Attack(m, level.Player)
		}
		if m.Hitpoints <= 0 {
			fmt.Println("Monster is dead")
			level.EnemiesForHealthBars = RemoveEnemyFromHealthArray(level.EnemiesForHealthBars, m)
			pos := m.Pos
			delete(level.Enemies, m.Pos)
			newPot := &HealthPot{10, pos}
			level.HealthPots[pos] = newPot
			fmt.Println(len(level.HealthPots))
		}
	}

}
