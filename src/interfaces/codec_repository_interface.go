package interfaces

type CodecRepositoryInterface interface {
	GetCodecs() map[string]struct {
		Containers []string
		Encoders   []string
	}
	GetContainers() map[string]struct {
		Containers []string
	}
	GetEncoders() map[string]struct {
		Presets []string
		Tune    []string
		Profile []string
		Level   []string
	}
}
