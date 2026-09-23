package handlers

import (
	"net/http"
	"time"

	"api-gin/models"
	"api-gin/storage"

	"github.com/gin-gonic/gin"
)

type TurmaHandler struct{}

func (h TurmaHandler) CriarTurma(c *gin.Context) {

	var turma models.Turma

	if err := c.ShouldBindJSON(&turma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "dados da turma inválidos",
		})
		return
	}

	if _, existe := storage.Turmas[turma.ID]; existe {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "já existe uma turma com esse identificador",
		})
		return
	}

	turma.Alunos = make([]string, 0)
	turma.Alocacao = nil
	turma.Ativa = true

	storage.Turmas[turma.ID] = turma

	c.JSON(http.StatusCreated, turma)
}

func (h TurmaHandler) ListarTurmas(c *gin.Context) {

	type TurmaResponse struct {
		ID                  string         `json:"id"`
		Nome                string         `json:"nome"`
		Disciplina          string         `json:"disciplina"`
		Professor           string         `json:"professor"`
		QuantidadeDeAlunos  int            `json:"quantidade_de_alunos"`
		Alocacao            *models.Alocacao `json:"alocacao"`
		Alocada             bool           `json:"alocada"`
		Ativa               bool           `json:"ativa"`
	}

	turmas := make([]TurmaResponse, 0, len(storage.Turmas))

	for _, turma := range storage.Turmas {

		turmaResponse := TurmaResponse{
			ID:                 turma.ID,
			Nome:               turma.Nome,
			Disciplina:         turma.Disciplina,
			Professor:          turma.Professor,
			QuantidadeDeAlunos: len(turma.Alunos),
			Alocacao:           turma.Alocacao,
			Alocada:            turma.Alocacao != nil,
			Ativa:              turma.Ativa,
		}

		turmas = append(turmas, turmaResponse)
	}

	c.JSON(http.StatusOK, turmas)
}

func (h TurmaHandler) AdicionarAluno(c *gin.Context) {

	turmaID := c.Param("id")

	var dados struct {
		AlunoID string `json:"aluno_id"`
	}

	if err := c.ShouldBindJSON(&dados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "dados da matrícula inválidos",
		})
		return
	}

	// Verifica se a turma existe
	turma, existe := storage.Turmas[turmaID]

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "turma não encontrada",
		})
		return
	}

	// Verifica se o aluno existe
	if _, existe := storage.Alunos[dados.AlunoID]; !existe {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "aluno não encontrado",
		})
		return
	}

	// Verifica se o aluno já está matriculado na turma
	for _, alunoID := range turma.Alunos {
		if alunoID == dados.AlunoID {
			c.JSON(http.StatusConflict, gin.H{
				"erro": "aluno já está matriculado nesta turma",
			})
			return
		}
	}

	// Adiciona o aluno à turma
	turma.Alunos = append(turma.Alunos, dados.AlunoID)

	// Atualiza a turma no storage
	storage.Turmas[turmaID] = turma

	c.JSON(http.StatusOK, gin.H{
		"mensagem": "aluno matriculado com sucesso",
		"turma":    turma,
	})
}

func (h TurmaHandler) ListarAlunosDaTurma(c *gin.Context) {

	turmaID := c.Param("id")

	// Verifica se a turma existe
	turma, existe := storage.Turmas[turmaID]

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "turma não encontrada",
		})
		return
	}

	// Busca os alunos cadastrados na turma
	alunos := make([]models.Aluno, 0, len(turma.Alunos))

	for _, alunoID := range turma.Alunos {

		aluno, existe := storage.Alunos[alunoID]

		if existe {
			alunos = append(alunos, aluno)
		}
	}

	c.JSON(http.StatusOK, alunos)
}

func (h TurmaHandler) AlocarSala(c *gin.Context) {

	turmaID := c.Param("id")

	var dados models.Alocacao

	if err := c.ShouldBindJSON(&dados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "dados da alocação inválidos",
		})
		return
	}

	// Verifica se a turma existe
	turma, existe := storage.Turmas[turmaID]

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "turma não encontrada",
		})
		return
	}

	// Verifica se a turma está ativa
	if !turma.Ativa {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro": "a turma está inativa",
		})
		return
	}

	// Verifica se a sala existe
	sala, existe := storage.Salas[dados.SalaID]

	if !existe {
		c.JSON(http.StatusNotFound, gin.H{
			"erro": "sala não encontrada",
		})
		return
	}

	// Verifica se a sala está ativa
	if !sala.Ativa {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro": "a sala está inativa",
		})
		return
	}

	// Valida os horários
	inicio, errInicio := time.Parse("15:04", dados.HorarioInicio)
	fim, errFim := time.Parse("15:04", dados.HorarioFim)

	if errInicio != nil || errFim != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "os horários devem estar no formato HH:MM",
		})
		return
	}

	if !inicio.Before(fim) {
		c.JSON(http.StatusBadRequest, gin.H{
			"erro": "o horário de início deve ser anterior ao horário de fim",
		})
		return
	}

	// Verifica se a sala possui capacidade suficiente
	if len(turma.Alunos) > sala.Capacidade {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"erro": "a capacidade da sala é insuficiente para a quantidade de alunos da turma",
		})
		return
	}

	// Verifica se a turma já possui uma alocação
	if turma.Alocacao != nil {
		c.JSON(http.StatusConflict, gin.H{
			"erro": "a turma já está alocada em uma sala",
		})
		return
	}

	// Verifica conflito de horário com outras turmas na mesma sala
	for _, outraTurma := range storage.Turmas {

		if outraTurma.Alocacao == nil {
			continue
		}

		alocacaoExistente := outraTurma.Alocacao

		if alocacaoExistente.SalaID != dados.SalaID {
			continue
		}

		if alocacaoExistente.DiaSemana != dados.DiaSemana {
			continue
		}

		existenteInicio, err1 := time.Parse(
			"15:04",
			alocacaoExistente.HorarioInicio,
		)

		existenteFim, err2 := time.Parse(
			"15:04",
			alocacaoExistente.HorarioFim,
		)

		if err1 != nil || err2 != nil {
			continue
		}

		// Regra de sobreposição:
		// novo início < fim existente
		// E novo fim > início existente
		if inicio.Before(existenteFim) && fim.After(existenteInicio) {

			c.JSON(http.StatusConflict, gin.H{
				"erro": "a sala já está ocupada nesse dia e horário",
			})
			return
		}
	}

	// Verifica conflito de horário dos alunos
	for _, alunoID := range turma.Alunos {

		for _, outraTurma := range storage.Turmas {

			if outraTurma.ID == turma.ID {
				continue
			}

			// A outra turma precisa estar alocada
			if outraTurma.Alocacao == nil {
				continue
			}

			// Verifica se o aluno também está na outra turma
			alunoNaOutraTurma := false

			for _, outroAlunoID := range outraTurma.Alunos {
				if outroAlunoID == alunoID {
					alunoNaOutraTurma = true
					break
				}
			}

			if !alunoNaOutraTurma {
				continue
			}

			alocacaoExistente := outraTurma.Alocacao

			// Só existe conflito se for no mesmo dia
			if alocacaoExistente.DiaSemana != dados.DiaSemana {
				continue
			}

			existenteInicio, err1 := time.Parse(
				"15:04",
				alocacaoExistente.HorarioInicio,
			)

			existenteFim, err2 := time.Parse(
				"15:04",
				alocacaoExistente.HorarioFim,
			)

			if err1 != nil || err2 != nil {
				continue
			}

			// Verifica sobreposição de horários
			if inicio.Before(existenteFim) && fim.After(existenteInicio) {

				c.JSON(http.StatusConflict, gin.H{
					"erro": "um dos alunos da turma já possui outra turma nesse horário",
				})
				return
			}
		}
	}

	// Todas as validações passaram.
	turma.Alocacao = &models.Alocacao{
		SalaID:        dados.SalaID,
		DiaSemana:     dados.DiaSemana,
		HorarioInicio: dados.HorarioInicio,
		HorarioFim:    dados.HorarioFim,
	}

	storage.Turmas[turmaID] = turma

	c.JSON(http.StatusOK, turma)
}