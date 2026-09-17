import { useState, useEffect, useCallback } from 'react';
import { useData } from '../../../context/DataContext';
import { flashcardsService } from '../../../services/flashcardsService';

export function useFlashcards() {
  const { activeContestId, alunoId } = useData();
  const [baralhos, setBaralhos] = useState<any[]>([]);
  const [cartoesPendentes, setCartoesPendentes] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [erro, setErro] = useState<string | null>(null);

  const carregar = useCallback(async () => {
    setLoading(true);
    setErro(null);
    try {
      const [resBaralhos, resCartoes] = await Promise.all([
        flashcardsService.listarBaralhos(alunoId, activeContestId),
        flashcardsService.listarCartoesPendentes(alunoId),
      ]);
      
      const listaBaralhos = Array.isArray(resBaralhos) ? resBaralhos : resBaralhos?.baralhos || [];
      const listaCartoes = Array.isArray(resCartoes) ? resCartoes : resCartoes?.cartoes || [];

      setBaralhos(listaBaralhos);
      setCartoesPendentes(listaCartoes);
    } catch (e: any) {
      setErro(e?.message || 'Erro ao carregar baralhos.');
    } finally {
      setLoading(false);
    }
  }, [activeContestId, alunoId]);

  useEffect(() => {
    carregar();
  }, [carregar]);

  const revisarCartao = async (cartaoId: string, qualidade: number) => {
    try {
      await flashcardsService.revisarCartao(alunoId, cartaoId, qualidade, activeContestId);
      await carregar();
    } catch (e: any) {
      alert(e?.message || 'Erro ao salvar revisão do cartão.');
    }
  };

  return {
    baralhos,
    cartoesPendentes,
    loading,
    erro,
    carregar,
    revisarCartao,
  };
}
