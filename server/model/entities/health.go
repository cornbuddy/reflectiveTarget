package entities

type HealthResponse struct {
	Connected   bool `json:"connected"`
	Connections int  `json:"connections"`
}
