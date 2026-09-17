import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../../services/api';
import { Card } from '../../../components/ui/Card';
import { Button } from '../../../components/ui/Button';
import { Badge } from '../../../components/ui/Badge';
import { Spinner } from '../../../components/ui/Spinner';
import { Modal } from '../../../components/ui/Modal';
import { Pagination } from '../../../components/ui/Pagination';

export const MentorAlunosPage: React.FC = () => {
  const navigate = useNavigate();
  const [alunos, setAlunos] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [pagina, setPagina] = useState(1);
  const [paginacao, setPaginacao] = useState({ pagina: 1, limite: 12, total: 0, totalPaginas: 1 });

  // States para edicao de aluno
  const [alunoEdicao, setAlunoEdicao] = useState<any | null>(null);
  const [dadosEdicao, setDadosEdicao] = useState({
    nome: '',
    senha: '',
    telefone: '',
    permiteCronogramaInteligente: false,
    dataExpiracaoPlano: '',
  });

  // Modal para cadastro de novo aluno
  const [isNovoAlunoModalOpen, setIsNovoAlunoModalOpen] = useState(false);
  const [alunoAtribuicao, setAlunoAtribuicao] = useState<any | null>(null);
  const [catalogoConcursos, setCatalogoConcursos] = useState<any[]>([]);
  const [catalogoEditais, setCatalogoEditais] = useState<any[]>([]);
  const [concursoSelecionado, setConcursoSelecionado] = useState('');
  const [editalSelecionado, setEditalSelecionado] = useState('');
  const [atribuicoesAtuais, setAtribuicoesAtuais] = useState<any[]>([]);
  const [editandoAtribuicaoId, setEditandoAtribuicaoId] = useState('');
  const [atribuindo, setAtribuindo] = useState(false);
  const [erroAtribuicao, setErroAtribuicao] = useState('');
  const [novoAluno, setNovoAluno] = useState({
    nome: '',
    email: '',
    senha: '',
    telefone: '',
    permiteCronogramaInteligente: true,
    dataExpiracaoPlano: '',
  });

  const carregarAlunos = (paginaSolicitada = pagina) => {
    setLoading(true);
    api.listarAlunos(paginaSolicitada, 12)
      .then((res) => {
        const lista = Array.isArray(res) ? res : res?.alunos || [];
        setAlunos(lista);
        setPaginacao(res?.paginacao || { pagina: 1, limite: 12, total: lista.length, totalPaginas: 1 });
      })
      .catch(() => setAlunos([]))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    carregarAlunos();
  }, [pagina]);

  const abrirEdicao = (aluno: any) => {
    setAlunoEdicao(aluno);
    setDadosEdicao({
      nome: aluno.nome || '',
      senha: '',
      telefone: aluno.telefone || '',
      permiteCronogramaInteligente: !!aluno.permiteCronogramaInteligente,
      dataExpiracaoPlano: aluno.dataExpiracaoPlano || '',
    });
  };

  const salvarEdicao = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!alunoEdicao) return;
    try {
      const payload: any = {
        nome: dadosEdicao.nome,
        telefone: dadosEdicao.telefone,
        dataExpiracaoPlano: dadosEdicao.dataExpiracaoPlano,
        permiteCronogramaInteligente: dadosEdicao.permiteCronogramaInteligente,
      };
      if (dadosEdicao.senha.trim()) {
        payload.senha = dadosEdicao.senha.trim();
      }

      await api.configurarAluno(alunoEdicao.id, payload);
      alert('Perfil do aluno atualizado com sucesso!');
      setAlunoEdicao(null);
      carregarAlunos();
    } catch (err: any) {
      alert(err?.message || 'Erro ao atualizar perfil do aluno.');
    }
  };

  const handleCriarAluno = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!novoAluno.nome || !novoAluno.email || !novoAluno.senha) {
      alert('Preencha nome, e-mail e senha.');
      return;
    }
    try {
      await api.criarAlunoComoMentor(novoAluno);
      alert('Aluno cadastrado com sucesso!');
      setIsNovoAlunoModalOpen(false);
      setNovoAluno({
        nome: '',
        email: '',
        senha: '',
        telefone: '',
        permiteCronogramaInteligente: true,
        dataExpiracaoPlano: '',
      });
      carregarAlunos();
    } catch (err: any) {
      alert(err?.message || 'Erro ao cadastrar aluno.');
    }
  };

  const acessarAreaAluno = async (alunoId: string) => {
    try {
      const res = await api.listarConcursos(alunoId);
      const concursos = Array.isArray(res) ? res : (res as any)?.concursos || [];
      navigate(concursos[0]?.id ? `/mentor/alunos/${alunoId}/concursos/${concursos[0].id}/dashboard` : `/mentor/alunos/${alunoId}/concursos`);
    } catch { navigate(`/mentor/alunos/${alunoId}/concursos`); }
  };

  const abrirAtribuicao = async (aluno: any) => {
    setAlunoAtribuicao(aluno); setConcursoSelecionado(''); setEditalSelecionado(''); setCatalogoEditais([]); setAtribuicoesAtuais([]); setEditandoAtribuicaoId(''); setErroAtribuicao('');
    try {
      const [catalogo, vinculados] = await Promise.all([api.listarConcursosMentor(), api.listarConcursos(aluno.id)]);
      setCatalogoConcursos(catalogo.concursos || []);
      const lista = Array.isArray(vinculados) ? vinculados : (vinculados as any)?.concursos || [];
      setAtribuicoesAtuais(await Promise.all(lista.map(async (concurso: any) => ({ ...concurso, editais: ((await api.listarEditais(aluno.id, concurso.id)).editais || []) }))));
    }
    catch { setErroAtribuicao('Não foi possível carregar o catálogo de concursos.'); }
  };
  const escolherConcurso = async (id: string) => {
    setConcursoSelecionado(id); setEditalSelecionado(''); setCatalogoEditais([]); if (!id) return;
    try { setCatalogoEditais((await api.listarEditaisMentor(id)).editais || []); }
    catch { setErroAtribuicao('Não foi possível carregar os editais deste concurso.'); }
  };
  const atribuir = async (e: React.FormEvent) => {
    e.preventDefault();
    const concurso = catalogoConcursos.find((item) => item.id === concursoSelecionado);
    if (!alunoAtribuicao || !concurso || (editandoAtribuicaoId && !editalSelecionado)) { setErroAtribuicao(editandoAtribuicaoId ? 'Selecione o edital para substituir o atual.' : 'Selecione um concurso.'); return; }
    setAtribuindo(true); setErroAtribuicao('');
    try {
      if (editandoAtribuicaoId) {
        await api.atribuirEdital(alunoAtribuicao.id, { editalId: editalSelecionado, concursoId: editandoAtribuicaoId });
      } else {
        const criado = await api.criarConcurso(alunoAtribuicao.id, { concursoId: concurso.id, nome: concurso.nome, banca: concurso.banca, cargo: concurso.cargo || '', logotipo: concurso.logotipo || '', salario: concurso.salario ?? null, dataProva: concurso.dataProva || '', preEdital: !!concurso.preEdital, grupo: 'foco', ordem: 1 });
        if (editalSelecionado) await api.atribuirEdital(alunoAtribuicao.id, { editalId: editalSelecionado, concursoId: criado.id });
      }
      await abrirAtribuicao(alunoAtribuicao);
    } catch (err: any) { setErroAtribuicao(err?.message || 'Não foi possível atribuir o concurso e o edital.'); }
    finally { setAtribuindo(false); }
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '60px' }}>
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="ui-container" style={{ padding: '24px 0' }}>
      <div className="ui-page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <h1 className="ui-page-title">Meus Alunos</h1>
          <p className="ui-page-subtitle">Gerencie e acompanhe seus alunos vinculados à mentoria.</p>
        </div>
        <Button variant="primary" onClick={() => setIsNovoAlunoModalOpen(true)}>
          ＋ Novo Aluno
        </Button>
      </div>

      {alunos.length === 0 ? (
        <Card style={{ textAlign: 'center', padding: '48px' }}>
          <p style={{ color: 'var(--text-secondary)', marginBottom: '16px' }}>Nenhum aluno encontrado na sua mentoria.</p>
          <Button variant="primary" onClick={() => setIsNovoAlunoModalOpen(true)}>Cadastrar Primeiro Aluno</Button>
        </Card>
      ) : (
        <>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(320px, 1fr))', gap: '20px' }}>
            {alunos.map((a) => (
            <Card key={a.id} variant="default" style={{ display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
                  <div>
                    <h3 style={{ fontSize: '1.1rem', fontWeight: 600, color: 'var(--text-primary)' }}>{a.nome}</h3>
                    <p style={{ fontSize: '0.875rem', color: 'var(--text-secondary)' }}>{a.email}</p>
                  </div>
                  <Badge variant={a.permiteCronogramaInteligente ? 'success' : 'neutral'}>
                    {a.permiteCronogramaInteligente ? 'Cronograma OK' : 'Básico'}
                  </Badge>
                </div>

                <div style={{ fontSize: '0.85rem', color: 'var(--text-muted)', display: 'flex', flexDirection: 'column', gap: '6px', marginBottom: '16px' }}>
                  <div>📱 Telefone: {a.telefone || 'Não informado'}</div>
                  <div>📅 Expiração: {a.dataExpiracaoPlano ? new Date(a.dataExpiracaoPlano).toLocaleDateString('pt-BR') : 'Sem limite'}</div>
                </div>
              </div>

              <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap', borderTop: '1px solid var(--border-color)', paddingTop: '12px' }}>
                <Button
                  variant="primary"
                  size="sm"
                  onClick={() => acessarAreaAluno(a.id)}
                  style={{ flex: 1 }}
                >
                  🚀 Acessar como Aluno
                </Button>
                <Button variant="secondary" size="sm" onClick={() => void abrirAtribuicao(a)}>Atribuir edital</Button>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={() => abrirEdicao(a)}
                >
                  ✏️ Editar
                </Button>
              </div>
            </Card>
            ))}
          </div>
          <Pagination pagina={paginacao.pagina} totalPaginas={paginacao.totalPaginas} total={paginacao.total} onChange={setPagina} />
        </>
      )}

      {/* MODAL DE EDIÇÃO DO ALUNO */}
      <Modal
        isOpen={!!alunoEdicao}
        onClose={() => setAlunoEdicao(null)}
        title={`Editar Perfil: ${alunoEdicao?.nome}`}
      >
        <form onSubmit={salvarEdicao} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Nome Completo do Aluno *
            </label>
            <input
              type="text"
              required
              value={dadosEdicao.nome}
              onChange={(e) => setDadosEdicao({ ...dadosEdicao, nome: e.target.value })}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Nova Senha (deixe em branco para manter a atual)
            </label>
            <input
              type="password"
              placeholder="Digite a nova senha..."
              value={dadosEdicao.senha}
              onChange={(e) => setDadosEdicao({ ...dadosEdicao, senha: e.target.value })}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Telefone / WhatsApp
            </label>
            <input
              type="text"
              value={dadosEdicao.telefone}
              onChange={(e) => setDadosEdicao({ ...dadosEdicao, telefone: e.target.value })}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Data de Expiração do Plano
            </label>
            <input
              type="date"
              value={dadosEdicao.dataExpiracaoPlano ? dadosEdicao.dataExpiracaoPlano.split('T')[0] : ''}
              onChange={(e) => setDadosEdicao({ ...dadosEdicao, dataExpiracaoPlano: e.target.value })}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginTop: '8px' }}>
            <input
              type="checkbox"
              id="permiteCronograma"
              checked={dadosEdicao.permiteCronogramaInteligente}
              onChange={(e) => setDadosEdicao({ ...dadosEdicao, permiteCronogramaInteligente: e.target.checked })}
            />
            <label htmlFor="permiteCronograma" style={{ fontSize: '0.9rem', color: 'var(--text-primary)', cursor: 'pointer' }}>
              Habilitar Módulo Cronograma Inteligente
            </label>
          </div>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px', marginTop: '16px' }}>
            <Button type="button" variant="ghost" onClick={() => setAlunoEdicao(null)}>Cancelar</Button>
            <Button type="submit" variant="primary">Salvar Alterações</Button>
          </div>
        </form>
      </Modal>

      <Modal isOpen={!!alunoAtribuicao} onClose={() => !atribuindo && setAlunoAtribuicao(null)} title={`Concursos e editais · ${alunoAtribuicao?.nome || ''}`}>
        <form onSubmit={atribuir} style={{ display: 'grid', gap: '16px' }}>
          <p style={{ margin: 0, color: 'var(--text-secondary)', fontSize: '13px' }}>Atribuições atuais do aluno.</p>
          {!!atribuicoesAtuais.length && <div style={{ display: 'grid', gap: '8px' }}>{atribuicoesAtuais.map((item) => <div key={item.id} className="card-base" style={{ padding: '10px' }}><strong>{item.nome}</strong><small style={{ display: 'block', color: 'var(--text3)', margin: '3px 0 8px' }}>Edital: {item.editais?.[0]?.nome || item.editais?.[0]?.titulo || 'Nenhum edital vinculado'}</small><div style={{ display: 'flex', gap: '8px' }}><Button type="button" variant="secondary" size="sm" onClick={() => { setEditandoAtribuicaoId(item.id); void escolherConcurso(item.id); }}>Trocar edital</Button><Button type="button" variant="ghost" size="sm" disabled={atribuindo} onClick={async () => { if (confirm(`Remover ${item.nome} deste aluno?`)) { await api.desativarConcurso(alunoAtribuicao.id, item.id); await abrirAtribuicao(alunoAtribuicao); } }}>Excluir</Button></div></div>)}</div>}
          <p style={{ margin: 0, fontWeight: 700 }}>{editandoAtribuicaoId ? 'Trocar edital vinculado' : 'Nova atribuição'}</p>
          <label><span className="field-label">Concurso</span><select className="form-control" required value={concursoSelecionado} disabled={!!editandoAtribuicaoId} onChange={(e) => void escolherConcurso(e.target.value)}><option value="">Selecione</option>{catalogoConcursos.filter((c) => editandoAtribuicaoId || !atribuicoesAtuais.some((v) => v.id === c.id)).map((c) => <option key={c.id} value={c.id}>{c.nome} · {c.banca}</option>)}</select></label>
          <label><span className="field-label">Edital</span><select className="form-control" value={editalSelecionado} onChange={(e) => setEditalSelecionado(e.target.value)} disabled={!concursoSelecionado || !catalogoEditais.length}><option value="">{editandoAtribuicaoId ? 'Selecione' : 'Usar o edital mais recente'}</option>{catalogoEditais.map((edital) => <option key={edital.id} value={edital.id}>{edital.nome || edital.titulo || 'Edital'}</option>)}</select></label>
          {erroAtribuicao && <p className="form-error">{erroAtribuicao}</p>}
          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}><Button type="button" variant="ghost" disabled={atribuindo} onClick={() => { setEditandoAtribuicaoId(''); setConcursoSelecionado(''); setEditalSelecionado(''); }}>Limpar</Button><Button type="submit" variant="primary" disabled={atribuindo}>{atribuindo ? 'Salvando…' : editandoAtribuicaoId ? 'Salvar edital' : 'Atribuir ao aluno'}</Button></div>
        </form>
      </Modal>

      {/* MODAL DE NOVO ALUNO */}
      <Modal
        isOpen={isNovoAlunoModalOpen}
        onClose={() => setIsNovoAlunoModalOpen(false)}
        title="Cadastrar Novo Aluno"
      >
        <form onSubmit={handleCriarAluno} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Nome Completo *
            </label>
            <input
              type="text"
              required
              value={novoAluno.nome}
              onChange={(e) => setNovoAluno({ ...novoAluno, nome: e.target.value })}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              E-mail de Acesso *
            </label>
            <input
              type="email"
              required
              value={novoAluno.email}
              onChange={(e) => setNovoAluno({ ...novoAluno, email: e.target.value })}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Senha Inicial *
            </label>
            <input
              type="password"
              required
              value={novoAluno.senha}
              onChange={(e) => setNovoAluno({ ...novoAluno, senha: e.target.value })}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Telefone / WhatsApp
            </label>
            <input
              type="text"
              value={novoAluno.telefone}
              onChange={(e) => setNovoAluno({ ...novoAluno, telefone: e.target.value })}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '0.875rem', color: 'var(--text-secondary)', marginBottom: '4px' }}>
              Data de Expiração do Plano
            </label>
            <input
              type="date"
              value={novoAluno.dataExpiracaoPlano}
              onChange={(e) => setNovoAluno({ ...novoAluno, dataExpiracaoPlano: e.target.value })}
              style={{ width: '100%', padding: '8px 12px', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff' }}
            />
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginTop: '8px' }}>
            <input
              type="checkbox"
              id="novoPermiteCronograma"
              checked={novoAluno.permiteCronogramaInteligente}
              onChange={(e) => setNovoAluno({ ...novoAluno, permiteCronogramaInteligente: e.target.checked })}
            />
            <label htmlFor="novoPermiteCronograma" style={{ fontSize: '0.9rem', color: 'var(--text-primary)', cursor: 'pointer' }}>
              Habilitar Módulo Cronograma Inteligente
            </label>
          </div>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px', marginTop: '16px' }}>
            <Button type="button" variant="ghost" onClick={() => setIsNovoAlunoModalOpen(false)}>Cancelar</Button>
            <Button type="submit" variant="primary">Cadastrar Aluno</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};
