package handlers

import (
	"net/http"

	"api-gin/models"
	"api-gin/storage"

	"github.com/gin-gonic/gin"
)

type SalaHandler struct{}

func (h SalaHandler) CriarSala(c *gin.Context) {

	var sala models.Sala

	if err := c.ShouldBindJSON(&sala); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "dados da sala inválidos",
		})
		return
	}

	if sala.Capacidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "a capacidade da sala deve ser maior que zero",
		})
		return
	}

	if _, existe := storage.Salas[sala.ID]; existe {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "já existe uma sala com esse identificador",
		})
		return
	}

	sala.Ativa = true

	storage.Salas[sala.ID] = sala

	c.JSON(http.StatusCreated, sala)
}

func (h SalaHandler) ListarSalas(c *gin.Context) {

	salas := make([]models.Sala, 0, len(storage.Salas))

	for _, sala := range storage.Salas {
		salas = append(salas, sala)
	}

	c.JSON(http.StatusOK, salas)
}