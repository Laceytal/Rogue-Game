package models

import (
	"encoding/json"
	"log"
	"os"
	"rogue/internal/domain"
)

type SessionStat struct {
	Treasures int `json:"treasures"`
	Level     int `json:"level"`
	Enemies   int `json:"enemies"`
	Food      int `json:"food"`
	Elixirs   int `json:"elixirs"`
	Scrolls   int `json:"scrolls"`
	Attacks   int `json:"attacks"`
	Missed    int `json:"missed"`
	Movies    int `json:"movies"`
}

func GetDataFromFile[T SessionStat | []SessionStat | domain.Game](filename string) (*T, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	var data T
	err = decoder.Decode(&data)
	return &data, err
}

func SaveDataToFile[T SessionStat | []SessionStat | domain.Game](filename string, stat T) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("SaveDataToFile: %v", err)
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")
	return encoder.Encode(stat)
}

func UpdateStatistic(pathScoreboard, pathStat string) error {
	scoreboard := make([]SessionStat, 0)
	array, _ := GetDataFromFile[[]SessionStat](pathScoreboard)
	if array == nil {
		array = &scoreboard
	}
	stat, err := GetDataFromFile[SessionStat](pathStat)
	if err != nil {
		return err
	}
	newArray := append(*array, *stat)
	return SaveDataToFile(pathScoreboard, newArray)
}
