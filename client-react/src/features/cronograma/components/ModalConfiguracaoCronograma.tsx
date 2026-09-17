import React, { useState } from 'react';
import { useData } from '../../../context/DataContext';

interface ModalConfiguracaoCronogramaProps {
  isOpen: boolean;
  onClose: () => void;
  configInicial: any;
  onGerar: (novaConfig: any) => Promise<void>;
}

const DIAS_SEMANA = [
  { key: 'seg', label: 'Segunda' },
  { key: 'ter', label: 'Terça' },
  { key: 'qua', label: 'Quarta' },
  { key: 'qui', label: 'Quinta' },
  { key: 'sex', label: 'Sexta' },
  { key: 'sab', label: 'Sábado' },
  { key: 'dom', label: 'Domingo' }
];

export const ModalConfiguracaoCronograma: React.FC<ModalConfiguracaoCronogramaProps> = ({
  isOpen,
  onClose,
  configInicial,
  onGerar
}) => {
  const { getArray } = useData();
  const materias = getArray('materias') || [];

  const [tipo, setTipo] = useState(configInicial?.tipo || 'ciclo_inteligente');
  const [horas, setHoras] = useState<Record<string, number>>({
    seg: configInicial?.horas?.seg ?? 3,
    ter: configInicial?.horas?.ter ?? 3,
    qua: configInicial?.horas?.qua ?? 3,
    qui: configInicial?.horas?.qui ?? 3,
    sex: configInicial?.horas?.sex ?? 3,
    sab: configInicial?.horas?.sab ?? 5,
    dom: configInicial?.horas?.dom ?? 4
  });

  const [minutosTopico, setMinutosTopico] = useState(configInicial?.minutosTopico || 60);
  const [maxTopicosDia, setMaxTopicosDia] = useState(configInicial?.maxTopicosDia || 4);
  const [metaAcertos, setMetaAcertos] = useState(configInicial?.metaAcertos || 80);

  // Matérias e Prioridades
  const [materiasSel, setMateriasSel] = useState<Record<string, boolean>>(() => {
    const obj: Record<string, boolean> = {};
    materias.forEach((m: any) => {
      obj[m.id] = configInicial?.materiasSelecionadas?.[m.id] !== false;
    });
    return obj;
  });

  const [materiasPrio, setMateriasPrio] = useState<Record<string, string>>(() => {
    const obj: Record<string, string> = {};
    materias.forEach((m: any) => {
      obj[m.id] = configInicial?.materiaPrioridades?.[m.id] || 'media';
    });
    return obj;
  });

  const [salvando, setSalvando] = useState(false);
  const [erro, setErro] = useState('');

  if (!isOpen) return null;

  const totalHorasSemanais = Object.values(horas).reduce((acc, v) => acc + (Number(v) || 0), 0);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErro('');
    if (totalHorasSemanais <= 0) {
      setErro('Informe pelo menos 1 hora de estudo em algum dia da semana.');
      return;
    }

    setSalvando(true);
    try {
      await onGerar({
        tipo,
        cicloModo: 'adaptativo',
        horas,
        minutosTopico: Number(minutosTopico),
        maxTopicosDia: Number(maxTopicosDia),
        metaAcertos: Number(metaAcertos),
        materiasSelecionadas: materiasSel,
        materiaPrioridades: materiasPrio
      });
      onClose();
    } catch (err: any) {
      setErro(err?.message || 'Erro ao gerar o cronograma.');
    } finally {
      setSalvando(false);
    }
  };

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <form
        className="modal-card"
        style={{ maxWidth: '720px', width: '95%' }}
        onClick={(e) => e.stopPropagation()}
        onSubmit={handleSubmit}
      >
        <div className="modal-heading">
          <div>
            <h2>⚙️ Configuração do Cronograma</h2>
            <p>Defina a sua disponibilidade e preferências para gerar o plano de estudos otimizado.</p>
          </div>
          <button type="button" className="icon-button" onClick={onClose}>×</button>
        </div>

        {erro && <div className="feedback-banner error" style={{ marginBottom: '16px' }}>{erro}</div>}

        <div className="form-fields" style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          {/* TIPO DE CRONOGRAMA */}
          <div>
            <label className="field-label" style={{ fontWeight: 700, fontSize: '13px', marginBottom: '6px' }}>
              Tipo de Cronograma
            </label>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
              <div
                style={{
                  padding: '12px 16px',
                  borderRadius: '8px',
                  border: `2px solid ${tipo === 'ciclo_inteligente' ? 'var(--accent)' : 'var(--border)'}`,
                  background: tipo === 'ciclo_inteligente' ? 'rgba(0,180,255,0.08)' : 'var(--bg-input)',
                  cursor: 'pointer'
                }}
                onClick={() => setTipo('ciclo_inteligente')}
              >
                <div style={{ fontWeight: 700, fontSize: '14px', marginBottom: '4px' }}>
                  🔄 Ciclo Adaptativo (Inteligente)
                </div>
                <small style={{ color: 'var(--text2)', fontSize: '12px' }}>
                  Avança por tópicos de acordo com seu ritmo e foca nas matérias de maior dificuldade.
                </small>
              </div>

              <div
                style={{
                  padding: '12px 16px',
                  borderRadius: '8px',
                  border: `2px solid ${tipo === 'agendado' ? 'var(--accent)' : 'var(--border)'}`,
                  background: tipo === 'agendado' ? 'rgba(0,180,255,0.08)' : 'var(--bg-input)',
                  cursor: 'pointer'
                }}
                onClick={() => setTipo('agendado')}
              >
                <div style={{ fontWeight: 700, fontSize: '14px', marginBottom: '4px' }}>
                  📅 Agendado Semanal
                </div>
                <small style={{ color: 'var(--text2)', fontSize: '12px' }}>
                  Distribui os tópicos em datas fixas do calendário semanal até a data da prova.
                </small>
              </div>
            </div>
          </div>

          {/* HORAS POR DIA DA SEMANA */}
          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
              <span className="field-label" style={{ fontWeight: 700, fontSize: '13px' }}>
                Disponibilidade Semanal de Estudos (Horas/Dia)
              </span>
              <span style={{ fontSize: '12px', fontWeight: 800, padding: '3px 10px', borderRadius: '12px', background: 'var(--accent)', color: '#fff' }}>
                Total: {totalHorasSemanais}h / semana
              </span>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(7, 1fr)', gap: '8px' }}>
              {DIAS_SEMANA.map((d) => (
                <div key={d.key} style={{ textAlign: 'center', background: 'var(--bg-input)', padding: '8px 4px', borderRadius: '8px', border: '1px solid var(--border)' }}>
                  <span style={{ fontSize: '11px', fontWeight: 700, color: 'var(--text2)', display: 'block', marginBottom: '4px' }}>
                    {d.label}
                  </span>
                  <input
                    type="number"
                    min="0"
                    max="16"
                    step="0.5"
                    className="form-control"
                    style={{ textAlign: 'center', padding: '6px 2px', fontWeight: 800, fontSize: '14px' }}
                    value={horas[d.key] ?? 0}
                    onChange={(e) => setHoras((prev) => ({ ...prev, [d.key]: Math.max(0, parseFloat(e.target.value) || 0) }))}
                  />
                  <small style={{ fontSize: '10px', color: 'var(--text3)' }}>horas</small>
                </div>
              ))}
            </div>
          </div>

          {/* METAS E DURAÇÃO */}
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '12px' }}>
            <label>
              <span className="field-label">Duração por Tópico</span>
              <select className="form-control" value={minutosTopico} onChange={(e) => setMinutosTopico(Number(e.target.value))}>
                <option value={30}>30 minutos</option>
                <option value={45}>45 minutos</option>
                <option value={60}>60 minutos (1h)</option>
                <option value={90}>90 minutos (1h30)</option>
                <option value={120}>120 minutos (2h)</option>
              </select>
            </label>

            <label>
              <span className="field-label">Máx. Tópicos / Dia</span>
              <input
                type="number"
                min="1"
                max="10"
                className="form-control"
                value={maxTopicosDia}
                onChange={(e) => setMaxTopicosDia(Number(e.target.value))}
              />
            </label>

            <label>
              <span className="field-label">Meta de Acertos (%)</span>
              <input
                type="number"
                min="50"
                max="100"
                className="form-control"
                value={metaAcertos}
                onChange={(e) => setMetaAcertos(Number(e.target.value))}
              />
            </label>
          </div>

          {/* MATÉRIAS E PRIORIDADES */}
          {materias.length > 0 && (
            <div>
              <span className="field-label" style={{ fontWeight: 700, fontSize: '13px', display: 'block', marginBottom: '8px' }}>
                Matérias do Edital & Prioridades
              </span>
              <div style={{ maxHeight: '180px', overflowY: 'auto', background: 'var(--bg-input)', padding: '10px', borderRadius: '8px', border: '1px solid var(--border)' }}>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                  {materias.map((m: any) => (
                    <div key={m.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'var(--bg-card)', padding: '8px 12px', borderRadius: '6px' }}>
                      <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer', flex: 1 }}>
                        <input
                          type="checkbox"
                          checked={materiasSel[m.id] !== false}
                          onChange={(e) => setMateriasSel((prev) => ({ ...prev, [m.id]: e.target.checked }))}
                        />
                        <span style={{ fontSize: '13px', fontWeight: 600 }}>{m.nome}</span>
                      </label>

                      {materiasSel[m.id] !== false && (
                        <select
                          className="form-control"
                          style={{ width: '130px', padding: '4px 8px', fontSize: '11px' }}
                          value={materiasPrio[m.id] || 'media'}
                          onChange={(e) => setMateriasPrio((prev) => ({ ...prev, [m.id]: e.target.value }))}
                        >
                          <option value="alta">🔴 Alta Prioridade</option>
                          <option value="media">🟡 Média Prioridade</option>
                          <option value="baixa">🟢 Baixa Prioridade</option>
                        </select>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>

        <div className="modal-actions" style={{ marginTop: '24px' }}>
          <button type="button" className="btn-secondary" onClick={onClose}>
            Cancelar
          </button>
          <button type="submit" className="btn-primary" disabled={salvando}>
            {salvando ? 'Gerando Cronograma…' : '⚡ Gerar / Atualizar Cronograma'}
          </button>
        </div>
      </form>
    </div>
  );
};
