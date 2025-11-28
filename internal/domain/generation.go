package domain

import (
	"math/rand"
	"time"
)

type generateConsumableFunc func(room *Room, player *Player)

func clearData(level *Level) {
	for room := range RoomsNum {
		level.Rooms[room].MonsterNum = 0
		level.Rooms[room].Consumables.FoodNum = 0
		level.Rooms[room].Consumables.WeaponNum = 0
		level.Rooms[room].Consumables.ElixirNum = 0
		level.Rooms[room].Consumables.ScrollNum = 0
	}
}

func GenerateNextLevel(level *Level, player *Player) {
	clearData(level)
	level.LevelNum++
	generateRooms(&level.Rooms)
	generatePassages(&level.Passages, &level.Rooms)
	playerRoom := generatePlayer(&level.Rooms, player)
	generateMonsters(level, playerRoom)
	generateConsumables(level, playerRoom, player, level.LevelNum)
	generateExit(level, playerRoom)
}

func generateRooms(rooms *[RoomsNum]Room) {
	for i := range RoomsNum {
		widthRoom := getRandomInRange(minRoomWidth, maxRoomWidth)
		heightRoom := getRandomInRange(minRoomHeight, maxRoomHeight)

		leftRangeCoord := (i%RoomsInWidth)*RegionWidth + 1
		rightRangeCoord := (i%RoomsInWidth+1)*RegionWidth - widthRoom - 1
		xCoord := getRandomInRange(leftRangeCoord, rightRangeCoord)

		upRangeCoord := (i/RoomsInWidth)*RegionHeight + 1
		bottomRangeCoord := (i/RoomsInWidth+1)*RegionHeight - heightRoom - 1
		yCoord := getRandomInRange(upRangeCoord, bottomRangeCoord)

		rooms[i].Coords.Size[X] = widthRoom
		rooms[i].Coords.Size[Y] = heightRoom

		rooms[i].Coords.Coordinates[X] = xCoord
		rooms[i].Coords.Coordinates[Y] = yCoord
	}
}

func generateEdgesForRooms(edges []edge, countEdges *int) {
	*countEdges = 0

	for i := range RoomsInHeight {
		for j := 0; j+1 < RoomsInWidth; j++ {
			currentRoom := i*RoomsInHeight + j
			edges[*countEdges].u = currentRoom
			edges[*countEdges].v = currentRoom + 1
			(*countEdges)++
		}
	}

	for i := 0; i+1 < RoomsInHeight; i++ {
		for j := range RoomsInWidth {
			currentRoom := i*RoomsInHeight + j
			edges[*countEdges].u = currentRoom
			edges[*countEdges].v = currentRoom + RoomsInWidth
			(*countEdges)++
		}
	}
}

func createPassage(coordX, coordY, width, height int, passages *Passages) {
	passages.Passages[passages.PassagesNum].Coordinates[X] = coordX - 1
	passages.Passages[passages.PassagesNum].Coordinates[Y] = coordY - 1
	passages.Passages[passages.PassagesNum].Size[X] = width + 2
	passages.Passages[passages.PassagesNum].Size[Y] = height + 2
	passages.PassagesNum++
}

func generateHorizontalPassage(firstRoom, secondRoom int, rooms *[RoomsNum]Room, passages *Passages) {
	firstCoords := rooms[firstRoom].Coords
	secondCoords := rooms[secondRoom].Coords

	firstX := firstCoords.Coordinates[X] + firstCoords.Size[X] - 1
	upRangeCoord := firstCoords.Coordinates[Y] + 1
	bottomRangeCoord := firstCoords.Coordinates[Y] + firstCoords.Size[Y] - 2
	firstY := getRandomInRange(upRangeCoord, bottomRangeCoord)

	secondX := secondCoords.Coordinates[X]
	upRangeCoord = secondCoords.Coordinates[Y] + 1
	bottomRangeCoord = secondCoords.Coordinates[Y] + secondCoords.Size[Y] - 2
	secondY := getRandomInRange(upRangeCoord, bottomRangeCoord)

	if firstY == secondY {
		createPassage(firstX, firstY, Abs(secondX-firstX)+1, 1, passages)
	} else {
		vertical := getRandomInRange(min(firstX, secondX)+1, max(firstX, secondX)-1)
		createPassage(firstX, firstY, Abs(vertical-firstX)+1, 1, passages)
		createPassage(vertical, min(firstY, secondY), 1, Abs(secondY-firstY)+1, passages)
		createPassage(vertical, secondY, Abs(secondX-vertical)+1, 1, passages)
	}
}

func generateVerticalPassage(firstRoom, secondRoom int, rooms *[RoomsNum]Room, passages *Passages) {
	firstCoords := rooms[firstRoom].Coords
	secondCoords := rooms[secondRoom].Coords

	firstY := firstCoords.Coordinates[Y] + firstCoords.Size[Y] - 1
	upRangeCoord := firstCoords.Coordinates[X] + 1
	bottomRangeCoord := firstCoords.Coordinates[X] + firstCoords.Size[X] - 2
	firstX := getRandomInRange(upRangeCoord, bottomRangeCoord)

	secondY := secondCoords.Coordinates[Y]
	upRangeCoord = secondCoords.Coordinates[X] + 1
	bottomRangeCoord = secondCoords.Coordinates[X] + secondCoords.Size[X] - 2
	secondX := getRandomInRange(upRangeCoord, bottomRangeCoord)

	if firstX == secondX {
		createPassage(firstX, firstY, 1, Abs(secondY-firstY)+1, passages)
	} else {
		horizont := getRandomInRange(min(firstY, secondY)+1, max(firstY, secondY)-1)
		createPassage(firstX, firstY, 1, Abs(horizont-firstY)+1, passages)
		createPassage(min(firstX, secondX), horizont, Abs(secondX-firstX)+1, 1, passages)
		createPassage(secondX, horizont, 1, Abs(secondY-horizont)+1, passages)
	}
}

func generatePassages(passages *Passages, rooms *[RoomsNum]Room) {
	passages.PassagesNum = 0
	var countPassages int
	edges := make([]edge, MaxPassagesNum)
	generateEdgesForRooms(edges, &countPassages)
	shuffleArray(edges[:countPassages])

	parent := make([]int, RoomsNum)
	rank := make([]int, RoomsNum)
	makeSets(parent, rank, RoomsNum)

	for i := range countPassages {
		if findSet(edges[i].u, parent) != findSet(edges[i].v, parent) {
			unionSets(edges[i].u, edges[i].v, parent, rank)
			if Abs(edges[i].u-edges[i].v) == 1 {
				generateHorizontalPassage(edges[i].u, edges[i].v, rooms, passages)
			} else {
				generateVerticalPassage(edges[i].u, edges[i].v, rooms, passages)
			}
		}
	}
}

func generateCoordsOfEntity(room *Room, coords *Object) {
	upperLeftX := room.Coords.Coordinates[X] + 1
	upperLeftY := room.Coords.Coordinates[Y] + 1

	bottomRightX := upperLeftX + room.Coords.Size[X] - 3
	bottomRightY := upperLeftY + room.Coords.Size[Y] - 3

	coords.Coordinates[X] = getRandomInRange(upperLeftX, bottomRightX)
	coords.Coordinates[Y] = getRandomInRange(upperLeftY, bottomRightY)

	coords.Size[X] = 1
	coords.Size[Y] = 1
}

func generatePlayer(rooms *[RoomsNum]Room, player *Player) int {
	playerRoom := getRandomInRange(0, RoomsNum-1)
	generateCoordsOfEntity(&rooms[playerRoom], &player.BaseStats.Coords)
	return playerRoom
}

func generateExit(level *Level, playerRoom int) {
	var exitRoom int
	for {
		exitRoom = getRandomInRange(0, RoomsNum-1)
		for exitRoom == playerRoom {
			exitRoom = getRandomInRange(0, RoomsNum-1)
		}

		upperLeftX := level.Rooms[exitRoom].Coords.Coordinates[X] + 2
		upperLeftY := level.Rooms[exitRoom].Coords.Coordinates[Y] + 2

		bottomRightX := upperLeftX + level.Rooms[exitRoom].Coords.Size[X] - 5
		bottomRightY := upperLeftY + level.Rooms[exitRoom].Coords.Size[Y] - 5

		level.EndOfLevel.Coordinates[X] = getRandomInRange(upperLeftX, bottomRightX)
		level.EndOfLevel.Coordinates[Y] = getRandomInRange(upperLeftY, bottomRightY)

		level.EndOfLevel.Size[X] = 1
		level.EndOfLevel.Size[Y] = 1

		if checkUnoccupiedRoom(&level.EndOfLevel, &level.Rooms[exitRoom]) {
			break
		}
	}
}

func generateMonsterData(monster *Monster, levelNum int) {
	monster.Type = MonsterType(getRandomInRange(0, int(MonsterTypeNum)-1))

	switch monster.Type {
	case Zombie:
		monster.hostility = average
		monster.BaseStats.Agility = 25
		monster.BaseStats.Strength = 125
		monster.BaseStats.Health = 50
	case Vampire:
		monster.hostility = high
		monster.BaseStats.Agility = 75
		monster.BaseStats.Strength = 125
		monster.BaseStats.Health = 50
	case Ghost:
		monster.hostility = low
		monster.BaseStats.Agility = 75
		monster.BaseStats.Strength = 25
		monster.BaseStats.Health = 75
	case Ogre:
		monster.hostility = average
		monster.BaseStats.Agility = 25
		monster.BaseStats.Strength = 100
		monster.BaseStats.Health = 150
	case Snake:
		monster.hostility = high
		monster.BaseStats.Agility = 100
		monster.BaseStats.Strength = 30
		monster.BaseStats.Health = 100
	case Mimic:
		monster.hostility = low
		monster.BaseStats.Agility = 75
		monster.BaseStats.Strength = 20
		monster.BaseStats.Health = 150
	}

	percentsUpdate := (percentsUpdateDifficultyMonsters * levelNum)
	monster.BaseStats.Agility += monster.BaseStats.Agility * percentsUpdate / 100
	monster.BaseStats.Strength += monster.BaseStats.Strength * percentsUpdate / 100
	monster.BaseStats.Health += monster.BaseStats.Health * float64(percentsUpdate) / 100

	monster.IsChasing = false
	monster.dir = stop
}
func generateMonsters(level *Level, playerRoom int) {
	maxMonsters := MaxMonstersPerRoom + level.LevelNum/levelUpdateDifficulty

	for room := range RoomsNum {
		if room == playerRoom {
			continue
		}

		countMonsters := getRandomInRange(0, maxMonsters)

		for i := range countMonsters {
			coords := &level.Rooms[room].Monsters[i].BaseStats.Coords

			for {
				generateCoordsOfEntity(&level.Rooms[room], coords)
				if checkUnoccupiedRoom(coords, &level.Rooms[room]) {
					break
				}
			}

			generateMonsterData(&level.Rooms[room].Monsters[i], level.LevelNum)
		}

		level.Rooms[room].MonsterNum = countMonsters
	}
}

func generateFoodData(food *Food, player *Player) {
	names := []string{
		"Ration of the Ironclad",
		"Crimson Berry Cluster",
		"Loaf of the Forgotten Baker",
		"Smoked Wyrm Jerky",
		"Golden Apple of Vitality",
		"Hardtack of the Endless March",
		"Spiced Venison Strips",
		"Honeyed Nectar Bread",
		"Dried Mushrooms of the Deep",
	}

	maxRegen := int(player.BaseStats.Health) * maxPercentFoodRegenFromHealth / 100
	food.ToRegen = getRandomInRange(1, maxRegen)

	namePos := getRandomInRange(0, len(names)-1)
	copy(food.Name[:], []rune(names[namePos]))
}

func generateFood(room *Room, player *Player) {
	countFood := room.Consumables.FoodNum
	coords := &room.Consumables.RoomFood[countFood].Geometry
	for {
		generateCoordsOfEntity(room, coords)
		if checkUnoccupiedRoom(coords, room) {
			break
		}
	}
	generateFoodData(&room.Consumables.RoomFood[countFood].Food, player)
	room.Consumables.FoodNum++
}

func generateElixirData(elixir *Elixir, player *Player) {
	names := []string{
		"Elixir of the Jade Serpent",
		"Potion of the Phantom's Breath",
		"Vial of Crimson Vitality",
		"Draught of the Frozen Star",
		"Elixir of the Shattered Mind",
		"Potion of the Wandering Soul",
		"Vial of Ember Essence",
		"Elixir of the Obsidian Veil",
		"Potion of the Howling Wind",
	}

	statType := StatType(getRandomInRange(0, int(StatTypeNum)-1))
	maxIncrease := 0
	switch statType {
	case Health:
		maxIncrease = int(player.RegenLimit) * maxPercentFoodRegenFromHealth / 100
	case Agility:
		maxIncrease = player.BaseStats.Agility * maxPercentAgilityIncrease / 100
	case Strength:
		maxIncrease = player.BaseStats.Strength * maxPercentStrengthIncrease / 100
	}

	elixir.Stat = statType
	elixir.Increase = getRandomInRange(1, maxIncrease)
	elixir.Duration = time.Duration(getRandomInRange(minElixirDurationSeconds, maxElixirDurationSeconds)) * time.Second

	namePos := getRandomInRange(0, len(names)-1)
	copy(elixir.Name[:], []rune(names[namePos]))
}

func generateElixir(room *Room, player *Player) {
	countElixirs := room.Consumables.ElixirNum
	coords := &room.Consumables.Elixirs[countElixirs].Geometry
	for {
		generateCoordsOfEntity(room, coords)
		if checkUnoccupiedRoom(coords, room) {
			break
		}
	}
	generateElixirData(&room.Consumables.Elixirs[countElixirs].Elixir, player)
	room.Consumables.ElixirNum++
}

func generateScrollData(scroll *Scroll, player *Player) {
	names := []string{
		"Scroll of Shadowstep",
		"Parchment of Eternal Flame",
		"Manuscript of Forgotten Truths",
		"Scroll of Iron Will",
		"Vellum of the Void",
		"Scroll of Whispers",
		"Tome of the Lost King",
		"Scroll of Unseen Paths",
		"Parchment of Thunderous Roar",
	}

	statType := StatType(getRandomInRange(0, int(StatTypeNum)-1))
	maxIncrease := 0
	switch statType {
	case Health:
		maxIncrease = int(player.RegenLimit) * maxPercentFoodRegenFromHealth / 100
	case Agility:
		maxIncrease = player.BaseStats.Agility * maxPercentAgilityIncrease / 100
	case Strength:
		maxIncrease = player.BaseStats.Strength * maxPercentStrengthIncrease / 100
	}

	scroll.Stat = statType
	scroll.Increase = getRandomInRange(1, maxIncrease)

	namePos := getRandomInRange(0, len(names)-1)
	copy(scroll.Name[:], []rune(names[namePos]))
}

func generateScroll(room *Room, player *Player) {
	countScrolls := room.Consumables.ScrollNum
	coords := &room.Consumables.Scrolls[countScrolls].Geometry
	for {
		generateCoordsOfEntity(room, coords)
		if checkUnoccupiedRoom(coords, room) {
			break
		}
	}
	generateScrollData(&room.Consumables.Scrolls[countScrolls].Scroll, player)
	room.Consumables.ScrollNum++
}

func generateWeaponData(weapon *Weapon, player *Player) {
	names := []string{
		"Blade of the Forgotten Dawn",
		"Obsidian Reaver",
		"Fang of the Shadow Wolf",
		"Ironclad Cleaver",
		"Crimson Talon",
		"Thunderstrike Maul",
		"Serpent's Kiss Dagger",
		"Voidrend Sword",
		"Ebonheart Spear",
	}

	maxStrength := maxWeaponStrength
	if player.Weapon.Strength < maxStrength {
		maxStrength = player.Weapon.Strength
	}
	weapon.Strength = getRandomInRange(minWeaponStrength, maxStrength)

	namePos := getRandomInRange(0, len(names)-1)
	copy(weapon.Name[:], []rune(names[namePos]))
}

func generateWeapon(room *Room, player *Player) {
	countWeapons := room.Consumables.WeaponNum
	coords := &room.Consumables.Weapons[countWeapons].Geometry
	for {
		generateCoordsOfEntity(room, coords)
		if checkUnoccupiedRoom(coords, room) {
			break
		}
	}
	generateWeaponData(&room.Consumables.Weapons[countWeapons].Weapon, player)
	room.Consumables.WeaponNum++
}

func generateConsumables(level *Level, playerRoom int, player *Player, levelNum int) {
	generateConsumable := []generateConsumableFunc{
		generateFood,
		generateElixir,
		generateScroll,
		generateWeapon,
	}

	maxConsumables := MaxConsumablesPerRoom - levelNum/levelUpdateDifficulty
	if maxConsumables < 1 {
		maxConsumables = 1
	}

	for room := range RoomsNum {
		if room == playerRoom {
			continue
		}

		countConsumables := getRandomInRange(0, maxConsumables)

		for range countConsumables {
			typeConsumable := getRandomInRange(0, len(generateConsumable)-1)
			generateConsumable[typeConsumable](&level.Rooms[room], player)
		}
	}
}

func getRandomInRange(min, max int) int {
	if min > max {
		min, max = max, min
	}
	if min == max {
		return min
	}

	rangeLen := max - min + 1
	return rand.Intn(rangeLen) + min
}

func shuffleArray[T any](arr []T) {
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}
