/* oxlint-disable react/only-export-components */
import React,{createContext,useContext,useMemo,useState}from'react';
interface Contexto{alunoSelecionadoId:string;selecionarAluno:(id:string)=>void}
const MentoriaContext=createContext<Contexto|undefined>(undefined);
export const MentoriaProvider:React.FC<{children:React.ReactNode}>=({children})=>{const[alunoSelecionadoId,selecionarAluno]=useState('');const valor=useMemo(()=>({alunoSelecionadoId,selecionarAluno}),[alunoSelecionadoId]);return <MentoriaContext.Provider value={valor}>{children}</MentoriaContext.Provider>};
export const useMentoria=()=>{const c=useContext(MentoriaContext);if(!c)throw new Error('MentoriaProvider ausente');return c};
