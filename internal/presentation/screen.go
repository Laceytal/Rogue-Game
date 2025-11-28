package presentation

import (
	"github.com/gbin/goncurses"
)

func StartScreen(stdscr *goncurses.Window) {
	stdscr.Clear()

	strings := []string{
		"     _____       _______       _____       _____       _____    ",
		"    /\\    \\     /::\\    \\     /\\    \\     /\\    \\     /\\    \\   ",
		"   /::\\    \\   /::::\\    \\   /::\\    \\   /::\\____\\   /::\\    \\  ",
		"  /::::\\    \\ /::::::\\    \\ /::::\\    \\ /:::/ /     /::::\\    \\ ",
		" /::::::\\    /::::::::\\    /::::::\\    /:::/ /     /::::::\\    \\",
		"/:::/\\:::\\  /:::/~~\\:::\\  /:::/\\:::\\  /:::/ /     /:::/\\:::\\   ",
		"/:::/__\\:::\\ /:::/ \\:::\\ /:::/ \\:::\\ /:::/ /     /:::/__\\:::\\  ",
		"/::::\\   \\:::/:::/ / \\:::/:::/ / \\:::/:::/ / _____ /::::::\\   \\:::",
		"/::::::\\   \\:::/____/   \\:::\\____/:/ / \\:::\\__/:::/____/ /\\   /:::/\\:::",
		"/:::/\\:::\\   \\:::|    | |:::| | /:::/ / ___\\:::| ||:::| / /::\\__/:::/",
		"/:::/ \\:::\\  /:::|____| |:::| |/:::/____/ /\\ /:::|____||:::|____\\ /:::/ /\\",
		"\\::/ |::::\\ /:::/ / \\:::\\ \\ /:::/ / \\::\\ /::\\ \\::/ / \\:::\\   /:::/ / \\:::",
		" \\/____|:::::\\/:::/ /   \\:::\\ /:::/ /   \\::/::\\ \\/____/ \\:::\\  /:::/ /  \\:::",
		"       |:::::::::/ /     \\:::/:::/ /     \\:::\\            \\:::\\ /:::/ /   \\:::",
		"       |::|\\::::/  /      \\::/ /        \\:::\\____\\        \\:::/:::/ /    \\:::",
		"       |::| \\::/____/       \\/__/         /:::/ /         \\::/ /      \\:::",
		"       |::|  ~|                          /:::/ /          \\/____/       \\:::",
		"       |::|   |                         /:::/ /                         \\:::",
		"       \\::|   |                        /:::/ /                          \\:::",
		"        \\:|   |                        \\::/____/                         \\:::/",
		"         \\|___|                         ~~                               \\/____/",
	}

	row, col := stdscr.MaxYX()
	artWidth := len(strings[0])
	artHeight := len(strings)
	shiftX := (col - artWidth) / 2
	shiftY := (row - artHeight) / 2

	for i, line := range strings {
		stdscr.MovePrint(shiftY+i, shiftX, line)
	}

	continueText := "Press any key to continue..."
	textX := (col - len(continueText)) / 2
	stdscr.MovePrint(shiftY+artHeight+1, textX, continueText)

	stdscr.GetChar()
	stdscr.Clear()
}
func MenuScreen(currentLine int, stdscr *goncurses.Window) {
	stdscr.Clear()

	strings := []string{
		"           GAME MENU        ",
		"+------------------------------+",
		"|                              |",
		"|          NEW GAME            |",
		"|          LOAD GAME           |",
		"|          SCOREBOARD          |",
		"|          EXIT GAME           |",
		"|                              |",
		"+------------------------------+",
	}

	row, col := stdscr.MaxYX()
	menuWidth := len(strings[0])
	menuHeight := len(strings)
	shiftX := (col - menuWidth) / 2
	shiftY := (row - menuHeight) / 2

	for i, line := range strings {
		stdscr.MovePrint(shiftY+i, shiftX, line)
	}

	stdscr.MovePrint(shiftY+currentLine+3, shiftX+5, "<<<")
	stdscr.MovePrint(shiftY+currentLine+3, shiftX+24, ">>>")
}

func DeadScreen(stdscr *goncurses.Window) {
	stdscr.Clear()

	strings := []string{
		"     ___         ___         ___         ___    ",
		"    /\\  \\       /\\  \\   ___     /\\__\\       /\\  \\   ___  ",
		"   ___ /::\\  \\       \\:\\  \\     /::\\  \\     /:/ _/_     /::\\  \\     /::\\  \\  ",
		"  /| | /:/\\:\\  \\       \\:\\  \\   /:/\\:\\  \\   /:/ /\\__\\   /:/\\:\\  \\   /:/\\:\\  \\ ",
		" |:| |/:/ \\:\\  \\   ___  \\:\\  \\ /:/ \\:\\\\__\\ /:/ /:/ _/_ /:/ /::\\  \\ /:/ \\:\\\\__\\",
		" |:| |/:/__/ \\:\\\\__\\ /\\  \\ \\:\\\\__/:/__/ \\:|__/:/_/:/ /\\__/:/_/:/\\:\\\\__/:/__/ \\:|__|",
		" __|:|__|\\:\\  \\ /:/ / \\:\\  \\ /:/ / \\:\\  \\ /:/ / \\:\\/:/ /:/ / \\:\\\\/:/ \\/__\\:\\  \\ /:/ /",
		"/:::\\  \\ \\:\\  /:/ /   \\:\\  /:/ /   \\:\\  /:/ /   \\::/_/:/ /   \\::/__/     \\:\\  /:/ / ",
		"~~~~\\:\\  \\ \\:\\\\/:/ /     \\:\\\\/:/ /     \\:\\\\/:/ /     \\:\\\\/:/ /     \\:\\  \\      \\:\\\\/:/ / ",
		"     \\:\\\\__\\ \\::/ /       \\::/ /       \\::/ /       \\::/ /       \\:\\\\__\\      \\::/ / ",
		"      \\/__/  \\/__/         \\/__/         \\/__/         \\/__/         \\/__/       \\/__/ ",
	}

	row, col := stdscr.MaxYX()
	artWidth := len(strings[0])
	artHeight := len(strings)
	shiftX := (col - artWidth) / 2
	shiftY := (row - artHeight) / 2

	for i, line := range strings {
		stdscr.MovePrint(shiftY+i, shiftX, line)
	}

	continueText := "Press any key to continue..."
	textX := (col - len(continueText)) / 2
	stdscr.MovePrint(shiftY+artHeight+1, textX, continueText)

	stdscr.GetChar()
	stdscr.Clear()
}

func EndgameScreen(stdscr *goncurses.Window) {
	stdscr.Clear()

	strings := []string{
		"     ___         ___         ___         ___         ___         ___    ",
		"    /\\__\\       /\\  \\   ___     /\\__\\       /\\  \\       /\\  \\       /\\__\\   ",
		"   /:/ _/_     \\:\\  \\     /::\\  \\     /:/ _/_     /::\\  \\     |::\\  \\     /:/ _/_  ",
		"  /:/ /\\__\\     \\:\\  \\   /:/\\:\\  \\   /:/ /\\ \\   /:/\\:\\  \\   |:|:\\  \\   /:/ /\\__\\ ",
		" /:/ /:/ _/_ _____\\:\\  \\ /:/ \\:\\\\__\\ /:/ /::\\ \\ /:/ /::\\  \\ __|:|\\:\\  \\ /:/ /:/ _/_ ",
		"/:/_/:/ /\\__/::::::::\\\\__/:/__/ \\:|__/:/__\\/\\:\\\\__/:/_/:/\\:\\\\__/::::|_\\:\\\\__/:/_/:/ /\\__\\",
		"\\:\\/:/ /:/ /\\:\\\\~~\\~~\\/__\\:\\  \\ /:/ / \\:\\  \\ /:/ / \\:\\ \\ /:/ / \\:\\\\~~\\ \\/__\\:\\/:/ /:/ /",
		" \\::/_/:/ /  \\:\\  \\       \\:\\  /:/ /   \\:\\  /:/ /   \\:\\ /:/ /   \\:\\  \\       \\::/_/:/ / ",
		"  \\:\\/:/ /    \\:\\  \\       \\:\\\\/:/ /     \\:\\\\/:/ /     \\:\\\\/:/ /     \\:\\  \\       \\:\\/:/ / ",
		"   \\::/ /      \\:\\\\__\\      \\::/ /       \\::/ /       \\::/ /       \\:\\\\__\\      \\::/ / ",
		"    \\/__/        \\/__/       \\/__/         \\/__/         \\/__/         \\/__/       \\/__/ ",
	}

	row, col := stdscr.MaxYX()
	artWidth := len(strings[0])
	artHeight := len(strings)
	shiftX := (col - artWidth) / 2
	shiftY := (row - artHeight) / 2

	for i, line := range strings {
		stdscr.MovePrint(shiftY+i, shiftX, line)
	}

	continueText := "Press any key to continue..."
	textX := (col - len(continueText)) / 2
	stdscr.MovePrint(shiftY+artHeight+1, textX, continueText)

	stdscr.GetChar()
	stdscr.Clear()
}
