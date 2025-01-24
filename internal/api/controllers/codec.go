package controllers

import (
	"encoding/json"
	"net/http"
	"transfigurr/internal/interfaces/repositories"
)

type CodecController struct {
	Repo repositories.CodecRepositoryInterface
}

func NewCodecController(repo repositories.CodecRepositoryInterface) *CodecController {
	return &CodecController{
		Repo: repo,
	}
}

func (ctrl CodecController) GetCodecs(w http.ResponseWriter, r *http.Request) {
	codecs := ctrl.Repo.GetCodecs()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).SetIndent("", "    ")
	json.NewEncoder(w).Encode(codecs)
}

func (ctrl CodecController) GetContainers(w http.ResponseWriter, r *http.Request) {
	containers := ctrl.Repo.GetContainers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).SetIndent("", "    ")
	json.NewEncoder(w).Encode(containers)
}

func (ctrl CodecController) GetEncoders(w http.ResponseWriter, r *http.Request) {
	encoders := ctrl.Repo.GetEncoders()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).SetIndent("", "    ")
	json.NewEncoder(w).Encode(encoders)
}
