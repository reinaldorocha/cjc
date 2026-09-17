import React, { useState } from 'react';
import { Card } from '../../../components/ui/Card';
import { Button } from '../../../components/ui/Button';

interface FlashcardViewerProps {
  cartao: any;
  onRevisar: (qualidade: number) => void;
  onVoltar: () => void;
}

export const FlashcardViewer: React.FC<FlashcardViewerProps> = ({ cartao, onRevisar, onVoltar }) => {
  const [revelado, setRevelado] = useState(false);

  return (
    <div style={{ maxWidth: '600px', margin: '0 auto' }}>
      <Button variant="ghost" onClick={onVoltar} style={{ marginBottom: '16px' }}>
        ← Voltar aos Baralhos
      </Button>

      <Card
        variant="glass"
        style={{
          minHeight: '260px',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          textAlign: 'center',
          cursor: 'pointer',
          padding: '32px',
        }}
        onClick={() => setRevelado(!revelado)}
      >
        <div style={{ fontSize: '0.8rem', color: 'var(--accent-primary)', textTransform: 'uppercase', marginBottom: '12px', fontWeight: 600 }}>
          {revelado ? 'Verso (Resposta)' : 'Frente (Pergunta) • Clique para virar'}
        </div>
        <p style={{ fontSize: '1.25rem', fontWeight: 500, lineHeight: 1.6 }}>
          {revelado ? (cartao.resposta || cartao.verso) : (cartao.pergunta || cartao.frente)}
        </p>
      </Card>

      {revelado && (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '12px', marginTop: '20px' }}>
          <Button variant="danger" onClick={() => onRevisar(1)}>Errei (1d)</Button>
          <Button variant="secondary" onClick={() => onRevisar(2)}>Difícil (2d)</Button>
          <Button variant="outline" onClick={() => onRevisar(3)}>Bom (4d)</Button>
          <Button variant="primary" onClick={() => onRevisar(4)}>Fácil (7d)</Button>
        </div>
      )}
    </div>
  );
};
