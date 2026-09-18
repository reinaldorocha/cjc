package transporte

import (
	"io"
	"net/http"
	"strings"

	"chega-junto-concurseiro-web/internal/arquivos"
	"chega-junto-concurseiro-web/internal/materiaisapoio"
	"chega-junto-concurseiro-web/internal/modelos"
)

func salvarMaterialArquivo(w http.ResponseWriter, r *http.Request, c modelos.ContextoAutenticado, s *Servidor) {
	r.Body = http.MaxBytesReader(w, r.Body, arquivos.LimiteMaterial+(1<<20))
	defer r.Body.Close()
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		responderErro(w, 400, "ENTRADA_INVALIDA", "Corpo do formulário inválido.")
		return
	}

	file, header, err := r.FormFile("arquivo")
	if err != nil && err != http.ErrMissingFile {
		responderErro(w, 400, "ENTRADA_INVALIDA", "Erro ao ler o arquivo enviado.")
		return
	}
	var arquivoNome string
	var arquivoConteudo io.Reader
	if err == nil {
		arquivoNome = header.Filename
		arquivoConteudo = file
	}
	entrada := materiaisapoio.EntradaArquivo{
		Titulo:      r.FormValue("titulo"),
		Descricao:   optionalString(r.FormValue("descricao")),
		Tipo:        "arquivo",
		Escopo:      r.FormValue("escopo"),
		EditalID:    optionalString(r.FormValue("editalId")),
		NomeArquivo: arquivoNome,
		Pasta:       optionalString(r.FormValue("pasta")),
	}

	if r.Method == http.MethodPatch {
		materialID := r.PathValue("materialId")
		errServico := s.materiaisApoio.AlterarArquivo(r.Context(), c.Usuario.ID, materialID, entrada, arquivoConteudo, s.cfg.UploadsPath)
		var errFechamento error
		if file != nil {
			errFechamento = file.Close()
		}
		if errServico != nil {
			responderErro(w, 400, "MATERIAL_INVALIDO", "Não foi possível atualizar o material de apoio de arquivo.")
			return
		}
		if errFechamento != nil {
			responderErro(w, 500, "ERRO_INTERNO", "Não foi possível finalizar a leitura do arquivo enviado.")
			return
		}
		responder(w, 200, map[string]any{"atualizado": true})
		return
	}

	if arquivoConteudo == nil {
		responderErro(w, 400, "ENTRADA_INVALIDA", "Arquivo não enviado.")
		return
	}

	id, err := s.materiaisApoio.CriarArquivo(r.Context(), c.Usuario.ID, entrada, arquivoConteudo, s.cfg.UploadsPath)
	errFechamento := file.Close()
	if err != nil {
		responderErro(w, 400, "MATERIAL_INVALIDO", "Não foi possível criar o material de apoio de arquivo.")
		return
	}
	if errFechamento != nil {
		responderErro(w, 500, "ERRO_INTERNO", "Não foi possível finalizar a leitura do arquivo enviado.")
		return
	}
	responder(w, 201, map[string]any{"id": id})
}

func optionalString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	v := strings.TrimSpace(value)
	return &v
}
