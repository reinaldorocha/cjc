import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { useData } from '../context/DataContext';
import { api } from '../services/api';
import { ModalConfirmacao } from '../components/ModalConfirmacao';

type Caderno = {
  id: string;
  alunoId: string;
  titulo: string;
  pasta: string;
  conteudo: string;
  editalId?: string;
  materiaId?: string;
  topicoId?: string;
  cor: string;
  ativo: boolean;
  criadoEm: string;
  atualizadoEm: string;
};

const PALETA_CORES = [
  { id: '#4f8ef7', label: 'Azul' },
  { id: '#7c5cfc', label: 'Roxo' },
  { id: '#3ecf8e', label: 'Verde' },
  { id: '#f5c842', label: 'Amarelo' },
  { id: '#f5874a', label: 'Laranja' },
  { id: '#f55a5a', label: 'Vermelho' },
  { id: '#00bcd4', label: 'Ciano' }
];

export const CadernosPage: React.FC = () => {
  const { alunoId, editalAtivo } = useData();



  const [cadernos, setCadernos] = useState<Caderno[]>([]);
  const [pastaAtiva, setPastaAtiva] = useState<string>('__TODAS__');
  const [busca, setBusca] = useState('');
  const [carregando, setCarregando] = useState(false);
  const [erro, setErro] = useState('');
  const [ocupado, setOcupado] = useState(false);

  // Estado do Editor Modal / Caderno Selecionado
  const [modalAberto, setModalAberto] = useState(false);
  const [editId, setEditId] = useState<string | null>(null);
  const [form, setForm] = useState({
    titulo: '',
    pasta: 'Geral',
    conteudo: '',
    cor: '#4f8ef7'
  });
  const [criarNovaPasta, setCriarNovaPasta] = useState(false);
  const [nomeNovaPasta, setNomeNovaPasta] = useState('');

  const carregar = useCallback(async () => {
    setErro('');
    setCarregando(true);
    try {
      const res = await api.listarCadernos(alunoId, editalAtivo?.id || '');
      setCadernos(res.cadernos || []);
    } catch (x: any) {
      setErro(x.message || 'Não foi possível carregar os cadernos de anotações.');
    } finally {
      setCarregando(false);
    }
  }, [alunoId, editalAtivo?.id]);

  useEffect(() => {
    carregar();
  }, [carregar]);

  // Lista única de pastas disponíveis
  const pastasUnicas = useMemo(() => {
    const conjunto = new Set<string>();
    cadernos.forEach((c) => {
      if (c.pasta?.trim()) conjunto.add(c.pasta.trim());
    });
    if (conjunto.size === 0) conjunto.add('Geral');
    return Array.from(conjunto).sort();
  }, [cadernos]);

  // Cadernos filtrados por pasta e busca
  const cadernosFiltrados = useMemo(() => {
    return cadernos.filter((c) => {
      const atendePasta = pastaAtiva === '__TODAS__' || (c.pasta?.trim() || 'Geral') === pastaAtiva;
      const q = busca.toLowerCase().trim();
      const atendeBusca = !q || c.titulo.toLowerCase().includes(q) || c.conteudo.toLowerCase().includes(q) || c.pasta.toLowerCase().includes(q);
      return atendePasta && atendeBusca;
    });
  }, [cadernos, pastaAtiva, busca]);

  const abrirNovo = (pastaPadrao?: string) => {
    setEditId(null);
    setCriarNovaPasta(false);
    setNomeNovaPasta('');
    setForm({
      titulo: '',
      pasta: pastaPadrao || (pastaAtiva !== '__TODAS__' ? pastaAtiva : pastasUnicas[0] || 'Geral'),
      conteudo: '',
      cor: '#4f8ef7'
    });
    setModalAberto(true);
    setErro('');
  };

  const abrirEdicao = (caderno: Caderno) => {
    setEditId(caderno.id);
    setCriarNovaPasta(false);
    setNomeNovaPasta('');
    setForm({
      titulo: caderno.titulo,
      pasta: caderno.pasta || 'Geral',
      conteudo: caderno.conteudo || '',
      cor: caderno.cor || '#4f8ef7'
    });
    setModalAberto(true);
    setErro('');
  };

  const salvarCaderno = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.titulo.trim()) {
      setErro('Informe um título para o caderno.');
      return;
    }
    const pastaFinal = criarNovaPasta ? (nomeNovaPasta.trim() || 'Geral') : (form.pasta.trim() || 'Geral');

    setOcupado(true);
    setErro('');
    try {
      const payload = {
        titulo: form.titulo.trim(),
        pasta: pastaFinal,
        conteudo: form.conteudo,
        cor: form.cor,
        editalId: editalAtivo?.id || null
      };

      if (editId) {
        await api.alterarCaderno(alunoId, editId, payload);
      } else {
        await api.criarCaderno(alunoId, payload);
      }

      await carregar();
      setModalAberto(false);
      setEditId(null);
    } catch (x: any) {
      setErro(x.message || 'Não foi possível salvar o caderno.');
    } finally {
      setOcupado(false);
    }
  };

  const [idParaDeletar, setIdParaDeletar] = useState<string | null>(null);

  const removerCaderno = (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    setIdParaDeletar(id);
  };

  const confirmarExclusaoCaderno = async () => {
    if (!idParaDeletar) return;
    setOcupado(true);
    try {
      await api.desativarCaderno(alunoId, idParaDeletar);
      setIdParaDeletar(null);
      await carregar();
    } catch (x: any) {
      setErro(x.message || 'Não foi possível excluir o caderno.');
    } finally {
      setOcupado(false);
    }
  };

  // Inserção de formatação rápida de texto livre no editor
  const inserirFormatacao = (prefixo: string, sufixo: string = '') => {
    const area = document.getElementById('caderno-textarea') as HTMLTextAreaElement;
    if (!area) return;
    const inicio = area.selectionStart;
    const fim = area.selectionEnd;
    const textoAtual = form.conteudo;
    const selecionado = textoAtual.substring(inicio, fim) || 'Texto';
    const novoTexto = textoAtual.substring(0, inicio) + prefixo + selecionado + sufixo + textoAtual.substring(fim);
    setForm({ ...form, conteudo: novoTexto });
    setTimeout(() => {
      area.focus();
      area.setSelectionRange(inicio + prefixo.length, fim + prefixo.length);
    }, 50);
  };

  const copiarTexto = (caderno: Caderno, e: React.MouseEvent) => {
    e.stopPropagation();
    navigator.clipboard.writeText(`${caderno.titulo}\n\n${caderno.conteudo}`);
    alert('Anotação copiada para a área de transferência!');
  };

  return (
    <div className="mentor-page" style={{ paddingBottom: '60px' }}>
      <div className="page-heading">
        <div>
          <h1>📖 Cadernos de Anotações & Resumos</h1>
          <p>Organize suas anotações, mapas mentais e resumos de estudo por matérias e pastas.</p>
        </div>
        <div className="heading-actions">
          <button className="btn-primary" onClick={() => abrirNovo()}>
            ＋ Novo Caderno
          </button>
        </div>
      </div>

      {erro && <div className="form-error">{erro}</div>}

      {/* BARRA DE PESQUISA E SELEÇÃO DE PASTAS */}
      <div style={{ display: 'flex', gap: '12px', marginBottom: '24px', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between' }}>
        <div className="fc-nav" style={{ gap: '6px', flexWrap: 'wrap', margin: 0 }}>
          <button
            type="button"
            className={pastaAtiva === '__TODAS__' ? 'active' : ''}
            onClick={() => setPastaAtiva('__TODAS__')}
            style={{ fontSize: '13px', padding: '8px 14px' }}
          >
            📚 Todas ({cadernos.length})
          </button>
          {pastasUnicas.map((p) => {
            const qtd = cadernos.filter((c) => (c.pasta?.trim() || 'Geral') === p).length;
            return (
              <button
                key={p}
                type="button"
                className={pastaAtiva === p ? 'active' : ''}
                onClick={() => setPastaAtiva(p)}
                style={{ fontSize: '13px', padding: '8px 14px' }}
              >
                📁 {p} ({qtd})
              </button>
            );
          })}
        </div>

        <div style={{ minWidth: '240px', flex: '0 1 300px' }}>
          <input
            className="form-control"
            placeholder="🔍 Buscar anotação ou resumo…"
            value={busca}
            onChange={(e) => setBusca(e.target.value)}
            style={{ borderRadius: '20px', padding: '8px 16px', fontSize: '13px' }}
          />
        </div>
      </div>

      {/* LISTAGEM EM CARDS DE CADERNOS */}
      {cadernosFiltrados.length === 0 ? (
        <div className="empty-state">
          <div style={{ fontSize: '48px', marginBottom: '12px' }}>📖</div>
          <h2>{busca ? 'Nenhuma anotação encontrada para essa busca' : 'Nenhum caderno criado nesta pasta'}</h2>
          <p>{carregando ? 'Carregando cadernos…' : 'Crie seu primeiro caderno de anotações para resumir a matéria de estudos.'}</p>
          <button type="button" className="btn-primary" onClick={() => abrirNovo(pastaAtiva === '__TODAS__' ? 'Geral' : pastaAtiva)} style={{ marginTop: '14px' }}>
            ＋ Criar Caderno Agora
          </button>
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(310px, 1fr))', gap: '18px' }}>
          {cadernosFiltrados.map((item) => {
            const palaCount = item.conteudo.trim() ? item.conteudo.trim().split(/\s+/).length : 0;
            return (
              <article
                key={item.id}
                className="card-base hover-card"
                onClick={() => abrirEdicao(item)}
                style={{
                  padding: '20px',
                  borderRadius: '14px',
                  cursor: 'pointer',
                  borderTop: `4px solid ${item.cor || '#4f8ef7'}`,
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'space-between',
                  minHeight: '190px'
                }}
              >
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: '8px', marginBottom: '8px' }}>
                    <span
                      style={{
                        background: 'var(--bg-tag)',
                        color: 'var(--text2)',
                        fontSize: '11px',
                        fontWeight: 700,
                        padding: '3px 8px',
                        borderRadius: '6px'
                      }}
                    >
                      📁 {item.pasta || 'Geral'}
                    </span>
                    <span style={{ fontSize: '11px', color: 'var(--text3)' }}>
                      {palaCount} {palaCount === 1 ? 'palavra' : 'palavras'}
                    </span>
                  </div>

                  <h3 style={{ margin: '4px 0 8px', fontSize: '17px', fontWeight: 800, color: 'var(--text)', lineHeight: 1.3 }}>
                    {item.titulo}
                  </h3>

                  <p
                    style={{
                      fontSize: '13px',
                      color: 'var(--text3)',
                      margin: 0,
                      display: '-webkit-box',
                      WebkitLineClamp: 4,
                      WebkitBoxOrient: 'vertical',
                      overflow: 'hidden',
                      whiteSpace: 'pre-wrap',
                      lineHeight: 1.5
                    }}
                  >
                    {item.conteudo ? (
                      item.conteudo
                        .replace(/^#+\s+/gm, '')
                        .replace(/(\*\*|__)(.*?)\1/g, '$2')
                        .replace(/(\*|_)(.*?)\1/g, '$2')
                        .replace(/`{1,3}[^`]*`{1,3}/g, '')
                        .replace(/•\s*•+/g, '•')
                        .replace(/^\s*[-*+]\s+/gm, '• ')
                        .replace(/^\s*>\s+/gm, '')
                        .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
                        .trim()
                    ) : 'Sem conteúdo inserido ainda. Clique para editar…'}
                  </p>
                </div>

                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    marginTop: '16px',
                    paddingTop: '12px',
                    borderTop: '1px solid var(--border)'
                  }}
                >
                  <span style={{ fontSize: '11px', color: 'var(--text3)' }}>
                    🕒 {item.atualizadoEm ? new Date(item.atualizadoEm).toLocaleDateString('pt-BR') : 'Hoje'}
                  </span>
                  <div style={{ display: 'flex', gap: '6px' }}>
                    <button
                      type="button"
                      className="btn-secondary"
                      onClick={(e) => copiarTexto(item, e)}
                      style={{ padding: '3px 8px', fontSize: '11px' }}
                      title="Copiar resumo"
                    >
                      📋 Copiar
                    </button>
                    <button
                      type="button"
                      className="btn-secondary danger-text"
                      onClick={(e) => removerCaderno(item.id, e)}
                      style={{ padding: '3px 8px', fontSize: '11px' }}
                      title="Excluir caderno"
                    >
                      🗑
                    </button>
                  </div>
                </div>
              </article>
            );
          })}
        </div>
      )}

      {/* MODAL EDITOR DE CADERNO / TEXTO LIVRE */}
      {modalAberto && (
        <div
          className="modal-backdrop"
          style={{ zIndex: 9999 }}
          onClick={() => {
            if (!ocupado) setModalAberto(false);
          }}
        >
          <form
            className="modal-card"
            onClick={(e) => e.stopPropagation()}
            onSubmit={salvarCaderno}
            style={{
              maxWidth: '96vw',
              width: '96vw',
              maxHeight: '94vh',
              height: '92vh',
              display: 'flex',
              flexDirection: 'column'
            }}
          >
            <div className="modal-heading" style={{ borderBottom: '1px solid var(--border)', paddingBottom: '12px' }}>
              <div>
                <h2>{editId ? '📖 Editar Caderno de Anotações' : '📖 Novo Caderno de Anotações'}</h2>
                <p>Escreva resumos, anotações de aulas ou pontos chaves do edital.</p>
              </div>
              <button type="button" className="icon-button" onClick={() => setModalAberto(false)}>
                ×
              </button>
            </div>

            <div style={{ display: 'grid', gap: '14px', marginTop: '16px', flex: 1, overflowY: 'auto', paddingRight: '4px' }}>
              <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '12px' }}>
                <label>
                  <span className="field-label">Título do Caderno *</span>
                  <input
                    className="form-control"
                    placeholder="ex: Resumo - Concordância Verbal e Nominal"
                    value={form.titulo}
                    onChange={(e) => setForm({ ...form, titulo: e.target.value })}
                    required
                    style={{ fontSize: '15px', fontWeight: 700 }}
                  />
                </label>

                <label>
                  <span className="field-label">Pasta / Categoria *</span>
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
                      {pastasUnicas.map((p) => (
                        <option key={p} value={p}>
                          📁 {p}
                        </option>
                      ))}
                      {!pastasUnicas.includes('Geral') && <option value="Geral">📁 Geral</option>}
                      <option value="__NOVA__">＋ Criar Nova Pasta…</option>
                    </select>
                  ) : (
                    <div style={{ display: 'flex', gap: '4px' }}>
                      <input
                        className="form-control"
                        placeholder="Nome da pasta (ex: Português)"
                        value={nomeNovaPasta}
                        onChange={(e) => setNomeNovaPasta(e.target.value)}
                        autoFocus
                      />
                      <button type="button" className="btn-secondary" style={{ padding: '0 6px', fontSize: '10px' }} onClick={() => setCriarNovaPasta(false)}>
                        ×
                      </button>
                    </div>
                  )}
                </label>
              </div>

              {/* BARRA DE CORES E FERRAMENTAS DO EDITOR */}
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '12px', flexWrap: 'wrap' }}>
                <div style={{ display: 'flex', gap: '6px', alignItems: 'center' }}>
                  <span style={{ fontSize: '12px', color: 'var(--text3)', fontWeight: 600 }}>Cor do Marcador:</span>
                  {PALETA_CORES.map((c) => (
                    <button
                      key={c.id}
                      type="button"
                      onClick={() => setForm({ ...form, cor: c.id })}
                      style={{
                        width: '22px',
                        height: '22px',
                        borderRadius: '50%',
                        background: c.id,
                        border: form.cor === c.id ? '2px solid #fff' : 'none',
                        boxShadow: form.cor === c.id ? '0 0 0 2px ' + c.id : 'none',
                        cursor: 'pointer'
                      }}
                      title={c.label}
                    />
                  ))}
                </div>

                <div style={{ display: 'flex', gap: '4px', background: 'var(--bg-tag)', padding: '4px', borderRadius: '8px', flexWrap: 'wrap' }}>
                  <button type="button" className="btn-secondary" style={{ padding: '3px 8px', fontSize: '12px', fontWeight: 800 }} onClick={() => inserirFormatacao('**', '**')}>
                    B
                  </button>
                  <button type="button" className="btn-secondary" style={{ padding: '3px 8px', fontSize: '12px', fontStyle: 'italic' }} onClick={() => inserirFormatacao('*', '*')}>
                    I
                  </button>
                  <button type="button" className="btn-secondary" style={{ padding: '3px 8px', fontSize: '12px' }} onClick={() => inserirFormatacao('# ')}>
                    H1
                  </button>
                  <button type="button" className="btn-secondary" style={{ padding: '3px 8px', fontSize: '12px' }} onClick={() => inserirFormatacao('## ')}>
                    H2
                  </button>
                  <button type="button" className="btn-secondary" style={{ padding: '3px 8px', fontSize: '12px' }} onClick={() => inserirFormatacao('• ')}>
                    • Lista
                  </button>
                  <button type="button" className="btn-secondary" style={{ padding: '3px 8px', fontSize: '12px' }} onClick={() => inserirFormatacao('> ')}>
                    " Citação
                  </button>
                  <button type="button" className="btn-secondary" style={{ padding: '3px 8px', fontSize: '12px' }} onClick={() => inserirFormatacao('📌 ATENÇÃO: ')}>
                    📌 Alerta
                  </button>
                </div>
              </div>

              {/* EDITOR DE TEXTO LIVRE */}
              <label style={{ display: 'flex', flexDirection: 'column', flex: 1, minHeight: '280px' }}>
                <span className="field-label" style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span>Conteúdo / Resumo Livre</span>
                  <span style={{ fontSize: '11px', color: 'var(--text3)' }}>
                    {form.conteudo.length} caracteres · {form.conteudo.trim() ? form.conteudo.trim().split(/\s+/).length : 0} palavras
                  </span>
                </span>
                <textarea
                  id="caderno-textarea"
                  className="form-control"
                  rows={18}
                  placeholder="Escreva ou cole seu resumo aqui. Utilize tópicos, esquemas, artigos de lei e observações de questões…"
                  value={form.conteudo}
                  onChange={(e) => setForm({ ...form, conteudo: e.target.value })}
                  style={{
                    fontFamily: 'system-ui, -apple-system, sans-serif',
                    fontSize: '14px',
                    lineHeight: '1.6',
                    flex: 1,
                    resize: 'vertical',
                    padding: '14px'
                  }}
                />
              </label>
            </div>

            <div className="modal-actions" style={{ marginTop: '16px', display: 'flex', justifyContent: 'flex-end', gap: '10px' }}>
              <button type="button" className="btn-secondary" onClick={() => setModalAberto(false)}>
                Cancelar
              </button>
              <button className="btn-primary" disabled={ocupado}>
                {ocupado ? 'Salvando…' : editId ? 'Salvar Caderno' : 'Criar Caderno'}
              </button>
            </div>
          </form>
        </div>
      )}

      <ModalConfirmacao
        aberto={!!idParaDeletar}
        titulo="Confirmar exclusão"
        mensagem="Tem certeza que deseja apagar?"
        textoConfirmar="Sim, apagar"
        textoCancelar="Cancelar"
        aoConfirmar={confirmarExclusaoCaderno}
        aoCancelar={() => setIdParaDeletar(null)}
      />
    </div>
  );
};
