import { useState, useEffect, useCallback } from 'react';
import { useData } from '../../../context/DataContext';
import { useAutenticacao } from '../../../context/AutenticacaoContext';
import { cronogramaService } from '../../../services/cronogramaService';

export interface CronogramaConfig {
  tipo: '' | 'agendado' | 'ciclo';
  cicloModo: 'sugerido' | 'livre';
  agendadoSetupDone: boolean;
  agendadoSetupStep: number;
  materiasSelecionadas: Record<string, boolean>;
  materiaAfinidade: Record<string, number>;
  materiaPrioridades: Record<string, string>;
  ritmo: 'conservador' | 'equilibrado' | 'intensivo';
  horas: Record<string, number>;
  minutosTopico: number;
  maxTopicosDia: number;
  repeticoesEdital: number;
  focusAlertsCustomEnabled: boolean;
  focusAlertsCustomGoal: number;
}

export const DEFAULT_CONFIG: CronogramaConfig = {
  tipo: '',
  cicloModo: 'sugerido',
  agendadoSetupDone: false,
  agendadoSetupStep: 1,
  materiasSelecionadas: {},
  materiaAfinidade: {},
  materiaPrioridades: {},
  ritmo: 'equilibrado',
  horas: { seg: 2, ter: 2, qua: 2, qui: 2, sex: 2, sab: 0, dom: 0 },
  minutosTopico: 60,
  maxTopicosDia: 0,
  repeticoesEdital: 1,
  focusAlertsCustomEnabled: false,
  focusAlertsCustomGoal: 80,
};

export function useCronograma() {
  const { activeContestId, alunoId } = useData();
  const { usuario } = useAutenticacao();
  const [config, setConfig] = useState<CronogramaConfig>(DEFAULT_CONFIG);
  const [cronograma, setCronograma] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [erro, setErro] = useState<string | null>(null);

  const carregar = useCallback(async () => {
    if (!activeContestId) {
      setLoading(false);
      return;
    }
    setLoading(true);
    setErro(null);
    try {
      const res = await cronogramaService.obterCronograma(alunoId, activeContestId);
      const dadosCronograma = res?.cronograma || res;
      if (dadosCronograma) {
        setCronograma(dadosCronograma);
        if (dadosCronograma.configuracao) {
          setConfig({ ...DEFAULT_CONFIG, ...dadosCronograma.configuracao });
        }
      }
    } catch (e: any) {
      setErro(e?.message || 'Erro ao carregar cronograma inteligente.');
    } finally {
      setLoading(false);
    }
  }, [activeContestId, alunoId]);

  useEffect(() => {
    carregar();
  }, [carregar]);

  const reprogramar = async () => {
    setLoading(true);
    try {
      await cronogramaService.reprogramarCronograma(alunoId, { concursoId: activeContestId });
      await carregar();
    } catch (e: any) {
      alert(e?.message || 'Erro ao reprogramar pendentes.');
    } finally {
      setLoading(false);
    }
  };

  const gerarCronograma = async (configParaGerar?: any) => {
    if (!activeContestId) throw new Error('Nenhum concurso selecionado');
    setLoading(true);
    setErro(null);
    try {
      const cfg = configParaGerar || config;
      const payload = {
        concursoId: activeContestId,
        tipo: cfg.tipo || 'ciclo_inteligente',
        configuracao: {
          tipo: cfg.tipo || 'ciclo_inteligente',
          cicloModo: cfg.cicloModo || 'adaptativo',
          materiasSelecionadas: cfg.materiasSelecionadas || {},
          materiaAfinidade: cfg.materiaAfinidade || {},
          materiaPrioridades: cfg.materiaPrioridades || {},
          ritmo: cfg.ritmo || 'moderado',
          horas: cfg.horas || { seg: 3, ter: 3, qua: 3, qui: 3, sex: 3, sab: 5, dom: 4 },
          minutosTopico: Number(cfg.minutosTopico) || 60,
          maxTopicosDia: Number(cfg.maxTopicosDia) || 4,
          alertaMetaHabilitado: Boolean(cfg.focusAlertsCustomEnabled),
          metaAcertos: Number(cfg.focusAlertsCustomGoal || cfg.metaAcertos) || 80,
          repeticoesEdital: Number(cfg.repeticoesEdital) || 1
        }
      };
      await cronogramaService.gerarCronograma(alunoId, payload);
      await carregar();
    } catch (e: any) {
      setErro(e?.message || 'Erro ao gerar cronograma.');
      throw e;
    } finally {
      setLoading(false);
    }
  };

  return {
    config,
    setConfig,
    cronograma,
    loading,
    erro,
    carregar,
    reprogramar,
    gerarCronograma,
    alunoId,
    activeContestId,
    usuario,
  };
}
