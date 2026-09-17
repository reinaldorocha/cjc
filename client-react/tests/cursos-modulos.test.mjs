import test from 'node:test';
import assert from 'node:assert/strict';
import { agruparModulos, aulasDosModulos, primeiraAulaDoModulo, resumoModulo, validarModulo } from '../src/features/cursos/modulos.ts';

test('módulo abre pela primeira aula disponível e resume bloqueio e progresso', () => {
  const aulas = [
    { id: 'a1', moduloId: 'm1', modulo: 'Módulo', titulo: '1', videoId: 'dQw4w9WgXcQ', bloqueada: false, progresso: { concluida: true } },
    { id: 'a2', moduloId: 'm1', modulo: 'Módulo', titulo: '2', videoId: 'dQw4w9WgXcQ', bloqueada: false },
  ];
  assert.equal(primeiraAulaDoModulo(aulas, 'm1'), 'a2');
  assert.deepEqual(resumoModulo(aulas), { total: 2, concluidas: 1, percentual: 50, bloqueado: false, motivoBloqueio: '' });
  const bloqueadas = aulas.map((a) => ({ ...a, bloqueada: true, motivoBloqueio: 'Disponível amanhã' }));
  assert.equal(primeiraAulaDoModulo(bloqueadas, 'm1'), undefined);
  assert.deepEqual(resumoModulo(bloqueadas), { total: 2, concluidas: 1, percentual: 50, bloqueado: true, motivoBloqueio: 'Disponível amanhã' });
});

test('remover capa e descrição não restaura os valores antigos das aulas', () => {
  const modulos = agruparModulos([{ id: 'a', moduloId: 'm', modulo: 'M', titulo: 'A', videoId: 'dQw4w9WgXcQ', moduloCapaUrl: '/api/v1/cursos/capas/id', moduloDescricao: 'Anterior' }]);
  modulos[0].capaUrl = ''; modulos[0].descricao = '';
  const [aula] = aulasDosModulos(modulos);
  assert.equal(aula.moduloCapaUrl, '');
  assert.equal(aula.moduloDescricao, '');
});

test('renomear e reordenar preserva identidade, questões e regras dos módulos', () => {
  const aulas = [
    { id: 'a', moduloId: 'm1', modulo: 'Inicial', titulo: 'A', videoId: 'dQw4w9WgXcQ', questoes: ['q1'] },
    { id: 'b', moduloId: 'm2', modulo: 'Avançado', titulo: 'B', videoId: 'dQw4w9WgXcQ', liberarEm: '2026-09-11', exigeAnterior: true, pdfId: 'pdf' },
  ];
  const modulos = agruparModulos(aulas);
  modulos[1].titulo = 'Novo título';
  modulos.reverse();
  const resultado = aulasDosModulos(modulos);
  assert.equal(resultado[0].id, 'b');
  assert.equal(resultado[0].moduloId, 'm2');
  assert.equal(resultado[0].modulo, 'Novo título');
  assert.equal(resultado[0].liberarEm, '2026-09-11');
  assert.equal(resultado[0].exigeAnterior, true);
  assert.equal(resultado[0].pdfId, 'pdf');
  assert.deepEqual(resultado[1].questoes, ['q1']);
});

test('100 aulas existentes viram 10 módulos sem perder PDFs ou ordem', () => {
  const aulas = Array.from({ length: 100 }, (_, i) => ({ titulo: `Aula ${i + 1}`, modulo: `Módulo ${Math.floor(i / 10) + 1}`, videoId: 'dQw4w9WgXcQ', pdfId: `pdf-${i}`, pdfNome: `${i}.pdf` }));
  const modulos = agruparModulos(aulas);
  assert.equal(modulos.length, 10);
  assert.equal(modulos[0].aulas.length, 10);
  assert.deepEqual(aulasDosModulos(modulos), aulas);
  modulos[0].titulo = 'Fundamentos';
  const resultado = aulasDosModulos(modulos);
  assert.equal(resultado[9].modulo, 'Fundamentos');
  assert.equal(resultado[9].pdfId, 'pdf-9');
  assert.equal(aulas[0].modulo, 'Módulo 1');
});

test('módulo exige título único e aulas preenchidas', () => {
  const modulo = { titulo: 'Direito', aulas: [{ titulo: 'Aula 1', modulo: 'Direito', videoId: 'dQw4w9WgXcQ' }] };
  assert.equal(validarModulo(modulo, []), '');
  assert.notEqual(validarModulo(modulo, [{ ...modulo, titulo: 'direito' }]), '');
  assert.notEqual(validarModulo({ ...modulo, aulas: [] }, []), '');
  assert.notEqual(validarModulo({ ...modulo, aulas: [{ titulo: '', videoId: '' }] }, []), '');
});
