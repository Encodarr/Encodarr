package models

type Codec struct {
	Containers []string `json:"containers"`
	Encoders   []string `json:"encoders"`
}

type Container struct {
	Extensions []string `json:"extensions"`
}

type Encoder struct {
	Presets []string `json:"presets"`
	Tune    []string `json:"tune"`
	Profile []string `json:"profile"`
	Level   []string `json:"level"`
}
