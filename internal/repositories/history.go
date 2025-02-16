package repositories

import (
	"database/sql"
	"transfigurr/internal/models"
)

type HistoryRepository struct {
	DB *sql.DB
}

func NewHistoryRepository(db *sql.DB) *HistoryRepository {
	return &HistoryRepository{
		DB: db,
	}
}

func (repo *HistoryRepository) GetHistories() ([]models.History, error) {
	rows, err := repo.DB.Query(`
        SELECT id, media_id, name, type, season_number,
        episode_number, profile_id, prev_codec, new_codec,
        prev_size, new_size, date
        FROM history
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []models.History
	for rows.Next() {
		var history models.History
		err := rows.Scan(
			&history.Id, &history.MediaId, &history.Name,
			&history.Type, &history.SeasonNumber, &history.EpisodeNumber,
			&history.ProfileId, &history.PrevCodec, &history.NewCodec,
			&history.PrevSize, &history.NewSize, &history.Date,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}

func (repo *HistoryRepository) GetHistoryById(id string) (models.History, error) {
	var history models.History
	err := repo.DB.QueryRow(`
        SELECT id, media_id, name, type, season_number,
        episode_number, profile_id, prev_codec, new_codec,
        prev_size, new_size, date
        FROM history WHERE id = ?`, id,
	).Scan(
		&history.Id, &history.MediaId, &history.Name,
		&history.Type, &history.SeasonNumber, &history.EpisodeNumber,
		&history.ProfileId, &history.PrevCodec, &history.NewCodec,
		&history.PrevSize, &history.NewSize, &history.Date,
	)

	if err == sql.ErrNoRows {
		return models.History{}, ErrRecordNotFound
	}
	if err != nil {
		return models.History{}, err
	}
	return history, nil
}

func (repo *HistoryRepository) UpsertHistoryById(history *models.History) error {
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM history WHERE id = ?)", history.Id).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = repo.DB.Exec(`
            UPDATE history SET
            media_id = ?, name = ?, type = ?,
            season_number = ?, episode_number = ?,
            profile_id = ?, prev_codec = ?,
            new_codec = ?, prev_size = ?,
            new_size = ?, date = ?
            WHERE id = ?`,
			history.MediaId, history.Name, history.Type,
			history.SeasonNumber, history.EpisodeNumber,
			history.ProfileId, history.PrevCodec,
			history.NewCodec, history.PrevSize,
			history.NewSize, history.Date,
			history.Id,
		)
	} else {
		_, err = repo.DB.Exec(`
            INSERT INTO history (
                media_id, name, type, season_number,
                episode_number, profile_id, prev_codec,
                new_codec, prev_size, new_size, date
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			history.MediaId, history.Name, history.Type,
			history.SeasonNumber, history.EpisodeNumber,
			history.ProfileId, history.PrevCodec,
			history.NewCodec, history.PrevSize,
			history.NewSize, history.Date,
		)
	}
	return err
}

func (repo *HistoryRepository) DeleteHistoryById(history *models.History) error {
	_, err := repo.DB.Exec("DELETE FROM history WHERE id = ?", history.Id)
	return err
}
