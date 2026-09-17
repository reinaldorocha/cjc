import React from 'react';
import { useBancoQuestoes } from '../hooks/useBancoQuestoes';
import { FiltroQuestoesBar } from '../components/FiltroQuestoesBar';
import { QuestaoCard } from '../components/QuestaoCard';
import { Spinner } from '../../../components/ui/Spinner';
import { EmptyState } from '../../../components/ui/EmptyState';
import { Card } from '../../../components/ui/Card';
import { useData } from '../../../context/DataContext';
import { Pagination } from '../../../components/ui/Pagination';

export const BancoQuestoesPage: React.FC = () => {
  const { getArray } = useData();
  const materias = getArray('materias');
  const { todasQuestoes, estatisticas, loading, erro, filtros, setFiltros, responder, historicoQuestao, pagina, paginacao, irParaPagina } = useBancoQuestoes();

  const limparFiltros = () => {
    setFiltros({ disciplina: '', assunto: '', tipo: '', situacao: '' });
  };

  const filtroAtivo = Boolean(filtros.disciplina || filtros.assunto || filtros.tipo || filtros.situacao);

  // Filtragem local resiliente e case-insensitive
  const questoesFiltradas = todasQuestoes.filter((q) => {
    if (filtros.disciplina && q.disciplina?.trim().toLowerCase() !== filtros.disciplina.trim().toLowerCase()) {
      return false;
    }
    if (filtros.assunto && q.assunto?.trim().toLowerCase() !== filtros.assunto.trim().toLowerCase()) {
      return false;
    }
    if (filtros.tipo && q.tipo !== filtros.tipo) {
      return false;
    }
    if (filtros.situacao === 'nao_respondidas' && q.respondida) {
      return false;
    }
    if (filtros.situacao === 'errei' && (!q.respondida || q.ultimoAcerto !== false)) {
      return false;
    }
    if (filtros.situacao === 'acertei' && (!q.respondida || q.ultimoAcerto !== true)) {
      return false;
    }
    return true;
  });

  const totalResolvidas = estatisticas?.totalResolvidas || estatisticas?.totalRespondidas || 0;
  const totalAcertos = estatisticas?.totalAcertos || 0;
  const taxaAcerto = estatisticas?.taxaAcerto !== undefined && estatisticas?.taxaAcerto !== null
    ? Number(estatisticas.taxaAcerto).toFixed(1)
    : '0.0';

  return (
    <div className="ui-container" style={{ padding: '24px 0' }}>
      <div className="ui-page-header">
        <h1 className="ui-page-title">Banco de Questões</h1>
        <p className="ui-page-subtitle">Pratique com questões reais de provas e acompanhe seu rendimento.</p>
      </div>

      {estatisticas && (
        <Card variant="glass" style={{ marginBottom: '24px', display: 'flex', gap: '24px', flexWrap: 'wrap', alignItems: 'center' }}>
          <div>Resolvidas: <strong>{totalResolvidas}</strong></div>
          <div>Acertos: <strong style={{ color: 'var(--status-success)' }}>{totalAcertos}</strong></div>
          <div>Aproveitamento: <strong style={{ color: 'var(--accent-primary)' }}>{taxaAcerto}%</strong></div>
        </Card>
      )}

      <FiltroQuestoesBar
        filtros={filtros}
        materias={materias}
        questoes={todasQuestoes}
        onChange={setFiltros}
        onLimpar={limparFiltros}
      />

      <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', color: 'var(--text-secondary)', fontSize: '0.9rem' }}>
        <span>
          Exibindo <strong>{questoesFiltradas.length}</strong> {questoesFiltradas.length === 1 ? 'questão' : 'questões'}
          {filtroAtivo && ` para o filtro selecionado`}
        </span>
        {filtroAtivo && (
          <button
            type="button"
            onClick={limparFiltros}
            style={{
              background: 'none',
              border: 'none',
              color: 'var(--accent-primary)',
              cursor: 'pointer',
              fontSize: '0.85rem',
              fontWeight: 600,
            }}
          >
            🔄 Limpar filtros
          </button>
        )}
      </div>

      {loading ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: '60px' }}>
          <Spinner size="lg" />
        </div>
      ) : erro ? (
        <Card style={{ textAlign: 'center', padding: '32px' }}>
          <p style={{ color: 'var(--status-danger)' }}>{erro}</p>
        </Card>
      ) : questoesFiltradas.length === 0 ? (
        <EmptyState
          icon={filtroAtivo ? "🎯" : "📚"}
          title={filtroAtivo ? "Nenhuma questão encontrada para o filtro selecionado" : "Nenhuma questão cadastrada no momento"}
          description={filtroAtivo ? "Tente selecionar outro assunto ou disciplina no filtro acima." : "As questões adicionadas pelo seu mentor aparecerão aqui."}
          actionLabel={filtroAtivo ? "Limpar Filtros" : undefined}
          onAction={filtroAtivo ? limparFiltros : undefined}
        />
      ) : (
        <div>
          {questoesFiltradas.map((q, idx) => (
            <QuestaoCard
              key={q.id || idx}
              questao={q}
              numeroQuestao={idx + 1}
              onResponder={responder}
              onHistorico={historicoQuestao}
            />
          ))}
          <Pagination pagina={pagina} totalPaginas={paginacao.totalPaginas} total={paginacao.total} onChange={irParaPagina} />
        </div>
      )}
    </div>
  );
};
