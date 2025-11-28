package domain

import "time"

type ConsumableTypes int

const (
	NoneConTyp ConsumableTypes = iota
	FoodConTyp
	ElixirConTyp
	WeaponConTyp
	ScrollConTyp
)

func CheckConsumable(player *Player, room *Room) {
	wasConsumable := false
	for i := 0; i < room.Consumables.ElixirNum && !wasConsumable && player.Backpack.ElixirNum < ConsumablesTypeMaxNum; i++ {
		if CheckEqualCoords(room.Consumables.Elixirs[i].Geometry.Coordinates, player.BaseStats.Coords.Coordinates) {
			takeElixir(&player.Backpack, room, &room.Consumables.Elixirs[i])
			player.Backpack.CurrentSize++
			wasConsumable = true
		}
	}

	for i := 0; i < room.Consumables.ScrollNum && !wasConsumable && player.Backpack.ScrollNum < ConsumablesTypeMaxNum; i++ {
		if CheckEqualCoords(room.Consumables.Scrolls[i].Geometry.Coordinates, player.BaseStats.Coords.Coordinates) {
			takeScroll(&player.Backpack, room, &room.Consumables.Scrolls[i])
			player.Backpack.CurrentSize++
			wasConsumable = true
		}
	}

	for i := 0; i < room.Consumables.FoodNum && !wasConsumable && player.Backpack.FoodNum < ConsumablesTypeMaxNum; i++ {
		if CheckEqualCoords(room.Consumables.RoomFood[i].Geometry.Coordinates, player.BaseStats.Coords.Coordinates) {
			takeFood(&player.Backpack, room, &room.Consumables.RoomFood[i])
			player.Backpack.CurrentSize++
			wasConsumable = true
		}
	}

	for i := 0; i < room.Consumables.WeaponNum && !wasConsumable && player.Backpack.WeaponNum < ConsumablesTypeMaxNum; i++ {
		if CheckEqualCoords(room.Consumables.Weapons[i].Geometry.Coordinates, player.BaseStats.Coords.Coordinates) {
			takeWeapon(player, room, &room.Consumables.Weapons[i])
			player.Backpack.CurrentSize++
			wasConsumable = true
		}
	}
}

func takeScroll(playersBackpack *Backpack, currentRoom *Room, scroll *ScrollRoom) {
	deleteFromRoom(currentRoom, &scroll.Geometry, ScrollConTyp)
	playersBackpack.Scrolls[playersBackpack.ScrollNum] = scroll.Scroll
	playersBackpack.ScrollNum++
}

func takeElixir(playersBackpack *Backpack, currentRoom *Room, elixir *ElixirRoom) {
	deleteFromRoom(currentRoom, &elixir.Geometry, ElixirConTyp)
	playersBackpack.Elixirs[playersBackpack.ElixirNum] = elixir.Elixir
	playersBackpack.ElixirNum++
}

func takeFood(playersBackpack *Backpack, currentRoom *Room, food *FoodRoom) {
	deleteFromRoom(currentRoom, &food.Geometry, FoodConTyp)
	playersBackpack.Foods[playersBackpack.FoodNum] = food.Food
	playersBackpack.FoodNum++
}

func takeWeapon(player *Player, currentRoom *Room, weapon *WeaponRoom) {
	deleteFromRoom(currentRoom, &weapon.Geometry, WeaponConTyp)
	player.Backpack.Weapons[player.Backpack.WeaponNum] = weapon.Weapon
	player.Backpack.WeaponNum++
}

func deleteFromRoom(room *Room, consumableCoords *Object, consumableType ConsumableTypes) {
	switch consumableType {
	case ElixirConTyp:
		i := 0
		for ; i < room.Consumables.ElixirNum && !CheckEqualCoords(consumableCoords.Coordinates, room.Consumables.Elixirs[i].Geometry.Coordinates); i++ {
		}
		if i != room.Consumables.ElixirNum {
			room.Consumables.Elixirs[i] = room.Consumables.Elixirs[room.Consumables.ElixirNum-1]
			room.Consumables.ElixirNum--
		}

	case FoodConTyp:
		i := 0
		for ; i < room.Consumables.FoodNum && !CheckEqualCoords(consumableCoords.Coordinates, room.Consumables.RoomFood[i].Geometry.Coordinates); i++ {
		}
		if i != room.Consumables.FoodNum {
			room.Consumables.RoomFood[i] = room.Consumables.RoomFood[room.Consumables.FoodNum-1]
			room.Consumables.FoodNum--
		}

	case ScrollConTyp:
		i := 0
		for ; i < room.Consumables.ScrollNum && !CheckEqualCoords(consumableCoords.Coordinates, room.Consumables.Scrolls[i].Geometry.Coordinates); i++ {
		}
		if i != room.Consumables.ScrollNum {
			room.Consumables.Scrolls[i] = room.Consumables.Scrolls[room.Consumables.ScrollNum-1]
			room.Consumables.ScrollNum--
		}

	case WeaponConTyp:
		i := 0
		for ; i < room.Consumables.WeaponNum && !CheckEqualCoords(consumableCoords.Coordinates, room.Consumables.Weapons[i].Geometry.Coordinates); i++ {
		}
		if i != ConsumablesTypeMaxNum {
			room.Consumables.Weapons[i] = room.Consumables.Weapons[room.Consumables.WeaponNum-1]
			room.Consumables.WeaponNum--
		}

	default:
	}
	room.ConsumablesNum--
}

func checkUnoccupiedLevel(coordinates *Object, level *Level) bool {
	if CheckEqualCoords(coordinates.Coordinates, level.EndOfLevel.Coordinates) {
		return false
	}

	for _, room := range level.Rooms {
		if !checkUnoccupiedRoom(coordinates, &room) {
			return false
		}
	}

	return true
}

func throwCurrentWeapon(player *Player, currentRoom *Room, oldWeapon Weapon) {
	i := 0
	for ; i < player.Backpack.WeaponNum && !equalWeapons(&oldWeapon, &player.Backpack.Weapons[i]); i++ {
	}

	currentWeapon := player.Backpack.Weapons[player.Backpack.WeaponNum-1]
	player.Backpack.Weapons[player.Backpack.WeaponNum-1] = player.Backpack.Weapons[i]
	player.Backpack.Weapons[i] = currentWeapon
	player.Backpack.WeaponNum--
	throwOnGround(player, currentRoom, oldWeapon)
	player.Backpack.CurrentSize--
}

func throwOnGround(player *Player, currentRoom *Room, weapon Weapon) {
	currentRoom.Consumables.Weapons[currentRoom.Consumables.WeaponNum].Geometry = player.BaseStats.Coords
	currentRoom.Consumables.Weapons[currentRoom.Consumables.WeaponNum].Weapon = weapon
	roomWeapon := currentRoom.Consumables.Weapons[currentRoom.Consumables.WeaponNum]
	moveCharacterByDirection(Right, &roomWeapon.Geometry)
	for direction := Directions(0); CharacterOutsideBorder(&roomWeapon.Geometry, &currentRoom.Coords) || !checkUnoccupiedRoom(&roomWeapon.Geometry, currentRoom); direction++ {
		roomWeapon = currentRoom.Consumables.Weapons[currentRoom.Consumables.WeaponNum]
		moveCharacterByDirection(direction, &roomWeapon.Geometry)
	}
	currentRoom.Consumables.Weapons[currentRoom.Consumables.WeaponNum].Geometry = roomWeapon.Geometry
	currentRoom.ConsumablesNum++
	currentRoom.Consumables.WeaponNum++
}

func checkUnoccupiedRoom(coordinates *Object, room *Room) bool {
	for i := range room.Consumables.ElixirNum {
		if CheckEqualCoords(coordinates.Coordinates, room.Consumables.Elixirs[i].Geometry.Coordinates) {
			return false
		}
	}

	for i := range room.Consumables.FoodNum {
		if CheckEqualCoords(coordinates.Coordinates, room.Consumables.RoomFood[i].Geometry.Coordinates) {
			return false
		}
	}

	for i := range room.Consumables.ScrollNum {
		if CheckEqualCoords(coordinates.Coordinates, room.Consumables.Scrolls[i].Geometry.Coordinates) {
			return false
		}
	}

	for i := range room.Consumables.WeaponNum {
		if CheckEqualCoords(coordinates.Coordinates, room.Consumables.Weapons[i].Geometry.Coordinates) {
			return false
		}
	}

	for i := range room.MonsterNum {
		if CheckEqualCoords(coordinates.Coordinates, room.Monsters[i].BaseStats.Coords.Coordinates) {
			return false
		}
	}

	return true
}

func UseConsumable(player *Player, consumableType ConsumableTypes, consumablePos int, room *Room, adapter *DifficultyAdapter) {
	var oldWeapon Weapon
	switch consumableType {
	case ScrollConTyp:
		readScroll(player, &player.Backpack.Scrolls[consumablePos])
		removeFromBackpack(&player.Backpack, consumablePos, ScrollConTyp)
		player.Backpack.CurrentSize--
	case ElixirConTyp:
		drinkElixir(player, &player.Backpack.Elixirs[consumablePos])
		removeFromBackpack(&player.Backpack, consumablePos, ElixirConTyp)
		player.Backpack.CurrentSize--
	case FoodConTyp:
		eatFood(player, &player.Backpack.Foods[consumablePos], adapter)
		removeFromBackpack(&player.Backpack, consumablePos, FoodConTyp)
		player.Backpack.CurrentSize--
	case WeaponConTyp:
		if consumablePos == -1 {
			player.Weapon.Strength = noWeapon
		} else if room != nil && !equalWeapons(&player.Backpack.Weapons[consumablePos], &player.Weapon) && player.Weapon.Strength != noWeapon {
			oldWeapon = player.Weapon
			player.Weapon = player.Backpack.Weapons[consumablePos]
			throwCurrentWeapon(player, room, oldWeapon)
		} else if player.Weapon.Strength == noWeapon {
			player.Weapon = player.Backpack.Weapons[consumablePos]
			removeFromBackpack(&player.Backpack, consumablePos, WeaponConTyp)
			player.Backpack.CurrentSize--
		} else {
			player.Weapon = player.Backpack.Weapons[0]
		}
	default:
	}
}

func eatFood(player *Player, food *Food, adapter *DifficultyAdapter) {
	healthBefore := player.BaseStats.Health

	if player.BaseStats.Health+float64(food.ToRegen) >= float64(player.RegenLimit) {
		player.BaseStats.Health = float64(player.RegenLimit)
	} else {
		player.BaseStats.Health += float64(food.ToRegen)
	}

	if adapter != nil {
		healthRestored := int(player.BaseStats.Health - healthBefore)
		adapter.TotalHealthUsed += healthRestored
	}
}

func drinkElixir(player *Player, elixir *Elixir) {
	currentTime := time.Now()

	switch elixir.Stat {
	case Health:
		player.ElixirBuffs.MaxHealth[player.ElixirBuffs.CurrentHealthBuffNum].StatIncrease += elixir.Increase
		player.ElixirBuffs.MaxHealth[player.ElixirBuffs.CurrentHealthBuffNum].EffectEnd = time.Duration(currentTime.Unix())*time.Second + elixir.Duration
		player.RegenLimit += elixir.Increase
		player.BaseStats.Health += float64(elixir.Increase)
		player.ElixirBuffs.CurrentHealthBuffNum++

	case Agility:
		player.ElixirBuffs.Agility[player.ElixirBuffs.CurrentAgilityBuffNum].StatIncrease += elixir.Increase
		player.ElixirBuffs.Agility[player.ElixirBuffs.CurrentAgilityBuffNum].EffectEnd = time.Duration(currentTime.Unix())*time.Second + elixir.Duration
		player.BaseStats.Agility += elixir.Increase
		player.ElixirBuffs.CurrentAgilityBuffNum++

	case Strength:
		player.ElixirBuffs.Strength[player.ElixirBuffs.CurrentStrengthBuffNum].StatIncrease += elixir.Increase
		player.ElixirBuffs.Strength[player.ElixirBuffs.CurrentStrengthBuffNum].EffectEnd = time.Duration(currentTime.Unix())*time.Second + elixir.Duration
		player.BaseStats.Strength += elixir.Increase
		player.ElixirBuffs.CurrentStrengthBuffNum++

	default:
	}
}

func readScroll(player *Player, scroll *Scroll) {
	switch scroll.Stat {
	case Health:
		player.RegenLimit += scroll.Increase
		player.BaseStats.Health += float64(scroll.Increase)

	case Agility:
		player.BaseStats.Agility += scroll.Increase

	case Strength:
		player.BaseStats.Strength += scroll.Increase

	default:
	}
}

func removeFromBackpack(playersBackpack *Backpack, pos int, consumableType ConsumableTypes) {
	switch consumableType {
	case ScrollConTyp:
		playersBackpack.Scrolls[pos] = playersBackpack.Scrolls[playersBackpack.ScrollNum-1]
		playersBackpack.ScrollNum--

	case ElixirConTyp:
		playersBackpack.Elixirs[pos] = playersBackpack.Elixirs[playersBackpack.ElixirNum-1]
		playersBackpack.ElixirNum--

	case FoodConTyp:
		playersBackpack.Foods[pos] = playersBackpack.Foods[playersBackpack.FoodNum-1]
		playersBackpack.FoodNum--

	default:
	}
}

func CheckTempEffectEnd(player *Player) {
	currentTime := time.Duration(time.Now().Unix()) * time.Second

	for i := 0; i < player.ElixirBuffs.CurrentHealthBuffNum; {
		if player.ElixirBuffs.MaxHealth[i].EffectEnd > currentTime {
			i++
		} else {
			player.RegenLimit -= player.ElixirBuffs.MaxHealth[i].StatIncrease
			player.BaseStats.Health -= float64(player.ElixirBuffs.MaxHealth[i].StatIncrease)
			if player.BaseStats.Health <= 0 {
				player.BaseStats.Health = 1
			}
			player.ElixirBuffs.MaxHealth[i] = player.ElixirBuffs.MaxHealth[player.ElixirBuffs.CurrentHealthBuffNum-1]
			player.ElixirBuffs.CurrentHealthBuffNum--
		}
	}

	for i := 0; i < player.ElixirBuffs.CurrentAgilityBuffNum; {
		if player.ElixirBuffs.Agility[i].EffectEnd > currentTime {
			i++
		} else {
			player.BaseStats.Agility -= player.ElixirBuffs.Agility[i].StatIncrease
			player.ElixirBuffs.Agility[i] = player.ElixirBuffs.Agility[player.ElixirBuffs.CurrentAgilityBuffNum-1]
			player.ElixirBuffs.CurrentAgilityBuffNum--
		}
	}

	for i := 0; i < player.ElixirBuffs.CurrentStrengthBuffNum; {
		if player.ElixirBuffs.Strength[i].EffectEnd > currentTime {
			i++
		} else {
			player.BaseStats.Strength -= player.ElixirBuffs.Strength[i].StatIncrease
			player.ElixirBuffs.Strength[i] = player.ElixirBuffs.Strength[player.ElixirBuffs.CurrentStrengthBuffNum-1]
			player.ElixirBuffs.CurrentStrengthBuffNum--
		}
	}
}

func equalWeapons(first *Weapon, second *Weapon) bool {
	return first.Name == second.Name && first.Strength == second.Strength
}
