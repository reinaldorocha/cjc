import React from 'react';
import { Card } from '../../../components/ui/Card';
import { Badge } from '../../../components/ui/Badge';

interface CicloEstudoGridProps {
  materias: any[];
}

export const CicloEstudoGrid: React.FC<CicloEstudoGridProps> = ({ materias }) => {
  if (!materias || materias.length === 0) {
    return (
      <Card style={{ textAlign: 'center', padding: '32px' }}>
        <p style={{ color: 'var(--text-muted)' }}>Nenhuma matéria cadastrada no edital para o ciclo.</p>
      </Card>
    );
  }

  return (
    <div className="ui-grid-auto">
      {materias.map((mat, idx) => (
        <Card key={mat.id || idx} variant="interactive">
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '12px' }}>
            <h4 style={{ fontSize: '1rem', fontWeight: 600, color: 'var(--text-primary)' }}>{mat.nome}</h4>
            <Badge variant={mat.peso > 2 ? 'danger' : 'info'}>Peso {mat.peso || 1}</Badge>
          </div>
          <div style={{ fontSize: '0.875rem', color: 'var(--text-secondary)', display: 'flex', flexDirection: 'column', gap: '6px' }}>
            <div>Tópicos concluídos: <strong>{mat.topicosConcluidos || 0} / {mat.totalTopicos || 0}</strong></div>
            <div>Horas dedicadas: <strong>{mat.horasEstudo || 0}h</strong></div>
          </div>
        </Card>
      ))}
    </div>
  );
};
