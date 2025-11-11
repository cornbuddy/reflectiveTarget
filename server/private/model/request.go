package model

type ShotsRequest struct {
	Shots []Shot `json:"shots"`
}

type Shot struct {
	X int `json:"x"`
	Y int `json:"y"`
}
