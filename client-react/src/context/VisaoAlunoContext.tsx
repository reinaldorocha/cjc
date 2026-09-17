import React, { createContext, useContext } from 'react';

const VisaoAlunoContext = createContext(false);

export const VisaoAlunoProvider: React.FC<{ ativa: boolean; children: React.ReactNode }> = ({ ativa, children }) => (
  <VisaoAlunoContext.Provider value={ativa}>{children}</VisaoAlunoContext.Provider>
);

export const useVisaoAluno = () => useContext(VisaoAlunoContext);
