import React, { createContext, useContext, useEffect, useState } from 'react';
import { api } from '../services/api';
import { useAutenticacao } from './AutenticacaoContext';

export type WhiteLabelConfig = {
  mentorId?: string;
  nomePlataforma: string;
  logoUrl: string | null;
  bannerUrl: string | null;
  corPrimaria: string;
  corSecundaria: string;
  mensagemBoasVindas: string | null;
};

const CONFIG_PADRAO: WhiteLabelConfig = {
  nomePlataforma: 'Chega Junto Concurseiro',
  logoUrl: null,
  bannerUrl: null,
  corPrimaria: '#4f8ef7',
  corSecundaria: '#7c5cfc',
  mensagemBoasVindas: 'Análise completa da sua preparação'
};

type WhiteLabelContextType = {
  config: WhiteLabelConfig;
  carregando: boolean;
  atualizarConfig: (novosDados: Partial<WhiteLabelConfig>) => void;
  recarregar: () => Promise<void>;
};

const WhiteLabelContext = createContext<WhiteLabelContextType | undefined>(undefined);

export const WhiteLabelProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { usuario } = useAutenticacao();
  const [config, setConfig] = useState<WhiteLabelConfig>(CONFIG_PADRAO);
  const [carregando, setCarregando] = useState(false);

  const aplicarCoresCSS = (corPrimaria: string, corSecundaria: string) => {
    const root = document.documentElement;
    root.style.setProperty('--primary', corPrimaria);
    root.style.setProperty('--primary-hover', corPrimaria);
    root.style.setProperty('--accent', corSecundaria);
    root.style.setProperty('--brand-color', corPrimaria);
  };

  const carregarWhiteLabel = async () => {
    if (!usuario) {
      setConfig(CONFIG_PADRAO);
      aplicarCoresCSS(CONFIG_PADRAO.corPrimaria, CONFIG_PADRAO.corSecundaria);
      return;
    }

    setCarregando(true);
    try {
      let res: any;
      if (usuario.papel === 'mentor') {
        res = await api.obterWhiteLabelMentor();
      } else {
        res = await api.obterWhiteLabelAluno('eu');
      }

      const conf = res?.configuracao || res?.dados?.configuracao;

      if (conf) {
        const novaConfig: WhiteLabelConfig = {
          nomePlataforma: conf.nomePlataforma || 'Chega Junto Concurseiro',
          logoUrl: conf.logoUrl || null,
          bannerUrl: conf.bannerUrl || null,
          corPrimaria: conf.corPrimaria || '#4f8ef7',
          corSecundaria: conf.corSecundaria || '#7c5cfc',
          mensagemBoasVindas: conf.mensagemBoasVindas || 'Análise completa da sua preparação'
        };
        setConfig(novaConfig);
        aplicarCoresCSS(novaConfig.corPrimaria, novaConfig.corSecundaria);
      } else {
        setConfig(CONFIG_PADRAO);
        aplicarCoresCSS(CONFIG_PADRAO.corPrimaria, CONFIG_PADRAO.corSecundaria);
      }
    } catch {
      setConfig(CONFIG_PADRAO);
      aplicarCoresCSS(CONFIG_PADRAO.corPrimaria, CONFIG_PADRAO.corSecundaria);
    } finally {
      setCarregando(false);
    }
  };

  useEffect(() => {
    carregarWhiteLabel();
  }, [usuario?.id, usuario?.papel]);

  const atualizarConfig = (novosDados: Partial<WhiteLabelConfig>) => {
    setConfig((prev) => {
      const proxima = { ...prev, ...novosDados };
      aplicarCoresCSS(proxima.corPrimaria, proxima.corSecundaria);
      return proxima;
    });
  };

  return (
    <WhiteLabelContext.Provider
      value={{
        config,
        carregando,
        atualizarConfig,
        recarregar: carregarWhiteLabel
      }}
    >
      {children}
    </WhiteLabelContext.Provider>
  );
};

export const useWhiteLabel = () => {
  const context = useContext(WhiteLabelContext);
  if (!context) {
    throw new Error('useWhiteLabel deve ser usado dentro de um WhiteLabelProvider');
  }
  return context;
};
