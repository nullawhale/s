package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Map struct {
	World struct {
		ScreenWidth  int32 `json:"screen_width"`
		ScreenHeight int32 `json:"screen_height"`
	} `json:"world"`
	Player struct {
		Coordinates struct {
			X float32 `json:"x"`
			Y float32 `json:"y"`
		} `json:"coordinates"`
		Angle float32 `json:"angle"`
	} `json:"player"`
	Objects []struct {
		Coordinates struct {
			X float32 `json:"x"`
			Y float32 `json:"y"`
		} `json:"coordinates"`
		Radius float32 `json:"radius"`
	} `json:"objects"`
}

func (m *Map) LoadMap() {
	file, _ := ioutil.ReadFile("assets/map.json")
	err := json.Unmarshal(file, &m)
	if err != nil {
		log.Print(err)
	}
}
