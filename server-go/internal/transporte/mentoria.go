package transporte

import (
	"net/http"
	"chega-junto-concurseiro-web/internal/modelos"
	"chega-junto-concurseiro-web/internal/usuarios"
)

func (s *Servidor) criarAlunoPeloMentor(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e struct {
		Nome                         string  `json:"nome"`
		Email                        string  `json:"email"`
		Senha                        string  `json:"senha"`
		Telefone                     string  `json:"telefone"`
		PermiteCronogramaInteligente bool    `json:"permiteCronogramaInteligente"`
		DataExpiracaoPlano           *string `json:"dataExpiracaoPlano"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	id, err := s.usuarios.CriarAlunoVinculado(r.Context(), usuarios.Criacao{Nome: e.Nome, Email: e.Email, Senha: e.Senha, Telefone: e.Telefone, Papel: "aluno"}, c.Usuario.ID, e.PermiteCronogramaInteligente, e.DataExpiracaoPlano)
	if err != nil {
		responderErro(w, 400, "ALUNO_INVALIDO", "Não foi possível cadastrar o aluno. Verifique os dados e se o e-mail já está em uso.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, id, "criar", "usuario_aluno", id, map[string]any{"vinculoAutomatico": true})
	responder(w, 201, map[string]any{"id": id})
}

func (s *Servidor) listarAlunos(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	lista, err := s.mentoria.ListarAlunos(r.Context(), c.Usuario.ID)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível listar alunos.")
		return
	}
	if pagina, limite, ok := paginacaoSolicitada(r); ok {
		itens, dados := paginar(lista, pagina, limite)
		responder(w, 200, map[string]any{"alunos": itens, "paginacao": dados})
		return
	}
	responder(w, 200, map[string]any{"alunos": lista})
}
func (s *Servidor) obterAluno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	aluno, err := s.mentoria.ObterAluno(r.Context(), c.Usuario.ID, r.PathValue("alunoId"))
	if err != nil {
		responderErro(w, 404, "ALUNO_NAO_ENCONTRADO", "Aluno não encontrado.")
		return
	}
	responder(w, 200, map[string]any{"aluno": aluno})
}
func (s *Servidor) visaoGeralAluno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	dados, err := s.mentoria.VisaoGeral(r.Context(), c.Usuario.ID, r.PathValue("alunoId"))
	if err != nil {
		responderErro(w, 404, "ALUNO_NAO_ENCONTRADO", "Aluno não encontrado.")
		return
	}
	responder(w, 200, dados)
}
func (s *Servidor) configurarAluno(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	var e struct {
		Nome                         *string `json:"nome"`
		Senha                        *string `json:"senha"`
		PermiteCronogramaInteligente *bool   `json:"permiteCronogramaInteligente"`
		DataExpiracaoPlano           *string `json:"dataExpiracaoPlano"`
		Telefone                     *string `json:"telefone"`
	}
	if !decodificar(w, r, &e) {
		return
	}
	if err := s.mentoria.Configurar(r.Context(), c.Usuario.ID, r.PathValue("alunoId"), e.PermiteCronogramaInteligente, e.DataExpiracaoPlano, e.Telefone, e.Nome, e.Senha); err != nil {
		responderErro(w, 400, "DADOS_INVALIDOS", "Não foi possível atualizar a configuração do aluno.")
		return
	}
	s.registrarAuditoria(r.Context(), c.Usuario.ID, r.PathValue("alunoId"), "alterar", "configuracao_aluno", r.PathValue("alunoId"), nil)
	responder(w, 200, map[string]any{"sucesso": true})
}

func (s *Servidor) radarAlunos(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado) {
	lista, err := s.mentoria.RadarAlunos(r.Context(), c.Usuario.ID)
	if err != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível carregar o radar de alunos.")
		return
	}
	responder(w, 200, map[string]any{"alunos": lista})
}
