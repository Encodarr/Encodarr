package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"transfigurr/internal/config"
	"transfigurr/internal/interfaces/repositories"
	"transfigurr/internal/models"
)

func ffmpegProbe(inputFile string) (models.ProbeData, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", inputFile)

	out, err := cmd.Output()
	if err != nil {
		return models.ProbeData{}, err
	}

	var probeData models.ProbeData
	err = json.Unmarshal(out, &probeData)
	if err != nil {
		return models.ProbeData{}, err
	}

	return probeData, nil
}

func AnalyzeMediaFile(filePath string) (string, error) {
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", err
	}

	// Check if the file is readable
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	file.Close()

	// Check if the file is writable
	file, err = os.OpenFile(filePath, os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}
	file.Close()

	// Run ffmpeg probe
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Extract the video codec
	codec := strings.TrimSpace(string(out))
	return codec, nil
}

func createFFMPEGFilter(profile models.Profile) []string {
	filters := []string{}

	if profile.Flipping {
		filters = append(filters, "hflip")
	}

	if profile.Rotation == 90 {
		filters = append(filters, "transpose=1")
	} else if profile.Rotation == -90 {
		filters = append(filters, "transpose=2")
	}

	if profile.Cropping == "conservative" {
		filters = append(filters, "cropdetect=24:16:0")
	} else if profile.Cropping == "automatic" {
		filters = append(filters, "cropdetect")
	}

	if profile.Anamorphic == "automatic" {
		filters = append(filters, "setsar=1")
	}

	if profile.Fill != "none" {
		if profile.Fill == "height" {
			filters = append(filters, fmt.Sprintf("pad=iw:ih+10:0:5:%s", profile.Color))
		} else if profile.Fill == "width" {
			filters = append(filters, fmt.Sprintf("pad=iw+10:ih:5:0:%s", profile.Color))
		} else if profile.Fill == "width & height" {
			filters = append(filters, fmt.Sprintf("pad=iw+10:ih+10:5:5:%s", profile.Color))
		}
	}

	if profile.Detelecine == "default" {
		filters = append(filters, "detelecine")
	}

	if profile.InterlaceDetection == "default" {
		filters = append(filters, "idet")
	} else if profile.InterlaceDetection == "less sensitive" {
		filters = append(filters, "idet=half_life=1:mult=2")
	}

	if profile.Deinterlace == "yadif" {
		if profile.DeinterlacePreset == "skip spatial check" {
			filters = append(filters, "yadif=mode=2")
		} else if profile.DeinterlacePreset == "bob" {
			filters = append(filters, "yadif=mode=1")
		} else {
			filters = append(filters, "yadif")
		}
	} else if profile.Deinterlace == "decomb" {
		if profile.DeinterlacePreset == "bob" {
			filters = append(filters, "w3fdif=complexity=2")
		} else if profile.DeinterlacePreset == "eedi2" {
			filters = append(filters, "w3fdif=complexity=1")
		} else if profile.DeinterlacePreset == "eedi2bob" {
			filters = append(filters, "w3fdif=complexity=0")
		} else {
			filters = append(filters, "w3fdif")
		}
	} else if profile.Deinterlace == "bwdif" {
		if profile.DeinterlacePreset == "bob" {
			filters = append(filters, "bwdif=mode=1")
		} else {
			filters = append(filters, "bwdif")
		}
	}
	if profile.Deblock != "off" {
		if profile.DeblockTune == "strong" {
			filters = append(filters, "deblock=3:3")
		} else if profile.DeblockTune == "weak" {
			filters = append(filters, "deblock=1:1")
		} else {
			filters = append(filters, "deblock")
		}
	}

	if profile.Denoise == "nlmeans" {
		if profile.DenoisePreset == "ultralight" {
			filters = append(filters, "nlmeans=s=1:d=1")
		} else if profile.DenoisePreset == "light" {
			filters = append(filters, "nlmeans=s=3:d=3")
		} else if profile.DenoisePreset == "medium" {
			filters = append(filters, "nlmeans=s=5:d=5")
		} else if profile.DenoisePreset == "strong" {
			filters = append(filters, "nlmeans=s=7:d=7")
		} else {
			filters = append(filters, "nlmeans")
		}

		if profile.DenoiseTune == "film" {
			filters[len(filters)-1] += ":p=film"
		} else if profile.DenoiseTune == "grain" {
			filters[len(filters)-1] += ":p=grain"
		} else if profile.DenoiseTune == "high motion" {
			filters[len(filters)-1] += ":p=highmotion"
		} else if profile.DenoiseTune == "animation" {
			filters[len(filters)-1] += ":p=animation"
		} else if profile.DenoiseTune == "tape" {
			filters[len(filters)-1] += ":p=tape"
		} else if profile.DenoiseTune == "sprite" {
			filters[len(filters)-1] += ":p=sprite"
		}
	} else if profile.Denoise == "hqdn3d" {
		if profile.DenoisePreset == "ultralight" {
			filters = append(filters, "hqdn3d=1:1:1:1")
		} else if profile.DenoisePreset == "light" {
			filters = append(filters, "hqdn3d=3:3:3:3")
		} else if profile.DenoisePreset == "medium" {
			filters = append(filters, "hqdn3d=5:5:5:5")
		} else if profile.DenoisePreset == "strong" {
			filters = append(filters, "hqdn3d=7:7:7:7")
		} else {
			filters = append(filters, "hqdn3d")
		}
	}

	if profile.ChromaSmooth == "ultralight" {
		filters = append(filters, "chroma_smooth=radius=1:strength=1")
	} else if profile.ChromaSmooth == "light" {
		filters = append(filters, "chroma_smooth=radius=2:strength=2")
	} else if profile.ChromaSmooth == "medium" {
		filters = append(filters, "chroma_smooth=radius=3:strength=3")
	} else if profile.ChromaSmooth == "strong" {
		filters = append(filters, "chroma_smooth=radius=4:strength=4")
	} else if profile.ChromaSmooth == "stronger" {
		filters = append(filters, "chroma_smooth=radius=5:strength=5")
	} else if profile.ChromaSmooth == "very strong" {
		filters = append(filters, "chroma_smooth=radius=6:strength=6")
	}

	if profile.ChromaSmoothTune == "tiny" {
		filters[len(filters)-1] += ":size=1"
	} else if profile.ChromaSmoothTune == "small" {
		filters[len(filters)-1] += ":size=2"
	} else if profile.ChromaSmoothTune == "medium" {
		filters[len(filters)-1] += ":size=3"
	} else if profile.ChromaSmoothTune == "wide" {
		filters[len(filters)-1] += ":size=4"
	} else if profile.ChromaSmoothTune == "very wide" {
		filters[len(filters)-1] += ":size=5"
	}

	if profile.Sharpen == "unsharp" {
		if profile.SharpenPreset == "ultralight" {
			filters = append(filters, "unsharp=3:3:0.3:3:3:0.3")
		} else if profile.SharpenPreset == "light" {
			filters = append(filters, "unsharp=5:5:0.5:5:5:0.5")
		} else if profile.SharpenPreset == "medium" {
			filters = append(filters, "unsharp=7:7:0.7:7:7:0.7")
		} else if profile.SharpenPreset == "strong" {
			filters = append(filters, "unsharp=9:9:0.9:9:9:0.9")
		} else if profile.SharpenPreset == "stronger" {
			filters = append(filters, "unsharp=11:11:1.1:11:11:1.1")
		} else if profile.SharpenPreset == "very strong" {
			filters = append(filters, "unsharp=13:13:1.3:13:13:1.3")
		} else {
			filters = append(filters, "unsharp")
		}
	} else if profile.Sharpen == "lapsharp" {
		if profile.SharpenPreset == "ultralight" {
			filters = append(filters, "lapsharp=c=0.3")
		} else if profile.SharpenPreset == "light" {
			filters = append(filters, "lapsharp=c=0.5")
		} else if profile.SharpenPreset == "medium" {
			filters = append(filters, "lapsharp=c=0.7")
		} else if profile.SharpenPreset == "strong" {
			filters = append(filters, "lapsharp=c=0.9")
		} else if profile.SharpenPreset == "stronger" {
			filters = append(filters, "lapsharp=c=1.1")
		} else if profile.SharpenPreset == "very strong" {
			filters = append(filters, "lapsharp=c=1.3")
		} else {
			filters = append(filters, "lapsharp")
		}
	}

	if profile.Colorspace == "bt.2020" {
		filters = append(filters, "colorspace=bt2020")
	} else if profile.Colorspace == "bt.709" {
		filters = append(filters, "colorspace=bt709")
	} else if profile.Colorspace == "bt.601" {
		filters = append(filters, "colorspace=bt601-6-525")
	} else if profile.Colorspace == "bt.601 smpte-c" {
		filters = append(filters, "colorspace=smpte170m")
	} else if profile.Colorspace == "bt.601 ebu" {
		filters = append(filters, "colorspace=bt601-6-625")
	}

	if profile.Grayscale {
		filters = append(filters, "format=gray")
	}

	if profile.Limit != "none" {
		filters = append(filters, fmt.Sprintf("scale=%s:{-1}", profile.Limit))
	}

	return filters
}

func getStreamIndices(inputFile string, streamType string, lang string) []int {
	probeOutput, _ := ffmpegProbe(inputFile)

	var indices []int
	for _, stream := range probeOutput.Streams {
		if stream.CodecType == streamType && (lang == "" || stream.Tags.Language == lang) {
			indices = append(indices, stream.Index)
		}
	}

	return indices
}

func parseTimeToSeconds(timeStr string) (float64, error) {
	layout := "15:04:05.000"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		layout = "15:04:05.00"
		t, err = time.Parse(layout, timeStr)
		if err != nil {
			return 0, fmt.Errorf("failed to parse time: %w", err)
		}
	}

	hours := t.Hour()
	minutes := t.Minute()
	seconds := t.Second()
	milliseconds := t.Nanosecond() / 1e6

	totalSeconds := float64(hours*3600+minutes*60+seconds) + float64(milliseconds)/1000
	return totalSeconds, nil
}

func createFFmpegCommand(inputFile string, outputFile string, profile models.Profile, hasAudio bool, hasSubtitle bool) []string {
	command := []string{"ffmpeg", "-y", "-i", inputFile}
	encoder := profile.Encoder
	if encoder != "" {
		command = append(command, "-vcodec", encoder)
		command = append(command, "-map", "0:v")
	}

	if hasAudio {
		audioLanguages := profile.ProfileAudioLanguages

		anyLanguage := false
		for _, lang := range audioLanguages {
			if lang.Language == "any" {
				anyLanguage = true
				break
			}
		}

		if anyLanguage {
			command = append(command, "-map", "0:a")
		} else {
			for _, lang := range audioLanguages {
				audioIndices := getStreamIndices(inputFile, "audio", lang.Language)
				for _, index := range audioIndices {
					command = append(command, "-map", fmt.Sprintf("0:%d", index))
				}
			}
			if profile.MapUntaggedAudioTracks {
				untaggedAudioIndices := getStreamIndices(inputFile, "audio", "")
				for _, index := range untaggedAudioIndices {
					command = append(command, "-map", fmt.Sprintf("0:%d", index))
				}
			}
		}
	}

	if hasSubtitle {
		subtitleLanguages := profile.ProfileSubtitleLanguages
		container := profile.Container
		allLanguage := false
		for _, lang := range subtitleLanguages {
			if lang.Language == "all" {
				allLanguage = true
				break
			}
		}

		if allLanguage {
			if container == "matroska" {
				command = append(command, "-scodec", "srt", "-map", "0:s")
			} else {
				command = append(command, "-scodec", "mov_text", "-map", "0:s")
			}
		} else {
			for _, lang := range subtitleLanguages {
				subtitleIndices := getStreamIndices(inputFile, "subtitle", lang.Language)
				for _, index := range subtitleIndices {
					if container == "matroska" {
						command = append(command, "-scodec", "srt", "-map", fmt.Sprintf("0:%d", index))
					} else {
						command = append(command, "-scodec", "mov_text", "-map", fmt.Sprintf("0:%d", index))
					}
				}
			}
			if profile.MapUntaggedSubtitleTracks {
				untaggedSubtitleIndices := getStreamIndices(inputFile, "subtitle", "")
				for _, index := range untaggedSubtitleIndices {
					if container == "matroska" {
						command = append(command, "-scodec", "srt", "-map", fmt.Sprintf("0:%d", index))
					} else {
						command = append(command, "-scodec", "mov_text", "-map", fmt.Sprintf("0:%d", index))
					}
				}
			}
		}
	}

	preset := profile.Preset
	codec := profile.Codec
	if preset != "" {
		if codec == "av1" {
			command = append(command, "-cpu-used", preset)
		} else {
			command = append(command, "-preset", preset)
		}
	}

	if profile.PassThruCommonMetadata {
		command = append(command, "-map_metadata", "0")
	}

	filters := createFFMPEGFilter(profile)
	if len(filters) > 0 {
		command = append(command, "-vf", strings.Join(filters, ","))
	}

	framerate := profile.Framerate
	framerateType := profile.FramerateType
	if framerate != "" && framerate != "same as source" {
		command = append(command, "-r", framerate)
		if framerateType == "peak framerate" {
			command = append(command, "-vsync", "2")
		} else {
			command = append(command, "-vsync", "1")
		}
	}

	qualityType := profile.QualityType
	constantQuality := profile.ConstantQuality
	averageBitrate := profile.AverageBitrate
	if qualityType != "" {
		if qualityType == "constant quality" {
			if codec == "mpeg4" {
				command = append(command, "-q:v", strconv.FormatFloat(float64(constantQuality), 'f', -1, 64))
			} else {
				command = append(command, "-crf", strconv.FormatFloat(float64(constantQuality), 'f', -1, 64))
			}
		} else if qualityType == "average bitrate" {
			command = append(command, "-b:v", strconv.FormatFloat(float64(averageBitrate), 'f', -1, 64))
		}
	}

	tune := profile.Tune
	fastDecode := profile.FastDecode
	if tune != "" && tune != "none" {
		if fastDecode && (codec == "h264" || codec == "av1") {
			command = append(command, "-tune", tune+",fastdecode")
		} else {
			command = append(command, "-tune", tune)
		}
	} else if fastDecode && (codec == "h264" || codec == "av1") {
		command = append(command, "-tune", "fastdecode")
	}

	if profile.Profile != "" && profile.Profile != "auto" {
		command = append(command, "-profile:v", profile.Profile)
	}

	level := profile.Level
	if level != "" && level != "auto" {
		command = append(command, "-level:v", level)
	}

	command = append(command, "-f", profile.Container, outputFile)
	return command
}

func moveOutputFile(inputFile, outputFile string) error {
	// Get the directory of the input file
	inputFileDir := filepath.Dir(inputFile)
	newOutputFilePath := filepath.Join(inputFileDir, filepath.Base(outputFile))

	// Create source file handle
	source, err := os.Open(outputFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer source.Close()

	// Create destination file handle
	destination, err := os.Create(newOutputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destination.Close()

	// Copy the contents
	if _, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	// Close files before removing
	source.Close()
	destination.Close()

	// Remove the original output file
	if err := os.Remove(outputFile); err != nil {
		return fmt.Errorf("failed to remove source file: %v", err)
	}

	return nil
}

func runFFmpeg(inputFile string, outputFile string, profile models.Profile, stage *string, progress *float64, eta *int, queueStatus *string, current *models.Item) bool {
	*stage = "Encoding"
	probe, err := ffmpegProbe(inputFile)
	if err != nil {
		return false
	}

	hasAudio, hasSubtitle := hasAudioAndSubtitleStreams(probe)
	totalDuration, _ := strconv.ParseFloat(probe.Format.Duration, 64)

	command := createFFmpegCommand(inputFile, outputFile, profile, hasAudio, hasSubtitle)
	cmd := exec.Command(command[0], command[1:]...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return false
	}
	startTime := time.Now()
	if err := cmd.Start(); err != nil {
		return false
	}

	buf := make([]byte, 1024)
	for {
		n, err := stderr.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return false
		}
		if n == 0 {
			break
		}
		output := string(buf[:n])
		re := regexp.MustCompile(`time=(\d+:\d+:\d+.\d+)`)
		match := re.FindStringSubmatch(output)
		if len(match) > 0 {
			currentTime := match[1]
			currentSeconds, err := parseTimeToSeconds(currentTime)
			if err != nil {
				return false
			}
			*progress = (currentSeconds / totalDuration) * 100

			elapsedTime := time.Since(startTime).Seconds()
			estimatedTotalTime := elapsedTime / (*progress / 100)
			*eta = int(estimatedTotalTime - elapsedTime)
		}
	}
	*stage = "Moving"
	moveOutputFile(inputFile, outputFile)
	*stage = "Idle"
	*eta = 0
	*progress = 0
	return true
}

func hasAudioAndSubtitleStreams(probe models.ProbeData) (bool, bool) {
	hasAudio := false
	hasSubtitle := false

	for _, stream := range probe.Streams {
		if stream.CodecType == "audio" {
			hasAudio = true
		} else if stream.CodecType == "subtitle" {
			hasSubtitle = true
		}

		if hasAudio && hasSubtitle {
			break
		}
	}

	return hasAudio, hasSubtitle
}

func EncodeMovie(item models.Item, movieRepo repositories.MovieRepositoryInterface, historyRepo repositories.HistoryRepositoryInterface, settingRepo repositories.SettingRepositoryInterface, profileRepo repositories.ProfileRepositoryInterface, stage *string, progress *float64, eta *int, queueStatus *string, current *models.Item) {
	movie, _ := movieRepo.GetMovieById(item.Id)
	profile, _ := profileRepo.GetProfileById(movie.ProfileID)
	prevSize := movie.Size

	extension := profile.Extension
	filename := strings.TrimSuffix(filepath.Base(movie.Filename), filepath.Ext(movie.Filename))
	inputFile := movie.Path
	outputFilename := filepath.Join(filename + "." + extension)
	outputFile := filepath.Join(config.TranscodeFolder, outputFilename)
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return
	}

	videoStream, _ := AnalyzeMediaFile(inputFile)

	if videoStream == profile.Codec {
		return
	}

	*current = models.Item{Id: movie.Id, Type: "movie", Name: movie.Name, ProfileId: movie.ProfileID}
	encodingSuccessful := runFFmpeg(inputFile, outputFile, profile, stage, progress, eta, queueStatus, current)
	*current = models.Item{}

	if !encodingSuccessful {
		return
	}
	ffmpegData, _ := ffmpegProbe(filepath.Join(filepath.Dir(movie.Path), outputFilename))

	size, _ := strconv.Atoi(ffmpegData.Format.Size)

	movie.Size = int(size)
	movie.SpaceSaved = movie.OriginalSize - movie.Size
	movie.VideoCodec = ffmpegData.Streams[0].CodecName
	movie.Missing = false
	movie.Filename = outputFilename
	movie.Path = filepath.Join(filepath.Dir(movie.Path), outputFilename)
	movieRepo.UpsertMovie(movie.Id, movie)
	history := &models.History{
		MediaId:   movie.Id,
		Name:      movie.Name,
		Type:      "movie",
		ProfileId: movie.ProfileID,
		PrevCodec: videoStream,
		NewCodec:  profile.Codec,
		PrevSize:  prevSize,
		NewSize:   movie.Size,
		Date:      time.Now().Format(time.RFC3339),
	}
	historyRepo.UpsertHistoryById(history)
}

func EncodeEpisode(item models.Item, seriesRepo repositories.SeriesRepositoryInterface, historyRepo repositories.HistoryRepositoryInterface, episodeRepo repositories.EpisodeRepositoryInterface, settingRepo repositories.SettingRepositoryInterface, profileRepo repositories.ProfileRepositoryInterface, stage *string, progress *float64, eta *int, queueStatus *string, current *models.Item) {

	episode, _ := episodeRepo.GetEpisodeById(item.Id)

	prevSize := episode.Size

	series, _ := seriesRepo.GetSeriesByID(episode.SeriesId)

	profile, _ := profileRepo.GetProfileById(series.ProfileID)

	extension := profile.Extension
	filename := strings.TrimSuffix(filepath.Base(episode.Filename), filepath.Ext(episode.Filename))
	inputFile := episode.Path
	outputFileName := filepath.Join(filename + "." + extension)
	outputFile := filepath.Join(config.TranscodeFolder, outputFileName)

	videoStream, _ := AnalyzeMediaFile(inputFile)

	if videoStream == profile.Codec {
		return
	}

	*current = models.Item{Id: series.Id, Type: "series", Name: series.Name, ProfileId: series.ProfileID, SeriesId: episode.SeriesId, SeasonNumber: episode.SeasonNumber, EpisodeNumber: episode.EpisodeNumber}
	encodingSuccessful := runFFmpeg(inputFile, outputFile, profile, stage, progress, eta, queueStatus, current)
	*current = models.Item{}

	if !encodingSuccessful {
		return
	}
	ffmpegData, _ := ffmpegProbe(filepath.Join(filepath.Dir(episode.Path), outputFileName))

	size, _ := strconv.Atoi(ffmpegData.Format.Size)

	episode.Size = size
	episode.SpaceSaved = episode.OriginalSize - episode.Size
	episode.VideoCodec = ffmpegData.Streams[0].CodecName
	episode.Missing = false
	episode.Filename = outputFileName
	episode.Path = filepath.Join(filepath.Dir(episode.Path), outputFileName)
	episodeRepo.UpsertEpisode(episode.SeriesId, episode.SeasonNumber, episode.EpisodeNumber, episode)
	history := &models.History{
		MediaId:       episode.Id,
		Name:          series.Name,
		Type:          "episode",
		SeasonNumber:  episode.SeasonNumber,
		EpisodeNumber: episode.EpisodeNumber,
		ProfileId:     series.ProfileID,
		PrevCodec:     videoStream,
		NewCodec:      profile.Codec,
		PrevSize:      prevSize,
		NewSize:       episode.Size,
		Date:          time.Now().Format(time.RFC3339),
	}
	historyRepo.UpsertHistoryById(history)
}
