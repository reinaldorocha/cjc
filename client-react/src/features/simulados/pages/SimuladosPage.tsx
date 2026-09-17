import React from 'react';
import { useSimulados } from '../hooks/useSimulados';
import { SimuladoCard } from '../components/SimuladoCard';
import { Spinner } from '../../../components/ui/Spinner';
import { EmptyState } from '../../../components/ui/EmptyState';

export const SimuladosPage: React.FC = () => {
  const { simulados, loading, erro } = useSimulados();

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '80px' }}>
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="ui-container" style={{ padding: '24px 0' }}>
      <div className="ui-page-header">
        <h1 className="ui-page-title">Simulados & Provas</h1>
        <p className="ui-page-subtitle">Testes práticos no formato real da prova do seu concurso.</p>
      </div>

      {erro ? (
        <div style={{ color: 'var(--status-danger)', textAlign: 'center', padding: '20px' }}>{erro}</div>
      ) : simulados.length === 0 ? (
        <EmptyState
          title="Nenhum simulado disponível"
          description="Você ainda não possui simulados cadastrados para este concurso."
        />
      ) : (
        <div>
          {simulados.map((s) => (
            <SimuladoCard key={s.id} simulado={s} onIniciar={(sim) => alert(`Iniciar simulado: ${sim.titulo}`)} />
          ))}
        </div>
      )}
    </div>
  );
};
