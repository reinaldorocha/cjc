import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

interface ExamCountdownBannerProps {
  activeContest: any;
  alunoId: string;
  aoAtualizarConcurso?: () => void;
}

export const ExamCountdownBanner: React.FC<ExamCountdownBannerProps> = ({
  activeContest,
  alunoId,
  aoAtualizarConcurso
}) => {
  const [tempoRestante, setTempoRestante] = useState<{ dias: number; horas: number; minutos: number; segundos: number } | null>(null);
  const [modalDataProva, setModalDataProva] = useState(false);
  const [novaDataProva, setNovaDataProva] = useState('');
  const [salvandoData, setSalvandoData] = useState(false);

  useEffect(() => {
    if (!activeContest?.dataProva) {
      setTempoRestante(null);
      return;
    }

    const dataAlvo = new Date(`${activeContest.dataProva}T08:00:00`).getTime();

    const calcular = () => {
      const agora = new Date().getTime();
      const diff = dataAlvo - agora;

      if (diff <= 0) {
        setTempoRestante({ dias: 0, horas: 0, minutos: 0, segundos: 0 });
        return;
      }

      const dias = Math.floor(diff / (1000 * 60 * 60 * 24));
      const horas = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
      const minutos = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
      const segundos = Math.floor((diff % (1000 * 60)) / 1000);

      setTempoRestante({ dias, horas, minutos, segundos });
    };

    calcular();
    const interval = setInterval(calcular, 1000);
    return () => clearInterval(interval);
  }, [activeContest?.dataProva]);

  const salvarDataProva = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeContest?.id || !novaDataProva) return;
    setSalvandoData(true);
    try {
      await api.alterarConcurso(alunoId, activeContest.id, { dataProva: novaDataProva });
      setModalDataProva(false);
      if (aoAtualizarConcurso) aoAtualizarConcurso();
    } catch (err: any) {
      alert(err.message || 'Erro ao atualizar data da prova');
    } finally {
      setSalvandoData(false);
    }
  };

  const dataProvaStr = activeContest?.dataProva || '';

  return (
    <div style={{ marginBottom: '24px' }}>
      <div
        className="card-base"
        style={{
          background: tempoRestante
            ? 'linear-gradient(135deg, #1e1b4b 0%, #312e81 50%, #4338ca 100%)'
            : 'linear-gradient(135deg, #18181b 0%, #27272a 100%)',
          color: '#ffffff',
          borderRadius: '16px',
          padding: '24px 28px',
          boxShadow: '0 10px 30px rgba(0,0,0,0.3)',
          border: '1px solid rgba(255,255,255,0.1)',
          position: 'relative',
          overflow: 'hidden'
        }}
      >
        <div style={{ position: 'relative', zIndex: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '16px' }}>
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '6px' }}>
              <span style={{ background: 'rgba(255,255,255,0.15)', backdropFilter: 'blur(8px)', color: '#a5b4fc', fontSize: '11px', fontWeight: 800, padding: '4px 12px', borderRadius: '20px', letterSpacing: '1px', textTransform: 'uppercase' }}>
                🏆 {activeContest?.nome || 'Concurso Alvo'} · {activeContest?.banca || 'Banca'}
              </span>
              {activeContest?.preEdital && (
                <span style={{ background: '#eab308', color: '#000', fontSize: '10px', fontWeight: 800, padding: '4px 8px', borderRadius: '12px' }}>
                  Pré-Edital
                </span>
              )}
            </div>
            <h2 style={{ fontSize: '24px', fontWeight: 900, margin: '4px 0 8px', color: '#ffffff', letterSpacing: '-0.5px' }}>
              {dataProvaStr ? `Contagem Regressiva para a Prova` : `Data da Prova não definida`}
            </h2>
            <p style={{ margin: 0, fontSize: '13px', color: '#c7d2fe' }}>
              {dataProvaStr
                ? `Data agendada: ${new Date(`${dataProvaStr}T12:00:00`).toLocaleDateString('pt-BR', { day: '2-digit', month: 'long', year: 'numeric' })}`
                : `Defina a data da prova para ativar a contagem regressiva.`}
            </p>
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
            {tempoRestante ? (
              <div style={{ display: 'flex', gap: '12px', textAlign: 'center' }}>
                <div style={{ background: 'rgba(0,0,0,0.3)', backdropFilter: 'blur(10px)', padding: '10px 16px', borderRadius: '12px', border: '1px solid rgba(255,255,255,0.15)' }}>
                  <strong style={{ fontSize: '28px', color: '#38bdf8', display: 'block', fontWeight: 900, lineHeight: 1 }}>{tempoRestante.dias}</strong>
                  <span style={{ fontSize: '10px', color: '#94a3b8', textTransform: 'uppercase', fontWeight: 700 }}>DIAS</span>
                </div>
                <div style={{ background: 'rgba(0,0,0,0.3)', backdropFilter: 'blur(10px)', padding: '10px 16px', borderRadius: '12px', border: '1px solid rgba(255,255,255,0.15)' }}>
                  <strong style={{ fontSize: '28px', color: '#818cf8', display: 'block', fontWeight: 900, lineHeight: 1 }}>{String(tempoRestante.horas).padStart(2, '0')}</strong>
                  <span style={{ fontSize: '10px', color: '#94a3b8', textTransform: 'uppercase', fontWeight: 700 }}>HORAS</span>
                </div>
                <div style={{ background: 'rgba(0,0,0,0.3)', backdropFilter: 'blur(10px)', padding: '10px 16px', borderRadius: '12px', border: '1px solid rgba(255,255,255,0.15)' }}>
                  <strong style={{ fontSize: '28px', color: '#c084fc', display: 'block', fontWeight: 900, lineHeight: 1 }}>{String(tempoRestante.minutos).padStart(2, '0')}</strong>
                  <span style={{ fontSize: '10px', color: '#94a3b8', textTransform: 'uppercase', fontWeight: 700 }}>MIN</span>
                </div>
                <div style={{ background: 'rgba(0,0,0,0.3)', backdropFilter: 'blur(10px)', padding: '10px 16px', borderRadius: '12px', border: '1px solid rgba(255,255,255,0.15)' }}>
                  <strong style={{ fontSize: '28px', color: '#f472b6', display: 'block', fontWeight: 900, lineHeight: 1 }}>{String(tempoRestante.segundos).padStart(2, '0')}</strong>
                  <span style={{ fontSize: '10px', color: '#94a3b8', textTransform: 'uppercase', fontWeight: 700 }}>SEG</span>
                </div>
              </div>
            ) : null}

            <button
              type="button"
              className="btn-secondary"
              style={{ background: 'rgba(255,255,255,0.15)', color: '#fff', border: '1px solid rgba(255,255,255,0.2)', padding: '10px 16px', fontSize: '12px', fontWeight: 700 }}
              onClick={() => {
                setNovaDataProva(dataProvaStr);
                setModalDataProva(true);
              }}
            >
              ✏️ {dataProvaStr ? 'Alterar Data' : 'Definir Data'}
            </button>
          </div>
        </div>
      </div>

      {modalDataProva && (
        <div className="modal-backdrop" onClick={() => setModalDataProva(false)}>
          <div className="modal-card small-modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-heading">
              <div>
                <h2>📅 Data da Prova do Concurso</h2>
                <p>{activeContest?.nome}</p>
              </div>
              <button className="icon-button" onClick={() => setModalDataProva(false)}>
                ×
              </button>
            </div>
            <form onSubmit={salvarDataProva}>
              <div className="form-fields">
                <label>
                  <span className="field-label">Data Oficial da Prova</span>
                  <input
                    required
                    type="date"
                    className="form-control"
                    value={novaDataProva}
                    onChange={(e) => setNovaDataProva(e.target.value)}
                  />
                </label>
              </div>
              <div className="modal-actions" style={{ marginTop: '20px' }}>
                <button type="button" className="btn-secondary" onClick={() => setModalDataProva(false)}>
                  Cancelar
                </button>
                <button type="submit" className="btn-primary" disabled={salvandoData}>
                  {salvandoData ? 'Salvando…' : 'Salvar Data'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
