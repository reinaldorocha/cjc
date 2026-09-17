import React, { useState } from 'react';
import { useFlashcards } from '../hooks/useFlashcards';
import { BaralhosGrid } from '../components/BaralhosGrid';
import { FlashcardViewer } from '../components/FlashcardViewer';
import { Spinner } from '../../../components/ui/Spinner';
import { EmptyState } from '../../../components/ui/EmptyState';
import { Card } from '../../../components/ui/Card';

export const FlashcardsPage: React.FC = () => {
  const { baralhos, cartoesPendentes, loading, revisarCartao } = useFlashcards();
  const [baralhoAtivo, setBaralhoAtivo] = useState<any>(null);
  const [cartaoIndex, setCartaoIndex] = useState(0);

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '80px' }}>
        <Spinner size="lg" />
      </div>
    );
  }

  const totalBaralhos = baralhos.length;
  const totalCartoes = baralhos.reduce((acc, b) => acc + (b.totalCartoes ?? b.cartoes?.length ?? 0), 0);
  const totalPendentes = cartoesPendentes.length;

  const cartoesBaralho = baralhoAtivo
    ? cartoesPendentes.filter((c) => c.baralhoId === baralhoAtivo.id)
    : cartoesPendentes;

  const cartaoAtual = cartoesBaralho[cartaoIndex];

  return (
    <div className="ui-container" style={{ padding: '24px 0' }}>
      <div className="ui-page-header">
        <h1 className="ui-page-title">Flashcards & Repetição Espaçada</h1>
        <p className="ui-page-subtitle">Memorize conceitos chaves e leis através de cartões inteligentes.</p>
      </div>

      {/* PAINEL DE ESTATÍSTICAS DOS FLASHCARDS */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '16px', marginBottom: '24px' }}>
        <Card variant="glass" style={{ padding: '16px 20px' }}>
          <div style={{ fontSize: '12px', fontWeight: 700, color: 'var(--text2)', marginBottom: '4px' }}>📚 Total de Baralhos</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: 'var(--text1)' }}>{totalBaralhos}</div>
        </Card>
        <Card variant="glass" style={{ padding: '16px 20px' }}>
          <div style={{ fontSize: '12px', fontWeight: 700, color: 'var(--text2)', marginBottom: '4px' }}>🎴 Total de Cartões no Catálogo</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: 'var(--accent)' }}>{totalCartoes}</div>
        </Card>
        <Card variant="glass" style={{ padding: '16px 20px' }}>
          <div style={{ fontSize: '12px', fontWeight: 700, color: 'var(--text2)', marginBottom: '4px' }}>⏰ Cartões Pendentes de Revisão</div>
          <div style={{ fontSize: '24px', fontWeight: 800, color: 'var(--gold)' }}>{totalPendentes}</div>
        </Card>
      </div>

      {baralhoAtivo ? (
        cartaoAtual ? (
          <FlashcardViewer
            cartao={cartaoAtual}
            onRevisar={async (qualidade) => {
              await revisarCartao(cartaoAtual.id, qualidade);
              setCartaoIndex((prev) => prev + 1);
            }}
            onVoltar={() => {
              setBaralhoAtivo(null);
              setCartaoIndex(0);
            }}
          />
        ) : (
          <EmptyState
            icon="🎉"
            title="Parabéns! Você revisou todos os cartões deste baralho!"
            actionLabel="Voltar aos Baralhos"
            onAction={() => {
              setBaralhoAtivo(null);
              setCartaoIndex(0);
            }}
          />
        )
      ) : (
        <BaralhosGrid baralhos={baralhos} onSelecionarBaralho={(b) => setBaralhoAtivo(b)} />
      )}
    </div>
  );
};
