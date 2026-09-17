import React, { useCallback, useEffect, useState } from 'react';
import { api } from '../services/api';
import { Pagination } from '../components/ui/Pagination';

const PROMPT_IA_QUESTOES = `Crie questões de prova no formato JSON estritamente válido.
Responda APENAS com o JSON no seguinte formato:

{
  "questoes": [
    {
      "disciplina": "Direito Constitucional",
      "assunto": "Direitos e Garantias Fundamentais",
      "tipo": "multipla_escolha",
      "enunciado": "Sobre os direitos individuais inscritos na Constituição Federal, assinale a afirmativa correta:",
      "alternativas": [
        "A) É livre a manifestação do pensamento, sendo permitido o anonimato.",
        "B) É assegurado o direito de resposta, proporcional ao agravo, além da indenização.",
        "C) É inviolável a liberdade de consciência e de crença, sendo vedada a prestação de assistência religiosa nas entidades de internação coletiva.",
        "D) A casa é asilo inviolável do indivíduo, nela ninguém podendo penetrar sem consentimento do morador, salvo em caso de flagrante delito ou desastre, ou durante a noite, por determinação judicial."
      ],
      "respostaCorreta": "B",
      "explicacao": "O art. 5º, V da CF/88 estabelece que é assegurado o direito de resposta, proporcional ao agravo, além da indenização por dano material, moral ou à imagem.",
      "alcance": "global"
    }
  ]
}`;

const questaoVazia = {
  disciplina: '',
  assunto: '',
  tipo: 'multipla_escolha',
  enunciado: '',
  alternativas: '',
  respostaCorreta: '',
  explicacao: '',
  alcance: 'global',
  concursoId: ''
};

export const MentorBancoQuestoesPage: React.FC = () => {
  const [questoes, setQuestoes] = useState<any[]>([]);
  const [concursos, setConcursos] = useState<any[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [erro, setErro] = useState('');
  const [sucesso, setSucesso] = useState('');
  const [ocupado, setOcupado] = useState(false);
  const [pagina, setPagina] = useState(1);

  // Filtros
  const [filtroDisciplina, setFiltroDisciplina] = useState('');
  const [filtroAlcance, setFiltroAlcance] = useState<'todos' | 'global' | 'concurso'>('todos');
  const [filtroConcursoId, setFiltroConcursoId] = useState('');

  // Modais
  const [modal, setModal] = useState<'' | 'questao' | 'importar'>('');
  const [editId, setEditId] = useState<string | null>(null);
  const [form, setForm] = useState({ ...questaoVazia });
  const [jsonImport, setJsonImport] = useState('');
  const [promptCopiado, setPromptCopiado] = useState(false);

  const carregarDados = useCallback(async () => {
    setCarregando(true);
    setErro('');
    try {
      const [resQ, resC] = await Promise.all([
        api.listarBancoQuestoes(),
        api.listarConcursosMentor()
      ]);

      const listaQ = Array.isArray(resQ) ? resQ : resQ?.questoes || [];
      const listaC = Array.isArray(resC) ? resC : resC?.concursos || [];

      setQuestoes(listaQ);
      setConcursos(listaC);
    } catch (e: any) {
      setErro(e.message || 'Erro ao carregar banco de questões.');
    } finally {
      setCarregando(false);
    }
  }, []);

  useEffect(() => {
    carregarDados();
  }, [carregarDados]);

  const executar = async (fn: () => Promise<any>, msggSucesso: string) => {
    setOcupado(true);
    setErro('');
    setSucesso('');
    try {
      await fn();
      setModal('');
      setEditId(null);
      setForm({ ...questaoVazia });
      setJsonImport('');
      setSucesso(msggSucesso);
      await carregarDados();
      setTimeout(() => setSucesso(''), 3000);
    } catch (e: any) {
      setErro(e.message || 'Ocorreu um erro ao processar a requisição.');
    } finally {
      setOcupado(false);
    }
  };

  const abrirModalCriar = () => {
    setEditId(null);
    setForm({ ...questaoVazia });
    setErro('');
    setModal('questao');
  };

  const abrirModalEditar = (q: any) => {
    setEditId(q.id);
    setForm({
      disciplina: q.disciplina || '',
      assunto: q.assunto || '',
      tipo: q.tipo || 'multipla_escolha',
      enunciado: q.enunciado || '',
      alternativas: Array.isArray(q.alternativas) ? q.alternativas.join('\n') : (q.alternativas || ''),
      respostaCorreta: q.respostaCorreta || '',
      explicacao: q.explicacao || '',
      alcance: q.alcance || (q.concursoId ? 'concurso' : 'global'),
      concursoId: q.concursoId || ''
    });
    setErro('');
    setModal('questao');
  };

  const salvarQuestao = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.disciplina.trim() || !form.enunciado.trim() || !form.respostaCorreta.trim()) {
      setErro('Informe ao menos a disciplina, o enunciado e a resposta correta.');
      return;
    }

    const alternativasList = form.tipo === 'multipla_escolha'
      ? (typeof form.alternativas === 'string' ? form.alternativas.split('\n').map((x) => x.trim()).filter(Boolean) : form.alternativas)
      : [];

    const payload: any = {
      disciplina: form.disciplina.trim(),
      assunto: form.assunto.trim() || 'Geral',
      tipo: form.tipo,
      enunciado: form.enunciado.trim(),
      alternativas: alternativasList,
      respostaCorreta: form.respostaCorreta.trim(),
      explicacao: form.explicacao.trim() || null,
      alcance: form.alcance,
      concursoId: form.alcance === 'concurso' ? (form.concursoId || null) : null
    };

    if (editId) {
      executar(() => api.alterarQuestaoBanco(editId, payload), 'Questão atualizada com sucesso!');
    } else {
      executar(() => api.criarQuestaoBanco(payload), 'Questão criada com sucesso!');
    }
  };

  const excluirQuestao = (id: string) => {
    if (!confirm('Tem certeza que deseja excluir esta questão do banco de questões?')) return;
    executar(() => api.desativarQuestaoBanco(id), 'Questão excluída com sucesso!');
  };

  const importarViaJson = () => {
    if (!jsonImport.trim()) {
      setErro('Cole o JSON de questões para importar.');
      return;
    }

    let parsed: any;
    try {
      parsed = JSON.parse(jsonImport.trim());
    } catch {
      setErro('O JSON fornecido é inválido. Verifique a sintaxe.');
      return;
    }

    const lista = Array.isArray(parsed) ? parsed : parsed.questoes || [];
    if (!lista.length) {
      setErro('Nenhuma questão encontrada no JSON.');
      return;
    }

    executar(() => api.importarQuestoesBanco({ questoes: lista }), `${lista.length} questão(ões) importada(s) com sucesso!`);
  };

  const copiarPrompt = () => {
    navigator.clipboard.writeText(PROMPT_IA_QUESTOES);
    setPromptCopiado(true);
    setTimeout(() => setPromptCopiado(false), 2000);
  };

  // Filtragem
  const questoesFiltradas = questoes.filter((q) => {
    if (filtroDisciplina && !q.disciplina?.toLowerCase().includes(filtroDisciplina.toLowerCase())) return false;
    if (filtroAlcance === 'global' && q.alcance !== 'global') return false;
    if (filtroAlcance === 'concurso') {
      if (q.alcance !== 'concurso') return false;
      if (filtroConcursoId && q.concursoId !== filtroConcursoId) return false;
    }
    return true;
  });
  const totalPaginas = Math.max(1, Math.ceil(questoesFiltradas.length / 12));
  const paginaAtual = Math.min(pagina, totalPaginas);
  const questoesDaPagina = questoesFiltradas.slice((paginaAtual - 1) * 12, paginaAtual * 12);

  return (
    <div className="mentor-banco-questoes">
      <div className="page-heading" style={{ marginBottom: '24px' }}>
        <div>
          <h2>❓ Banco de Questões da Mentoria</h2>
          <p>Cadastre, importe via IA e gerencie as questões disponibilizadas para os seus alunos.</p>
        </div>
        <div className="heading-actions" style={{ display: 'flex', gap: '8px' }}>
          <button className="btn-secondary" onClick={() => { setModal('importar'); setErro(''); }}>
            ⚡ Importar por IA / JSON
          </button>
          <button className="btn-primary" onClick={abrirModalCriar}>
            ＋ Criar Questão
          </button>
        </div>
      </div>

      {erro && <div className="feedback-banner error" style={{ marginBottom: '20px' }}>{erro}</div>}
      {sucesso && <div className="feedback-banner success" style={{ marginBottom: '20px' }}>{sucesso}</div>}

      {/* FILTROS */}
      <div className="card-base" style={{ padding: '16px', marginBottom: '24px', display: 'flex', gap: '12px', flexWrap: 'wrap', alignItems: 'center' }}>
        <div style={{ flex: 1, minWidth: '200px' }}>
          <label style={{ fontSize: '12px', fontWeight: 700, color: 'var(--text2)', display: 'block', marginBottom: '4px' }}>
            Filtrar por Disciplina
          </label>
          <input
            className="form-control"
            placeholder="ex: Direito Constitucional, Português..."
            value={filtroDisciplina}
            onChange={(e) => setFiltroDisciplina(e.target.value)}
          />
        </div>

        <div style={{ minWidth: '180px' }}>
          <label style={{ fontSize: '12px', fontWeight: 700, color: 'var(--text2)', display: 'block', marginBottom: '4px' }}>
            Disponibilização / Alcance
          </label>
          <select className="form-control" value={filtroAlcance} onChange={(e: any) => setFiltroAlcance(e.target.value)}>
            <option value="todos">🌐 + 🎯 Todos os Alcances</option>
            <option value="global">🌐 Geral (Para todos os alunos)</option>
            <option value="concurso">🎯 Por Concurso Específico</option>
          </select>
        </div>

        {filtroAlcance === 'concurso' && (
          <div style={{ minWidth: '200px' }}>
            <label style={{ fontSize: '12px', fontWeight: 700, color: 'var(--text2)', display: 'block', marginBottom: '4px' }}>
              Selecionar Concurso
            </label>
            <select className="form-control" value={filtroConcursoId} onChange={(e) => setFiltroConcursoId(e.target.value)}>
              <option value="">Todos os Concursos</option>
              {concursos.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.nome}
                </option>
              ))}
            </select>
          </div>
        )}
      </div>

      {/* LISTAGEM DE QUESTÕES */}
      {carregando ? (
        <div style={{ padding: '40px', textAlign: 'center', color: 'var(--text3)' }}>Carregando questões do banco…</div>
      ) : questoesFiltradas.length === 0 ? (
        <div className="empty-state">
          <h2>Nenhuma questão encontrada</h2>
          <p>Nenhuma questão corresponde aos filtros selecionados ou nenhuma questão foi cadastrada ainda.</p>
        </div>
      ) : (
        <>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
          {questoesDaPagina.map((q, idx) => {
            const concursoAssociado = concursos.find((c) => c.id === q.concursoId);

            return (
              <div key={q.id || idx} className="card-base" style={{ padding: '20px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: '16px', marginBottom: '10px' }}>
                  <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap', alignItems: 'center' }}>
                    <span style={{ fontSize: '11px', fontWeight: 800, padding: '3px 8px', borderRadius: '4px', background: 'var(--accent)', color: '#fff' }}>
                      #{(paginaAtual - 1) * 12 + idx + 1}
                    </span>
                    <span style={{ fontSize: '11px', fontWeight: 700, padding: '3px 8px', borderRadius: '4px', background: 'var(--bg-tag)', color: 'var(--text2)' }}>
                      📚 {q.disciplina} {q.assunto ? `• ${q.assunto}` : ''}
                    </span>
                    <span style={{ fontSize: '11px', fontWeight: 700, padding: '3px 8px', borderRadius: '4px', background: q.alcance === 'concurso' ? 'rgba(230, 160, 0, 0.15)' : 'rgba(0, 200, 100, 0.15)', color: q.alcance === 'concurso' ? 'var(--gold)' : 'var(--green)' }}>
                      {q.alcance === 'concurso' ? `🎯 Concurso: ${concursoAssociado?.nome || 'Específico'}` : '🌐 Geral (Todos os Alunos)'}
                    </span>
                  </div>

                  <div style={{ display: 'flex', gap: '6px' }}>
                    <button
                      type="button"
                      className="btn-secondary"
                      style={{ padding: '4px 10px', fontSize: '12px' }}
                      onClick={() => abrirModalEditar(q)}
                    >
                      ✏️ Editar
                    </button>
                    <button
                      type="button"
                      className="btn-secondary danger-text"
                      style={{ padding: '4px 10px', fontSize: '12px' }}
                      onClick={() => excluirQuestao(q.id)}
                    >
                      🗑️ Excluir
                    </button>
                  </div>
                </div>

                <div style={{ fontWeight: 700, fontSize: '15px', color: 'var(--text1)', marginBottom: '12px', lineHeight: 1.5 }}>
                  {q.enunciado}
                </div>

                {/* ALTERNATIVAS OU GABARITO */}
                {q.tipo === 'multipla_escolha' && Array.isArray(q.alternativas) && q.alternativas.length > 0 && (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '6px', marginBottom: '12px', background: 'var(--bg-input)', padding: '12px', borderRadius: '8px' }}>
                    {q.alternativas.map((alt: string, aIdx: number) => {
                      const letra = String.fromCharCode(65 + aIdx);
                      const gabUpper = (q.respostaCorreta || '').trim().toUpperCase();
                      const altTrim = (alt || '').trim();
                      const ehGabarito =
                        gabUpper === letra ||
                        gabUpper === altTrim.toUpperCase() ||
                        (altTrim.toUpperCase().startsWith(`${letra})`) && gabUpper === letra) ||
                        gabUpper.startsWith(`${letra})`);
                      return (
                        <div
                          key={aIdx}
                          style={{
                            fontSize: '13px',
                            color: ehGabarito ? 'var(--green)' : 'var(--text2)',
                            fontWeight: ehGabarito ? 700 : 400
                          }}
                        >
                          {alt} {ehGabarito && ' ✅ (Gabarito)'}
                        </div>
                      );
                    })}
                  </div>
                )}

                <div style={{ fontSize: '13px', color: 'var(--text2)', background: 'var(--bg-card)', padding: '10px 12px', borderRadius: '6px', border: '1px solid var(--border)' }}>
                  ✅ <strong>Gabarito:</strong> {q.respostaCorreta}
                  {q.explicacao && (
                    <div style={{ marginTop: '6px', fontSize: '12px', color: 'var(--text3)' }}>
                      💡 <strong>Comentário:</strong> {q.explicacao}
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
        <Pagination pagina={paginaAtual} totalPaginas={totalPaginas} total={questoesFiltradas.length} onChange={setPagina} />
        </>
      )}

      {/* MODAL DE CRIAR / EDITAR QUESTÃO */}
      {modal === 'questao' && (
        <div className="modal-backdrop" onClick={() => setModal('')}>
          <form
            className="modal-card"
            style={{ maxWidth: '650px' }}
            onClick={(e) => e.stopPropagation()}
            onSubmit={salvarQuestao}
          >
            <div className="modal-heading">
              <h2>{editId ? '✏️ Editar Questão' : '＋ Criar Nova Questão'}</h2>
              <button type="button" className="icon-button" onClick={() => setModal('')}>×</button>
            </div>

            {erro && <div className="feedback-banner error">{erro}</div>}

            <div className="form-fields">
              {/* ALCANCE */}
              <div style={{ background: 'var(--bg-input)', padding: '12px', borderRadius: '8px', marginBottom: '12px' }}>
                <label style={{ fontSize: '12px', fontWeight: 700, color: 'var(--text2)', display: 'block', marginBottom: '6px' }}>
                  Disponibilização / Alcance da Questão *
                </label>
                <div style={{ display: 'flex', gap: '16px' }}>
                  <label style={{ display: 'flex', alignItems: 'center', gap: '6px', cursor: 'pointer', fontSize: '13px' }}>
                    <input
                      type="radio"
                      name="alcance"
                      value="global"
                      checked={form.alcance === 'global'}
                      onChange={() => setForm((x) => ({ ...x, alcance: 'global', concursoId: '' }))}
                    />
                    🌐 Para todos os alunos (Geral / Global)
                  </label>
                  <label style={{ display: 'flex', alignItems: 'center', gap: '6px', cursor: 'pointer', fontSize: '13px' }}>
                    <input
                      type="radio"
                      name="alcance"
                      value="concurso"
                      checked={form.alcance === 'concurso'}
                      onChange={() => setForm((x) => ({ ...x, alcance: 'concurso' }))}
                    />
                    🎯 Somente para um concurso específico
                  </label>
                </div>
              </div>

              {form.alcance === 'concurso' && (
                <label>
                  <span className="field-label">Selecione o Concurso *</span>
                  <select required className="form-control" value={form.concursoId} onChange={(e) => setForm((x) => ({ ...x, concursoId: e.target.value }))}>
                    <option value="">Selecione o Concurso</option>
                    {concursos.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.nome} ({c.banca})
                      </option>
                    ))}
                  </select>
                </label>
              )}

              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <label>
                  <span className="field-label">Disciplina *</span>
                  <input required className="form-control" placeholder="ex: Direito Administrativo" value={form.disciplina} onChange={(e) => setForm((x) => ({ ...x, disciplina: e.target.value }))} />
                </label>
                <label>
                  <span className="field-label">Assunto</span>
                  <input className="form-control" placeholder="ex: Atos Administrativos" value={form.assunto} onChange={(e) => setForm((x) => ({ ...x, assunto: e.target.value }))} />
                </label>
              </div>

              <label>
                <span className="field-label">Tipo de Questão</span>
                <select className="form-control" value={form.tipo} onChange={(e) => setForm((x) => ({ ...x, tipo: e.target.value }))}>
                  <option value="multipla_escolha">Múltipla Escolha (A, B, C, D, E)</option>
                  <option value="certo_errado">Certo / Errado</option>
                </select>
              </label>

              <label>
                <span className="field-label">Enunciado da Questão *</span>
                <textarea required rows={4} className="form-control" placeholder="Digite o enunciado completo da questão..." value={form.enunciado} onChange={(e) => setForm((x) => ({ ...x, enunciado: e.target.value }))} />
              </label>

              {form.tipo === 'multipla_escolha' && (
                <label>
                  <span className="field-label">Alternativas (uma por linha)</span>
                  <textarea rows={4} className="form-control" placeholder={`A) Primeira opção\nB) Segunda opção\nC) Terceira opção\nD) Quarta opção`} value={form.alternativas} onChange={(e) => setForm((x) => ({ ...x, alternativas: e.target.value }))} />
                </label>
              )}

              <label>
                <span className="field-label">Resposta Correta (Gabarito) *</span>
                {form.tipo === 'certo_errado' ? (
                  <select required className="form-control" value={form.respostaCorreta} onChange={(e) => setForm((x) => ({ ...x, respostaCorreta: e.target.value }))}>
                    <option value="">Selecione</option>
                    <option value="Certo">Certo</option>
                    <option value="Errado">Errado</option>
                  </select>
                ) : (
                  <input required className="form-control" placeholder="ex: A, B, C ou texto exato da alternativa" value={form.respostaCorreta} onChange={(e) => setForm((x) => ({ ...x, respostaCorreta: e.target.value }))} />
                )}
              </label>

              <label>
                <span className="field-label">Explicação / Comentário</span>
                <textarea rows={3} className="form-control" placeholder="Fundamentação teórica ou jurídica do gabarito..." value={form.explicacao} onChange={(e) => setForm((x) => ({ ...x, explicacao: e.target.value }))} />
              </label>
            </div>

            <div className="modal-actions">
              <button type="button" className="btn-secondary" onClick={() => setModal('')}>
                Cancelar
              </button>
              <button className="btn-primary" disabled={ocupado}>
                {ocupado ? 'Salvando…' : editId ? 'Salvar Alterações' : 'Criar Questão'}
              </button>
            </div>
          </form>
        </div>
      )}

      {/* MODAL DE IMPORTAÇÃO VIA IA / JSON */}
      {modal === 'importar' && (
        <div className="modal-backdrop" onClick={() => setModal('')}>
          <div
            className="modal-card"
            style={{ maxWidth: '750px' }}
            onClick={(e) => e.stopPropagation()}
          >
            <div className="modal-heading">
              <h2>⚡ Importar Questões via IA / JSON</h2>
              <button type="button" className="icon-button" onClick={() => setModal('')}>×</button>
            </div>

            {erro && <div className="feedback-banner error">{erro}</div>}

            <div style={{ marginBottom: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                <span style={{ fontSize: '13px', fontWeight: 700, color: 'var(--text2)' }}>
                  1. Copie o Prompt abaixo para gerar questões no ChatGPT / Gemini:
                </span>
                <button type="button" className="btn-secondary" style={{ padding: '4px 10px', fontSize: '11px' }} onClick={copiarPrompt}>
                  {promptCopiado ? '✓ Copiado!' : '📋 Copiar Prompt'}
                </button>
              </div>
              <pre style={{ background: 'var(--bg-input)', padding: '12px', borderRadius: '8px', fontSize: '11px', maxHeight: '140px', overflowY: 'auto', border: '1px solid var(--border)', margin: 0 }}>
                {PROMPT_IA_QUESTOES}
              </pre>
            </div>

            <div>
              <span style={{ fontSize: '13px', fontWeight: 700, color: 'var(--text2)', display: 'block', marginBottom: '8px' }}>
                2. Cole o JSON gerado pela IA no campo abaixo:
              </span>
              <textarea
                rows={8}
                className="form-control"
                placeholder='{ "questoes": [ { "disciplina": "...", "enunciado": "...", "respostaCorreta": "..." } ] }'
                value={jsonImport}
                onChange={(e) => setJsonImport(e.target.value)}
              />
            </div>

            <div className="modal-actions" style={{ marginTop: '20px' }}>
              <button type="button" className="btn-secondary" onClick={() => setModal('')}>
                Cancelar
              </button>
              <button type="button" className="btn-primary" disabled={ocupado || !jsonImport.trim()} onClick={importarViaJson}>
                {ocupado ? 'Importando…' : '⚡ Importar Questões'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
