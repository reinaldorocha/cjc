import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../services/api';

interface ItemRadar {
  alunoId: string;
  alunoNome: string;
  alunoEmail: string;
  telefone?: string;
  concursoNome: string;
  ultimoEstudoEm?: string;
  diasSemEstudar: number;
  segundosSemana: number;
  segundosMes: number;
  questoesResolvidas: number;
  questoesAcertos: number;
  taxaAcerto: number;
  revisoesPendentes: number;
  nivelRisco: 'vermelho' | 'amarelo' | 'verde';
  motivoRisco: string;
}

export const RadarMentoriaPage: React.FC = () => {
  const navigate = useNavigate();
  const [alunos, setAlunos] = useState<ItemRadar[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [erro, setErro] = useState('');
  const [filtroRisco, setFiltroRisco] = useState<'todos' | 'vermelho' | 'amarelo' | 'verde'>('todos');
  const [busca, setBusca] = useState('');

  const carregar = useCallback(async () => {
    setCarregando(true);
    setErro('');
    try {
      const res = await api.obterRadarAlunos();
      setAlunos(res.alunos || []);
    } catch (err: any) {
      setErro(err.message || 'Erro ao carregar o radar da mentoria.');
    } finally {
      setCarregando(false);
    }
  }, []);

  useEffect(() => {
    carregar();
  }, [carregar]);

  const formatHoras = (sec: number) => {
    const h = Math.floor(sec / 3600);
    const m = Math.floor((sec % 3600) / 60);
    return h ? `${h}h ${m}m` : `${m}min`;
  };

  const alunosFiltrados = alunos.filter((a) => {
    const batemRisco = filtroRisco === 'todos' || a.nivelRisco === filtroRisco;
    const batemBusca = !busca.trim() || a.alunoNome.toLowerCase().includes(busca.toLowerCase()) || a.alunoEmail.toLowerCase().includes(busca.toLowerCase()) || a.concursoNome.toLowerCase().includes(busca.toLowerCase());
    return batemRisco && batemBusca;
  });

  const countVermelho = alunos.filter((a) => a.nivelRisco === 'vermelho').length;
  const countAmarelo = alunos.filter((a) => a.nivelRisco === 'amarelo').length;
  const countVerde = alunos.filter((a) => a.nivelRisco === 'verde').length;

  const abrirWhatsApp = (a: ItemRadar) => {
    const num = (a.telefone || '').replace(/\D/g, '');
    let msg = '';
    if (a.nivelRisco === 'vermelho') {
      msg = `Olá ${a.alunoNome}, tudo bem? Notei na plataforma Chega Junto Concurseiro que você está sem registrar estudos há ${a.diasSemEstudar > 30 ? 'alguns' : a.diasSemEstudar} dias. Como posso te ajudar a retomar o ritmo do ${a.concursoNome}?`;
    } else if (a.nivelRisco === 'amarelo') {
      msg = `Olá ${a.alunoNome}! Passando para acompanhar seus estudos do ${a.concursoNome}. Vi que temos algumas pendências/revisões para alinhar. Como foi seu dia de estudos?`;
    } else {
      msg = `Parabéns ${a.alunoNome}! Seu ritmo de estudos no ${a.concursoNome} está constante e excelente. Continue assim! 🔥`;
    }

    const url = num ? `https://wa.me/55${num}?text=${encodeURIComponent(msg)}` : `https://wa.me/?text=${encodeURIComponent(msg)}`;
    window.open(url, '_blank');
  };

  if (carregando) return <div className="loading-state">Carregando Radar da Mentoria…</div>;

  return (
    <div>
      {erro && <div className="form-error">{erro}</div>}

      <div className="page-heading">
        <div>
          <h1>🚨 Radar da Mentoria</h1>
          <p>Painel de acompanhamento em tempo real e monitoramento de risco dos alunos</p>
        </div>
        <button className="btn-secondary" onClick={carregar}>
          🔄 Atualizar Radar
        </button>
      </div>

      {/* Cards de Resumo dos Alertas */}
      <div className="stat-grid" style={{ marginBottom: '24px' }}>
        <div
          className="stat-card"
          style={{ cursor: 'pointer', borderLeft: filtroRisco === 'todos' ? '4px solid var(--accent)' : '1px solid var(--border)' }}
          onClick={() => setFiltroRisco('todos')}
        >
          <span>TOTAL DE ALUNOS</span>
          <strong>{alunos.length}</strong>
          <small>na sua mentoria</small>
        </div>
        <div
          className="stat-card"
          style={{ cursor: 'pointer', borderLeft: filtroRisco === 'vermelho' ? '4px solid var(--red)' : '1px solid var(--border)', background: countVermelho > 0 ? 'rgba(245, 90, 90, 0.08)' : undefined }}
          onClick={() => setFiltroRisco('vermelho')}
        >
          <span style={{ color: 'var(--red)', fontWeight: 800 }}>🔴 EM RISCO ALTO</span>
          <strong style={{ color: 'var(--red)' }}>{countVermelho}</strong>
          <small>2+ dias inativos ou baixo acerto</small>
        </div>
        <div
          className="stat-card"
          style={{ cursor: 'pointer', borderLeft: filtroRisco === 'amarelo' ? '4px solid var(--yellow)' : '1px solid var(--border)', background: countAmarelo > 0 ? 'rgba(245, 200, 66, 0.08)' : undefined }}
          onClick={() => setFiltroRisco('amarelo')}
        >
          <span style={{ color: 'var(--yellow)', fontWeight: 800 }}>🟡 EM ATENÇÃO</span>
          <strong style={{ color: 'var(--yellow)' }}>{countAmarelo}</strong>
          <small>1 dia sem estudar ou revisões pendentes</small>
        </div>
        <div
          className="stat-card"
          style={{ cursor: 'pointer', borderLeft: filtroRisco === 'verde' ? '4px solid var(--green)' : '1px solid var(--border)' }}
          onClick={() => setFiltroRisco('verde')}
        >
          <span style={{ color: 'var(--green)', fontWeight: 800 }}>🟢 RITMO IDEAL</span>
          <strong style={{ color: 'var(--green)' }}>{countVerde}</strong>
          <small>estudando diariamente</small>
        </div>
      </div>

      {/* Filtros e Busca */}
      <div className="card-base" style={{ padding: '16px 20px', marginBottom: '24px', display: 'flex', gap: '16px', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between' }}>
        <div className="period-chips" style={{ margin: 0 }}>
          <button className={filtroRisco === 'todos' ? 'active' : ''} onClick={() => setFiltroRisco('todos')}>
            Todos ({alunos.length})
          </button>
          <button className={filtroRisco === 'vermelho' ? 'active' : ''} style={filtroRisco === 'vermelho' ? { background: 'var(--red)', borderColor: 'var(--red)' } : {}} onClick={() => setFiltroRisco('vermelho')}>
            🔴 Em Risco ({countVermelho})
          </button>
          <button className={filtroRisco === 'amarelo' ? 'active' : ''} style={filtroRisco === 'amarelo' ? { background: 'var(--yellow)', borderColor: 'var(--yellow)', color: '#000' } : {}} onClick={() => setFiltroRisco('amarelo')}>
            🟡 Atenção ({countAmarelo})
          </button>
          <button className={filtroRisco === 'verde' ? 'active' : ''} style={filtroRisco === 'verde' ? { background: 'var(--green)', borderColor: 'var(--green)' } : {}} onClick={() => setFiltroRisco('verde')}>
            🟢 Ritmo Ideal ({countVerde})
          </button>
        </div>
        <div style={{ minWidth: '240px', flex: '1', maxWidth: '360px' }}>
          <input
            className="form-control"
            placeholder="🔍 Buscar aluno ou concurso..."
            value={busca}
            onChange={(e) => setBusca(e.target.value)}
          />
        </div>
      </div>

      {/* Grid de Cards dos Alunos no Radar */}
      <div className="catalog-grid" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(340px, 1fr))' }}>
        {alunosFiltrados.map((a) => {
          const corBorda = a.nivelRisco === 'vermelho' ? 'var(--red)' : a.nivelRisco === 'amarelo' ? 'var(--yellow)' : 'var(--green)';
          const bgBadge = a.nivelRisco === 'vermelho' ? 'rgba(245, 90, 90, 0.15)' : a.nivelRisco === 'amarelo' ? 'rgba(245, 200, 66, 0.15)' : 'rgba(62, 207, 142, 0.15)';
          const textoBadge = a.nivelRisco === 'vermelho' ? '🔴 RISCO ALTO' : a.nivelRisco === 'amarelo' ? '🟡 ATENÇÃO' : '🟢 RITMO IDEAL';

          return (
            <article
              className="premium-card"
              key={a.alunoId}
              style={{
                borderLeft: `4px solid ${corBorda}`,
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'space-between'
              }}
            >
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
                  <span
                    style={{
                      background: bgBadge,
                      color: corBorda,
                      fontSize: '11px',
                      fontWeight: 800,
                      padding: '4px 10px',
                      borderRadius: '12px',
                      letterSpacing: '0.5px'
                    }}
                  >
                    {textoBadge}
                  </span>
                  <span style={{ fontSize: '11px', color: 'var(--text3)' }}>
                    {a.diasSemEstudar === 0
                      ? '⏱️ Estudou hoje'
                      : a.diasSemEstudar === 1
                      ? '⏱️ Estudou ontem'
                      : a.diasSemEstudar >= 999
                      ? '⚠️ Sem registros'
                      : `⚠️ Há ${a.diasSemEstudar}d sem estudar`}
                  </span>
                </div>

                <h2 style={{ fontSize: '18px', margin: '0 0 4px', fontWeight: 800, color: 'var(--text1)' }}>{a.alunoNome}</h2>
                <p style={{ fontSize: '12px', color: 'var(--accent)', margin: '0 0 12px', fontWeight: 700 }}>🎯 {a.concursoNome}</p>

                {/* Motivo de Risco / Alerta */}
                <div
                  style={{
                    background: 'var(--bg3)',
                    padding: '8px 12px',
                    borderRadius: '8px',
                    fontSize: '12px',
                    marginBottom: '14px',
                    color: 'var(--text2)',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '6px'
                  }}
                >
                  <span>📌</span>
                  <strong>{a.motivoRisco}</strong>
                </div>

                {/* Indicadores do Aluno */}
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px', fontSize: '12px', marginBottom: '14px' }}>
                  <div style={{ background: 'var(--bg2)', padding: '8px', borderRadius: '6px' }}>
                    <span style={{ fontSize: '10px', color: 'var(--text3)', display: 'block' }}>TEMPO NA SEMANA</span>
                    <strong style={{ color: 'var(--text1)' }}>{formatHoras(a.segundosSemana)}</strong>
                  </div>
                  <div style={{ background: 'var(--bg2)', padding: '8px', borderRadius: '6px' }}>
                    <span style={{ fontSize: '10px', color: 'var(--text3)', display: 'block' }}>TAXA DE ACERTO</span>
                    <strong style={{ color: a.taxaAcerto >= 70 ? 'var(--green)' : a.taxaAcerto >= 50 ? 'var(--yellow)' : 'var(--red)' }}>
                      {a.questoesResolvidas > 0 ? `${Math.round(a.taxaAcerto)}%` : '—'}
                    </strong>
                  </div>
                </div>
              </div>

              {/* Botões de Ação Rápida */}
              <div style={{ display: 'flex', gap: '8px', marginTop: '12px', paddingTop: '12px', borderTop: '1px solid var(--border)' }}>
                <button
                  type="button"
                  className="btn-secondary"
                  style={{ flex: 1, padding: '8px', fontSize: '12px', fontWeight: 700 }}
                  onClick={() => navigate(`/aluno/${a.alunoId}`)}
                >
                  👁️ Ver Perfil
                </button>
                <button
                  type="button"
                  className="btn-primary"
                  style={{
                    flex: 1,
                    padding: '8px',
                    fontSize: '12px',
                    fontWeight: 700,
                    background: '#25D366',
                    borderColor: '#25D366',
                    color: '#fff'
                  }}
                  onClick={() => abrirWhatsApp(a)}
                >
                  📱 WhatsApp
                </button>
              </div>
            </article>
          );
        })}

        {!alunosFiltrados.length && (
          <div className="compact-empty" style={{ gridColumn: '1 / -1', padding: '40px 20px', textAlign: 'center' }}>
            Nenhum aluno encontrado para os filtros selecionados.
          </div>
        )}
      </div>
    </div>
  );
};
