import { request, jsonBody } from './client';

export const concursoService = {
  listarConcursos: (alunoId = 'eu') => request<any[]>(`/alunos/${alunoId}/concursos`),
  criarConcurso: (alunoId: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/concursos`, { method: 'POST', body: jsonBody(dados) }),
  alterarConcurso: (alunoId: string, concursoId: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/concursos/${concursoId}`, { method: 'PATCH', body: jsonBody(dados) }),
  desativarConcurso: (alunoId: string, concursoId: string) =>
    request<any>(`/alunos/${alunoId}/concursos/${concursoId}`, { method: 'DELETE' }),
  reordenarConcursos: (alunoId: string, itens: any[]) =>
    request<any>(`/alunos/${alunoId}/concursos/ordem`, { method: 'PUT', body: jsonBody({ itens }) }),
  listarEditais: (alunoId = 'eu', concursoId = '') =>
    request<any>(`/alunos/${alunoId}/edital${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
};
