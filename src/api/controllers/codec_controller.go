package controllers

import (
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
)

type CodecController struct {
	Repo *repository.CodecRepository
}

func NewCodecController(repo *repository.CodecRepository) *CodecController {
	return &CodecController{
		Repo: repo,
	}
}

func (ctrl CodecController) GetCodecs(c *gin.Context) {
	codecs := ctrl.Repo.GetCodecs()
	c.IndentedJSON(200, codecs)
}

func (ctrl CodecController) GetContainers(c *gin.Context) {
	containers := ctrl.Repo.GetContainers()
	c.IndentedJSON(200, containers)
}

func (ctrl CodecController) GetEncoders(c *gin.Context) {
	encoders := ctrl.Repo.GetEncoders()
	c.IndentedJSON(200, encoders)
}
