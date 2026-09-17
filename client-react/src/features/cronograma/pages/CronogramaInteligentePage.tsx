import React, { useState } from 'react';
import { useCronograma } from '../hooks/useCronograma';
import { CronogramaHeader } from '../components/CronogramaHeader';
import { CicloEstudoGrid } from '../components/CicloEstudoGrid';
import { KanbanAtividades } from '../components/KanbanAtividades';
import { ModalConclusaoAtividade } from '../../../components/ModalConclusaoAtividade';
import { ModalConfiguracaoCronograma } from '../components/ModalConfiguracaoCronograma';
import { Spinner } from '../../../components/ui/Spinner';
import { Card } from '../../../components/ui/Card';
import { useData } from '../../../context/DataContext';

export const CronogramaInteligentePage: React.FC = () => {
  const { config, cronograma, loading, erro, carregar, reprogramar, gerarCronograma } = useCronograma();
  const { alunoId, activeContestId, getArray } = useData();
  const [selectedItem, setSelectedItem] = useState<any>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isConfigModalOpen, setIsConfigModalOpen] = useState(false);

  const materiasData = getArray('materias') || [];

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '80px' }}>
        <Spinner size="lg" />
      </div>
    );
  }

  if (erro) {
    return (
      <Card style={{ padding: '32px', textAlign: 'center', margin: '40px auto', maxWidth: '500px' }}>
        <p style={{ color: 'var(--status-danger)', marginBottom: '16px' }}>{erro}</p>
        <button className="ui-btn ui-btn-primary" onClick={carregar}>Tentar Novamente</button>
      </Card>
    );
  }

  const totalHoras = Object.values(config.horas || {}).reduce((acc, h) => acc + (Number(h) || 0), 0);
  const atividades = cronograma?.atividades || cronograma?.itens || [];
  const materias = cronograma?.materias || [];
  const temCronograma = Boolean(cronograma && (atividades.length > 0 || materias.length > 0));

  return (
    <div className="ui-container" style={{ padding: '24px 0' }}>
      <CronogramaHeader
        tipo={config.tipo || 'ciclo_inteligente'}
        totalHorasSemanais={totalHoras}
        onConfigurar={() => setIsConfigModalOpen(true)}
        onReprogramar={reprogramar}
      />

      {!temCronograma ? (
        <Card style={{ padding: '48px', textAlign: 'center', marginTop: '24px' }}>
          <h2>⚡ Nenhum Cronograma Gerado Ainda</h2>
          <p style={{ margin: '12px 0 24px 0', color: 'var(--text2)' }}>
            Configure a disponibilidade de estudos semanal e as matérias do edital para gerar o seu cronograma personalizado.
          </p>
          <button
            className="btn-primary"
            style={{ padding: '12px 24px', fontSize: '15px' }}
            onClick={() => setIsConfigModalOpen(true)}
          >
            ⚙️ Configurar e Gerar Cronograma Agora
          </button>
        </Card>
      ) : (
        <>
          <h2 style={{ fontSize: '1.25rem', fontWeight: 600, margin: '24px 0 16px', color: 'var(--text-primary)' }}>
            Ciclo de Estudos Otimizado
          </h2>
          <CicloEstudoGrid materias={materias} />
          <KanbanAtividades
            atividades={atividades}
            onConcluir={(item) => {
              setSelectedItem(item);
              setIsModalOpen(true);
            }}
          />
        </>
      )}

      {isModalOpen && selectedItem && (
        <ModalConclusaoAtividade
          item={selectedItem}
          alunoId={alunoId}
          activeContestId={activeContestId}
          materias={materiasData}
          aoFechar={() => setIsModalOpen(false)}
          aoConcluirSucesso={() => carregar()}
        />
      )}

      <ModalConfiguracaoCronograma
        isOpen={isConfigModalOpen}
        onClose={() => setIsConfigModalOpen(false)}
        configInicial={config}
        onGerar={async (novaConfig) => {
          await gerarCronograma(novaConfig);
        }}
      />
    </div>
  );
};
