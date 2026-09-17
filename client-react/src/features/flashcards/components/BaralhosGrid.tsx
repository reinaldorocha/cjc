import React from 'react';
import { Card } from '../../../components/ui/Card';
import { Badge } from '../../../components/ui/Badge';
import { Button } from '../../../components/ui/Button';

interface BaralhosGridProps {
  baralhos: any[];
  onSelecionarBaralho: (baralho: any) => void;
}

export const BaralhosGrid: React.FC<BaralhosGridProps> = ({ baralhos, onSelecionarBaralho }) => {
  if (baralhos.length === 0) {
    return (
      <Card style={{ padding: '32px', textAlign: 'center' }}>
        <p style={{ color: 'var(--text-muted)' }}>Nenhum baralho cadastrado ainda.</p>
      </Card>
    );
  }

  return (
    <div className="ui-grid-auto">
      {baralhos.map((b) => {
        const cartoes = b.cartoes || [];
        const total = b.totalCartoes ?? cartoes.length ?? 0;
        const dominados = cartoes.filter((c: any) => Boolean(c.ultimaRevisao) && (c.repeticoes || 0) > 0).length;
        const pct = total > 0 ? Math.round((dominados / total) * 100) : 0;

        return (
          <Card key={b.id} variant="interactive">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '8px' }}>
              <div>
                <h3 style={{ fontSize: '1.1rem', fontWeight: 700, margin: 0 }}>
                  {b.icone ? `${b.icone} ` : ''}{b.titulo || b.nome}
                </h3>
                <small style={{ color: 'var(--text3)', fontSize: '12px' }}>
                  {b.alcance === 'global' ? '🌐 Todos os alunos' : '🎯 Concurso Específico'}
                </small>
              </div>
              <Badge variant={pct === 100 ? 'success' : pct > 0 ? 'info' : 'neutral'}>
                {total} {total === 1 ? 'cartão' : 'cartões'}
              </Badge>
            </div>

            <p style={{ fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '14px', minHeight: '36px' }}>
              {b.descricao || 'Baralho de fixação de tópicos e legislação.'}
            </p>

            {/* BARRA DE DOMÍNIO / PROGRESSO */}
            <div style={{ marginBottom: '16px', background: 'var(--bg-input)', padding: '10px 12px', borderRadius: '8px', border: '1px solid var(--border)' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '6px', fontSize: '12px' }}>
                <span style={{ fontWeight: 700, color: pct > 0 ? 'var(--green)' : 'var(--text2)' }}>
                  {pct > 0 ? `🟢 ${pct}% Dominado` : '⚪ 0% Estudado'}
                </span>
                <span style={{ color: 'var(--text3)', fontWeight: 600 }}>
                  {dominados}/{total} aprendidos
                </span>
              </div>
              <div style={{ height: '8px', background: 'var(--bg-card)', borderRadius: '4px', overflow: 'hidden' }}>
                <div
                  style={{
                    height: '100%',
                    width: `${pct}%`,
                    background: pct === 100 ? 'var(--green)' : 'linear-gradient(90deg, var(--accent) 0%, var(--green) 100%)',
                    borderRadius: '4px',
                    transition: 'width 0.3s ease'
                  }}
                />
              </div>
            </div>

            <Button variant="primary" onClick={() => onSelecionarBaralho(b)} style={{ width: '100%' }}>
              🎴 Estudar Baralho
            </Button>
          </Card>
        );
      })}
    </div>
  );
};
