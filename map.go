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
		Angle float64 `json:"angle"`
	} `json:"player"`
	Const struct {
		ObjectSize    int32   `json:"object_size"`
		ObjectSpeed   float32 `json:"object_speed"`
		BulletSpeed   float32 `json:"bullet_speed"`
		RotationSpeed float64 `json:"rotation_speed"`
	} `json:"world_constants"`
	Asteroids []struct {
		Coordinates struct {
			X float32 `json:"x"`
			Y float32 `json:"y"`
		} `json:"coordinates"`
		Radius float32 `json:"radius"`
	} `json:"asteroids"`
	Walls []struct {
		Start struct {
			X float32 `json:"x"`
			Y float32 `json:"y"`
		} `json:"start"`
		End struct {
			X float32 `json:"x"`
			Y float32 `json:"y"`
		} `json:"end"`
	} `json:"walls"`
}

func (m *Map) LoadMap() {
	file, _ := ioutil.ReadFile("assets/map.json")
	err := json.Unmarshal(file, &m)
	if err != nil {
		log.Print(err)
	}
}
