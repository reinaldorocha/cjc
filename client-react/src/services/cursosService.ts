import { request, jsonBody } from './client';

export interface ProgressoAula { posicao: number; duracao: number; concluida: boolean; atualizadoEm: string }
export interface ResumoCurso { total: number; concluidas: number; percentual: number; retomarAulaId: string; ultimaAtividade: string }
export interface AulaCurso { id?: string; moduloId?: string; moduloCapaUrl?: string; moduloDescricao?: string; titulo: string; modulo: string; videoId: string; pdfId?: string; pdfNome?: string; liberarEm?: string; exigeAnterior?: boolean; questoes?: string[]; bloqueada?: boolean; motivoBloqueio?: string; progresso?: ProgressoAula }
export interface QuestaoCurso { id: string; disciplina: string; assunto: string; tipo: string; enunciado: string; alternativas: string[] }
export interface ResultadoQuestao { correta: boolean; respostaCorreta: string; explicacao: string }
export interface AlunoCurso { id: string; nome: string; email: string; resumo: ResumoCurso; situacao: string; respostas: number; acertos: number }
export interface Curso {
  id: string; titulo: string; descricao: string; categoria: string; capaUrl: string;
  modoExibicao?: 'curso' | 'modulos';
  escopo: 'global' | 'alunos' | 'concursos'; destinatarios: string[]; aulas: AulaCurso[];
  revisao?: number; publicado?: boolean; alteracoesPendentes?: boolean; resumo?: ResumoCurso;
}
export const cursosService = {
  enviarCapa: (arquivo: File) => { const body = new FormData(); body.append('arquivo', arquivo); return request<{ url: string }>('/mentor/cursos/capas', { method: 'POST', body }); },
  enviarPDF: (arquivo: File) => { const body = new FormData(); body.append('arquivo', arquivo); return request<{ id: string; nome: string }>('/mentor/cursos/pdfs', { method: 'POST', body }); },
  pdfURL: (cursoId: string, pdfId: string, mentor: boolean) => `/api/v1/${mentor ? 'mentor' : 'alunos/eu'}/cursos/${encodeURIComponent(cursoId)}/pdfs/${encodeURIComponent(pdfId)}`,
  listar: (mentor: boolean, alunoId = 'eu') => request<{ cursos: Curso[] }>(mentor ? '/mentor/cursos' : `/alunos/${encodeURIComponent(alunoId)}/cursos`),
  obter: (id: string, mentor: boolean, alunoId = 'eu') => request<{ curso: Curso }>(`${mentor ? '/mentor' : `/alunos/${encodeURIComponent(alunoId)}`}/cursos/${encodeURIComponent(id)}`),
  salvar: (curso: Curso) => request<{ id: string; curso: Curso }>(`/mentor/cursos${curso.id ? `/${encodeURIComponent(curso.id)}` : ''}`, { method: curso.id ? 'PUT' : 'POST', body: jsonBody(curso) }),
  publicar: (curso: Curso, publicar: boolean) => request<{ curso: Curso }>(`/mentor/cursos/${encodeURIComponent(curso.id)}/publicacao`, { method: 'POST', body: jsonBody({ revisao: curso.revisao, publicar }) }),
  previa: (id: string) => request<{ curso: Curso }>(`/mentor/cursos/${encodeURIComponent(id)}/previa`),
  progresso: (id: string, aulaId: string, dados: { posicao: number; duracao: number; concluida?: boolean }, keepalive = false) => request<{ curso: Curso }>(`/cursos/${encodeURIComponent(id)}/aulas/${encodeURIComponent(aulaId)}/progresso`, { method: 'PUT', body: jsonBody(dados), keepalive }),
  acompanhar: (id: string) => request<{ alunos: AlunoCurso[] }>(`/mentor/cursos/${encodeURIComponent(id)}/acompanhamento`),
  catalogoQuestoes: (busca: string) => request<{ questoes: QuestaoCurso[] }>(`/mentor/cursos/questoes?busca=${encodeURIComponent(busca)}`),
  questoes: (id: string, aulaId: string) => request<{ questoes: QuestaoCurso[] }>(`/cursos/${encodeURIComponent(id)}/aulas/${encodeURIComponent(aulaId)}/questoes`),
  responder: (id: string, aulaId: string, questaoId: string, resposta: string) => request<ResultadoQuestao>(`/cursos/${encodeURIComponent(id)}/aulas/${encodeURIComponent(aulaId)}/questoes/${encodeURIComponent(questaoId)}/responder`, { method: 'POST', body: jsonBody({ resposta }) }),
  excluir: (id: string) => request(`/mentor/cursos/${encodeURIComponent(id)}`, { method: 'DELETE' }),
};
