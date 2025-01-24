package models

// Define a struct to hold the return values
type QueueStatus struct {
	Queue       []Item  `json:"queue"`
	Progress    float64 `json:"progress"`
	Stage       string  `json:"stage"`
	ETA         int     `json:"eta"`
	Current     Item    `json:"current"`
	QueueStatus string  `json:"queueStatus"`
}
