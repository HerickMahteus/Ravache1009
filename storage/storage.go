package storage

import "api-gin/models"

var Salas = make(map[string]models.Sala)
var Alunos = make(map[string]models.Aluno)
var Turmas = make(map[string]models.Turma)