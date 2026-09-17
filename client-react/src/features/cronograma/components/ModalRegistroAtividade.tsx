import React, { useState } from 'react';
import { Modal } from '../../../components/ui/Modal';
import { Button } from '../../../components/ui/Button';

interface ModalRegistroAtividadeProps {
  isOpen: boolean;
  onClose: () => void;
  item: any;
  onSalvar: (dados: { tempoMinutos: number; questoesResolvidas: number; acertos: number }) => void;
}

export const ModalRegistroAtividade: React.FC<ModalRegistroAtividadeProps> = ({
  isOpen,
  onClose,
  item,
  onSalvar,
}) => {
  const [tempo, setTempo] = useState(60);
  const [questoes, setQuestoes] = useState(10);
  const [acertos, setAcertos] = useState(8);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSalvar({ tempoMinutos: tempo, questoesResolvidas: questoes, acertos });
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Registrar Conclusão: ${item?.materiaNome || 'Estudo'}`}>
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <div>
          <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
            Tempo Estudado (minutos)
          </label>
          <input
            type="number"
            value={tempo}
            onChange={(e) => setTempo(Number(e.target.value))}
            style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
          />
        </div>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Questões Resolvidas
            </label>
            <input
              type="number"
              value={questoes}
              onChange={(e) => setQuestoes(Number(e.target.value))}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>
          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Acertos
            </label>
            <input
              type="number"
              value={acertos}
              onChange={(e) => setAcertos(Number(e.target.value))}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>
        </div>
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px', marginTop: '12px' }}>
          <Button variant="ghost" onClick={onClose} type="button">Cancelar</Button>
          <Button variant="primary" type="submit">Salvar Registro</Button>
        </div>
      </form>
    </Modal>
  );
};
