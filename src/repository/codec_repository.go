package repository

type CodecRepository struct{}

func NewCodecRepository() *CodecRepository {
	return &CodecRepository{}
}

func (repo *CodecRepository) GetCodecs() map[string]struct {
	Containers []string
	Encoders   []string
} {
	return map[string]struct {
		Containers []string
		Encoders   []string
	}{
		"Any":   {Containers: []string{}, Encoders: []string{}},
		"h264":  {Containers: []string{"mp4", "matroska"}, Encoders: []string{"libx264", "h264"}},
		"hevc":  {Containers: []string{"mp4", "matroska"}, Encoders: []string{"libx265"}},
		"mpeg4": {Containers: []string{"mp4", "matroska"}, Encoders: []string{"mpeg4"}},
		"vp8":   {Containers: []string{"mp4", "matroska"}, Encoders: []string{"libvpx"}},
		"vp9":   {Containers: []string{"mp4", "matroska"}, Encoders: []string{"libvpx-vp9"}},
		"av1":   {Containers: []string{"mp4", "matroska"}, Encoders: []string{"libaom-av1"}},
	}
}

func (repo *CodecRepository) GetContainers() map[string]struct {
	Containers []string
} {
	return map[string]struct {
		Containers []string
	}{
		"mp4":      {Containers: []string{"mp4", "m4a", "m4v", "f4v", "f4a", "m4b", "m4r", "f4b", "mov"}},
		"matroska": {Containers: []string{"mkv", "mk3d", "mka", "mks"}},
	}
}

func (repo *CodecRepository) GetEncoders() map[string]struct {
	Presets []string
	Tune    []string
	Profile []string
	Level   []string
} {
	return map[string]struct {
		Presets []string
		Tune    []string
		Profile []string
		Level   []string
	}{
		"h264": {
			Presets: []string{"ultrafast", "superfast", "veryfast", "faster", "fast", "medium", "slow", "slower", "veryslow", "placebo"},
			Tune:    []string{"none", "film", "animation", "grain", "stillimage", "psnr", "ssim", "zerolatency"},
			Profile: []string{"auto", "baseline", "main", "high", "high422", "high444"},
			Level:   []string{"auto", "1.0", "1b", "1.1", "1.2", "1.3", "2.0", "2.1", "2.2", "3.0", "3.1", "3.2", "4.0", "4.1", "4.2", "5.0", "5.1", "5.2", "6.0", "6.1", "6.2"},
		},
		"libx264": {
			Presets: []string{"ultrafast", "superfast", "veryfast", "faster", "fast", "medium", "slow", "slower", "veryslow", "placebo"},
			Tune:    []string{"none", "film", "animation", "grain", "stillimage", "psnr", "ssim", "zerolatency"},
			Profile: []string{"auto", "baseline", "main", "high", "high422", "high444"},
			Level:   []string{"auto", "1.0", "1b", "1.1", "1.2", "1.3", "2.0", "2.1", "2.2", "3.0", "3.1", "3.2", "4.0", "4.1", "4.2", "5.0", "5.1", "5.2", "6.0", "6.1", "6.2"},
		},
		"libx265": {
			Presets: []string{"ultrafast", "superfast", "veryfast", "faster", "fast", "medium", "slow", "slower", "veryslow", "placebo"},
			Tune:    []string{"none", "psnr", "ssim", "grain", "zerolatency", "fastdecode", "animation"},
			Profile: []string{"auto", "main", "main10", "mainstillpicture", "main444-8", "main444-intra"},
			Level:   []string{"auto", "1.0", "2.0", "2.1", "3.0", "3.1", "4.0", "4.1", "5.0", "5.1", "5.2", "6.0", "6.1", "6.2"},
		},
		"mpeg4": {
			Presets: []string{},
			Tune:    []string{},
			Profile: []string{},
			Level:   []string{},
		},
		"libvpx": {
			Presets: []string{"ultrafast", "superfast", "veryfast", "faster", "fast", "medium", "slow", "slower", "veryslow", "placebo"},
			Tune:    []string{"none"},
			Profile: []string{"auto"},
			Level:   []string{"auto"},
		},
		"libvpx-vp9": {
			Presets: []string{"ultrafast", "superfast", "veryfast", "faster", "fast", "medium", "slow", "slower", "veryslow", "placebo"},
			Tune:    []string{"none"},
			Profile: []string{"auto"},
			Level:   []string{"auto"},
		},
		"libaom-av1": {
			Presets: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
			Tune:    []string{"none", "psnr", "ssim"},
			Profile: []string{"auto", "main"},
			Level:   []string{"auto", "2.0", "2.1", "2.2", "3.0", "3.1", "4.0", "4.1", "5.0", "5.1", "5.2", "6.0", "6.1", "6.2"},
		},
	}
}
