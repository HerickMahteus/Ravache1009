package models

type Turma struct {
	ID         string     `json:"id"`
	Nome       string     `json:"nome"`
	Disciplina string     `json:"disciplina"`
	Professor  string     `json:"professor"`
	Alunos     []string   `json:"alunos"`
	Alocacao   *Alocacao  `json:"alocacao"`
	Ativa      bool       `json:"ativa"`
}