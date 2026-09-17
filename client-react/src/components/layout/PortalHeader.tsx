import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAutenticacao } from '../../context/AutenticacaoContext';
import { ContaModal } from '../profile/ContaModal';

export const PortalHeader: React.FC = () => {
  const { usuario } = useAutenticacao();
  const [conta, setConta] = useState(false);

  return (
    <>
      <header className="portal-header">
        <Link to={usuario?.papel === 'mestre' ? '/administracao' : '/mentor/alunos'} className="portal-logo">
          🎯 <strong>Chega Junto Concurseiro</strong>
        </Link>
        <div>
          <span>{usuario?.papel === 'mestre' ? 'Painel Mestre' : 'Área do Mentor'}</span>
          <button className="profile-button" onClick={() => setConta(true)}>
            👤 {usuario?.nome}
          </button>
        </div>
      </header>
      {conta && <ContaModal onClose={() => setConta(false)} />}
    </>
  );
};
