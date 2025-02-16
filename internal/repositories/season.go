package repositories

import (
	"database/sql"
	"strconv"
	"transfigurr/internal/models"
)

type SeasonRepository struct {
	DB *sql.DB
}

func NewSeasonRepository(db *sql.DB) *SeasonRepository {
	return &SeasonRepository{
		DB: db,
	}
}

func (repo *SeasonRepository) GetSeasons(seriesId string) ([]models.Season, error) {
	rows, err := repo.DB.Query(`
        SELECT id, name, season_number, episode_count, size,
        series_id, space_saved, missing_episodes
        FROM seasons WHERE series_id = ?`, seriesId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seasons []models.Season
	for rows.Next() {
		var season models.Season
		err := rows.Scan(
			&season.Id, &season.Name, &season.SeasonNumber,
			&season.EpisodeCount, &season.Size, &season.SeriesId,
			&season.SpaceSaved, &season.MissingEpisodes,
		)
		if err != nil {
			return nil, err
		}

		// Load episodes
		episodes, err := repo.getEpisodesBySeason(season.Id)
		if err != nil {
			return nil, err
		}
		season.Episodes = episodes
		seasons = append(seasons, season)
	}
	return seasons, nil
}

func (repo *SeasonRepository) UpsertSeason(seriesId string, seasonNumber int, inputSeason models.Season) (models.Season, error) {
	inputSeason.Id = seriesId + strconv.Itoa(seasonNumber)
	inputSeason.SeriesId = seriesId

	// Check if season exists
	var exists bool
	err := repo.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM seasons WHERE series_id = ? AND season_number = ?)",
		seriesId, seasonNumber,
	).Scan(&exists)
	if err != nil {
		return models.Season{}, err
	}

	if exists {
		_, err = repo.DB.Exec(`
            UPDATE seasons SET
            name = ?, episode_count = ?, size = ?,
            space_saved = ?, missing_episodes = ?
            WHERE series_id = ? AND season_number = ?`,
			inputSeason.Name, inputSeason.EpisodeCount,
			inputSeason.Size, inputSeason.SpaceSaved,
			inputSeason.MissingEpisodes, seriesId, seasonNumber,
		)
	} else {
		_, err = repo.DB.Exec(`
            INSERT INTO seasons (
                id, name, season_number, episode_count, size,
                series_id, space_saved, missing_episodes
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			inputSeason.Id, inputSeason.Name, seasonNumber,
			inputSeason.EpisodeCount, inputSeason.Size,
			seriesId, inputSeason.SpaceSaved,
			inputSeason.MissingEpisodes,
		)
	}
	if err != nil {
		return models.Season{}, err
	}

	return repo.GetSeasonById(seriesId, seasonNumber)
}

func (repo *SeasonRepository) GetSeasonById(seriesId string, seasonNumber int) (models.Season, error) {
	var season models.Season
	err := repo.DB.QueryRow(`
        SELECT id, name, season_number, episode_count, size,
        series_id, space_saved, missing_episodes
        FROM seasons WHERE series_id = ? AND season_number = ?`,
		seriesId, seasonNumber,
	).Scan(
		&season.Id, &season.Name, &season.SeasonNumber,
		&season.EpisodeCount, &season.Size, &season.SeriesId,
		&season.SpaceSaved, &season.MissingEpisodes,
	)
	if err != nil {
		return models.Season{}, err
	}

	// Load episodes
	episodes, err := repo.getEpisodesBySeason(season.Id)
	if err != nil {
		return models.Season{}, err
	}
	season.Episodes = episodes

	return season, nil
}

func (repo *SeasonRepository) DeleteSeasonById(seriesId string, seasonNumber int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete episodes first
	_, err = tx.Exec(
		"DELETE FROM episodes WHERE series_id = ? AND season_number = ?",
		seriesId, seasonNumber,
	)
	if err != nil {
		return err
	}

	// Delete season
	_, err = tx.Exec(
		"DELETE FROM seasons WHERE series_id = ? AND season_number = ?",
		seriesId, seasonNumber,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *SeasonRepository) getEpisodesBySeason(seasonId string) ([]models.Episode, error) {
	rows, err := repo.DB.Query(`
        SELECT id, series_id, season_id, episode_number,
        season_name, season_number, filename, episode_name,
        video_codec, air_date, size, space_saved, original_size,
        path, missing
        FROM episodes WHERE season_id = ?`, seasonId)
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
