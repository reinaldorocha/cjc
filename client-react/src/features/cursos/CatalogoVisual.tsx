import { useId, useRef, useState, type ReactNode } from 'react';
import type { Curso } from '../../services/cursosService';

export function BannerCursos({ cursos, mentor, ocupado, abrir }: { cursos: Curso[]; mentor: boolean; ocupado: boolean; abrir: (curso: Curso) => void }) {
  const [selecionado, setSelecionado] = useState('');
  const [detalhes, setDetalhes] = useState(false);
  const destaques = cursos.slice(0, 5);
  const curso = destaques.find((c) => c.id === selecionado) || destaques[0];
  if (!curso) return null;
  const modulos = new Set(curso.aulas.map((a) => a.moduloId || a.modulo)).size;
  return <section className="curso-banner" aria-label="Cursos em destaque">
    {curso.capaUrl && <img key={curso.capaUrl} className="curso-banner-imagem" src={curso.capaUrl} alt="" fetchPriority="high" onError={(e) => { e.currentTarget.style.display = 'none'; }} />}
    <div className="curso-banner-conteudo">
      <span className="curso-banner-selo"><span aria-hidden="true">▶</span> A SUA PRÓXIMA CONQUISTA</span>
      <p className="curso-banner-categoria">{curso.categoria} <span>•</span> Em destaque</p>
      <h2>{curso.titulo}</h2>
      <div className="curso-banner-meta"><strong>{modulos} {modulos === 1 ? 'módulo' : 'módulos'}</strong><span>{curso.aulas.length} aulas</span><span className="curso-banner-tag">{mentor ? curso.publicado ? 'Publicado' : 'Rascunho' : 'Incluso na sua mentoria'}</span></div>
      <p className="curso-banner-descricao">{curso.descricao || 'Dê o próximo passo na sua preparação. Aulas e materiais organizados para você estudar no seu ritmo.'}</p>
      <div className="curso-banner-acoes"><button className="curso-banner-assistir" disabled={ocupado} onClick={() => abrir(curso)}><span aria-hidden="true">▶</span>{mentor ? 'Prévia do curso' : curso.resumo?.ultimaAtividade ? 'Continuar assistindo' : 'Começar a assistir'}</button><button className="curso-banner-info" aria-expanded={detalhes} onClick={() => setDetalhes(!detalhes)}><span aria-hidden="true">ⓘ</span>Mais informações</button></div>
      {detalhes && <div className="curso-banner-detalhes">{curso.descricao && <p>{curso.descricao}</p>}<p>{curso.aulas.filter((a) => a.pdfId).length} aulas com PDF · {curso.aulas.filter((a) => a.questoes?.length).length} aulas com exercícios</p><p>{mentor ? 'Use a gestão do curso abaixo para editar módulos, materiais e acesso.' : `${curso.resumo?.concluidas || 0} de ${curso.aulas.length} aulas concluídas. Materiais de módulos bloqueados aparecem após a liberação.`}</p></div>}
    </div>
    {destaques.length > 1 && <div className="curso-banner-seletor" role="group" aria-label="Selecionar destaque">{destaques.map((c, i) => <button key={c.id} aria-label={`Destaque ${i + 1}: ${c.titulo}`} aria-pressed={c.id === curso.id} onClick={() => { setSelecionado(c.id); setDetalhes(false); }}><span /></button>)}</div>}
    <span className="curso-banner-rodape">APRENDA NO SEU RITMO</span>
  </section>;
}

export function TrilhoCursos({ titulo, children }: { titulo: string; children: ReactNode }) {
  const ref = useRef<HTMLDivElement>(null);
  const id = useId();
  const mover = (direcao: number) => ref.current?.scrollBy({ left: direcao * ref.current.clientWidth * .85, behavior: window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'auto' : 'smooth' });
  return <section className="curso-colecao"><header className="curso-trilho-heading"><h2 id={id}>{titulo}</h2><div><button aria-label={`Voltar em ${titulo}`} onClick={() => mover(-1)}>‹</button><button aria-label={`Avançar em ${titulo}`} onClick={() => mover(1)}>›</button></div></header><div ref={ref} className="cursos-trilho" role="region" aria-labelledby={id} tabIndex={0}>{children}</div></section>;
}
