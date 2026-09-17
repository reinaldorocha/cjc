import React, { useEffect, useState } from 'react';
import { useData } from '../../../context/DataContext';
import { api } from '../../../services/api';
import { Card } from '../../../components/ui/Card';
import { Badge } from '../../../components/ui/Badge';
import { Spinner } from '../../../components/ui/Spinner';
import { useWhiteLabel } from '../../../context/WhiteLabelContext';

interface DiaEstudoSemana {
  key: string;
  label: string;
  estudado: boolean;
  horas: string;
}

export const DashboardPage: React.FC = () => {
  const { activeContest, activeContestId, alunoId, getArray } = useData();
  const { config } = useWhiteLabel();
  const materiasFallback = getArray('materias');
  const [metricas, setMetricas] = useState<any>(null);
  const [linhaTempo, setLinhaTempo] = useState<any[]>([]);
  const [materiasMetricas, setMateriasMetricas] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!activeContestId) {
      setLoading(false);
      return;
    }
    setLoading(true);

    Promise.all([
      api.obterResumoMetricas(alunoId, activeContestId),
      api.obterLinhaDoTempoMetricas(alunoId, activeContestId).catch(() => ({ dias: [] })),
      api.obterMateriasMetricas(alunoId, activeContestId).catch(() => ({ materias: [] })),
    ])
      .then(([resResumo, resLinhaTempo, resMaterias]) => {
        setMetricas(resResumo);
        setLinhaTempo(Array.isArray(resLinhaTempo?.dias) ? resLinhaTempo.dias : []);
        const listaMaterias = Array.isArray(resMaterias)
          ? resMaterias
          : resMaterias?.materias || [];
        setMateriasMetricas(listaMaterias);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [activeContestId, alunoId]);

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '80px' }}>
        <Spinner size="lg" />
      </div>
    );
  }

  // Cálculo de dias até a prova
  let diasParaProva: number | null = null;
  if (activeContest?.dataProva) {
    const diff = new Date(activeContest.dataProva).getTime() - new Date().getTime();
    diasParaProva = Math.ceil(diff / (1000 * 60 * 60 * 24));
  }

  // Cálculo real dos dias da semana atual (Segunda a Domingo)
  const nomesDias = [
    { key: 'seg', label: 'Seg', idxSemana: 1 },
    { key: 'ter', label: 'Ter', idxSemana: 2 },
    { key: 'qua', label: 'Qua', idxSemana: 3 },
    { key: 'qui', label: 'Qui', idxSemana: 4 },
    { key: 'sex', label: 'Sex', idxSemana: 5 },
    { key: 'sab', label: 'Sáb', idxSemana: 6 },
    { key: 'dom', label: 'Dom', idxSemana: 0 },
  ];

  // Mapear datas reais de estudo da semana atual
  const agora = new Date();
  const diaSemanaAtual = agora.getDay(); // 0 = Dom, 1 = Seg ...
  const distSegunda = diaSemanaAtual === 0 ? -6 : 1 - diaSemanaAtual;
  const segundaFeira = new Date(agora);
  segundaFeira.setDate(agora.getDate() + distSegunda);

  const diasSemanaCalculados: DiaEstudoSemana[] = nomesDias.map((d, i) => {
    const dataDia = new Date(segundaFeira);
    dataDia.setDate(segundaFeira.getDate() + i);
    const dataIso = dataDia.toISOString().split('T')[0];

    const registro = linhaTempo.find((item: any) => item.data === dataIso);
    const estudado = Boolean(registro && (registro.segundos > 0 || registro.questoes > 0 || registro.diaAtivo));
    const horasNum = registro?.segundos ? (registro.segundos / 3600).toFixed(1) : '0';

    return {
      key: d.key,
      label: d.label,
      estudado,
      horas: estudado ? `${horasNum}h` : '0h',
    };
  });

  const horasEstudadasFormatada = metricas?.segundosEstudados
    ? (metricas.segundosEstudados / 3600).toFixed(1)
    : metricas?.horasEstudadas || '0';

  const taxaAcertos = metricas?.percentualAcertos !== undefined
    ? Math.round(metricas.percentualAcertos)
    : metricas?.taxaAcertos || 0;

  // Lista final de matérias para renderizar na tabela
  const listaTabelaMaterias = materiasMetricas.length > 0 ? materiasMetricas : materiasFallback;

  return (
    <div className="ui-container" style={{ padding: '24px 0' }}>
      {config.bannerUrl && <img src={config.bannerUrl} alt="Banner da mentoria" style={{ display: 'block', width: '100%', maxHeight: '260px', minHeight: '110px', objectFit: 'cover', borderRadius: 'var(--radius-lg)', marginBottom: '24px', border: '1px solid var(--border-color)' }} />}
      {/* HEADER DO CONCURSO ATIVO */}
      {activeContest && (
        <Card variant="glass" style={{ marginBottom: '24px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '16px' }}>
            <div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <span style={{ fontSize: '24px' }}>🎯</span>
                <h1 className="ui-page-title">{activeContest.nome || 'Concurso Ativo'}</h1>
                <Badge variant="info">{activeContest.banca || 'Banca N/I'}</Badge>
              </div>
              <p className="ui-page-subtitle">
                {activeContest.cargo ? `Cargo: ${activeContest.cargo}` : ''} {activeContest.salario ? `• Salário: R$ ${activeContest.salario}` : ''}
              </p>
            </div>
            {diasParaProva !== null && (
              <div
                style={{
                  textAlign: 'center',
                  background: 'var(--bg-secondary)',
                  padding: '12px 20px',
                  borderRadius: 'var(--radius-md)',
                  border: '1px solid var(--border-color)',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '12px',
                }}
              >
                <span style={{ fontSize: '1.8rem' }}>⏳</span>
                <div style={{ textAlign: 'left' }}>
                  <div
                    style={{
                      fontSize: '1.4rem',
                      fontWeight: 800,
                      color: diasParaProva < 30 ? 'var(--status-warning)' : 'var(--accent-primary)',
                      lineHeight: 1.1,
                    }}
                  >
                    {diasParaProva > 0 ? `${diasParaProva} ${diasParaProva === 1 ? 'dia' : 'dias'}` : 'Hoje!'}
                  </div>
                  <div style={{ fontSize: '0.85rem', fontWeight: 600, color: 'var(--text-secondary)', marginTop: '2px' }}>
                    {diasParaProva > 0 ? 'para a prova' : 'Dia da prova!'}
                  </div>
                </div>
              </div>
            )}
          </div>
        </Card>
      )}

      {/* METRICAS PRINCIPAIS (STAT CARDS) */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px', marginBottom: '28px' }}>
        <Card variant="default">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>Horas Estudadas</span>
            <span style={{ fontSize: '20px' }}>⏱️</span>
          </div>
          <div style={{ fontSize: '1.8rem', fontWeight: 700, color: 'var(--accent-primary)', marginTop: '8px' }}>
            {horasEstudadasFormatada}h
          </div>
        </Card>

        <Card variant="default">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>Questões Resolvidas</span>
            <span style={{ fontSize: '20px' }}>📝</span>
          </div>
          <div style={{ fontSize: '1.8rem', fontWeight: 700, color: '#a855f7', marginTop: '8px' }}>
            {metricas?.questoesResolvidas || 0}
          </div>
        </Card>

        <Card variant="default">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>Taxa de Acertos</span>
            <span style={{ fontSize: '20px' }}>🎯</span>
          </div>
          <div style={{ fontSize: '1.8rem', fontWeight: 700, color: 'var(--status-success)', marginTop: '8px' }}>
            {taxaAcertos}%
          </div>
        </Card>

        <Card variant="default">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>Último Simulado</span>
            <span style={{ fontSize: '20px' }}>📊</span>
          </div>
          <div style={{ fontSize: '1.8rem', fontWeight: 700, color: 'var(--status-warning)', marginTop: '8px' }}>
            {metricas?.ultimoSimulado !== undefined && metricas?.ultimoSimulado !== null ? `${metricas.ultimoSimulado}%` : 'N/I'}
          </div>
        </Card>

      </div>

      {/* SEÇÃO PRINCIPAL: TABELA DE MATÉRIAS + CRONOGRAMA SEMANAL */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))', gap: '24px', marginBottom: '28px' }}>
        
        {/* TABELA DETALHADA DAS MATÉRIAS */}
        <Card variant="glass" style={{ gridColumn: 'span 2' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
            <div>
              <h3 style={{ fontSize: '1.1rem', fontWeight: 600 }}>Matérias do Edital</h3>
              <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>Progresso de cobertura por disciplina</p>
            </div>
          </div>

          {listaTabelaMaterias.length === 0 ? (
            <p style={{ color: 'var(--text-muted)', fontSize: '0.9rem', padding: '16px 0' }}>
              Nenhuma matéria cadastrada no edital deste concurso.
            </p>
          ) : (
            <div style={{ overflowX: 'auto' }}>
              <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.875rem' }}>
                <thead>
                  <tr style={{ borderBottom: '1px solid var(--border-color)', textAlign: 'left', color: 'var(--text-secondary)' }}>
                    <th style={{ padding: '10px' }}>Matéria</th>
                    <th style={{ padding: '10px' }}>Tópicos</th>
                    <th style={{ padding: '10px' }}>Tempo</th>
                    <th style={{ padding: '10px' }}>Acertos</th>
                    <th style={{ padding: '10px' }}>Progresso</th>
                  </tr>
                </thead>
                <tbody>
                  {listaTabelaMaterias.map((mat: any, idx: number) => {
                    const concluidos = mat.itensConcluidos !== undefined ? mat.itensConcluidos : (mat.topicosConcluidos || 0);
                    const total = mat.itensTotal !== undefined ? mat.itensTotal : (mat.totalTopicos || 1);
                    const pct = total > 0 ? Math.min(100, Math.round((concluidos / total) * 100)) : (Math.round(mat.percentualEdital || 0));

                    let horasStr = '0h';
                    if (mat.segundos) {
                      horasStr = `${(mat.segundos / 3600).toFixed(1)}h`;
                    } else if (mat.horasEstudo) {
                      horasStr = `${mat.horasEstudo}h`;
                    }

                    const acertosStr = mat.percentualAcertos != null 
                      ? `${Math.round(mat.percentualAcertos)}%` 
                      : (mat.acertos !== undefined && mat.questoes ? `${Math.round((mat.acertos / mat.questoes) * 100)}%` : '-');

                    return (
                      <tr key={mat.id || idx} style={{ borderBottom: '1px solid rgba(255,255,255,0.04)' }}>
                        <td style={{ padding: '12px 10px', fontWeight: 600 }}>{mat.nome}</td>
                        <td style={{ padding: '12px 10px', color: 'var(--text-secondary)' }}>
                          {concluidos} / {total}
                        </td>
                        <td style={{ padding: '12px 10px', color: 'var(--text-secondary)' }}>
                          {horasStr}
                        </td>
                        <td style={{ padding: '12px 10px', color: 'var(--text-secondary)' }}>
                          {acertosStr}
                        </td>
                        <td style={{ padding: '12px 10px', width: '140px' }}>
                          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                            <div style={{ flex: 1, height: '6px', background: 'var(--bg-secondary)', borderRadius: '3px', overflow: 'hidden' }}>
                              <div style={{ width: `${pct}%`, height: '100%', background: 'var(--accent-primary)', borderRadius: '3px' }} />
                            </div>
                            <span style={{ fontSize: '0.75rem', fontWeight: 600 }}>{pct}%</span>
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </Card>

        {/* PAINEL LATERAL: CRONOGRAMA SEMANAL */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          
          {/* CRONOGRAMA SEMANAL REAL */}
          <Card variant="glass">
            <h3 style={{ fontSize: '1.1rem', fontWeight: 600, marginBottom: '12px' }}>📅 Estudo Semanal</h3>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: '6px', textAlign: 'center' }}>
              {diasSemanaCalculados.map((d) => (
                <div
                  key={d.key}
                  style={{
                    padding: '8px 4px',
                    borderRadius: 'var(--radius-sm)',
                    background: d.estudado ? 'var(--status-success-bg)' : 'var(--bg-secondary)',
                    border: `1px solid ${d.estudado ? 'var(--status-success)' : 'var(--border-color)'}`,
                  }}
                >
                  <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{d.label}</div>
                  <div style={{ fontSize: '0.85rem', fontWeight: 600, color: d.estudado ? 'var(--status-success)' : 'var(--text-secondary)' }}>
                    {d.estudado ? '✓' : '-'}
                  </div>
                </div>
              ))}
            </div>
          </Card>

        </div>
      </div>
    </div>
  );
};
