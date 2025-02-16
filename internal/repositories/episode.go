package repositories

import (
	"database/sql"
	"strconv"
	"transfigurr/internal/models"
)

type EpisodeRepository struct {
	DB *sql.DB
}

func NewEpisodeRepository(db *sql.DB) *EpisodeRepository {
	return &EpisodeRepository{
		DB: db,
	}
}

func (repo *EpisodeRepository) GetEpisodes(seriesId string, seasonNumber int) ([]models.Episode, error) {
	rows, err := repo.DB.Query(`
        SELECT id, series_id, season_id, episode_number, season_name, 
        season_number, filename, episode_name, video_codec, air_date, 
        size, space_saved, original_size, path, missing
        FROM episodes 
        WHERE series_id = ? AND season_number = ?`,
		seriesId, seasonNumber,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []models.Episode
	for rows.Next() {
		var episode models.Episode
		err := rows.Scan(
			&episode.Id, &episode.SeriesId, &episode.SeasonId,
			&episode.EpisodeNumber, &episode.SeasonName,
			&episode.SeasonNumber, &episode.Filename,
			&episode.EpisodeName, &episode.VideoCodec,
			&episode.AirDate, &episode.Size, &episode.SpaceSaved,
			&episode.OriginalSize, &episode.Path, &episode.Missing,
		)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, episode)
	}
	return episodes, nil
}

func (repo *EpisodeRepository) UpsertEpisode(seriesId string, seasonNumber int, episodeNumber int, inputEpisode models.Episode) (models.Episode, error) {
	inputEpisode.Id = seriesId + strconv.Itoa(seasonNumber) + strconv.Itoa(episodeNumber)
	inputEpisode.SeriesId = seriesId

	var exists bool
	err := repo.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM episodes WHERE id = ?)",
		inputEpisode.Id,
	).Scan(&exists)
	if err != nil {
		return models.Episode{}, err
	}

	if exists {
		_, err = repo.DB.Exec(`
            UPDATE episodes SET
            season_id = ?, episode_number = ?, season_name = ?,
            season_number = ?, filename = ?, episode_name = ?,
            video_codec = ?, air_date = ?, size = ?,
            space_saved = ?, original_size = ?, path = ?,
            missing = ?
            WHERE id = ?`,
			inputEpisode.SeasonId, inputEpisode.EpisodeNumber,
			inputEpisode.SeasonName, inputEpisode.SeasonNumber,
			inputEpisode.Filename, inputEpisode.EpisodeName,
			inputEpisode.VideoCodec, inputEpisode.AirDate,
			inputEpisode.Size, inputEpisode.SpaceSaved,
			inputEpisode.OriginalSize, inputEpisode.Path,
			inputEpisode.Missing, inputEpisode.Id,
		)
	} else {
		_, err = repo.DB.Exec(`
            INSERT INTO episodes (
                id, series_id, season_id, episode_number,
                season_name, season_number, filename,
                episode_name, video_codec, air_date, size,
                space_saved, original_size, path, missing
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			inputEpisode.Id, inputEpisode.SeriesId,
			inputEpisode.SeasonId, inputEpisode.EpisodeNumber,
			inputEpisode.SeasonName, inputEpisode.SeasonNumber,
			inputEpisode.Filename, inputEpisode.EpisodeName,
			inputEpisode.VideoCodec, inputEpisode.AirDate,
			inputEpisode.Size, inputEpisode.SpaceSaved,
			inputEpisode.OriginalSize, inputEpisode.Path,
			inputEpisode.Missing,
		)
	}
	if err != nil {
		return models.Episode{}, err
	}

	return repo.GetEpisodeById(inputEpisode.Id)
}

func (repo *EpisodeRepository) GetEpisodeById(episodeId string) (models.Episode, error) {
	var episode models.Episode
	err := repo.DB.QueryRow(`
        SELECT id, series_id, season_id, episode_number,
        season_name, season_number, filename, episode_name,
        video_codec, air_date, size, space_saved,
        original_size, path, missing
        FROM episodes WHERE id = ?`,
		episodeId,
	).Scan(
		&episode.Id, &episode.SeriesId, &episode.SeasonId,
		&episode.EpisodeNumber, &episode.SeasonName,
		&episode.SeasonNumber, &episode.Filename,
		&episode.EpisodeName, &episode.VideoCodec,
		&episode.AirDate, &episode.Size, &episode.SpaceSaved,
		&episode.OriginalSize, &episode.Path, &episode.Missing,
	)
	if err != nil {
		return models.Episode{}, err
	}
	return episode, nil
}

func (repo *EpisodeRepository) GetEpisodeBySeriesSeasonEpisode(seriesId string, seasonNumber int, episodeNumber int) (models.Episode, error) {
	var episode models.Episode
	err := repo.DB.QueryRow(`
        SELECT id, series_id, season_id, episode_number,
        season_name, season_number, filename, episode_name,
        video_codec, air_date, size, space_saved,
        original_size, path, missing
        FROM episodes 
        WHERE series_id = ? AND season_number = ? AND episode_number = ?`,
		seriesId, seasonNumber, episodeNumber,
	).Scan(
		&episode.Id, &episode.SeriesId, &episode.SeasonId,
		&episode.EpisodeNumber, &episode.SeasonName,
		&episode.SeasonNumber, &episode.Filename,
		&episode.EpisodeName, &episode.VideoCodec,
		&episode.AirDate, &episode.Size, &episode.SpaceSaved,
		&episode.OriginalSize, &episode.Path, &episode.Missing,
	)
	if err != nil {
		return models.Episode{}, err
	}
	return episode, nil
}

func (repo *EpisodeRepository) DeleteEpisodeById(seriesId string, seasonNumber int, episodeNumber int) error {
	_, err := repo.DB.Exec(`
        DELETE FROM episodes 
        WHERE series_id = ? AND season_number = ? AND episode_number = ?`,
		seriesId, seasonNumber, episodeNumber,
	)
	return err
}
