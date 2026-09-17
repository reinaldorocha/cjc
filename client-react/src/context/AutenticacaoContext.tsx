/* oxlint-disable react/only-export-components */
import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';
import { api, ErroApi } from '../services/api';

export type Papel = 'mestre' | 'mentor' | 'aluno';
export interface UsuarioAutenticado { id:string; nome:string; email:string; papel:Papel; ativo:boolean; expirado:boolean; permiteCronogramaInteligente:boolean; dataExpiracaoPlano?:string }
interface Contexto { usuario:UsuarioAutenticado|null; carregando:boolean; erroAcesso:string; entrar:(email:string,senha:string)=>Promise<void>; sair:()=>Promise<void>; atualizarNome:(nome:string)=>Promise<void>; alterarSenha:(atual:string,nova:string)=>Promise<void> }
const AutenticacaoContext=createContext<Contexto|undefined>(undefined);

export const AutenticacaoProvider:React.FC<{children:React.ReactNode}>=({children})=>{
  const[usuario,setUsuario]=useState<UsuarioAutenticado|null>(null);const[carregando,setCarregando]=useState(true);const[erroAcesso,setErroAcesso]=useState('');
  useEffect(()=>{api.obterSessao().then(d=>setUsuario(d.usuario)).catch((e:ErroApi)=>{setUsuario(null);if(e.codigo==='PLANO_EXPIRADO')setErroAcesso(e.message)}).finally(()=>setCarregando(false))},[]);
  const entrar=async(email:string,senha:string)=>{setErroAcesso('');const d=await api.entrar(email,senha);setUsuario(d.usuario)};
  const sair=async()=>{await api.sair().catch(()=>undefined);setUsuario(null)};
  const atualizarNome=async(nome:string)=>{await api.alterarNome(nome);setUsuario(u=>u?{...u,nome}:u)};
  const alterarSenha=async(atual:string,nova:string)=>{await api.alterarSenha(atual,nova)};
  const valor=useMemo(()=>({usuario,carregando,erroAcesso,entrar,sair,atualizarNome,alterarSenha}),[usuario,carregando,erroAcesso]);
  return <AutenticacaoContext.Provider value={valor}>{children}</AutenticacaoContext.Provider>;
};
export const useAutenticacao=()=>{const c=useContext(AutenticacaoContext);if(!c)throw new Error('AutenticacaoProvider ausente');return c};
