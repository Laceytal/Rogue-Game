package domain

import (
	"math/rand"
)

const (
	ogreStep              = 2
	simpleDirections      = 4
	diagonalDirections    = 4
	allDirections         = 8
	simpleToDiagonalShift = 4
	maxTriesToMove        = 16
)

type npcMovementFunc func(monster *Monster, level *Level) *Vector

func CharacterOutsideBorder(characterCoords *Object, room *Object) bool {
	charX := characterCoords.Coordinates[X]
	charY := characterCoords.Coordinates[Y]
	charWidth := characterCoords.Size[X]
	charHeight := characterCoords.Size[Y]

	roomX := room.Coordinates[X]
	roomY := room.Coordinates[Y]
	roomWidth := room.Size[X]
	roomHeight := room.Size[Y]

	leftBorder := charX < roomX+1
	rightBorder := charX+charWidth+1 > roomX+roomWidth
	topBorder := charY < roomY+1
	bottomBorder := charY+charHeight+1 > roomY+roomHeight

	return leftBorder || rightBorder || topBorder || bottomBorder
}

func moveCharacterByDirection(direction Directions, characterGeometry *Object) {
	x := X
	y := Y
	switch direction {
	case Forward:
		characterGeometry.Coordinates[y]--

	case Left:
		characterGeometry.Coordinates[x]--

	case Right:
		characterGeometry.Coordinates[x]++

	case Back:
		characterGeometry.Coordinates[y]++

	case diagonallyForwardLeft:
		characterGeometry.Coordinates[x]--
		characterGeometry.Coordinates[y]--

	case diagonallyForwardRight:
		characterGeometry.Coordinates[x]++
		characterGeometry.Coordinates[y]--

	case diagonallyBackLeft:
		characterGeometry.Coordinates[x]--
		characterGeometry.Coordinates[y]++

	case diagonallyBackRight:
		characterGeometry.Coordinates[x]++
		characterGeometry.Coordinates[y]++

	case stop:
	}
}

func moveCharacterByPath(path *Vector, characterGeometry *Object) {
	if path != nil {
		for _, direction := range path.Data {
			moveCharacterByDirection(direction, characterGeometry)
		}
	}

}

func MoveMonster(monster *Monster, playerCoordinates *Object, level *Level) {
	npcMovementFunctions := map[MonsterType]npcMovementFunc{
		Zombie:  patternZombie,
		Vampire: patternVampire,
		Ghost:   patternGhost,
		Ogre:    patternOgre,
		Snake:   patternSnake,
		Mimic:   patternMimic,
	}

	originalHostility := monster.hostility

	if monster.Type == Mimic && monster.IsChasing {
		monster.hostility = average
	}

	var path *Vector
	if isPlayerNear(playerCoordinates, monster, level) {
		path = distAndNextPosToTarget(&monster.BaseStats.Coords, playerCoordinates, level)
		if path != nil {
			path.Data = path.Data[:1]
		}
		monster.IsChasing = true
	}

	if monster.Type == Mimic && !monster.IsChasing {
		monster.hostility = originalHostility
	}

	if path == nil {
		path = npcMovementFunctions[monster.Type](monster, level)
	}

	coords := monster.BaseStats.Coords
	if path != nil && len(path.Data) > 0 {
		moveCharacterByPath(path, &coords)
		if isValidMonsterPosition(&coords, level, monster.IsChasing) && !CheckEqualCoords(coords.Coordinates, playerCoordinates.Coordinates) {
			moveCharacterByPath(path, &monster.BaseStats.Coords)
		}
		monster.dir = path.Data[len(path.Data)-1]
	}
	path = nil
}

func isValidMonsterPosition(coords *Object, level *Level, isChasing bool) bool {
	if !isChasing && isInCorridor(coords, level) {
		return false
	}

	if checkOutsideBorder(coords, level) {
		return false
	}

	if !checkUnoccupiedLevel(coords, level) {
		return false
	}

	return true
}

func isInCorridor(coords *Object, level *Level) bool {
	for i := range level.Passages.PassagesNum {
		if !CharacterOutsideBorder(coords, &level.Passages.Passages[i]) {
			return true
		}
	}
	return false
}

func patternMimic(monster *Monster, level *Level) *Vector {
	if !monster.IsChasing {
		return nil
	}

	path := createVector()
	for try := 0; try < maxTriesToMove && len(path.Data) == 0; try++ {
		currentCoords := monster.BaseStats.Coords
		currentDirection := Directions(rand.Intn(simpleDirections))
		moveCharacterByDirection(currentDirection, &currentCoords)
		if isValidMonsterPosition(&currentCoords, level, monster.IsChasing) {
			path.pushBack(currentDirection)
		}
	}
	return path
}

func patternZombie(monster *Monster, level *Level) *Vector {
	path := createVector()
	for try := 0; try < maxTriesToMove && len(path.Data) == 0; try++ {
		currentCoords := monster.BaseStats.Coords
		currentDirection := Directions(rand.Intn(simpleDirections))
		moveCharacterByDirection(currentDirection, &currentCoords)
		if isValidMonsterPosition(&currentCoords, level, monster.IsChasing) {
			path.pushBack(currentDirection)
		}
	}
	return path
}

func patternVampire(monster *Monster, level *Level) *Vector {
	path := createVector()
	for try := 0; try < maxTriesToMove && len(path.Data) == 0; try++ {
		currentCoords := monster.BaseStats.Coords
		currentDirection := Directions(rand.Intn(allDirections))
		moveCharacterByDirection(currentDirection, &currentCoords)
		if isValidMonsterPosition(&currentCoords, level, monster.IsChasing) {
			path.pushBack(currentDirection)
		}
	}
	return path
}

func patternGhost(monster *Monster, level *Level) *Vector {
	var path *Vector
	var room *Room
	for i := 0; i < RoomsNum && room == nil; i++ {
		if !CharacterOutsideBorder(&monster.BaseStats.Coords, &level.Rooms[i].Coords) {
			room = &level.Rooms[i]
		}
	}

	for try := 0; try < maxTriesToMove && path == nil && room != nil; try++ {
		currentCoords := Object{
			Size: [2]int{1, 1},
			Coordinates: [2]int{
				room.Coords.Coordinates[X] + rand.Intn(room.Coords.Size[X]-2) + 1,
				room.Coords.Coordinates[Y] + rand.Intn(room.Coords.Size[Y]-2) + 1,
			},
		}
		if isValidMonsterPosition(&currentCoords, level, monster.IsChasing) {
			path = distAndNextPosToTarget(&monster.BaseStats.Coords, &currentCoords, level)
		}
	}
	return path
}

func patternOgre(monster *Monster, level *Level) *Vector {
	path := createVector()
	for try := 0; try < maxTriesToMove && len(path.Data) == 0; try++ {
		currentCoords := monster.BaseStats.Coords
		currentDirection := Directions(rand.Intn(simpleDirections))

		move := true
		for range ogreStep {
			moveCharacterByDirection(currentDirection, &currentCoords)
			if !isValidMonsterPosition(&currentCoords, level, monster.IsChasing) {
				move = false
				break
			}
		}
		if move {
			for range ogreStep {
				path.pushBack(currentDirection)
			}
		}
	}
	return path
}

func patternSnake(monster *Monster, level *Level) *Vector {
	path := createVector()
	for try := 0; try < maxTriesToMove && len(path.Data) == 0; try++ {
		currentCoords := monster.BaseStats.Coords
		currentDirection := Directions(simpleToDiagonalShift + rand.Intn(diagonalDirections))
		moveCharacterByDirection(currentDirection, &currentCoords)
		if currentDirection == monster.dir {
			continue
		}
		if isValidMonsterPosition(&currentCoords, level, monster.IsChasing) {
			path.pushBack(currentDirection)
		}
	}

	if len(path.Data) == 0 {
		currentCoords := monster.BaseStats.Coords
		moveCharacterByDirection(monster.dir, &currentCoords)
		if isValidMonsterPosition(&currentCoords, level, monster.IsChasing) {
			path.pushBack(monster.dir)
		}
	}

	return path
}

func isPlayerNear(playerCoordinates *Object, monster *Monster, level *Level) bool {

	if !isPlayerInRoom(playerCoordinates, level) && !monster.IsChasing {
		return false
	}
	dist := Abs(playerCoordinates.Coordinates[X] - monster.BaseStats.Coords.Coordinates[X])
	dist += Abs(playerCoordinates.Coordinates[Y] - monster.BaseStats.Coords.Coordinates[Y])

	playerNear := false
	switch monster.hostility {
	case low:
		if dist <= lowHostilityRadius {
			playerNear = true
		}

	case average:
		if dist <= averageHostilityRadius {
			playerNear = true
		}

	case high:
		if dist <= highHostilityRadius {
			playerNear = true
		}
	}
	return playerNear
}

func isPlayerInRoom(playerCoords *Object, level *Level) bool {
	for i := range RoomsNum {
		if !CharacterOutsideBorder(playerCoords, &level.Rooms[i].Coords) {
			return true
		}
	}
	return false
}

func MovePlayer(player *Player, level *Level, chosenDirection Directions) {

	currentCoords := [2]int{player.BaseStats.Coords.Coordinates[X], player.BaseStats.Coords.Coordinates[Y]}
	switch chosenDirection {
	case Forward:
		player.BaseStats.Coords.Coordinates[Y]--
	case Back:
		player.BaseStats.Coords.Coordinates[Y]++
	case Right:
		player.BaseStats.Coords.Coordinates[X]++
	case Left:
		player.BaseStats.Coords.Coordinates[X]--
	default:
	}

	if checkOutsideBorder(&player.BaseStats.Coords, level) {
		player.BaseStats.Coords.Coordinates[X] = currentCoords[X]
		player.BaseStats.Coords.Coordinates[Y] = currentCoords[Y]
	}

}

func checkOutsideBorder(playerCoordinates *Object, level *Level) bool {

	for _, room := range level.Rooms {
		if !CharacterOutsideBorder(playerCoordinates, &room.Coords) {
			return false
		}
	}

	for i := range level.Passages.PassagesNum {
		if !CharacterOutsideBorder(playerCoordinates, &level.Passages.Passages[i]) {
			return false
		}
	}

	return true
}

func FindCurrentRoom(playerCoords *Object, level *Level) *Room {
	for i := range RoomsNum {
		if !CharacterOutsideBorder(playerCoords, &level.Rooms[i].Coords) {
			return &level.Rooms[i]
		}
	}
	return nil
}
