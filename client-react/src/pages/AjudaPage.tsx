import React, { useState } from 'react';

const topics = [
  { icon: '⏱', title: 'Cronômetro e Pomodoro', text: 'Use o botão flutuante para alternar entre tempo livre e blocos de foco. A sessão pode ser vinculada até o nível de subtópico.', items: ['Escolha blocos de 15, 25, 45 ou 60 minutos.', 'O alerta sonoro avisa quando o Pomodoro termina.', 'Questões, acertos e anotações alimentam suas estatísticas.'] },
  { icon: '📝', title: 'Edital e subtópicos', text: 'O Edital Verticalizado organiza matéria, tópico e subtópico com progresso, materiais e desempenho.', items: ['Use as setas para ajustar a ordem dos conteúdos.', 'Marcar um item como estudado agenda sua primeira revisão.', 'Registre horas e questões diretamente em qualquer item.'] },
  { icon: '✨', title: 'Importação assistida por IA', text: 'Na tela do Edital, copie a instrução, envie o edital à sua IA preferida e cole o JSON devolvido.', items: ['Revise a estrutura antes de começar os estudos.', 'Matérias, tópicos e subtópicos são criados de uma só vez.', 'O importador preserva a ordem recebida.'] },
  { icon: '🔗', title: 'Materiais e Smart Links', text: 'Vincule sites, vídeos, PDFs e cadernos de questões aos itens do edital.', items: ['Tópicos com o mesmo nome em outros concursos são detectados.', 'Importe os links encontrados sem duplicar URLs.', 'Os materiais continuam associados ao perfil e ao backup.'] },
  { icon: '🔄', title: 'Revisão espaçada', text: 'A tela de Revisões reúne o que venceu, o que é para hoje e os próximos agendamentos.', items: ['Registre o desempenho ao concluir.', 'Use +1d quando precisar reagendar.', 'Acompanhe a taxa de conclusão no resumo.'] },
  { icon: '🎴', title: 'Flashcards', text: 'Crie perguntas e respostas vinculadas às matérias e pratique com recordação ativa.', items: ['Abra a sessão de revisão para estudar os cards.', 'Marque se lembrou ou não da resposta.', 'O histórico mostra seu desempenho acumulado.'] },
  { icon: '📁', title: 'Perfis e backups', text: 'Cada perfil mantém seus próprios concursos. O menu do perfil permite alternar, exportar e importar seus dados.', items: ['O app funciona localmente se o servidor estiver indisponível.', 'Com o servidor ativo, alterações geram backups automáticos.', 'Exporte um JSON para manter uma cópia externa.'] }
];

export const AjudaPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const filtered = topics.filter(t => `${t.title} ${t.text} ${t.items.join(' ')}`.toLowerCase().includes(query.toLowerCase()));
  return <div>
    <div className="page-heading"><div><h1>Ajuda e documentação</h1><p>Aprenda os fluxos essenciais do Chega Junto Concurseiro.</p></div></div>
    <input className="form-control help-search" type="search" value={query} onChange={e => setQuery(e.target.value)} placeholder="Buscar na ajuda…" aria-label="Buscar na ajuda" />
    <div className="help-grid">{filtered.map(topic => <article className="card-base" key={topic.title}>
      <h2><span>{topic.icon}</span>{topic.title}</h2><p>{topic.text}</p><ul>{topic.items.map(item => <li key={item}>{item}</li>)}</ul>
    </article>)}</div>
    {!filtered.length && <div className="empty-state"><span>🔎</span><h2>Nenhum tópico encontrado</h2><p>Tente uma busca mais curta.</p></div>}
  </div>;
};
