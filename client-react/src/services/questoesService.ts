import { request, jsonBody } from './client';

export const questoesService = {
  listarQuestoes: (alunoId = 'eu', concursoId = '') =>
    request<any>(`/alunos/${alunoId}/registros-questoes${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
  criarQuestoes: (alunoId: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/registros-questoes`, { method: 'POST', body: jsonBody(dados) }),
  alterarQuestoes: (alunoId: string, id: string, dados: any) =>
    request<any>(`/alunos/${alunoId}/registros-questoes/${id}`, { method: 'PATCH', body: jsonBody(dados) }),
  desativarQuestoes: (alunoId: string, id: string) =>
    request<any>(`/alunos/${alunoId}/registros-questoes/${id}`, { method: 'DELETE' }),
  listarBancoQuestoes: (params?: { disciplina?: string; assunto?: string; tipo?: string; concursoId?: string; alunoId?: string; pagina?: number; limite?: number }) => {
    const q = new URLSearchParams();
    if (params?.disciplina) q.set('disciplina', params.disciplina);
    if (params?.assunto) q.set('assunto', params.assunto);
    if (params?.tipo) q.set('tipo', params.tipo);
    if (params?.concursoId) q.set('concursoId', params.concursoId);
    if (params?.alunoId) q.set('alunoId', params.alunoId);
    if (params?.pagina) q.set('pagina', String(params.pagina));
    if (params?.limite) q.set('limite', String(params.limite));
    const queryStr = q.toString();
    return request<any>(`/banco-questoes${queryStr ? `?${queryStr}` : ''}`);
  },
  responderQuestaoBanco: (alunoId = 'eu', dados: { questaoId: string; respostaAluno: string; concursoId?: string }) =>
    request<any>(`/alunos/${alunoId}/banco-questoes/responder`, { method: 'POST', body: jsonBody(dados) }),
  historicoQuestaoBanco: (alunoId = 'eu', questaoId: string, concursoId = '') =>
    request<any>(`/alunos/${alunoId}/banco-questoes/${encodeURIComponent(questaoId)}/historico${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
  estatisticasBancoQuestoes: (alunoId = 'eu', concursoId = '') =>
    request<any>(`/alunos/${alunoId}/banco-questoes/estatisticas${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
};
