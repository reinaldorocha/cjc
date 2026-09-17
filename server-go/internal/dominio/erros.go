package dominio

import "errors"

var (
	ErrEntradaInvalida = errors.New("entrada inválida")
	ErrNaoEncontrado   = errors.New("não encontrado")
	ErrConflito        = errors.New("conflito")
	ErrNaoAutorizado   = errors.New("não autorizado")
)
