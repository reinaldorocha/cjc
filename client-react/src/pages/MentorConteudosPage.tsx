import React, { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../services/api';
import { MentorBancoQuestoesPage } from './MentorBancoQuestoesPage';

type Aba = 'concursos' | 'flashcards' | 'banco_questoes';

const PROMPT_IA_FLASHCARDS = `Crie baralhos de flashcards para estudo no formato JSON estritamente válido.
Responda APENAS com o JSON no seguinte formato:

{
  "nome": "Baralho de Direito Constitucional - Direitos Fundamentais",
  "alcance": "global",
  "cartoes": [
    {
      "tipo": "basico",
      "frente": "O que diz o art. 5º, V da CF/88?",
      "verso": "É assegurado o direito de resposta, proporcional ao agravo, além da indenização por dano material, moral ou à imagem."
    },
    {
      "tipo": "certo_errado",
      "frente": "É livre a manifestação do pensamento, sendo permitido o anonimato.",
      "respostaCorreta": "Errado",
      "explicacao": "É VEDADO o anonimato, nos termos do art. 5º, IV da CF/88."
    },
    {
      "tipo": "multipla_escolha",
      "frente": "Assinale a regra de inviolabilidade do domicílio:",
      "alternativas": [
        "A) Permite entrada à noite por determinação judicial",
        "B) Permite entrada durante o dia por determinação judicial",
        "C) Nenhuma das anteriores"
      ],
      "respostaCorreta": "B",
      "explicacao": "Art. 5º, XI da CF/88."
    }
  ]
}`;

export const MentorConteudosPage: React.FC<{ abaInicial?: Aba }> = ({ abaInicial }) => {
  const navigate = useNavigate();
  const [abaState, setAba] = useState<Aba>(abaInicial || 'concursos');
  const aba = abaInicial || abaState;

  useEffect(() => {
    if (abaInicial) {
      setAba(abaInicial);
    }
  }, [abaInicial]);
  const [concursos, setConcursos] = useState<any[]>([]);
  const [editais, setEditais] = useState<any[]>([]);
  const [baralhos, setBaralhos] = useState<any[]>([]);
  const [erro, setErro] = useState('');
  const [modal, setModal] = useState('');
  const [ocupado, setOcupado] = useState(false);
  const [editId, setEditId] = useState<string | null>(null);

  const [baralhoExpandidoId, setBaralhoExpandidoId] = useState<string | null>(null);
  const [filtroBaralhoAlcance, setFiltroBaralhoAlcance] = useState<'todos' | 'global' | 'concurso'>('todos');
  const [jsonFlashcardsImport, setJsonFlashcardsImport] = useState('');
  const [promptFlashcardsCopiado, setPromptFlashcardsCopiado] = useState(false);
  const [baralho, setBaralho] = useState({ nome: '', alcance: 'global', editalId: '', concursoId: '', descricao: '' });
  const [cartao, setCartao] = useState({ baralhoId: '', tipo: 'basico', frente: '', verso: '', alternativas: '', respostaCorreta: '', explicacao: '' });

  const carregar = useCallback(async () => {
    setErro('');
    try {
      const [c, e, b] = await Promise.all([
        api.listarConcursosMentor(),
        api.listarEditaisMentor(),
        api.listarBaralhosMentor()
      ]);
      setConcursos(c.concursos || []);
      setEditais(e.editais || []);
      setBaralhos(b.baralhos || []);
    } catch (x: any) {
      setErro(x.message);
    }
  }, []);

  useEffect(() => {
    carregar();
  }, [carregar]);

  const executar = async (fn: () => Promise<any>) => {
    setOcupado(true);
    setErro('');
    try {
      await fn();
      setModal('');
      setEditId(null);
      await carregar();
    } catch (x: any) {
      setErro(x.message);
    } finally {
      setOcupado(false);
    }
  };

  const abrirModalConcurso = (c?: any) => {
    if (c?.id) {
      navigate(`/mentor/concursos/${c.id}/editar`);
    } else {
      navigate('/mentor/concursos/novo');
    }
  };

  const abrirModalBaralho = (b?: any) => {
    setEditId(b?.id || null);
    setBaralho(
      b
        ? { nome: b.nome, alcance: b.alcance || 'global', editalId: b.editalId || '', concursoId: b.concursoId || '', descricao: b.descricao || '' }
        : { nome: '', alcance: 'global', editalId: '', concursoId: '', descricao: '' }
    );
    setModal('baralho');
  };

  const abrirModalCartao = (bId: string, c?: any) => {
    setEditId(c?.id || null);
    setCartao(
      c
        ? {
            baralhoId: bId,
            tipo: c.tipo || 'basico',
            frente: c.frente || '',
            verso: c.verso || '',
            alternativas: Array.isArray(c.alternativas) ? c.alternativas.join('\n') : (c.alternativas || ''),
            respostaCorreta: c.respostaCorreta || '',
            explicacao: c.explicacao || ''
          }
        : {
            baralhoId: bId,
            tipo: 'basico',
            frente: '',
            verso: '',
            alternativas: '',
            respostaCorreta: '',
            explicacao: ''
          }
    );
    setModal('cartao');
  };

  const salvarBaralho = (e: React.FormEvent) => {
    e.preventDefault();
    const payload = {
      nome: baralho.nome.trim(),
      descricao: baralho.descricao || null,
      alcance: baralho.alcance,
      editalId: baralho.alcance === 'edital' ? baralho.editalId : null
    };
    executar(() => (editId ? api.alterarBaralhoMentor(editId, payload) : api.criarBaralhoMentor(payload)));
  };

  const salvarCartao = (e: React.FormEvent) => {
    e.preventDefault();
    const alternativas = cartao.tipo === 'multipla_escolha'
      ? (typeof cartao.alternativas === 'string' ? cartao.alternativas.split('\n').map((x) => x.trim()).filter(Boolean) : cartao.alternativas)
      : [];
    const payload = {
      tipo: cartao.tipo,
      frente: cartao.frente.trim(),
      verso: cartao.verso ? cartao.verso.trim() : null,
      alternativas,
      respostaCorreta: cartao.respostaCorreta ? cartao.respostaCorreta.trim() : null,
      explicacao: cartao.explicacao ? cartao.explicacao.trim() : null,
      etiquetas: []
    };
    executar(() =>
      editId
        ? api.alterarCartaoMentor(editId, payload)
        : api.criarCartaoMentor(cartao.baralhoId, payload)
    );
  };

  const desativarBaralho = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!confirm('Desativar este baralho do catálogo?')) return;
    executar(() => api.desativarBaralhoMentor(id));
  };

  const desativarCartao = async (id: string, e?: React.MouseEvent) => {
    if (e) e.stopPropagation();
    if (!confirm('Deseja excluir este cartão do catálogo?')) return;
    executar(() => api.desativarCartaoMentor(id));
  };

  return (
    <div className="mentor-page">
      {!abaInicial && (
        <div className="page-heading">
          <div>
            <h1>Catálogo de Conteúdos da Mentoria</h1>
            <p>Cadastre concursos e seus editais unificados para atribuir aos seus alunos em 1 clique.</p>
          </div>
        </div>
      )}

      {erro && <div className="auth-error">{erro}</div>}

      {!abaInicial && (
        <div className="fc-nav">
          <button className={aba === 'concursos' ? 'active' : ''} onClick={() => setAba('concursos')}>
            🏆 Concursos & Editais
          </button>
          <button className={aba === 'flashcards' ? 'active' : ''} onClick={() => setAba('flashcards')}>
            🎴 Flashcards
          </button>
          <button className={aba === 'banco_questoes' ? 'active' : ''} onClick={() => setAba('banco_questoes')}>
            ❓ Banco de Questões
          </button>
        </div>
      )}

      {aba === 'concursos' && (
        <section>
          <div className="page-heading">
            <div>
              <h2>Concursos & Editais no Catálogo</h2>
              <p>Cada concurso do catálogo já traz seu edital verticalizado pronto para ser importado.</p>
            </div>
            <button className="btn-primary" onClick={() => navigate('/mentor/concursos/novo')}>
              ＋ Cadastrar Concurso com Edital
            </button>
          </div>
          <div className="catalog-grid">
            {concursos.map((c) => {
              const editalVinculado = editais.find((e) => e.concursoId === c.id);
              const totalMaterias = editalVinculado?.materias?.length || 0;

              return (
                <article className="premium-card" key={c.id}>
                  <div>
                    <div className="card-kicker">
                      <span>{c.banca}</span>
                      {c.preEdital && <strong>Pré-edital</strong>}
                    </div>
                    <h2>{c.nome}</h2>
                    <p>{c.cargo || 'Cargo não informado'}</p>
                    <div style={{ marginTop: '10px', fontSize: '11px', color: totalMaterias ? 'var(--green)' : 'var(--text3)' }}>
                      📝 {totalMaterias ? `${totalMaterias} matéria(s) no edital` : 'Sem edital verticalizado anexado'}
                    </div>
                    <div style={{ marginTop: '4px', fontSize: '11px', color: 'var(--text3)' }}>
                      🔄 Revisões automáticas: <strong>{c.prazosRevisao || '1,7,30'} dias</strong>
                    </div>
                  </div>
                  <div style={{ display: 'flex', gap: '8px', marginTop: '14px' }}>
                    <button className="btn-secondary" style={{ padding: '6px 12px', fontSize: '11px', fontWeight: 800 }} onClick={() => abrirModalConcurso(c)}>
                      ✏ Editar Concurso & Edital
                    </button>
                  </div>
                </article>
              );
            })}
            {!concursos.length && (
              <div className="empty-state">
                <h2>Nenhum concurso cadastrado no catálogo</h2>
                <p>Clique acima para cadastrar seu primeiro concurso com edital.</p>
              </div>
            )}
          </div>
        </section>
      )}

      {aba === 'flashcards' && (
        <section>
          <div className="page-heading">
            <div>
              <h2>Flashcards do mentor</h2>
              <p>Publique baralhos globais ou associados aos editais da mentoria e gerencie seus cartões.</p>
            </div>
            <div className="heading-actions" style={{ display: 'flex', gap: '8px' }}>
              <button className="btn-secondary" onClick={() => { setModal('importar_flashcards'); setErro(''); }}>
                ⚡ Importar por IA / JSON
              </button>
              <button className="btn-secondary" disabled={!baralhos.length} onClick={() => abrirModalCartao(baralhos[0]?.id || '')}>
                ＋ Cartão
              </button>
              <button className="btn-primary" onClick={() => abrirModalBaralho()}>
                ＋ Baralho
              </button>
            </div>
          </div>

          {/* FILTRO DE BARALHOS POR ALCANCE */}
          <div className="card-base" style={{ padding: '12px 16px', marginBottom: '20px', display: 'flex', gap: '16px', alignItems: 'center' }}>
            <span style={{ fontSize: '12px', fontWeight: 700, color: 'var(--text2)' }}>Filtrar Baralhos:</span>
            <div style={{ display: 'flex', gap: '8px' }}>
              <button
                type="button"
                className={`btn-secondary ${filtroBaralhoAlcance === 'todos' ? 'active' : ''}`}
                style={{ padding: '4px 12px', fontSize: '12px', background: filtroBaralhoAlcance === 'todos' ? 'var(--accent)' : undefined, color: filtroBaralhoAlcance === 'todos' ? '#fff' : undefined }}
                onClick={() => setFiltroBaralhoAlcance('todos')}
              >
                🌐 + 🎯 Todos
              </button>
              <button
                type="button"
                className={`btn-secondary ${filtroBaralhoAlcance === 'global' ? 'active' : ''}`}
                style={{ padding: '4px 12px', fontSize: '12px', background: filtroBaralhoAlcance === 'global' ? 'var(--accent)' : undefined, color: filtroBaralhoAlcance === 'global' ? '#fff' : undefined }}
                onClick={() => setFiltroBaralhoAlcance('global')}
              >
                🌐 Geral (Todos os alunos)
              </button>
              <button
                type="button"
                className={`btn-secondary ${filtroBaralhoAlcance === 'concurso' ? 'active' : ''}`}
                style={{ padding: '4px 12px', fontSize: '12px', background: filtroBaralhoAlcance === 'concurso' ? 'var(--accent)' : undefined, color: filtroBaralhoAlcance === 'concurso' ? '#fff' : undefined }}
                onClick={() => setFiltroBaralhoAlcance('concurso')}
              >
                🎯 Por Concurso Específico
              </button>
            </div>
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            {baralhos
              .filter((b) => {
                if (filtroBaralhoAlcance === 'global' && b.alcance !== 'global') return false;
                if (filtroBaralhoAlcance === 'concurso' && b.alcance === 'global') return false;
                return true;
              })
              .map((b) => {
              const estaExpandido = baralhoExpandidoId === b.id;
              const totalCartoes = b.totalCartoes ?? b.cartoes?.length ?? 0;
              const dominados = (b.cartoes || []).filter((c: any) => Boolean(c.ultimaRevisao) && (c.repeticoes || 0) > 0).length;
              const pct = totalCartoes > 0 ? Math.round((dominados / totalCartoes) * 100) : 0;

              return (
                <div key={b.id} className="card-base" style={{ padding: '20px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '12px' }}>
                    <div
                      style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer', flex: 1, minWidth: '240px' }}
                      onClick={() => setBaralhoExpandidoId(estaExpandido ? null : b.id)}
                    >
                      <span style={{ fontSize: '28px' }}>{b.icone || '📚'}</span>
                      <div>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                          <h3 style={{ margin: 0, fontSize: '16px', fontWeight: 700 }}>{b.nome}</h3>
                          {pct > 0 && (
                            <span style={{ fontSize: '10px', fontWeight: 800, padding: '2px 6px', borderRadius: '4px', background: 'var(--green)', color: '#fff' }}>
                              🟢 {pct}% Dominado
                            </span>
                          )}
                        </div>
                        <small style={{ color: 'var(--text3)', fontSize: '12px' }}>
                          {b.alcance === 'global' ? '🌐 Todos os alunos' : `🎯 Concurso: ${editais.find((e) => e.id === b.editalId)?.nome || 'Especificado'}`} · <strong>{totalCartoes} cartão(ões)</strong>
                        </small>
                      </div>
                    </div>

                    <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                      <button
                        type="button"
                        className="btn-secondary"
                        style={{ padding: '6px 12px', fontSize: '12px', fontWeight: 700 }}
                        onClick={() => setBaralhoExpandidoId(estaExpandido ? null : b.id)}
                      >
                        {estaExpandido ? '▲ Ocultar Cartões' : `👁 Ver Cartões (${totalCartoes})`}
                      </button>
                      <button
                        type="button"
                        className="btn-primary"
                        style={{ padding: '6px 12px', fontSize: '12px' }}
                        onClick={() => abrirModalCartao(b.id)}
                      >
                        ＋ Cartão
                      </button>
                      <button
                        type="button"
                        className="btn-secondary"
                        style={{ padding: '6px 10px', fontSize: '12px' }}
                        onClick={() => abrirModalBaralho(b)}
                        title="Editar Baralho"
                      >
                        ✏️
                      </button>
                      <button
                        type="button"
                        className="btn-secondary danger-text"
                        style={{ padding: '6px 10px', fontSize: '12px' }}
                        onClick={(e) => desativarBaralho(b.id, e)}
                        title="Excluir Baralho"
                      >
                        🗑
                      </button>
                    </div>
                  </div>

                  {/* PAINEL EXPANDIDO DE CARTÕES DO BARALHO */}
                  {estaExpandido && (
                    <div style={{ marginTop: '20px', paddingTop: '16px', borderTop: '1px solid var(--border)' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '14px' }}>
                        <h4 style={{ margin: 0, fontSize: '14px', color: 'var(--accent)' }}>
                          🎴 Cartões Cadastrados no Baralho ({totalCartoes})
                        </h4>
                        <button
                          type="button"
                          className="btn-primary"
                          style={{ padding: '4px 10px', fontSize: '11px' }}
                          onClick={() => abrirModalCartao(b.id)}
                        >
                          ＋ Adicionar Cartão neste Baralho
                        </button>
                      </div>

                      {!b.cartoes || b.cartoes.length === 0 ? (
                        <div style={{ padding: '24px', textAlign: 'center', background: 'var(--bg-card)', borderRadius: '8px', color: 'var(--text3)', fontSize: '13px' }}>
                          Nenhum cartão cadastrado neste baralho ainda.<br />
                          <button
                            type="button"
                            className="btn-secondary"
                            style={{ marginTop: '10px', padding: '6px 12px', fontSize: '12px' }}
                            onClick={() => abrirModalCartao(b.id)}
                          >
                            ＋ Cadastrar Primeiro Cartão
                          </button>
                        </div>
                      ) : (
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
                          {b.cartoes.map((c: any, index: number) => (
                            <div
                              key={c.id || index}
                              style={{
                                background: 'var(--bg-card)',
                                border: '1px solid var(--border)',
                                borderRadius: '8px',
                                padding: '14px 16px',
                                display: 'flex',
                                justifyContent: 'space-between',
                                gap: '16px',
                                alignItems: 'flex-start'
                              }}
                            >
                              <div style={{ flex: 1, minWidth: 0 }}>
                                <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginBottom: '6px' }}>
                                  <span style={{ fontSize: '11px', fontWeight: 800, padding: '2px 8px', borderRadius: '4px', background: 'var(--accent)', color: '#fff' }}>
                                    #{index + 1}
                                  </span>
                                  <span style={{ fontSize: '11px', fontWeight: 700, color: 'var(--text2)', textTransform: 'uppercase' }}>
                                    {c.tipo === 'multipla_escolha' ? 'Múltipla Escolha' : c.tipo === 'certo_errado' ? 'Certo / Errado' : 'Discursivo'}
                                  </span>
                                </div>
                                <div style={{ fontWeight: 700, fontSize: '14px', marginBottom: '6px', color: 'var(--text1)' }}>
                                  ❓ {c.frente}
                                </div>
                                {c.tipo === 'basico' ? (
                                  <div style={{ fontSize: '13px', color: 'var(--text2)', background: 'rgba(255,255,255,0.02)', padding: '8px', borderRadius: '6px' }}>
                                    💡 <strong>Resposta:</strong> {c.verso}
                                  </div>
                                ) : (
                                  <div style={{ fontSize: '13px', color: 'var(--text2)', background: 'rgba(255,255,255,0.02)', padding: '8px', borderRadius: '6px' }}>
                                    ✅ <strong>Gabarito:</strong> {c.respostaCorreta}
                                    {Array.isArray(c.alternativas) && c.alternativas.length > 0 && (
                                      <div style={{ marginTop: '4px', fontSize: '12px', color: 'var(--text3)' }}>
                                        <strong>Alternativas:</strong> {c.alternativas.join(' | ')}
                                      </div>
                                    )}
                                    {c.explicacao && (
                                      <div style={{ marginTop: '4px', fontSize: '12px', color: 'var(--text3)' }}>
                                        💬 <em>{c.explicacao}</em>
                                      </div>
                                    )}
                                  </div>
                                )}
                              </div>

                              <div style={{ display: 'flex', gap: '6px', flexShrink: 0 }}>
                                <button
                                  type="button"
                                  className="btn-secondary"
                                  style={{ padding: '5px 10px', fontSize: '11px' }}
                                  onClick={() => abrirModalCartao(b.id, c)}
                                >
                                  ✏ Editar
                                </button>
                                <button
                                  type="button"
                                  className="btn-secondary danger-text"
                                  style={{ padding: '5px 10px', fontSize: '11px' }}
                                  onClick={(e) => desativarCartao(c.id, e)}
                                >
                                  🗑 Excluir
                                </button>
                              </div>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  )}
                </div>
              );
            })}

            {!baralhos.length && (
              <div className="empty-state">
                <h2>Nenhum baralho cadastrado no catálogo</h2>
                <p>Clique em "+ Baralho" para criar seu primeiro baralho de flashcards.</p>
              </div>
            )}
          </div>
        </section>
      )}

      {aba === 'banco_questoes' && (
        <section>
          <MentorBancoQuestoesPage />
        </section>
      )}

      {(modal === 'baralho' || modal === 'cartao') && (
        <div className="modal-backdrop" onClick={() => { setModal(''); setEditId(null); }}>
          <form
            className="modal-card small-modal"
            style={{ maxWidth: '520px' }}
            onClick={(e) => e.stopPropagation()}
            onSubmit={modal === 'baralho' ? salvarBaralho : salvarCartao}
          >
            <div className="modal-heading">
              <h2>
                {modal === 'baralho'
                  ? editId ? 'Editar Baralho' : 'Criar Baralho'
                  : 'Criar Cartão'}
              </h2>
              <button type="button" className="icon-button" onClick={() => { setModal(''); setEditId(null); }}>
                ×
              </button>
            </div>
            <div className="form-fields">

              {modal === 'baralho' && (
                <>
                  <label>
                    <span className="field-label">Nome *</span>
                    <input required className="form-control" value={baralho.nome} onChange={(e) => setBaralho((x) => ({ ...x, nome: e.target.value }))} />
                  </label>
                  <label>
                    <span className="field-label">Disponibilizar / Alcance *</span>
                    <select className="form-control" value={baralho.alcance} onChange={(e) => setBaralho((x) => ({ ...x, alcance: e.target.value }))}>
                      <option value="global">🌐 Para todos os alunos (Geral / Global)</option>
                      <option value="edital">🎯 Somente para um concurso específico</option>
                    </select>
                  </label>
                  {baralho.alcance === 'edital' && (
                    <label>
                      <span className="field-label">Concurso *</span>
                      <select required className="form-control" value={baralho.editalId} onChange={(e) => setBaralho((x) => ({ ...x, editalId: e.target.value }))}>
                        <option value="">Selecione um concurso</option>
                        {editais.map((x) => (
                          <option key={x.id} value={x.id}>
                            {x.nome}
                          </option>
                        ))}
                      </select>
                    </label>
                  )}
                </>
              )}

              {modal === 'cartao' && (
                <>
                  <label>
                    <span className="field-label">Baralho *</span>
                    <select required className="form-control" value={cartao.baralhoId} onChange={(e) => setCartao((x) => ({ ...x, baralhoId: e.target.value }))}>
                      <option value="">Selecione</option>
                      {baralhos.map((b) => (
                        <option key={b.id} value={b.id}>
                          {b.nome}
                        </option>
                      ))}
                    </select>
                  </label>
                  <label>
                    <span className="field-label">Tipo</span>
                    <select className="form-control" value={cartao.tipo} onChange={(e) => setCartao((x) => ({ ...x, tipo: e.target.value }))}>
                      <option value="basico">Discursivo</option>
                      <option value="multipla_escolha">Múltipla escolha</option>
                      <option value="certo_errado">Certo / Errado</option>
                    </select>
                  </label>
                  <label>
                    <span className="field-label">Pergunta *</span>
                    <textarea required className="form-control" value={cartao.frente} onChange={(e) => setCartao((x) => ({ ...x, frente: e.target.value }))} />
                  </label>
                  {cartao.tipo === 'basico' ? (
                    <label>
                      <span className="field-label">Resposta *</span>
                      <textarea required className="form-control" value={cartao.verso} onChange={(e) => setCartao((x) => ({ ...x, verso: e.target.value }))} />
                    </label>
                  ) : (
                    <>
                      <label>
                        <span className="field-label">Gabarito *</span>
                        <input required className="form-control" value={cartao.respostaCorreta} onChange={(e) => setCartao((x) => ({ ...x, respostaCorreta: e.target.value }))} />
                      </label>
                      {cartao.tipo === 'multipla_escolha' && (
                        <label>
                          <span className="field-label">Alternativas, uma por linha</span>
                          <textarea className="form-control" value={cartao.alternativas} onChange={(e) => setCartao((x) => ({ ...x, alternativas: e.target.value }))} />
                        </label>
                      )}
                      <label>
                        <span className="field-label">Explicação</span>
                        <textarea className="form-control" value={cartao.explicacao} onChange={(e) => setCartao((x) => ({ ...x, explicacao: e.target.value }))} />
                      </label>
                    </>
                  )}
                </>
              )}
            </div>
            <div className="modal-actions">
              <button type="button" className="btn-secondary" onClick={() => { setModal(''); setEditId(null); }}>
                Cancelar
              </button>
              <button className="btn-primary" disabled={ocupado}>
                {ocupado ? 'Salvando…' : editId ? 'Salvar alterações' : 'Salvar'}
              </button>
            </div>
          </form>
        </div>
      )}
      {/* MODAL DE IMPORTAÇÃO DE FLASHCARDS VIA IA / JSON */}
      {modal === 'importar_flashcards' && (
        <div className="modal-backdrop" onClick={() => setModal('')}>
          <div
            className="modal-card"
            style={{ maxWidth: '750px' }}
            onClick={(e) => e.stopPropagation()}
          >
            <div className="modal-heading">
              <h2>⚡ Importar Baralho de Flashcards via IA / JSON</h2>
              <button type="button" className="icon-button" onClick={() => setModal('')}>×</button>
            </div>

            {erro && <div className="feedback-banner error">{erro}</div>}

            <div style={{ marginBottom: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                <span style={{ fontSize: '13px', fontWeight: 700, color: 'var(--text2)' }}>
                  1. Copie o Prompt abaixo para gerar baralhos e cartões no ChatGPT / Gemini:
                </span>
                <button
                  type="button"
                  className="btn-secondary"
                  style={{ padding: '4px 10px', fontSize: '11px' }}
                  onClick={() => {
                    navigator.clipboard.writeText(PROMPT_IA_FLASHCARDS);
                    setPromptFlashcardsCopiado(true);
                    setTimeout(() => setPromptFlashcardsCopiado(false), 2000);
                  }}
                >
                  {promptFlashcardsCopiado ? '✓ Copiado!' : '📋 Copiar Prompt'}
                </button>
              </div>
              <pre style={{ background: 'var(--bg-input)', padding: '12px', borderRadius: '8px', fontSize: '11px', maxHeight: '140px', overflowY: 'auto', border: '1px solid var(--border)', margin: 0 }}>
                {PROMPT_IA_FLASHCARDS}
              </pre>
            </div>

            <div>
              <span style={{ fontSize: '13px', fontWeight: 700, color: 'var(--text2)', display: 'block', marginBottom: '8px' }}>
                2. Cole o JSON do baralho com os cartões no campo abaixo:
              </span>
              <textarea
                rows={8}
                className="form-control"
                placeholder='{ "nome": "...", "alcance": "global", "cartoes": [ { "tipo": "basico", "frente": "...", "verso": "..." } ] }'
                value={jsonFlashcardsImport}
                onChange={(e) => setJsonFlashcardsImport(e.target.value)}
              />
            </div>

            <div className="modal-actions" style={{ marginTop: '20px' }}>
              <button type="button" className="btn-secondary" onClick={() => setModal('')}>
                Cancelar
              </button>
              <button
                type="button"
                className="btn-primary"
                disabled={ocupado || !jsonFlashcardsImport.trim()}
                onClick={async () => {
                  if (!jsonFlashcardsImport.trim()) return;
                  let parsed: any;
                  try {
                    parsed = JSON.parse(jsonFlashcardsImport.trim());
                  } catch {
                    setErro('O JSON informado é inválido.');
                    return;
                  }
                  const listaBaralhos = Array.isArray(parsed) ? parsed : (parsed.baralhos || [parsed]);
                  if (!listaBaralhos.length) {
                    setErro('Nenhum baralho encontrado no JSON.');
                    return;
                  }

                  setOcupado(true);
                  setErro('');
                  try {
                    for (const bItem of listaBaralhos) {
                      if (!bItem.nome) continue;
                      const editalItem = editais.find(x => x.concursoId === bItem.concursoId || x.id === bItem.editalId);
                      const payloadBaralho = {
                        nome: String(bItem.nome).trim(),
                        descricao: bItem.descricao || null,
                        alcance: bItem.alcance || (bItem.concursoId || bItem.editalId ? 'edital' : 'global'),
                        concursoId: bItem.concursoId || editalItem?.concursoId || null,
                        editalId: bItem.editalId || editalItem?.id || null
                      };

                      const criado = await api.criarBaralhoMentor(payloadBaralho);
                      const bId = criado.id;

                      if (Array.isArray(bItem.cartoes) && bId) {
                        for (const cItem of bItem.cartoes) {
                          if (!cItem.frente) continue;
                          const alternativasList = cItem.tipo === 'multipla_escolha'
                            ? (Array.isArray(cItem.alternativas) ? cItem.alternativas : (typeof cItem.alternativas === 'string' ? cItem.alternativas.split('\n').map((x: string) => x.trim()).filter(Boolean) : []))
                            : [];

                          await api.criarCartaoMentor(bId, {
                            tipo: cItem.tipo || 'basico',
                            frente: String(cItem.frente).trim(),
                            verso: cItem.verso ? String(cItem.verso).trim() : null,
                            alternativas: alternativasList,
                            respostaCorreta: cItem.respostaCorreta ? String(cItem.respostaCorreta).trim() : null,
                            explicacao: cItem.explicacao ? String(cItem.explicacao).trim() : null,
                            etiquetas: []
                          });
                        }
                      }
                    }
                    setModal('');
                    setJsonFlashcardsImport('');
                    await carregar();
                  } catch (e: any) {
                    setErro(e.message || 'Erro ao importar baralhos.');
                  } finally {
                    setOcupado(false);
                  }
                }}
              >
                {ocupado ? 'Importando…' : '⚡ Importar Baralho & Cartões'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
