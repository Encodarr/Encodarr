package repositories

import (
	"database/sql"
	"transfigurr/internal/models"
)

type SeriesRepository struct {
	DB *sql.DB
}

func NewSeriesRepository(db *sql.DB) *SeriesRepository {
	return &SeriesRepository{
		DB: db,
	}
}

func (repo *SeriesRepository) GetSeries() ([]models.Series, error) {
	rows, err := repo.DB.Query(`
        SELECT id, name, release_date, genre, status, last_air_date, 
        networks, overview, profile_id, monitored, episode_count, 
        size, seasons_count, space_saved, missing_episodes, runtime
        FROM series
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []models.Series
	for rows.Next() {
		var series models.Series
		err := rows.Scan(
			&series.Id, &series.Name, &series.ReleaseDate, &series.Genre,
			&series.Status, &series.LastAirDate, &series.Networks,
			&series.Overview, &series.ProfileID, &series.Monitored,
			&series.EpisodeCount, &series.Size, &series.SeasonsCount,
			&series.SpaceSaved, &series.MissingEpisodes, &series.Runtime,
		)
		if err != nil {
			return nil, err
		}

		// Load seasons and episodes
		seasons, err := repo.getSeasonsBySeries(series.Id)
		if err != nil {
			return nil, err
		}
		series.Seasons = seasons
		seriesList = append(seriesList, series)
	}
	return seriesList, nil
}

func (repo *SeriesRepository) UpsertSeries(id string, inputSeries models.Series) (models.Series, error) {
	// Check if series exists
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM series WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return models.Series{}, err
	}

	if exists {
		// Update
		_, err = repo.DB.Exec(`
            UPDATE series SET
            name = ?, release_date = ?, genre = ?, status = ?,
            last_air_date = ?, networks = ?, overview = ?,
            profile_id = ?, monitored = ?, episode_count = ?,
            size = ?, seasons_count = ?, space_saved = ?,
            missing_episodes = ?, runtime = ?
            WHERE id = ?`,
			inputSeries.Name, inputSeries.ReleaseDate, inputSeries.Genre,
			inputSeries.Status, inputSeries.LastAirDate, inputSeries.Networks,
			inputSeries.Overview, inputSeries.ProfileID, inputSeries.Monitored,
			inputSeries.EpisodeCount, inputSeries.Size, inputSeries.SeasonsCount,
			inputSeries.SpaceSaved, inputSeries.MissingEpisodes, inputSeries.Runtime,
			id,
		)
	} else {
		// Insert
		_, err = repo.DB.Exec(`
            INSERT INTO series (
                id, name, release_date, genre, status, last_air_date,
                networks, overview, profile_id, monitored, episode_count,
                size, seasons_count, space_saved, missing_episodes, runtime
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, inputSeries.Name, inputSeries.ReleaseDate, inputSeries.Genre,
			inputSeries.Status, inputSeries.LastAirDate, inputSeries.Networks,
			inputSeries.Overview, inputSeries.ProfileID, inputSeries.Monitored,
			inputSeries.EpisodeCount, inputSeries.Size, inputSeries.SeasonsCount,
			inputSeries.SpaceSaved, inputSeries.MissingEpisodes, inputSeries.Runtime,
		)
	}
	if err != nil {
		return models.Series{}, err
	}

	return repo.GetSeriesByID(id)
}

func (repo *SeriesRepository) GetSeriesByID(id string) (models.Series, error) {
	var series models.Series
	err := repo.DB.QueryRow(`
        SELECT id, name, release_date, genre, status, last_air_date,
        networks, overview, profile_id, monitored, episode_count,
        size, seasons_count, space_saved, missing_episodes, runtime
        FROM series WHERE id = ?`, id,
	).Scan(
		&series.Id, &series.Name, &series.ReleaseDate, &series.Genre,
		&series.Status, &series.LastAirDate, &series.Networks,
		&series.Overview, &series.ProfileID, &series.Monitored,
		&series.EpisodeCount, &series.Size, &series.SeasonsCount,
		&series.SpaceSaved, &series.MissingEpisodes, &series.Runtime,
	)
	if err != nil {
		return models.Series{}, err
	}

	// Load seasons and episodes
	seasons, err := repo.getSeasonsBySeries(id)
	if err != nil {
		return models.Series{}, err
	}
	series.Seasons = seasons

	return series, nil
}

func (repo *SeriesRepository) DeleteSeriesByID(id string) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete episodes first
	_, err = tx.Exec("DELETE FROM episodes WHERE series_id = ?", id)
	if err != nil {
		return err
	}

	// Delete seasons
	_, err = tx.Exec("DELETE FROM seasons WHERE series_id = ?", id)
	if err != nil {
		return err
	}

	// Delete series
	_, err = tx.Exec("DELETE FROM series WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *SeriesRepository) getSeasonsBySeries(seriesID string) ([]models.Season, error) {
	rows, err := repo.DB.Query(`
        SELECT id, name, season_number, episode_count, size,
        series_id, space_saved, missing_episodes
        FROM seasons WHERE series_id = ?`, seriesID,
	)
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

		// Load episodes for each season
		episodes, err := repo.getEpisodesBySeason(season.Id)
		if err != nil {
			return nil, err
		}
		season.Episodes = episodes
		seasons = append(seasons, season)
	}
	return seasons, nil
}

func (repo *SeriesRepository) getEpisodesBySeason(seasonID string) ([]models.Episode, error) {
	rows, err := repo.DB.Query(`
        SELECT id, series_id, season_id, episode_number,
        season_name, season_number, filename, episode_name,
        video_codec, air_date, size, space_saved, original_size,
        path, missing
        FROM episodes WHERE season_id = ?`, seasonID,
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
