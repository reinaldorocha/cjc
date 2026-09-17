package modelos

type Usuario struct {
	ID                           string  `json:"id"`
	Nome                         string  `json:"nome"`
	Email                        string  `json:"email"`
	Papel                        string  `json:"papel"`
	Ativo                        bool    `json:"ativo"`
	Expirado                     bool    `json:"expirado"`
	PermiteCronogramaInteligente bool    `json:"permiteCronogramaInteligente"`
	DataExpiracaoPlano           *string `json:"dataExpiracaoPlano,omitempty"`
}

type Sessao struct {
	Usuario   Usuario
	Token     string
	TokenCSRF string
}

type ContextoAutenticado struct {
	Usuario  Usuario
	CSRFHash string
}
