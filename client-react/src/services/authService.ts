import { request, jsonBody, renovarSessao, setTokenCsrf } from './client';

export const authService = {
  entrar: async (email: string, senha: string) => {
    const dados = await request<any>('/autenticacao/entrar', {
      method: 'POST',
      body: jsonBody({ email, senha }),
    });
    if (dados?.tokenCsrf) setTokenCsrf(dados.tokenCsrf);
    return dados;
  },
  obterSessao: renovarSessao,
  sair: async () => {
    try {
      return await request<any>('/autenticacao/sair', { method: 'POST', body: '{}' });
    } finally {
      setTokenCsrf('');
    }
  },
  sairDeTodas: async () => {
    try {
      return await request<any>('/autenticacao/sair-de-todas', { method: 'POST', body: '{}' });
    } finally {
      setTokenCsrf('');
    }
  },
  alterarNome: (nome: string) =>
    request<any>('/autenticacao/perfil', { method: 'PATCH', body: jsonBody({ nome }) }),
  alterarSenha: (senhaAtual: string, novaSenha: string) =>
    request<any>('/autenticacao/alterar-senha', {
      method: 'POST',
      body: jsonBody({ senhaAtual, novaSenha }),
    }),
};
