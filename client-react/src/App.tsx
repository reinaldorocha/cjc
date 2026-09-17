import React from 'react';
import { HashRouter, Routes, Route, Navigate, useParams } from 'react-router-dom';
import { DataProvider } from './context/DataContext';
import { useData } from './context/DataContext';
import { Sidebar } from './components/layout/Sidebar';
import { Navbar } from './components/layout/Navbar';
import { TimerModal } from './components/timer/TimerModal';

import { ConcursosPage } from './pages/ConcursosPage';
import { DashboardPage } from './pages/DashboardPage';
import { EditalPage } from './pages/EditalPage';
import { MateriaisApoioPage } from './pages/MateriaisApoioPage';
import { CronogramaInteligentePage } from './pages/CronogramaInteligentePage';
import { SimuladosPage } from './pages/SimuladosPage';
import { RevisoesPage } from './pages/RevisoesPage';
import { FlashcardsPage } from './pages/FlashcardsPage';
import { BancoQuestoesPage } from './pages/BancoQuestoesPage';
import { CadernosPage } from './pages/CadernosPage';
import { HistoricoPage } from './pages/HistoricoPage';
import { AjudaPage } from './pages/AjudaPage';
import { CronogramaAlunoPage } from './pages/CronogramaAlunoPage';
import { AutenticacaoProvider, useAutenticacao } from './context/AutenticacaoContext';
import { WhiteLabelProvider } from './context/WhiteLabelContext';
import { VisaoAlunoProvider } from './context/VisaoAlunoContext';
import { MentoriaProvider } from './context/MentoriaContext';
import { EntrarPage } from './pages/EntrarPage';
import { MentorAlunosPage } from './pages/MentorAlunosPage';
import { MentorConteudosPage } from './pages/MentorConteudosPage';
import { MentorWhiteLabelPage } from './pages/MentorWhiteLabelPage';
import { RadarMentoriaPage } from './pages/RadarMentoriaPage';
import { AdministracaoPage } from './pages/AdministracaoPage';
import { NovoConcursoPage } from './pages/NovoConcursoPage';
import { PortalHeader } from './components/layout/PortalHeader';
import { CursosPage } from './features/cursos/CursosPage';
import { CursoDetalhePage } from './features/cursos/CursoDetalhePage';

const RequireContest: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { activeContestId, loading } = useData();
  if (loading) return <div className="loading-state">Carregando seus dados…</div>;
  return activeContestId ? (
    <>{children}</>
  ) : (
    <div className="empty-state card-base" style={{ margin: '40px auto', maxWidth: '520px', padding: '40px 24px', textAlign: 'center' }}>
      <span style={{ fontSize: '48px' }}>🎯</span>
      <h2 style={{ marginTop: '14px', fontSize: '18px' }}>Aguardando Atribuição de Concurso</h2>
      <p style={{ color: 'var(--text3)', fontSize: '12px', marginTop: '6px', lineHeight: 1.6 }}>
        Seu mentor ainda não atribuiu um concurso à sua conta. Entre em contato com seu mentor para disponibilizar seu edital e plano de estudos.
      </p>
    </div>
  );
};
const RequireCronograma: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { usuario } = useAutenticacao();
  return usuario?.permiteCronogramaInteligente ? children : <Navigate to="/dashboard" replace />;
};

const AplicacaoAluno: React.FC = () => (
  <DataProvider>
    <div className="app-layout">
      <Sidebar />
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0 }}>
        <div style={{ padding: '16px 24px 0' }}>
          <Navbar />
        </div>
        <main className="main-content">
          <Routes>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/concursos" element={<Navigate to="/dashboard" replace />} />
            <Route path="/concursos/novo" element={<NovoConcursoPage />} />
            <Route path="/concursos/:concursoId/editar" element={<NovoConcursoPage />} />
            <Route path="/dashboard" element={<RequireContest><DashboardPage /></RequireContest>} />
            <Route path="/edital" element={<RequireContest><EditalPage /></RequireContest>} />
            <Route path="/materiais-apoio" element={<RequireContest><MateriaisApoioPage /></RequireContest>} />
            <Route path="/cadernos" element={<RequireContest><CadernosPage /></RequireContest>} />
            <Route path="/cronograma" element={<RequireContest><CronogramaAlunoPage /></RequireContest>} />
            <Route path="/cronograma-inteligente" element={<RequireContest><RequireCronograma><CronogramaInteligentePage /></RequireCronograma></RequireContest>} />
            <Route path="/simulados" element={<RequireContest><SimuladosPage /></RequireContest>} />
            <Route path="/revisoes" element={<RequireContest><RevisoesPage /></RequireContest>} />
            <Route path="/flashcards" element={<RequireContest><FlashcardsPage /></RequireContest>} />
            <Route path="/questoes" element={<RequireContest><BancoQuestoesPage /></RequireContest>} />
            <Route path="/historico" element={<RequireContest><HistoricoPage /></RequireContest>} />
            <Route path="/ajuda" element={<AjudaPage />} />
            <Route path="/cursos" element={<CursosPage />} />
            <Route path="/cursos/:cursoId" element={<CursoDetalhePage />} />
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Routes>
        </main>
      </div>
      <TimerModal />
    </div>
  </DataProvider>
);

const LayoutMentorAluno: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { usuario } = useAutenticacao();
  const { alunoId } = useParams();
  const visaoAluno = usuario?.papel === 'mentor' && !!alunoId;
  return (
    <VisaoAlunoProvider ativa={visaoAluno}>
      <div className="app-layout" style={{ minHeight: '100vh' }}>
        <Sidebar />
        <main className="main-content" style={{ flex: 1 }}>
          <Navbar />
          {children}
        </main>
      </div>
    </VisaoAlunoProvider>
  );
};

const ConcursosMentor: React.FC = () => {
  const { alunoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''}>
      <LayoutMentorAluno>
        <ConcursosPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const CursosAlunoMentor: React.FC = () => {
  const { alunoId } = useParams();
  return <DataProvider alunoId={alunoId || ''}><LayoutMentorAluno><CursosPage /></LayoutMentorAluno></DataProvider>;
};
const CursoDetalheAlunoMentor: React.FC = () => {
  const { alunoId } = useParams();
  return <DataProvider alunoId={alunoId || ''}><LayoutMentorAluno><CursoDetalhePage /></LayoutMentorAluno></DataProvider>;
};
const ConcursosNovoMentor: React.FC = () => {
  const { alunoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''}>
      <LayoutMentorAluno>
        <NovoConcursoPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const EditalMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <EditalPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const CronogramaMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <CronogramaAlunoPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const CronogramaInteligenteMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <CronogramaInteligentePage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const HistoricoMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <HistoricoPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const RevisoesMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <RevisoesPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const FlashcardsMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <FlashcardsPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const DashboardMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <DashboardPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const SimuladosMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <SimuladosPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const MateriaisApoioMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <MateriaisApoioPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const BancoQuestoesMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <BancoQuestoesPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};
const CadernosMentor: React.FC = () => {
  const { alunoId, concursoId } = useParams();
  return (
    <DataProvider alunoId={alunoId || ''} concursoInicialId={concursoId || ''}>
      <LayoutMentorAluno>
        <CadernosPage />
      </LayoutMentorAluno>
    </DataProvider>
  );
};

const AplicacaoAutenticada: React.FC = () => {
  const { usuario, carregando } = useAutenticacao();
  if (carregando)
    return (
      <div className="auth-loading">
        🎯<span>Carregando sessão…</span>
      </div>
    );
  if (!usuario)
    return (
      <Routes>
        <Route path="*" element={<EntrarPage />} />
      </Routes>
    );
  if (usuario.papel === 'aluno') return <AplicacaoAluno />;
  if (usuario.papel === 'mentor')
    return (
      <DataProvider alunoId="eu">
        <MentoriaProvider>
          <Routes>
            <Route path="/mentor/cursos" element={<LayoutMentorAluno><CursosPage /></LayoutMentorAluno>} />
            <Route path="/mentor/cursos/:cursoId" element={<LayoutMentorAluno><CursoDetalhePage /></LayoutMentorAluno>} />
            <Route
              path="/mentor/alunos"
              element={
                <LayoutMentorAluno>
                  <MentorAlunosPage />
                </LayoutMentorAluno>
              }
            />
            <Route
              path="/mentor/conteudos"
              element={<Navigate to="/mentor/concursos-catalogo" replace />}
            />
            <Route
              path="/mentor/concursos-catalogo"
              element={
                <LayoutMentorAluno>
                  <MentorConteudosPage key="concursos" abaInicial="concursos" />
                </LayoutMentorAluno>
              }
            />
            <Route
              path="/mentor/flashcards-catalogo"
              element={
                <LayoutMentorAluno>
                  <MentorConteudosPage key="flashcards" abaInicial="flashcards" />
                </LayoutMentorAluno>
              }
            />
            <Route
              path="/mentor/banco-questoes"
              element={
                <LayoutMentorAluno>
                  <MentorConteudosPage key="banco_questoes" abaInicial="banco_questoes" />
                </LayoutMentorAluno>
              }
            />
            <Route
              path="/questoes"
              element={
                <LayoutMentorAluno>
                  <BancoQuestoesPage />
                </LayoutMentorAluno>
              }
            />
            <Route
              path="/mentor/materiais-apoio"
              element={
                <DataProvider alunoId="eu">
                  <LayoutMentorAluno>
                    <MateriaisApoioPage />
                  </LayoutMentorAluno>
                </DataProvider>
              }
            />
            <Route
              path="/mentor/radar"
              element={
                <DataProvider alunoId="eu">
                  <LayoutMentorAluno>
                    <RadarMentoriaPage />
                  </LayoutMentorAluno>
                </DataProvider>
              }
            />
            <Route
              path="/mentor/whitelabel"
              element={
                <DataProvider alunoId="eu">
                  <LayoutMentorAluno>
                    <MentorWhiteLabelPage />
                  </LayoutMentorAluno>
                </DataProvider>
              }
            />
            <Route path="/mentor/alunos/:alunoId/concursos" element={<ConcursosMentor />} />
            <Route path="/mentor/alunos/:alunoId/cursos" element={<CursosAlunoMentor />} />
            <Route path="/mentor/alunos/:alunoId/cursos/:cursoId" element={<CursoDetalheAlunoMentor />} />
            <Route
              path="/mentor/concursos/novo"
              element={
                <LayoutMentorAluno>
                  <NovoConcursoPage />
                </LayoutMentorAluno>
              }
            />
            <Route
              path="/mentor/concursos/:concursoId/editar"
              element={
                <LayoutMentorAluno>
                  <NovoConcursoPage />
                </LayoutMentorAluno>
              }
            />
            <Route path="/mentor/alunos/:alunoId/concursos/novo" element={<ConcursosNovoMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/editar" element={<ConcursosNovoMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/edital" element={<EditalMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/cronograma" element={<CronogramaMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/cronograma-inteligente" element={<CronogramaInteligenteMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/cronograma-inteligente" element={<CronogramaInteligenteMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/revisoes" element={<RevisoesMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/historico" element={<HistoricoMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/flashcards" element={<FlashcardsMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/questoes" element={<BancoQuestoesMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/materiais-apoio" element={<MateriaisApoioMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/dashboard" element={<DashboardMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/simulados" element={<SimuladosMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/:concursoId/cadernos" element={<CadernosMentor />} />
            <Route path="/mentor/alunos/:alunoId/concursos/cadernos" element={<CadernosMentor />} />
            <Route path="*" element={<Navigate to="/mentor/alunos" replace />} />
          </Routes>
        </MentoriaProvider>
      </DataProvider>
    );
  return (
    <>
      <PortalHeader />
      <main className="portal-content">
        <Routes>
          <Route path="/administracao" element={<AdministracaoPage />} />
          <Route path="*" element={<Navigate to="/administracao" replace />} />
        </Routes>
      </main>
    </>
  );
};

export const App: React.FC = () => (
  <HashRouter>
    <AutenticacaoProvider>
      <WhiteLabelProvider>
        <AplicacaoAutenticada />
      </WhiteLabelProvider>
    </AutenticacaoProvider>
  </HashRouter>
);

export default App;
