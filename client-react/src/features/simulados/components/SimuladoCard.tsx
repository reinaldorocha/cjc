import React from 'react';
import { Card } from '../../../components/ui/Card';
import { Badge } from '../../../components/ui/Badge';
import { Button } from '../../../components/ui/Button';

interface SimuladoCardProps {
  simulado: any;
  onIniciar: (simulado: any) => void;
}

export const SimuladoCard: React.FC<SimuladoCardProps> = ({ simulado, onIniciar }) => {
  const tituloFormatado =
    simulado.titulo ||
    (simulado.id ? `Simulado Oficial #${simulado.id.slice(0, 6)}` : 'Simulado Oficial');

  return (
    <Card variant="interactive" style={{ marginBottom: '16px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
        <h3 style={{ fontSize: '1.1rem', fontWeight: 600 }}>{tituloFormatado}</h3>
        <Badge variant={simulado.concluido ? 'success' : 'warning'}>
          {simulado.concluido ? 'Concluído' : 'Pendente'}
        </Badge>
      </div>
      <div style={{ fontSize: '0.875rem', color: 'var(--text-secondary)', display: 'flex', gap: '24px', marginBottom: '16px' }}>
        <div>Questões: <strong>{simulado.totalQuestoes || 60}</strong></div>
        <div>Tempo sugerido: <strong>{simulado.tempoMinutos || 240} min</strong></div>
        {simulado.nota !== undefined && (
          <div>Nota: <strong>{simulado.nota}%</strong></div>
        )}
      </div>
      <Button variant="primary" onClick={() => onIniciar(simulado)}>
        {simulado.concluido ? '📊 Ver Relatório' : '📝 Iniciar Simulado'}
      </Button>
    </Card>
  );
};
