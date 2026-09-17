import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { useAutenticacao } from '../context/AutenticacaoContext';
import { useData } from '../context/DataContext';
import { api } from '../services/api';
import { Pagination } from '../components/ui/Pagination';

type TipoMaterial = 'youtube' | 'arquivo' | 'texto' | 'link';
type Escopo = 'global' | 'edital';

type MaterialApoio = {
  id: string;
  titulo: string;
  descricao?: string;
  tipo: TipoMaterial;
  url?: string;
  texto?: string;
  escopo: Escopo;
  editalId?: string;
  arquivoNome?: string;
  pasta?: string;
};

const abas = [
  { key: 'youtube' as const, label: 'Vídeos (YouTube)' },
  { key: 'arquivo' as const, label: 'Arquivos' },
  { key: 'texto' as const, label: 'Textos' },
  { key: 'link' as const, label: 'Links' }
];

const icones: Record<TipoMaterial, string> = {
  youtube: '🎥',
  arquivo: '📎',
  texto: '📝',
  link: '🔗'
};

const vazioMaterial = () => ({
  id: '',
  titulo: '',
  descricao: '',
  tipo: 'youtube' as TipoMaterial,
  url: '',
  texto: '',
  escopo: 'global' as Escopo,
  editalId: '',
  pasta: 'Geral'
});

export const MateriaisApoioPage: React.FC = () => {
  const { usuario } = useAutenticacao();
  const { alunoId, editalAtivo } = useData();
  const mentor = usuario?.papel === 'mentor';
  const isMentorPage = mentor && alunoId === 'eu';
  const [mentorEditais, setMentorEditais] = useState<any[]>([]);
  const [mentorConcursos, setMentorConcursos] = useState<any[]>([]);
  const [materiais, setMateriais] = useState<MaterialApoio[]>([]);
  const [pagina, setPagina] = useState(1);

  // Aba ativa (Padrão: Vídeos)
  const [aba, setAba] = useState<TipoMaterial>('youtube');
  // Pasta selecionada dentro da aba (null = exibição em grade de pastas, ou nome da pasta)
  const [pastaSelecionada, setPastaSelecionada] = useState<string | null>(null);

  const [form, setForm] = useState(() => vazioMaterial());
  const [criarNovaPasta, setCriarNovaPasta] = useState(false);
  const [nomeNovaPasta, setNomeNovaPasta] = useState('');

  const [editId, setEditId] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [erro, setErro] = useState('');
  const [carregando, setCarregando] = useState(false);
  const [ocupado, setOcupado] = useState(false);
  const [arquivo, setArquivo] = useState<File | null>(null);
  const [videoExpandido, setVideoExpandido] = useState<{ url: string; titulo: string } | null>(null);

  const carregar = useCallback(async () => {
    setErro('');
    setCarregando(true);
    try {
      const dados = isMentorPage
        ? await api.listarMateriaisApoioMentor()
        : await api.listarMateriaisApoio(alunoId, editalAtivo?.id || '');
      setMateriais(dados.materiais || []);
    } catch (x: any) {
      setErro(x.message || 'Não foi possível carregar os materiais.');
    } finally {
      setCarregando(false);
    }
  }, [alunoId, editalAtivo?.id, isMentorPage]);

  useEffect(() => {
    carregar();
  }, [carregar]);

  useEffect(() => {
    if (!isMentorPage) return;
    Promise.all([api.listarConcursosMentor(), api.listarEditaisMentor()])
      .then(([cDados, eDados]) => {
        setMentorConcursos(cDados.concursos || []);
        setMentorEditais(eDados.editais || []);
      })
      .catch(() => {
        setMentorConcursos([]);
        setMentorEditais([]);
      });
  }, [isMentorPage]);

  // Materiais filtrados apenas pelo Tipo ativo
  const materiaisDoTipo = useMemo(() => {
    return materiais.filter((m) => m.tipo === aba);
  }, [aba, materiais]);

  // Extração única de pastas para a aba atual
  const pastasDaAba = useMemo(() => {
    const conjunto = new Set<string>();
    materiaisDoTipo.forEach((m) => {
      const p = m.pasta?.trim() || 'Geral';
      conjunto.add(p);
    });
    if (conjunto.size === 0) {
      conjunto.add('Geral');
    }
    return Array.from(conjunto).sort();
  }, [materiaisDoTipo]);

  // Materiais da pasta atualmente selecionada
  const materiaisDaPasta = useMemo(() => {
    if (!pastaSelecionada) return [];
    if (pastaSelecionada === '__TODAS__') return materiaisDoTipo;
    return materiaisDoTipo.filter((m) => (m.pasta?.trim() || 'Geral') === pastaSelecionada);
  }, [materiaisDoTipo, pastaSelecionada]);
  const totalPaginas = Math.max(1, Math.ceil(materiaisDaPasta.length / 12));
  const paginaAtual = Math.min(pagina, totalPaginas);
  const materiaisDaPagina = materiaisDaPasta.slice((paginaAtual - 1) * 12, paginaAtual * 12);

  const abrirNovo = (pastaPadrao?: string) => {
    setEditId(null);
    setArquivo(null);
    setCriarNovaPasta(false);
    setNomeNovaPasta('');

    const possuiConcursos = mentorConcursos.length > 0;
    const primeiroTarget = mentorConcursos[0] ? (mentorEditais.find((e) => e.concursoId === mentorConcursos[0].id)?.id || mentorConcursos[0].id) : '';

    setForm({
      ...vazioMaterial(),
      tipo: aba,
      pasta: pastaPadrao || pastaSelecionada || (pastasDaAba[0] || 'Geral'),
      escopo: isMentorPage && possuiConcursos ? 'global' : editalAtivo ? 'edital' : 'global',
      editalId: isMentorPage && possuiConcursos ? primeiroTarget : ''
    });
    setShowForm(true);
    setErro('');
  };

  const abrirEdicao = (material: MaterialApoio) => {
    setEditId(material.id);
    setArquivo(null);
    setCriarNovaPasta(false);
    setNomeNovaPasta('');

    setForm({
      id: material.id,
      titulo: material.titulo,
      descricao: material.descricao || '',
      tipo: material.tipo,
      url: material.url || '',
      texto: material.texto || '',
      escopo: material.escopo,
      editalId: material.editalId || '',
      pasta: material.pasta?.trim() || 'Geral'
    });
    setShowForm(true);
    setErro('');
  };

  const salvarMaterial = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!form.titulo.trim()) {
      setErro('Informe o título do material.');
      return;
    }
    if (form.tipo === 'texto' && !form.texto.trim()) {
      setErro('Informe o texto do material.');
      return;
    }
    if (form.tipo === 'arquivo') {
      if (!editId && !arquivo) {
        setErro('Selecione um arquivo para upload.');
        return;
      }
    } else if (!form.url.trim()) {
      setErro('Informe a URL do material.');
      return;
    }

    const pastaFinal = criarNovaPasta ? (nomeNovaPasta.trim() || 'Geral') : (form.pasta.trim() || 'Geral');

    setOcupado(true);
    setErro('');
    try {
      let payload: any;
      if (form.tipo === 'arquivo' && (arquivo || !editId)) {
        const data = new FormData();
        data.append('titulo', form.titulo.trim());
        if (form.descricao.trim()) data.append('descricao', form.descricao.trim());
        data.append('tipo', form.tipo);
        data.append('escopo', form.escopo);
        data.append('pasta', pastaFinal);
        if (form.escopo === 'edital') {
          data.append('editalId', isMentorPage ? form.editalId || '' : editalAtivo?.id || '');
        }
        if (arquivo) {
          data.append('arquivo', arquivo, arquivo.name);
        }
        payload = data;
      } else {
        payload = {
          titulo: form.titulo.trim(),
          descricao: form.descricao.trim() || null,
          tipo: form.tipo,
          url: form.tipo === 'texto' ? null : form.url.trim(),
          texto: form.tipo === 'texto' ? form.texto.trim() : null,
          escopo: form.escopo,
          editalId: form.escopo === 'edital' ? (isMentorPage ? form.editalId || null : editalAtivo?.id || null) : null,
          pasta: pastaFinal
        };
      }

      if (editId) {
        await api.alterarMaterialApoioMentor(editId, payload);
      } else {
        await api.criarMaterialApoioMentor(payload);
      }

      await carregar();
      setPastaSelecionada(pastaFinal);
      setShowForm(false);
      setEditId(null);
      setArquivo(null);
      setForm(vazioMaterial());
    } catch (x: any) {
      setErro(x.message || 'Não foi possível salvar o material.');
    } finally {
      setOcupado(false);
    }
  };

  const removerMaterial = async (id: string) => {
    if (!confirm('Remover este material de apoio?')) return;
    setOcupado(true);
    setErro('');
    try {
      await api.desativarMaterialApoioMentor(id);
      await carregar();
    } catch (x: any) {
      setErro(x.message || 'Não foi possível remover o material.');
    } finally {
      setOcupado(false);
    }
  };

  const youtubeEmbedUrl = (url?: string) => {
    if (!url) return null;
    try {
      const parsed = new URL(url);
      const host = parsed.hostname.replace('www.', '').toLowerCase();
      if (host === 'youtu.be') {
        return `https://www.youtube.com/embed/${parsed.pathname.slice(1)}`;
      }
      if (host === 'youtube.com' || host === 'm.youtube.com') {
        if (parsed.pathname === '/watch') {
          const id = parsed.searchParams.get('v');
          return id ? `https://www.youtube.com/embed/${id}` : null;
        }
        if (parsed.pathname.startsWith('/embed/') || parsed.pathname.startsWith('/shorts/')) {
          return `https://www.youtube.com/embed/${parsed.pathname.split('/')[2]}`;
        }
      }
    } catch {
      return null;
    }
    return null;
  };

  const labelAbaAtiva = abas.find((a) => a.key === aba)?.label || 'Materiais';

  return (
    <div className="mentor-page">
      <div className="page-heading">
        <div>
          <h1>Materiais de Apoio</h1>
          <p>Organize seus conteúdos por tipo e em pastas temáticas (ex: Português, Direito, etc.).</p>
        </div>
        {isMentorPage && (
          <div className="heading-actions">
            <button className="btn-primary" onClick={() => abrirNovo()}>
              ＋ Novo Material
            </button>
          </div>
        )}
      </div>

      {erro && <div className="form-error">{erro}</div>}

      {/* ABAS SUPERIORES DE TIPOS (SEM A ABA TODOS) */}
      <div className="fc-nav" style={{ flexWrap: 'wrap', gap: '8px', marginBottom: '24px' }}>
        {abas.map((tabItem) => (
          <button
            key={tabItem.key}
            type="button"
            className={aba === tabItem.key ? 'active' : ''}
            onClick={() => {
              setAba(tabItem.key);
              setPastaSelecionada(null); // Volta para a grade de pastas ao trocar de aba
            }}
            style={{ fontSize: '14px', fontWeight: 700, padding: '10px 18px' }}
          >
            {icones[tabItem.key]} {tabItem.label}
          </button>
        ))}
      </div>

      {/* BREADCRUMB E NAVEGAÇÃO DE PASTAS */}
      <div style={{ marginBottom: '20px', display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: '12px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '15px', fontWeight: 700 }}>
          <span style={{ color: 'var(--text2)' }}>{icones[aba]} {labelAbaAtiva}</span>
          {pastaSelecionada && (
            <>
              <span style={{ color: 'var(--text3)' }}>›</span>
              <span style={{ color: 'var(--accent)' }}>
                📁 {pastaSelecionada === '__TODAS__' ? 'Todas as Pastas' : pastaSelecionada}
              </span>
            </>
          )}
        </div>

        {pastaSelecionada ? (
          <div style={{ display: 'flex', gap: '8px' }}>
            <button
              type="button"
              className="btn-secondary"
              onClick={() => setPastaSelecionada(null)}
              style={{ fontSize: '12px', padding: '6px 12px' }}
            >
              ← Voltar para Grade de Pastas
            </button>
            {isMentorPage && (
              <button
                type="button"
                className="btn-primary"
                onClick={() => abrirNovo(pastaSelecionada === '__TODAS__' ? 'Geral' : pastaSelecionada)}
                style={{ fontSize: '12px', padding: '6px 12px' }}
              >
                ＋ Adicionar nesta Pasta
              </button>
            )}
          </div>
        ) : (
          <div style={{ display: 'flex', gap: '8px' }}>
            <button
              type="button"
              className="btn-secondary"
              onClick={() => setPastaSelecionada('__TODAS__')}
              style={{ fontSize: '12px', padding: '6px 12px' }}
            >
              📄 Ver Todos de {labelAbaAtiva}
            </button>
          </div>
        )}
      </div>

      {/* VISUALIZAÇÃO 1: GRADE DE PASTAS (Quando pastaSelecionada === null) */}
      {!pastaSelecionada && (
        <div>
          {pastasDaAba.length === 0 ? (
            <div className="empty-state">
              <h2>Nenhuma pasta criada em {labelAbaAtiva}</h2>
              <p>Cadastre um material para criar sua primeira pasta nesta categoria.</p>
              {isMentorPage && (
                <button type="button" className="btn-primary" onClick={() => abrirNovo()} style={{ marginTop: '12px' }}>
                  ＋ Cadastrar Primeiro Conteúdo
                </button>
              )}
            </div>
          ) : (
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))', gap: '16px' }}>
              {pastasDaAba.map((nomePasta) => {
                const qtd = materiaisDoTipo.filter((m) => (m.pasta?.trim() || 'Geral') === nomePasta).length;
                return (
                  <div
                    key={nomePasta}
                    className="card-base hover-card"
                    onClick={() => setPastaSelecionada(nomePasta)}
                    style={{
                      padding: '20px',
                      borderRadius: '12px',
                      cursor: 'pointer',
                      display: 'flex',
                      alignItems: 'center',
                      gap: '16px',
                      border: '1px solid var(--border)',
                      transition: 'all 0.2s ease',
                      background: 'var(--bg-card)'
                    }}
                  >
                    <div
                      style={{
                        fontSize: '32px',
                        width: '54px',
                        height: '54px',
                        borderRadius: '12px',
                        background: 'rgba(var(--accent-rgb, 99, 102, 241), 0.1)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        flexShrink: 0
                      }}
                    >
                      📁
                    </div>
                    <div style={{ minWidth: 0, flex: 1 }}>
                      <h3 style={{ margin: 0, fontSize: '16px', fontWeight: 700, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                        {nomePasta}
                      </h3>
                      <span style={{ fontSize: '12px', color: 'var(--text3)', display: 'block', marginTop: '4px' }}>
                        {qtd} {qtd === 1 ? 'conteúdo' : 'conteúdos'}
                      </span>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}

      {/* VISUALIZAÇÃO 2: MATERIAIS DA PASTA SELECIONADA (Quando pastaSelecionada !== null) */}
      {pastaSelecionada && (
        <div style={{ marginTop: '10px' }}>
          {materiaisDaPasta.length === 0 ? (
            <div className="empty-state">
              <h2>Nenhum conteúdo nesta pasta</h2>
              <p>{carregando ? 'Carregando materiais…' : 'Esta pasta ainda não possui conteúdos cadastrados.'}</p>
              {isMentorPage && (
                <button
                  type="button"
                  className="btn-primary"
                  onClick={() => abrirNovo(pastaSelecionada === '__TODAS__' ? 'Geral' : pastaSelecionada)}
                  style={{ marginTop: '12px' }}
                >
                  ＋ Adicionar Conteúdo Aqui
                </button>
              )}
            </div>
          ) : (
            <>
            <div style={{ display: 'grid', gap: '16px', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))' }}>
              {materiaisDaPagina.map((material) => (
                <article key={material.id} className="card-base" style={{ padding: '20px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', gap: '14px', alignItems: 'flex-start' }}>
                    <div style={{ display: 'flex', gap: '12px', minWidth: 0, flex: 1 }}>
                      <div style={{ fontSize: '26px', lineHeight: 1 }}>{icones[material.tipo]}</div>
                      <div style={{ minWidth: 0 }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>
                          <h2 style={{ margin: 0, fontSize: '17px', lineHeight: 1.2 }}>{material.titulo}</h2>
                          <span style={{ background: 'var(--bg-tag)', padding: '2px 8px', borderRadius: '4px', color: 'var(--text2)', fontSize: '11px', fontWeight: 600 }}>
                            📁 {material.pasta || 'Geral'}
                          </span>
                          <span style={{ color: 'var(--text3)', fontSize: '11px' }}>
                            {mentor ? (material.escopo === 'global' ? 'Global' : 'Concurso específico') : (material.escopo === 'edital' ? 'Concurso atual' : '')}
                          </span>
                        </div>
                        {material.descricao && <p style={{ margin: '10px 0 0', color: 'var(--text3)', whiteSpace: 'pre-wrap', fontSize: '13px' }}>{material.descricao}</p>}

                        {material.tipo === 'texto' ? (
                          <div style={{ marginTop: '12px', color: 'var(--text)', whiteSpace: 'pre-wrap', fontSize: '13px', background: 'rgba(255,255,255,0.03)', padding: '12px', borderRadius: '8px' }}>
                            {material.texto}
                          </div>
                        ) : material.tipo === 'arquivo' ? (
                          <a
                            href={isMentorPage ? `/api/v1/mentor/materiais-apoio/${encodeURIComponent(material.id)}/arquivo` : `/api/v1/alunos/${encodeURIComponent(alunoId)}/materiais-apoio/${encodeURIComponent(material.id)}/arquivo`}
                            style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', marginTop: '12px', color: 'var(--accent)', fontWeight: 700 }}
                          >
                            💾 Baixar {material.arquivoNome || 'arquivo'}
                          </a>
                        ) : material.tipo === 'youtube' ? (
                          <>
                            {(() => {
                              const embedUrl = youtubeEmbedUrl(material.url || undefined);
                              return embedUrl ? (
                                <div style={{ marginTop: '16px', width: '100%' }}>
                                  <div
                                    style={{
                                      position: 'relative',
                                      width: '100%',
                                      height: 'clamp(220px, 24vw, 280px)',
                                      borderRadius: '14px',
                                      overflow: 'hidden',
                                      border: '1px solid var(--border)',
                                      boxShadow: '0 10px 24px rgba(0, 0, 0, 0.18)',
                                      background: '#000'
                                    }}
                                  >
                                    <iframe
                                      title={material.titulo}
                                      src={embedUrl}
                                      frameBorder="0"
                                      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                                      allowFullScreen
                                      style={{ position: 'absolute', top: 0, left: 0, width: '100%', height: '100%' }}
                                    />
                                  </div>
                                  <div style={{ marginTop: '12px', display: 'flex', gap: '10px', alignItems: 'center', flexWrap: 'wrap' }}>
                                    <button
                                      type="button"
                                      className="btn-secondary"
                                      style={{ fontSize: '12px', padding: '6px 14px', fontWeight: 700 }}
                                      onClick={() => setVideoExpandido({ url: embedUrl, titulo: material.titulo })}
                                    >
                                      🎬 Modo Cinema
                                    </button>
                                    {material.url && (
                                      <a
                                        href={material.url}
                                        target="_blank"
                                        rel="noreferrer"
                                        style={{ fontSize: '12px', color: 'var(--accent)', display: 'inline-flex', alignItems: 'center', gap: '6px', fontWeight: 700 }}
                                      >
                                        ▶ Abrir no YouTube
                                      </a>
                                    )}
                                  </div>
                                </div>
                              ) : (
                                <a href={material.url} target="_blank" rel="noreferrer" style={{ display: 'inline-block', marginTop: '12px', color: 'var(--accent)', fontWeight: 700 }}>
                                  ▶ Abrir vídeo no YouTube
                                </a>
                              );
                            })()}
                          </>
                        ) : (
                          material.url && (
                            <a href={material.url} target="_blank" rel="noreferrer" style={{ display: 'inline-block', marginTop: '12px', color: 'var(--accent)', fontWeight: 700 }}>
                              🔗 Abrir link externo
                            </a>
                          )
                        )}
                      </div>
                    </div>
                    {mentor && (
                      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: '8px' }}>
                        <button className="btn-secondary" type="button" onClick={() => abrirEdicao(material)} style={{ fontSize: '12px', padding: '4px 10px' }}>
                          ✎ Editar
                        </button>
                        <button className="btn-secondary danger" type="button" onClick={() => removerMaterial(material.id)} style={{ fontSize: '12px', padding: '4px 10px' }}>
                          × Excluir
                        </button>
                      </div>
                    )}
                  </div>
                </article>
              ))}
            </div>
            <Pagination pagina={paginaAtual} totalPaginas={totalPaginas} total={materiaisDaPasta.length} onChange={setPagina} />
            </>
          )}
        </div>
      )}

      {/* MODAL DE CADASTRO E EDIÇÃO DE MATERIAL */}
      {mentor && (
        <div
          className="modal-backdrop"
          style={{ display: showForm ? 'block' : 'none' }}
          onClick={() => {
            if (!ocupado) {
              setShowForm(false);
              setEditId(null);
              setForm(vazioMaterial());
              setErro('');
            }
          }}
        >
          <form
            className="modal-card ed-material-modal"
            onClick={(event) => event.stopPropagation()}
            onSubmit={salvarMaterial}
            style={{ maxWidth: '640px', width: '100%' }}
          >
            <div className="modal-heading">
              <div>
                <h2>{editId ? 'Editar Material de Apoio' : 'Novo Material de Apoio'}</h2>
                <p>Cadastre materiais e organize em pastas de estudos.</p>
              </div>
              <button
                type="button"
                className="icon-button"
                onClick={() => {
                  setShowForm(false);
                  setEditId(null);
                  setForm(vazioMaterial());
                  setErro('');
                }}
              >
                ×
              </button>
            </div>

            <div className="form-grid" style={{ display: 'grid', gap: '16px' }}>
              <label>
                <span className="field-label">Título *</span>
                <input
                  className="form-control"
                  placeholder="Ex: Aula 01 — Concordância Verbal"
                  value={form.titulo}
                  onChange={(e) => setForm({ ...form, titulo: e.target.value })}
                  required
                />
              </label>

              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <label>
                  <span className="field-label">Tipo de Conteúdo *</span>
                  <select
                    className="form-control"
                    value={form.tipo}
                    onChange={(e) => setForm({ ...form, tipo: e.target.value as TipoMaterial, url: '', texto: '' })}
                  >
                    <option value="youtube">🎥 Vídeo (YouTube)</option>
                    <option value="arquivo">📎 Arquivo (PDF / Doc)</option>
                    <option value="texto">📝 Texto Informativo</option>
                    <option value="link">🔗 Link Externo</option>
                  </select>
                </label>

                <label>
                  <span className="field-label">Pasta / Disciplina *</span>
                  {!criarNovaPasta ? (
                    <select
                      className="form-control"
                      value={form.pasta}
                      onChange={(e) => {
                        if (e.target.value === '__NOVA__') {
                          setCriarNovaPasta(true);
                          setNomeNovaPasta('');
                        } else {
                          setForm({ ...form, pasta: e.target.value });
                        }
                      }}
                    >
                      {pastasDaAba.map((p) => (
                        <option key={p} value={p}>
                          📁 {p}
                        </option>
                      ))}
                      {!pastasDaAba.includes('Geral') && <option value="Geral">📁 Geral</option>}
                      <option value="__NOVA__">＋ Criar Nova Pasta…</option>
                    </select>
                  ) : (
                    <div style={{ display: 'flex', gap: '6px' }}>
                      <input
                        className="form-control"
                        placeholder="Nome da pasta (ex: Português)"
                        value={nomeNovaPasta}
                        onChange={(e) => setNomeNovaPasta(e.target.value)}
                        autoFocus
                      />
                      <button
                        type="button"
                        className="btn-secondary"
                        style={{ padding: '0 8px', fontSize: '11px' }}
                        onClick={() => setCriarNovaPasta(false)}
                      >
                        Cancelar
                      </button>
                    </div>
                  )}
                </label>
              </div>

              <label>
                <span className="field-label">Escopo *</span>
                <select
                  className="form-control"
                  value={form.escopo}
                  onChange={(e) => setForm({ ...form, escopo: e.target.value as Escopo })}
                >
                  <option value="global">Global (Visível em todos os concursos)</option>
                  <option value="edital" disabled={isMentorPage ? mentorConcursos.length === 0 : !editalAtivo}>
                    {isMentorPage ? 'Específico para um Concurso' : 'Concurso atual'}
                  </option>
                </select>
              </label>

              {form.escopo === 'edital' && isMentorPage && (
                <label>
                  <span className="field-label">Concurso Vinculado *</span>
                  <select
                    className="form-control"
                    value={form.editalId}
                    onChange={(e) => setForm({ ...form, editalId: e.target.value })}
                  >
                    <option value="">Selecione um concurso</option>
                    {mentorConcursos.map((c: any) => {
                      const editalRelacionado = mentorEditais.find((e: any) => e.concursoId === c.id);
                      const targetId = editalRelacionado?.id || c.id;
                      return (
                        <option key={c.id} value={targetId}>
                          {c.nome} {c.banca ? `(${c.banca})` : ''}
                        </option>
                      );
                    })}
                  </select>
                </label>
              )}

              <label>
                <span className="field-label">Descrição / Observações</span>
                <textarea
                  className="form-control"
                  rows={2}
                  placeholder="Resumo do conteúdo ou recomendações de estudo…"
                  value={form.descricao}
                  onChange={(e) => setForm({ ...form, descricao: e.target.value })}
                />
              </label>

              {form.tipo === 'arquivo' ? (
                <label>
                  <span className="field-label">Arquivo *</span>
                  <input className="form-control" type="file" accept="*/*" onChange={(e) => setArquivo(e.target.files?.[0] || null)} />
                  {editId && !arquivo && <small style={{ color: 'var(--text3)', display: 'block', marginTop: '4px' }}>Arquivo atual mantido.</small>}
                </label>
              ) : form.tipo !== 'texto' ? (
                <label>
                  <span className="field-label">{form.tipo === 'youtube' ? 'Link do Vídeo no YouTube *' : 'URL do Link *'}</span>
                  <input
                    className="form-control"
                    type="url"
                    placeholder={form.tipo === 'youtube' ? 'https://www.youtube.com/watch?v=...' : 'https://...'}
                    value={form.url}
                    onChange={(e) => setForm({ ...form, url: e.target.value })}
                    required
                  />
                </label>
              ) : (
                <label>
                  <span className="field-label">Conteúdo do Texto *</span>
                  <textarea
                    className="form-control"
                    rows={6}
                    placeholder="Escreva ou cole o resumo/material em texto aqui…"
                    value={form.texto}
                    onChange={(e) => setForm({ ...form, texto: e.target.value })}
                    required
                  />
                </label>
              )}
            </div>

            <div className="modal-actions" style={{ marginTop: '20px', display: 'flex', justifyContent: 'flex-end', gap: '10px' }}>
              <button
                type="button"
                className="btn-secondary"
                onClick={() => {
                  setShowForm(false);
                  setEditId(null);
                  setForm(vazioMaterial());
                  setErro('');
                }}
              >
                Cancelar
              </button>
              <button className="btn-primary" disabled={ocupado}>
                {ocupado ? 'Salvando…' : editId ? 'Salvar Alterações' : 'Criar Material'}
              </button>
            </div>
          </form>
        </div>
      )}

      {/* MODAL MODO CINEMA PARA YOUTUBE */}
      {videoExpandido && (
        <div className="modal-backdrop" style={{ zIndex: 9999, background: 'rgba(0, 0, 0, 0.88)' }} onClick={() => setVideoExpandido(null)}>
          <div
            style={{
              width: 'min(1280px, 95vw)',
              maxHeight: '94vh',
              display: 'flex',
              flexDirection: 'column',
              gap: '12px'
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', color: '#fff' }}>
              <h2 style={{ fontSize: '18px', margin: 0, fontWeight: 700 }}>🎬 {videoExpandido.titulo}</h2>
              <button
                type="button"
                className="icon-button"
                style={{ background: 'rgba(255, 255, 255, 0.2)', color: '#fff', fontSize: '22px', width: '36px', height: '36px' }}
                onClick={() => setVideoExpandido(null)}
              >
                ×
              </button>
            </div>
            <div style={{ width: '100%', height: 'min(740px, 80vh)', borderRadius: '16px', overflow: 'hidden', border: '1px solid #333', background: '#000', boxShadow: '0 20px 50px rgba(0,0,0,0.8)' }}>
              <iframe
                title={videoExpandido.titulo}
                src={videoExpandido.url}
                frameBorder="0"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                allowFullScreen
                style={{ width: '100%', height: '100%' }}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
