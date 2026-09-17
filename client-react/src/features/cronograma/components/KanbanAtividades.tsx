import React from 'react';
import { Card } from '../../../components/ui/Card';
import { Button } from '../../../components/ui/Button';

interface KanbanAtividadesProps {
  atividades: any[];
  onConcluir: (item: any) => void;
}

export const KanbanAtividades: React.FC<KanbanAtividadesProps> = ({ atividades, onConcluir }) => {
  return (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '20px', marginTop: '24px' }}>
      <Card variant="glass">
        <h3 style={{ fontSize: '1.1rem', fontWeight: 600, marginBottom: '16px', borderBottom: '1px solid var(--border-color)', paddingBottom: '8px' }}>
          📋 Atividades de Hoje ({atividades.length})
        </h3>
        {atividades.length === 0 ? (
          <p style={{ color: 'var(--text-muted)', fontSize: '0.875rem' }}>Nenhuma atividade pendente para hoje.</p>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {atividades.map((item, idx) => (
              <div
                key={item.id || idx}
                style={{
                  background: 'var(--bg-secondary)',
                  padding: '12px 16px',
                  borderRadius: 'var(--radius-md)',
                  border: '1px solid var(--border-color)',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                }}
              >
                <div>
                  <div style={{ fontWeight: 600, fontSize: '0.9rem' }}>{item.materiaNome || 'Estudo'}</div>
                  <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>{item.topicoNome || 'Revisão'}</div>
                </div>
                <Button size="sm" variant="primary" onClick={() => onConcluir(item)}>
                  Concluir
                </Button>
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
};
