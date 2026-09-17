/* oxlint-disable react/only-export-components */
import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { api } from '../services/api';

interface DataContextType {
  alunoId: string; activeContestId: string; setActiveContestId: (id: string) => void; activeContest: any;
  getArray: (key: string) => any[]; loading: boolean; recarregarConcursos: () => Promise<void>;
  editalAtivo: any | null; recarregarEdital: () => Promise<void>;
}
const DataContext = createContext<DataContextType | undefined>(undefined);
const normalizarConcurso=(c:any)=>({...c,grupoMeusConcursos:c.grupo,ordemMeusConcursos:c.ordem,realizado:c.grupo==='realizado',logoBase64:c.logotipo||'',nomeado:{ativo:!!c.nomeado,data:c.dataNomeacao||''}});

export const DataProvider:React.FC<{children:React.ReactNode;alunoId?:string;concursoInicialId?:string}>=({children,alunoId='eu',concursoInicialId=''})=>{
  const[activeContestId,setActiveContestId]=useState(concursoInicialId),[concursos,setConcursos]=useState<any[]>([]),[editalAtivo,setEditalAtivo]=useState<any|null>(null),[loading,setLoading]=useState(true);
  const recarregarConcursos=useCallback(async()=>{const d: any=await api.listarConcursos(alunoId);setConcursos((Array.isArray(d) ? d : d?.concursos||[]).map(normalizarConcurso))},[alunoId]);
  const recarregarEdital=useCallback(async()=>{if(!activeContestId){setEditalAtivo(null);return}const d: any=await api.listarEditais(alunoId,activeContestId);setEditalAtivo(d?.editais?.[0]||null)},[alunoId,activeContestId]);
  useEffect(()=>{setLoading(true);recarregarConcursos().finally(()=>setLoading(false))},[recarregarConcursos]);
  useEffect(()=>{recarregarEdital().catch(()=>setEditalAtivo(null))},[recarregarEdital]);
  useEffect(()=>{if(concursos.length&&(!activeContestId||!concursos.some(c=>c.id===activeContestId)))setActiveContestId(concursos[0].id);if(!concursos.length&&activeContestId)setActiveContestId('')},[activeContestId,concursos]);
  const getArray=useCallback((key:string):any[]=>{if(key==='concursos')return concursos;if(key==='materias')return(editalAtivo?.materias||[]).map((m:any)=>normalizarItemEdital(m,'materia',activeContestId));if(key==='topicos')return(editalAtivo?.materias||[]).flatMap((m:any)=>(m.topicos||[]).map((t:any)=>normalizarItemEdital({...t,materiaId:m.id},'topico',activeContestId)));if(key==='subtopicos')return(editalAtivo?.materias||[]).flatMap((m:any)=>(m.topicos||[]).flatMap((t:any)=>(t.subtopicos||[]).map((s:any)=>normalizarItemEdital({...s,materiaId:m.id,topicoId:t.id},'subtopico',activeContestId))));return[]},[concursos,editalAtivo,activeContestId]);
  const activeContest=concursos.find(c=>c.id===activeContestId);
  return <DataContext.Provider value={{alunoId,activeContestId,setActiveContestId,activeContest,getArray,loading,recarregarConcursos,editalAtivo,recarregarEdital}}>{children}</DataContext.Provider>;
};
function normalizarItemEdital(item:any,tipo:string,concursoId:string){const materiais=(Array.isArray(item.materiais)?item.materiais:[]).map((m:any)=>({...m,concluido:!!item.progressoMateriais?.[m.id]?.concluido,concluidoEm:item.progressoMateriais?.[m.id]?.concluidoEm}));return{...item,tipo,concursoId,cadernos:materiais,estudado:!!item.estudado,concluido:!!item.estudado,estudadoEm:item.concluidoEm||null,subtopId:item.subtopicoId}}
export const useData=()=>{const c=useContext(DataContext);if(!c)throw new Error('useData must be used within DataProvider');return c};
