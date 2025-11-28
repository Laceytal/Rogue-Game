package presentation

import (
	"log"
	"rogue/internal/domain"
	"rogue/internal/models"

	"github.com/gbin/goncurses"
)

const (
	Escape = 27
)

func checkDeath(player *domain.Player) bool {
	return player.BaseStats.Health <= 0
}

func checkEnd(level *domain.Level) bool {
	return level.LevelNum > domain.LevelNum
}

func checkLevelEnd(level *domain.Level, player *domain.Player) bool {
	return domain.CheckEqualCoords(level.EndOfLevel.Coordinates, player.BaseStats.Coords.Coordinates)
}

func printWeaponMenu(player *domain.Player, stdscr *goncurses.Window) {
	row, col := stdscr.MaxYX()
	shiftX := (col - 30) / 2
	shiftY := (row - 10) / 2
	countWeapon := player.Backpack.WeaponNum

	stdscr.MovePrint(shiftY-1, shiftX, "Choose weapon:")

	if countWeapon > 0 {
		stdscr.MovePrint(shiftY, shiftX, "0. Without weapon")
	}

	for i := 1; i <= countWeapon; i++ {
		stdscr.MovePrintf(shiftY+i, shiftX, "%d. %s %+d strength", i,
			runeArrayToString(player.Backpack.Weapons[i-1].Name[:]),
			player.Backpack.Weapons[i-1].Strength)
	}

	if countWeapon == 0 {
		stdscr.MovePrint(shiftY+1, shiftX, "You haven't weapon!")
		stdscr.MovePrint(shiftY+2, shiftX, "Press any key to continue...")
	} else {
		stdscr.MovePrintf(shiftY+countWeapon+1, shiftX,
			"Press 1-%d key to choose weapon or any key to continue", countWeapon)
	}
}

func printFoodMenu(player *domain.Player, stdscr *goncurses.Window) {
	row, col := stdscr.MaxYX()
	shiftX := (col - 30) / 2
	shiftY := (row - 10) / 2
	countFood := player.Backpack.FoodNum

	stdscr.MovePrint(shiftY, shiftX, "Choose food:")

	for i := 1; i <= countFood; i++ {
		stdscr.MovePrintf(shiftY+i, shiftX, "%d. %s %+d health", i,
			runeArrayToString(player.Backpack.Foods[i-1].Name[:]),
			player.Backpack.Foods[i-1].ToRegen)
	}

	if countFood == 0 {
		stdscr.MovePrint(shiftY+1, shiftX, "You haven't food!")
		stdscr.MovePrint(shiftY+2, shiftX, "Press any key to continue...")
	} else {
		stdscr.MovePrintf(shiftY+countFood+1, shiftX,
			"Press 1-%d key to choose food or any key to continue", countFood)
	}
}

func printScrollMenu(player *domain.Player, stdscr *goncurses.Window) {
	scrollType := []string{"health", "agility", "strength"}
	row, col := stdscr.MaxYX()
	shiftX := (col - 30) / 2
	shiftY := (row - 10) / 2
	countScroll := player.Backpack.ScrollNum

	stdscr.MovePrint(shiftY, shiftX, "Choose scroll:")

	for i := 1; i <= countScroll; i++ {
		statStr := scrollType[player.Backpack.Scrolls[i-1].Stat]
		stdscr.MovePrintf(shiftY+i, shiftX, "%d. %s %+d %s", i,
			runeArrayToString(player.Backpack.Scrolls[i-1].Name[:]),
			player.Backpack.Scrolls[i-1].Increase, statStr)
	}

	if countScroll == 0 {
		stdscr.MovePrint(shiftY+1, shiftX, "You haven't scroll!")
		stdscr.MovePrint(shiftY+2, shiftX, "Press any key to continue...")
	} else {
		stdscr.MovePrintf(shiftY+countScroll+1, shiftX,
			"Press 1-%d key to choose scroll or any key to continue", countScroll)
	}
}

func printElixirMenu(player *domain.Player, stdscr *goncurses.Window) {
	scrollType := []string{"health", "agility", "strength"}
	row, col := stdscr.MaxYX()
	shiftX := (col - 30) / 2
	shiftY := (row - 10) / 2
	countElixir := player.Backpack.ElixirNum

	stdscr.MovePrint(shiftY, shiftX, "Choose elixir:")

	for i := 1; i <= countElixir; i++ {
		statStr := scrollType[player.Backpack.Elixirs[i-1].Stat]
		stdscr.MovePrintf(shiftY+i, shiftX, "%d. %s %+d %s for %d seconds", i,
			runeArrayToString(player.Backpack.Elixirs[i-1].Name[:]),
			player.Backpack.Elixirs[i-1].Increase, statStr,
			int(player.Backpack.Elixirs[i-1].Duration.Seconds()))
	}

	if countElixir == 0 {
		stdscr.MovePrint(shiftY+1, shiftX, "You haven't elixir!")
		stdscr.MovePrint(shiftY+2, shiftX, "Press any key to continue...")
	} else {
		stdscr.MovePrintf(shiftY+countElixir+1, shiftX,
			"Press 1-%d key to choose elixir or any key to continue", countElixir)
	}
}

func chooseConsumable(player *domain.Player, consType domain.ConsumableTypes, room *domain.Room, adapter *domain.DifficultyAdapter, stdscr *goncurses.Window) {
	stdscr.Clear()
	countConsumable := 0
	switch consType {
	case domain.WeaponConTyp:
		printWeaponMenu(player, stdscr)
		countConsumable = player.Backpack.WeaponNum
	case domain.FoodConTyp:
		printFoodMenu(player, stdscr)
		countConsumable = player.Backpack.FoodNum
	case domain.ScrollConTyp:
		printScrollMenu(player, stdscr)
		countConsumable = player.Backpack.ScrollNum
	case domain.ElixirConTyp:
		printElixirMenu(player, stdscr)
		countConsumable = player.Backpack.ElixirNum
	default:
		return
	}

	key := stdscr.GetChar() - '0'
	if key >= 1 && key <= goncurses.Key(countConsumable) && (consType != domain.WeaponConTyp || room != nil) {
		domain.UseConsumable(player, consType, int(key)-1, room, adapter)
	}

	if key == 0 && consType == domain.WeaponConTyp {
		domain.UseConsumable(player, consType, -1, room, adapter)
	}

	stdscr.Clear()
}

func attackMonsterUI(monster *domain.Monster, wasAttack bool, stdscr *goncurses.Window) {
	monsterTypes := []string{"Zombie", "Vampire", "Ghost", "Ogre", "Snake", "Mimic"}
	row, col := stdscr.MaxYX()
	shiftX := (col - domain.MapWidth) / 2
	shiftY := (row - domain.MapHeight) / 2

	stdscr.Move(shiftY-2, shiftX)
	if wasAttack {
		stdscr.Printf("You attacked %s!!!", monsterTypes[monster.Type])
	} else {
		stdscr.Print("You missed...")
	}
}

func checkConsumableUI(player *domain.Player, level *domain.Level, stdscr *goncurses.Window) {
	oldBackpack := player.Backpack

	for room := range domain.RoomsNum {
		domain.CheckConsumable(player, &level.Rooms[room])
	}

	consType := "unknown"
	if oldBackpack.FoodNum != player.Backpack.FoodNum {
		consType = "food"
	} else if oldBackpack.WeaponNum != player.Backpack.WeaponNum {
		consType = "weapon"
	} else if oldBackpack.ScrollNum != player.Backpack.ScrollNum {
		consType = "scroll"
	} else if oldBackpack.ElixirNum != player.Backpack.ElixirNum {
		consType = "elixir"
	}

	if consType != "unknown" {
		row, col := stdscr.MaxYX()
		shiftX := (col - domain.MapWidth) / 2
		shiftY := (row - domain.MapHeight) / 2
		stdscr.MovePrintf(shiftY-2, shiftX, "You take the %s!!!", consType)
	}
}

func processPlayerMoveUI(player *domain.Player, level *domain.Level, battles []domain.BattleInfo, adapter *domain.DifficultyAdapter, direction domain.Directions, stat *models.SessionStat, stdscr *goncurses.Window) {
	attacked := false
	for i := range domain.MaximumFights {
		if battles[i].IsFight {
			monsterHealthBefore := battles[i].Enemy.BaseStats.Health
			if domain.CheckPlayerAttack(player, &battles[i], direction, adapter) {
				stat.Attacks++
				attacked = true
				attackMonsterUI(battles[i].Enemy, battles[i].Enemy.BaseStats.Health < monsterHealthBefore, stdscr)
				if battles[i].Enemy.BaseStats.Health <= 0 {
					stat.Enemies++
				}
			}
		}
	}
	if !attacked {
		domain.MovePlayer(player, level, direction)
		checkConsumableUI(player, level, stdscr)
	}
}

func printDataAboutMonsterAttack(monster *domain.Monster, wasAttack bool, stdscr *goncurses.Window) {
	monsterTypes := []string{"Zombie", "Vampire", "Ghost", "Ogre", "Snake", "Mimic"}
	row, col := stdscr.MaxYX()
	shiftX := (col - domain.MapWidth) / 2
	shiftY := (row - domain.MapHeight) / 2

	stdscr.Move(shiftY-1, shiftX)
	if wasAttack {
		stdscr.Printf("%s attacked!!!", monsterTypes[monster.Type])
	} else if monster.Type != domain.Ogre {
		stdscr.Printf("%s missed!!!", monsterTypes[monster.Type])
	}
}

func processMonstersMoveUI(player *domain.Player, level *domain.Level, battles []domain.BattleInfo, adapter *domain.DifficultyAdapter, stat *models.SessionStat, stdscr *goncurses.Window) {
	for i := range domain.RoomsNum {
		for j := range level.Rooms[i].MonsterNum {
			if domain.CheckUnique(&level.Rooms[i].Monsters[j], battles) {
				domain.MoveMonster(&level.Rooms[i].Monsters[j], &player.BaseStats.Coords, level)
			}
		}
	}

	domain.RemoveDeadMonsters(level)
	for i := range domain.MaximumFights {
		if battles[i].IsFight {
			playerHealthBefore := player.BaseStats.Health
			domain.Attack(player, &battles[i], domain.MonsterTurn, adapter)
			printDataAboutMonsterAttack(battles[i].Enemy, player.BaseStats.Health < playerHealthBefore, stdscr)
			if player.BaseStats.Health < playerHealthBefore {
				stat.Missed++
			}
		}
	}
}

func processUserMove(player *domain.Player, level *domain.Level, battles []domain.BattleInfo, adapter *domain.DifficultyAdapter, direction domain.Directions, stat *models.SessionStat, stdscr *goncurses.Window) {
	domain.UpdateFightStatus(&player.BaseStats.Coords, level, battles)
	processPlayerMoveUI(player, level, battles, adapter, direction, stat, stdscr)
	processMonstersMoveUI(player, level, battles, adapter, stat, stdscr)
	domain.CheckTempEffectEnd(player)
}

func processUserInput(player *domain.Player, level *domain.Level, battles []domain.BattleInfo, adapter *domain.DifficultyAdapter, filename string, stdscr *goncurses.Window) bool {
	stat, err := models.GetDataFromFile[models.SessionStat](filename)
	if err != nil || stat == nil {
		stat = &models.SessionStat{}
	}

	stat.Treasures = player.Backpack.Treasures.Value
	stat.Level = level.LevelNum
	key := stdscr.GetChar()
	row, col := stdscr.MaxYX()
	shiftX := (col - domain.MapWidth) / 2
	shiftY := (row - domain.MapHeight) / 2
	stdscr.Move(shiftY-1, shiftX)
	stdscr.ClearToEOL()
	stdscr.Move(shiftY-2, shiftX)
	stdscr.ClearToEOL()

	quit := false
	currentCount := 0

	switch key {
	case 'W', 'w':
		processUserMove(player, level, battles, adapter, domain.Forward, stat, stdscr)
		stat.Movies++
	case 'A', 'a':
		processUserMove(player, level, battles, adapter, domain.Left, stat, stdscr)
		stat.Movies++
	case 'S', 's':
		processUserMove(player, level, battles, adapter, domain.Back, stat, stdscr)
		stat.Movies++
	case 'D', 'd':
		processUserMove(player, level, battles, adapter, domain.Right, stat, stdscr)
		stat.Movies++
	case 'H', 'h':
		chooseConsumable(player, domain.WeaponConTyp, domain.FindCurrentRoom(&player.BaseStats.Coords, level), adapter, stdscr)
	case 'J', 'j':
		currentCount = player.Backpack.FoodNum
		chooseConsumable(player, domain.FoodConTyp, nil, adapter, stdscr)
		if currentCount != player.Backpack.FoodNum {
			stat.Food++
		}
	case 'K', 'k':
		currentCount = player.Backpack.ElixirNum
		chooseConsumable(player, domain.ElixirConTyp, nil, adapter, stdscr)
		if currentCount != player.Backpack.ElixirNum {
			stat.Elixirs++
		}
	case 'E', 'e':
		currentCount = player.Backpack.ScrollNum
		chooseConsumable(player, domain.ScrollConTyp, nil, adapter, stdscr)
		if currentCount != player.Backpack.ScrollNum {
			stat.Scrolls++
		}
	case Escape:
		quit = true
	}

	err = models.SaveDataToFile[models.SessionStat](filename, *stat)
	if err != nil {
		log.Fatal(err)
	}

	return quit
}

func GameCycle(player *domain.Player, level *domain.Level, mapData *domain.Map, battles []domain.BattleInfo, adapter *domain.DifficultyAdapter, save, score, stat string, stdscr *goncurses.Window) {
	running := true
	normalExit := false

	for running {
		DisplayMap(mapData, level, player, battles, stdscr)

		if processUserInput(player, level, battles, adapter, stat, stdscr) {
			running = false
			normalExit = true
		}

		if checkLevelEnd(level, player) {
			if adapter != nil {
				adapter.AdjustDifficulty(player)
				adapter.ResetLevelStats()
			}

			ClearMapData(mapData)
			domain.GenerateNextLevel(level, player)
		}

		if checkDeath(player) {
			running = false
			err := models.UpdateStatistic(score, stat)
			if err != nil {
				log.Fatal(err)
			}
			models.GetStandartSave(save, stat)
			DeadScreen(stdscr)
		}

		if checkEnd(level) {
			running = false
			err := models.UpdateStatistic(score, stat)
			if err != nil {
				log.Fatal(err)
			}
			models.GetStandartSave(save, stat)
			EndgameScreen(stdscr)
		}
	}

	if normalExit {
		err := models.SaveDataToFile(save, domain.Game{
			Player:            player,
			Level:             level,
			Map:               mapData,
			DifficultyAdapter: adapter,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runeArrayToString(runes []rune) string {
	for i, r := range runes {
		if r == 0 {
			return string(runes[:i])
		}
	}
	return string(runes)
}
