import React, { useCallback, useEffect, useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useData } from '../context/DataContext';
import { api } from '../services/api';

const SPELL_INTERVALS = [1, 7, 30, 90];

function dateKey(d: Date) {
  return d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0') + '-' + String(d.getDate()).padStart(2, '0');
}

function formatDate(str: string) {
  if (!str) return '—';
  const parts = str.split('-');
  if (parts.length === 3) return `${parts[2]}/${parts[1]}/${parts[0]}`;
  return str;
}

function diffDays(dateStr: string) {
  const today = new Date(); today.setHours(0, 0, 0, 0);
  const target = new Date(dateStr + 'T00:00:00');
  return Math.round((target.getTime() - today.getTime()) / 86400000);
}

function pctColor(p: number | null) {
  if (p == null) return 'var(--text3)';
  if (p >= 93) return '#34f58d';
  if (p >= 80) return '#9cf7cc';
  if (p >= 75) return '#ffe45c';
  if (p >= 50) return 'var(--orange)';
  return 'var(--red)';
}

interface Revisao {
  id: string;
  concursoId: string;
  topico: string;
  materia: string;
  materiaId?: string;
  topicoId?: string;
  subtopicoId?: string;
  subtopId?: string;
  data?: string;
  dataProxima: string;
  cicloAtual: number;
  pct?: number | null;
  feito?: boolean;
  criadoEm?: string;
  concluidaEm?: string;
}

function classifyRevisoes(revisoes: Revisao[]) {
  const todayKey = dateKey(new Date());
  const atrasadas: Revisao[] = [];
  const hoje: Revisao[] = [];
  const futuras: Revisao[] = [];
  const feitas: Revisao[] = [];

  revisoes.forEach(r => {
    if (r.feito) { feitas.push(r); return; }
    const d = r.dataProxima || r.data || '';
    if (!d) { atrasadas.push(r); return; }
    if (d < todayKey) atrasadas.push(r);
    else if (d === todayKey) hoje.push(r);
    else futuras.push(r);
  });

  return { atrasadas, hoje, futuras, feitas };
}

export const RevisoesPage: React.FC = () => {
  const { alunoId, activeContestId, setActiveContestId, getArray } = useData();
  const navigate = useNavigate();
  const [allRevisoes,setAllRevisoes]=useState<Revisao[]>([]);
  const [questoesApi,setQuestoesApi]=useState<any[]>([]);
  const [carregando,setCarregando]=useState(true);
  const [erro,setErro]=useState('');
  const allMaterias = getArray('materias');
  const concursos = getArray('concursos');
  const carregar=useCallback(async()=>{setErro('');try{const[r,q]=await Promise.all([api.listarRevisoes(alunoId),api.listarQuestoes(alunoId)]);setAllRevisoes((r.revisoes||[]).map((x:any)=>({...x,concursoId:x.concursoId||'',data:x.proximaData||'',dataProxima:x.proximaData||'',pct:x.percentualAnterior,feito:!!x.concluida,subtopId:x.subtopicoId})));setQuestoesApi(q.registros||[])}catch(e:any){setErro(e.message)}finally{setCarregando(false)}},[alunoId]);
  useEffect(()=>{carregar()},[carregar]);

  // Filter by active contest, or all if no filter
  const [filtro, setFiltro] = useState<string>(activeContestId || 'todos');
  const revisoesConcurso = useMemo(() => {
    if (filtro === 'todos') return allRevisoes;
    return allRevisoes.filter((r: Revisao) => r.concursoId === filtro);
  }, [allRevisoes, filtro]);

  const classified = useMemo(() => classifyRevisoes(revisoesConcurso), [revisoesConcurso]);

  // Next 7 days agenda
  const nextDays = useMemo(() => {
    const days: { date: Date; key: string; count: number; colors: string[] }[] = [];
    const today = new Date(); today.setHours(0, 0, 0, 0);
    for (let i = 0; i < 7; i++) {
      const d = new Date(today); d.setDate(today.getDate() + i);
      const k = dateKey(d);
      const items = allRevisoes.filter((r: Revisao) => !r.feito && (r.dataProxima || r.data) === k);
      days.push({ date: d, key: k, count: items.length, colors: items.slice(0, 4).map(() => 'var(--accent)') });
    }
    return days;
  }, [allRevisoes]);

  // Performance alerts
  const alertas = useMemo(() => {
    return allMaterias
      .filter((m: any) => m.concursoId === filtro || filtro === 'todos')
      .map((m: any) => {
        const questoes = questoesApi.filter((q: any) => q.materiaId === m.id);
        const res = questoes.reduce((acc: number, q: any) => acc + (Number(q.resolvidas) || 0), 0);
        const ac = questoes.reduce((acc: number, q: any) => acc + (Number(q.acertos) || 0), 0);
        const pct = res >= 10 ? Math.round((ac / res) * 100) : null;
        return { ...m, pct, resolvidas: res };
      })
      .filter((m: any) => m.pct != null && m.pct < 65)
      .sort((a: any, b: any) => (a.pct || 0) - (b.pct || 0));
  }, [allMaterias, filtro, questoesApi]);

  const [showModal, setShowModal] = useState(false);
  const [selectedRevisao, setSelectedRevisao] = useState<Revisao | null>(null);
  const [acertos, setAcertos] = useState('');
  const [nextDate, setNextDate] = useState('');
  const [modalMode, setModalMode] = useState<'review' | 'postpone'>('review');

  const targetIds = (r: Revisao) => ({ materiaId: r.materiaId || '', topicoId: r.topicoId || '', subtopicoId: r.subtopicoId || r.subtopId || '' });
  const openTarget = (r: Revisao) => {
    const target = targetIds(r); if (r.concursoId) setActiveContestId(r.concursoId);
    const params = new URLSearchParams();
    if (target.subtopicoId) params.set('subtopico', target.subtopicoId);
    else if (target.topicoId) params.set('topico', target.topicoId);
    if (target.topicoId) params.set('topicoId', target.topicoId);
    if (target.materiaId) params.set('materiaId', target.materiaId);
    navigate(alunoId==='eu'?`/edital?${params.toString()}`:`/mentor/alunos/${alunoId}/concursos/${r.concursoId}/edital?${params.toString()}`);
  };
  const startTimer = (r: Revisao) => {
    if (r.concursoId) setActiveContestId(r.concursoId);
    window.dispatchEvent(new CustomEvent('ct:open-timer', { detail: { ...targetIds(r), autoStart: true } }));
  };

  const openRevisar = (r: Revisao) => {
    setSelectedRevisao(r);
    setAcertos('');
    setModalMode('review');
    const d = new Date(); d.setDate(d.getDate() + SPELL_INTERVALS[Math.min(r.cicloAtual || 0, SPELL_INTERVALS.length - 1)]); setNextDate(dateKey(d));
    setShowModal(true);
  };

  const openAdiar = (r: Revisao) => {
    setSelectedRevisao(r); setModalMode('postpone');
    const d = new Date(); d.setDate(d.getDate() + 7); setNextDate(dateKey(d)); setShowModal(true);
  };

  const confirmarRevisao = async (continueReview = true) => {
    if (!selectedRevisao) return;
    const ciclo = (selectedRevisao.cicloAtual || 0) + 1;
    const pctNum = acertos !== '' ? Math.min(100, Math.max(0, Number(acertos))) : null;
    setErro('');try{await api.alterarRevisao(alunoId,selectedRevisao.id,{cicloAtual:ciclo,proximaData:continueReview?nextDate:'',percentualAnterior:pctNum,concluida:!continueReview});setShowModal(false);await carregar()}catch(e:any){setErro(e.message)}
  };

  const confirmarAdiamento = async () => {
    if (!selectedRevisao || !nextDate) return;
    setErro('');try{await api.alterarRevisao(alunoId,selectedRevisao.id,{proximaData:nextDate,concluida:false});setShowModal(false);await carregar()}catch(e:any){setErro(e.message)}
  };

  const deletarRevisao = async (id: string) => {
    if (!confirm('Desativar esta revisão? O histórico será preservado.')) return;
    setErro('');try{await api.desativarRevisao(alunoId,id);await carregar()}catch(e:any){setErro(e.message)}
  };


  function RevCard({ r, tipo }: { r: Revisao; tipo: 'atrasada' | 'hoje' | 'futura' }) {
    const reviewDate = r.dataProxima || r.data || '';
    const diff = reviewDate ? diffDays(reviewDate) : null;
    const concurso = concursos.find((c: any) => c.id === r.concursoId);
    const borderColor = tipo === 'atrasada' ? 'var(--red)' : tipo === 'hoje' ? 'var(--yellow)' : 'var(--border2)';

    return (
      <div style={{ background: 'var(--bg2)', border: '1px solid var(--border)', borderLeft: `3px solid ${borderColor}`, borderRadius: '10px', padding: '14px 16px', display: 'flex', alignItems: 'center', gap: '12px', transition: 'all .15s', flexWrap: 'wrap' }}>
        <div style={{ flex: 1, minWidth: 0 }}>
          <div style={{ fontSize: '11px', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '1px', color: 'var(--text3)', marginBottom: '3px' }}>
            {concurso?.nome || r.concursoId} {r.materia && <span>• {r.materia}</span>}
          </div>
          <button onClick={() => openTarget(r)} title="Abrir no edital" style={{ display: 'block', width: '100%', padding: 0, border: 0, background: 'transparent', textAlign: 'left', cursor: 'pointer', fontSize: '14px', fontWeight: 600, color: 'var(--text)', marginBottom: '3px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{r.topico || 'Revisão'}</button>
          <div style={{ fontSize: '12px', color: 'var(--text3)', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <span>Ciclo {r.cicloAtual || 0} de {SPELL_INTERVALS.length}</span>
            {r.pct != null && <span style={{ color: pctColor(r.pct), fontWeight: 700, fontFamily: 'monospace' }}>{r.pct}%</span>}
          </div>
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: '6px', flexShrink: 0 }}>
          <span style={{
            padding: '3px 10px', borderRadius: '20px', fontSize: '12px', fontWeight: 700, fontFamily: 'monospace',
            background: tipo === 'atrasada' ? 'rgba(245,90,90,.12)' : tipo === 'hoje' ? 'rgba(255,215,0,.12)' : 'rgba(79,142,247,.08)',
            color: tipo === 'atrasada' ? 'var(--red)' : tipo === 'hoje' ? 'var(--yellow)' : 'var(--text3)',
          }}>
            {tipo === 'atrasada' ? `${Math.abs(diff || 0)}d atraso` : tipo === 'hoje' ? 'Hoje' : reviewDate ? formatDate(reviewDate) : '—'}
          </span>
          <div style={{ display: 'flex', gap: '6px', alignItems: 'center' }}>
            {alunoId==='eu'&&<button onClick={() => startTimer(r)} title="Iniciar cronômetro nesta revisão" style={{ padding: '6px 9px', borderRadius: '7px', border: '1px solid var(--border2)', background: 'var(--bg3)', color: 'var(--text2)', cursor: 'pointer' }}>⏱</button>}
            <button onClick={() => openRevisar(r)} style={{ padding: '6px 14px', borderRadius: '7px', border: 'none', background: 'linear-gradient(135deg, var(--accent), var(--accent2))', color: '#fff', fontSize: '12px', fontWeight: 700, cursor: 'pointer' }}>
              {tipo === 'hoje' ? '✅ Revisar' : '📖 Revisar'}
            </button>
            {tipo === 'atrasada' && (
              <button onClick={() => openAdiar(r)} style={{ padding: '6px 10px', borderRadius: '7px', border: '1px solid var(--border2)', background: 'var(--bg3)', color: 'var(--text3)', fontSize: '12px', cursor: 'pointer' }}>Adiar</button>
            )}
            <button onClick={() => deletarRevisao(r.id)} style={{ padding: '5px 8px', border: '1px solid var(--border)', borderRadius: '6px', background: 'transparent', color: 'var(--text3)', fontSize: '12px', cursor: 'pointer' }}>✕</button>
          </div>
        </div>
      </div>
    );
  }

  if(carregando)return <div className="loading-state">Carregando revisões…</div>;
  return (
    <div>
      {erro&&<div className="form-error">{erro}</div>}
      {/* Header */}
      <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', marginBottom: '24px', flexWrap: 'wrap', gap: '12px' }}>
        <div>
          <h1 style={{ fontSize: '24px', fontWeight: 800, letterSpacing: '-.5px' }}>🔄 Revisões Programadas</h1>
          <p style={{ color: 'var(--text3)', fontSize: '13px', marginTop: '4px' }}>Ciclo de repetição espaçada: 1d → 7d → 30d → 90d</p>
        </div>
        {/* Filtro por concurso */}
        {concursos.length > 1 && (
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>
            <span style={{ fontSize: '12px', fontWeight: 600, color: 'var(--text3)', textTransform: 'uppercase', letterSpacing: '1px' }}>Concurso:</span>
            <button onClick={() => setFiltro('todos')} style={{ padding: '5px 14px', borderRadius: '20px', border: '1px solid var(--border2)', background: filtro === 'todos' ? 'rgba(79,142,247,.12)' : 'var(--bg3)', color: filtro === 'todos' ? 'var(--accent)' : 'var(--text2)', fontSize: '12px', fontWeight: 600, cursor: 'pointer' }}>Todos</button>
            {concursos.map((c: any) => (
              <button key={c.id} onClick={() => {setFiltro(c.id);setActiveContestId(c.id)}} style={{ padding: '5px 14px', borderRadius: '20px', border: '1px solid var(--border2)', background: filtro === c.id ? 'rgba(79,142,247,.12)' : 'var(--bg3)', color: filtro === c.id ? 'var(--accent)' : 'var(--text2)', fontSize: '12px', fontWeight: 600, cursor: 'pointer' }}>{c.nome}</button>
            ))}
          </div>
        )}
      </div>

      {/* Performance alerts */}
      {alertas.length > 0 && (
        <div style={{ background: 'rgba(245,90,90,.06)', border: '1px dashed var(--red)', borderRadius: '14px', padding: '18px', marginBottom: '24px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '12px' }}>
            <span style={{ fontSize: '18px' }}>⚠️</span>
            <div style={{ fontSize: '14px', fontWeight: 800, color: 'var(--red)', textTransform: 'uppercase', letterSpacing: '.5px' }}>Atenção: Rendimento Baixo em Questões</div>
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            {alertas.map((m: any) => (
              <div key={m.id} style={{ display: 'flex', alignItems: 'center', gap: '10px', fontSize: '12px', color: 'var(--text2)', background: 'rgba(245,90,90,.04)', border: '1px solid rgba(245,90,90,.15)', borderRadius: '8px', padding: '10px 12px' }}>
                <span style={{ color: 'var(--red)', fontFamily: 'monospace', fontWeight: 700, fontSize: '14px' }}>{m.pct}%</span>
                <span style={{ fontWeight: 700, color: 'var(--text)' }}>{m.nome}</span>
                <span style={{ color: 'var(--text3)' }}>({m.resolvidas} questões resolvidas)</span>
                <span style={{ marginLeft: 'auto', color: 'var(--red)', fontSize: '11px', fontWeight: 700 }}>REVISAR COM PRIORIDADE</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Summary */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '12px', marginBottom: '24px' }}>
        {[
          { val: classified.atrasadas.length, label: '🔴 Atrasadas', bg: 'rgba(245,90,90,.08)', border: 'rgba(245,90,90,.2)', color: 'var(--red)' },
          { val: classified.hoje.length, label: '🟡 Para hoje', bg: 'rgba(255,215,0,.08)', border: 'rgba(255,215,0,.2)', color: 'var(--yellow)' },
          { val: classified.futuras.length, label: '📅 Futuras', bg: 'rgba(79,142,247,.08)', border: 'rgba(79,142,247,.2)', color: 'var(--accent)' },
          { val: classified.feitas.length, label: '✅ Concluídas', bg: 'rgba(62,207,142,.08)', border: 'rgba(62,207,142,.2)', color: 'var(--green)' },
        ].map((c, i) => (
          <div key={i} style={{ background: c.bg, border: `1px solid ${c.border}`, borderRadius: '12px', padding: '16px 20px', textAlign: 'center' }}>
            <div style={{ fontFamily: 'monospace', fontSize: '28px', fontWeight: 700, color: c.color }}>{c.val}</div>
            <div style={{ fontSize: '12px', color: 'var(--text3)', marginTop: '4px' }}>{c.label}</div>
          </div>
        ))}
      </div>

      {/* 7-day agenda */}
      <div className="card-base" style={{ padding: '16px 18px', marginBottom: '24px' }}>
        <div style={{ fontSize: '13px', fontWeight: 800, marginBottom: '12px' }}>📅 Agenda — Próximos 7 dias</div>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, minmax(0,1fr))', gap: '8px' }}>
          {nextDays.map((d, i) => {
            const isToday = i === 0;
            return (
              <div key={i} style={{
                background: isToday ? 'rgba(79,142,247,.05)' : 'var(--bg2)',
                border: `1px solid ${isToday ? 'rgba(79,142,247,.4)' : 'var(--border)'}`,
                borderRadius: '10px', padding: '10px 8px', textAlign: 'center',
              }}>
                <div style={{ fontSize: '11px', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '1px', color: isToday ? 'var(--accent)' : 'var(--text3)', marginBottom: '4px' }}>
                  {['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb'][d.date.getDay()]}
                </div>
                <div style={{ fontFamily: 'monospace', fontSize: '16px', fontWeight: 700, color: isToday ? 'var(--accent)' : 'var(--text)', marginBottom: '8px' }}>{d.date.getDate()}</div>
                <div style={{ fontFamily: 'monospace', fontSize: d.count > 0 ? '18px' : '14px', fontWeight: 700, color: d.count > 0 ? 'var(--text)' : 'var(--text3)', lineHeight: 1 }}>{d.count > 0 ? d.count : '·'}</div>
                {d.count > 0 && <div style={{ fontSize: '11px', color: 'var(--text3)', marginTop: '3px' }}>revisão{d.count !== 1 ? 'ões' : ''}</div>}
              </div>
            );
          })}
        </div>
      </div>

      {/* Sections */}
      {classified.atrasadas.length === 0 && classified.hoje.length === 0 && classified.futuras.length === 0 && (
        <div style={{ textAlign: 'center', padding: '80px 20px', color: 'var(--text3)' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>✅</div>
          <div style={{ fontSize: '18px', fontWeight: 700, color: 'var(--text2)', marginBottom: '8px' }}>Nenhuma revisão pendente</div>
          <div style={{ fontSize: '14px', lineHeight: 1.6 }}>Quando você marcar um tópico como estudado<br />e agendar revisão, ele aparecerá aqui.</div>
        </div>
      )}

      {classified.atrasadas.length > 0 && (
        <div style={{ marginBottom: '28px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '14px' }}>
            <div style={{ width: '10px', height: '10px', borderRadius: '50%', background: 'var(--red)', flexShrink: 0 }} />
            <div style={{ fontSize: '14px', fontWeight: 700 }}>Atrasadas</div>
            <div style={{ fontSize: '12px', color: 'var(--text3)', background: 'var(--bg3)', border: '1px solid var(--border)', padding: '2px 8px', borderRadius: '20px', fontFamily: 'monospace' }}>{classified.atrasadas.length}</div>
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            {classified.atrasadas.map(r => <RevCard key={r.id} r={r} tipo="atrasada" />)}
          </div>
        </div>
      )}

      {classified.hoje.length > 0 && (
        <div style={{ marginBottom: '28px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '14px' }}>
            <div style={{ width: '10px', height: '10px', borderRadius: '50%', background: 'var(--yellow)', flexShrink: 0 }} />
            <div style={{ fontSize: '14px', fontWeight: 700 }}>Para hoje</div>
            <div style={{ fontSize: '12px', color: 'var(--text3)', background: 'var(--bg3)', border: '1px solid var(--border)', padding: '2px 8px', borderRadius: '20px', fontFamily: 'monospace' }}>{classified.hoje.length}</div>
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            {classified.hoje.map(r => <RevCard key={r.id} r={r} tipo="hoje" />)}
          </div>
        </div>
      )}

      {classified.futuras.length > 0 && (
        <div style={{ marginBottom: '28px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '14px' }}>
            <div style={{ width: '10px', height: '10px', borderRadius: '50%', background: 'var(--accent)', flexShrink: 0 }} />
            <div style={{ fontSize: '14px', fontWeight: 700 }}>Próximas revisões</div>
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            {classified.futuras.slice(0, 20).map(r => <RevCard key={r.id} r={r} tipo="futura" />)}
            {classified.futuras.length > 20 && (
              <div style={{ textAlign: 'center', fontSize: '12px', color: 'var(--text3)', padding: '12px' }}>... e mais {classified.futuras.length - 20} futuras</div>
            )}
          </div>
        </div>
      )}

      {classified.feitas.length>0&&<details className="card-base" style={{marginBottom:'24px',padding:'16px 18px'}}><summary style={{cursor:'pointer',fontWeight:800}}>Histórico de revisões concluídas ({classified.feitas.length})</summary><div className="history-log-list" style={{marginTop:'12px'}}>{classified.feitas.map(r=><div key={r.id}><span>✅</span><div><strong>{r.topico||'Revisão'}</strong><small>{r.materia||'Conteúdo do edital'} · ciclo {r.cicloAtual||0}</small></div>{r.pct!=null&&<b style={{color:pctColor(r.pct)}}>{r.pct}%</b>}<button className="icon-button danger" onClick={()=>deletarRevisao(r.id)}>×</button></div>)}</div></details>}

      {/* Modal revisar */}
      {showModal && selectedRevisao && (
        <div onClick={() => setShowModal(false)} style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,.6)', backdropFilter: 'blur(4px)', zIndex: 200, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <div onClick={e => e.stopPropagation()} style={{ background: 'var(--bg2)', border: '1px solid var(--border2)', borderRadius: '16px', padding: '28px', maxWidth: '420px', width: '90%', boxShadow: '0 24px 60px rgba(0,0,0,.5)' }}>
            <div style={{ fontSize: '32px', textAlign: 'center', marginBottom: '12px' }}>📖</div>
            <div style={{ fontSize: '17px', fontWeight: 700, color: 'var(--text)', textAlign: 'center', marginBottom: '6px' }}>{modalMode === 'postpone' ? 'Adiar revisão' : 'Registrar Revisão'}</div>
            <div style={{ fontSize: '12px', color: 'var(--accent)', textAlign: 'center', marginBottom: '16px', fontWeight: 600, background: 'rgba(79,142,247,.08)', padding: '6px 12px', borderRadius: '20px', display: 'block' }}>{selectedRevisao.topico}</div>

            <div style={{ fontSize: '12px', color: 'var(--text3)', lineHeight: 1.6, marginBottom: '16px', background: 'var(--bg3)', borderRadius: '8px', padding: '10px 14px', borderLeft: '3px solid var(--accent2)' }}>
              Ciclo atual: <strong style={{ color: 'var(--text)' }}>{selectedRevisao.cicloAtual || 0}</strong> de {SPELL_INTERVALS.length}.<br />
              Próximos intervalos: {SPELL_INTERVALS.slice((selectedRevisao.cicloAtual || 0) + 1).map(n => `${n}d`).join(' → ')}
            </div>

            {modalMode === 'review' && <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '8px', marginBottom: '12px' }}>
              {[70, 80, 90, 100].map(pct => (
                <button key={pct} onClick={() => setAcertos(String(pct))} style={{ padding: '10px', borderRadius: '9px', border: `2px solid ${acertos === String(pct) ? 'var(--accent)' : 'var(--border2)'}`, background: acertos === String(pct) ? 'rgba(79,142,247,.1)' : 'var(--bg3)', cursor: 'pointer', textAlign: 'center' }}>
                  <div style={{ fontFamily: 'monospace', fontSize: '20px', fontWeight: 700, color: 'var(--accent)', lineHeight: 1 }}>{pct}%</div>
                  <div style={{ fontSize: '12px', color: 'var(--text3)', marginTop: '3px' }}>{pct >= 90 ? 'Ótimo' : pct >= 80 ? 'Bom' : pct >= 70 ? 'Regular' : 'Difícil'}</div>
                </button>
              ))}
            </div>}

            {modalMode === 'review' && <div style={{ marginBottom: '12px' }}>
              <input type="number" min={0} max={100} value={acertos} onChange={e => setAcertos(e.target.value)} placeholder="% de acertos (opcional)" style={{ width: '100%', background: 'var(--bg3)', border: '1px solid var(--border2)', borderRadius: '8px', padding: '9px 12px', color: 'var(--text)', fontSize: '12px', outline: 'none' }} />
            </div>}
            <label style={{ display: 'block', marginBottom: '16px' }}><span style={{ display: 'block', fontSize: '12px', color: 'var(--text3)', marginBottom: '6px' }}>{modalMode === 'postpone' ? 'Nova data' : 'Próxima revisão'}</span><input type="date" min={dateKey(new Date())} value={nextDate} onChange={e => setNextDate(e.target.value)} style={{ width: '100%', background: 'var(--bg3)', border: '1px solid var(--border2)', borderRadius: '8px', padding: '9px 12px', color: 'var(--text)' }} /></label>

            <div style={{ display: 'flex', gap: '8px' }}>
              <button onClick={() => setShowModal(false)} style={{ flex: 1, padding: '10px', borderRadius: '8px', border: '1px solid var(--border2)', background: 'transparent', color: 'var(--text3)', fontSize: '12px', cursor: 'pointer' }}>Cancelar</button>
              {modalMode === 'review' && <button onClick={() => confirmarRevisao(false)} className="btn-secondary" style={{ flex: 1, padding: '10px', fontSize: '12px' }}>Não revisar mais</button>}
              <button onClick={modalMode === 'postpone' ? confirmarAdiamento : () => confirmarRevisao(true)} className="btn-primary" style={{ flex: 2, padding: '10px', fontSize: '12px' }}>{modalMode === 'postpone' ? 'Salvar data' : '✅ Confirmar Revisão'}</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
