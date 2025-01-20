package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/repository"
)

type EventController struct {
	Repo interfaces.EventRepositoryInterface
}

func NewEventController(repo interfaces.EventRepositoryInterface) *EventController {
	return &EventController{
		Repo: repo,
	}
}

func (ctrl EventController) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := ctrl.Repo.GetEvents()
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			http.Error(w, "Events not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (ctrl EventController) UpsertEvent(w http.ResponseWriter, r *http.Request, eventId string) {
	var inputEvent models.Event
	if err := json.NewDecoder(r.Body).Decode(&inputEvent); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event, err := ctrl.Repo.GetEventById(eventId)
	if err != nil && errors.Is(err, repository.ErrRecordNotFound) {
		event = inputEvent
	} else if err != nil {
		http.Error(w, "Error retrieving event", http.StatusInternalServerError)
		return
	}

	if err := ctrl.Repo.UpsertEventById(event); err != nil {
		http.Error(w, "Error upserting event", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(event)
}

func (ctrl EventController) GetEventById(w http.ResponseWriter, r *http.Request, eventId string) {
	event, err := ctrl.Repo.GetEventById(eventId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			http.Error(w, "Event not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving event", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(event)
}

func (ctrl EventController) DeleteEventById(w http.ResponseWriter, r *http.Request, eventId string) {
	event, err := ctrl.Repo.GetEventById(eventId)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	if err := ctrl.Repo.DeleteEventById(event); err != nil {
		http.Error(w, "Error deleting event", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Event deleted successfully"})
}
