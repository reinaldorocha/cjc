import React from 'react';
import { Card } from '../../../components/ui/Card';
import { Badge } from '../../../components/ui/Badge';
import { Button } from '../../../components/ui/Button';

interface CronogramaHeaderProps {
  tipo: string;
  totalHorasSemanais: number;
  onConfigurar: () => void;
  onReprogramar: () => void;
}

export const CronogramaHeader: React.FC<CronogramaHeaderProps> = ({
  tipo,
  totalHorasSemanais,
  onConfigurar,
  onReprogramar,
}) => {
  return (
    <Card variant="glass" className="ui-page-header">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <h1 className="ui-page-title">Cronograma Inteligente</h1>
            <Badge variant={tipo === 'ciclo' ? 'info' : 'success'}>
              {tipo === 'ciclo' ? 'Ciclo Adaptativo' : tipo === 'agendado' ? 'Agendado Semanal' : 'Não Configurado'}
            </Badge>
          </div>
          <p className="ui-page-subtitle">
            Plano de estudos otimizado com base no seu ritmo ({totalHorasSemanais}h semanais dedicadas).
          </p>
        </div>
        <div style={{ display: 'flex', gap: '12px' }}>
          <Button variant="secondary" onClick={onReprogramar}>
            🔄 Reprogramar
          </Button>
          <Button variant="primary" onClick={onConfigurar}>
            ⚙️ Ajustar Configurações
          </Button>
        </div>
      </div>
    </Card>
  );
};
