import React, { useState } from 'react';
import { useNavigate, Link, useParams } from 'react-router-dom';
import { api } from '../services/api';
import { useData } from '../context/DataContext';
import { useAutenticacao } from '../context/AutenticacaoContext';

const PROMPT_IA_EXEMPLO = `Você é um especialista em estruturação de editais para concursos públicos.
Por favor, analise o edital/conteúdo programático fornecido abaixo e extraia todas as disciplinas (matérias), tópicos e subtópicos no formato JSON estritamente válido, seguindo exatamente este modelo:

{
  "nome": "Nome do Concurso / Órgão",
  "banca": "Nome da Banca Organizadora",
  "cargo": "Nome do Cargo",
  "salario": 12500,
  "dataProva": "2026-10-15",
  "preEdital": false,
  "prazosRevisao": "1,7,30",
  "materias": [
    {
      "nome": "Língua Portuguesa",
      "peso": 1,
      "topicos": [
        {
          "nome": "Compreensão e interpretação de textos",
          "subtopicos": [
            "Tipologia textual",
            "Mecanismos de coesão textual"
          ]
        },
        {
          "nome": "Ortografia oficial",
          "subtopicos": [
            "Acentuação gráfica"
          ]
        }
      ]
    },
    {
      "nome": "Direito Constitucional",
      "peso": 1,
      "topicos": [
        {
          "nome": "Direitos e garantias fundamentais",
          "subtopicos": [
            "Direitos e deveres individuais e coletivos"
          ]
        }
      ]
    }
  ]
}

Regras Obrigatórias:
1. Retorne APENAS o código JSON puro, sem blocos de explicações ou saudações.
2. Mantenha os nomes dos tópicos e subtópicos limpos e organizados hierarquicamente.
3. Se a data da prova não for conhecida, coloque preEdital como true e dataProva como "".

Segue o conteúdo programático do edital:
[COLE AQUI O TEXTO DO SEU EDITAL OU EDITAL EM PDF]`;

const JSON_EXEMPLO = {
  nome: 'Polícia Federal - Agente',
  banca: 'Cebraspe',
  cargo: 'Agente de Polícia Federal',
  salario: 12522.5,
  dataProva: '',
  preEdital: true,
  prazosRevisao: '1,7,30',
  materias: [
    {
      nome: 'Língua Portuguesa',
      peso: 1,
      topicos: [
        {
          nome: 'Compreensão e interpretação de textos',
          subtopicos: ['Tipologia textual', 'Mecanismos de coesão textual']
        },
        {
          nome: 'Domínio da estrutura morfossintática',
          subtopicos: ['Relações de regência e concordância', 'Emprego do sinal indicativo de crase']
        }
      ]
    },
    {
      nome: 'Raciocínio Lógico',
      peso: 1,
      topicos: [
        {
          nome: 'Lógica de argumentação',
          subtopicos: ['Diagramas lógicos', 'Tabelas verdade']
        }
      ]
    },
    {
      nome: 'Direito Administrativo',
      peso: 1,
      topicos: [
        {
          nome: 'Noções de organização administrativa',
          subtopicos: ['Administração direta e indireta', 'Agentes públicos']
        }
      ]
    }
  ]
};

export const NovoConcursoPage: React.FC = () => {
  const navigate = useNavigate();
  const { concursoId: concursoIdParam } = useParams<{ concursoId?: string }>();
  const modoEdicao = !!concursoIdParam;

  const { alunoId, setActiveContestId, activeContestId, recarregarConcursos, getArray } = useData();
  const { usuario } = useAutenticacao();

  const [modo, setModo] = useState<'json' | 'manual'>('manual');
  const [erro, setErro] = useState('');
  const [sucesso, setSucesso] = useState('');
  const [salvando, setSalvando] = useState(false);
  const [carregando, setCarregando] = useState(modoEdicao);
  const [copiadoPrompt, setCopiadoPrompt] = useState(false);
  const [copiadoExemplo, setCopiadoExemplo] = useState(false);

  const [mostrarPrompt, setMostrarPrompt] = useState(true);
  const [mostrarExemplo, setMostrarExemplo] = useState(false);

  // Estado para modo JSON
  const [jsonText, setJsonText] = useState('');

  // Estado para modo Manual
  const [formManual, setFormManual] = useState({
    nome: '',
    banca: '',
    cargo: '',
    salario: '',
    dataProva: '',
    preEdital: true,
    prazosRevisao: '1,7,30',
    materias: [
      {
        nome: 'Língua Portuguesa',
        peso: 1,
        topicos: [
          {
            nome: 'Compreensão de Texto',
            subtopicos: ['Tipos de Texto']
          }
        ]
      }
    ]
  });

  // Carrega dados existentes em modo edição
  React.useEffect(() => {
    if (!modoEdicao || !concursoIdParam) return;
    const carregar = async () => {
      setCarregando(true);
      try {
        let c: any = null;

        // 1. Tenta buscar da lista local de concursos primeiro
        const concursosLocais = getArray('concursos');
        c = concursosLocais.find((x: any) => x.id === concursoIdParam);

        // 2. Se não achou localmente (ex: mentor no catálogo ou página recarregada), busca da API
        if (!c) {
          if (usuario?.papel === 'mentor') {
            const resC: any = await api.listarConcursosMentor().catch(() => null);
            const lista = Array.isArray(resC) ? resC : resC?.concursos || [];
            c = lista.find((x: any) => x.id === concursoIdParam);
          }
          if (!c) {
            const resC: any = await api.listarConcursos(alunoId).catch(() => null);
            const lista = Array.isArray(resC) ? resC : resC?.concursos || [];
            c = lista.find((x: any) => x.id === concursoIdParam);
          }
        }

        // 3. Busca o edital existente (mentor ou aluno)
        let edital: any = null;
        if (usuario?.papel === 'mentor') {
          const resE = await api.listarEditaisMentor(concursoIdParam).catch(() => null);
          edital = (resE?.editais || resE)?.[0] || null;
        }
        if (!edital) {
          const resE = await api.listarEditais(alunoId, concursoIdParam).catch(() => null);
          const lista = Array.isArray(resE) ? resE : resE?.editais || (resE ? [resE] : []);
          edital = lista[0] || null;
        }

        if (c) {
          const materias = (edital?.materias || []).map((m: any) => ({
            nome: m.nome || '',
            peso: m.peso || 1,
            topicos: (m.topicos || []).map((t: any) => ({
              nome: t.nome || '',
              subtopicos: (t.subtopicos || []).map((s: any) =>
                typeof s === 'string' ? s : (s.nome || '')
              )
            }))
          }));

          setFormManual({
            nome: c.nome || '',
            banca: c.banca || '',
            cargo: c.cargo || '',
            salario: c.salario != null && c.salario !== '' ? String(c.salario) : '',
            dataProva: c.dataProva || '',
            preEdital: !!c.preEdital,
            prazosRevisao: c.prazosRevisao || '1,7,30',
            materias: materias.length > 0 ? materias : [{ nome: '', peso: 1, topicos: [{ nome: '', subtopicos: [''] }] }]
          });

          // Também gera o JSON formatado caso o usuário prefira editar via JSON
          const jsonVal = {
            nome: c.nome || '',
            banca: c.banca || '',
            cargo: c.cargo || '',
            salario: c.salario ? Number(c.salario) : undefined,
            dataProva: c.dataProva || '',
            preEdital: !!c.preEdital,
            prazosRevisao: c.prazosRevisao || '1,7,30',
            materias
          };
          setJsonText(JSON.stringify(jsonVal, null, 2));

          setModo('manual');
        } else {
          setErro('Concurso não encontrado.');
        }
      } catch (e: any) {
        setErro('Erro ao carregar dados do concurso.');
      } finally {
        setCarregando(false);
      }
    };
    carregar();
  }, [concursoIdParam, modoEdicao, alunoId, usuario?.papel]);

  const copiarTexto = (texto: string, setFn: (v: boolean) => void) => {
    navigator.clipboard.writeText(texto);
    setFn(true);
    setTimeout(() => setFn(false), 2000);
  };

  const preencherComExemplo = () => {
    setJsonText(JSON.stringify(JSON_EXEMPLO, null, 2));
    setErro('');
  };

  const salvarViaJson = async () => {
    if (!jsonText.trim()) {
      setErro('Cole o JSON do concurso no campo abaixo para prosseguir.');
      return;
    }

    let payload: any;
    try {
      payload = JSON.parse(jsonText.trim());
    } catch (e: any) {
      setErro('O texto informado não é um JSON válido. Verifique a sintaxe e tente novamente.');
      return;
    }

    if (!payload.nome || !payload.banca) {
      setErro('O JSON deve conter ao menos os campos "nome" e "banca".');
      return;
    }

    setSalvando(true);
    setErro('');
    try {
      const dadosConcurso = {
        nome: String(payload.nome).trim(),
        banca: String(payload.banca).trim(),
        cargo: payload.cargo ? String(payload.cargo).trim() : '',
        salario: payload.salario ? Number(payload.salario) : null,
        limparSalario: !payload.salario,
        dataProva: payload.preEdital ? '' : payload.dataProva || '',
        preEdital: Boolean(payload.preEdital),
        prazosRevisao: payload.prazosRevisao || '1,7,30'
      };

      let concursoId = concursoIdParam || '';

      if (modoEdicao && concursoIdParam) {
        if (usuario?.papel === 'mentor') {
          await api.alterarConcurso('eu', concursoIdParam, dadosConcurso);
        } else {
          await api.alterarConcurso(alunoId, concursoIdParam, dadosConcurso);
        }
      } else {
        if (usuario?.papel === 'mentor') {
          const criado = await api.criarConcursoMentor(dadosConcurso);
          concursoId = criado.id;
        } else {
          const criado = await api.criarConcurso(alunoId, dadosConcurso);
          concursoId = criado.id;
        }
      }

      // Se houver matérias no JSON, salvar edital
      if (Array.isArray(payload.materias) && payload.materias.length > 0 && concursoId) {
        const materiasConvertidas = payload.materias.map((m: any, idx: number) => ({
          nome: String(m.nome || `Matéria ${idx + 1}`).trim(),
          ordem: idx + 1,
          peso: Number(m.peso) || 1,
          topicos: Array.isArray(m.topicos)
            ? m.topicos.map((t: any, tIdx: number) => ({
              nome: String(t.nome || `Tópico ${tIdx + 1}`).trim(),
              ordem: tIdx + 1,
              subtopicos: Array.isArray(t.subtopicos)
                ? t.subtopicos.map((s: any, sIdx: number) => ({
                  nome: typeof s === 'string' ? s.trim() : String(s?.nome || `Subtópico ${sIdx + 1}`).trim(),
                  ordem: sIdx + 1
                }))
                : []
            }))
            : []
        }));

        if (usuario?.papel === 'mentor') {
          await api.criarEditalMentor({
            concursoId,
            nome: `Edital ${dadosConcurso.nome}`,
            versao: '1.0',
            materias: materiasConvertidas
          }).catch(() => undefined);
        } else {
          await api.atribuirEdital(alunoId, {
            concursoId,
            nome: `Edital ${dadosConcurso.nome}`,
            materias: materiasConvertidas
          }).catch(() => undefined);
        }
      }

      await recarregarConcursos();
      setActiveContestId(concursoId);
      setSucesso(modoEdicao ? 'Concurso atualizado com sucesso!' : 'Concurso e edital cadastrados com sucesso!');

      setTimeout(() => {
        if (usuario?.papel === 'mentor') {
          navigate('/mentor/concursos-catalogo');
        } else {
          navigate('/dashboard');
        }
      }, 1000);
    } catch (err: any) {
      setErro(err.message || 'Erro ao salvar concurso.');
    } finally {
      setSalvando(false);
    }
  };

  const salvarViaManual = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formManual.nome.trim() || !formManual.banca.trim()) {
      setErro('Informe ao menos o Nome do Concurso e a Banca Organizadora.');
      return;
    }

    setSalvando(true);
    setErro('');
    try {
      const dadosConcurso = {
        nome: formManual.nome.trim(),
        banca: formManual.banca.trim(),
        cargo: formManual.cargo.trim(),
        salario: formManual.salario ? Number(formManual.salario) : null,
        limparSalario: !formManual.salario,
        dataProva: formManual.preEdital ? '' : formManual.dataProva || '',
        preEdital: formManual.preEdital,
        prazosRevisao: formManual.prazosRevisao || '1,7,30'
      };

      const materiasConvertidas = formManual.materias
        .filter(m => m.nome.trim())
        .map((m, idx) => ({
          nome: m.nome.trim(),
          ordem: idx + 1,
          peso: m.peso || 1,
          topicos: m.topicos.map((t, tIdx) => ({
            nome: t.nome.trim() || `Tópico ${tIdx + 1}`,
            ordem: tIdx + 1,
            subtopicos: t.subtopicos.map((s, sIdx) => ({
              nome: s.trim() || `Subtópico ${sIdx + 1}`,
              ordem: sIdx + 1
            }))
          }))
        }));

      let concursoId = concursoIdParam || '';

      if (modoEdicao && concursoIdParam) {
        // MODO EDIÇÃO
        if (usuario?.papel === 'mentor') {
          await api.alterarConcurso('eu', concursoIdParam, dadosConcurso);
        } else {
          await api.alterarConcurso(alunoId, concursoIdParam, dadosConcurso);
        }
        if (materiasConvertidas.length > 0) {
          if (usuario?.papel === 'mentor') {
            await api.criarEditalMentor({
              concursoId: concursoIdParam,
              nome: `Edital ${dadosConcurso.nome}`,
              versao: '1.0',
              materias: materiasConvertidas
            }).catch(() => undefined);
          } else {
            await api.atribuirEdital(alunoId, {
              concursoId: concursoIdParam,
              nome: `Edital ${dadosConcurso.nome}`,
              materias: materiasConvertidas
            }).catch(() => undefined);
          }
        }
        setSucesso('Concurso atualizado com sucesso!');
      } else {
        // MODO CRIAÇÃO
        if (usuario?.papel === 'mentor') {
          const criado = await api.criarConcursoMentor(dadosConcurso);
          concursoId = criado.id;
        } else {
          const criado = await api.criarConcurso(alunoId, dadosConcurso);
          concursoId = criado.id;
        }
        if (materiasConvertidas.length > 0 && concursoId) {
          if (usuario?.papel === 'mentor') {
            await api.criarEditalMentor({
              concursoId,
              nome: `Edital ${dadosConcurso.nome}`,
              versao: '1.0',
              materias: materiasConvertidas
            }).catch(() => undefined);
          } else {
            await api.atribuirEdital(alunoId, {
              concursoId,
              nome: `Edital ${dadosConcurso.nome}`,
              materias: materiasConvertidas
            }).catch(() => undefined);
          }
        }
        setSucesso('Concurso cadastrado com sucesso!');
      }

      await recarregarConcursos();
      setActiveContestId(concursoId);

      setTimeout(() => {
        if (usuario?.papel === 'mentor') {
          navigate('/mentor/concursos-catalogo');
        } else {
          navigate('/dashboard');
        }
      }, 1000);
    } catch (err: any) {
      setErro(err.message || 'Erro ao salvar concurso.');
    } finally {
      setSalvando(false);
    }
  };

  // Funções para manipulação dinâmica do edital manual
  const adicionarMateria = () => {
    setFormManual((prev) => ({
      ...prev,
      materias: [
        ...prev.materias,
        {
          nome: '',
          peso: 1,
          topicos: [{ nome: '', subtopicos: [''] }]
        }
      ]
    }));
  };

  const removerMateria = (mIdx: number) => {
    setFormManual((prev) => ({
      ...prev,
      materias: prev.materias.filter((_, i) => i !== mIdx)
    }));
  };

  const adicionarTopico = (mIdx: number) => {
    setFormManual((prev) => {
      const novas = [...prev.materias];
      novas[mIdx].topicos.push({ nome: '', subtopicos: [''] });
      return { ...prev, materias: novas };
    });
  };

  const removerTopico = (mIdx: number, tIdx: number) => {
    setFormManual((prev) => {
      const novas = [...prev.materias];
      novas[mIdx].topicos = novas[mIdx].topicos.filter((_, i) => i !== tIdx);
      return { ...prev, materias: novas };
    });
  };

  const adicionarSubtopico = (mIdx: number, tIdx: number) => {
    setFormManual((prev) => {
      const novas = [...prev.materias];
      novas[mIdx].topicos[tIdx].subtopicos.push('');
      return { ...prev, materias: novas };
    });
  };

  const removerSubtopico = (mIdx: number, tIdx: number, sIdx: number) => {
    setFormManual((prev) => {
      const novas = [...prev.materias];
      novas[mIdx].topicos[tIdx].subtopicos = novas[mIdx].topicos[tIdx].subtopicos.filter((_, i) => i !== sIdx);
      return { ...prev, materias: novas };
    });
  };

  if (carregando) {
    return (
      <div className="page-container" style={{ maxWidth: '1000px', margin: '0 auto', paddingTop: '60px', textAlign: 'center' }}>
        <div className="loading-state">Carregando dados do concurso…</div>
      </div>
    );
  }



  const excluirConcurso = async () => {
    if (!concursoIdParam) return;
    if (!confirm(`Tem certeza que deseja excluir/desativar o concurso "${formManual.nome || 'este concurso'}"?`)) {
      return;
    }

    setSalvando(true);
    setErro('');
    try {
      if (usuario?.papel === 'mentor') {
        await api.desativarConcurso('eu', concursoIdParam);
      } else {
        await api.desativarConcurso(alunoId, concursoIdParam);
      }
      await recarregarConcursos();
      if (concursoIdParam === activeContestId) {
        setActiveContestId('');
      }
      setSucesso('Concurso excluído com sucesso!');
      setTimeout(() => {
        if (usuario?.papel === 'mentor') {
          navigate('/mentor/concursos-catalogo');
        } else {
          navigate('/dashboard');
        }
      }, 1000);
    } catch (err: any) {
      setErro(err.message || 'Erro ao excluir concurso.');
    } finally {
      setSalvando(false);
    }
  };

  return (
    <div className="page-container" style={{ maxWidth: '1000px', margin: '0 auto', paddingBottom: '60px' }}>
      <div className="page-header flex-header" style={{ marginBottom: '24px' }}>
        <div>
          <span className="eyebrow">{modoEdicao ? 'Editar Concurso' : 'Novo Concurso'}</span>
          <h1>{modoEdicao ? '✏️ Editar Concurso' : '🎯 Cadastrar Novo Concurso'}</h1>
          <p>{modoEdicao ? 'Atualize os dados do concurso e o conteúdo programático do edital.' : 'Importe a estrutura completa via IA (JSON) ou cadastre as disciplinas manualmente.'}</p>
        </div>
        <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
          {modoEdicao && (
            <button
              type="button"
              className="btn-secondary danger-text"
              onClick={excluirConcurso}
              disabled={salvando}
              style={{ color: 'var(--red)', borderColor: 'var(--red)' }}
            >
              🗑 Excluir Concurso
            </button>
          )}
          <Link className="btn-secondary" to={usuario?.papel === 'mentor' ? '/mentor/concursos-catalogo' : '/dashboard'}>
            ← Voltar
          </Link>
        </div>
      </div>

      {erro && <div className="feedback-banner error" style={{ marginBottom: '20px' }}>{erro}</div>}
      {sucesso && <div className="feedback-banner success" style={{ marginBottom: '20px' }}>{sucesso}</div>}

      {/* TABS DE SELEÇÃO DE MODO */}
      <div className="card-base" style={{ padding: '6px', marginBottom: '24px', display: 'flex', gap: '8px' }}>
        <button
          type="button"
          onClick={() => { setModo('json'); setErro(''); }}
          style={{
            flex: 1,
            padding: '12px 18px',
            borderRadius: '8px',
            border: 'none',
            background: modo === 'json' ? 'var(--accent)' : 'transparent',
            color: modo === 'json' ? '#fff' : 'var(--text2)',
            fontWeight: 700,
            fontSize: '14px',
            cursor: 'pointer',
            transition: 'all 0.2s ease'
          }}
        >
          ⚡ Importar por JSON (Recomendado via IA)
        </button>
        <button
          type="button"
          onClick={() => { setModo('manual'); setErro(''); }}
          style={{
            flex: 1,
            padding: '12px 18px',
            borderRadius: '8px',
            border: 'none',
            background: modo === 'manual' ? 'var(--accent)' : 'transparent',
            color: modo === 'manual' ? '#fff' : 'var(--text2)',
            fontWeight: 700,
            fontSize: '14px',
            cursor: 'pointer',
            transition: 'all 0.2s ease'
          }}
        >
          ✍️ Cadastrar Manualmente
        </button>
      </div>

      {/* ABA 1: IMPORTAR VIA JSON */}
      {modo === 'json' && (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          {/* ACCORDION 1: PROMPT PARA IA */}
          <div className="card-base" style={{ padding: '20px' }}>
            <div
              style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', cursor: 'pointer' }}
              onClick={() => setMostrarPrompt(!mostrarPrompt)}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <span style={{ fontSize: '20px' }}>🤖</span>
                <div>
                  <h3 style={{ margin: 0, fontSize: '15px' }}>Prompt de IA para extrair edital (Copie e Cole no ChatGPT / Gemini)</h3>
                  <small style={{ color: 'var(--text3)' }}>Copie este comando para transformar qualquer PDF ou edital em JSON instantaneamente.</small>
                </div>
              </div>
              <span style={{ fontSize: '14px', color: 'var(--text2)' }}>{mostrarPrompt ? '▲ Recolher' : '▼ Expandir'}</span>
            </div>

            {mostrarPrompt && (
              <div style={{ marginTop: '16px', paddingTop: '16px', borderTop: '1px solid var(--border)' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                  <span style={{ fontSize: '12px', fontWeight: 700, color: 'var(--accent)' }}>COMANDO PRONTO PARA IA:</span>
                  <button
                    type="button"
                    className="btn-secondary"
                    onClick={() => copiarTexto(PROMPT_IA_EXEMPLO, setCopiadoPrompt)}
                    style={{ padding: '6px 12px', fontSize: '11px' }}
                  >
                    {copiadoPrompt ? '✓ Copiado!' : '📋 Copiar Prompt'}
                  </button>
                </div>
                <textarea
                  className="form-control"
                  readOnly
                  rows={10}
                  value={PROMPT_IA_EXEMPLO}
                  style={{ fontFamily: 'monospace', fontSize: '11.5px', background: 'var(--bg-card)', color: 'var(--text1)' }}
                />
              </div>
            )}
          </div>

          {/* ACCORDION 2: EXEMPLO DE JSON */}
          <div className="card-base" style={{ padding: '20px' }}>
            <div
              style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', cursor: 'pointer' }}
              onClick={() => setMostrarExemplo(!mostrarExemplo)}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <span style={{ fontSize: '20px' }}>📋</span>
                <div>
                  <h3 style={{ margin: 0, fontSize: '15px' }}>Modelo & Exemplo de JSON de Concurso</h3>
                  <small style={{ color: 'var(--text3)' }}>Veja a estrutura aceita ou teste preenchendo com o exemplo pronto.</small>
                </div>
              </div>
              <span style={{ fontSize: '14px', color: 'var(--text2)' }}>{mostrarExemplo ? '▲ Recolher' : '▼ Expandir'}</span>
            </div>

            {mostrarExemplo && (
              <div style={{ marginTop: '16px', paddingTop: '16px', borderTop: '1px solid var(--border)' }}>
                <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end', marginBottom: '8px' }}>
                  <button
                    type="button"
                    className="btn-secondary"
                    onClick={() => copiarTexto(JSON.stringify(JSON_EXEMPLO, null, 2), setCopiadoExemplo)}
                    style={{ padding: '6px 12px', fontSize: '11px' }}
                  >
                    {copiadoExemplo ? '✓ Copiado!' : '📋 Copiar Exemplo'}
                  </button>
                  <button
                    type="button"
                    className="btn-primary"
                    onClick={preencherComExemplo}
                    style={{ padding: '6px 12px', fontSize: '11px' }}
                  >
                    ✨ Testar com Exemplo
                  </button>
                </div>
                <textarea
                  className="form-control"
                  readOnly
                  rows={10}
                  value={JSON.stringify(JSON_EXEMPLO, null, 2)}
                  style={{ fontFamily: 'monospace', fontSize: '11.5px', background: 'var(--bg-card)', color: 'var(--text1)' }}
                />
              </div>
            )}
          </div>

          {/* ÁREA DE ENTRADA DO JSON */}
          <div className="card-base" style={{ padding: '24px' }}>
            <h2 style={{ fontSize: '16px', marginBottom: '12px' }}>📥 Cole o JSON do Concurso abaixo:</h2>
            <textarea
              className="form-control"
              rows={14}
              placeholder="Cole o código JSON gerado pela IA aqui…"
              value={jsonText}
              onChange={(e) => { setJsonText(e.target.value); setErro(''); }}
              style={{ fontFamily: 'monospace', fontSize: '12.5px', lineHeight: 1.5 }}
            />

            <div style={{ marginTop: '20px', display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
              <button
                type="button"
                className="btn-secondary"
                onClick={() => setJsonText('')}
                disabled={!jsonText || salvando}
              >
                Limpar Campo
              </button>
              <button
                type="button"
                className="btn-primary"
                disabled={salvando || !jsonText.trim()}
                onClick={salvarViaJson}
                style={{ padding: '12px 24px', fontSize: '14px', fontWeight: 800 }}
              >
                {salvando ? 'Importando & Criando…' : '⚡ Criar Concurso via JSON'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ABA 2: CADASTRO MANUAL */}
      {modo === 'manual' && (
        <form onSubmit={salvarViaManual} className="card-base" style={{ padding: '28px' }}>
          <h2 style={{ fontSize: '18px', marginBottom: '20px' }}>📝 Dados Principais do Concurso</h2>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', marginBottom: '24px' }}>
            <label style={{ gridColumn: 'span 2' }}>
              <span className="field-label">Nome do Concurso *</span>
              <input
                className="form-control"
                placeholder="Ex: Receita Federal — Auditor Fiscal"
                value={formManual.nome}
                onChange={(e) => setFormManual({ ...formManual, nome: e.target.value })}
                required
              />
            </label>

            <label>
              <span className="field-label">Banca Organizadora *</span>
              <input
                className="form-control"
                placeholder="Ex: FGV, Cebraspe, Vunesp"
                value={formManual.banca}
                onChange={(e) => setFormManual({ ...formManual, banca: e.target.value })}
                required
              />
            </label>

            <label>
              <span className="field-label">Cargo</span>
              <input
                className="form-control"
                placeholder="Ex: Auditor Fiscal"
                value={formManual.cargo}
                onChange={(e) => setFormManual({ ...formManual, cargo: e.target.value })}
              />
            </label>

            <label>
              <span className="field-label">Salário Previsto (R$)</span>
              <input
                className="form-control"
                type="number"
                step="0.01"
                placeholder="Ex: 12500"
                value={formManual.salario}
                onChange={(e) => setFormManual({ ...formManual, salario: e.target.value })}
              />
            </label>

            <label>
              <span className="field-label">Prazos de Revisão (dias)</span>
              <input
                className="form-control"
                placeholder="1,7,30"
                value={formManual.prazosRevisao}
                onChange={(e) => setFormManual({ ...formManual, prazosRevisao: e.target.value })}
              />
            </label>

            <label style={{ display: 'flex', alignItems: 'center', gap: '8px', marginTop: '12px' }}>
              <input
                type="checkbox"
                checked={formManual.preEdital}
                onChange={(e) => setFormManual({ ...formManual, preEdital: e.target.checked })}
              />
              <span style={{ fontSize: '13px', fontWeight: 600 }}>Concurso em Pré-Edital</span>
            </label>

            {!formManual.preEdital && (
              <label>
                <span className="field-label">Data da Prova</span>
                <input
                  className="form-control"
                  type="date"
                  value={formManual.dataProva}
                  onChange={(e) => setFormManual({ ...formManual, dataProva: e.target.value })}
                />
              </label>
            )}
          </div>

          <hr style={{ border: 'none', borderTop: '1px solid var(--border)', margin: '24px 0' }} />

          {/* CONSTRUTOR DINÂMICO DE EDITAL */}
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
            <div>
              <h2 style={{ fontSize: '18px', margin: 0 }}>📚 Disciplinas & Conteúdo Programático</h2>
              <p style={{ fontSize: '12px', color: 'var(--text3)', margin: '2px 0 0' }}>Adicione as matérias, tópicos e subtópicos do edital.</p>
            </div>
            <button type="button" className="btn-secondary" onClick={adicionarMateria}>
              ＋ Adicionar Matéria
            </button>
          </div>

          {formManual.materias.map((materia, mIdx) => (
            <div
              key={mIdx}
              style={{
                background: 'var(--bg-card)',
                border: '1px solid var(--border)',
                borderRadius: '10px',
                padding: '20px',
                marginBottom: '20px'
              }}
            >
              <div style={{ display: 'flex', gap: '12px', alignItems: 'center', marginBottom: '16px' }}>
                <input
                  className="form-control"
                  placeholder={`Nome da Matéria ${mIdx + 1} (ex: Língua Portuguesa)`}
                  value={materia.nome}
                  onChange={(e) => {
                    const novas = [...formManual.materias];
                    novas[mIdx].nome = e.target.value;
                    setFormManual({ ...formManual, materias: novas });
                  }}
                  style={{ fontWeight: 700, fontSize: '15px' }}
                />
                <button
                  type="button"
                  className="icon-button danger"
                  onClick={() => removerMateria(mIdx)}
                  title="Remover Matéria"
                >
                  ×
                </button>
              </div>

              {/* TÓPICOS DA MATÉRIA */}
              <div style={{ paddingLeft: '16px', borderLeft: '3px solid var(--accent)' }}>
                {materia.topicos.map((topico, tIdx) => (
                  <div key={tIdx} style={{ marginBottom: '16px', background: 'rgba(255,255,255,0.03)', padding: '12px', borderRadius: '8px' }}>
                    <div style={{ display: 'flex', gap: '10px', alignItems: 'center', marginBottom: '8px' }}>
                      <span style={{ fontSize: '12px', color: 'var(--accent)', fontWeight: 700 }}>📖 Tópico {tIdx + 1}:</span>
                      <input
                        className="form-control"
                        placeholder="Nome do Tópico (ex: Compreensão de Texto)"
                        value={topico.nome}
                        onChange={(e) => {
                          const novas = [...formManual.materias];
                          novas[mIdx].topicos[tIdx].nome = e.target.value;
                          setFormManual({ ...formManual, materias: novas });
                        }}
                      />
                      <button
                        type="button"
                        className="icon-button danger"
                        onClick={() => removerTopico(mIdx, tIdx)}
                        title="Remover Tópico"
                      >
                        ×
                      </button>
                    </div>

                    {/* SUBTÓPICOS DO TÓPICO */}
                    <div style={{ paddingLeft: '20px' }}>
                      {topico.subtopicos.map((sub, sIdx) => (
                        <div key={sIdx} style={{ display: 'flex', gap: '8px', alignItems: 'center', marginBottom: '6px' }}>
                          <span style={{ fontSize: '11px', color: 'var(--text3)' }}>↳ Subtópico:</span>
                          <input
                            className="form-control"
                            placeholder="Nome do Subtópico (opcional)"
                            value={sub}
                            onChange={(e) => {
                              const novas = [...formManual.materias];
                              novas[mIdx].topicos[tIdx].subtopicos[sIdx] = e.target.value;
                              setFormManual({ ...formManual, materias: novas });
                            }}
                            style={{ fontSize: '12px' }}
                          />
                          <button
                            type="button"
                            className="icon-button danger"
                            style={{ padding: '2px 6px', fontSize: '10px' }}
                            onClick={() => removerSubtopico(mIdx, tIdx, sIdx)}
                          >
                            ×
                          </button>
                        </div>
                      ))}
                      <button
                        type="button"
                        className="btn-secondary"
                        onClick={() => adicionarSubtopico(mIdx, tIdx)}
                        style={{ padding: '3px 8px', fontSize: '10px', marginTop: '4px' }}
                      >
                        ＋ Subtópico
                      </button>
                    </div>
                  </div>
                ))}

                <button
                  type="button"
                  className="btn-secondary"
                  onClick={() => adicionarTopico(mIdx)}
                  style={{ padding: '4px 10px', fontSize: '11px' }}
                >
                  ＋ Adicionar Tópico nesta Matéria
                </button>
              </div>
            </div>
          ))}

          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '24px' }}>
            {modoEdicao ? (
              <button
                type="button"
                className="btn-secondary danger-text"
                onClick={excluirConcurso}
                disabled={salvando}
                style={{ color: 'var(--red)', borderColor: 'var(--red)' }}
              >
                🗑 Excluir Concurso
              </button>
            ) : <div />}
            <div style={{ display: 'flex', gap: '12px' }}>
              <Link className="btn-secondary" to={usuario?.papel === 'mentor' ? '/mentor/concursos-catalogo' : '/dashboard'}>
                Cancelar
              </Link>
              <button
                type="submit"
                className="btn-primary"
                disabled={salvando || !formManual.nome.trim() || !formManual.banca.trim()}
                style={{ padding: '12px 24px', fontSize: '14px', fontWeight: 800 }}
              >
                {salvando ? (modoEdicao ? 'Salvando…' : 'Cadastrando…') : (modoEdicao ? '💾 Salvar Alterações' : '💾 Criar Concurso & Edital')}
              </button>
            </div>
          </div>
        </form>
      )}
    </div>
  );
};
