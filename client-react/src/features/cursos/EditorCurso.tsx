import { useEffect, useState } from 'react';
import { api } from '../../services/api';
import { cursosService, type AulaCurso, type Curso } from '../../services/cursosService';
import { Janela } from './CursoDialog';
import { agruparModulos, aulasDosModulos, validarModulo, type ModuloCurso } from './modulos';
import { SeletorQuestoes } from './SeletorQuestoes';

type Opcao = { id: string; nome: string };
type Rascunho = { indice: number; modulo: ModuloCurso };
const mensagem = (e: unknown) => e instanceof Error ? e.message : 'Não foi possível concluir a operação.';
const aulaVazia = (): AulaCurso => ({ id: crypto.randomUUID(), titulo: '', modulo: '', videoId: '', questoes: [] });

export function Editor({ curso, alunoId, salvo, fechar }: { curso: Curso; alunoId: string; salvo: () => void; fechar: () => void }) {
  const [form, setForm] = useState<Curso>(() => structuredClone(curso));
  const [modulos, setModulos] = useState(() => agruparModulos(curso.aulas));
  const [rascunho, setRascunho] = useState<Rascunho | null>(null);
  const [aulaAtiva, setAulaAtiva] = useState(0);
  const [alunos, setAlunos] = useState<Opcao[]>([]);
  const [concursos, setConcursos] = useState<Opcao[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [erro, setErro] = useState('');
  const [aviso, setAviso] = useState('');
  const [salvando, setSalvando] = useState(false);
  const [enviando, setEnviando] = useState(false);
  const [ultimaVersao, setUltimaVersao] = useState(() => JSON.stringify({ form: curso, modulos: agruparModulos(curso.aulas) }));
  const alterado = !!rascunho || JSON.stringify({ form, modulos }) !== ultimaVersao;
  const bloqueado = salvando || enviando;
  useEffect(() => {
    let ativo = true;
    Promise.all([api.listarAlunos(), api.listarConcursosMentor()]).then(([a, c]) => {
      if (ativo) { setAlunos(a.alunos || []); setConcursos(c.concursos || []); }
    }).catch((e) => { if (ativo) setErro(mensagem(e)); }).finally(() => { if (ativo) setCarregando(false); });
    return () => { ativo = false; };
  }, []);
  const opcoes = form.escopo === 'alunos' ? alunos : concursos;
  const alterar = (dados: Partial<Curso>) => setForm((f) => ({ ...f, ...dados }));
  const editarModulo = (indice: number) => {
    setErro(''); setAviso(''); setAulaAtiva(0);
    setRascunho({ indice, modulo: indice < 0 ? { id: crypto.randomUUID(), titulo: '', descricao: '', capaUrl: '', aulas: [aulaVazia()], liberarEm: '', exigeAnterior: false } : structuredClone(modulos[indice]) });
  };
  const alterarModulo = (dados: Partial<ModuloCurso>) => setRascunho((r) => r && ({ ...r, modulo: { ...r.modulo, ...dados } }));
  const alterarAula = (dados: Partial<AulaCurso>) => setRascunho((r) => r && ({ ...r, modulo: { ...r.modulo, aulas: r.modulo.aulas.map((a, i) => i === aulaAtiva ? { ...a, ...dados } : a) } }));
  const cancelarModulo = () => { if (window.confirm('Descartar as alterações não salvas deste módulo?')) { setRascunho(null); setErro(''); } };
  const fecharEditor = () => { if (!bloqueado && (!alterado || window.confirm('Descartar as alterações ainda não salvas?'))) fechar(); };
  useEffect(() => { const avisar = (e: BeforeUnloadEvent) => { if (alterado) e.preventDefault(); }; window.addEventListener('beforeunload', avisar); return () => window.removeEventListener('beforeunload', avisar); }, [alterado]);
  const salvar = async (publicar = false) => {
    setErro(''); setAviso('');
    let proximos = modulos;
    if (rascunho) {
      const invalid = validarModulo(rascunho.modulo, modulos.filter((_, i) => i !== rascunho.indice));
      if (invalid) { setErro(invalid); return; }
      proximos = [...modulos];
      if (rascunho.indice < 0) proximos.push(rascunho.modulo); else proximos[rascunho.indice] = rascunho.modulo;
    }
    const aulas = aulasDosModulos(proximos);
    if (aulas.length > 300 || (publicar && !aulas.length)) { setErro('Para publicar, adicione entre 1 e 300 aulas.'); return; }
    setSalvando(true);
    try {
      let { curso: gravado } = await cursosService.salvar({ ...form, aulas });
      const aplicar = (c: Curso) => { setForm(c); const m = agruparModulos(c.aulas); setModulos(m); setUltimaVersao(JSON.stringify({ form: c, modulos: m })); setRascunho(null); };
      aplicar(gravado);
      if (publicar) { ({ curso: gravado } = await cursosService.publicar(gravado, true)); aplicar(gravado); salvo(); }
      else setAviso(gravado.publicado ? 'Alterações salvas no rascunho. Publique a atualização quando estiver pronta.' : 'Rascunho salvo. Os alunos só terão acesso depois da publicação.');
    } catch (e) { setErro(mensagem(e)); } finally { setSalvando(false); }
  };
  const aula = rascunho?.modulo.aulas[aulaAtiva];
  const restantes = modulos.filter((_, i) => i !== rascunho?.indice).reduce((total, m) => total + m.aulas.length, 0);
  return <Janela titulo={form.id ? 'Editar curso e acesso' : 'Novo curso'} fechar={fecharEditor}>
    <form className="curso-form" onSubmit={(e) => { e.preventDefault(); void salvar((e.nativeEvent as SubmitEvent).submitter?.getAttribute('data-publicar') === 'true'); }}>
      <fieldset disabled={bloqueado}>
        <p className="curso-notice">{form.publicado ? 'Publicado · as edições ficam em rascunho até você publicar uma atualização.' : 'Rascunho · prepare o curso e publique quando estiver pronto.'}</p>
        <div className="curso-form-grid">
          <label>Título do curso<input required maxLength={180} value={form.titulo} onChange={(e) => alterar({ titulo: e.target.value })} /></label>
          <label>Categoria<input maxLength={80} placeholder="Ex.: Direito Constitucional" value={form.categoria} onChange={(e) => alterar({ categoria: e.target.value })} /></label>
        </div>
        <label>Descrição<textarea maxLength={16000} rows={2} value={form.descricao} onChange={(e) => alterar({ descricao: e.target.value })} /></label>
        <div className="curso-capa-editor">
          <label>Enviar imagem de capa (opcional)<input type="file" accept="image/jpeg,image/png,.jpg,.jpeg,.png" onChange={async (e) => {
            const file = e.target.files?.[0]; e.target.value = ''; if (!file) return;
            if (file.size > 5 * 1024 * 1024 || !/\.(jpe?g|png)$/i.test(file.name)) { setErro('Selecione uma imagem JPG ou PNG de até 5 MB.'); return; }
            setEnviando(true); setErro('');
            try { const capa = await cursosService.enviarCapa(file); alterar({ capaUrl: capa.url }); setAviso('Capa enviada. Salve o rascunho ou publique para confirmar.'); } catch (error) { setErro(mensagem(error)); } finally { setEnviando(false); }
          }} /></label>
          <small>JPG ou PNG, até 5 MB. Recomendado: formato horizontal 16:9.</small>
          {form.capaUrl.startsWith('/api/v1/cursos/capas/') ? <p>Imagem enviada do computador.</p> : <label>Ou usar link da imagem<input type="url" placeholder="https://…" maxLength={2048} value={form.capaUrl} onChange={(e) => alterar({ capaUrl: e.target.value })} /></label>}
          {form.capaUrl && <><img key={form.capaUrl} src={form.capaUrl} alt="Prévia da capa do curso" style={{ width: '100%', maxWidth: 360, aspectRatio: '16 / 9', objectFit: 'cover', borderRadius: 12 }} /><button type="button" onClick={() => alterar({ capaUrl: '' })}>Remover capa</button></>}
          {enviando && <p role="status">Enviando arquivo…</p>}
        </div>
        <p className="curso-muted">O catálogo mostra um card por curso. Os cards dos módulos aparecem ao abrir o curso.</p>
        <h3>Quem pode assistir</h3>
        <select aria-label="Disponibilidade do curso" value={form.escopo} onChange={(e) => alterar({ escopo: e.target.value as Curso['escopo'], destinatarios: e.target.value === 'alunos' && alunoId ? [alunoId] : [] })}>
          <option value="global">Todos os alunos da minha mentoria</option><option value="alunos">Alunos selecionados</option><option value="concursos">Alunos dos concursos selecionados</option>
        </select>
        {form.escopo !== 'global' && <div className="curso-destinatarios">
          {carregando ? <p>Carregando destinatários…</p> : opcoes.length === 0 ? <p>Nenhum destinatário disponível. Cadastre-o na mentoria primeiro.</p> : opcoes.map((o) => <label key={o.id}><input type="checkbox" checked={form.destinatarios.includes(o.id)} onChange={(e) => alterar({ destinatarios: e.target.checked ? [...form.destinatarios, o.id] : form.destinatarios.filter((id) => id !== o.id) })} />{o.nome}</label>)}
          {form.destinatarios.filter((id) => !opcoes.some((o) => o.id === id)).map((id) => <label key={id}><input type="checkbox" checked onChange={() => alterar({ destinatarios: form.destinatarios.filter((v) => v !== id) })} />Destinatário indisponível — desmarque para remover</label>)}
        </div>}
        <div className="curso-modulos-heading"><div><h3>Módulos do curso</h3><p className="curso-muted">{modulos.length} módulos · {modulos.reduce((n, m) => n + m.aulas.length, 0)} aulas salvas</p></div>{!rascunho && <button type="button" onClick={() => editarModulo(-1)}>+ Criar módulo</button>}</div>
        {!rascunho ? <div className="curso-modulos-lista">
          {!modulos.length && <p className="curso-modulo-vazio">Crie um módulo, informe seu título e adicione as aulas.</p>}
          {modulos.map((m, i) => <article className="curso-modulo-resumo" key={i}>
            <span className="curso-modulo-numero">{String(i + 1).padStart(2, '0')}</span><div><strong>{m.titulo}</strong><small>{m.aulas.length} aulas · {m.aulas.filter((a) => a.pdfId).length} PDFs{m.capaUrl ? ' · Capa própria' : ''}</small></div>
            <div className="curso-actions"><button type="button" disabled={i === 0} aria-label={`Subir módulo ${m.titulo}`} onClick={() => { const m = [...modulos]; [m[i], m[i - 1]] = [m[i - 1], m[i]]; setModulos(m); setAviso('Ordem alterada. Salve o curso para confirmar.'); }}>↑</button><button type="button" disabled={i === modulos.length - 1} aria-label={`Descer módulo ${m.titulo}`} onClick={() => { const m = [...modulos]; [m[i], m[i + 1]] = [m[i + 1], m[i]]; setModulos(m); setAviso('Ordem alterada. Salve o curso para confirmar.'); }}>↓</button><button type="button" onClick={() => editarModulo(i)}>Editar módulo</button><button type="button" aria-label={`Remover módulo ${m.titulo}`} onClick={() => { if (window.confirm(`Remover o módulo “${m.titulo}” e suas ${m.aulas.length} aulas?`)) { setModulos(modulos.filter((_, j) => j !== i)); setAviso('Módulo removido da edição. Salve o curso para confirmar.'); } }}>✕</button></div>
          </article>)}
        </div> : <section className="curso-modulo-editor">
          <div className="curso-modulos-heading"><h3>{rascunho.indice < 0 ? 'Criar módulo' : 'Editar módulo'}</h3><button type="button" onClick={cancelarModulo}>Cancelar edição</button></div>
          <label>Título do módulo<input required autoFocus maxLength={100} placeholder="Ex.: Módulo 1 — Fundamentos" value={rascunho.modulo.titulo} onChange={(e) => alterarModulo({ titulo: e.target.value })} /></label>
          <label>Descrição do módulo (opcional)<input maxLength={250} placeholder="Breve resumo do que o aluno aprenderá neste módulo…" value={rascunho.modulo.descricao || ''} onChange={(e) => alterarModulo({ descricao: e.target.value })} /></label>
          <div className="curso-capa-editor" style={{ marginBottom: '12px' }}>
            <label>Capa do módulo (opcional — usa a capa do curso se não informada)<input type="file" accept="image/jpeg,image/png,.jpg,.jpeg,.png" onChange={async (e) => {
              const file = e.target.files?.[0]; e.target.value = ''; if (!file) return;
              if (file.size > 5 * 1024 * 1024 || !/\.(jpe?g|png)$/i.test(file.name)) { setErro('Selecione uma imagem JPG ou PNG de até 5 MB.'); return; }
              setEnviando(true); setErro('');
              try { const capa = await cursosService.enviarCapa(file); alterarModulo({ capaUrl: capa.url }); } catch (error) { setErro(mensagem(error)); } finally { setEnviando(false); }
            }} /></label>
            <small>Ou link direto da imagem:</small>
            {rascunho.modulo.capaUrl?.startsWith('/api/v1/cursos/capas/') ? <p>Imagem enviada do computador.</p> : <input type="url" aria-label="Link da capa do módulo" placeholder="https://…" maxLength={2048} value={rascunho.modulo.capaUrl || ''} onChange={(e) => alterarModulo({ capaUrl: e.target.value })} />}
            {rascunho.modulo.capaUrl && <div style={{ marginTop: '6px', display: 'flex', alignItems: 'center', gap: '10px' }}>
              <img src={rascunho.modulo.capaUrl} alt="Capa do módulo" style={{ width: 140, height: 78, objectFit: 'cover', borderRadius: 8, border: '1px solid var(--border)' }} />
              <button type="button" onClick={() => alterarModulo({ capaUrl: '' })}>Remover capa do módulo</button>
            </div>}
          </div>
          <div className="curso-form-grid"><label>Liberar a partir de (opcional)<input type="date" value={rascunho.modulo.liberarEm || ''} onChange={(e) => alterarModulo({ liberarEm: e.target.value })} /></label><label className="curso-checkbox"><input type="checkbox" checked={!!rascunho.modulo.exigeAnterior} onChange={(e) => alterarModulo({ exigeAnterior: e.target.checked })} />Exigir conclusão dos módulos anteriores</label></div>
          <p className="curso-muted">Quando as duas regras estão configuradas, ambas precisam ser atendidas.</p>
          <div className="curso-modulo-aulas">
            <nav aria-label="Aulas do módulo" className="curso-aulas-seletor">
              {rascunho.modulo.aulas.map((a, i) => <button key={i} type="button" aria-current={aulaAtiva === i ? 'true' : undefined} onClick={() => setAulaAtiva(i)}><span>{i + 1}</span><span>{a.titulo || 'Nova aula'}{a.pdfId && <small>PDF anexado</small>}</span></button>)}
              <button type="button" disabled={restantes + rascunho.modulo.aulas.length >= 300} onClick={() => { alterarModulo({ aulas: [...rascunho.modulo.aulas, aulaVazia()] }); setAulaAtiva(rascunho.modulo.aulas.length); }}>+ Adicionar aula</button>
            </nav>
            {aula && <div className="curso-aula-campos" key={aulaAtiva}>
              <h4>Aula {aulaAtiva + 1}</h4>
              <label>Título da aula<input required maxLength={180} value={aula.titulo} onChange={(e) => alterarAula({ titulo: e.target.value })} /></label>
              <label>Link do YouTube<input required placeholder="https://www.youtube.com/watch?v=…" value={aula.videoId} onChange={(e) => alterarAula({ videoId: e.target.value })} /></label>
              <SeletorQuestoes key={aula.id || aulaAtiva} ids={aula.questoes || []} alterar={(questoes) => alterarAula({ questoes })} />
              <div className="curso-pdf-anexo"><label>PDF da aula (opcional, até 10 MB)<input type="file" accept="application/pdf,.pdf" onChange={async (e) => {
                const file = e.target.files?.[0]; e.target.value = ''; if (!file) return;
                if (file.size > 10 * 1024 * 1024 || !file.name.toLowerCase().endsWith('.pdf')) { setErro('Selecione um PDF de até 10 MB.'); return; }
                setEnviando(true); setErro('');
                try { const pdf = await cursosService.enviarPDF(file); alterarAula({ pdfId: pdf.id, pdfNome: pdf.nome }); } catch (error) { setErro(mensagem(error)); } finally { setEnviando(false); }
              }} /></label>{enviando && <p role="status">Enviando PDF…</p>}{aula.pdfId && <div className="curso-actions"><span>📄 {aula.pdfNome || 'PDF anexado'}</span><button type="button" onClick={() => alterarAula({ pdfId: undefined, pdfNome: undefined })}>Remover PDF</button></div>}</div>
              <div className="curso-actions"><button type="button" disabled={aulaAtiva === 0} onClick={() => { const a = [...rascunho.modulo.aulas]; [a[aulaAtiva], a[aulaAtiva - 1]] = [a[aulaAtiva - 1], a[aulaAtiva]]; alterarModulo({ aulas: a }); setAulaAtiva(aulaAtiva - 1); }}>↑ Subir</button><button type="button" disabled={aulaAtiva === rascunho.modulo.aulas.length - 1} onClick={() => { const a = [...rascunho.modulo.aulas]; [a[aulaAtiva], a[aulaAtiva + 1]] = [a[aulaAtiva + 1], a[aulaAtiva]]; alterarModulo({ aulas: a }); setAulaAtiva(aulaAtiva + 1); }}>↓ Descer</button><button type="button" disabled={rascunho.modulo.aulas.length === 1} onClick={() => { if (window.confirm('Remover esta aula do módulo?')) { alterarModulo({ aulas: rascunho.modulo.aulas.filter((_, i) => i !== aulaAtiva) }); setAulaAtiva(Math.max(0, aulaAtiva - 1)); } }}>Remover aula</button></div>
            </div>}
          </div>
        </section>}
        {erro && <p role="alert" className="curso-alerta">{erro}</p>}
        {aviso && <p role="status" className="curso-notice">{aviso}</p>}
        <footer className="curso-actions"><button type="button" onClick={fecharEditor}>Fechar</button><button disabled={carregando || bloqueado}>{salvando ? 'Salvando…' : rascunho ? 'Salvar módulo' : 'Salvar rascunho'}</button><button data-publicar="true" className="curso-primary" disabled={carregando || bloqueado || (form.escopo !== 'global' && form.destinatarios.length === 0)}>{form.publicado ? 'Publicar atualização' : 'Publicar curso'}</button></footer>
      </fieldset>
    </form>
  </Janela>;
}
