package main

import (
	"rogue/internal/domain"
	"rogue/internal/models"
	"rogue/internal/presentation"

	"github.com/gbin/goncurses"
)

const (
	SavePath       = "./internal/models/data/save.json"
	StatisticsPath = "./internal/models/data/statistics.json"
	ScoreboardPath = "./internal/models/data/scoreboard.json"
)

func main() {
	stdscr, err := goncurses.Init()
	if err != nil {
		panic(err)
	}
	defer goncurses.End()

	presentation.InitPresentation(stdscr)
	game := domain.NewGame()
	battles := domain.NewBattles()

	presentation.StartScreen(stdscr)

	runningMenu := true
	currentOption := 0
	for runningMenu {
		presentation.MenuScreen(currentOption, stdscr)
		key := stdscr.GetChar()

		switch key {
		case '\n': // Enter
			switch currentOption {
			case 0: // New Game
				models.InitLevel(game, battles, StatisticsPath)
				if game.DifficultyAdapter == nil {
					game.DifficultyAdapter = domain.NewDifficultyAdapter()
				}
				presentation.GameCycle(game.Player, game.Level, game.Map, battles[:], game.DifficultyAdapter, SavePath, ScoreboardPath, StatisticsPath, stdscr)
			case 1: // Load Game
				models.LoadData(game, battles, SavePath, StatisticsPath)
				if game.DifficultyAdapter == nil {
					game.DifficultyAdapter = domain.NewDifficultyAdapter()
				}
				presentation.GameCycle(game.Player, game.Level, game.Map, battles[:], game.DifficultyAdapter, SavePath, ScoreboardPath, StatisticsPath, stdscr)
			case 2: // Scoreboard
				presentation.DisplayScoreboard(ScoreboardPath, stdscr)
				for stdscr.GetChar() != presentation.Escape {
				}
			case 3: // Exit
				runningMenu = false
			}
		case 'W', 'w':
			currentOption = max(0, currentOption-1)
		case 'S', 's':
			currentOption = min(3, currentOption+1)
		}
	}
}
