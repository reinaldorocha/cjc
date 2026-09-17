import React, { useMemo } from 'react';
import { Card } from '../../../components/ui/Card';
import { Button } from '../../../components/ui/Button';
import type { FiltrosQuestoes } from '../hooks/useBancoQuestoes';

interface FiltroQuestoesBarProps {
  filtros: FiltrosQuestoes;
  materias: any[];
  questoes: any[];
  onChange: (novosFiltros: FiltrosQuestoes) => void;
  onLimpar: () => void;
}

export const FiltroQuestoesBar: React.FC<FiltroQuestoesBarProps> = ({
  filtros,
  materias,
  questoes,
  onChange,
  onLimpar,
}) => {
  // Lista única de disciplinas do edital + questões
  const disciplinasDisponiveis = useMemo(() => {
    const setDisc = new Set<string>();
    (materias || []).forEach((m) => {
      if (m.nome) setDisc.add(m.nome);
    });
    (questoes || []).forEach((q) => {
      if (q.disciplina) setDisc.add(q.disciplina);
    });
    return Array.from(setDisc).sort();
  }, [materias, questoes]);

  // Lista única de assuntos da disciplina selecionada
  const assuntosDisponiveis = useMemo(() => {
    const setAssunto = new Set<string>();

    if (filtros.disciplina) {
      // Buscar tópicos da matéria correspondente
      const mat = (materias || []).find((m) => m.nome?.trim().toLowerCase() === filtros.disciplina.trim().toLowerCase());
      if (mat && Array.isArray(mat.topicos)) {
        mat.topicos.forEach((t: any) => {
          if (t.nome) setAssunto.add(t.nome.trim());
        });
      }
    }

    // Também incluir assuntos existentes nas questões retornadas
    (questoes || []).forEach((q) => {
      if (!filtros.disciplina || q.disciplina?.trim().toLowerCase() === filtros.disciplina.trim().toLowerCase()) {
        if (q.assunto) setAssunto.add(q.assunto.trim());
      }
    });

    return Array.from(setAssunto).sort();
  }, [filtros.disciplina, materias, questoes]);

  return (
    <Card variant="glass" style={{ marginBottom: '24px' }}>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', marginBottom: '16px' }}>
        {[
          { valor: '', rotulo: 'Todas' },
          { valor: 'nao_respondidas', rotulo: 'Não respondidas' },
          { valor: 'errei', rotulo: 'Errei' },
          { valor: 'acertei', rotulo: 'Acertei' },
        ].map((opcao) => {
          const selecionado = filtros.situacao === opcao.valor;
          return (
            <button
              key={opcao.valor || 'todas'}
              type="button"
              aria-pressed={selecionado}
              onClick={() => onChange({ ...filtros, situacao: opcao.valor as FiltrosQuestoes['situacao'] })}
              style={{
                padding: '8px 12px',
                borderRadius: '999px',
                border: selecionado ? '1px solid var(--accent-primary)' : '1px solid var(--border-color)',
                background: selecionado ? 'var(--accent-light)' : 'var(--bg-secondary)',
                color: 'var(--text-primary)',
                cursor: 'pointer',
                fontSize: '0.85rem',
                fontWeight: 600,
              }}
            >
              {opcao.rotulo}
            </button>
          );
        })}
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px', alignItems: 'end' }}>
        {/* Dropdown Disciplina */}
        <div>
          <label style={{ display: 'block', fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '6px', fontWeight: 600 }}>
            📚 Disciplina
          </label>
          <select
            value={filtros.disciplina}
            onChange={(e) => onChange({ ...filtros, disciplina: e.target.value, assunto: '' })}
            style={{
              width: '100%',
              padding: '10px 12px',
              background: 'var(--bg-secondary)',
              border: '1px solid var(--border-color)',
              borderRadius: 'var(--radius-sm)',
              color: '#fff',
              cursor: 'pointer',
            }}
          >
            <option value="">Todas as Disciplinas</option>
            {disciplinasDisponiveis.map((d) => (
              <option key={d} value={d}>
                {d}
              </option>
            ))}
          </select>
        </div>

        {/* Dropdown Assunto */}
        <div>
          <label style={{ display: 'block', fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '6px', fontWeight: 600 }}>
            📌 Assunto / Tópico
          </label>
          <select
            value={filtros.assunto}
            onChange={(e) => onChange({ ...filtros, assunto: e.target.value })}
            disabled={!filtros.disciplina && assuntosDisponiveis.length === 0}
            style={{
              width: '100%',
              padding: '10px 12px',
              background: 'var(--bg-secondary)',
              border: '1px solid var(--border-color)',
              borderRadius: 'var(--radius-sm)',
              color: '#fff',
              cursor: 'pointer',
            }}
          >
            <option value="">
              {!filtros.disciplina ? 'Selecione uma Disciplina primeiro...' : 'Todos os Assuntos'}
            </option>
            {assuntosDisponiveis.map((a) => (
              <option key={a} value={a}>
                {a}
              </option>
            ))}
          </select>
        </div>

        {/* Dropdown Tipo */}
        <div>
          <label style={{ display: 'block', fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '6px', fontWeight: 600 }}>
            ⚙️ Tipo de Questão
          </label>
          <select
            value={filtros.tipo}
            onChange={(e) => onChange({ ...filtros, tipo: e.target.value })}
            style={{
              width: '100%',
              padding: '10px 12px',
              background: 'var(--bg-secondary)',
              border: '1px solid var(--border-color)',
              borderRadius: 'var(--radius-sm)',
              color: '#fff',
              cursor: 'pointer',
            }}
          >
            <option value="">Todos os Tipos</option>
            <option value="multipla_escolha">Múltipla Escolha</option>
            <option value="certo_errado">Certo / Errado</option>
          </select>
        </div>

        {/* Botão Limpar */}
        <div>
          <Button variant="secondary" onClick={onLimpar} style={{ width: '100%' }}>
            🔄 Limpar Filtros
          </Button>
        </div>
      </div>
    </Card>
  );
};
