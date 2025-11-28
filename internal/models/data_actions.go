package models

import (
	"log"
	"rogue/internal/domain"
	"time"
)

func LoadData(game *domain.Game, battles *[domain.MaximumFights]domain.BattleInfo, savePath string, statPath string) {
	if game == nil || battles == nil || savePath == "" || statPath == "" {
		log.Fatalln("Load Data: Game or Battles and SavePath are nil.")
	}

	for i := range domain.MaximumFights {
		(*battles)[i].IsFight = false
	}

	gameFromFile, err := GetDataFromFile[domain.Game](savePath)
	if err != nil {
		InitLevel(game, battles, statPath)
	} else {
		*game = *gameFromFile

		if game.DifficultyAdapter != nil {
			game.DifficultyAdapter.LastLevelStartTime = time.Now()
		} else {
			game.DifficultyAdapter = domain.NewDifficultyAdapter()
		}
	}
}
func InitLevel(game *domain.Game, battles *[domain.MaximumFights]domain.BattleInfo, statPath string) {
	if game == nil || battles == nil || statPath == "" {
		log.Fatalln("Init Level: Game or Battles and SavePath are nil.")
	}
	domain.InitPlayer(game.Player)

	domain.InitBattles(battles)

	domain.InitLevel(game.Level, 0)

	domain.GenerateNextLevel(game.Level, game.Player)

	domain.InitMap(game.Map)

	stat := *new(SessionStat)
	err := SaveDataToFile(statPath, stat)
	if err != nil {
		log.Fatalln("Init Level:", err)
	}
}

func GetStandartSave(savePath string, statPath string) {
	game := domain.NewGame()
	game.Player = domain.NewPlayer()
	game.Level = domain.NewLevel(0)
	game.Map = domain.NewMap()
	game.DifficultyAdapter = domain.NewDifficultyAdapter()

	battles := domain.NewBattles()
	InitLevel(game, battles, statPath)

	err := SaveDataToFile(savePath, *game)
	if err != nil {
		log.Fatalln("GetStandartSave:", err)
	}
}
