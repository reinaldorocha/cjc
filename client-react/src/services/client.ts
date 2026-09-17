const API_BASE = '/api/v1';

export class ErroApi extends Error {
  status: number;
  codigo: string;
  constructor(status: number, codigo: string, mensagem: string) {
    super(mensagem);
    this.status = status;
    this.codigo = codigo;
  }
}

let tokenCsrf = '';
let renovacaoPendente: Promise<any> | null = null;

const eMetodoEscrita = (metodo = 'GET') => !['GET', 'HEAD', 'OPTIONS'].includes(metodo.toUpperCase());

export async function request<T>(path: string, init: RequestInit = {}, repeticao = false): Promise<T> {
  const metodo = init.method || 'GET';
  const headers = new Headers(init.headers || {});
  
  if (eMetodoEscrita(metodo) && tokenCsrf) {
    headers.set('X-Token-CSRF', tokenCsrf);
  }
  if (!(init.body instanceof FormData) && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    credentials: 'include',
    headers,
  });

  const corpo = await response.json().catch(() => ({}));

  if (!response.ok) {
    const erro = corpo?.erro || {};
    if (!repeticao && eMetodoEscrita(metodo) && erro.codigo === 'CSRF_INVALIDO') {
      await renovarSessao();
      return request<T>(path, init, true);
    }
    throw new ErroApi(
      response.status,
      erro.codigo || 'ERRO_HTTP',
      erro.mensagem || `Falha na comunicação com o servidor (${response.status}).`
    );
  }

  if (corpo && typeof corpo === 'object' && 'tokenCsrf' in corpo && corpo.tokenCsrf) {
    tokenCsrf = corpo.tokenCsrf as string;
  }

  return (corpo?.dados !== undefined ? corpo.dados : corpo) as T;
}

export function setTokenCsrf(token: string) {
  tokenCsrf = token;
}

export function getTokenCsrf() {
  return tokenCsrf;
}

export function renovarSessao(): Promise<any> {
  if (renovacaoPendente) return renovacaoPendente;
  renovacaoPendente = request<any>('/autenticacao/eu', {}, true)
    .then((dados) => {
      if (dados?.tokenCsrf) tokenCsrf = dados.tokenCsrf;
      return dados;
    })
    .finally(() => {
      renovacaoPendente = null;
    });
  return renovacaoPendente;
}

export const jsonBody = (valor: unknown) => JSON.stringify(valor);
