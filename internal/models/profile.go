package models

type Profile struct {
	Id                        int    `json:"id"`
	Name                      string `json:"name"`
	Container                 string `json:"container"`
	Extension                 string `json:"extension"`
	PassThruCommonMetadata    bool   `json:"passThruCommonMetadata"`
	Flipping                  bool   `json:"flipping"`
	Rotation                  int    `json:"rotation"`
	Cropping                  string `json:"cropping"`
	Limit                     string `json:"limit"`
	Anamorphic                string `json:"anamorphic"`
	Fill                      string `json:"fill"`
	Color                     string `json:"color"`
	Detelecine                string `json:"detelecine"`
	InterlaceDetection        string `json:"interlaceDetection"`
	Deinterlace               string `json:"deinterlace"`
	DeinterlacePreset         string `json:"deinterlacePreset"`
	Deblock                   string `json:"deblock"`
	DeblockTune               string `json:"deblockTune"`
	Denoise                   string `json:"denoise"`
	DenoisePreset             string `json:"denoisePreset"`
	DenoiseTune               string `json:"denoiseTune"`
	ChromaSmooth              string `json:"chromaSmooth"`
	ChromaSmoothTune          string `json:"chromaSmoothTune"`
	Sharpen                   string `json:"sharpen"`
	SharpenPreset             string `json:"sharpenPreset"`
	SharpenTune               string `json:"sharpenTune"`
	Colorspace                string `json:"colorspace"`
	Grayscale                 bool   `json:"grayscale"`
	Codec                     string `json:"codec"`
	Encoder                   string `json:"encoder"`
	Framerate                 string `json:"framerate"`
	FramerateType             string `json:"framerateType"`
	QualityType               string `json:"qualityType"`
	ConstantQuality           int    `json:"constantQuality"`
	AverageBitrate            int    `json:"averageBitrate"`
	MultipassEncoding         bool   `json:"multipassEncoding"`
	Preset                    string `json:"preset"`
	Tune                      string `json:"tune"`
	Profile                   string `json:"profile"`
	Level                     string `json:"level"`
	FastDecode                bool   `json:"fastDecode"`
	MapUntaggedAudioTracks    bool   `json:"mapUntaggedAudioTracks"`
	MapUntaggedSubtitleTracks bool   `json:"mapUntaggedSubtitleTracks"`

	ProfileAudioLanguages    []ProfileAudioLanguage    `json:"profileAudioLanguages"`
	ProfileSubtitleLanguages []ProfileSubtitleLanguage `json:"profileSubtitleLanguages"`
	ProfileCodecs            []ProfileCodec            `json:"codecs"`
}
