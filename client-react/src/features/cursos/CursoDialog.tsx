import { useEffect, useRef, type ReactNode } from 'react';

export function Janela({ titulo, fechar, children }: { titulo: string; fechar: () => void; children: ReactNode }) {
  const dialog = useRef<HTMLDialogElement>(null);
  useEffect(() => { dialog.current?.showModal(); }, []);
  return <dialog ref={dialog} className="cursos-dialog" aria-label={titulo} onCancel={(e) => { e.preventDefault(); fechar(); }}>
    <header><h2>{titulo}</h2><button type="button" aria-label="Fechar" onClick={fechar}>✕</button></header>
    {children}
  </dialog>;
}
