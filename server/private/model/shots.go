package model

type Shot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Shots []Shot

type ShotsRequest struct {
	Shots `json:"shots"`
}

type ShotsResponse struct {
	Shots `json:"shots"`
}
