import React, { useEffect, useState } from 'react';
import { useData } from '../../context/DataContext';
import { ContaModal } from '../profile/ContaModal';
import { useAutenticacao } from '../../context/AutenticacaoContext';
import { useNavigate, useParams, useLocation } from 'react-router-dom';
import { useVisaoAluno } from '../../context/VisaoAlunoContext';

export const Navbar: React.FC = () => {
  const { activeContestId, setActiveContestId, getArray } = useData();
  const { usuario } = useAutenticacao();
  const visaoAluno = useVisaoAluno();
  const navigate = useNavigate();
  const params = useParams();
  const location = useLocation();
  const concursos = getArray('concursos').filter((c) => !c.realizado);
  const [showProfiles, setShowProfiles] = useState(false);
  const [theme, setTheme] = useState<'dark' | 'light'>(() =>
    matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark'
  );

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
  }, [theme]);

  const toggleTheme = () => {
    const next = theme === 'dark' ? 'light' : 'dark';
    setTheme(next);
  };

  const { alunoId } = useData();

  const mudarConcurso = (novoId: string) => {
    setActiveContestId(novoId);
    const targetAluno = params.alunoId || (usuario?.papel === 'mentor' && alunoId && alunoId !== 'eu' ? alunoId : '');
    if (usuario?.papel === 'mentor' && targetAluno) {
      if (params.concursoId) {
        const novaRota = location.pathname.replace(params.concursoId, novoId);
        navigate(novaRota);
      } else {
        navigate(`/mentor/alunos/${targetAluno}/concursos/${novoId}/dashboard`);
      }
    }
  };

  return (
    <>
      <header className="header-bar">
        <div className="contest-switcher" style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
          <button
            className="menu-toggle"
            onClick={() => window.dispatchEvent(new Event('ct:toggle-sidebar'))}
            aria-label="Alternar Menu"
            title="Alternar Menu"
          >
            ☰
          </button>
          <span className="sync-dot online" />
          <span className="sync-label">Sessão segura</span>
          {usuario?.papel === 'mentor' && !visaoAluno && (
            <span
              style={{
                fontSize: '11px',
                fontWeight: 700,
                padding: '2px 8px',
                borderRadius: '6px',
                background: 'var(--accent-light)',
                color: 'var(--accent-primary)',
                border: '1px solid var(--border-active)',
                display: 'inline-flex',
                alignItems: 'center',
                gap: '4px'
              }}
            >
              <span>🎯</span> Área do Mentor
            </span>
          )}
          {concursos.length > 0 && (
            <select
              value={activeContestId}
              onChange={(e) => mudarConcurso(e.target.value)}
              className="nav-select"
              aria-label="Concurso ativo"
            >
              {concursos.map((c: any) => (
                <option key={c.id} value={c.id}>
                  {c.nome} ({c.banca || 'Sem banca'})
                </option>
              ))}
            </select>
          )}
        </div>
        <div className="nav-actions">
          <button onClick={toggleTheme} className="nav-button" aria-label="Alternar tema">
            {theme === 'dark' ? '☾' : '☀'} <span>{theme === 'dark' ? 'Escuro' : 'Claro'}</span>
          </button>
          <button onClick={() => setShowProfiles(true)} className="profile-button">
            👤 {usuario?.nome}
          </button>
        </div>
      </header>
      {showProfiles && <ContaModal onClose={() => setShowProfiles(false)} />}
    </>
  );
};
