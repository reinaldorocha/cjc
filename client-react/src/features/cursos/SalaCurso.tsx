import { useRef, useState } from 'react';
import { cursosService, type Curso } from '../../services/cursosService';
import { Janela } from './CursoDialog';
import { agruparModulos } from './modulos';
import { YouTubePlayer } from './YouTubePlayer';
import { ExerciciosAula } from './ExerciciosAula';

export function SalaCurso({ inicial, previa, aulaInicialId, moduloInicialId, fechar }: { inicial: Curso; previa: boolean; aulaInicialId?: string; moduloInicialId?: string; fechar: () => void }) {
  const [curso, setCurso] = useState(inicial);
  const aulaAlvo = aulaInicialId || (moduloInicialId ? inicial.aulas.find((a) => (a.moduloId === moduloInicialId || a.modulo === moduloInicialId) && !a.progresso?.concluida)?.id || inicial.aulas.find((a) => a.moduloId === moduloInicialId || a.modulo === moduloInicialId)?.id : '') || inicial.resumo?.retomarAulaId || inicial.aulas[0]?.id || '';
  const [aulaId, setAulaId] = useState(aulaAlvo);
  const [erro, setErro] = useState('');
  const [status, setStatus] = useState('');
  const [ocupado, setOcupado] = useState(false);
  const latest = useRef<Record<string, { posicao: number; duracao: number }>>({});
  const fila = useRef<Promise<void>>(Promise.resolve());
  const ativo = curso.aulas.find((a) => a.id === aulaId);
  const resumo = curso.resumo;
  const gravar = (id: string, posicao: number, duracao: number, concluida?: boolean) => {
    if (previa) return Promise.resolve();
    latest.current[id] = { posicao, duracao };
    const op = fila.current.catch(() => undefined).then(async () => {
      setStatus('Salvando progresso…');
      const data = await cursosService.progresso(curso.id, id, { posicao, duracao, ...(concluida !== undefined ? { concluida } : {}) }, true);
      setCurso(data.curso); setErro(''); setStatus('Progresso salvo');
    });
    fila.current = op;
    void op.catch((e: Error) => { setErro(e.message); setStatus('Progresso não salvo'); });
    return op;
  };
  const salvarAtual = async () => { const p = latest.current[aulaId]; if (ativo?.id && p && !ativo.bloqueada) await gravar(ativo.id, p.posicao, p.duracao); else await fila.current; };
  const sair = async () => { setOcupado(true); try { await salvarAtual(); fechar(); } catch { /* erro já exibido */ } finally { setOcupado(false); } };
  const modulos = agruparModulos(curso.aulas);
  return <Janela titulo={curso.titulo} fechar={() => { if (!ocupado) void sair(); }}>
    {previa && <p className="curso-notice">Prévia do rascunho como aluno iniciante. Regras de liberação aplicadas; progresso e respostas não são gravados.</p>}
    <div className="curso-sala-resumo"><div><strong>{resumo?.percentual || 0}% concluído</strong><span>{resumo?.concluidas || 0} de {curso.aulas.length} aulas</span></div><progress value={resumo?.percentual || 0} max={100} aria-label="Progresso do curso" /></div>
    {erro && <p className="curso-alerta" role="alert">{erro} <button onClick={() => void salvarAtual().catch(() => undefined)}>Tentar salvar novamente</button><button onClick={fechar}>Fechar sem salvar</button></p>}
    <div className="curso-sala"><div>
      {!ativo ? <p className="cursos-vazio">Este curso ainda não possui aulas. Adicione um módulo antes de publicar.</p> : ativo.bloqueada ? <div className="cursos-vazio"><span>🔒</span><h3>{ativo.titulo}</h3><p>{ativo.motivoBloqueio}</p></div> : <>
        <YouTubePlayer key={`${ativo.id}/${ativo.videoId}`} videoId={ativo.videoId} titulo={ativo.titulo} inicio={ativo.progresso?.posicao || 0} onTempo={(posicao, duracao) => { latest.current[ativo.id!] = { posicao, duracao }; }} onSalvar={(p, d) => { void gravar(ativo.id!, p, d).catch(() => undefined); }} />
        <h3 className="curso-aula-titulo">{ativo.titulo}</h3>
        <div className="curso-actions"><button className={ativo.progresso?.concluida ? '' : 'curso-primary'} disabled={previa || ocupado} onClick={async () => { setOcupado(true); const p = latest.current[ativo.id!] || ativo.progresso || { posicao: 0, duracao: 0 }; try { await gravar(ativo.id!, p.posicao, p.duracao, !ativo.progresso?.concluida); } catch { /* erro já exibido */ } finally { setOcupado(false); } }}>{ativo.progresso?.concluida ? '✓ Concluída · desmarcar' : 'Marcar aula como concluída'}</button><small role="status">{previa ? 'Conclusão desativada na prévia' : status}</small></div>
        {ativo.pdfId && <a className="curso-pdf-download" href={cursosService.pdfURL(curso.id, ativo.pdfId, previa)}>📄 Baixar PDF — {ativo.pdfNome || 'Material da aula'}</a>}
        <div className="curso-actions"><button disabled={ocupado || curso.aulas.findIndex((a) => a.id === aulaId) <= 0} onClick={() => setAulaId(curso.aulas[curso.aulas.findIndex((a) => a.id === aulaId) - 1].id!)}>← Aula anterior</button><button disabled={ocupado || curso.aulas.findIndex((a) => a.id === aulaId) >= curso.aulas.length - 1} onClick={() => setAulaId(curso.aulas[curso.aulas.findIndex((a) => a.id === aulaId) + 1].id!)}>Próxima aula →</button></div>
        {!!ativo.questoes?.length && <ExerciciosAula key={ativo.id} cursoId={curso.id} aulaId={ativo.id!} previa={previa} />}
      </>}
    </div><nav aria-label="Módulos do curso" className="curso-playlist">{modulos.map((m) => <details key={m.id} open={ativo?.moduloId === m.id}>
      <summary>{m.titulo}<small>{m.aulas.filter((a) => a.progresso?.concluida).length}/{m.aulas.length} concluídas</small></summary>
      {m.aulas.map((a, i) => <button key={a.id} disabled={ocupado} aria-current={aulaId === a.id ? 'true' : undefined} onClick={() => setAulaId(a.id!)}><span>{a.bloqueada ? '🔒' : a.progresso?.concluida ? '✓' : String(i + 1).padStart(2, '0')}</span><span>{a.titulo}{a.pdfId && ' · PDF'}{a.bloqueada && <small>{a.motivoBloqueio}</small>}</span></button>)}
    </details>)}</nav></div>
  </Janela>;
}
