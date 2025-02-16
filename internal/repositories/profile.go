package repositories

import (
	"database/sql"
	"transfigurr/internal/models"
)

type ProfileRepository struct {
	DB *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{
		DB: db,
	}
}

func (repo *ProfileRepository) GetAllProfiles() ([]models.Profile, error) {
	rows, err := repo.DB.Query(`
        SELECT id, name, container, extension, pass_thru_common_metadata,
        flipping, rotation, cropping, limit_value, anamorphic, fill,
        color, detelecine, interlace_detection, deinterlace,
        deinterlace_preset, deblock, deblock_tune, denoise,
        denoise_preset, denoise_tune, chroma_smooth,
        chroma_smooth_tune, sharpen, sharpen_preset, sharpen_tune,
        colorspace, grayscale, codec, encoder, framerate,
        framerate_type, quality_type, constant_quality,
        average_bitrate, multipass_encoding, preset, tune,
        profile, level, fast_decode, map_untagged_audio_tracks,
        map_untagged_subtitle_tracks
        FROM profiles
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var p models.Profile
		err := rows.Scan(
			&p.Id, &p.Name, &p.Container, &p.Extension,
			&p.PassThruCommonMetadata, &p.Flipping, &p.Rotation,
			&p.Cropping, &p.Limit, &p.Anamorphic, &p.Fill,
			&p.Color, &p.Detelecine, &p.InterlaceDetection,
			&p.Deinterlace, &p.DeinterlacePreset, &p.Deblock,
			&p.DeblockTune, &p.Denoise, &p.DenoisePreset,
			&p.DenoiseTune, &p.ChromaSmooth, &p.ChromaSmoothTune,
			&p.Sharpen, &p.SharpenPreset, &p.SharpenTune,
			&p.Colorspace, &p.Grayscale, &p.Codec, &p.Encoder,
			&p.Framerate, &p.FramerateType, &p.QualityType,
			&p.ConstantQuality, &p.AverageBitrate,
			&p.MultipassEncoding, &p.Preset, &p.Tune,
			&p.Profile, &p.Level, &p.FastDecode,
			&p.MapUntaggedAudioTracks, &p.MapUntaggedSubtitleTracks,
		)
		if err != nil {
			return nil, err
		}

		// Load audio languages
		audioLangs, err := repo.getProfileAudioLanguages(p.Id)
		if err != nil {
			return nil, err
		}
		p.ProfileAudioLanguages = audioLangs

		// Load subtitle languages
		subtitleLangs, err := repo.getProfileSubtitleLanguages(p.Id)
		if err != nil {
			return nil, err
		}
		p.ProfileSubtitleLanguages = subtitleLangs

		// Load codecs
		codecs, err := repo.getProfileCodecs(p.Id)
		if err != nil {
			return nil, err
		}
		p.ProfileCodecs = codecs

		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (repo *ProfileRepository) UpsertProfile(profileId int, inputProfile models.Profile) (models.Profile, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		return models.Profile{}, err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM profiles WHERE id = ?)", profileId).Scan(&exists)
	if err != nil {
		return models.Profile{}, err
	}

	if exists {
		_, err = tx.Exec(`
            UPDATE profiles SET
            name = ?, container = ?, extension = ?,
            pass_thru_common_metadata = ?, flipping = ?,
            rotation = ?, cropping = ?, limit_value = ?,
            anamorphic = ?, fill = ?, color = ?,
            detelecine = ?, interlace_detection = ?,
            deinterlace = ?, deinterlace_preset = ?,
            deblock = ?, deblock_tune = ?, denoise = ?,
            denoise_preset = ?, denoise_tune = ?,
            chroma_smooth = ?, chroma_smooth_tune = ?,
            sharpen = ?, sharpen_preset = ?, sharpen_tune = ?,
            colorspace = ?, grayscale = ?, codec = ?,
            encoder = ?, framerate = ?, framerate_type = ?,
            quality_type = ?, constant_quality = ?,
            average_bitrate = ?, multipass_encoding = ?,
            preset = ?, tune = ?, profile = ?, level = ?,
            fast_decode = ?, map_untagged_audio_tracks = ?,
            map_untagged_subtitle_tracks = ?
            WHERE id = ?`,
			inputProfile.Name, inputProfile.Container,
			inputProfile.Extension, inputProfile.PassThruCommonMetadata,
			inputProfile.Flipping, inputProfile.Rotation,
			inputProfile.Cropping, inputProfile.Limit,
			inputProfile.Anamorphic, inputProfile.Fill,
			inputProfile.Color, inputProfile.Detelecine,
			inputProfile.InterlaceDetection, inputProfile.Deinterlace,
			inputProfile.DeinterlacePreset, inputProfile.Deblock,
			inputProfile.DeblockTune, inputProfile.Denoise,
			inputProfile.DenoisePreset, inputProfile.DenoiseTune,
			inputProfile.ChromaSmooth, inputProfile.ChromaSmoothTune,
			inputProfile.Sharpen, inputProfile.SharpenPreset,
			inputProfile.SharpenTune, inputProfile.Colorspace,
			inputProfile.Grayscale, inputProfile.Codec,
			inputProfile.Encoder, inputProfile.Framerate,
			inputProfile.FramerateType, inputProfile.QualityType,
			inputProfile.ConstantQuality, inputProfile.AverageBitrate,
			inputProfile.MultipassEncoding, inputProfile.Preset,
			inputProfile.Tune, inputProfile.Profile,
			inputProfile.Level, inputProfile.FastDecode,
			inputProfile.MapUntaggedAudioTracks,
			inputProfile.MapUntaggedSubtitleTracks,
			profileId,
		)
	} else {
		_, err = tx.Exec(`
            INSERT INTO profiles (
                id, name, container, extension,
                pass_thru_common_metadata, flipping,
                rotation, cropping, limit_value,
                anamorphic, fill, color, detelecine,
                interlace_detection, deinterlace,
                deinterlace_preset, deblock,
                deblock_tune, denoise, denoise_preset,
                denoise_tune, chroma_smooth,
                chroma_smooth_tune, sharpen,
                sharpen_preset, sharpen_tune,
                colorspace, grayscale, codec,
                encoder, framerate, framerate_type,
                quality_type, constant_quality,
                average_bitrate, multipass_encoding,
                preset, tune, profile, level,
                fast_decode, map_untagged_audio_tracks,
                map_untagged_subtitle_tracks
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			profileId, inputProfile.Name, inputProfile.Container,
			inputProfile.Extension, inputProfile.PassThruCommonMetadata,
			inputProfile.Flipping, inputProfile.Rotation,
			inputProfile.Cropping, inputProfile.Limit,
			inputProfile.Anamorphic, inputProfile.Fill,
			inputProfile.Color, inputProfile.Detelecine,
			inputProfile.InterlaceDetection, inputProfile.Deinterlace,
			inputProfile.DeinterlacePreset, inputProfile.Deblock,
			inputProfile.DeblockTune, inputProfile.Denoise,
			inputProfile.DenoisePreset, inputProfile.DenoiseTune,
			inputProfile.ChromaSmooth, inputProfile.ChromaSmoothTune,
			inputProfile.Sharpen, inputProfile.SharpenPreset,
			inputProfile.SharpenTune, inputProfile.Colorspace,
			inputProfile.Grayscale, inputProfile.Codec,
			inputProfile.Encoder, inputProfile.Framerate,
			inputProfile.FramerateType, inputProfile.QualityType,
			inputProfile.ConstantQuality, inputProfile.AverageBitrate,
			inputProfile.MultipassEncoding, inputProfile.Preset,
			inputProfile.Tune, inputProfile.Profile,
			inputProfile.Level, inputProfile.FastDecode,
			inputProfile.MapUntaggedAudioTracks,
			inputProfile.MapUntaggedSubtitleTracks,
		)
	}
	if err != nil {
		return models.Profile{}, err
	}

	// Update audio languages
	_, err = tx.Exec("DELETE FROM profile_audio_languages WHERE profile_id = ?", profileId)
	if err != nil {
		return models.Profile{}, err
	}

	for _, lang := range inputProfile.ProfileAudioLanguages {
		_, err = tx.Exec(
			"INSERT INTO profile_audio_languages (profile_id, language) VALUES (?, ?)",
			profileId, lang.Language,
		)
		if err != nil {
			return models.Profile{}, err
		}
	}

	// Update subtitle languages
	_, err = tx.Exec("DELETE FROM profile_subtitle_languages WHERE profile_id = ?", profileId)
	if err != nil {
		return models.Profile{}, err
	}

	for _, lang := range inputProfile.ProfileSubtitleLanguages {
		_, err = tx.Exec(
			"INSERT INTO profile_subtitle_languages (profile_id, language) VALUES (?, ?)",
			profileId, lang.Language,
		)
		if err != nil {
			return models.Profile{}, err
		}
	}

	// Update codecs
	_, err = tx.Exec("DELETE FROM profile_codecs WHERE profile_id = ?", profileId)
	if err != nil {
		return models.Profile{}, err
	}

	for _, codec := range inputProfile.ProfileCodecs {
		_, err = tx.Exec(
			"INSERT INTO profile_codecs (profile_id, codec_id) VALUES (?, ?)",
			profileId, codec.CodecId,
		)
		if err != nil {
			return models.Profile{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return models.Profile{}, err
	}

	return repo.GetProfileById(profileId)
}

func (repo *ProfileRepository) GetProfileById(profileId int) (models.Profile, error) {
	var p models.Profile
	err := repo.DB.QueryRow(`
		SELECT id, name, container, extension, pass_thru_common_metadata,
		flipping, rotation, cropping, limit_value, anamorphic, fill,
		color, detelecine, interlace_detection, deinterlace,
		deinterlace_preset, deblock, deblock_tune, denoise,
		denoise_preset, denoise_tune, chroma_smooth,
		chroma_smooth_tune, sharpen, sharpen_preset, sharpen_tune,
		colorspace, grayscale, codec, encoder, framerate,
		framerate_type, quality_type, constant_quality,
		average_bitrate, multipass_encoding, preset, tune,
		profile, level, fast_decode, map_untagged_audio_tracks,
		map_untagged_subtitle_tracks
		FROM profiles WHERE id = ?`, profileId).Scan(
		&p.Id, &p.Name, &p.Container, &p.Extension,
		&p.PassThruCommonMetadata, &p.Flipping, &p.Rotation,
		&p.Cropping, &p.Limit, &p.Anamorphic, &p.Fill,
		&p.Color, &p.Detelecine, &p.InterlaceDetection,
		&p.Deinterlace, &p.DeinterlacePreset, &p.Deblock,
		&p.DeblockTune, &p.Denoise, &p.DenoisePreset,
		&p.DenoiseTune, &p.ChromaSmooth, &p.ChromaSmoothTune,
		&p.Sharpen, &p.SharpenPreset, &p.SharpenTune,
		&p.Colorspace, &p.Grayscale, &p.Codec, &p.Encoder,
		&p.Framerate, &p.FramerateType, &p.QualityType,
		&p.ConstantQuality, &p.AverageBitrate,
		&p.MultipassEncoding, &p.Preset, &p.Tune,
		&p.Profile, &p.Level, &p.FastDecode,
		&p.MapUntaggedAudioTracks, &p.MapUntaggedSubtitleTracks,
	)
	if err != nil {
		return models.Profile{}, err
	}

	// Load audio languages
	audioLangs, err := repo.getProfileAudioLanguages(p.Id)
	if err != nil {
		return models.Profile{}, err
	}
	p.ProfileAudioLanguages = audioLangs

	// Load subtitle languages
	subtitleLangs, err := repo.getProfileSubtitleLanguages(p.Id)
	if err != nil {
		return models.Profile{}, err
	}
	p.ProfileSubtitleLanguages = subtitleLangs

	// Load codecs
	codecs, err := repo.getProfileCodecs(p.Id)
	if err != nil {
		return models.Profile{}, err
	}
	p.ProfileCodecs = codecs

	return p, nil
}

func (repo *ProfileRepository) DeleteProfileById(profileId int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete associated data
	_, err = tx.Exec("DELETE FROM profile_codecs WHERE profile_id = ?", profileId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM profile_audio_languages WHERE profile_id = ?", profileId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM profile_subtitle_languages WHERE profile_id = ?", profileId)
	if err != nil {
		return err
	}

	// Delete the profile
	_, err = tx.Exec("DELETE FROM profiles WHERE id = ?", profileId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *ProfileRepository) getProfileAudioLanguages(profileId int) ([]models.ProfileAudioLanguage, error) {
	rows, err := repo.DB.Query(
		"SELECT id, profile_id, language FROM profile_audio_languages WHERE profile_id = ?",
		profileId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []models.ProfileAudioLanguage
	for rows.Next() {
		var lang models.ProfileAudioLanguage
		err := rows.Scan(&lang.Id, &lang.ProfileId, &lang.Language)
		if err != nil {
			return nil, err
		}
		languages = append(languages, lang)
	}
	return languages, nil
}

func (repo *ProfileRepository) getProfileSubtitleLanguages(profileId int) ([]models.ProfileSubtitleLanguage, error) {
	rows, err := repo.DB.Query(
		"SELECT id, profile_id, language FROM profile_subtitle_languages WHERE profile_id = ?",
		profileId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []models.ProfileSubtitleLanguage
	for rows.Next() {
		var lang models.ProfileSubtitleLanguage
		err := rows.Scan(&lang.Id, &lang.ProfileId, &lang.Language)
		if err != nil {
			return nil, err
		}
		languages = append(languages, lang)
	}
	return languages, nil
}

func (repo *ProfileRepository) getProfileCodecs(profileId int) ([]models.ProfileCodec, error) {
	rows, err := repo.DB.Query(
		"SELECT id, profile_id, codec_id FROM profile_codecs WHERE profile_id = ?",
		profileId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codecs []models.ProfileCodec
	for rows.Next() {
		var codec models.ProfileCodec
		err := rows.Scan(&codec.Id, &codec.ProfileId, &codec.CodecId)
		if err != nil {
			return nil, err
		}
		codecs = append(codecs, codec)
	}
	return codecs, nil
}
