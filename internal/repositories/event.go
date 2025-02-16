package repositories

import (
	"database/sql"
	"time"
	"transfigurr/internal/models"
)

type EventRepository struct {
	DB *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{DB: db}
}

func (repo *EventRepository) GetEvents() ([]models.Event, error) {
	rows, err := repo.DB.Query(`
        SELECT id, timestamp, level, service, message 
        FROM events
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.Id, &event.Timestamp, &event.Level,
			&event.Service, &event.Message,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (repo *EventRepository) GetEventById(id string) (models.Event, error) {
	var event models.Event
	err := repo.DB.QueryRow(`
        SELECT id, timestamp, level, service, message 
        FROM events WHERE id = ?`, id,
	).Scan(&event.Id, &event.Timestamp, &event.Level, &event.Service, &event.Message)

	if err == sql.ErrNoRows {
		return models.Event{}, ErrRecordNotFound
	}
	if err != nil {
		return models.Event{}, err
	}
	return event, nil
}

func (repo *EventRepository) UpsertEventById(event models.Event) error {
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM events WHERE id = ?)", event.Id).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = repo.DB.Exec(`
            UPDATE events SET 
            timestamp = ?, level = ?, service = ?, message = ?
            WHERE id = ?`,
			event.Timestamp, event.Level, event.Service, event.Message, event.Id,
		)
	} else {
		_, err = repo.DB.Exec(`
            INSERT INTO events (timestamp, level, service, message)
            VALUES (?, ?, ?, ?)`,
			event.Timestamp, event.Level, event.Service, event.Message,
		)
	}
	return err
}

func (repo *EventRepository) DeleteEventById(event models.Event) error {
	_, err := repo.DB.Exec("DELETE FROM events WHERE id = ?", event.Id)
	return err
}

func (repo *EventRepository) Log(level, service, message string) error {
	eventEntry := models.Event{
		Timestamp: time.Now().Format("2006-01-02T15:04:05.000"),
		Level:     level,
		Service:   service,
		Message:   message,
	}

	_, err := repo.DB.Exec(`
        INSERT INTO events (timestamp, level, service, message)
        VALUES (?, ?, ?, ?)`,
		eventEntry.Timestamp, eventEntry.Level,
		eventEntry.Service, eventEntry.Message,
	)
	return err
}
