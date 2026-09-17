import { useCallback, useEffect, useState } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { useAutenticacao } from '../../context/AutenticacaoContext';
import { useData } from '../../context/DataContext';
import { useVisaoAluno } from '../../context/VisaoAlunoContext';
import { cursosService, type Curso } from '../../services/cursosService';
import { AcompanhamentoCurso } from './AcompanhamentoCurso';
import { Janela } from './CursoDialog';
import { Editor } from './EditorCurso';
import { SalaCurso } from './SalaCurso';
import { TrilhoCursos } from './CatalogoVisual';
import { agruparModulos, primeiraAulaDoModulo, resumoModulo } from './modulos';
import './cursos.css';
import './catalogo.css';

const mensagem = (e: unknown) => e instanceof Error ? e.message : 'Não foi possível concluir a operação.';

export function CursoDetalhePage() {
  const { usuario } = useAutenticacao();
  const visaoAluno = useVisaoAluno();
  const { alunoId: alunoEmVisao } = useData();
  const mentor = usuario?.papel === 'mentor' && !visaoAluno;
  const { cursoId = '' } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const alunoId = new URLSearchParams(location.search).get('alunoId') || '';
  const alunoConsulta = visaoAluno ? alunoEmVisao : 'eu';
  const [curso, setCurso] = useState<Curso | null>(null);
  const [carregando, setCarregando] = useState(true);
  const [erro, setErro] = useState('');
  const [ocupado, setOcupado] = useState(false);
  const [editor, setEditor] = useState<Curso | null>(null);
  const [acompanhamento, setAcompanhamento] = useState(false);
  const [publicacao, setPublicacao] = useState<boolean | null>(null);
  const [excluir, setExcluir] = useState(false);
  const [assistindo, setAssistindo] = useState<{ aulaId?: string; moduloId?: string } | null>(null);
  const catalogo = visaoAluno ? `/mentor/alunos/${encodeURIComponent(alunoEmVisao)}/cursos` : `${mentor ? '/mentor/cursos' : '/cursos'}${mentor && alunoId ? `?alunoId=${encodeURIComponent(alunoId)}` : ''}`;

  const carregar = useCallback(async () => {
    if (!cursoId) return;
    setCarregando(true); setErro('');
    try { setCurso((mentor ? await cursosService.previa(cursoId) : await cursosService.obter(cursoId, false, alunoConsulta)).curso); }
    catch (e) { setErro(mensagem(e)); }
    finally { setCarregando(false); }
  }, [cursoId, mentor, alunoConsulta]);

  useEffect(() => { void carregar(); }, [carregar]);
  const abrirSala = (aulaId?: string, moduloId?: string) => setAssistindo({ aulaId, moduloId });
  const editar = async () => {
    if (!curso) return;
    setOcupado(true); setErro('');
    try { setEditor((await cursosService.obter(curso.id, true)).curso); }
    catch (e) { setErro(mensagem(e)); }
    finally { setOcupado(false); }
  };
  const alterarPublicacao = async () => {
    if (!curso || publicacao === null) return;
    setOcupado(true); setErro('');
    try { setCurso((await cursosService.publicar(curso, publicacao)).curso); setPublicacao(null); }
    catch (e) { setErro(mensagem(e)); }
    finally { setOcupado(false); }
  };
  const remover = async () => {
    if (!curso) return;
    setOcupado(true); setErro('');
    try { await cursosService.excluir(curso.id); navigate(catalogo); }
    catch (e) { setErro(mensagem(e)); setOcupado(false); }
  };

  if (carregando) return <section className="curso-detalhe"><p className="cursos-vazio" role="status">Carregando curso…</p></section>;
  if (!curso) return <section className="curso-detalhe"><button className="curso-detalhe-voltar" onClick={() => navigate(catalogo)}>← Voltar para cursos</button><div className="cursos-vazio"><h1>Curso indisponível</h1><p>{erro || 'Este curso não existe ou não está disponível para você.'}</p></div></section>;

  const modulos = agruparModulos(curso.aulas);
  const retomar = curso.resumo?.retomarAulaId;
  return <section className="curso-detalhe">
    <button className="curso-detalhe-voltar" onClick={() => navigate(catalogo)}>← Todos os cursos</button>
    {erro && <p className="curso-alerta" role="alert">{erro} <button onClick={() => void carregar()}>Tentar novamente</button></p>}
    <header className="curso-detalhe-banner">
      {curso.capaUrl && <img src={curso.capaUrl} alt="" onError={(e) => { e.currentTarget.style.display = 'none'; }} />}
      <div className="curso-detalhe-banner-conteudo">
        <span className="curso-eyebrow">{mentor ? 'GESTÃO DO CURSO' : 'CURSO DISPONÍVEL PARA VOCÊ'}</span>
        <p className="curso-detalhe-categoria">{curso.categoria || 'Curso'} <span>•</span> {modulos.length} {modulos.length === 1 ? 'módulo' : 'módulos'}</p>
        <h1>{curso.titulo}</h1>
        <p className="curso-detalhe-descricao">{curso.descricao || 'Aulas e materiais organizados para você estudar no seu ritmo.'}</p>
        <div className="curso-detalhe-banner-acoes">
          <button className="curso-detalhe-assistir" disabled={ocupado || !curso.aulas.length} onClick={() => abrirSala(retomar)}>{mentor ? '▶ Prévia do curso' : retomar ? '▶ Continuar assistindo' : '▶ Começar curso'}</button>
          {!mentor && <span className="curso-detalhe-progresso">{curso.resumo?.concluidas || 0}/{curso.aulas.length} aulas concluídas · {curso.resumo?.percentual || 0}%</span>}
        </div>
      </div>
    </header>
    {mentor && <div className="curso-detalhe-gestao">
      <span className={`curso-status ${curso.publicado ? 'publicado' : 'rascunho'}`}>{curso.publicado ? curso.alteracoesPendentes ? 'Publicado · alterações pendentes' : 'Publicado' : 'Rascunho'}</span>
      <div className="curso-actions"><button disabled={ocupado} onClick={() => void editar()}>Editar curso</button><button disabled={ocupado} onClick={() => abrirSala()}>Prévia</button><button disabled={ocupado} onClick={() => setAcompanhamento(true)}>Acompanhar alunos</button><button className="curso-primary" disabled={ocupado} onClick={() => setPublicacao(!curso.publicado)}>{curso.publicado ? 'Desativar curso' : 'Publicar curso'}</button><button disabled={ocupado} onClick={() => setExcluir(true)}>Excluir</button></div>
    </div>}
    <TrilhoCursos titulo="Módulos do curso">
      {modulos.map((modulo, indice) => {
        const resumo = resumoModulo(modulo.aulas);
        const aulaId = primeiraAulaDoModulo(curso.aulas, modulo.id || modulo.titulo);
        return <article key={modulo.id || modulo.titulo} className={`curso-modulo-card${resumo.bloqueado ? ' bloqueado' : ''}`}>
          <button className="curso-modulo-card-capa" disabled={resumo.bloqueado || ocupado || !aulaId} onClick={() => abrirSala(aulaId, modulo.id || modulo.titulo)} aria-label={`Abrir módulo ${modulo.titulo}`}>
            {modulo.capaUrl && <img src={modulo.capaUrl} alt="" loading="lazy" onError={(e) => { e.currentTarget.style.display = 'none'; }} />}
            <span className="curso-modulo-card-numero">MÓDULO {String(indice + 1).padStart(2, '0')}</span><span className="curso-play">{resumo.bloqueado ? '🔒' : '▶'}</span>
          </button>
          <div className="curso-modulo-card-info"><span>{resumo.total} {resumo.total === 1 ? 'aula' : 'aulas'} · {resumo.concluidas} concluídas</span><h2>{modulo.titulo}</h2><p>{resumo.bloqueado ? resumo.motivoBloqueio : modulo.descricao || 'Acesse as aulas deste módulo.'}</p>{!resumo.bloqueado && <progress value={resumo.percentual} max={100} aria-label={`Progresso do módulo ${modulo.titulo}`} />}</div>
        </article>;
      })}
    </TrilhoCursos>
    {!modulos.length && <div className="cursos-vazio"><h2>Este curso ainda não tem módulos</h2><p>{mentor ? 'Edite o curso para criar o primeiro módulo.' : 'Seu mentor ainda está preparando este conteúdo.'}</p></div>}
    {assistindo && <SalaCurso inicial={curso} previa={mentor || visaoAluno} aulaInicialId={assistindo.aulaId} moduloInicialId={assistindo.moduloId} fechar={() => { setAssistindo(null); void carregar(); }} />}
    {editor && <Editor curso={editor} alunoId={alunoId} fechar={() => { setEditor(null); void carregar(); }} salvo={() => { setEditor(null); void carregar(); }} />}
    {acompanhamento && <AcompanhamentoCurso curso={curso} fechar={() => setAcompanhamento(false)} />}
    {publicacao !== null && <Janela titulo={publicacao ? 'Publicar curso' : 'Desativar curso'} fechar={() => setPublicacao(null)}><p>{publicacao ? 'Os alunos selecionados poderão acessar a versão publicada deste curso.' : 'Os alunos deixarão de acessar este curso. O conteúdo será preservado.'}</p><div className="curso-actions"><button disabled={ocupado} onClick={() => setPublicacao(null)}>Cancelar</button><button className="curso-primary" disabled={ocupado} onClick={() => void alterarPublicacao()}>{publicacao ? 'Publicar agora' : 'Desativar'}</button></div></Janela>}
    {excluir && <Janela titulo="Excluir curso" fechar={() => setExcluir(false)}><p>Excluir “{curso.titulo}”? Esta ação não pode ser desfeita.</p><div className="curso-actions"><button disabled={ocupado} onClick={() => setExcluir(false)}>Cancelar</button><button className="curso-danger" disabled={ocupado} onClick={() => void remover()}>Excluir curso</button></div></Janela>}
  </section>;
}
