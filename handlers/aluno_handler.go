package handlers

import (
	"net/http"

	"api-gin/models"
	"api-gin/storage"

	"github.com/gin-gonic/gin"
)

type AlunoHandler struct{}

func (h AlunoHandler) CriarAluno(c *gin.Context) {

	var aluno models.Aluno

	if err := c.ShouldBindJSON(&aluno); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "dados do aluno inválidos",
		})
		return
	}

	if _, existe := storage.Alunos[aluno.ID]; existe {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "já existe um aluno com esse identificador",
		})
		return
	}

	storage.Alunos[aluno.ID] = aluno

	c.JSON(http.StatusCreated, aluno)
}

func (h AlunoHandler) ListarAlunos(c *gin.Context) {

	alunos := make([]models.Aluno, 0, len(storage.Alunos))

	for _, aluno := range storage.Alunos {
		alunos = append(alunos, aluno)
	}

	c.JSON(http.StatusOK, alunos)
}

func (h AlunoHandler) BuscarAluno(c *gin.Context) {

	id := c.Param("id")

	aluno, existe := storage.Alunos[id]

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "aluno não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, aluno)
}