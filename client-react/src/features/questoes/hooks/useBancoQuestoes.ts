import { useState, useEffect, useCallback } from 'react';
import { useData } from '../../../context/DataContext';
import { questoesService } from '../../../services/questoesService';

export interface FiltrosQuestoes {
  disciplina: string;
  assunto: string;
  tipo: string;
  situacao: '' | 'nao_respondidas' | 'errei' | 'acertei';
}

export function useBancoQuestoes() {
  const { activeContestId, alunoId } = useData();
  const [todasQuestoes, setTodasQuestoes] = useState<any[]>([]);
  const [estatisticas, setEstatisticas] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [erro, setErro] = useState<string | null>(null);
  const [pagina, setPagina] = useState(1);
  const [paginacao, setPaginacao] = useState({ pagina: 1, limite: 12, total: 0, totalPaginas: 1 });
  const [filtros, setFiltros] = useState<FiltrosQuestoes>({
    disciplina: '',
    assunto: '',
    tipo: '',
    situacao: '',
  });

  const carregar = useCallback(async () => {
    setLoading(true);
    setErro(null);
    try {
      const [resQuestoes, resStats] = await Promise.all([
        questoesService.listarBancoQuestoes({
          concursoId: activeContestId,
          alunoId,
          pagina,
          limite: 12,
        }),
        questoesService.estatisticasBancoQuestoes(alunoId, activeContestId),
      ]);
      
      const lista = Array.isArray(resQuestoes) 
        ? resQuestoes 
        : resQuestoes?.questoes || resQuestoes?.itens || [];
        
      setTodasQuestoes(lista);
      setPaginacao(resQuestoes?.paginacao || { pagina: 1, limite: 12, total: lista.length, totalPaginas: 1 });
      setEstatisticas(resStats);
    } catch (e: any) {
      setErro(e?.message || 'Erro ao carregar banco de questões.');
    } finally {
      setLoading(false);
    }
  }, [activeContestId, alunoId, pagina]);

  useEffect(() => {
    carregar();
  }, [carregar]);

  const responder = async (questaoId: string, respostaAluno: string) => {
    try {
      const res = await questoesService.responderQuestaoBanco(alunoId, {
        questaoId,
        respostaAluno,
        concursoId: activeContestId,
      });

      // Atualizar localmente a questão respondida
      setTodasQuestoes((prev) =>
        prev.map((q) =>
          q.id === questaoId
            ? {
                ...q,
                respondida: true,
                ultimaResposta: respostaAluno,
                ultimoAcerto: res?.correto,
              }
            : q
        )
      );

      // Atualiza estatísticas em background
      questoesService.estatisticasBancoQuestoes(alunoId, activeContestId)
        .then((stats) => setEstatisticas(stats))
        .catch(() => undefined);

      return res;
    } catch (e: any) {
      alert(e?.message || 'Erro ao responder questão.');
    }
  };

  const historicoQuestao = async (questaoId: string) => {
    const resposta = await questoesService.historicoQuestaoBanco(alunoId, questaoId, activeContestId);
    return Array.isArray(resposta) ? resposta : resposta?.historico || [];
  };

  return {
    questoes: todasQuestoes,
    todasQuestoes,
    estatisticas,
    loading,
    erro,
    filtros,
    setFiltros,
    carregar,
    responder,
    historicoQuestao,
    pagina: paginacao.pagina,
    paginacao,
    irParaPagina: setPagina,
  };
}
