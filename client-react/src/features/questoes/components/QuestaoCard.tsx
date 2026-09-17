import React, { useState, useEffect } from 'react';
import { Card } from '../../../components/ui/Card';
import { Badge } from '../../../components/ui/Badge';
import { Button } from '../../../components/ui/Button';

interface QuestaoCardProps {
  questao: any;
  numeroQuestao: number;
  onResponder: (questaoId: string, resposta: string) => Promise<any>;
  onHistorico: (questaoId: string) => Promise<any[]>;
}

export const QuestaoCard: React.FC<QuestaoCardProps> = ({
  questao,
  numeroQuestao,
  onResponder,
  onHistorico,
}) => {
  const [respostaSelecionada, setRespostaSelecionada] = useState<string>('');
  const [resultado, setResultado] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [historico, setHistorico] = useState<any[] | null>(null);
  const [carregandoHistorico, setCarregandoHistorico] = useState(false);
  const [erroHistorico, setErroHistorico] = useState<string | null>(null);

  // Sincroniza estado caso o histórico ou a questão mude
  useEffect(() => {
    setRespostaSelecionada('');
    setResultado(null);
    setHistorico(null);
    setErroHistorico(null);
  }, [questao.id]);

  const handleSubmeter = async () => {
    if (!respostaSelecionada) return;
    setLoading(true);
    const res = await onResponder(questao.id, respostaSelecionada);
    if (res) {
      setResultado(res);
      setHistorico(null);
    }
    setLoading(false);
  };

  const handleRefazer = () => {
    setResultado(null);
    setRespostaSelecionada('');
  };

  const alternarHistorico = async () => {
    if (historico !== null) {
      setHistorico(null);
      return;
    }
    setCarregandoHistorico(true);
    setErroHistorico(null);
    try {
      setHistorico(await onHistorico(questao.id));
    } catch (e: any) {
      setErroHistorico(e?.message || 'Não foi possível carregar o histórico.');
    } finally {
      setCarregandoHistorico(false);
    }
  };

  // Normalização flexível das alternativas com remoção de prefixos repetidos
  let listaAlternativas: { letra: string; texto: string }[] = [];
  if (Array.isArray(questao.alternativas) && questao.alternativas.length > 0) {
    const letras = ['A', 'B', 'C', 'D', 'E'];
    listaAlternativas = questao.alternativas.map((alt: any, idx: number) => {
      const letraPadrao = letras[idx] || `${idx + 1}`;
      let texto = typeof alt === 'string' ? alt : alt.texto || String(alt);
      const letra = (typeof alt === 'object' && alt?.letra) ? alt.letra : letraPadrao;

      // Remove prefixos como "A) ", "A - ", "A. " se já estiverem presentes no texto
      const match = texto.match(/^[A-Ea-e][\)\.\-:]\s*(.*)/);
      if (match && match[1]) {
        texto = match[1];
      }

      return { letra, texto };
    });
  } else if (questao.tipo === 'certo_errado') {
    listaAlternativas = [
      { letra: 'Certo', texto: 'Certo' },
      { letra: 'Errado', texto: 'Errado' },
    ];
  } else {
    listaAlternativas = [
      { letra: 'A', texto: questao.opcaoA || 'Opção A' },
      { letra: 'B', texto: questao.opcaoB || 'Opção B' },
      { letra: 'C', texto: questao.opcaoC || 'Opção C' },
      { letra: 'D', texto: questao.opcaoD || 'Opção D' },
      { letra: 'E', texto: questao.opcaoE || 'Opção E' },
    ].filter((a) => a.texto);
  }

  const explicacao = questao.explicacao || questao.gabaritoComentado;
  const respostaGabarito = resultado?.respostaCorreta || questao.respostaCorreta || '';

  // Determina se uma alternativa é o gabarito oficial
  const ehGabaritoOficial = (alt: { letra: string; texto: string }, idx: number) => {
    if (!respostaGabarito) return false;
    const gabUpper = respostaGabarito.trim().toUpperCase();
    const letraAlt = alt.letra.trim().toUpperCase();
    const idxLetra = String.fromCharCode(65 + idx);

    if (gabUpper === letraAlt || gabUpper === idxLetra) return true;
    if (gabUpper === alt.texto.trim().toUpperCase()) return true;
    if (gabUpper.startsWith(`${letraAlt})`) || gabUpper.startsWith(`${idxLetra})`)) return true;
    if (questao.tipo === 'certo_errado') {
      if ((gabUpper === 'CERTO' || gabUpper === 'C') && (letraAlt === 'CERTO' || letraAlt === 'C')) return true;
      if ((gabUpper === 'ERRADO' || gabUpper === 'E') && (letraAlt === 'ERRADO' || letraAlt === 'E')) return true;
    }
    return false;
  };

  return (
    <Card variant="default" style={{ marginBottom: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '8px', marginBottom: '16px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>
          <span style={{ fontWeight: 700, fontSize: '0.9rem', color: 'var(--accent-primary)' }}>
            Questão {numeroQuestao}
          </span>
          {questao.disciplina && <Badge variant="info">{questao.disciplina}</Badge>}
          {questao.assunto && <Badge variant="neutral">{questao.assunto}</Badge>}
          {questao.tipo && (
            <Badge variant="warning">
              {questao.tipo === 'certo_errado' ? 'Certo / Errado' : 'Múltipla Escolha'}
            </Badge>
          )}
          {resultado && (
            <Badge variant={resultado.correto ? 'success' : 'danger'}>
              {resultado.correto ? '✓ Acerto' : '✕ Erro'}
            </Badge>
          )}
        </div>
        <button
          type="button"
          onClick={() => void alternarHistorico()}
          disabled={carregandoHistorico}
          style={{ background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', color: 'var(--text-primary)', padding: '6px 10px', borderRadius: 'var(--radius-sm)', cursor: carregandoHistorico ? 'wait' : 'pointer', fontSize: '0.8rem', fontWeight: 600 }}
        >
          {carregandoHistorico ? 'Carregando...' : historico !== null ? 'Ocultar histórico' : 'Histórico'}
        </button>
        {questao.banca && (
          <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>
            {questao.banca} {questao.ano ? `• ${questao.ano}` : ''}
          </span>
        )}
      </div>

      <p style={{ fontSize: '1rem', color: 'var(--text-primary)', marginBottom: '20px', lineHeight: 1.65, whiteSpace: 'pre-line' }}>
        {questao.enunciado || questao.texto}
      </p>

      {erroHistorico && <p style={{ color: 'var(--status-danger)', marginBottom: '16px' }}>{erroHistorico}</p>}
      {historico !== null && (
        <div style={{ marginBottom: '20px', padding: '12px', borderRadius: 'var(--radius-md)', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)' }}>
          <strong style={{ display: 'block', marginBottom: '8px' }}>Tentativas anteriores</strong>
          {historico.length === 0 ? (
            <span style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>Você ainda não respondeu esta questão.</span>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              {historico.map((tentativa, indice) => (
                <div key={tentativa.id || indice} style={{ display: 'flex', justifyContent: 'space-between', gap: '12px', flexWrap: 'wrap', fontSize: '0.9rem' }}>
                  <span>
                    <strong style={{ color: tentativa.correto ? 'var(--status-success)' : 'var(--status-danger)' }}>
                      {tentativa.correto ? 'Acerto' : 'Erro'}
                    </strong>
                    {` · resposta: ${tentativa.respostaAluno}`}
                  </span>
                  <span style={{ color: 'var(--text-secondary)' }}>
                    {tentativa.respondidoEm ? new Date(tentativa.respondidoEm).toLocaleString('pt-BR') : ''}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', gap: '10px', marginBottom: '20px' }}>
        {listaAlternativas.map((alt, idx) => {
          const isSelected = respostaSelecionada === alt.letra;
          const isGabarito = resultado ? ehGabaritoOficial(alt, idx) : false;
          const isMarcadaErrada = resultado && isSelected && !resultado.correto;

          // Definição de cores baseada no resultado
          let bg = 'var(--bg-secondary)';
          let border = '1px solid var(--border-color)';
          let badgeBg = 'rgba(255, 255, 255, 0.08)';
          let badgeColor = '#fff';

          if (resultado) {
            if (isGabarito) {
              bg = 'rgba(46, 204, 113, 0.15)';
              border = '1px solid #2ecc71';
              badgeBg = '#2ecc71';
            } else if (isMarcadaErrada) {
              bg = 'rgba(231, 76, 60, 0.15)';
              border = '1px solid #e74c3c';
              badgeBg = '#e74c3c';
            } else {
              bg = 'var(--bg-secondary)';
              border = '1px solid rgba(255, 255, 255, 0.05)';
            }
          } else if (isSelected) {
            bg = 'var(--accent-light)';
            border = '1px solid var(--accent-primary)';
            badgeBg = 'var(--accent-primary)';
          }

          const isCertoErrado = questao.tipo === 'certo_errado' || (questao.alternativas?.length === 2 && ['certo', 'errado', 'c', 'e'].includes(alt.letra?.toLowerCase()));
          const letraBadge =
            alt.letra?.toLowerCase() === 'certo' || (isCertoErrado && alt.letra?.toLowerCase() === 'c') ? 'C' :
            alt.letra?.toLowerCase() === 'errado' || (isCertoErrado && alt.letra?.toLowerCase() === 'e') ? 'E' :
            alt.letra;

          const textoExibicao =
            isCertoErrado && (alt.texto?.toLowerCase() === 'c' || alt.texto?.toLowerCase() === 'certo') ? 'Certo' :
            isCertoErrado && (alt.texto?.toLowerCase() === 'e' || alt.texto?.toLowerCase() === 'errado') ? 'Errado' :
            alt.texto;

          return (
            <button
              key={alt.letra}
              onClick={() => !resultado && setRespostaSelecionada(alt.letra)}
              disabled={!!resultado}
              style={{
                textAlign: 'left',
                padding: '12px 16px',
                borderRadius: 'var(--radius-md)',
                background: bg,
                border,
                color: 'var(--text-primary)',
                cursor: resultado ? 'default' : 'pointer',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                gap: '12px',
                transition: 'all var(--transition-fast)',
                opacity: resultado && !isGabarito && !isMarcadaErrada ? 0.7 : 1,
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
                <span
                  style={{
                    width: '28px',
                    height: '28px',
                    borderRadius: '50%',
                    background: badgeBg,
                    color: badgeColor,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontWeight: 600,
                    fontSize: '0.85rem',
                    flexShrink: 0,
                  }}
                >
                  {letraBadge}
                </span>
                <span style={{ fontSize: '0.95rem', lineHeight: 1.5 }}>{textoExibicao}</span>
              </div>
              {resultado && isGabarito && (
                <span style={{ color: '#2ecc71', fontWeight: 700, fontSize: '0.85rem' }}>✓ Correta</span>
              )}
              {resultado && isMarcadaErrada && (
                <span style={{ color: '#e74c3c', fontWeight: 700, fontSize: '0.85rem' }}>✕ Sua resposta</span>
              )}
            </button>
          );
        })}
      </div>

      {!resultado ? (
        <Button
          variant="primary"
          onClick={handleSubmeter}
          isLoading={loading}
          disabled={!respostaSelecionada}
        >
          Submeter Resposta
        </Button>
      ) : (
        <div
          style={{
            padding: '16px',
            borderRadius: 'var(--radius-md)',
            background: resultado.correto ? 'rgba(46, 204, 113, 0.1)' : 'rgba(231, 76, 60, 0.1)',
            border: `1px solid ${resultado.correto ? '#2ecc71' : '#e74c3c'}`,
            color: 'var(--text-primary)',
          }}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '8px', marginBottom: '8px' }}>
            <div style={{ fontWeight: 700, fontSize: '1rem', color: resultado.correto ? '#2ecc71' : '#e74c3c' }}>
              {resultado.correto
                ? '✅ Resposta Correta!'
                : `❌ Resposta Incorreta (Gabarito: ${respostaGabarito})`}
            </div>
            <button
              type="button"
              onClick={handleRefazer}
              style={{
                background: 'var(--bg-secondary)',
                border: '1px solid var(--border-color)',
                color: 'var(--text-primary)',
                padding: '6px 12px',
                borderRadius: 'var(--radius-sm)',
                cursor: 'pointer',
                fontSize: '0.85rem',
                fontWeight: 600,
              }}
            >
              🔄 Refazer Questão
            </button>
          </div>

          {explicacao && (
            <div style={{ marginTop: '10px', fontSize: '0.9rem', color: 'var(--text-secondary)', borderTop: '1px solid var(--border-color)', paddingTop: '10px', lineHeight: 1.5 }}>
              <strong style={{ color: 'var(--text-primary)' }}>💡 Comentário / Justificativa:</strong>
              <p style={{ marginTop: '4px' }}>{explicacao}</p>
            </div>
          )}
        </div>
      )}
    </Card>
  );
};
