import { useState, useEffect, useCallback } from 'react';
import { useData } from '../../../context/DataContext';
import { api } from '../../../services/api';

export function useSimulados() {
  const { activeContestId, alunoId } = useData();
  const [simulados, setSimulados] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [erro, setErro] = useState<string | null>(null);

  const carregar = useCallback(async () => {
    setLoading(true);
    setErro(null);
    try {
      const res = await api.listarSimulados(alunoId, activeContestId);
      const lista = Array.isArray(res) ? res : res?.simulados || [];
      setSimulados(lista);
    } catch (e: any) {
      setErro(e?.message || 'Erro ao carregar simulados.');
    } finally {
      setLoading(false);
    }
  }, [activeContestId, alunoId]);

  useEffect(() => {
    carregar();
  }, [carregar]);

  return {
    simulados,
    loading,
    erro,
    carregar,
  };
}
