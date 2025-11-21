package valueobjects

type Shot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Shots []Shot
