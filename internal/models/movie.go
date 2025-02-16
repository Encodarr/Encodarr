package models

type Movie struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	ReleaseDate  string `json:"releaseDate"`
	Genre        string `json:"genre"`
	Status       string `json:"status"`
	Filename     string `json:"filename"`
	VideoCodec   string `json:"videoCodec"`
	Overview     string `json:"overview"`
	Size         int    `json:"size"`
	SpaceSaved   int    `json:"spaceSaved"`
	ProfileID    int    `json:"profileId"`
	Monitored    bool   `json:"monitored"`
	Missing      bool   `json:"missing"`
	Studio       string `json:"studio"`
	OriginalSize int    `json:"originalSize"`
	Path         string `json:"path"`
	Runtime      int    `json:"runtime"`
}
