package model

type ShotsRequest struct {
	Shots `json:"shots"`
}

type Shot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Shots []Shot
