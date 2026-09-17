import React, { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useData } from '../context/DataContext';
import { useAutenticacao } from '../context/AutenticacaoContext';
import { api } from '../services/api';
import { useVisaoAluno } from '../context/VisaoAlunoContext';

const vazio = { nome: '', banca: '', cargo: '', salario: '', dataProva: '', preEdital: false, grupoMeusConcursos: 'foco', realizado: false, resultado: 'aguardando', classificacao: '', notaFinal: '', logoBase64: '', nomeado: { ativo: false, data: '' } };

export const ConcursosPage: React.FC = () => {
  const { getArray, setActiveContestId, activeContestId, alunoId, recarregarConcursos } = useData(); const { usuario } = useAutenticacao(); const visaoAluno = useVisaoAluno(); const navigate = useNavigate(); const concursos = getArray('concursos'); const somenteLeitura = visaoAluno;
  const [aberto, setAberto] = useState(false); const editId: string | null = null; const [form, setForm] = useState<any>({ ...vazio, concursoId: '' }); const [erro, setErro] = useState(''); const [ocupado, setOcupado] = useState(false); const [catalogo, setCatalogo] = useState<any[]>([]);
  useEffect(() => { if (usuario?.papel === 'mentor' && !somenteLeitura) api.listarConcursosMentor().then(d => setCatalogo(d.concursos || [])).catch(() => setCatalogo([])) }, [usuario?.papel, somenteLeitura]);
  const grupos = useMemo(() => ({ foco: concursos.filter(c => !c.realizado && (c.grupoMeusConcursos === 'foco' || !c.grupoMeusConcursos)).sort(ordem), mira: concursos.filter(c => !c.realizado && c.grupoMeusConcursos === 'mira').sort(ordem), realizados: concursos.filter(c => c.realizado).sort(ordem) }), [concursos]);
  const editar = (c?: any) => {
    if (c?.id) {
      if (usuario?.papel === 'mentor' && alunoId && alunoId !== 'eu') {
        navigate(`/mentor/alunos/${alunoId}/concursos/${c.id}/editar`);
      } else if (usuario?.papel === 'mentor') {
        navigate(`/mentor/concursos/${c.id}/editar`);
      } else {
        navigate(`/concursos/${c.id}/editar`);
      }
    } else {
      setForm({ ...vazio, concursoId: '' });
      setErro('');
      setAberto(true);
    }
  }; const mudar = (k: string, v: any) => setForm((x: any) => ({ ...x, [k]: v }));
  const logo = (arquivo?: File) => { if (!arquivo) return; if (arquivo.size > 2 * 1024 * 1024) { setErro('A imagem deve ter no máximo 2 MB.'); return } const leitor = new FileReader(); leitor.onload = () => mudar('logoBase64', leitor.result); leitor.readAsDataURL(arquivo) };
  const salvar = async (e: React.FormEvent) => { e.preventDefault(); if (!editId && usuario?.papel === 'mentor' && !form.concursoId) { setErro('Selecione um concurso do seu catálogo.'); return } if (!form.nome.trim() || !form.banca.trim()) { setErro('Informe o nome e a banca do concurso.'); return } setOcupado(true); setErro(''); const dados = { ...dadosFormulario(form, concursos.length + 1), concursoId: form.concursoId }; try { let id = editId; if (id) await api.alterarConcurso(alunoId, id, dados); else { const criado = await api.criarConcurso(alunoId, dados); id = criado.id; if (form.realizado) await api.alterarConcurso(alunoId, id!, dados); if (form.concursoId) { const resE = await api.listarEditaisMentor(form.concursoId).catch(() => ({ editais: [] })); const editalCatalogo = resE.editais?.[0]; if (editalCatalogo) { await api.atribuirEdital(alunoId, { editalId: editalCatalogo.id, concursoId: id }).catch(() => undefined) } } } await recarregarConcursos(); setActiveContestId(id!); setAberto(false) } catch (x: any) { setErro(x.message || 'Não foi possível salvar o concurso.') } finally { setOcupado(false) } };
  const alterar = async (id: string, dados: any) => { setErro(''); try { await api.alterarConcurso(alunoId, id, dados); await recarregarConcursos() } catch (x: any) { setErro(x.message) } };
  const reordenar = async (c: any, direcao: -1 | 1) => { const pares = concursos.filter(x => !!x.realizado === !!c.realizado && (c.realizado || (x.grupoMeusConcursos || 'foco') === (c.grupoMeusConcursos || 'foco'))).sort(ordem); const i = pares.findIndex(x => x.id === c.id), outro = pares[i + direcao]; if (!outro) return; try { await api.reordenarConcursos(alunoId, [{ id: c.id, ordem: outro.ordemMeusConcursos || i + direcao + 1 }, { id: outro.id, ordem: c.ordemMeusConcursos || i + 1 }]); await recarregarConcursos() } catch (x: any) { setErro(x.message) } };
  const remover = async (c: any) => { if (!confirm(`Desativar “${c.nome}” para este aluno? O histórico será preservado.`)) return; try { await api.desativarConcurso(alunoId, c.id); if (activeContestId === c.id) setActiveContestId(''); await recarregarConcursos() } catch (x: any) { setErro(x.message) } };

  const abrirPainel = (id: string) => {
    setActiveContestId(id);
    if (usuario?.papel === 'mentor' && alunoId && alunoId !== 'eu') {
      navigate(`/mentor/alunos/${alunoId}/concursos/${id}/dashboard`);
    } else {
      navigate('/dashboard');
    }
  };

  const secao = (titulo: string, subtitulo: string, itens: any[], mensagem: string) => (
    <section className="contest-section">
      <div className="section-heading">
        <div>
          <h2>{titulo}</h2>
          <p>{subtitulo}</p>
        </div>
        <span>{itens.length}</span>
      </div>
      {itens.length ? (
        <div className="contest-grid">
          {itens.map((c) => (
            <article key={c.id} className={`contest-card ${c.id === activeContestId ? 'active' : ''}`}>
              <div className="contest-card-top">
                <div className="contest-logo">{c.logoBase64 ? <img src={c.logoBase64} alt="" /> : '🎯'}</div>
                <div>
                  <span className="eyebrow">{c.banca}</span>
                  <h3>{c.nome}</h3>
                  <p>{c.cargo || 'Cargo não informado'}</p>
                </div>
              </div>
              <div className="contest-meta">
                <span>{c.preEdital ? '🚀 Pré-edital' : c.dataProva ? `📅 ${dataBr(c.dataProva)}` : '📅 Sem data'}</span>
                {c.salario != null && <span>{Number(c.salario).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })}</span>}
                {c.realizado && c.resultado && <span className={`result-tag ${c.resultado}`}>{rotuloResultado(c.resultado)}</span>}
                {c.nomeado?.ativo && <span>🎉 Nomeado{c.nomeado.data ? ` em ${dataBr(c.nomeado.data)}` : ''}</span>}
              </div>

              <div className="contest-actions">
                {somenteLeitura ? (
                  <button className="btn-primary" onClick={() => abrirPainel(c.id)}>
                    {c.realizado ? 'Ver histórico' : 'Abrir painel'}
                  </button>
                ) : (
                  <>
                    <button className="icon-button" title="Subir" onClick={() => reordenar(c, -1)}>
                      ↑
                    </button>
                    <button className="icon-button" title="Descer" onClick={() => reordenar(c, 1)}>
                      ↓
                    </button>
                    <button className="icon-button" title="Editar" onClick={() => editar(c)}>
                      ✎
                    </button>
                    {!c.realizado && (
                      <button
                        className="icon-button"
                        title="Mover entre grupos"
                        onClick={() => alterar(c.id, { grupo: c.grupoMeusConcursos === 'mira' ? 'foco' : 'mira' })}
                      >
                        ↔
                      </button>
                    )}
                    <button className="icon-button danger" title="Desativar" onClick={() => remover(c)}>
                      ×
                    </button>
                  </>
                )}
              </div>
            </article>
          ))}
        </div>
      ) : (
        <div className="group-empty">{mensagem}</div>
      )}
    </section>
  );

  return (
    <div>
      <div className="page-heading">
        <div>
          <h1>{somenteLeitura ? 'Meus concursos' : 'Concursos do aluno'}</h1>
          <p>{somenteLeitura ? 'Concursos atribuídos pelo seu mentor.' : 'Atribua concursos do catálogo geral da mentoria a este aluno.'}</p>
        </div>
        {!somenteLeitura && (
          <button className="btn-primary" onClick={() => editar()}>
            ＋ Atribuir Concurso do Catálogo
          </button>
        )}
      </div>
      {erro && !aberto && <div className="form-error">{erro}</div>}
      {secao('🎯 Foco principal', 'Objetivos que concentram a preparação.', grupos.foco, 'Nenhum concurso em foco.')}
      {secao('🔭 Na mira', 'Oportunidades acompanhadas.', grupos.mira, 'Nenhuma oportunidade na mira.')}
      {secao('🏁 Realizados', 'Resultados e provas anteriores.', grupos.realizados, 'Nenhum concurso realizado.')}
      {aberto && (
        <div className="modal-backdrop" onClick={() => setAberto(false)}>
          <form className="modal-card contest-form-modal" onClick={(e) => e.stopPropagation()} onSubmit={salvar}>
            <div className="modal-heading">
              <div>
                <h2>{editId ? 'Editar concurso' : 'Atribuir concurso do catálogo'}</h2>
                <p>Dados de planejamento e resultado.</p>
              </div>
              <button type="button" className="icon-button" onClick={() => setAberto(false)}>
                ×
              </button>
            </div>
            <div className="contest-form-grid">
              <div className="form-fields">
                {!editId && usuario?.papel === 'mentor' && (
                  <Campo nome="Concurso do catálogo *">
                    <select
                      className="form-control"
                      value={form.concursoId}
                      onChange={(e) => {
                        const c = catalogo.find((x) => x.id === e.target.value);
                        setForm((f: any) => ({
                          ...f,
                          concursoId: e.target.value,
                          nome: c?.nome || '',
                          banca: c?.banca || '',
                          cargo: c?.cargo || '',
                          preEdital: !!c?.preEdital,
                          dataProva: c?.dataProva || '',
                          salario: c?.salario ?? ''
                        }));
                      }}
                    >
                      <option value="">Selecione</option>
                      {catalogo
                        .filter((c) => !concursos.some((x) => x.id === c.id))
                        .map((c) => (
                          <option key={c.id} value={c.id}>
                            {c.nome} · {c.banca}
                          </option>
                        ))}
                    </select>
                  </Campo>
                )}
                <Campo nome="Nome *">
                  <input className="form-control" value={form.nome} onChange={(e) => mudar('nome', e.target.value)} />
                </Campo>
                <Campo nome="Banca *">
                  <input className="form-control" value={form.banca} onChange={(e) => mudar('banca', e.target.value)} />
                </Campo>
                <Campo nome="Cargo">
                  <input className="form-control" value={form.cargo} onChange={(e) => mudar('cargo', e.target.value)} />
                </Campo>
                <Campo nome="Salário">
                  <input className="form-control" value={form.salario ?? ''} onChange={(e) => mudar('salario', e.target.value)} placeholder="10.000,00" />
                </Campo>
                <Campo nome="Data da prova">
                  <input className="form-control" type="date" disabled={form.preEdital} value={form.dataProva || ''} onChange={(e) => mudar('dataProva', e.target.value)} />
                </Campo>
                <Campo nome="Grupo">
                  <select className="form-control" value={form.grupoMeusConcursos} onChange={(e) => mudar('grupoMeusConcursos', e.target.value)}>
                    <option value="foco">Foco principal</option>
                    <option value="mira">Na mira</option>
                  </select>
                </Campo>
              </div>
              <aside className="logo-field">
                <div className="logo-preview">{form.logoBase64 ? <img src={form.logoBase64} alt="Prévia" /> : '🎯'}</div>
                <label className="btn-secondary file-button">
                  Escolher imagem
                  <input type="file" accept="image/png,image/jpeg,image/webp" onChange={(e) => logo(e.target.files?.[0])} />
                </label>
                {form.logoBase64 && (
                  <button type="button" className="danger-link" onClick={() => mudar('logoBase64', '')}>
                    Remover imagem
                  </button>
                )}
              </aside>
            </div>
            <div className="toggle-list">
              <label>
                <input type="checkbox" checked={form.preEdital} onChange={(e) => mudar('preEdital', e.target.checked)} /> Concurso em pré-edital
              </label>

              <label>
                <input type="checkbox" checked={form.realizado} onChange={(e) => mudar('realizado', e.target.checked)} /> Concurso já realizado
              </label>
            </div>
            {form.realizado && (
              <>
                <div className="result-fields">
                  <select className="form-control" value={form.resultado} onChange={(e) => mudar('resultado', e.target.value)}>
                    <option value="aguardando">Aguardando resultado</option>
                    <option value="aprovado">Aprovado</option>
                    <option value="cadastro_reserva">Cadastro reserva</option>
                    <option value="reprovado">Reprovado</option>
                    <option value="eliminado">Eliminado</option>
                  </select>
                  <input className="form-control" type="number" value={form.classificacao || ''} onChange={(e) => mudar('classificacao', e.target.value)} placeholder="Classificação" />
                  <input className="form-control" type="number" step="0.01" value={form.notaFinal || ''} onChange={(e) => mudar('notaFinal', e.target.value)} placeholder="Nota final" />
                </div>
                <div className="toggle-list">
                  <label>
                    <input type="checkbox" checked={!!form.nomeado?.ativo} onChange={(e) => mudar('nomeado', { ...form.nomeado, ativo: e.target.checked })} /> Fui nomeado(a)
                  </label>
                  {form.nomeado?.ativo && (
                    <input className="form-control" type="date" value={form.nomeado?.data || ''} onChange={(e) => mudar('nomeado', { ...form.nomeado, data: e.target.value })} />
                  )}
                </div>
              </>
            )}
            {erro && <div className="form-error">{erro}</div>}
            <div className="modal-actions">
              <button type="button" className="btn-secondary" onClick={() => setAberto(false)}>
                Cancelar
              </button>
              <button className="btn-primary" disabled={ocupado}>
                {ocupado ? 'Salvando…' : editId ? 'Salvar alterações' : 'Criar concurso'}
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
};

const Campo: React.FC<{ nome: string; children: React.ReactNode }> = ({ nome, children }) => <label><span className="field-label">{nome}</span>{children}</label>;
const ordem = (a: any, b: any) => (a.ordemMeusConcursos || 0) - (b.ordemMeusConcursos || 0); const dataBr = (v: string) => new Date(`${v}T12:00:00`).toLocaleDateString('pt-BR');
const rotuloResultado = (v: string) => ({ aguardando: 'Aguardando', aprovado: 'Aprovado', cadastro_reserva: 'Cadastro reserva', reprovado: 'Reprovado', eliminado: 'Eliminado' } as any)[v] || v;
const numero = (v: any) => { if (v === '' || v == null) return null; const texto = String(v).trim(); const normal = texto.includes(',') ? texto.replace(/\./g, '').replace(',', '.') : texto; const n = Number(normal); return Number.isFinite(n) ? n : null };
const dadosFormulario = (f: any, ordemPadrao: number) => ({ nome: f.nome.trim(), banca: f.banca.trim(), cargo: f.cargo.trim(), logotipo: f.logoBase64 || '', salario: numero(f.salario), limparSalario: !f.salario, dataProva: f.preEdital ? '' : f.dataProva || '', preEdital: !!f.preEdital, grupo: f.realizado ? 'realizado' : f.grupoMeusConcursos, ordem: f.ordemMeusConcursos || ordemPadrao, resultado: f.realizado ? f.resultado : '', classificacao: f.classificacao ? Number(f.classificacao) : null, limparClassificacao: !f.classificacao, notaFinal: f.notaFinal ? Number(f.notaFinal) : null, limparNotaFinal: !f.notaFinal, nomeado: !!f.nomeado?.ativo, dataNomeacao: f.nomeado?.ativo ? f.nomeado.data || '' : '' });
