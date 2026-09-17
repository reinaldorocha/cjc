import { useCallback, useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAutenticacao } from '../../context/AutenticacaoContext';
import { useData } from '../../context/DataContext';
import { useVisaoAluno } from '../../context/VisaoAlunoContext';
import { Editor } from './EditorCurso';
import { cursosService, type Curso } from '../../services/cursosService';
import { Pagination } from '../../components/ui/Pagination';
import './cursos.css';
import './catalogo.css';

const mensagem = (e: unknown) => e instanceof Error ? e.message : 'NÃ£o foi possÃ­vel concluir a operaÃ§Ã£o.';
const novoCurso = (alunoId: string): Curso => ({ id: '', titulo: '', descricao: '', categoria: '', capaUrl: '', modoExibicao: 'curso', escopo: alunoId ? 'alunos' : 'global', destinatarios: alunoId ? [alunoId] : [], aulas: [] });

export function CursosPage() {
  const { usuario } = useAutenticacao();
  const visaoAluno = useVisaoAluno();
  const { alunoId: alunoEmVisao } = useData();
  const mentor = usuario?.papel === 'mentor' && !visaoAluno;
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const alunoId = params.get('alunoId') || '';
  const alunoConsulta = visaoAluno ? alunoEmVisao : 'eu';
  const [lista, setLista] = useState<Curso[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [erro, setErro] = useState('');
  const [busca, setBusca] = useState('');
  const [categoria, setCategoria] = useState('');
  const [filtroEstado, setFiltroEstado] = useState('');
  const [pagina, setPagina] = useState(1);
  const [editor, setEditor] = useState<Curso | null>(null);
  const carregar = useCallback(async () => {
    setCarregando(true); setErro('');
    try { setLista((await cursosService.listar(mentor, alunoConsulta)).cursos); } catch (e) { setErro(mensagem(e)); } finally { setCarregando(false); }
  }, [mentor, alunoConsulta]);
  useEffect(() => { void carregar(); }, [carregar]);
  const abrir = (curso: Curso) => navigate(`${visaoAluno ? `/mentor/alunos/${encodeURIComponent(alunoEmVisao)}/cursos` : mentor ? '/mentor/cursos' : '/cursos'}/${encodeURIComponent(curso.id)}${alunoId && mentor ? `?alunoId=${encodeURIComponent(alunoId)}` : ''}`);
  const cursos = lista.filter((c) => (!categoria || c.categoria === categoria) && (!filtroEstado || (mentor ? filtroEstado === 'publicados' ? c.publicado : !c.publicado : filtroEstado === 'concluidos' ? c.resumo?.percentual === 100 : (c.resumo?.percentual || 0) < 100)) && `${c.titulo} ${c.descricao} ${c.categoria}`.toLocaleLowerCase().includes(busca.toLocaleLowerCase()));
  const totalPaginas = Math.max(1, Math.ceil(cursos.length / 12));
  const paginaAtual = Math.min(pagina, totalPaginas);
  const cursosDaPagina = cursos.slice((paginaAtual - 1) * 12, paginaAtual * 12);
  const categorias = [...new Set(lista.map((c) => c.categoria))];
  return <section className="cursos-page cursos-catalogo">
    <header className="cursos-heading"><div className="cursos-marca"><span className="cursos-marca-icone" aria-hidden="true">â–¶</span><div><span className="curso-eyebrow">SUA PLATAFORMA DE APRENDIZAGEM</span><h1>{mentor ? 'Cursos da mentoria' : 'Meus cursos'}</h1></div></div><div className="curso-actions"><span className="cursos-catalogo-total">{lista.length} {lista.length === 1 ? 'curso' : 'cursos'} no catÃ¡logo</span>{mentor && <button className="curso-primary" onClick={() => setEditor(novoCurso(alunoId))}>+ Novo curso</button>}</div></header>
    {alunoId && mentor && <p className="curso-notice">Ao criar um curso, este aluno jÃ¡ estarÃ¡ selecionado.</p>}
    {erro && <div role="alert" className="curso-alerta">{erro} <button onClick={() => void carregar()}>Tentar novamente</button></div>}
    {carregando ? <p role="status" className="cursos-vazio">Carregando seus cursosâ€¦</p> : <>
      <div className="cursos-filtros"><input aria-label="Buscar cursos" placeholder="Buscar curso, assunto ou palavra-chaveâ€¦" value={busca} onChange={(e) => setBusca(e.target.value)} /><select aria-label="Filtrar categoria" value={categoria} onChange={(e) => setCategoria(e.target.value)}><option value="">Todas as categorias</option>{categorias.map((c) => <option key={c}>{c}</option>)}</select></div>
      <div className="curso-filtro-abas" role="group" aria-label="SituaÃ§Ã£o dos cursos">{(mentor ? [['', 'Todos'], ['publicados', 'Publicados'], ['rascunhos', 'Rascunhos']] : [['', 'Todos'], ['pendentes', 'A concluir'], ['concluidos', 'ConcluÃ­dos']]).map(([value, label]) => <button key={value} aria-pressed={filtroEstado === value} onClick={() => setFiltroEstado(value)}>{label}</button>)}</div>
      {!cursos.length && <div className="cursos-vazio"><span>â–·</span><h2>{lista.length ? 'Nenhum curso encontrado' : mentor ? 'Monte seu primeiro curso' : 'Seus cursos aparecerÃ£o aqui'}</h2><p>{lista.length ? 'Experimente outra busca ou categoria.' : mentor ? 'Adicione mÃ³dulos e escolha quem poderÃ¡ assistir.' : 'Seu mentor ainda nÃ£o disponibilizou cursos para vocÃª.'}</p></div>}
      {!!cursos.length && <><div className="cursos-grade" aria-label="Cursos">{cursosDaPagina.map((c) => <CursoCard key={c.id} curso={c} mentor={mentor} abrir={abrir} />)}</div><Pagination pagina={paginaAtual} totalPaginas={totalPaginas} total={cursos.length} onChange={setPagina} /></>}
    </>}
    {editor && <Editor curso={editor} alunoId={alunoId} fechar={() => { setEditor(null); void carregar(); }} salvo={() => { setEditor(null); void carregar(); }} />}
  </section>;
}

function CursoCard({ curso, mentor, abrir }: { curso: Curso; mentor: boolean; abrir: (curso: Curso) => void }) {
  const modulos = new Set(curso.aulas.map((a) => a.moduloId || a.modulo)).size;
  return <article className="curso-card">
    <button className="curso-card-capa" aria-label={`Abrir curso ${curso.titulo}`} onClick={() => abrir(curso)}>{curso.capaUrl && <img src={curso.capaUrl} alt="" loading="lazy" onError={(e) => { e.currentTarget.style.display = 'none'; }} />}<span className="curso-play">â–¶</span><span className="curso-count">{modulos} {modulos === 1 ? 'mÃ³dulo' : 'mÃ³dulos'}</span></button>
    <div className="curso-card-info"><h3>{curso.titulo}</h3><p>{curso.descricao || curso.categoria}</p>{mentor ? <><span className={`curso-status ${curso.publicado ? 'publicado' : 'rascunho'}`}>{curso.publicado ? curso.alteracoesPendentes ? 'Publicado Â· alteraÃ§Ãµes pendentes' : 'Publicado' : 'Rascunho'}</span><button className="curso-card-abrir" onClick={() => abrir(curso)}>Abrir curso</button></> : <div className="curso-card-progresso"><progress value={curso.resumo?.percentual || 0} max={100} aria-label={`Progresso de ${curso.titulo}`} /><small>{curso.resumo?.concluidas || 0}/{curso.aulas.length} aulas Â· {curso.resumo?.percentual || 0}%</small></div>}</div>
  </article>;
}
