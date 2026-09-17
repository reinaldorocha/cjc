import { request, jsonBody } from './client';

export const cronogramaService = {
  obterCronograma: (alunoId = 'eu', concursoId = '') =>
    request<any>(`/alunos/${alunoId}/cronograma${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
  salvarCronograma: (alunoId: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/cronograma`, { method: 'PUT', body: jsonBody(dados) }),
  gerarCronograma: (alunoId: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/cronograma/gerar`, { method: 'POST', body: jsonBody(dados) }),
  criarItemCronograma: (alunoId: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/cronograma/itens`, { method: 'POST', body: jsonBody(dados) }),
  reprogramarCronograma: (alunoId: string, dados?: any) =>
    request<any>(`/alunos/${alunoId}/cronograma/reprogramar`, { method: 'POST', body: jsonBody(dados || {}) }),
  alterarItemCronograma: (alunoId: string, itemId: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/cronograma/itens/${itemId}`, { method: 'PATCH', body: jsonBody(dados) }),
  calendarioCronograma: (alunoId: string, inicio = '', fim = '', concursoId = '') =>
    request<any>(`/alunos/${alunoId}/cronograma/calendario?inicio=${encodeURIComponent(inicio)}&fim=${encodeURIComponent(fim)}${concursoId ? `&concursoId=${encodeURIComponent(concursoId)}` : ''}`),
};
