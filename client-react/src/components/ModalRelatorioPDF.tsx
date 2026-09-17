import React, { useState } from 'react';

interface ModalRelatorioPDFProps {
  alunoNome?: string;
  concursoNome?: string;
  sessoes: any[];
  questoes: any[];
  materias: any[];
  topicos: any[];
  subtopicos: any[];
  aoFechar: () => void;
}

export const ModalRelatorioPDF: React.FC<ModalRelatorioPDFProps> = ({
  alunoNome = 'Aluno',
  concursoNome = 'Concurso',
  sessoes = [],
  questoes = [],
  materias = [],
  topicos = [],
  subtopicos = [],
  aoFechar
}) => {
  const [periodo, setPeriodo] = useState<'7d' | '14d' | '30d' | 'todos'>('7d');
  const [parecerMentor, setParecerMentor] = useState('');

  const dateKey = (d: Date) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;

  const dataInicio = new Date();
  if (periodo === '7d') dataInicio.setDate(dataInicio.getDate() - 6);
  else if (periodo === '14d') dataInicio.setDate(dataInicio.getDate() - 13);
  else if (periodo === '30d') dataInicio.setDate(dataInicio.getDate() - 29);
  else dataInicio.setFullYear(2020, 0, 1);

  const dataInicioStr = dateKey(dataInicio);

  const logDate = (item: any) => item.data || String(item.estudadoEm || item.registradoEm || item.criadoEm || '').slice(0, 10);

  const sessoesFiltradas = sessoes.filter((x) => periodo === 'todos' || logDate(x) >= dataInicioStr);
  const questoesFiltradas = questoes.filter((x) => periodo === 'todos' || logDate(x) >= dataInicioStr);

  const totalSegundos = sessoesFiltradas.reduce((acc, x) => acc + (Number(x.segundos) || 0), 0);
  const totalHorasStr = `${Math.floor(totalSegundos / 3600)}h ${Math.floor((totalSegundos % 3600) / 60)}min`;

  const totalResolvidas = questoesFiltradas.reduce((acc, x) => acc + (Number(x.resolvidas) || 0), 0);
  const totalAcertos = questoesFiltradas.reduce((acc, x) => acc + (Number(x.acertos) || 0), 0);
  const taxaAcerto = totalResolvidas > 0 ? Math.round((totalAcertos / totalResolvidas) * 100) : 0;

  // Cobertura do edital
  const coberturaItens = [...topicos.filter((t) => !subtopicos.some((s) => s.topicoId === t.id)), ...subtopicos];
  const itensConcluidos = coberturaItens.filter((x) => x.estudado).length;
  const pctCobertura = coberturaItens.length > 0 ? Math.round((itensConcluidos / coberturaItens.length) * 100) : 0;

  // Estatísticas por matéria
  const materiasStats = materias.map((m) => {
    const sMat = sessoesFiltradas.filter((x) => x.materiaId === m.id);
    const qMat = questoesFiltradas.filter((x) => x.materiaId === m.id);
    const seg = sMat.reduce((acc, x) => acc + (Number(x.segundos) || 0), 0);
    const res = qMat.reduce((acc, x) => acc + (Number(x.resolvidas) || 0), 0);
    const ace = qMat.reduce((acc, x) => acc + (Number(x.acertos) || 0), 0);
    const pct = res > 0 ? Math.round((ace / res) * 100) : null;

    const tMat = topicos.filter((t) => t.materiaId === m.id);
    const sMatSub = subtopicos.filter((s) => s.materiaId === m.id);
    const cMatItens = [...tMat.filter((t) => !sMatSub.some((s) => s.topicoId === t.id)), ...sMatSub];
    const cMatDone = cMatItens.filter((x) => x.estudado).length;

    return {
      id: m.id,
      nome: m.nome,
      segundos: seg,
      tempoStr: `${Math.floor(seg / 3600)}h ${Math.floor((seg % 3600) / 60)}m`,
      resolvidas: res,
      acertos: ace,
      pct,
      done: cMatDone,
      total: cMatItens.length
    };
  });

  const imprimirPDF = () => {
    window.print();
  };

  const copiarTextoWhatsApp = () => {
    const periodoLabel = periodo === '7d' ? 'Semanal (7 dias)' : periodo === '14d' ? '14 dias' : periodo === '30d' ? 'Mensal (30 dias)' : 'Completo';
    let txt = `📊 *RELATÓRIO DE DESEMPENHO — CHEGA JUNTO CONCURSEIRO*\n`;
    txt += `👤 *Aluno:* ${alunoNome}\n`;
    txt += `🏆 *Concurso:* ${concursoNome}\n`;
    txt += `📅 *Período:* ${periodoLabel}\n`;
    txt += `-----------------------------------\n`;
    txt += `⏱️ *Tempo Estudado:* ${totalHorasStr}\n`;
    txt += `📝 *Questões Resolvidas:* ${totalResolvidas} (${taxaAcerto}% de acertos)\n`;
    txt += `📚 *Edital Coberto:* ${pctCobertura}% (${itensConcluidos}/${coberturaItens.length} tópicos)\n`;
    txt += `-----------------------------------\n`;
    txt += `📌 *DESEMPENHO POR MATÉRIA:*\n`;
    materiasStats.forEach((m) => {
      if (m.segundos > 0 || m.resolvidas > 0) {
        txt += `• *${m.nome}:* ${m.tempoStr} | ${m.resolvidas}q (${m.pct ?? 0}% acertos)\n`;
      }
    });
    if (parecerMentor.trim()) {
      txt += `\n💬 *PARECER DA MENTORIA:*\n${parecerMentor.trim()}\n`;
    }
    navigator.clipboard.writeText(txt);
    alert('Relatório formatado copiado para a área de transferência!');
  };

  return (
    <div className="modal-backdrop" style={{ zIndex: 9999, background: 'rgba(0,0,0,0.85)', overflowY: 'auto' }} onClick={aoFechar}>
      <div
        className="modal-card print-report-container"
        style={{
          maxWidth: '860px',
          width: '95%',
          background: '#ffffff',
          color: '#1a1d24',
          borderRadius: '16px',
          padding: '28px',
          margin: '30px auto',
          boxShadow: '0 20px 50px rgba(0,0,0,0.5)'
        }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Barra de Controles (Não é impressa no PDF) */}
        <div className="no-print" style={{ marginBottom: '24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '12px', borderBottom: '1px solid #e0e4ec', paddingBottom: '16px' }}>
          <div>
            <h3 style={{ margin: 0, fontSize: '18px', color: '#1a1d24', fontWeight: 800 }}>📄 Relatório de Desempenho em PDF</h3>
            <p style={{ margin: '4px 0 0', fontSize: '12px', color: '#667085' }}>Selecione o período e baixe/imprima o relatório para o aluno</p>
          </div>
          <div style={{ display: 'flex', gap: '8px', alignItems: 'center', flexWrap: 'wrap' }}>
            <select className="form-control" style={{ width: 'auto', background: '#f2f4f7', color: '#101828', border: '1px solid #d0d5dd' }} value={periodo} onChange={(e) => setPeriodo(e.target.value as any)}>
              <option value="7d">Relatório Semanal (7 dias)</option>
              <option value="14d">Relatório 14 dias</option>
              <option value="30d">Relatório Mensal (30 dias)</option>
              <option value="todos">Relatório Completo</option>
            </select>
            <button type="button" className="btn-secondary" style={{ padding: '8px 14px', fontSize: '12px', background: '#2563eb', color: '#fff', border: 0 }} onClick={imprimirPDF}>
              🖨️ Baixar / Imprimir PDF
            </button>
            <button type="button" className="btn-secondary" style={{ padding: '8px 14px', fontSize: '12px' }} onClick={copiarTextoWhatsApp}>
              📱 Copiar p/ WhatsApp
            </button>
            <button type="button" className="icon-button" style={{ background: '#f2f4f7', color: '#667085' }} onClick={aoFechar}>
              ×
            </button>
          </div>
        </div>

        {/* Parecer do Mentor (Campo editável pré-impressão, não impresso se vazio) */}
        <div className="no-print" style={{ marginBottom: '20px' }}>
          <label style={{ display: 'block', fontSize: '12px', fontWeight: 700, color: '#344054', marginBottom: '6px' }}>
            💬 Adicionar Parecer / Comentários da Mentoria (Opcional):
          </label>
          <textarea
            className="form-control"
            rows={2}
            style={{ background: '#f8fafc', color: '#0f172a', border: '1px solid #cbd5e1' }}
            placeholder="Digite aqui observações e orientações para o aluno neste período…"
            value={parecerMentor}
            onChange={(e) => setParecerMentor(e.target.value)}
          />
        </div>

        {/* ================= CORPO DO RELATÓRIO IMPRESSO ================= */}
        <div id="printable-report-area">
          {/* Cabeçalho */}
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', borderBottom: '3px solid #2563eb', paddingBottom: '16px', marginBottom: '24px' }}>
            <div>
              <span style={{ fontSize: '11px', fontWeight: 800, color: '#2563eb', letterSpacing: '1px', textTransform: 'uppercase' }}>Chega Junto Concurseiro</span>
              <h1 style={{ margin: '4px 0 0', fontSize: '24px', color: '#0f172a', fontWeight: 900 }}>Relatório de Desempenho</h1>
              <p style={{ margin: '4px 0 0', fontSize: '13px', color: '#475569' }}>
                Concurso: <strong>{concursoNome}</strong>
              </p>
            </div>
            <div style={{ textAlign: 'right' }}>
              <div style={{ fontSize: '14px', fontWeight: 800, color: '#0f172a' }}>{alunoNome}</div>
              <div style={{ fontSize: '11px', color: '#64748b', marginTop: '2px' }}>
                Período: {periodo === '7d' ? 'Semanal (7 dias)' : periodo === '14d' ? '14 dias' : periodo === '30d' ? 'Mensal (30 dias)' : 'Completo'}
              </div>
              <div style={{ fontSize: '10px', color: '#94a3b8', marginTop: '2px' }}>Emitido em {new Date().toLocaleDateString('pt-BR')}</div>
            </div>
          </div>

          {/* Cards de Métricas Principais */}
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '12px', marginBottom: '28px' }}>
            <div style={{ background: '#f8fafc', border: '1px solid #e2e8f0', borderRadius: '10px', padding: '14px', textAlign: 'center' }}>
              <span style={{ fontSize: '11px', color: '#64748b', display: 'block', fontWeight: 600 }}>TEMPO ESTUDADO</span>
              <strong style={{ fontSize: '20px', color: '#0f172a', display: 'block', marginTop: '4px' }}>{totalHorasStr}</strong>
            </div>
            <div style={{ background: '#f8fafc', border: '1px solid #e2e8f0', borderRadius: '10px', padding: '14px', textAlign: 'center' }}>
              <span style={{ fontSize: '11px', color: '#64748b', display: 'block', fontWeight: 600 }}>QUESTÕES RESOLVIDAS</span>
              <strong style={{ fontSize: '20px', color: '#0f172a', display: 'block', marginTop: '4px' }}>{totalResolvidas}</strong>
            </div>
            <div style={{ background: '#f8fafc', border: '1px solid #e2e8f0', borderRadius: '10px', padding: '14px', textAlign: 'center' }}>
              <span style={{ fontSize: '11px', color: '#64748b', display: 'block', fontWeight: 600 }}>TAXA DE ACERTO</span>
              <strong style={{ fontSize: '20px', color: taxaAcerto >= 70 ? '#16a34a' : '#d97706', display: 'block', marginTop: '4px' }}>{taxaAcerto}%</strong>
            </div>
            <div style={{ background: '#f8fafc', border: '1px solid #e2e8f0', borderRadius: '10px', padding: '14px', textAlign: 'center' }}>
              <span style={{ fontSize: '11px', color: '#64748b', display: 'block', fontWeight: 600 }}>COBERTURA EDITAL</span>
              <strong style={{ fontSize: '20px', color: '#2563eb', display: 'block', marginTop: '4px' }}>{pctCobertura}%</strong>
            </div>
          </div>

          {/* Tabela de Desempenho por Matéria */}
          <div style={{ marginBottom: '28px' }}>
            <h3 style={{ fontSize: '15px', color: '#0f172a', fontWeight: 800, marginBottom: '12px' }}>📊 Detalhamento por Matéria</h3>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '12px' }}>
              <thead>
                <tr style={{ background: '#f1f5f9', color: '#334155', textAlign: 'left', borderBottom: '2px solid #cbd5e1' }}>
                  <th style={{ padding: '10px 12px' }}>Matéria</th>
                  <th style={{ padding: '10px 12px' }}>Horas Estudadas</th>
                  <th style={{ padding: '10px 12px' }}>Questões</th>
                  <th style={{ padding: '10px 12px' }}>Acertos (%)</th>
                  <th style={{ padding: '10px 12px' }}>Progresso Edital</th>
                </tr>
              </thead>
              <tbody>
                {materiasStats.map((m) => (
                  <tr key={m.id} style={{ borderBottom: '1px solid #e2e8f0' }}>
                    <td style={{ padding: '10px 12px', fontWeight: 700, color: '#0f172a' }}>{m.nome}</td>
                    <td style={{ padding: '10px 12px', color: '#334155' }}>{m.tempoStr}</td>
                    <td style={{ padding: '10px 12px', color: '#334155' }}>{m.resolvidas}</td>
                    <td style={{ padding: '10px 12px', fontWeight: 700, color: m.pct == null ? '#94a3b8' : m.pct >= 70 ? '#16a34a' : '#d97706' }}>
                      {m.pct == null ? '—' : `${m.pct}%`}
                    </td>
                    <td style={{ padding: '10px 12px', color: '#334155' }}>
                      {m.total > 0 ? `${m.done}/${m.total} (${Math.round((m.done / m.total) * 100)}%)` : '—'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Parecer do Mentor Impresso */}
          {parecerMentor.trim() && (
            <div style={{ background: '#f8fafc', borderLeft: '4px solid #2563eb', padding: '16px', borderRadius: '4px', marginBottom: '28px' }}>
              <h4 style={{ margin: '0 0 6px', fontSize: '13px', color: '#1e293b', fontWeight: 800 }}>💬 Parecer & Orientações da Mentoria:</h4>
              <p style={{ margin: 0, fontSize: '12px', color: '#334155', whiteSpace: 'pre-wrap', lineHeight: 1.5 }}>{parecerMentor.trim()}</p>
            </div>
          )}

          {/* Rodapé */}
          <div style={{ borderTop: '1px solid #e2e8f0', paddingTop: '12px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', fontSize: '10px', color: '#94a3b8' }}>
            <span>Chega Junto Concurseiro — Plataforma de Mentoria e Acompanhamento de Concursos</span>
            <span>Página 1 de 1</span>
          </div>
        </div>
      </div>
    </div>
  );
};
