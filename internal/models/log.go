package models

type Log struct {
	ID        uint   `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	Message   string `json:"message"`
}
