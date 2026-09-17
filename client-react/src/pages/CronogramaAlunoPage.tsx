import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../services/api';
import { useAutenticacao } from '../context/AutenticacaoContext';
import { useData } from '../context/DataContext';
import { ModalConclusaoAtividade } from '../components/ModalConclusaoAtividade';
import { useVisaoAluno } from '../context/VisaoAlunoContext';

const MESES = [
  'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
  'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
];

const dataBr = (v?: string) => (v ? new Date(`${v}T12:00:00`).toLocaleDateString('pt-BR') : 'Sem data');
const rotulo = (v: string) => ({ pendente: 'Pendente', concluido: 'Concluído', ignorado: 'Ignorado' } as any)[v] || v;

export const CronogramaAlunoPage: React.FC = () => {
  const { usuario } = useAutenticacao();
  const visaoAluno = useVisaoAluno();
  const { alunoId, activeContestId, activeContest, getArray } = useData();
  const [cronograma, setCronograma] = useState<any>(undefined);
  const [revisoes, setRevisoes] = useState<any[]>([]);

  const [erro, setErro] = useState('');
  const [ocupado, setOcupado] = useState('');
  const [filtro, setFiltro] = useState('todos');
  const [modoVisao, setModoVisao] = useState<'lista' | 'calendario'>('calendario');
  const [mesOffset, setMesOffset] = useState(0);

  // Estado do Modal de Conclusão Manual (unificado)
  const [concluirModalItem, setConcluirModalItem] = useState<any>(null);
  // Estado do Modal de Detalhes do Dia (ao clicar no calendário)
  const [diaDetalhesModal, setDiaDetalhesModal] = useState<{
    dateKey: string;
    dataFormatada: string;
    dayItens: any[];
  } | null>(null);

  const carregar = useCallback(async () => {
    if (!activeContestId) return;
    try {
      const [resCrono, resRev] = await Promise.all([
        api.obterCronograma(alunoId, activeContestId),
        api.listarRevisoes(alunoId, activeContestId).catch(() => ({ revisoes: [] }))
      ]);
      setCronograma(resCrono.cronograma || null);
      setRevisoes(resRev.revisoes || []);
    } catch (e: any) {
      setErro(e.message);
    }
  }, [alunoId, activeContestId]);

  useEffect(() => {
    carregar();
  }, [carregar]);

  useEffect(() => {
    const handleUpdate = () => {
      carregar();
    };
    window.addEventListener('ct:dados-estudo-alterados', handleUpdate);
    return () => window.removeEventListener('ct:dados-estudo-alterados', handleUpdate);
  }, [carregar]);

  const materias = getArray('materias') || [];
  const topicos = getArray('topicos') || [];
  const subtopicos = getArray('subtopicos') || [];

  const itens = useMemo(() => {
    const lista = cronograma?.itens || [];
    return lista.filter((x: any) => filtro === 'todos' || x.situacao === filtro);
  }, [cronograma, filtro]);

  const dataFechamentoEdital = useMemo(() => {
    const datas = itens
      .map((x: any) => x.dataPlanejada?.slice(0, 10))
      .filter((v: string | undefined): v is string => !!v);
    if (!datas.length) return null;
    return datas.reduce((a: string, b: string) => (a > b ? a : b));
  }, [itens]);

  const grupos = useMemo(() => {
    if (cronograma?.tipo === 'ciclo_inteligente') return [];
    const mapa = new Map<string, any[]>();
    for (const x of itens) {
      const chave = x.dataPlanejada || 'sem-data';
      mapa.set(chave, [...(mapa.get(chave) || []), x]);
    }
    for (const r of revisoes) {
      const chave = (r.proximaData || r.data || '').slice(0, 10) || 'sem-data';
      mapa.set(chave, [...(mapa.get(chave) || []), { ...r, tipoItem: 'revisao' }]);
    }
    return [...mapa.entries()];
  }, [cronograma, itens, revisoes]);

  const calendarData = useMemo(() => {
    const base = new Date();
    base.setDate(1);
    base.setMonth(base.getMonth() + mesOffset);
    const year = base.getFullYear();
    const month = base.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const startOffset = (firstDay.getDay() + 6) % 7; // Seg-Dom
    const days: (Date | null)[] = [];
    for (let i = 0; i < startOffset; i++) days.push(null);
    for (let d = 1; d <= lastDay.getDate(); d++) days.push(new Date(year, month, d));
    return { days, year, month, monthName: MESES[month] };
  }, [mesOffset]);

  const itensPorData = useMemo(() => {
    const mapa = new Map<string, any[]>();
    for (const item of itens) {
      if (item.dataPlanejada) {
        const key = item.dataPlanejada.slice(0, 10);
        mapa.set(key, [...(mapa.get(key) || []), item]);
      }
    }
    for (const r of revisoes) {
      const key = (r.proximaData || r.data || '').slice(0, 10);
      if (key) {
        mapa.set(key, [...(mapa.get(key) || []), { ...r, tipoItem: 'revisao' }]);
      }
    }
    if (activeContest?.dataProva) {
      const key = activeContest.dataProva.slice(0, 10);
      mapa.set(key, [...(mapa.get(key) || []), { tipoItem: 'prova', nome: `🚩 PROVA: ${activeContest.nome}` }]);
    }
    return mapa;
  }, [itens, revisoes, activeContest]);

  const itensSemData = useMemo(() => {
    return itens.filter((x: any) => !x.dataPlanejada);
  }, [itens]);

  const alterar = async (item: any, situacao: string) => {
    setOcupado(item.id);
    try {
      await api.alterarItemCronograma(alunoId, item.id, { situacao });
      if (situacao === 'pendente') {
        const targetId = item.subtopicoId || item.topicoId;
        const targetTipo = item.subtopicoId ? 'subtopico' : 'topico';
        if (targetId) {
          await api.atualizarProgressoEdital(alunoId, targetId, { tipo: targetTipo, estudado: false }).catch(() => undefined);
        }
      }
      window.dispatchEvent(new CustomEvent('ct:dados-estudo-alterados'));
      await carregar();
    } catch (e: any) {
      alert(e.message);
    } finally {
      setOcupado('');
    }
  };

  const abrirModalConclusao = (item: any, e?: React.MouseEvent) => {
    if (e) e.stopPropagation();
    setConcluirModalItem(item);
  };

  const abrirTimer = (item: any, e?: React.MouseEvent) => {
    if (e) e.stopPropagation();
    window.dispatchEvent(
      new CustomEvent('ct:open-timer', {
        detail: {
          materiaId: item.materiaId,
          topicoId: item.topicoId,
          subtopicoId: item.subtopicoId,
          cronogramaItemId: item.id,
          autoStart: true
        }
      })
    );
  };

  const reprogramarPendentes = async () => {
    setOcupado('reprogramar');
    try {
      await api.reprogramarCronograma(alunoId, { concursoId: activeContestId });
      window.dispatchEvent(new CustomEvent('ct:dados-estudo-alterados'));
      await carregar();
    } catch (e: any) {
      alert(e.message || 'Erro ao reprogramar pendentes');
    } finally {
      setOcupado('');
    }
  };

  if (cronograma === undefined) {
    return <div className="empty-state">Carregando cronograma…</div>;
  }

  const todayStr = new Date().toISOString().slice(0, 10);
  const concluidosItens = cronograma?.itens?.filter((x: any) => x.situacao === 'concluido').length || 0;
  const restantesItens = cronograma?.itens?.filter((x: any) => x.situacao === 'pendente').length || 0;
  const totalItens = cronograma?.itens?.length || 0;
  const pctConclusao = totalItens ? Math.round((concluidosItens / totalItens) * 100) : 0;

  return (
    <div className="schedule-page">
      <div className="page-heading">
        <div>
          <h1>Cronograma de Estudos</h1>
          <p>Acompanhe suas tarefas diárias e o avanço no ciclo.</p>
        </div>
        <div style={{ display: 'flex', gap: '8px' }}>
          {usuario?.papel === 'mentor' && !visaoAluno && (
            <Link
              className="btn-secondary"
              to={`/mentor/alunos/${alunoId}/concursos/${activeContestId}/cronograma-inteligente`}
            >
              ⚙ Configurar / Gerar
            </Link>
          )}
          {usuario?.papel === 'aluno' && (
            <Link className="btn-secondary" to="/cronograma-inteligente">
              ⚡ Gerar Cronograma Inteligente
            </Link>
          )}
        </div>
      </div>

      {erro && <div className="feedback-banner error">{erro}</div>}

      {!cronograma ? (
        <div className="empty-state card-base" style={{ padding: '36px', textAlign: 'center' }}>
          <h2>Nenhum cronograma ativo para este concurso</h2>
          <p style={{ margin: '12px 0 24px 0', color: 'var(--text2)' }}>
            Gere o seu cronograma personalizado para organizar seus estudos por dia ou por ciclo.
          </p>

          {usuario?.papel === 'mentor' && !visaoAluno ? (
            <Link
              className="btn-primary"
              to={`/mentor/alunos/${alunoId}/concursos/${activeContestId}/cronograma-inteligente`}
            >
              ⚡ Gerar Cronograma Agora
            </Link>
          ) : (
            <Link className="btn-primary" to="/cronograma-inteligente">
              ⚡ Gerar Meu Cronograma
            </Link>
          )}
        </div>
      ) : (
        <>
          <div className="stat-grid" style={{ marginBottom: '24px' }}>
            <div className="stat-card">
              <span>Data fechamento edital</span>
              <strong>{dataFechamentoEdital ? dataBr(dataFechamentoEdital) : 'Sem data'}</strong>
              <small>último dia planejado do cronograma</small>
            </div>
            <div className="stat-card">
              <span>Concluídas</span>
              <strong style={{ color: 'var(--green)' }}>{concluidosItens}</strong>
              <small>tópicos finalizados</small>
            </div>
            <div className="stat-card">
              <span>Restantes</span>
              <strong style={{ color: restantesItens > 0 ? 'var(--yellow)' : 'var(--green)' }}>{restantesItens}</strong>
              <small>pendentes de estudo</small>
            </div>
            <div className="stat-card">
              <span>Progresso geral</span>
              <strong style={{ color: 'var(--accent)' }}>{pctConclusao}%</strong>
              <small>concluído</small>
            </div>
          </div>

          <div className="schedule-toolbar card-base" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <div className="period-chips">
              <button className={filtro === 'todos' ? 'active' : ''} onClick={() => setFiltro('todos')}>
                Todas
              </button>
              <button className={filtro === 'pendente' ? 'active' : ''} onClick={() => setFiltro('pendente')}>
                Pendentes
              </button>
              <button className={filtro === 'concluido' ? 'active' : ''} onClick={() => setFiltro('concluido')}>
                Concluídas
              </button>
            </div>

            <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
              <button
                type="button"
                className="btn-secondary"
                disabled={ocupado === 'reprogramar'}
                onClick={reprogramarPendentes}
                style={{ padding: '6px 12px', fontSize: '12px', whiteSpace: 'nowrap' }}
                title="Reprograma em cascata todas as tarefas pendentes a partir de hoje respeitando sua carga horária diária"
              >
                {ocupado === 'reprogramar' ? 'Reprogramando…' : '⚡ Reprogramar Pendentes'}
              </button>

              <div className="schedule-view-toggle">
                <button
                  type="button"
                  className={modoVisao === 'lista' ? 'active' : ''}
                  onClick={() => setModoVisao('lista')}
                >
                  📋 Lista
                </button>
                <button
                  type="button"
                  className={modoVisao === 'calendario' ? 'active' : ''}
                  onClick={() => setModoVisao('calendario')}
                >
                  📅 Calendário
                </button>
              </div>
            </div>
          </div>

          {modoVisao === 'lista' ? (
            cronograma.tipo === 'ciclo_inteligente' ? (
              <section className="card-base student-schedule">
                <h2>Ordem do ciclo</h2>
                {itens.map((item: any) => (
                  <Atividade key={item.id} {...{ item, materias, topicos, subtopicos, ocupado, alterar, abrirModalConclusao, abrirTimer }} />
                ))}
              </section>
            ) : (
              <div className="schedule-days">
                {grupos.map(([data, lista]: any) => (
                  <section className="card-base student-schedule" key={data}>
                    <h2>{data === 'sem-data' ? 'Sem data definida' : dataBr(data)}</h2>
                    {lista.map((item: any, idx: number) => (
                      <Atividade key={item.id || idx} {...{ item, materias, topicos, subtopicos, ocupado, alterar, abrirModalConclusao, abrirTimer }} />
                    ))}
                  </section>
                ))}
              </div>
            )
          ) : (
            <div className="schedule-calendar-container card-base" style={{ padding: '18px' }}>
              <div className="schedule-calendar-header">
                <h2>
                  {calendarData.monthName} {calendarData.year}
                </h2>
                <div className="schedule-calendar-nav">
                  <button type="button" onClick={() => setMesOffset((m) => m - 1)}>
                    ◀ Mês anterior
                  </button>
                  <button type="button" onClick={() => setMesOffset(0)}>
                    Hoje
                  </button>
                  <button type="button" onClick={() => setMesOffset((m) => m + 1)}>
                    Próximo mês ▶
                  </button>
                </div>
              </div>

              <div style={{ overflowX: 'auto' }}>
                <div style={{ minWidth: '700px' }}>
                  <div className="schedule-calendar-weekdays">
                    {['Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb', 'Dom'].map((d) => (
                      <div key={d}>{d}</div>
                    ))}
                  </div>

                  <div className="schedule-calendar-grid">
                    {calendarData.days.map((day, i) => {
                      if (!day) return <div key={i} className="schedule-calendar-day empty" style={{ background: 'transparent', border: 'none' }} />;
                      const y = day.getFullYear();
                      const m = String(day.getMonth() + 1).padStart(2, '0');
                      const d = String(day.getDate()).padStart(2, '0');
                      const dateKey = `${y}-${m}-${d}`;
                      const isToday = dateKey === todayStr;
                      const dayItens = itensPorData.get(dateKey) || [];

                      return (
                        <div
                          key={i}
                          className={`schedule-calendar-day ${isToday ? 'is-today' : ''}`}
                          style={{ cursor: dayItens.length > 0 ? 'pointer' : 'default' }}
                          onClick={() => {
                            if (dayItens.length > 0) {
                              setDiaDetalhesModal({
                                dateKey,
                                dataFormatada: dataBr(dateKey),
                                dayItens
                              });
                            }
                          }}
                        >
                          <div className="schedule-calendar-day-header">
                            <span>{day.getDate()}</span>
                            {isToday && <small style={{ fontSize: '8px', color: 'var(--accent)', textTransform: 'uppercase' }}>Hoje</small>}
                            {dayItens.length > 0 && (
                              <small style={{ fontSize: '9px', color: 'var(--accent)', fontWeight: 800 }} title="Clique para ver lista detalhada do dia">
                                🔍 {dayItens.length}
                              </small>
                            )}
                          </div>
                          <div className="schedule-calendar-items">
                            {dayItens.map((item: any, idx: number) => {
                              if (item.tipoItem === 'prova') {
                                return (
                                  <div
                                    key={idx}
                                    style={{
                                      background: '#ef4444',
                                      color: '#fff',
                                      padding: '4px 6px',
                                      borderRadius: '6px',
                                      fontSize: '9px',
                                      fontWeight: 800,
                                      marginBottom: '4px'
                                    }}
                                  >
                                    {item.nome}
                                  </div>
                                );
                              }

                              if (item.tipoItem === 'revisao') {
                                const mat = materias.find((m: any) => m.id === item.materiaId);
                                const top = topicos.find((t: any) => t.id === item.topicoId);
                                const sub = subtopicos.find((s: any) => s.id === item.subtopicoId);
                                const mNome = item.materiaNome || item.materia || mat?.nome || 'Matéria';
                                const tNome = item.topicoNome || item.topico || top?.nome;
                                const sNome = item.subtopicoNome || item.subtopico || sub?.nome;

                                return (
                                  <div
                                    key={idx}
                                    style={{
                                      background: 'rgba(168, 85, 247, 0.18)',
                                      borderLeft: '3px solid #a855f7',
                                      color: '#e9d5ff',
                                      padding: '5px 7px',
                                      borderRadius: '5px',
                                      fontSize: '9.5px',
                                      marginBottom: '4px'
                                    }}
                                  >
                                    <strong style={{ color: '#c084fc', display: 'block', fontSize: '9px' }}>
                                      🔄 Revisão ({item.cicloAtual || 1}º Ciclo)
                                    </strong>
                                    <div style={{ fontWeight: 800, color: '#ffffff' }} title={mNome}>
                                      {mNome}
                                    </div>
                                    {tNome && (
                                      <div style={{ opacity: 0.9, color: 'var(--accent)', fontWeight: 600 }} title={tNome}>
                                        📖 {tNome}
                                      </div>
                                    )}
                                    {sNome && (
                                      <div style={{ opacity: 0.85, color: '#e2e8f0', fontSize: '9px' }} title={sNome}>
                                        ↳ {sNome}
                                      </div>
                                    )}
                                  </div>
                                );
                              }

                              // Item do Cronograma
                              const materia = materias.find((mat: any) => mat.id === item.materiaId);
                              const topico = topicos.find((t: any) => t.id === item.topicoId);
                              const subtopico = subtopicos.find((s: any) => s.id === item.subtopicoId);
                              const matNome = materia?.nome || item.materiaNome || 'Matéria';
                              const topNome = topico?.nome || item.topicoNome || (!item.subtopicoId ? item.nome : '');
                              const subNome = subtopico?.nome || item.subtopicoNome || (item.subtopicoId ? item.nome : '');

                              return (
                                <div key={item.id || idx} className={`schedule-calendar-item ${item.situacao}`}>
                                  <span className="schedule-calendar-item-materia" title={matNome}>
                                    {matNome}
                                  </span>
                                  <strong className="schedule-calendar-item-title" title={topNome || item.nome}>
                                    {topNome || item.nome || 'Atividade'}
                                  </strong>
                                  {subNome && (
                                    <div style={{ fontSize: '8.5px', color: 'var(--text2)', opacity: 0.9 }} title={subNome}>
                                      ↳ {subNome}
                                    </div>
                                  )}
                                  <div className="schedule-calendar-item-footer">
                                    <span className="schedule-calendar-item-duration">
                                      {item.duracaoMinutos || 0} min
                                    </span>
                                    <div className="schedule-calendar-item-actions" onClick={(e) => e.stopPropagation()}>
                                      {item.situacao !== 'concluido' && (
                                        <>
                                          <button
                                            type="button"
                                            className="btn-secondary"
                                            style={{ padding: '3px 6px', fontSize: '9px', fontWeight: 800 }}
                                            title="Iniciar temporizador deste tópico"
                                            onClick={(e) => abrirTimer(item, e)}
                                          >
                                            ⏱
                                          </button>
                                          <button
                                            type="button"
                                            className="btn-primary"
                                            style={{ padding: '3px 6px', fontSize: '9px' }}
                                            disabled={ocupado === item.id}
                                            onClick={(e) => abrirModalConclusao(item, e)}
                                          >
                                            ✓
                                          </button>
                                        </>
                                      )}
                                      {item.situacao === 'concluido' && (
                                        <button
                                          type="button"
                                          className="btn-secondary"
                                          style={{ padding: '3px 6px', fontSize: '9px' }}
                                          disabled={ocupado === item.id}
                                          onClick={() => alterar(item, 'pendente')}
                                        >
                                          ↩
                                        </button>
                                      )}
                                    </div>
                                  </div>
                                </div>
                              );
                            })}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              </div>

              {itensSemData.length > 0 && (
                <section className="student-schedule" style={{ marginTop: '20px', padding: '0' }}>
                  <h2>Atividades do Ciclo / Sem data definida</h2>
                  {itensSemData.map((item: any) => (
                    <Atividade key={item.id} {...{ item, materias, topicos, subtopicos, ocupado, alterar, abrirModalConclusao, abrirTimer }} />
                  ))}
                </section>
              )}
            </div>
          )}
        </>
      )}

      {/* MODAL DETALHADO AO CLICAR EM UM DIA DO CALENDÁRIO */}
      {diaDetalhesModal && (
        <div className="modal-backdrop" onClick={() => setDiaDetalhesModal(null)}>
          <div className="modal-card" style={{ maxWidth: '680px', width: '92%', maxHeight: '85vh', display: 'flex', flexDirection: 'column' }} onClick={(e) => e.stopPropagation()}>
            <div className="modal-heading" style={{ marginBottom: '16px', flexShrink: 0 }}>
              <div>
                <h2>📅 Tarefas & Revisões — {diaDetalhesModal.dataFormatada}</h2>
                <p>Confira a lista detalhada das atividades agendadas para este dia.</p>
              </div>
              <button className="icon-button" onClick={() => setDiaDetalhesModal(null)}>
                ×
              </button>
            </div>

            <div style={{ overflowY: 'auto', flex: 1, paddingRight: '4px' }}>
              {diaDetalhesModal.dayItens.length === 0 ? (
                <div className="empty-state" style={{ padding: '24px' }}>
                  Nenhuma atividade agendada para esta data.
                </div>
              ) : (
                diaDetalhesModal.dayItens.map((item: any, idx: number) => {
                  if (item.tipoItem === 'prova') {
                    return (
                      <div key={idx} style={{ background: '#ef4444', color: '#fff', padding: '12px 16px', borderRadius: '8px', marginBottom: '8px', fontWeight: 800 }}>
                        {item.nome}
                      </div>
                    );
                  }
                  return (
                    <Atividade
                      key={item.id || idx}
                      item={item}
                      materias={materias}
                      topicos={topicos}
                      subtopicos={subtopicos}
                      ocupado={ocupado}
                      alterar={alterar}
                      abrirModalConclusao={abrirModalConclusao}
                      abrirTimer={abrirTimer}
                    />
                  );
                })
              )}
            </div>
          </div>
        </div>
      )}

      {concluirModalItem && (
        <ModalConclusaoAtividade
          item={concluirModalItem}
          alunoId={alunoId}
          activeContestId={activeContestId}
          materias={materias}
          aoFechar={() => setConcluirModalItem(null)}
          aoConcluirSucesso={carregar}
        />
      )}
    </div>
  );
};

function Atividade({ item, materias, topicos, subtopicos, ocupado, alterar, abrirModalConclusao, abrirTimer }: any) {
  const mat = (materias || []).find((m: any) => m.id === item.materiaId);
  const top = (topicos || []).find((t: any) => t.id === item.topicoId);
  const sub = (subtopicos || []).find((s: any) => s.id === item.subtopicoId);

  const matNome = item.materiaNome || item.materia || mat?.nome || 'Matéria';
  const topNome = item.topicoNome || item.topico || top?.nome || (item.tipoItem !== 'revisao' && !item.subtopicoId ? item.nome : '');
  const subNome = item.subtopicoNome || item.subtopico || sub?.nome || (item.tipoItem !== 'revisao' && item.subtopicoId ? item.nome : '');

  if (item.tipoItem === 'revisao') {
    return (
      <article
        className="revisao-card"
        style={{
          borderLeft: '4px solid #a855f7',
          background: 'rgba(168, 85, 247, 0.08)',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          padding: '14px 18px',
          borderRadius: '10px',
          marginBottom: '8px'
        }}
      >
        <div>
          <small style={{ color: '#a855f7', fontWeight: 800, display: 'block', marginBottom: '3px', textTransform: 'uppercase', letterSpacing: '0.4px' }}>
            🔄 REVISÃO PROGRAMADA ({item.cicloAtual || 1}º CICLO) · 30 min previstos
          </small>
          <strong style={{ color: 'var(--text1)', fontSize: '15px', display: 'block' }}>{matNome}</strong>
          {topNome && <div style={{ fontSize: '13.5px', fontWeight: 700, color: 'var(--accent)', marginTop: '2px' }}>📖 Tópico: {topNome}</div>}
          {subNome && <div style={{ fontSize: '13px', color: 'var(--text2)', marginTop: '2px', fontWeight: 600 }}>↳ Subtópico: {subNome}</div>}
          {!topNome && !subNome && <div style={{ fontSize: '13px', color: 'var(--text2)', marginTop: '2px' }}>{item.nome || 'Assunto a Revisar'}</div>}
        </div>
        <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
          <button type="button" className="btn-secondary" onClick={(e) => abrirTimer(item, e)}>
            ⏱ Timer
          </button>
          <button
            type="button"
            className="btn-primary"
            style={{ background: '#a855f7', borderColor: '#a855f7' }}
            onClick={(e) => abrirModalConclusao(item, e)}
          >
            ✓ Concluir
          </button>
        </div>
      </article>
    );
  }

  return (
    <article className={item.situacao} style={{ padding: '14px 18px', borderRadius: '10px', marginBottom: '8px' }}>
      <div>
        <small style={{ fontWeight: 800, color: 'var(--accent)' }}>{matNome}</small>
        <strong style={{ display: 'block', fontSize: '15px', marginTop: '2px' }}>
          {topNome || item.nome || 'Atividade de estudo'}
        </strong>
        {subNome && (
          <div style={{ fontSize: '13px', color: 'var(--text2)', marginTop: '2px', fontWeight: 600 }}>
            ↳ Subtópico: {subNome}
          </div>
        )}
        <span style={{ marginTop: '6px', display: 'block', fontSize: '12px', color: 'var(--text3)' }}>
          {item.dataPlanejada
            ? dataBr(item.dataPlanejada)
            : `Posição ${item.posicaoCiclo || item.ordem} do ciclo`}{' '}
          · prioridade {item.prioridade || 0}
        </span>
      </div>
      <b>{item.duracaoMinutos || 0} min</b>
      <small>{rotulo(item.situacao)}</small>
      <div className="heading-actions">
        {item.situacao !== 'concluido' && (
          <>
            <button type="button" className="btn-secondary" onClick={(e) => abrirTimer(item, e)}>
              ⏱ Timer
            </button>
            <button
              type="button"
              className="btn-primary"
              disabled={ocupado === item.id}
              onClick={(e) => abrirModalConclusao(item, e)}
            >
              ✓ Concluir
            </button>
          </>
        )}

        {item.situacao === 'concluido' && (
          <button type="button" className="btn-secondary" disabled={ocupado === item.id} onClick={() => alterar(item, 'pendente')}>
            Reabrir
          </button>
        )}
      </div>
    </article>
  );
}
