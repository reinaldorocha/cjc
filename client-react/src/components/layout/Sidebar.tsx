import React, { useEffect, useState } from 'react';
import { Link, NavLink, useLocation, useNavigate, useParams } from 'react-router-dom';
import { useData } from '../../context/DataContext';
import { useAutenticacao } from '../../context/AutenticacaoContext';
import { useWhiteLabel } from '../../context/WhiteLabelContext';
import { api } from '../../services/api';
import { useVisaoAluno } from '../../context/VisaoAlunoContext';

export const Sidebar: React.FC = () => {
  const [visivel, setVisivel] = useState(() => localStorage.getItem('ct:sidebar_visivel') !== 'false');

  const { activeContest, activeContestId, setActiveContestId, getArray, alunoId } = useData();
  const { usuario } = useAutenticacao();
  const visaoAluno = useVisaoAluno();
  const { config: wlConfig } = useWhiteLabel();
  const navigate = useNavigate();
  const location = useLocation();
  const params = useParams();
  const isMentor = usuario?.papel === 'mentor';
  const targetAlunoId = params.alunoId || (isMentor && alunoId && alunoId !== 'eu' ? alunoId : '');

  const concursos = getArray('concursos').filter((c: any) => !c.realizado);

  const mudarConcurso = (novoId: string) => {
    setActiveContestId(novoId);
    if (isMentor && targetAlunoId) {
      if (params.concursoId) {
        const novaRota = location.pathname.replace(params.concursoId, novoId);
        navigate(novaRota);
      } else {
        navigate(`/mentor/alunos/${targetAlunoId}/concursos/${novoId}/dashboard`);
      }
    }
  };

  useEffect(() => {
    const handleToggle = () => {
      setVisivel((prev) => {
        const next = !prev;
        localStorage.setItem('ct:sidebar_visivel', String(next));
        return next;
      });
    };
    window.addEventListener('ct:toggle-sidebar', handleToggle);
    return () => window.removeEventListener('ct:toggle-sidebar', handleToggle);
  }, []);

  // Fechar menu no mobile ao mudar de página
  useEffect(() => {
    if (window.innerWidth <= 900) {
      setVisivel(false);
    }
  }, [location.pathname]);

  const [alunoInfo, setAlunoInfo] = useState<any>(null);

  useEffect(() => {
    if (isMentor && targetAlunoId) {
      api
        .obterAluno(targetAlunoId)
        .then((d) => setAlunoInfo(d.aluno))
        .catch(() => undefined);
    }
  }, [isMentor, targetAlunoId]);

  const today = new Date().toISOString().slice(0, 10);
  const reviewsDue = getArray('revisoes').filter(
    (r: any) => r.concursoId === activeContestId && !r.feito && (!r.dataProxima || r.dataProxima <= today)
  ).length;

  const contestId = activeContestId || params.concursoId || '';

  const getPath = (studentPath: string, mentorSubPath: string) => {
    if (!isMentor) return studentPath;
    const base = `/mentor/alunos/${targetAlunoId}/concursos`;
    if (mentorSubPath === 'concursos') return base;
    return contestId ? `${base}/${contestId}/${mentorSubPath}` : base;
  };

  const navItems = visaoAluno
      ? [
        { path: getPath('/dashboard', 'dashboard'), label: 'Dashboard', icon: '📊', disabled: !contestId },
        { path: `/mentor/alunos/${encodeURIComponent(targetAlunoId)}/cursos`, label: 'Meus Cursos', icon: '▶', disabled: false },
        { path: getPath('/edital', 'edital'), label: 'Edital Verticalizado', icon: '📝', disabled: !contestId },
        { path: getPath('/materiais-apoio', 'materiais-apoio'), label: 'Materiais de Apoio', icon: '📁', disabled: !contestId },
        { path: getPath('/cadernos', 'cadernos'), label: 'Cadernos & Resumos', icon: '📖', disabled: !contestId },
        { path: getPath('/cronograma', 'cronograma'), label: 'Meu Cronograma', icon: '📅', disabled: !contestId },
        { path: getPath('/revisoes', 'revisoes'), label: 'Revisões', icon: '🔄', badge: reviewsDue, disabled: !contestId },
        { path: getPath('/flashcards', 'flashcards'), label: 'Flashcards', icon: '🃏', disabled: !contestId },
        { path: getPath('/questoes', 'questoes'), label: 'Banco de Questões', icon: '❓', disabled: !contestId },
        { path: getPath('/simulados', 'simulados'), label: 'Simulados e Raio-X', icon: '🎯', disabled: !contestId },
        { path: getPath('/historico', 'historico'), label: 'Métricas', icon: '📈', disabled: !contestId }
      ]
    : isMentor
    ? targetAlunoId
      ? [
          { path: '/mentor/alunos', label: '← Voltar para Alunos', icon: '👥', disabled: false },
          { path: `/mentor/cursos?alunoId=${encodeURIComponent(targetAlunoId)}`, label: 'Disponibilizar Cursos', icon: '▶', disabled: false },
          { path: getPath('/concursos', 'concursos'), label: 'Concursos do Aluno', icon: '🏆', disabled: false },
          { path: getPath('/dashboard', 'dashboard'), label: 'Dashboard do Aluno', icon: '📊', disabled: !contestId },
          { path: getPath('/edital', 'edital'), label: 'Edital Verticalizado', icon: '📝', disabled: !contestId },
          { path: getPath('/materiais-apoio', 'materiais-apoio'), label: 'Materiais de Apoio', icon: '📁', disabled: !contestId },
          { path: getPath('/simulados', 'simulados'), label: 'Simulados e Raio-X', icon: '🎯', disabled: !contestId },
          { path: getPath('/cronograma', 'cronograma'), label: 'Cronograma do Aluno', icon: '📅', disabled: !contestId },
          { path: getPath('/revisoes', 'revisoes'), label: 'Revisões', icon: '🔄', badge: reviewsDue, disabled: !contestId },
          { path: getPath('/flashcards', 'flashcards'), label: 'Flashcards', icon: '🎴', disabled: !contestId },
          { path: getPath('/questoes', 'questoes'), label: 'Banco de Questões', icon: '❓', disabled: !contestId },
          { path: getPath('/historico', 'historico'), label: 'Métricas', icon: '📈', disabled: !contestId },
          { path: getPath('/cadernos', 'cadernos'), label: 'Cadernos & Resumos', icon: '📖', disabled: !contestId }
        ]
      : [
          { path: '/mentor/alunos', label: 'Alunos da Mentoria', icon: '👥', disabled: false },
          { path: '/mentor/radar', label: 'Radar da Mentoria', icon: '🚨', disabled: false },
          { path: '/mentor/cursos', label: 'Cursos da Mentoria', icon: '▶', disabled: false },
          { path: '/mentor/concursos-catalogo', label: 'Concursos & Editais', icon: '🏆', disabled: false },
          { path: '/mentor/flashcards-catalogo', label: 'Flashcards da Mentoria', icon: '🎴', disabled: false },
          { path: '/mentor/banco-questoes', label: 'Banco de Questões', icon: '❓', disabled: false },
          { path: '/mentor/materiais-apoio', label: 'Materiais de Apoio', icon: '📁', disabled: false },
          { path: '/mentor/whitelabel', label: 'Marca da Mentoria', icon: '⚙️', disabled: false }
        ]
    : [
      { path: '/dashboard', label: 'Dashboard', icon: '📊', disabled: false },
      { path: '/cursos', label: 'Meus Cursos', icon: '▶', disabled: false },
      { path: '/edital', label: 'Edital Verticalizado', icon: '📝', disabled: false },
      { path: '/materiais-apoio', label: 'Materiais de Apoio', icon: '📁', disabled: false },
      { path: '/cadernos', label: 'Cadernos & Resumos', icon: '📖', disabled: false },
      { path: '/cronograma', label: 'Meu Cronograma', icon: '📅', disabled: false },
      { path: '/revisoes', label: 'Revisões', icon: '🔄', badge: reviewsDue, disabled: false },
      { path: '/flashcards', label: 'Flashcards', icon: '🎴', disabled: false },
      { path: '/questoes', label: 'Banco de Questões', icon: '❓', disabled: false },
      { path: '/simulados', label: 'Simulados e Raio-X', icon: '🎯', disabled: false },
      { path: '/historico', label: 'Métricas', icon: '📈', disabled: false },
      { path: '/ajuda', label: 'Ajuda', icon: '❓', disabled: false }
    ];

  if (!visivel) return null;

  return (
    <>
      <div
        className="sidebar-backdrop"
        onClick={() => setVisivel(false)}
      />
      <aside className="sidebar">
        {/* LOGO DO MENU */}
        <div style={{ marginBottom: '16px', padding: '0 4px' }}>
          {isMentor && !visaoAluno ? (
            targetAlunoId ? (
              <div style={{ fontSize: '11px', fontWeight: 800, color: 'var(--accent)' }}>
                Aluno em Acompanhamento
              </div>
            ) : (
              <Link to="/mentor/alunos" className="sidebar-logo" style={{ marginBottom: 0, padding: 0 }}>
                {wlConfig.logoUrl ? (
                  <img src={wlConfig.logoUrl} alt="Logo" style={{ maxHeight: '28px', maxWidth: '80px', objectFit: 'contain' }} />
                ) : (
                  <span>🏆</span>
                )}
                <span>{wlConfig.nomePlataforma || 'Mentor'}</span>
              </Link>
            )
          ) : (
            <Link to={visaoAluno ? (contestId ? getPath('/dashboard', 'dashboard') : getPath('/concursos', 'concursos')) : '/dashboard'} className="sidebar-logo" style={{ marginBottom: 0, padding: 0 }}>
              {wlConfig.logoUrl ? (
                <img src={wlConfig.logoUrl} alt="Logo" style={{ maxHeight: '28px', maxWidth: '80px', objectFit: 'contain' }} />
              ) : (
                <span>🏆</span>
              )}
              <span style={{ fontSize: '1rem', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                {wlConfig.nomePlataforma || 'Track Concursos'}
              </span>
            </Link>
          )}
        </div>

        {/* DETALHES DO ALUNO E CONCURSO ATIVO */}
        {isMentor && targetAlunoId && (
          <div
            className="sidebar-contest"
            style={{
              background: 'var(--accent-soft)',
              borderColor: 'var(--accent-border-soft)',
              marginBottom: '12px'
            }}
          >
            <strong>{alunoInfo?.nome || 'Carregando...'}</strong>
            <span style={{ fontSize: '9px', color: 'var(--text3)' }}>{alunoInfo?.email || ''}</span>
          </div>
        )}

        {(activeContest || concursos.length > 0) && (
          <div className="sidebar-contest" style={{ marginBottom: '16px' }}>
            <div className="sidebar-contest-title">Concurso ativo</div>
            {concursos.length > 1 ? (
              <select
                value={activeContestId}
                onChange={(e) => mudarConcurso(e.target.value)}
                className="nav-select"
                style={{
                  width: '100%',
                  marginTop: '4px',
                  marginBottom: '6px',
                  padding: '6px 8px',
                  borderRadius: '6px',
                  background: 'var(--bg-card)',
                  color: 'var(--text)',
                  border: '1px solid var(--border)',
                  fontSize: '12px',
                  fontWeight: 700
                }}
              >
                {concursos.map((c: any) => (
                  <option key={c.id} value={c.id}>
                    {c.nome} ({c.banca || 'Sem banca'})
                  </option>
                ))}
              </select>
            ) : (
              <strong>{activeContest?.nome || 'Sem concurso ativo'}</strong>
            )}
          </div>
        )}

        {/* NAVEGAÇÃO */}
        <nav className="sidebar-nav" style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
          {navItems.map((item) =>
            item.disabled ? (
              <div
                key={item.label}
                className="nav-link"
                style={{
                  opacity: 0.4,
                  cursor: 'not-allowed',
                  padding: '10px 14px',
                }}
                title={item.label}
              >
                <span style={{ fontSize: '18px', flexShrink: 0 }}>{item.icon}</span>
                <span style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{item.label}</span>
              </div>
            ) : (
              <NavLink
                key={item.path}
                to={item.path}
                className={({ isActive }) =>
                  `nav-link ${isActive ? 'active' : ''} ${'featured' in item && item.featured ? 'featured' : ''}`
                }
                style={{
                  padding: '10px 14px',
                  position: 'relative',
                }}
              >
                <span style={{ fontSize: '18px', flexShrink: 0 }}>{item.icon}</span>
                <span style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{item.label}</span>
                {!!item.badge && <b>{item.badge}</b>}
              </NavLink>
            )
          )}
        </nav>

        {/* FOOTER */}
        <div
          className="sidebar-footer"
          style={{
            marginTop: 'auto',
            paddingTop: '14px',
            borderTop: '1px solid var(--border)',
            fontSize: '11px',
            color: 'var(--text-muted)',
          }}
        >
          v2.5.0 · Chega Junto
        </div>
      </aside>
    </>
  );
};
