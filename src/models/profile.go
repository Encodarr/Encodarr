package models

type Profile struct {
	BaseModel
	Id                        int    `gorm:"primary_key;autoIncrement" json:"id"`
	Name                      string `gorm:"type:varchar(255)" json:"name"`
	Container                 string `gorm:"type:varchar(255)" json:"container"`
	Extension                 string `gorm:"type:varchar(255)" json:"extension"`
	PassThruCommonMetadata    bool   `gorm:"type:boolean" json:"passThruCommonMetadata"`
	Flipping                  bool   `gorm:"type:boolean" json:"flipping"`
	Rotation                  int    `gorm:"type:int" json:"rotation"`
	Cropping                  string `gorm:"type:varchar(255)" json:"cropping"`
	Limit                     string `gorm:"type:varchar(255)" json:"limit"`
	Anamorphic                string `gorm:"type:varchar(255)" json:"anamorphic"`
	Fill                      string `gorm:"type:varchar(255)" json:"fill"`
	Color                     string `gorm:"type:varchar(255)" json:"color"`
	Detelecine                string `gorm:"type:varchar(255)" json:"detelecine"`
	InterlaceDetection        string `gorm:"type:varchar(255)" json:"interlaceDetection"`
	Deinterlace               string `gorm:"type:varchar(255)" json:"deinterlace"`
	DeinterlacePreset         string `gorm:"type:varchar(255)" json:"deinterlacePreset"`
	Deblock                   string `gorm:"type:varchar(255)" json:"deblock"`
	DeblockTune               string `gorm:"type:varchar(255)" json:"deblockTune"`
	Denoise                   string `gorm:"type:varchar(255)" json:"denoise"`
	DenoisePreset             string `gorm:"type:varchar(255)" json:"denoisePreset"`
	DenoiseTune               string `gorm:"type:varchar(255)" json:"denoiseTune"`
	ChromaSmooth              string `gorm:"type:varchar(255)" json:"chromaSmooth"`
	ChromaSmoothTune          string `gorm:"type:varchar(255)" json:"chromaSmoothTune"`
	Sharpen                   string `gorm:"type:varchar(255)" json:"sharpen"`
	SharpenPreset             string `gorm:"type:varchar(255)" json:"sharpenPreset"`
	SharpenTune               string `gorm:"type:varchar(255)" json:"sharpenTune"`
	Colorspace                string `gorm:"type:varchar(255)" json:"colorspace"`
	Grayscale                 bool   `gorm:"type:boolean" json:"grayscale"`
	Codec                     string `gorm:"type:varchar(255)" json:"codec"`
	Encoder                   string `gorm:"type:varchar(255)" json:"encoder"`
	Framerate                 string `gorm:"type:varchar(255)" json:"framerate"`
	FramerateType             string `gorm:"type:varchar(255)" json:"framerateType"`
	QualityType               string `gorm:"type:varchar(255)" json:"qualityType"`
	ConstantQuality           int    `gorm:"type:int" json:"constantQuality"`
	AverageBitrate            int    `gorm:"type:int" json:"averageBitrate"`
	MultipassEncoding         bool   `gorm:"type:boolean" json:"multipassEncoding"`
	Preset                    string `gorm:"type:varchar(255)" json:"preset"`
	Tune                      string `gorm:"type:varchar(255)" json:"tune"`
	Profile                   string `gorm:"type:varchar(255)" json:"profile"`
	Level                     string `gorm:"type:varchar(255)" json:"level"`
	FastDecode                bool   `gorm:"type:boolean" json:"fastDecode"`
	MapUntaggedAudioTracks    bool   `gorm:"type:boolean" json:"mapUntaggedAudioTracks"`
	MapUntaggedSubtitleTracks bool   `gorm:"type:boolean" json:"mapUntaggedSubtitleTracks"`

	ProfileAudioLanguages    []ProfileAudioLanguage    `gorm:"foreignkey:ProfileId" json:"profileAudioLanguages"`
	ProfileSubtitleLanguages []ProfileSubtitleLanguage `gorm:"foreignkey:ProfileId" json:"profileSubtitleLanguages"`
	ProfileCodecs            []ProfileCodec            `gorm:"foreignkey:ProfileID" json:"codecs"`
}
