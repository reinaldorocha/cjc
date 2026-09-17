import { useEffect, useState } from 'react';
import { cursosService, type QuestaoCurso, type ResultadoQuestao } from '../../services/cursosService';

export function ExerciciosAula({ cursoId, aulaId, previa }: { cursoId: string; aulaId: string; previa: boolean }) {
  const [questoes, setQuestoes] = useState<QuestaoCurso[]>([]);
  const [loading, setLoading] = useState(true);
  const [erro, setErro] = useState('');
  const [respostas, setRespostas] = useState<Record<string, string>>({});
  const [resultados, setResultados] = useState<Record<string, ResultadoQuestao>>({});
  const [enviando, setEnviando] = useState('');
  const [tentativa, setTentativa] = useState(0);
  useEffect(() => {
    let ativo = true; setLoading(true); setErro('');
    cursosService.questoes(cursoId, aulaId).then((d) => { if (ativo) setQuestoes(d.questoes); }).catch((e: Error) => { if (ativo) setErro(e.message); }).finally(() => { if (ativo) setLoading(false); });
    return () => { ativo = false; };
  }, [cursoId, aulaId, tentativa]);
  return <section className="curso-exercicios"><h3>Pratique o que aprendeu</h3>{previa && <p className="curso-muted">Prévia: suas respostas não serão registradas no desempenho dos alunos.</p>}
    {loading && <p role="status">Carregando exercícios…</p>}{erro && <p role="alert" className="curso-alerta">{erro} <button onClick={() => setTentativa(tentativa + 1)}>Tentar novamente</button></p>}
    {!loading && !erro && !questoes.length && <p>Nenhum exercício disponível nesta aula.</p>}
    {questoes.map((q, i) => { const alternativas = q.tipo === 'certo_errado' ? ['Certo', 'Errado'] : q.alternativas || []; const resultado = resultados[q.id]; return <article key={q.id}>
      <small>Questão {i + 1} · {q.disciplina}</small><p className="curso-enunciado">{q.enunciado}</p>
      <fieldset disabled={!!enviando || !!resultado}><legend className="sr-only">Escolha sua resposta</legend>{alternativas.map((texto, indice) => { const valor = q.tipo === 'certo_errado' ? texto : String.fromCharCode(65 + indice); return <label key={valor}><input type="radio" name={`q-${q.id}`} value={valor} checked={respostas[q.id] === valor} onChange={() => setRespostas({ ...respostas, [q.id]: valor })} /><span>{q.tipo !== 'certo_errado' && <b>{valor}. </b>}{texto}</span></label>; })}</fieldset>
      {!resultado ? <button disabled={!respostas[q.id] || !!enviando} onClick={async () => { setEnviando(q.id); setErro(''); try { const result = await cursosService.responder(cursoId, aulaId, q.id, respostas[q.id]); setResultados((r) => ({ ...r, [q.id]: result })); } catch (e) { setErro(e instanceof Error ? e.message : 'Falha ao responder.'); } finally { setEnviando(''); } }}>{enviando === q.id ? 'Corrigindo…' : 'Conferir resposta'}</button> : <div className={`curso-resultado ${resultado.correta ? 'acerto' : 'erro'}`} role="status"><strong>{resultado.correta ? 'Resposta correta!' : `Resposta correta: ${resultado.respostaCorreta}`}</strong><p>{resultado.explicacao || 'Sem comentário cadastrado.'}</p><button onClick={() => { setResultados((r) => { const next = { ...r }; delete next[q.id]; return next; }); }}>Tentar outra vez</button></div>}
    </article>; })}
  </section>;
}
