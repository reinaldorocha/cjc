import { authService } from './authService';
import { concursoService } from './concursoService';
import { cronogramaService } from './cronogramaService';
import { questoesService } from './questoesService';
import { flashcardsService } from './flashcardsService';
import { request, jsonBody, ErroApi } from './client';

export { ErroApi };

function avisarReplanejamento<T extends { avisoReplanejamento?: string }>(dados: T): T {
  if (dados.avisoReplanejamento) window.alert(dados.avisoReplanejamento);
  return dados;
}

export const api = {
  ...authService,
  ...concursoService,
  ...cronogramaService,
  ...questoesService,
  ...flashcardsService,

  // Mentor / Admin & Extra services
  listarAlunos: (pagina?: number, limite?: number) => request<any>(`/mentor/alunos${pagina ? `?pagina=${pagina}&limite=${limite || 12}` : ''}`),
  obterRadarAlunos: () => request<any>('/mentor/radar-alunos'),
  criarAlunoComoMentor: (dados: any) => request<any>('/mentor/alunos', { method: 'POST', body: jsonBody(dados) }),
  listarConcursosMentor: () => request<any>('/mentor/concursos'),
  criarConcursoMentor: (dados: any) => request<any>('/mentor/concursos', { method: 'POST', body: jsonBody(dados) }),
  listarEditaisMentor: (concursoId = '') => request<any>(`/mentor/editais${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
  criarEditalMentor: (dados: any) => request<any>('/mentor/editais', { method: 'POST', body: jsonBody(dados) }),
  listarBaralhosMentor: () => request<any>('/mentor/baralhos-cartoes'),
  criarBaralhoMentor: (dados: any) => request<any>('/mentor/baralhos-cartoes', { method: 'POST', body: jsonBody(dados) }),
  alterarBaralhoMentor: (id: string, dados: any) => request<any>(`/mentor/baralhos-cartoes/${id}`, { method: 'PATCH', body: jsonBody(dados) }),
  desativarBaralhoMentor: (id: string) => request<any>(`/mentor/baralhos-cartoes/${id}`, { method: 'DELETE' }),
  criarCartaoMentor: (baralhoId: string, dados: any) => request<any>(`/mentor/baralhos-cartoes/${baralhoId}/cartoes`, { method: 'POST', body: jsonBody(dados) }),
  alterarCartaoMentor: (id: string, dados: any) => request<any>(`/mentor/cartoes/${id}`, { method: 'PATCH', body: jsonBody(dados) }),
  desativarCartaoMentor: (id: string) => request<any>(`/mentor/cartoes/${id}`, { method: 'DELETE' }),
  obterAluno: (id: string) => request<any>(`/mentor/alunos/${id}`),
  visaoGeralAluno: (id: string) => request<any>(`/mentor/alunos/${id}/visao-geral`),
  configurarAluno: (id: string, dados: any) => request<any>(`/mentor/alunos/${id}/configuracoes`, { method: 'PATCH', body: jsonBody(dados) }),

  listarMateriaisApoio: (alunoId = 'eu', editalId = '') => request<any>(`/alunos/${alunoId}/materiais-apoio${editalId ? `?editalId=${encodeURIComponent(editalId)}` : ''}`),
  listarMateriaisApoioMentor: (editalId = '') => request<any>(`/mentor/materiais-apoio${editalId ? `?editalId=${encodeURIComponent(editalId)}` : ''}`),
  catalogarEditais: (concursoId: string) => request<any>(`/editais?concursoId=${encodeURIComponent(concursoId)}`),
  atribuirEdital: (alunoId: string, dados: any) => request<any>(`/alunos/${alunoId}/edital`, { method: 'PUT', body: jsonBody(dados) }),
  criarItemEdital: (alunoId: string, dados: any) => request<any>(`/alunos/${alunoId}/edital/itens`, { method: 'POST', body: jsonBody(dados) }),
  alterarItemEdital: (alunoId: string, itemId: string, dados: any) => request<any>(`/alunos/${alunoId}/edital/itens/${itemId}`, { method: 'PATCH', body: jsonBody(dados) }),
  reordenarItensEdital: (alunoId: string, itens: any[]) => request<any>(`/alunos/${alunoId}/edital/ordem`, { method: 'PUT', body: jsonBody({ itens }) }),
  atualizarProgressoEdital: (alunoId: string, itemId: string, dados: any) => request<any>(`/alunos/${alunoId}/edital/progresso/${itemId}`, { method: 'PATCH', body: jsonBody(dados) }).then(avisarReplanejamento),
  atualizarProgressoMaterial: (alunoId: string, materialId: string, dados: any) => request<any>(`/alunos/${alunoId}/edital/materiais/${encodeURIComponent(materialId)}/progresso`, { method: 'PATCH', body: jsonBody(dados) }),
  criarMaterialApoioMentor: (dados: any) => request<any>('/mentor/materiais-apoio', { method: 'POST', body: dados instanceof FormData ? dados : jsonBody(dados) }),
  alterarMaterialApoioMentor: (id: string, dados: any) => request<any>(`/mentor/materiais-apoio/${id}`, { method: 'PATCH', body: dados instanceof FormData ? dados : jsonBody(dados) }),
  desativarMaterialApoioMentor: (id: string) => request<any>(`/mentor/materiais-apoio/${id}`, { method: 'DELETE' }),
  listarSessoes: (alunoId = 'eu', concursoId = '') => request<any>(`/alunos/${alunoId}/sessoes-estudo${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
  criarSessao: (alunoId: string, dados: any) => request<any>(`/alunos/${alunoId}/sessoes-estudo`, { method: 'POST', body: jsonBody(dados) }),
  alterarSessao: (alunoId: string, id: string, dados: any) => request<any>(`/alunos/${alunoId}/sessoes-estudo/${id}`, { method: 'PATCH', body: jsonBody(dados) }),
  desativarSessao: (alunoId: string, id: string) => request<any>(`/alunos/${alunoId}/sessoes-estudo/${id}`, { method: 'DELETE' }),

  criarQuestaoBanco: (dados: any) => request<any>('/banco-questoes', { method: 'POST', body: jsonBody(dados) }),
  importarQuestoesBanco: (dados: any) => request<any>('/banco-questoes/importar', { method: 'POST', body: jsonBody(dados) }),
  alterarQuestaoBanco: (id: string, dados: any) => request<any>(`/banco-questoes/${id}`, { method: 'PATCH', body: jsonBody(dados) }),
  desativarQuestaoBanco: (id: string) => request<any>(`/banco-questoes/${id}`, { method: 'DELETE' }),
  listarRevisoes: (alunoId = 'eu', concursoId = '') => request<any>(`/alunos/${alunoId}/revisoes${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
  criarRevisao: (alunoId: string, dados: any) => request<any>(`/alunos/${alunoId}/revisoes`, { method: 'POST', body: jsonBody(dados) }).then(avisarReplanejamento),
  alterarRevisao: (alunoId: string, id: string, dados: any) => request<any>(`/alunos/${alunoId}/revisoes/${id}`, { method: 'PATCH', body: jsonBody(dados) }).then(avisarReplanejamento),
  desativarRevisao: (alunoId: string, id: string) => request<any>(`/alunos/${alunoId}/revisoes/${id}`, { method: 'DELETE' }),
  listarSimulados: (alunoId = 'eu', concursoId = '') => request<any>(`/alunos/${alunoId}/simulados${concursoId ? `?concursoId=${encodeURIComponent(concursoId)}` : ''}`),
  criarSimulado: (alunoId: string, dados: any) => request<any>(`/alunos/${alunoId}/simulados`, { method: 'POST', body: jsonBody(dados) }),
  alterarSimulado: (alunoId: string, id: string, dados: any) => request<any>(`/alunos/${alunoId}/simulados/${id}`, { method: 'PUT', body: jsonBody(dados) }),
  desativarSimulado: (alunoId: string, id: string) => request<any>(`/alunos/${alunoId}/simulados/${id}`, { method: 'DELETE' }),
  obterConfiguracaoProva: (alunoId = 'eu', concursoId: string) => request<any>(`/alunos/${alunoId}/configuracao-prova/${encodeURIComponent(concursoId)}`),
  salvarConfiguracaoProva: (alunoId: string, concursoId: string, configuracao: any) => request<any>(`/alunos/${alunoId}/configuracao-prova/${encodeURIComponent(concursoId)}`, { method: 'PUT', body: jsonBody({ configuracao }) }),
  obterResumoMetricas: (alunoId = 'eu', concursoId = '', inicio = '', fim = '') => request<any>(`/alunos/${alunoId}/metricas/resumo${concursoId || inicio || fim ? `?concursoId=${encodeURIComponent(concursoId)}&inicio=${encodeURIComponent(inicio)}&fim=${encodeURIComponent(fim)}` : ''}`),
  obterLinhaDoTempoMetricas: (alunoId = 'eu', concursoId = '', inicio = '', fim = '') => request<any>(`/alunos/${alunoId}/metricas/linha-do-tempo${concursoId || inicio || fim ? `?concursoId=${encodeURIComponent(concursoId)}&inicio=${encodeURIComponent(inicio)}&fim=${encodeURIComponent(fim)}` : ''}`),
  obterMateriasMetricas: (alunoId = 'eu', concursoId = '', inicio = '', fim = '') => request<any>(`/alunos/${alunoId}/metricas/materias${concursoId || inicio || fim ? `?concursoId=${encodeURIComponent(concursoId)}&inicio=${encodeURIComponent(inicio)}&fim=${encodeURIComponent(fim)}` : ''}`),

  listarCadernos: (alunoId = 'eu', editalId = '') => request<any>(`/alunos/${alunoId}/cadernos${editalId ? `?editalId=${encodeURIComponent(editalId)}` : ''}`),
  obterCaderno: (alunoId = 'eu', id: string) => request<any>(`/alunos/${alunoId}/cadernos/${id}`),
  criarCaderno: (alunoId: string, dados: any) => request<any>(`/alunos/${alunoId}/cadernos`, { method: 'POST', body: jsonBody(dados) }),
  alterarCaderno: (alunoId: string, id: string, dados: any) => request<any>(`/alunos/${alunoId}/cadernos/${id}`, { method: 'PATCH', body: jsonBody(dados) }),
  desativarCaderno: (alunoId: string, id: string) => request<any>(`/alunos/${alunoId}/cadernos/${id}`, { method: 'DELETE' }),

  listarUsuarios: () => request<any>('/administracao/usuarios'),
  criarUsuario: (dados: any) => request<any>('/administracao/usuarios', { method: 'POST', body: jsonBody(dados) }),
  alterarUsuario: (id: string, dados: any) => request<any>(`/administracao/usuarios/${id}`, { method: 'PATCH', body: jsonBody(dados) }),
  reativarUsuario: (id: string) => request<any>(`/administracao/usuarios/${id}/reativar`, { method: 'POST', body: '{}' }),
  redefinirSenha: (id: string, novaSenha: string) => request<any>(`/administracao/usuarios/${id}/redefinir-senha`, { method: 'POST', body: jsonBody({ novaSenha }) }),
  definirMentor: (alunoId: string, mentorId: string) => request<any>(`/administracao/alunos/${alunoId}/mentor`, { method: 'PUT', body: jsonBody({ mentorId }) }),

  obterWhiteLabelMentor: () => request<any>('/mentor/whitelabel'),
  salvarWhiteLabelMentor: (dados: any) => request<any>('/mentor/whitelabel', { method: 'POST', body: jsonBody(dados) }),
  uploadLogoWhiteLabel: (file: File) => {
    const form = new FormData();
    form.append('logo', file);
    return request<any>('/mentor/whitelabel/logo', { method: 'POST', body: form });
  },
  uploadBannerWhiteLabel: (file: File) => {
    const form = new FormData();
    form.append('banner', file);
    return request<any>('/mentor/whitelabel/banner', { method: 'POST', body: form });
  },
  obterWhiteLabelMentorPorId: (mentorId: string) => request<any>(`/mentores/${mentorId}/whitelabel`),
  obterWhiteLabelAluno: (alunoId = 'eu') => request<any>(`/alunos/${alunoId}/whitelabel`),
};
