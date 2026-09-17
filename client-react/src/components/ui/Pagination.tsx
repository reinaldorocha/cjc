import React from 'react';

interface PaginationProps {
  pagina: number;
  totalPaginas: number;
  total?: number;
  onChange: (pagina: number) => void;
}

export const Pagination: React.FC<PaginationProps> = ({ pagina, totalPaginas, total, onChange }) => {
  if (totalPaginas <= 1) return null;
  return (
    <nav aria-label="Paginação" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '12px', flexWrap: 'wrap', marginTop: '24px' }}>
      <button type="button" disabled={pagina <= 1} onClick={() => onChange(pagina - 1)} style={estiloBotao(pagina <= 1)}>Anterior</button>
      <span style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>
        Página <strong>{pagina}</strong> de <strong>{totalPaginas}</strong>{typeof total === 'number' ? ` · ${total} itens` : ''}
      </span>
      <button type="button" disabled={pagina >= totalPaginas} onClick={() => onChange(pagina + 1)} style={estiloBotao(pagina >= totalPaginas)}>Próxima</button>
    </nav>
  );
};

function estiloBotao(desabilitado: boolean): React.CSSProperties {
  return { padding: '8px 12px', borderRadius: 'var(--radius-sm)', border: '1px solid var(--border-color)', background: 'var(--bg-secondary)', color: 'var(--text-primary)', cursor: desabilitado ? 'not-allowed' : 'pointer', opacity: desabilitado ? 0.5 : 1, fontWeight: 600 };
}
