import { request, jsonBody } from './client';

export const flashcardsService = {
  listarBaralhos: (alunoId = 'eu', concursoId = '') =>
    request<any>(`/alunos/${alunoId}/baralhos-cartoes${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
  criarBaralho: (alunoId: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/baralhos-cartoes`, { method: 'POST', body: jsonBody(dados) }),
  alterarBaralho: (alunoId: string, id: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/baralhos-cartoes/${id}`, { method: 'PATCH', body: jsonBody(dados) }),
  desativarBaralho: (alunoId: string, id: string) =>
    request<any>(`/alunos/${alunoId}/baralhos-cartoes/${id}`, { method: 'DELETE' }),
  criarCartao: (alunoId: string, baralhoId: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/baralhos-cartoes/${baralhoId}/cartoes`, { method: 'POST', body: jsonBody(dados) }),
  alterarCartao: (alunoId: string, id: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/cartoes/${id}`, { method: 'PATCH', body: jsonBody(dados) }),
  desativarCartao: (alunoId: string, id: string) =>
    request<any>(`/alunos/${alunoId}/cartoes/${id}`, { method: 'DELETE' }),
  revisarCartao: (alunoId: string, id: string, qualidade: number, concursoId = '') =>
    request<any>(`/alunos/${alunoId}/cartoes/${id}/revisar`, { method: 'POST', body: jsonBody({ qualidade, concursoId }) }),
  listarCartoesPendentes: (alunoId = 'eu') =>
    request<any>(`/alunos/${alunoId}/cartoes/pendentes`),
  listarHistoricoCartoes: (alunoId = 'eu', concursoId = '') =>
    request<any>(`/alunos/${alunoId}/cartoes/historico${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
};
