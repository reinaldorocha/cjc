import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useData } from '../../context/DataContext';
import { api } from '../../services/api';

const formatTime = (totalSec: number) => {
  const safe = Math.max(0, totalSec);
  const h = Math.floor(safe / 3600);
  const m = Math.floor((safe % 3600) / 60);
  const s = safe % 60;
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
};

const beep = () => {
  try {
    const AudioCtx = window.AudioContext || (window as any).webkitAudioContext;
    const ctx = new AudioCtx();
    const gain = ctx.createGain();
    const osc = ctx.createOscillator();
    osc.frequency.value = 880;
    gain.gain.setValueAtTime(0.12, ctx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.7);
    osc.connect(gain);
    gain.connect(ctx.destination);
    osc.start();
    osc.stop(ctx.currentTime + 0.7);
  } catch {
    /* áudio opcional */
  }
};

export const TimerModal: React.FC = () => {
  const { activeContestId, alunoId, getArray } = useData();
  const materias = getArray('materias').filter((m: any) => m.concursoId === activeContestId);
  const topicos = getArray('topicos');
  const subtopicos = getArray('subtopicos');

  const [isOpen, setIsOpen] = useState(false);
  const [isRunning, setIsRunning] = useState(false);
  const [seconds, setSeconds] = useState(0);
  const [mode, setMode] = useState<'livre' | 'pomodoro'>('livre');
  const [pomodoroMinutes, setPomodoroMinutes] = useState(25);
  const [cycles, setCycles] = useState(0);
  const [selectedMateriaId, setSelectedMateriaId] = useState('');
  const [selectedTopicoId, setSelectedTopicoId] = useState('');
  const [selectedSubtopicoId, setSelectedSubtopicoId] = useState('');
  const [cronogramaItemId, setCronogramaItemId] = useState('');
  const [resolvedQuestions, setResolvedQuestions] = useState('');
  const [correctAnswers, setCorrectAnswers] = useState('');
  const [notes, setNotes] = useState('');
  const [showPrompt, setShowPrompt] = useState(false);

  const completedRef = useRef(false);
  const currentTopicos = topicos.filter((t: any) => t.materiaId === selectedMateriaId);
  const currentSubs = subtopicos.filter((s: any) => s.topicoId === selectedTopicoId);
  const remaining = mode === 'pomodoro' ? pomodoroMinutes * 60 - seconds : seconds;

  useEffect(() => {
    if (!isRunning) return;
    const interval = window.setInterval(() => setSeconds((s) => s + 1), 1000);
    return () => window.clearInterval(interval);
  }, [isRunning]);

  useEffect(() => {
    if (mode === 'pomodoro' && isRunning && remaining <= 0 && !completedRef.current) {
      completedRef.current = true;
      setIsRunning(false);
      setCycles((c) => c + 1);
      beep();
      setIsOpen(true);
    }
  }, [isRunning, mode, remaining]);

  useEffect(() => {
    completedRef.current = false;
    setSeconds(0);
    setIsRunning(false);
  }, [mode, pomodoroMinutes]);

  useEffect(() => {
    const openPreset = (event: Event) => {
      const preset = (event as CustomEvent).detail || {};
      setMode('livre');
      setSelectedMateriaId(String(preset.materiaId || ''));
      setSelectedTopicoId(String(preset.topicoId || ''));
      setSelectedSubtopicoId(String(preset.subtopicoId || ''));
      setCronogramaItemId(String(preset.cronogramaItemId || ''));
      setIsOpen(true);
      if (preset.autoStart) setIsRunning(true);
    };
    window.addEventListener('ct:open-timer', openPreset);
    return () => window.removeEventListener('ct:open-timer', openPreset);
  }, []);

  const reset = () => {
    setIsRunning(false);
    setSeconds(0);
    completedRef.current = false;
  };

  const handleFinishClick = () => {
    if (seconds <= 0 || !selectedMateriaId) return;
    setIsRunning(false);
    if (selectedTopicoId || cronogramaItemId) {
      setShowPrompt(true);
    } else {
      executeFinish(false);
    }
  };

  const executeFinish = async (concluirTopico: boolean) => {
    if (seconds <= 0 || !selectedMateriaId) return;
    const instante = new Date().toISOString();
    const common = {
      concursoId: activeContestId,
      materiaId: selectedMateriaId,
      topicoId: selectedTopicoId || null,
      subtopicoId: selectedSubtopicoId || null,
      origem: 'cronometro'
    };
    await api.criarSessao(alunoId, {
      ...common,
      segundos: seconds,
      modo: mode,
      observacoes: notes,
      estudadoEm: instante
    });
    const resolved = Math.max(0, Number(resolvedQuestions) || 0);
    const correct = Math.min(resolved, Math.max(0, Number(correctAnswers) || 0));
    if (resolved) {
      await api.criarQuestoes(alunoId, {
        ...common,
        resolvidas: resolved,
        acertos: correct,
        erros: resolved - correct,
        registradoEm: instante
      });
    }

    if (concluirTopico) {
      if (cronogramaItemId) {
        await api.alterarItemCronograma(alunoId, cronogramaItemId, { situacao: 'concluido', criarSessao: false }).catch(() => undefined);
      }
      const targetId = selectedSubtopicoId || selectedTopicoId;
      const targetTipo = selectedSubtopicoId ? 'subtopico' : 'topico';
      if (targetId) {
        await api.atualizarProgressoEdital(alunoId, targetId, { tipo: targetTipo, estudado: true }).catch(() => undefined);
      }
    }

    window.dispatchEvent(new CustomEvent('ct:dados-estudo-alterados'));
    beep();
    reset();
    setResolvedQuestions('');
    setCorrectAnswers('');
    setNotes('');
    setCronogramaItemId('');
    setShowPrompt(false);
    setIsOpen(false);
  };

  const label = useMemo(
    () => (mode === 'pomodoro' ? `Pomodoro ${pomodoroMinutes}min` : 'Cronômetro livre'),
    [mode, pomodoroMinutes]
  );

  if (!activeContestId) return null;

  return (
    <>
      <div className="timer-launcher">
        <button onClick={() => setIsOpen(true)} className={isRunning ? 'running' : ''}>
          <span>{isRunning ? '●' : '⏱'} {label}</span>
          {(seconds > 0 || isRunning) && <strong>{formatTime(remaining)}</strong>}
        </button>
      </div>

      {isOpen && (
        <div className="modal-backdrop timer-backdrop" onClick={() => setIsOpen(false)}>
          <section
            className="modal-card timer-modal"
            onClick={(e) => e.stopPropagation()}
            role="dialog"
            aria-modal="true"
            aria-label="Cronômetro de estudos"
          >
            <div className="modal-heading">
              <div>
                <h2>⏱ Sessão de estudo</h2>
                <p>
                  {cycles
                    ? `${cycles} pomodoro(s) concluído(s) nesta sessão`
                    : 'Concentre-se em uma tarefa por vez.'}
                </p>
              </div>
              <button className="icon-button" onClick={() => setIsOpen(false)}>
                ×
              </button>
            </div>

            <div className="timer-mode-tabs">
              <button className={mode === 'livre' ? 'active' : ''} onClick={() => setMode('livre')}>
                Cronômetro livre
              </button>
              <button className={mode === 'pomodoro' ? 'active' : ''} onClick={() => setMode('pomodoro')}>
                Pomodoro
              </button>
            </div>

            {mode === 'pomodoro' && (
              <div className="pomodoro-options">
                {[15, 25, 45, 60].map((value) => (
                  <button
                    key={value}
                    className={pomodoroMinutes === value ? 'active' : ''}
                    onClick={() => setPomodoroMinutes(value)}
                  >
                    {value} min
                  </button>
                ))}
              </div>
            )}

            <div className={`timer-display ${isRunning ? 'running' : ''}`}>{formatTime(remaining)}</div>

            {mode === 'pomodoro' && (
              <div className="timer-progress">
                <div style={{ width: `${Math.min(100, (seconds / (pomodoroMinutes * 60)) * 100)}%` }} />
              </div>
            )}

            <div className="timer-controls">
              {isRunning ? (
                <button className="btn-secondary pause" onClick={() => setIsRunning(false)}>
                  Ⅱ Pausar
                </button>
              ) : (
                <button className="btn-primary" onClick={() => setIsRunning(true)}>
                  ▶ {seconds ? 'Continuar' : 'Iniciar'}
                </button>
              )}
              <button className="btn-secondary" onClick={reset} disabled={!seconds}>
                ↺ Zerar
              </button>
            </div>

            <div className="timer-form">
              <label>
                <span className="field-label">Matéria *</span>
                <select
                  className="form-control"
                  value={selectedMateriaId}
                  onChange={(e) => {
                    setSelectedMateriaId(e.target.value);
                    setSelectedTopicoId('');
                    setSelectedSubtopicoId('');
                  }}
                >
                  <option value="">Selecione…</option>
                  {materias.map((m: any) => (
                    <option key={m.id} value={m.id}>
                      {m.nome}
                    </option>
                  ))}
                </select>
              </label>

              <label>
                <span className="field-label">Tópico</span>
                <select
                  className="form-control"
                  value={selectedTopicoId}
                  onChange={(e) => {
                    setSelectedTopicoId(e.target.value);
                    setSelectedSubtopicoId('');
                  }}
                  disabled={!selectedMateriaId}
                >
                  <option value="">Toda a matéria</option>
                  {currentTopicos.map((t: any) => (
                    <option key={t.id} value={t.id}>
                      {t.nome}
                    </option>
                  ))}
                </select>
              </label>

              <label>
                <span className="field-label">Subtópico</span>
                <select
                  className="form-control"
                  value={selectedSubtopicoId}
                  onChange={(e) => setSelectedSubtopicoId(e.target.value)}
                  disabled={!selectedTopicoId || !currentSubs.length}
                >
                  <option value="">Nenhum</option>
                  {currentSubs.map((s: any) => (
                    <option key={s.id} value={s.id}>
                      {s.nome}
                    </option>
                  ))}
                </select>
              </label>

              <label>
                <span className="field-label">Questões</span>
                <input
                  className="form-control"
                  type="number"
                  min="0"
                  value={resolvedQuestions}
                  onChange={(e) => setResolvedQuestions(e.target.value)}
                />
              </label>

              <label>
                <span className="field-label">Acertos</span>
                <input
                  className="form-control"
                  type="number"
                  min="0"
                  max={resolvedQuestions || undefined}
                  value={correctAnswers}
                  onChange={(e) => setCorrectAnswers(e.target.value)}
                />
              </label>

              <label className="timer-notes">
                <span className="field-label">Anotações</span>
                <input
                  className="form-control"
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  placeholder="O que foi estudado?"
                />
              </label>
            </div>

            <button
              className="btn-primary save-timer"
              onClick={handleFinishClick}
              disabled={!seconds || !selectedMateriaId}
            >
              💾 Salvar sessão no histórico
            </button>
          </section>
        </div>
      )}

      {showPrompt && (
        <div className="modal-backdrop" onClick={() => setShowPrompt(false)} style={{ zIndex: 1000 }}>
          <div
            className="modal-card small-modal"
            onClick={(e) => e.stopPropagation()}
            style={{ padding: '22px', textAlign: 'center' }}
          >
            <h2 style={{ fontSize: '15px', marginBottom: '8px' }}>O tópico foi concluído?</h2>
            <p style={{ color: 'var(--text3)', fontSize: '11.5px', marginBottom: '20px', lineHeight: 1.5 }}>
              Deseja atualizar o status deste tópico para <strong>Concluído</strong> no seu cronograma e edital?
            </p>
            <div style={{ display: 'flex', gap: '10px', justifyContent: 'center' }}>
              <button
                type="button"
                className="btn-secondary"
                onClick={() => executeFinish(false)}
                style={{ padding: '9px 15px', fontSize: '11.5px' }}
              >
                Manter pendente
              </button>
              <button
                type="button"
                className="btn-primary"
                onClick={() => executeFinish(true)}
                style={{ padding: '9px 16px', fontSize: '11.5px' }}
              >
                ✓ Sim, concluir tópico
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};
