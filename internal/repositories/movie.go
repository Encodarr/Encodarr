package repositories

import (
	"database/sql"
	"transfigurr/internal/models"
)

type MovieRepository struct {
	DB *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		DB: db,
	}
}

func (repo *MovieRepository) GetMovies() ([]models.Movie, error) {
	rows, err := repo.DB.Query(`
        SELECT id, name, release_date, genre, status,
        filename, video_codec, overview, size, space_saved,
        profile_id, monitored, missing, studio,
        original_size, path, runtime
        FROM movies
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.Id, &movie.Name, &movie.ReleaseDate,
			&movie.Genre, &movie.Status, &movie.Filename,
			&movie.VideoCodec, &movie.Overview, &movie.Size,
			&movie.SpaceSaved, &movie.ProfileID, &movie.Monitored,
			&movie.Missing, &movie.Studio, &movie.OriginalSize,
			&movie.Path, &movie.Runtime,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (repo *MovieRepository) UpsertMovie(id string, movie models.Movie) (models.Movie, error) {
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM movies WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return models.Movie{}, err
	}

	if exists {
		_, err = repo.DB.Exec(`
            UPDATE movies SET
            name = ?, release_date = ?, genre = ?,
            status = ?, filename = ?, video_codec = ?,
            overview = ?, size = ?, space_saved = ?,
            profile_id = ?, monitored = ?, missing = ?,
            studio = ?, original_size = ?, path = ?,
            runtime = ?
            WHERE id = ?`,
			movie.Name, movie.ReleaseDate, movie.Genre,
			movie.Status, movie.Filename, movie.VideoCodec,
			movie.Overview, movie.Size, movie.SpaceSaved,
			movie.ProfileID, movie.Monitored, movie.Missing,
			movie.Studio, movie.OriginalSize, movie.Path,
			movie.Runtime, id,
		)
	} else {
		_, err = repo.DB.Exec(`
            INSERT INTO movies (
                id, name, release_date, genre, status,
                filename, video_codec, overview, size,
                space_saved, profile_id, monitored,
                missing, studio, original_size,
                path, runtime
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, movie.Name, movie.ReleaseDate, movie.Genre,
			movie.Status, movie.Filename, movie.VideoCodec,
			movie.Overview, movie.Size, movie.SpaceSaved,
			movie.ProfileID, movie.Monitored, movie.Missing,
			movie.Studio, movie.OriginalSize, movie.Path,
			movie.Runtime,
		)
	}
	if err != nil {
		return models.Movie{}, err
	}

	return repo.GetMovieById(id)
}

func (repo *MovieRepository) GetMovieById(id string) (models.Movie, error) {
	var movie models.Movie
	err := repo.DB.QueryRow(`
        SELECT id, name, release_date, genre, status,
        filename, video_codec, overview, size, space_saved,
        profile_id, monitored, missing, studio,
        original_size, path, runtime
        FROM movies WHERE id = ?`, id,
	).Scan(
		&movie.Id, &movie.Name, &movie.ReleaseDate,
		&movie.Genre, &movie.Status, &movie.Filename,
		&movie.VideoCodec, &movie.Overview, &movie.Size,
		&movie.SpaceSaved, &movie.ProfileID, &movie.Monitored,
		&movie.Missing, &movie.Studio, &movie.OriginalSize,
		&movie.Path, &movie.Runtime,
	)
	if err != nil {
		return models.Movie{}, err
	}
	return movie, nil
}

func (repo *MovieRepository) DeleteMovieById(id string) error {
	_, err := repo.DB.Exec("DELETE FROM movies WHERE id = ?", id)
	return err
}
