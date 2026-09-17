import React, { useState } from 'react';
import { api } from '../services/api';

interface ModalConclusaoAtividadeProps {
  item: any;
  alunoId: string;
  activeContestId: string;
  materias: any[];
  aoFechar: () => void;
  aoConcluirSucesso?: () => void;
}

export const ModalConclusaoAtividade: React.FC<ModalConclusaoAtividadeProps> = ({
  item,
  alunoId,
  activeContestId,
  materias,
  aoFechar,
  aoConcluirSucesso
}) => {
  const [modalMinutos, setModalMinutos] = useState(item?.duracaoMinutos || (item?.tipoItem === 'revisao' ? 30 : 60));
  const [modalQuestoes, setModalQuestoes] = useState('');
  const [modalAcertos, setModalAcertos] = useState('');
  const [modalNotas, setModalNotas] = useState('');
  const [salvando, setSalvando] = useState(false);

  if (!item) return null;

  const materiaObj = materias.find((m: any) => m.id === item.materiaId);
  const materiaNome = item.materiaNome || materiaObj?.nome || 'Matéria';
  const topicoNome = item.topicoNome || item.nome || 'Tópico de Estudo';

  const getLocalYYYYMMDD = (d: Date = new Date()) => {
    const y = d.getFullYear();
    const m = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    return `${y}-${m}-${day}`;
  };

  const todayStr = getLocalYYYYMMDD();
  const rawDate = item.dataPlanejada || item.proximaData || item.data || '';
  const itemData = rawDate.slice(0, 10);
  const ehOutroDia = Boolean(itemData && itemData !== todayStr);

  const confirmarConclusao = async (marcarConcluido: boolean) => {
    if (!modalMinutos) return;
    setSalvando(true);
    try {
      const segundos = Math.max(1, modalMinutos) * 60;
      const instante = new Date().toISOString();
      const commonPayload = {
        concursoId: activeContestId,
        materiaId: item.materiaId || null,
        topicoId: item.topicoId || null,
        subtopicoId: item.subtopicoId || null,
        origem: item.tipoItem === 'revisao' ? 'revisao' : 'cronograma'
      };

      // 1. Registrar sessão de estudo
      await api.criarSessao(alunoId, {
        ...commonPayload,
        segundos,
        modo: 'manual',
        observacoes: modalNotas || 'Estudo registrado via modal',
        estudadoEm: instante
      }).catch(() => undefined);

      // 2. Registrar questões resolvidas se informado
      const resolved = Math.max(0, Number(modalQuestoes) || 0);
      const correct = Math.min(resolved, Math.max(0, Number(modalAcertos) || 0));
      if (resolved > 0) {
        await api.criarQuestoes(alunoId, {
          ...commonPayload,
          resolvidas: resolved,
          acertos: correct,
          erros: resolved - correct,
          registradoEm: instante
        }).catch(() => undefined);
      }

      const hojeObj = new Date();
      const amanhaObj = new Date(hojeObj.getFullYear(), hojeObj.getMonth(), hojeObj.getDate() + 1);
      const amanhaStr = getLocalYYYYMMDD(amanhaObj);

      // 3. Atualizar situação de acordo com a escolha do aluno
      if (marcarConcluido) {
        // Conclusão total do tópico/atividade
        const dadosRev: any = { concluida: true };
        const dadosItem: any = { situacao: 'concluido', criarSessao: false };
        if (ehOutroDia) {
          dadosRev.proximaData = todayStr;
          dadosItem.dataPlanejada = todayStr;
        }

        if (item.tipoItem === 'revisao' && item.id) {
          await api.alterarRevisao(alunoId, item.id, dadosRev).catch(() => undefined);
        } else if (item.id) {
          await api.alterarItemCronograma(alunoId, item.id, dadosItem).catch(() => undefined);
        }

        const targetId = item.subtopicoId || item.topicoId;
        const targetTipo = item.subtopicoId ? 'subtopico' : 'topico';
        if (targetId && item.tipoItem !== 'revisao') {
          await api.atualizarProgressoEdital(alunoId, targetId, { tipo: targetTipo, estudado: true });
        }

        // Reservar o tempo das revisões e redistribuir assuntos futuros, preservando hoje.
        await api.reprogramarCronograma(alunoId, { aPartirDe: amanhaStr, concursoId: activeContestId }).catch((erro: Error) => window.alert(`Estudo salvo, mas o cronograma não foi reorganizado: ${erro.message}. Use Reprogramar Pendentes.`));
      } else {
        // Estudo Parcial / Continuar no próximo dia útil:
        if (item.tipoItem === 'revisao' && item.id) {
          await api.alterarRevisao(alunoId, item.id, { proximaData: amanhaStr, concluida: false }).catch(() => undefined);
          await api.reprogramarCronograma(alunoId, { aPartirDe: amanhaStr, concursoId: activeContestId }).catch((erro: Error) => window.alert(`Estudo salvo, mas o cronograma não foi reorganizado: ${erro.message}. Use Reprogramar Pendentes.`));
        } else if (item.id) {
          await api.alterarItemCronograma(alunoId, item.id, {
            situacao: 'concluido',
            criarSessao: false,
            ...(ehOutroDia ? { dataPlanejada: todayStr } : {})
          }).catch(() => undefined);

          await api.criarItemCronograma(alunoId, {
            materiaId: item.materiaId || null,
            topicoId: item.topicoId || null,
            subtopicoId: item.subtopicoId || null,
            duracaoMinutos: item.duracaoMinutos || 60,
            prioridade: item.prioridade || 0,
            dataPlanejada: amanhaStr,
            situacao: 'pendente'
          }).catch(() => undefined);

          await api.reprogramarCronograma(alunoId, { aPartirDe: amanhaStr, concursoId: activeContestId }).catch((erro: Error) => window.alert(`Estudo salvo, mas o cronograma não foi reorganizado: ${erro.message}. Use Reprogramar Pendentes.`));
        }
      }

      // 4. Notificar a aplicação para recarregar dados
      window.dispatchEvent(new CustomEvent('ct:dados-estudo-alterados'));

      if (aoConcluirSucesso) aoConcluirSucesso();
      aoFechar();
    } catch (err: any) {
      alert(err.message || 'Erro ao registrar conclusão da atividade');
    } finally {
      setSalvando(false);
    }
  };

  return (
    <div className="modal-backdrop" onClick={aoFechar}>
      <div className="modal-card small-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-heading">
          <div>
            <h2>✓ Registrar Estudo</h2>
            <p>Informe o tempo e as questões praticadas nesta atividade.</p>
          </div>
          <button className="icon-button" onClick={aoFechar}>
            ×
          </button>
        </div>

        {ehOutroDia && (
          <div style={{ background: 'rgba(79, 142, 247, 0.12)', borderLeft: '4px solid var(--accent)', padding: '10px 14px', borderRadius: '8px', marginTop: '12px', fontSize: '12px', color: 'var(--text1)', textAlign: 'left' }}>
            <strong>📅 Atividade de outro dia:</strong> Esta tarefa estava agendada para <b>{new Date(`${itemData}T12:00:00`).toLocaleDateString('pt-BR')}</b>. Ao concluí-la, ela será <b>puxada para o dia de hoje ({new Date().toLocaleDateString('pt-BR')})</b>.
          </div>
        )}

        <div className="timer-form" style={{ gridTemplateColumns: '1fr 1fr', borderTop: 'none', paddingTop: '10px' }}>
          <label style={{ gridColumn: 'span 2' }}>
            <span className="field-label">Matéria & Tópico</span>
            <input
              className="form-control"
              disabled
              value={`${materiaNome} — ${topicoNome}`}
            />
          </label>

          <label style={{ gridColumn: 'span 2' }}>
            <span className="field-label">Tempo Estudado (minutos) *</span>
            <input
              className="form-control"
              type="number"
              min="1"
              value={modalMinutos}
              onChange={(e) => setModalMinutos(Math.max(1, Number(e.target.value) || 0))}
            />
          </label>

          <label>
            <span className="field-label">Questões Resolvidas</span>
            <input
              className="form-control"
              type="number"
              min="0"
              value={modalQuestoes}
              placeholder="0"
              onChange={(e) => setModalQuestoes(e.target.value)}
            />
          </label>

          <label>
            <span className="field-label">Acertos</span>
            <input
              className="form-control"
              type="number"
              min="0"
              max={modalQuestoes || undefined}
              placeholder="0"
              value={modalAcertos}
              onChange={(e) => setModalAcertos(e.target.value)}
            />
          </label>

          <label className="timer-notes" style={{ gridColumn: 'span 2' }}>
            <span className="field-label">Anotações / Observações</span>
            <textarea
              className="form-control"
              rows={2}
              value={modalNotas}
              placeholder="Observações sobre o estudo (opcional)…"
              onChange={(e) => setModalNotas(e.target.value)}
            />
          </label>
        </div>

        <div style={{ marginTop: '16px', textAlign: 'center' }}>
          <p style={{ color: 'var(--text3)', fontSize: '11.5px', marginBottom: '12px' }}>
            O estudo deste assunto foi totalmente finalizado?
          </p>
          <div style={{ display: 'flex', gap: '8px', justifyContent: 'center' }}>
            <button
              type="button"
              className="btn-secondary"
              disabled={salvando || !modalMinutos}
              onClick={() => confirmarConclusao(false)}
              style={{ padding: '9px 12px', fontSize: '11px', background: 'rgba(234, 179, 8, 0.15)', color: '#facc15', borderColor: 'rgba(234, 179, 8, 0.3)' }}
              title="Registra a meta de hoje e reinsere o tópico para continuar no próximo dia útil"
            >
              ⏳ Estudo Parcial (Adiar restante)
            </button>
            <button
              type="button"
              className="btn-primary"
              disabled={salvando || !modalMinutos}
              onClick={() => confirmarConclusao(true)}
              style={{ padding: '9px 14px', fontSize: '11px' }}
            >
              {salvando ? 'Salvando…' : '✓ Concluir Tópico'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
