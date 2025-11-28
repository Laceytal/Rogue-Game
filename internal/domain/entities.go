package domain

import (
	"time"
)

const (
	RoomsInWidth  = 3
	RoomsInHeight = 3
	RoomsNum      = RoomsInWidth * RoomsInHeight

	MaxPassageParts = 3
	MaxPassagesNum  = (RoomsNum - 1) * MaxPassageParts
)

const (
	RegionWidth   = 27
	RegionHeight  = 10
	minRoomWidth  = 6
	maxRoomWidth  = RegionWidth - 3
	minRoomHeight = 5
	maxRoomHeight = RegionHeight - 3
)

const (
	ConsumablesTypeMaxNum = 9
)

const (
	noWeapon   = 0
	MaxNameLen = 32 + 1
)

const (
	maxPercentFoodRegenFromHealth = 20
	maxPercentAgilityIncrease     = 10
	maxPercentStrengthIncrease    = 10
)

const (
	minElixirDurationSeconds = 30
	maxElixirDurationSeconds = 60
)

const (
	minWeaponStrength = 30
	maxWeaponStrength = 50
)

const (
	lowHostilityRadius     = 2
	averageHostilityRadius = 4
	highHostilityRadius    = 6
)

const (
	MaxMonstersPerRoom    = 2
	MaxConsumablesPerRoom = 3
)

const (
	levelUpdateDifficulty            = 10
	percentsUpdateDifficultyMonsters = 2
)

const (
	LevelNum = 21
)

type Dimension int

const (
	X Dimension = iota
	Y
	CoordsNum
)

type MonsterType int

const (
	Zombie MonsterType = iota
	Vampire
	Ghost
	Ogre
	Snake
	Mimic
	MonsterTypeNum
)

type hostilityType int

const (
	low hostilityType = iota
	average
	high
)

type StatType int

const (
	Health StatType = iota
	Agility
	Strength
	StatTypeNum
)

type Directions int

const (
	Forward Directions = iota
	Back
	Left
	Right
	diagonallyForwardLeft
	diagonallyForwardRight
	diagonallyBackLeft
	diagonallyBackRight
	stop
)

type Coordinates [CoordsNum]int
type Sizes [CoordsNum]int

type Object struct {
	Coordinates Coordinates
	Size        Sizes
}

type character struct {
	Coords   Object
	Health   float64
	Agility  int
	Strength int
}

type Monster struct {
	BaseStats character
	Type      MonsterType
	hostility hostilityType
	IsChasing bool
	dir       Directions
}

type Treasure struct {
	Value int
}

type Food struct {
	ToRegen int
	Name    [MaxNameLen]rune
}

type Elixir struct {
	Duration time.Duration
	Stat     StatType
	Increase int
	Name     [MaxNameLen]rune
}

type Scroll struct {
	Stat     StatType
	Increase int
	Name     [MaxNameLen]rune
}

type Weapon struct {
	Strength int
	Name     [MaxNameLen]rune
}

type FoodRoom struct {
	Food     Food
	Geometry Object
}

type ElixirRoom struct {
	Elixir   Elixir
	Geometry Object
}

type ScrollRoom struct {
	Scroll   Scroll
	Geometry Object
}

type WeaponRoom struct {
	Weapon   Weapon
	Geometry Object
}

type ConsumablesRoom struct {
	RoomFood [MaxConsumablesPerRoom]FoodRoom
	FoodNum  int

	Elixirs   [MaxConsumablesPerRoom]ElixirRoom
	ElixirNum int

	Scrolls   [MaxConsumablesPerRoom]ScrollRoom
	ScrollNum int

	Weapons   [MaxConsumablesPerRoom + ConsumablesTypeMaxNum]WeaponRoom
	WeaponNum int
}

type Backpack struct {
	CurrentSize int

	Foods   [ConsumablesTypeMaxNum]Food
	FoodNum int

	Elixirs   [ConsumablesTypeMaxNum]Elixir
	ElixirNum int

	Scrolls   [ConsumablesTypeMaxNum]Scroll
	ScrollNum int

	Treasures Treasure

	Weapons   [ConsumablesTypeMaxNum]Weapon
	WeaponNum int
}

type Buff struct {
	StatIncrease int
	EffectEnd    time.Duration
}

type Buffs struct {
	MaxHealth            [ConsumablesTypeMaxNum]Buff
	CurrentHealthBuffNum int

	Agility               [ConsumablesTypeMaxNum]Buff
	CurrentAgilityBuffNum int

	Strength               [ConsumablesTypeMaxNum]Buff
	CurrentStrengthBuffNum int
}

type Player struct {
	BaseStats   character
	RegenLimit  int
	Backpack    Backpack
	Weapon      Weapon
	ElixirBuffs Buffs
}

type Room struct {
	Coords         Object
	Consumables    ConsumablesRoom
	ConsumablesNum int
	Monsters       [MaxMonstersPerRoom]Monster
	MonsterNum     int
}

type Passage = Object

type Passages struct {
	Passages    [MaxPassagesNum]Passage
	PassagesNum int
}

type Level struct {
	Coords     Object
	Rooms      [RoomsNum]Room
	Passages   Passages
	LevelNum   int
	EndOfLevel Object
}

const (
	ChanceUnvisibleGhost = 80
	MapHeight            = RoomsInHeight * RegionHeight
	MapWidth             = RoomsInWidth * RegionWidth
	MaxScoreboardSize    = 10
)

type Font int

const (
	WhiteFont Font = iota + 1
	RedFont
	GreenFont
	BlueFont
	YellowFont
	CyanFont
)

type Cell struct {
	Char  rune
	Color Font
}

type Map struct {
	MapData         [MapHeight][MapWidth]Cell
	VisibleRooms    [RoomsNum]bool
	VisiblePassages [MaxPassagesNum]bool
}

type Game struct {
	Player            *Player            `json:"player"`
	Level             *Level             `json:"level"`
	Map               *Map               `json:"map"`
	DifficultyAdapter *DifficultyAdapter `json:"difficulty_adapter"`
}

func NewMap() *Map {
	return &Map{
		MapData:         [MapHeight][MapWidth]Cell{},
		VisibleRooms:    [RoomsNum]bool{},
		VisiblePassages: [MaxPassagesNum]bool{},
	}
}

func NewPlayer() *Player {
	return &Player{
		BaseStats: character{
			Coords: Object{
				Coordinates: Coordinates{0, 0},
				Size:        Sizes{1, 1},
			},
			Health:   500.0,
			Agility:  70,
			Strength: 70,
		},
		RegenLimit: 100,
		Backpack: Backpack{
			CurrentSize: 0,
			FoodNum:     0,
			ElixirNum:   0,
			ScrollNum:   0,
			WeaponNum:   0,
			Treasures:   Treasure{Value: 0},
		},
		Weapon: Weapon{
			Strength: noWeapon,
			Name:     [MaxNameLen]rune{},
		},
		ElixirBuffs: Buffs{
			CurrentHealthBuffNum:   0,
			CurrentAgilityBuffNum:  0,
			CurrentStrengthBuffNum: 0,
		},
	}
}

func NewLevel(levelNum int) *Level {
	return &Level{
		Coords: Object{
			Coordinates: Coordinates{0, 0},
			Size:        Sizes{MapWidth, MapHeight},
		},
		Rooms:    [RoomsNum]Room{},
		Passages: Passages{PassagesNum: 0},
		LevelNum: levelNum,
		EndOfLevel: Object{
			Coordinates: Coordinates{0, 0},
			Size:        Sizes{1, 1},
		},
	}
}

func NewGame() *Game {
	return &Game{
		Player:            NewPlayer(),
		Level:             NewLevel(1),
		Map:               NewMap(),
		DifficultyAdapter: NewDifficultyAdapter(),
	}
}

func NewRoom() Room {
	return Room{
		Coords: Object{
			Coordinates: Coordinates{0, 0},
			Size:        Sizes{0, 0},
		},
		Consumables: ConsumablesRoom{
			FoodNum:   0,
			ElixirNum: 0,
			ScrollNum: 0,
			WeaponNum: 0,
		},
		ConsumablesNum: 0,
		MonsterNum:     0,
	}
}

func NewMonster(monsterType MonsterType, hostility hostilityType) Monster {
	return Monster{
		BaseStats: character{
			Coords: Object{
				Coordinates: Coordinates{0, 0},
				Size:        Sizes{1, 1},
			},
			Health:   50.0,
			Agility:  5,
			Strength: 5,
		},
		Type:      monsterType,
		hostility: hostility,
		IsChasing: false,
		dir:       stop,
	}
}

func NewBackpack() Backpack {
	return Backpack{
		CurrentSize: 0,
		FoodNum:     0,
		ElixirNum:   0,
		ScrollNum:   0,
		WeaponNum:   0,
		Treasures:   Treasure{Value: 0},
	}
}

func NewBattleInfo() BattleInfo {
	return BattleInfo{
		IsFight:            false,
		Enemy:              nil,
		VampireFirstAttack: false,
		OgreCooldown:       false,
		PlayerAsleep:       false,
	}
}

func NewBattles() *[MaximumFights]BattleInfo {
	var battles [MaximumFights]BattleInfo
	for i := range battles {
		battles[i] = NewBattleInfo()
	}
	return &battles
}

func InitPlayer(player *Player) {
	player.BaseStats = character{
		Coords: Object{
			Coordinates: Coordinates{0, 0},
			Size:        Sizes{1, 1},
		},
		Health:   500,
		Agility:  70,
		Strength: 70,
	}
	player.RegenLimit = 500
	player.Backpack = Backpack{
		CurrentSize: 0,
		FoodNum:     0,
		ElixirNum:   0,
		ScrollNum:   0,
		WeaponNum:   0,
		Treasures:   Treasure{Value: 0},
	}
	player.Weapon = Weapon{
		Strength: noWeapon,
		Name:     [MaxNameLen]rune{},
	}
	player.ElixirBuffs = Buffs{
		CurrentHealthBuffNum:   0,
		CurrentAgilityBuffNum:  0,
		CurrentStrengthBuffNum: 0,
	}
}

func InitLevel(level *Level, levelNum int) {
	level.Coords = Object{
		Coordinates: Coordinates{0, 0},
		Size:        Sizes{MapWidth, MapHeight},
	}
	level.Rooms = [RoomsNum]Room{}
	level.Passages = Passages{PassagesNum: 0}
	level.LevelNum = levelNum
	level.EndOfLevel = Object{
		Coordinates: Coordinates{0, 0},
		Size:        Sizes{1, 1},
	}
}

func InitMap(mapData *Map) {
	mapData.MapData = [MapHeight][MapWidth]Cell{}
	mapData.VisibleRooms = [RoomsNum]bool{}
	mapData.VisiblePassages = [MaxPassagesNum]bool{}
}

func InitBattles(battles *[MaximumFights]BattleInfo) {
	for i := range MaximumFights {
		battles[i] = BattleInfo{
			IsFight:            false,
			Enemy:              nil,
			VampireFirstAttack: false,
			OgreCooldown:       false,
			PlayerAsleep:       false,
		}
	}
}
