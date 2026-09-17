import type { AulaCurso } from '../../services/cursosService';

export type ModuloCurso = { id?: string; titulo: string; descricao?: string; capaUrl?: string; aulas: AulaCurso[]; liberarEm?: string; exigeAnterior?: boolean };

export type ResumoModulo = { total: number; concluidas: number; percentual: number; bloqueado: boolean; motivoBloqueio: string };

export function primeiraAulaDoModulo(aulas: AulaCurso[], moduloId: string): string | undefined {
  const modulo = aulas.filter((a) => a.moduloId === moduloId || a.modulo === moduloId);
  return modulo.find((a) => !a.bloqueada && !a.progresso?.concluida)?.id || modulo.find((a) => !a.bloqueada)?.id;
}

export function resumoModulo(aulas: AulaCurso[]): ResumoModulo {
  const total = aulas.length;
  const concluidas = aulas.filter((a) => a.progresso?.concluida).length;
  const bloqueado = total > 0 && aulas.every((a) => a.bloqueada);
  return { total, concluidas, percentual: total ? Math.round(concluidas * 100 / total) : 0, bloqueado, motivoBloqueio: bloqueado ? aulas.find((a) => a.motivoBloqueio)?.motivoBloqueio || 'Módulo bloqueado' : '' };
}

export function agruparModulos(aulas: AulaCurso[]): ModuloCurso[] {
  const modulos: ModuloCurso[] = [];
  for (const aula of aulas) {
    const titulo = aula.modulo.trim() || 'Aulas';
    let modulo = modulos.find((m) => aula.moduloId ? m.id === aula.moduloId : m.titulo === titulo);
    if (!modulo) {
      modulo = {
        id: aula.moduloId,
        titulo,
        descricao: aula.moduloDescricao || '',
        capaUrl: aula.moduloCapaUrl || '',
        aulas: [],
        liberarEm: aula.liberarEm,
        exigeAnterior: aula.exigeAnterior
      };
      modulos.push(modulo);
    }
    if (!modulo.capaUrl && aula.moduloCapaUrl) modulo.capaUrl = aula.moduloCapaUrl;
    if (!modulo.descricao && aula.moduloDescricao) modulo.descricao = aula.moduloDescricao;
    modulo.aulas.push({ ...aula, modulo: titulo });
  }
  return modulos;
}

export function aulasDosModulos(modulos: ModuloCurso[]): AulaCurso[] {
  return modulos.flatMap((m) => m.aulas.map((a) => ({
    ...a,
    modulo: m.titulo.trim(),
    ...(m.id ? { moduloId: m.id } : {}),
    ...(m.liberarEm !== undefined ? { liberarEm: m.liberarEm } : {}),
    ...(m.exigeAnterior !== undefined ? { exigeAnterior: m.exigeAnterior } : {}),
    ...(m.capaUrl || a.moduloCapaUrl !== undefined ? { moduloCapaUrl: m.capaUrl || '' } : {}),
    ...(m.descricao || a.moduloDescricao !== undefined ? { moduloDescricao: m.descricao || '' } : {})
  })));
}

export function validarModulo(modulo: ModuloCurso, outros: ModuloCurso[]): string {
  const titulo = modulo.titulo.trim();
  if (!titulo || titulo.length > 100) return 'Informe o título do módulo (até 100 caracteres).';
  if (outros.some((m) => m.titulo.toLocaleLowerCase() === titulo.toLocaleLowerCase())) return 'Já existe um módulo com esse título.';
  if (!modulo.aulas.length) return 'Adicione pelo menos uma aula ao módulo.';
  const incompleta = modulo.aulas.findIndex((a) => !a.titulo.trim() || !a.videoId.trim());
  if (incompleta >= 0) return `Preencha o título e o link do YouTube da aula ${incompleta + 1}.`;
  return '';
}
