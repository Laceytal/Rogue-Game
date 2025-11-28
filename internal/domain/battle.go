package domain

import (
	"math/rand"
)

type turn int

const (
	playerTurn turn = iota
	MonsterTurn
)

const (
	initialHitChance   = 70.0
	standartAgility    = 50
	agilityFactor      = 0.3
	initialDamage      = 30.0
	standartStrength   = 50
	strengthFactor     = 0.3
	strengthAddition   = 65
	sleepChance        = 15
	maxHpPart          = 10
	lootAgilityFactor  = 0.2
	lootHpFactor       = 0.5
	lootStrengthFactor = 0.5
	MaximumFights      = 8
)

type BattleInfo struct {
	IsFight            bool
	Enemy              *Monster
	VampireFirstAttack bool
	OgreCooldown       bool
	PlayerAsleep       bool
}

type damageFormulasFunc func(*BattleInfo) float64

func Attack(player *Player, battleInfo *BattleInfo, currTurn turn, adapter *DifficultyAdapter) {
	switch currTurn {
	case playerTurn:
		if checkHit(player, battleInfo.Enemy, playerTurn) {
			damage := calculateDamage(player, battleInfo, playerTurn)
			battleInfo.Enemy.BaseStats.Health -= damage

			if adapter != nil {
				adapter.TotalDamageDealт += damage
			}
		}

		if battleInfo.Enemy.BaseStats.Health <= 0 {
			player.Backpack.Treasures.Value += int(calculateLoot(battleInfo.Enemy))

			if adapter != nil {
				adapter.TotalFightsWon++
			}
		}

	case MonsterTurn:
		if checkHit(player, battleInfo.Enemy, MonsterTurn) {
			damage := calculateDamage(player, battleInfo, MonsterTurn)
			player.BaseStats.Health -= damage

			if adapter != nil {
				adapter.TotalDamageReceived += damage
			}
		}
	}
}

func checkHit(player *Player, monster *Monster, currTurn turn) bool {
	wasHit := false
	chance := initialHitChance

	switch currTurn {
	case playerTurn:
		chance += hitChanceFormula(player.BaseStats.Agility, monster.BaseStats.Agility)

	case MonsterTurn:
		chance += hitChanceFormula(monster.BaseStats.Agility, player.BaseStats.Agility)
	}

	if rand.Intn(100) < int(chance) || monster.Type == Ogre {
		wasHit = true
	}
	return wasHit
}

func calculateDamage(player *Player, battleInfo *BattleInfo, currTurn turn) float64 {
	damage := initialDamage
	monsterDamageFormulas := map[MonsterType]damageFormulasFunc{
		Zombie:  zombieGhostMimicDamageFormula,
		Vampire: nil,
		Ghost:   zombieGhostMimicDamageFormula,
		Ogre:    ogreDamageFormula,
		Snake:   snakeDamageFormula,
		Mimic:   zombieGhostMimicDamageFormula,
	}

	switch currTurn {
	case playerTurn:
		if !(battleInfo.Enemy.Type == Vampire && battleInfo.VampireFirstAttack) && !(battleInfo.Enemy.Type == Snake && battleInfo.PlayerAsleep) {
			if player.Weapon.Strength == noWeapon {
				damage += float64(player.BaseStats.Strength-standartStrength) * strengthFactor
			} else {
				damage = float64(player.Weapon.Strength) * float64(player.BaseStats.Strength+strengthAddition) / 100
			}
		} else if battleInfo.Enemy.Type == Vampire && battleInfo.VampireFirstAttack {
			battleInfo.VampireFirstAttack = false
		} else {
			battleInfo.PlayerAsleep = false
		}

	case MonsterTurn:
		if battleInfo.Enemy.Type == Vampire {
			damage = vampireDamageFormula(player)
		} else {
			damage = monsterDamageFormulas[battleInfo.Enemy.Type](battleInfo)
		}
	}
	return damage
}

func calculateLoot(monster *Monster) uint {
	loot := uint(float64(monster.BaseStats.Agility)*lootAgilityFactor) +
		uint(monster.BaseStats.Health*lootHpFactor) +
		uint(float64(monster.BaseStats.Strength)*lootStrengthFactor) +
		uint(rand.Intn(20))
	return loot
}

func deleteMonsterInfo(room *Room, monster *Monster) {
	currentPos := 0
	for i := range room.MonsterNum {
		room.Monsters[currentPos] = room.Monsters[i]
		if !CheckEqualCoords(monster.BaseStats.Coords.Coordinates, room.Monsters[i].BaseStats.Coords.Coordinates) {
			currentPos++
		}
	}
	room.MonsterNum--
}

func UpdateFightStatus(playerCoordinates *Object, level *Level, battlesArray []BattleInfo) {
	for i := range RoomsNum {
		for j := range level.Rooms[i].MonsterNum {
			if checkContact(playerCoordinates, &level.Rooms[i].Monsters[j]) && CheckUnique(&level.Rooms[i].Monsters[j], battlesArray) {
				initBattle(&level.Rooms[i].Monsters[j], battlesArray)
			}
		}
	}

	for i := range MaximumFights {
		if battlesArray[i].IsFight && (!checkContact(playerCoordinates, battlesArray[i].Enemy) || battlesArray[i].Enemy.BaseStats.Health <= 0) {
			battlesArray[i].IsFight = false
		}
	}
}

func initBattle(monster *Monster, battlesArray []BattleInfo) {
	for i := range MaximumFights {
		if !battlesArray[i].IsFight {
			battlesArray[i].IsFight = true
			battlesArray[i].OgreCooldown = false
			battlesArray[i].VampireFirstAttack = true
			battlesArray[i].PlayerAsleep = false
			battlesArray[i].Enemy = monster
			break
		}
	}
}

func checkContact(playerCoordinates *Object, monster *Monster) bool {
	isContact := checkIfNeighborTile(playerCoordinates.Coordinates, monster.BaseStats.Coords.Coordinates)
	if !isContact {
		isContact = (monster.Type == Snake && checkIfDiagonallyNeighbourTile(playerCoordinates.Coordinates, monster.BaseStats.Coords.Coordinates))
	}

	return isContact
}

func CheckPlayerAttack(player *Player, battle *BattleInfo, playerChosenDirection Directions, adapter *DifficultyAdapter) bool {
	playerIsAttacking := false
	oldCoords := player.BaseStats.Coords
	moveCharacterByDirection(playerChosenDirection, &oldCoords)
	if CheckEqualCoords(oldCoords.Coordinates, battle.Enemy.BaseStats.Coords.Coordinates) {
		Attack(player, battle, playerTurn, adapter)
		playerIsAttacking = true
	}

	return playerIsAttacking
}

func RemoveDeadMonsters(level *Level) {
	for i := range RoomsNum {
		for j := range level.Rooms[i].MonsterNum {
			if level.Rooms[i].Monsters[j].BaseStats.Health <= 0 {
				deleteMonsterInfo(&level.Rooms[i], &level.Rooms[i].Monsters[j])
			}
		}
	}
}

func CheckEqualCoords(firstCoords, secondCoords Coordinates) bool {
	return firstCoords[X] == secondCoords[X] &&
		firstCoords[Y] == secondCoords[Y]
}

func checkIfNeighborTile(firstCoords Coordinates, secondCoords Coordinates) bool {
	return (firstCoords[X] == secondCoords[X] && Abs(firstCoords[Y]-secondCoords[Y]) == 1) ||
		(firstCoords[Y] == secondCoords[Y] && Abs(firstCoords[X]-secondCoords[X]) == 1)
}

func checkIfDiagonallyNeighbourTile(firstCoords Coordinates, secondCoords Coordinates) bool {
	return Abs(firstCoords[X]-secondCoords[X]) == 1 && Abs(firstCoords[Y]-secondCoords[Y]) == 1
}

func CheckUnique(monster *Monster, battlesArray []BattleInfo) bool {
	isUnique := true
	for i := 0; i < MaximumFights && isUnique; i++ {
		if battlesArray[i].IsFight && CheckEqualCoords(battlesArray[i].Enemy.BaseStats.Coords.Coordinates, monster.BaseStats.Coords.Coordinates) {
			isUnique = false
		}
	}
	return isUnique
}

func hitChanceFormula(attackerAgility int, targetAgility int) float64 {
	return float64(attackerAgility-targetAgility-standartAgility) * agilityFactor
}

func vampireDamageFormula(player *Player) float64 {

	damage := float64(player.RegenLimit) / maxHpPart

	return damage
}

func zombieGhostMimicDamageFormula(battleInfo *BattleInfo) float64 {
	return initialDamage + float64(battleInfo.Enemy.BaseStats.Strength-standartStrength)*strengthFactor
}

func ogreDamageFormula(battleInfo *BattleInfo) float64 {
	damage := 0.0
	if !battleInfo.OgreCooldown {
		damage = float64(battleInfo.Enemy.BaseStats.Strength-standartStrength) * strengthFactor
		battleInfo.OgreCooldown = true
	} else {
		battleInfo.OgreCooldown = false
	}
	return damage
}

func snakeDamageFormula(battleInfo *BattleInfo) float64 {
	if rand.Intn(100) < sleepChance {
		battleInfo.PlayerAsleep = true
	}
	return zombieGhostMimicDamageFormula(battleInfo)
}
