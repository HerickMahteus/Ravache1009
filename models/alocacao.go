package models

type Alocacao struct {
	SalaID      string `json:"sala_id"`
	DiaSemana   string `json:"dia_semana"`
	HorarioInicio string `json:"horario_inicio"`
	HorarioFim    string `json:"horario_fim"`
}