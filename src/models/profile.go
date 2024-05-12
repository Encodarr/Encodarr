package models

type Profile struct {
	BaseModel
	Id                        int    `gorm:"primary_key"`
	Name                      string `gorm:"type:varchar(255)"`
	Container                 string `gorm:"type:varchar(255)"`
	Extension                 string `gorm:"type:varchar(255)"`
	PassThruCommonMetadata    bool   `gorm:"type:boolean"`
	Flipping                  bool   `gorm:"type:boolean"`
	Rotation                  int    `gorm:"type:int"`
	Cropping                  string `gorm:"type:varchar(255)"`
	Limit                     string `gorm:"type:varchar(255)"`
	Anamorphic                string `gorm:"type:varchar(255)"`
	Fill                      string `gorm:"type:varchar(255)"`
	Color                     string `gorm:"type:varchar(255)"`
	Detelecine                string `gorm:"type:varchar(255)"`
	InterlaceDetection        string `gorm:"type:varchar(255)"`
	Deinterlace               string `gorm:"type:varchar(255)"`
	DeinterlacePreset         string `gorm:"type:varchar(255)"`
	Deblock                   string `gorm:"type:varchar(255)"`
	DeblockTune               string `gorm:"type:varchar(255)"`
	Denoise                   string `gorm:"type:varchar(255)"`
	DenoisePreset             string `gorm:"type:varchar(255)"`
	DenoiseTune               string `gorm:"type:varchar(255)"`
	ChromaSmooth              string `gorm:"type:varchar(255)"`
	ChromaSmoothTune          string `gorm:"type:varchar(255)"`
	Sharpen                   string `gorm:"type:varchar(255)"`
	SharpenPreset             string `gorm:"type:varchar(255)"`
	SharpenTune               string `gorm:"type:varchar(255)"`
	Colorspace                string `gorm:"type:varchar(255)"`
	Grayscale                 bool   `gorm:"type:boolean"`
	Codec                     string `gorm:"type:varchar(255)"`
	Encoder                   string `gorm:"type:varchar(255)"`
	Framerate                 string `gorm:"type:varchar(255)"`
	FramerateType             string `gorm:"type:varchar(255)"`
	QualityType               string `gorm:"type:varchar(255)"`
	ConstantQuality           int    `gorm:"type:int"`
	AverageBitrate            int    `gorm:"type:int"`
	MultipassEncoding         bool   `gorm:"type:boolean"`
	Preset                    string `gorm:"type:varchar(255)"`
	Tune                      string `gorm:"type:varchar(255)"`
	Profile                   string `gorm:"type:varchar(255)"`
	Level                     string `gorm:"type:varchar(255)"`
	FastDecode                bool   `gorm:"type:boolean"`
	MapUntaggedAudioTracks    bool   `gorm:"type:boolean"`
	MapUntaggedSubtitleTracks bool   `gorm:"type:boolean"`
}
