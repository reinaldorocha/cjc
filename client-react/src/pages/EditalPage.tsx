import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { useData } from '../context/DataContext';
import { api } from '../services/api';
import { Card } from '../components/ui/Card';
import { Badge } from '../components/ui/Badge';
import { Button } from '../components/ui/Button';
import { Spinner } from '../components/ui/Spinner';

const hoje = () => new Date().toISOString().slice(0, 10);
const dataHora = (data: string) => `${data}T12:00:00Z`;

const formatarTempo = (segundos: number) => {
  if (!segundos || segundos <= 0) return '0min';
  const h = Math.floor(segundos / 3600);
  const m = Math.floor((segundos % 3600) / 60);
  if (h > 0) return `${h}h ${m > 0 ? `${m}m` : ''}`;
  return `${m}min`;
};

export const EditalPage: React.FC = () => {
  const { alunoId, activeContestId, activeContest, editalAtivo, recarregarEdital } = useData();
  const [materiaId, setMateriaId] = useState('');
  const [busca, setBusca] = useState('');
  const [loading, setLoading] = useState(false);
  const [erro, setErro] = useState('');

  // Histórico de estudo e questões
  const [sessoes, setSessoes] = useState<any[]>([]);
  const [questoes, setQuestoes] = useState<any[]>([]);

  // Modais de ação rápida
  const [modalRegistro, setModalRegistro] = useState<any | null>(null);
  const [registroForm, setRegistroForm] = useState({
    horas: '0',
    minutos: '30',
    acertos: '0',
    erros: '0',
    data: hoje(),
  });

  const materias = useMemo(() => editalAtivo?.materias || [], [editalAtivo]);
  const materiaSelecionada = materias.find((m: any) => m.id === materiaId) || materias[0];

  useEffect(() => {
    if (materiaSelecionada && materiaSelecionada.id !== materiaId) {
      setMateriaId(materiaSelecionada.id);
    }
  }, [materiaSelecionada, materiaId]);

  useEffect(() => {
    setMateriaId('');
    setErro('');
  }, [activeContestId]);

  const carregarHistorico = useCallback(async () => {
    if (!activeContestId) return;
    setLoading(true);
    try {
      const [resSessoes, resQuestoes] = await Promise.all([
        api.listarSessoes(alunoId, activeContestId).catch(() => ({ sessoes: [] })),
        api.listarQuestoes(alunoId, activeContestId).catch(() => ({ registros: [] })),
      ]);
      setSessoes(resSessoes.sessoes || []);
      setQuestoes(resQuestoes.registros || []);
    } catch {
      // Ignorar erros secundários de histórico
    } finally {
      setLoading(false);
    }
  }, [activeContestId, alunoId]);

  useEffect(() => {
    carregarHistorico();
  }, [carregarHistorico]);

  // Estatísticas gerais do edital
  const itensFinais = useMemo(() => {
    return materias.flatMap((m: any) =>
      (m.topicos || []).flatMap((t: any) => (t.subtopicos?.length ? t.subtopicos : [t]))
    );
  }, [materias]);

  const concluidos = itensFinais.filter((x: any) => x.estudado).length;
  const percentual = itensFinais.length ? Math.round((concluidos / itensFinais.length) * 100) : 0;
  const segundosTotal = sessoes.reduce((n: number, x: any) => n + Number(x.segundos || 0), 0);
  const questoesTotal = questoes.reduce((n: number, x: any) => n + Number(x.resolvidas || 0), 0);
  const acertosTotal = questoes.reduce((n: number, x: any) => n + Number(x.acertos || 0), 0);
  const taxaAcertosGlobal = questoesTotal ? Math.round((acertosTotal / questoesTotal) * 100) : 0;

  // Mapa de estatísticas por tópico/subtópico
  const getEstatisticasItem = (topicoId: string, subtopicoId?: string) => {
    const sessoesFiltradas = sessoes.filter((s: any) =>
      subtopicoId ? s.subtopicoId === subtopicoId : s.topicoId === topicoId
    );
    const segundos = sessoesFiltradas.reduce((acc: number, s: any) => acc + (Number(s.segundos) || 0), 0);

    const questoesFiltradas = questoes.filter((q: any) =>
      subtopicoId ? q.subtopicoId === subtopicoId : q.topicoId === topicoId
    );
    const res = questoesFiltradas.reduce((acc: number, q: any) => acc + (Number(q.resolvidas) || 0), 0);
    const ac = questoesFiltradas.reduce((acc: number, q: any) => acc + (Number(q.acertos) || 0), 0);

    return { segundos, resolvidas: res, acertos: ac };
  };

  // Alterar progresso (Concluído / Pendente)
  const alternarProgresso = async (item: any, tipo: 'topico' | 'subtopico') => {
    try {
      await api.atualizarProgressoEdital(alunoId, item.id, {
        tipo,
        estudado: !item.estudado,
      });
      await recarregarEdital();
    } catch (e: any) {
      setErro(e?.message || 'Erro ao atualizar status.');
    }
  };

  // Abrir Cronômetro
  const abrirTimer = (item: any, tipo: 'topico' | 'subtopico') => {
    window.dispatchEvent(
      new CustomEvent('ct:open-timer', {
        detail: {
          concursoId: activeContestId,
          materiaId: materiaSelecionada?.id,
          topicoId: tipo === 'topico' ? item.id : item.topicoId || '',
          subtopicoId: tipo === 'subtopico' ? item.id : '',
        },
      })
    );
  };

  // Salvar Lançamento Manual de Estudo / Questões
  const handleSalvarRegistro = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!modalRegistro) return;
    const segundos = (Number(registroForm.horas) || 0) * 3600 + (Number(registroForm.minutos) || 0) * 60;
    const acertos = Number(registroForm.acertos) || 0;
    const erros = Number(registroForm.erros) || 0;
    const resolvidas = acertos + erros;

    const basePayload = {
      concursoId: activeContestId,
      materiaId: materiaSelecionada?.id,
      topicoId: modalRegistro.tipo === 'topico' ? modalRegistro.item.id : modalRegistro.item.topicoId || null,
      subtopicoId: modalRegistro.tipo === 'subtopico' ? modalRegistro.item.id : null,
      origem: 'edital',
    };

    try {
      if (segundos > 0) {
        await api.criarSessao(alunoId, {
          ...basePayload,
          segundos,
          modo: 'manual',
          estudadoEm: dataHora(registroForm.data),
        });
      }
      if (resolvidas > 0) {
        await api.criarQuestoes(alunoId, {
          ...basePayload,
          resolvidas,
          acertos,
          erros,
          registradoEm: dataHora(registroForm.data),
        });
      }
      setModalRegistro(null);
      setRegistroForm({ horas: '0', minutos: '30', acertos: '0', erros: '0', data: hoje() });
      await Promise.all([recarregarEdital(), carregarHistorico()]);
    } catch (err: any) {
      alert(err?.message || 'Erro ao salvar registro de estudo.');
    }
  };

  if (!editalAtivo) {
    return (
      <div className="ui-container" style={{ padding: '40px 0', textAlign: 'center' }}>
        <h1 className="ui-page-title">Edital Verticalizado</h1>
        <p style={{ color: 'var(--text-muted)', marginTop: '8px' }}>
          Nenhum edital cadastrado para o concurso selecionado.
        </p>
      </div>
    );
  }

  // Filtrar tópicos da matéria atual
  const topicosFiltrados = (materiaSelecionada?.topicos || []).filter((t: any) => {
    if (!busca) return true;
    const texto = `${t.nome} ${(t.subtopicos || []).map((s: any) => s.nome).join(' ')}`.toLowerCase();
    return texto.includes(busca.toLowerCase());
  });

  return (
    <div className="ui-container" style={{ padding: '24px 0' }}>
      {/* HEADER DO EDITAL */}
      <div className="ui-page-header" style={{ marginBottom: '24px' }}>
        <div>
          <h1 className="ui-page-title">Edital Verticalizado</h1>
          <p className="ui-page-subtitle">
            {editalAtivo.nome || 'Edital Oficial'} • {activeContest?.nome || 'Concurso'}
          </p>
        </div>
      </div>

      {erro && (
        <Card style={{ marginBottom: '16px', borderLeft: '4px solid var(--status-danger)' }}>
          <p style={{ color: 'var(--status-danger)' }}>{erro}</p>
        </Card>
      )}

      {/* STATS DE COBERTURA */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '16px', marginBottom: '24px' }}>
        <Card variant="glass">
          <div style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>Progresso Geral</div>
          <div style={{ fontSize: '1.6rem', fontWeight: 700, color: 'var(--accent-primary)', marginTop: '4px' }}>
            {percentual}%
          </div>
          <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{concluidos} de {itensFinais.length} tópicos</div>
        </Card>

        <Card variant="default">
          <div style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>Tempo Estudado</div>
          <div style={{ fontSize: '1.6rem', fontWeight: 700, color: 'var(--status-info)', marginTop: '4px' }}>
            {formatarTempo(segundosTotal)}
          </div>
          <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{sessoes.length} sessões registradas</div>
        </Card>

        <Card variant="default">
          <div style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>Questões Resolvidas</div>
          <div style={{ fontSize: '1.6rem', fontWeight: 700, color: '#a855f7', marginTop: '4px' }}>
            {questoesTotal}
          </div>
          <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{taxaAcertosGlobal}% de aproveitamento</div>
        </Card>
      </div>

      {/* BARRA DE BUSCA SIMPLES */}
      <Card variant="glass" style={{ marginBottom: '24px' }}>
        <input
          type="text"
          placeholder="🔍 Buscar tópico ou subtópico..."
          value={busca}
          onChange={(e) => setBusca(e.target.value)}
          style={{
            width: '100%',
            padding: '12px 16px',
            background: 'var(--bg-secondary)',
            border: '1px solid var(--border-color)',
            borderRadius: 'var(--radius-sm)',
            color: '#fff',
            fontSize: '0.95rem',
          }}
        />
      </Card>

      {/* LAYOUT PRINCIPAL: MATÉRIAS + CONTEÚDO */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '24px' }}>
        {/* NAVEGAÇÃO LATERAL DE MATÉRIAS */}
        <div>
          <Card variant="default" style={{ padding: '16px' }}>
            <h3 style={{ fontSize: '0.95rem', fontWeight: 700, color: 'var(--text-secondary)', marginBottom: '12px', textTransform: 'uppercase', letterSpacing: '0.5px' }}>
              📚 Matérias ({materias.length})
            </h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              {materias.map((m: any) => {
                const isActive = m.id === materiaSelecionada?.id;
                const totalT = m.topicos?.length || 0;
                const concT = (m.topicos || []).filter((t: any) => t.estudado).length;

                return (
                  <button
                    key={m.id}
                    onClick={() => setMateriaId(m.id)}
                    style={{
                      textAlign: 'left',
                      padding: '12px 14px',
                      borderRadius: 'var(--radius-sm)',
                      background: isActive ? 'var(--accent-light)' : 'transparent',
                      border: `1px solid ${isActive ? 'var(--accent-primary)' : 'var(--border-color)'}`,
                      color: 'var(--text-primary)',
                      cursor: 'pointer',
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center',
                      transition: 'all 0.15s ease',
                    }}
                  >
                    <div>
                      <div style={{ fontWeight: isActive ? 700 : 500, fontSize: '0.95rem' }}>{m.nome}</div>
                      <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginTop: '2px' }}>
                        {concT} / {totalT} tópicos concluídos
                      </div>
                    </div>
                  </button>
                );
              })}
            </div>
          </Card>
        </div>

        {/* LISTA DE TÓPICOS DA MATÉRIA SELECIONADA */}
        <div style={{ gridColumn: 'span 2' }}>
          {materiaSelecionada ? (
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                <h2 style={{ fontSize: '1.2rem', fontWeight: 700 }}>{materiaSelecionada.nome}</h2>
                <Badge variant="neutral">{topicosFiltrados.length} tópicos</Badge>
              </div>

              {loading ? (
                <div style={{ display: 'flex', justifyContent: 'center', padding: '40px' }}>
                  <Spinner size="md" />
                </div>
              ) : topicosFiltrados.length === 0 ? (
                <Card style={{ textAlign: 'center', padding: '32px' }}>
                  <p style={{ color: 'var(--text-muted)' }}>Nenhum tópico encontrado para o filtro selecionado.</p>
                </Card>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                  {topicosFiltrados.map((t: any) => {
                    const stats = getEstatisticasItem(t.id);
                    const isConcluido = Boolean(t.estudado);

                    return (
                      <Card
                        key={t.id}
                        variant="default"
                        style={{
                          borderLeft: `4px solid ${isConcluido ? 'var(--status-success)' : 'var(--border-color)'}`,
                          padding: '16px',
                        }}
                      >
                        <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', flexWrap: 'wrap', gap: '12px' }}>
                          
                          {/* CHECKBOX + NOME DO TÓPICO */}
                          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', flex: 1, minWidth: '240px' }}>
                            <button
                              onClick={() => alternarProgresso(t, 'topico')}
                              style={{
                                width: '24px',
                                height: '24px',
                                borderRadius: '6px',
                                background: isConcluido ? 'var(--status-success)' : 'var(--bg-secondary)',
                                border: `2px solid ${isConcluido ? 'var(--status-success)' : 'var(--border-color)'}`,
                                color: '#fff',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                cursor: 'pointer',
                                fontSize: '0.85rem',
                                fontWeight: 700,
                                flexShrink: 0,
                              }}
                            >
                              {isConcluido ? '✓' : ''}
                            </button>

                            <div>
                              <div style={{ fontWeight: 600, fontSize: '1rem', color: 'var(--text-primary)', textDecoration: isConcluido ? 'line-through' : 'none', opacity: isConcluido ? 0.8 : 1 }}>
                                {t.nome}
                              </div>

                              {/* METRICAS DO TÓPICO (TEMPO + QUESTÕES) */}
                              <div style={{ display: 'flex', gap: '12px', marginTop: '6px', fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
                                <span>⏱️ <strong>{formatarTempo(stats.segundos)}</strong> estudados</span>
                                <span>•</span>
                                <span>📝 <strong>{stats.resolvidas}</strong> questões {stats.resolvidas > 0 ? `(${Math.round((stats.acertos / stats.resolvidas) * 100)}% acertos)` : ''}</span>
                              </div>
                            </div>
                          </div>

                          {/* BOTÕES DE AÇÃO RÁPIDA */}
                          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                            <Button
                              variant="secondary"
                              onClick={() => setModalRegistro({ item: t, tipo: 'topico' })}
                              style={{ fontSize: '0.8rem', padding: '6px 10px' }}
                            >
                              ＋ Lançar Estudo
                            </Button>

                            <Button
                              variant="secondary"
                              onClick={() => abrirTimer(t, 'topico')}
                              style={{ fontSize: '0.8rem', padding: '6px 10px' }}
                            >
                              ⏱️ Timer
                            </Button>
                          </div>
                        </div>

                        {/* SUBTÓPICOS SE HOUVER */}
                        {t.subtopicos && t.subtopicos.length > 0 && (
                          <div style={{ marginTop: '16px', paddingTop: '12px', borderTop: '1px solid rgba(255,255,255,0.06)', display: 'flex', flexDirection: 'column', gap: '8px' }}>
                            {t.subtopicos.map((s: any) => {
                              const subStats = getEstatisticasItem(t.id, s.id);
                              const subConcluido = Boolean(s.estudado);

                              return (
                                <div
                                  key={s.id}
                                  style={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'space-between',
                                    padding: '8px 12px',
                                    borderRadius: 'var(--radius-sm)',
                                    background: 'var(--bg-secondary)',
                                  }}
                                >
                                  <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                    <button
                                      onClick={() => alternarProgresso(s, 'subtopico')}
                                      style={{
                                        width: '20px',
                                        height: '20px',
                                        borderRadius: '4px',
                                        background: subConcluido ? 'var(--status-success)' : 'transparent',
                                        border: `2px solid ${subConcluido ? 'var(--status-success)' : 'var(--border-color)'}`,
                                        color: '#fff',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        cursor: 'pointer',
                                        fontSize: '0.75rem',
                                      }}
                                    >
                                      {subConcluido ? '✓' : ''}
                                    </button>
                                    <span style={{ fontSize: '0.875rem', textDecoration: subConcluido ? 'line-through' : 'none', opacity: subConcluido ? 0.7 : 1 }}>
                                      {s.nome}
                                    </span>
                                  </div>

                                  <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', display: 'flex', gap: '10px' }}>
                                    <span>⏱️ {formatarTempo(subStats.segundos)}</span>
                                    <span>📝 {subStats.resolvidas} quest.</span>
                                  </div>
                                </div>
                              );
                            })}
                          </div>
                        )}
                      </Card>
                    );
                  })}
                </div>
              )}
            </div>
          ) : (
            <Card style={{ textAlign: 'center', padding: '40px' }}>
              <p style={{ color: 'var(--text-muted)' }}>Selecione uma matéria no menu ao lado.</p>
            </Card>
          )}
        </div>
      </div>

      {/* MODAL LIMPO DE LANÇAMENTO RÁPIDO DE ESTUDO */}
      {modalRegistro && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            background: 'rgba(0,0,0,0.7)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 1000,
            padding: '16px',
          }}
          onClick={() => setModalRegistro(null)}
        >
          <div
            onClick={(e) => e.stopPropagation()}
            style={{
              background: 'var(--bg-primary)',
              border: '1px solid var(--border-color)',
              borderRadius: 'var(--radius-md)',
              width: '100%',
              maxWidth: '460px',
              padding: '24px',
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
              <h3 style={{ fontSize: '1.1rem', fontWeight: 700 }}>
                Registrar Estudo em "{modalRegistro.item.nome}"
              </h3>
              <button
                onClick={() => setModalRegistro(null)}
                style={{ background: 'none', border: 'none', color: '#fff', fontSize: '1.2rem', cursor: 'pointer' }}
              >
                ✕
              </button>
            </div>

            <form onSubmit={handleSalvarRegistro} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <div>
                  <label style={{ display: 'block', fontSize: '0.8rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
                    Horas Estudadas
                  </label>
                  <input
                    type="number"
                    min="0"
                    value={registroForm.horas}
                    onChange={(e) => setRegistroForm({ ...registroForm, horas: e.target.value })}
                    style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
                  />
                </div>
                <div>
                  <label style={{ display: 'block', fontSize: '0.8rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
                    Minutos Estudados
                  </label>
                  <input
                    type="number"
                    min="0"
                    max="59"
                    value={registroForm.minutos}
                    onChange={(e) => setRegistroForm({ ...registroForm, minutos: e.target.value })}
                    style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
                  />
                </div>
              </div>

              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <div>
                  <label style={{ display: 'block', fontSize: '0.8rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
                    Questões Certas
                  </label>
                  <input
                    type="number"
                    min="0"
                    value={registroForm.acertos}
                    onChange={(e) => setRegistroForm({ ...registroForm, acertos: e.target.value })}
                    style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
                  />
                </div>
                <div>
                  <label style={{ display: 'block', fontSize: '0.8rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
                    Questões Erradas
                  </label>
                  <input
                    type="number"
                    min="0"
                    value={registroForm.erros}
                    onChange={(e) => setRegistroForm({ ...registroForm, erros: e.target.value })}
                    style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
                  />
                </div>
              </div>

              <div>
                <label style={{ display: 'block', fontSize: '0.8rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
                  Data do Estudo
                </label>
                <input
                  type="date"
                  value={registroForm.data}
                  onChange={(e) => setRegistroForm({ ...registroForm, data: e.target.value })}
                  style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
                />
              </div>

              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px', marginTop: '8px' }}>
                <Button variant="secondary" onClick={() => setModalRegistro(null)}>
                  Cancelar
                </Button>
                <Button variant="primary" type="submit">
                  Salvar Registro
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
