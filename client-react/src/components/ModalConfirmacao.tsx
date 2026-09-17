import React from 'react';

type ModalConfirmacaoProps = {
  aberto: boolean;
  titulo?: string;
  mensagem?: string;
  textoConfirmar?: string;
  textoCancelar?: string;
  aoConfirmar: () => void;
  aoCancelar: () => void;
};

export const ModalConfirmacao: React.FC<ModalConfirmacaoProps> = ({
  aberto,
  titulo = 'Confirmar exclusão',
  mensagem = 'Tem certeza que deseja apagar?',
  textoConfirmar = 'Sim, apagar',
  textoCancelar = 'Cancelar',
  aoConfirmar,
  aoCancelar
}) => {
  if (!aberto) return null;

  return (
    <div className="modal-backdrop" style={{ zIndex: 10000 }} onClick={aoCancelar}>
      <div
        className="modal-card small-modal"
        style={{ maxWidth: '420px', padding: '24px', textAlign: 'center' }}
        onClick={(e) => e.stopPropagation()}
      >
        <div style={{ fontSize: '40px', marginBottom: '12px' }}>⚠️</div>
        <h2 style={{ fontSize: '18px', margin: '0 0 8px 0', fontWeight: 800 }}>{titulo}</h2>
        <p style={{ fontSize: '14px', color: 'var(--text2)', margin: '0 0 20px 0', lineHeight: '1.4' }}>
          {mensagem}
        </p>

        <div style={{ display: 'flex', gap: '10px', justifyContent: 'center' }}>
          <button type="button" className="btn-secondary" onClick={aoCancelar} style={{ flex: 1 }}>
            {textoCancelar}
          </button>
          <button
            type="button"
            className="btn-primary danger"
            onClick={aoConfirmar}
            style={{ flex: 1, background: 'var(--red)', borderColor: 'var(--red)', color: '#fff' }}
          >
            {textoConfirmar}
          </button>
        </div>
      </div>
    </div>
  );
};
