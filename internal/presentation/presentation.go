package presentation

import (
	"math/rand"
	"rogue/internal/domain"
	"rogue/internal/models"

	"github.com/gbin/goncurses"
)

func InitPresentation(stdscr *goncurses.Window) {
	goncurses.Echo(false)
	goncurses.Cursor(0)
	goncurses.StartColor()
	stdscr.Keypad(true)

	goncurses.InitPair(int16(domain.WhiteFont), goncurses.C_WHITE, goncurses.C_BLACK)
	goncurses.InitPair(int16(domain.RedFont), goncurses.C_RED, goncurses.C_BLACK)
	goncurses.InitPair(int16(domain.GreenFont), goncurses.C_GREEN, goncurses.C_BLACK)
	goncurses.InitPair(int16(domain.BlueFont), goncurses.C_BLUE, goncurses.C_BLACK)
	goncurses.InitPair(int16(domain.YellowFont), goncurses.C_YELLOW, goncurses.C_BLACK)
	goncurses.InitPair(int16(domain.CyanFont), goncurses.C_CYAN, goncurses.C_BLACK)
}

func ClearMapData(mapData *domain.Map) {
	for i := range domain.MapHeight {
		for j := range domain.MapWidth {
			mapData.MapData[i][j] = domain.Cell{Char: ' ', Color: domain.WhiteFont}
		}
	}

	for i := range domain.RoomsNum {
		mapData.VisibleRooms[i] = false
	}

	for i := range domain.MaxPassagesNum {
		mapData.VisiblePassages[i] = false
	}
}

func getRoomByCoord(coords *domain.Object, rooms *[domain.RoomsNum]domain.Room) int {
	x := coords.Coordinates[domain.X]
	y := coords.Coordinates[domain.Y]

	for room := 0; room < domain.RoomsNum; room++ {
		xRoom := rooms[room].Coords.Coordinates[domain.X]
		yRoom := rooms[room].Coords.Coordinates[domain.Y]
		xsize := rooms[room].Coords.Size[domain.X]
		ysize := rooms[room].Coords.Size[domain.Y]

		checkX := (x >= xRoom) && (x < xRoom+xsize)
		checkY := (y >= yRoom) && (y < yRoom+ysize)

		if checkY && checkX {
			return room
		}
	}

	return -1
}

func roomsToMap(mapData *domain.Map, rooms *[domain.RoomsNum]domain.Room, player *domain.Player) {
	for i := range domain.RoomsNum {
		if !mapData.VisibleRooms[i] && getRoomByCoord(&player.BaseStats.Coords, rooms) != i {
			continue
		}

		x1 := rooms[i].Coords.Coordinates[domain.X]
		y1 := rooms[i].Coords.Coordinates[domain.Y]
		xsize := rooms[i].Coords.Size[domain.X]
		ysize := rooms[i].Coords.Size[domain.Y]

		for y := range domain.MapHeight {
			for x := range domain.MapWidth {
				checkX := (x == x1 || x == x1+xsize-1) && (y1 <= y && y < y1+ysize)
				checkY := (y == y1 || y == y1+ysize-1) && (x1 <= x && x < x1+xsize)

				if checkY {
					mapData.MapData[y][x] = domain.Cell{Char: '-', Color: domain.WhiteFont}
				} else if checkX {
					mapData.MapData[y][x] = domain.Cell{Char: '|', Color: domain.WhiteFont}
				}
			}
		}

		mapData.VisibleRooms[i] = true
	}
}

func passagesToMap(mapData *domain.Map, passages *domain.Passages, rooms *[domain.RoomsNum]domain.Room, player *domain.Player) {
	for i := range passages.PassagesNum {
		visible := true
		if !mapData.VisiblePassages[i] && domain.CharacterOutsideBorder(&player.BaseStats.Coords, &passages.Passages[i]) {
			visible = false
		}

		x1 := passages.Passages[i].Coordinates[domain.X]
		y1 := passages.Passages[i].Coordinates[domain.Y]
		xsize := passages.Passages[i].Size[domain.X]
		ysize := passages.Passages[i].Size[domain.Y]

		for y := range domain.MapHeight {
			for x := range domain.MapWidth {
				checkX := (x1 < x && x < x1+xsize-1) && (y1 < y && y < y1+ysize-1)
				coords := domain.Object{
					Coordinates: domain.Coordinates{x, y},
				}
				room := getRoomByCoord(&coords, rooms)

				if checkX && visible {
					if room != -1 {
						mapData.MapData[y][x] = domain.Cell{Char: '+', Color: domain.WhiteFont}
					} else {
						mapData.MapData[y][x] = domain.Cell{Char: '#', Color: domain.WhiteFont}
					}
				} else if checkX && room != -1 && mapData.VisibleRooms[room] {
					mapData.MapData[y][x] = domain.Cell{Char: '+', Color: domain.WhiteFont}
				}
			}
		}

		mapData.VisiblePassages[i] = visible
	}
}

func unvisibleGhost(monster *domain.Monster, battles []domain.BattleInfo) bool {
	unvisible := false
	if rand.Intn(100) < domain.ChanceUnvisibleGhost {
		unvisible = true
	}
	if !domain.CheckUnique(monster, battles) {
		unvisible = false
	}
	return unvisible
}

func onTheSameRoomOrPassage(level *domain.Level, characterCoords, monsterCoords *domain.Object) bool {
	same := false

	for i := 0; i < domain.RoomsNum && !same; i++ {
		if !domain.CharacterOutsideBorder(characterCoords, &level.Rooms[i].Coords) &&
			!domain.CharacterOutsideBorder(monsterCoords, &level.Rooms[i].Coords) {
			same = true
		}
	}

	for i := 0; i < level.Passages.PassagesNum && !same; i++ {
		if !domain.CharacterOutsideBorder(characterCoords, &level.Passages.Passages[i]) &&
			!domain.CharacterOutsideBorder(monsterCoords, &level.Passages.Passages[i]) {
			same = true
		}
	}

	return same
}

func monstersToMap(mapData *domain.Map, level *domain.Level, player *domain.Player, battles []domain.BattleInfo) {
	monsterLetters := "zvgOsm"
	monsterColors := []domain.Font{
		domain.GreenFont,
		domain.RedFont,
		domain.WhiteFont,
		domain.YellowFont,
		domain.WhiteFont,
		domain.WhiteFont,
	}

	for i := range domain.RoomsNum {
		for j := range level.Rooms[i].MonsterNum {
			if onTheSameRoomOrPassage(level, &player.BaseStats.Coords, &level.Rooms[i].Monsters[j].BaseStats.Coords) ||
				!domain.CheckUnique(&level.Rooms[i].Monsters[j], battles) {

				x := level.Rooms[i].Monsters[j].BaseStats.Coords.Coordinates[domain.X]
				y := level.Rooms[i].Monsters[j].BaseStats.Coords.Coordinates[domain.Y]
				monsterType := level.Rooms[i].Monsters[j].Type

				if monsterType == domain.Mimic && !level.Rooms[i].Monsters[j].IsChasing {
					seed := x*31 + y*17
					disguiseIndex := seed % 4
					disguises := []struct {
						char  rune
						color domain.Font
					}{
						{'f', domain.WhiteFont},
						{'e', domain.WhiteFont},
						{'S', domain.WhiteFont},
						{'w', domain.WhiteFont},
					}
					mapData.MapData[y][x] = domain.Cell{
						Char:  disguises[disguiseIndex].char,
						Color: disguises[disguiseIndex].color,
					}
					continue
				}

				mapData.MapData[y][x] = domain.Cell{
					Char:  rune(monsterLetters[monsterType]),
					Color: monsterColors[monsterType],
				}

				if monsterType == domain.Ghost && unvisibleGhost(&level.Rooms[i].Monsters[j], battles) {
					mapData.MapData[y][x].Char = ' '
				}
			}
		}
	}
}

func consumablesToMap(mapData *domain.Map, rooms *[domain.RoomsNum]domain.Room) {
	for i := range domain.RoomsNum {
		for j := 0; j < rooms[i].Consumables.FoodNum && mapData.VisibleRooms[i]; j++ {
			x := rooms[i].Consumables.RoomFood[j].Geometry.Coordinates[domain.X]
			y := rooms[i].Consumables.RoomFood[j].Geometry.Coordinates[domain.Y]
			mapData.MapData[y][x] = domain.Cell{Char: 'f', Color: domain.WhiteFont}
		}

		for j := 0; j < rooms[i].Consumables.ElixirNum && mapData.VisibleRooms[i]; j++ {
			x := rooms[i].Consumables.Elixirs[j].Geometry.Coordinates[domain.X]
			y := rooms[i].Consumables.Elixirs[j].Geometry.Coordinates[domain.Y]
			mapData.MapData[y][x] = domain.Cell{Char: 'e', Color: domain.WhiteFont}
		}

		for j := 0; j < rooms[i].Consumables.ScrollNum && mapData.VisibleRooms[i]; j++ {
			x := rooms[i].Consumables.Scrolls[j].Geometry.Coordinates[domain.X]
			y := rooms[i].Consumables.Scrolls[j].Geometry.Coordinates[domain.Y]
			mapData.MapData[y][x] = domain.Cell{Char: 'S', Color: domain.WhiteFont}
		}

		for j := 0; j < rooms[i].Consumables.WeaponNum && mapData.VisibleRooms[i]; j++ {
			x := rooms[i].Consumables.Weapons[j].Geometry.Coordinates[domain.X]
			y := rooms[i].Consumables.Weapons[j].Geometry.Coordinates[domain.Y]
			mapData.MapData[y][x] = domain.Cell{Char: 'w', Color: domain.WhiteFont}
		}
	}
}

func exitToMap(mapData *domain.Map, level *domain.Level) {
	x := level.EndOfLevel.Coordinates[domain.X]
	y := level.EndOfLevel.Coordinates[domain.Y]
	roomEndOfLevel := getRoomByCoord(&level.EndOfLevel, &level.Rooms)

	if roomEndOfLevel >= 0 && mapData.VisibleRooms[roomEndOfLevel] {
		mapData.MapData[y][x] = domain.Cell{Char: 'E', Color: domain.WhiteFont}
	}
}

func playerToMap(mapData *domain.Map, player *domain.Player) {
	x := player.BaseStats.Coords.Coordinates[domain.X]
	y := player.BaseStats.Coords.Coordinates[domain.Y]
	mapData.MapData[y][x] = domain.Cell{Char: '@', Color: domain.WhiteFont}
}

func fillRoomByFog(mapData *domain.Map, room *domain.Room) {
	xRoom := room.Coords.Coordinates[domain.X]
	yRoom := room.Coords.Coordinates[domain.Y]
	xsize := room.Coords.Size[domain.X]
	ysize := room.Coords.Size[domain.Y]

	for x := xRoom + 1; x < xRoom+xsize-1; x++ {
		for y := yRoom + 1; y < yRoom+ysize-1; y++ {
			mapData.MapData[y][x] = domain.Cell{Char: '.', Color: domain.WhiteFont}
		}
	}
}

func isVerticalDirectionFog(coords, roomCoords *domain.Object) bool {
	newCoords := *coords
	newCoords.Coordinates[domain.X]++
	if !domain.CharacterOutsideBorder(&newCoords, roomCoords) {
		return false
	}

	newCoords.Coordinates[domain.X] -= 2
	if !domain.CharacterOutsideBorder(&newCoords, roomCoords) {
		return false
	}

	return true
}

func fillRoomByPartFog(mapData *domain.Map, player *domain.Player, room *domain.Room) {
	playerX := player.BaseStats.Coords.Coordinates[domain.X]
	playerY := player.BaseStats.Coords.Coordinates[domain.Y]
	isVertical := isVerticalDirectionFog(&player.BaseStats.Coords, &room.Coords)

	for y := range domain.MapHeight {
		for x := range domain.MapWidth {
			coords := domain.Object{
				Coordinates: domain.Coordinates{x, y},
				Size:        domain.Sizes{1, 1},
			}

			if domain.CharacterOutsideBorder(&coords, &room.Coords) {
				continue
			}

			newX := x - playerX
			newY := y - playerY

			if isVertical && domain.Abs(newX) >= domain.Abs(newY) {
				mapData.MapData[y][x] = domain.Cell{Char: '.', Color: domain.WhiteFont}
			}
			if !isVertical && domain.Abs(newX) <= domain.Abs(newY) {
				mapData.MapData[y][x] = domain.Cell{Char: '.', Color: domain.WhiteFont}
			}
		}
	}
}

func fogOfWarToMap(mapData *domain.Map, level *domain.Level, player *domain.Player) {
	roomPlayer := getRoomByCoord(&player.BaseStats.Coords, &level.Rooms)

	for room := range domain.RoomsNum {
		if room != roomPlayer && mapData.VisibleRooms[room] {
			fillRoomByFog(mapData, &level.Rooms[room])
		}

		if room == roomPlayer && getRoomByCoord(&player.BaseStats.Coords, &level.Rooms) != -1 &&
			domain.CharacterOutsideBorder(&player.BaseStats.Coords, &level.Rooms[room].Coords) {
			fillRoomByPartFog(mapData, player, &level.Rooms[room])
		}
	}
}

func CreateNewMap(mapData *domain.Map, level *domain.Level, player *domain.Player, battles []domain.BattleInfo) {
	for i := range domain.MapHeight {
		for j := range domain.MapWidth {
			mapData.MapData[i][j] = domain.Cell{Char: ' ', Color: domain.WhiteFont}
		}
	}

	roomsToMap(mapData, &level.Rooms, player)
	passagesToMap(mapData, &level.Passages, &level.Rooms, player)
	monstersToMap(mapData, level, player, battles)
	consumablesToMap(mapData, &level.Rooms)
	exitToMap(mapData, level)
	playerToMap(mapData, player)
	fogOfWarToMap(mapData, level, player)
}

func DisplayMap(mapData *domain.Map, level *domain.Level, player *domain.Player, battles []domain.BattleInfo, stdscr *goncurses.Window) {
	CreateNewMap(mapData, level, player, battles)

	row, col := stdscr.MaxYX()
	shiftX := (col - domain.MapWidth) / 2
	shiftY := (row - domain.MapHeight) / 2

	for i := range domain.MapHeight {
		stdscr.Move(shiftY+i, shiftX)
		for j := range domain.MapWidth {
			cell := mapData.MapData[i][j]

			stdscr.AttrOn(goncurses.ColorPair(int16(cell.Color)))
			stdscr.Printf("%c", cell.Char)
			stdscr.AttrOff(goncurses.ColorPair(int16(cell.Color)))
		}
	}

	stdscr.MovePrintf(shiftY+domain.MapHeight, shiftX,
		"Level: %-8d Gold: %-8d Health: %.2f/%-8d Agility: %-6d Strength: %d(+%d) ",
		level.LevelNum, player.Backpack.Treasures.Value, player.BaseStats.Health, player.RegenLimit,
		player.BaseStats.Agility, player.BaseStats.Strength, player.Weapon.Strength)

	stdscr.Move(row, col)
	stdscr.Refresh()
}

func DisplayScoreboard(pathScoreboard string, stdscr *goncurses.Window) {
	row, col := stdscr.MaxYX()
	stdscr.Clear()
	stats := make([]models.SessionStat, 0)
	statsP, err := models.GetDataFromFile[[]models.SessionStat](pathScoreboard)
	if err == nil {
		stats = append(stats, *statsP...)
	}

	sizeArray := len(stats)
	if sizeArray > domain.MaxScoreboardSize {
		sizeArray = domain.MaxScoreboardSize
	}

	fieldSize := 10
	tableWidth := fieldSize * 10
	tableHeight := sizeArray*2 + 3
	shiftX := (col - tableWidth) / 2
	shiftY := (row - tableHeight) / 2

	printLineCh(shiftY-2, shiftX, tableWidth, stdscr)

	headers := []string{"treasures", "level", "enemies", "food", "elixirs", "scrolls", "attacks", "missed", "moves"}
	stdscr.Move(shiftY-1, shiftX)
	for _, header := range headers {
		stdscr.Printf("|%-*s", fieldSize, header)
	}
	stdscr.Print("|\n")

	for i := range sizeArray {
		printLineCh(shiftY+2*i, shiftX, tableWidth, stdscr)

		stdscr.MovePrintf(shiftY+2*i+1, shiftX,
			"|%*d|%*d|%*d|%*d|%*d|%*d|%*d|%*d|%*d|\n",
			fieldSize, (stats)[i].Treasures,
			fieldSize, (stats)[i].Level,
			fieldSize, (stats)[i].Enemies,
			fieldSize, (stats)[i].Food,
			fieldSize, (stats)[i].Elixirs,
			fieldSize, (stats)[i].Scrolls,
			fieldSize, (stats)[i].Attacks,
			fieldSize, (stats)[i].Missed,
			fieldSize, (stats)[i].Movies)
	}

	printLineCh(shiftY+2*sizeArray, shiftX, tableWidth, stdscr)
	stdscr.MovePrint(shiftY+2*(sizeArray+1), (col-20)/2, "Press ESCAPE to exit.")
	stdscr.Refresh()
}

func printLineCh(y, x, width int, stdscr *goncurses.Window) {
	stdscr.Move(y, x)
	for range width {
		stdscr.AddChar('-')
	}
}
