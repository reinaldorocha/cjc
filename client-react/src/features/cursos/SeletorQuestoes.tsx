import { useEffect, useState } from 'react';
import { cursosService, type QuestaoCurso } from '../../services/cursosService';

export function SeletorQuestoes({ ids, alterar }: { ids: string[]; alterar: (ids: string[]) => void }) {
  const [busca, setBusca] = useState('');
  const [questoes, setQuestoes] = useState<QuestaoCurso[]>([]);
  const [loading, setLoading] = useState(false);
  const [erro, setErro] = useState('');
  const [aberto, setAberto] = useState(false);
  const [tentativa, setTentativa] = useState(0);
  useEffect(() => {
    if (!aberto) return;
    let ativo = true;
    const timer = window.setTimeout(() => {
      setLoading(true); setErro('');
      cursosService.catalogoQuestoes(busca).then((d) => { if (ativo) setQuestoes(d.questoes); }).catch(() => { if (ativo) setErro('Não foi possível carregar as questões.'); }).finally(() => { if (ativo) setLoading(false); });
    }, 250);
    return () => { ativo = false; window.clearTimeout(timer); };
  }, [busca, aberto, tentativa]);
  return <section className="curso-questoes-seletor">
    <div className="curso-actions"><strong>Exercícios da aula · {ids.length}/50</strong><button type="button" onClick={() => setAberto(!aberto)}>{aberto ? 'Recolher' : 'Vincular questões'}</button></div>
    {aberto && <><label>Buscar no meu banco de questões<input value={busca} onChange={(e) => setBusca(e.target.value)} placeholder="Disciplina, assunto ou enunciado" /></label>
      {loading && <p role="status">Buscando…</p>}{erro && <p role="alert">{erro} <button type="button" onClick={() => setTentativa(tentativa + 1)}>Tentar novamente</button></p>}
      <div className="curso-questoes-opcoes">{questoes.map((q) => <label key={q.id}><input type="checkbox" checked={ids.includes(q.id)} disabled={!ids.includes(q.id) && ids.length >= 50} onChange={(e) => alterar(e.target.checked ? [...ids, q.id] : ids.filter((id) => id !== q.id))} /><span><strong>{q.disciplina} · {q.assunto}</strong>{q.enunciado}</span></label>)}</div>
      {!loading && !erro && !questoes.length && <p className="curso-muted">Nenhuma questão encontrada. Cadastre questões na aba Banco de Questões da mentoria.</p>}
      {ids.filter((id) => !questoes.some((q) => q.id === id)).length > 0 && <p className="curso-muted">Há questões selecionadas fora desta busca. <button type="button" onClick={() => alterar([])}>Limpar seleção inteira</button></p>}
    </>}
  </section>;
}
