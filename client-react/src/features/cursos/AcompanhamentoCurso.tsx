import { useEffect, useState } from 'react';
import { cursosService, type AlunoCurso, type Curso } from '../../services/cursosService';
import { Janela } from './CursoDialog';

export function AcompanhamentoCurso({ curso, fechar }: { curso: Curso; fechar: () => void }) {
  const [alunos, setAlunos] = useState<AlunoCurso[]>([]);
  const [loading, setLoading] = useState(true);
  const [erro, setErro] = useState('');
  const [busca, setBusca] = useState('');
  const [situacao, setSituacao] = useState('');
  const [tentativa, setTentativa] = useState(0);
  useEffect(() => { let ativo = true; setLoading(true); setErro(''); cursosService.acompanhar(curso.id).then((d) => { if (ativo) setAlunos(d.alunos); }).catch((e: Error) => { if (ativo) setErro(e.message); }).finally(() => { if (ativo) setLoading(false); }); return () => { ativo = false; }; }, [curso.id, tentativa]);
  const filtrados = alunos.filter((a) => `${a.nome} ${a.email}`.toLocaleLowerCase().includes(busca.toLocaleLowerCase()) && (!situacao || a.situacao === situacao));
  const concluidos = alunos.filter((a) => a.situacao === 'Concluído').length;
  const parados = alunos.filter((a) => a.situacao === 'Sem atividade há 7 dias').length;
  const iniciados = alunos.filter((a) => a.situacao !== 'Não iniciado').length;
  return <Janela titulo={`Acompanhamento · ${curso.titulo}`} fechar={fechar}>
    <p className="curso-muted">Dados da versão publicada, considerando alunos ativos com acesso atual. Conclusões são marcadas pelo aluno; não comprovam tempo assistido.</p>
    {loading ? <p role="status">Carregando acompanhamento…</p> : erro ? <p className="curso-alerta" role="alert">{erro} <button onClick={() => setTentativa(tentativa + 1)}>Tentar novamente</button></p> : <>
      <div className="curso-indicadores">{[['Com acesso', alunos.length], ['Começaram', iniciados], ['Concluíram', concluidos], ['Sem atividade há 7 dias', parados]].map(([label, n]) => <div key={label}><strong>{n}</strong><span>{label}</span></div>)}</div>
      <div className="cursos-filtros"><input aria-label="Buscar aluno" value={busca} onChange={(e) => setBusca(e.target.value)} placeholder="Nome ou e-mail" /><select aria-label="Filtrar situação" value={situacao} onChange={(e) => setSituacao(e.target.value)}><option value="">Todas as situações</option>{['Não iniciado', 'Em andamento', 'Concluído', 'Sem atividade há 7 dias'].map((s) => <option key={s}>{s}</option>)}</select></div>
      <div className="curso-tabela"><table><thead><tr><th>Aluno</th><th>Progresso</th><th>Situação</th><th>Última atividade</th><th>Exercícios</th></tr></thead><tbody>{filtrados.map((a) => <tr key={a.id}><td><strong>{a.nome}</strong><small>{a.email}</small></td><td><progress value={a.resumo.percentual} max={100} aria-label={`Progresso de ${a.nome}`} /><small>{a.resumo.concluidas}/{a.resumo.total} aulas · {a.resumo.percentual}%</small></td><td>{a.situacao}</td><td>{a.resumo.ultimaAtividade ? new Date(a.resumo.ultimaAtividade).toLocaleString('pt-BR') : 'Ainda não acessou'}</td><td>{a.respostas ? `${a.acertos}/${a.respostas} acertos` : 'Sem respostas'}</td></tr>)}</tbody></table></div>
      {!filtrados.length && <p className="cursos-vazio">{alunos.length ? 'Nenhum aluno corresponde ao filtro.' : 'Nenhum aluno ativo com acesso à versão publicada.'}</p>}
    </>}
  </Janela>;
}
