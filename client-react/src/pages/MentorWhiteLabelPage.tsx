import React, { useEffect, useState } from 'react';
import { api } from '../services/api';
import { useWhiteLabel } from '../context/WhiteLabelContext';

const PRESETS_CORES = [
  { nome: 'Azul Padrão', primaria: '#4f8ef7', secundaria: '#7c5cfc' },
  { nome: 'Laranja Vibrante', primaria: '#ff6b35', secundaria: '#f7c59f' },
  { nome: 'Esmeralda', primaria: '#10b981', secundaria: '#059669' },
  { nome: 'Roxo Imperial', primaria: '#8b5cf6', secundaria: '#ec4899' },
  { nome: 'Vermelho Foco', primaria: '#ef4444', secundaria: '#f97316' },
  { nome: 'Ciano Cyber', primaria: '#06b6d4', secundaria: '#3b82f6' }
];

export const MentorWhiteLabelPage: React.FC = () => {
  const { config, atualizarConfig, recarregar } = useWhiteLabel();

  const [form, setForm] = useState({
    nomePlataforma: config.nomePlataforma || 'Chega Junto Concurseiro',
    logoUrl: config.logoUrl || '',
    bannerUrl: config.bannerUrl || '',
    corPrimaria: config.corPrimaria || '#4f8ef7',
    corSecundaria: config.corSecundaria || '#7c5cfc',
    mensagemBoasVindas: config.mensagemBoasVindas || 'Análise completa da sua preparação'
  });

  const [salvando, setSalvando] = useState(false);
  const [enviandoLogo, setEnviandoLogo] = useState(false);
  const [enviandoBanner, setEnviandoBanner] = useState(false);
  const [mensagemSucesso, setMensagemSucesso] = useState('');
  const [erro, setErro] = useState('');

  useEffect(() => {
    setForm({
      nomePlataforma: config.nomePlataforma || 'Chega Junto Concurseiro',
      logoUrl: config.logoUrl || '',
      bannerUrl: config.bannerUrl || '',
      corPrimaria: config.corPrimaria || '#4f8ef7',
      corSecundaria: config.corSecundaria || '#7c5cfc',
      mensagemBoasVindas: config.mensagemBoasVindas || 'Análise completa da sua preparação'
    });
  }, [config]);

  const handleChange = (campo: string, valor: string) => {
    const novoForm = { ...form, [campo]: valor };
    setForm(novoForm);
    // Atualização em tempo real na tela!
    atualizarConfig(novoForm);
  };

  const aplicarPreset = (primaria: string, secundaria: string) => {
    const novoForm = { ...form, corPrimaria: primaria, corSecundaria: secundaria };
    setForm(novoForm);
    atualizarConfig(novoForm);
  };

  const handleUploadLogo = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 5 * 1024 * 1024) {
      setErro('A imagem do logo deve ter no máximo 5MB.');
      return;
    }

    setEnviandoLogo(true);
    setErro('');
    try {
      const res = await api.uploadLogoWhiteLabel(file);
      const url = res?.url || res?.dados?.url;
      if (url) {
        handleChange('logoUrl', url);
      }
    } catch {
      setErro('Não foi possível enviar o logo. Tente novamente.');
    } finally {
      setEnviandoLogo(false);
    }
  };

  const handleUploadBanner = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    if (file.size > 5 * 1024 * 1024 || !file.type.startsWith('image/')) {
      setErro('Selecione uma imagem de até 5 MB para o banner.');
      return;
    }
    setEnviandoBanner(true); setErro('');
    try {
      const res = await api.uploadBannerWhiteLabel(file);
      const url = res?.url || res?.dados?.url;
      if (url) handleChange('bannerUrl', url);
    } catch {
      setErro('Não foi possível enviar o banner. Tente novamente.');
    } finally { setEnviandoBanner(false); }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSalvando(true);
    setMensagemSucesso('');
    setErro('');

    try {
      const res = await api.salvarWhiteLabelMentor({
        nomePlataforma: form.nomePlataforma,
        logoUrl: form.logoUrl || null,
        bannerUrl: form.bannerUrl || null,
        corPrimaria: form.corPrimaria,
        corSecundaria: form.corSecundaria,
        mensagemBoasVindas: form.mensagemBoasVindas || null
      });

      const conf = res?.configuracao || res?.dados?.configuracao;
      if (conf) {
        atualizarConfig(conf);
      }

      setMensagemSucesso('Marca da mentoria salva com sucesso! Seus alunos já estão vendo a nova identidade.');
      await recarregar();
      setTimeout(() => setMensagemSucesso(''), 4000);
    } catch (err: any) {
      setErro(err.message || 'Erro ao salvar a marca da mentoria.');
    } finally {
      setSalvando(false);
    }
  };

  return (
    <div style={{ maxWidth: '1000px', margin: '0 auto', paddingBottom: '40px' }}>
      <div className="page-heading">
        <div>
          <h1>⚙️ Marca & White Label da Mentoria</h1>
          <p>Personalize o nome da plataforma, logo e paleta de cores exclusiva para os seus alunos.</p>
        </div>
      </div>

      {mensagemSucesso && <div className="success-banner">{mensagemSucesso}</div>}
      {erro && <div className="form-error">{erro}</div>}

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 380px', gap: '24px', alignItems: 'start' }}>
        {/* FORMULÁRIO DE CUSTOMIZAÇÃO */}
        <form className="card-base" onSubmit={handleSubmit} style={{ padding: '24px' }}>
          <h2 style={{ fontSize: '18px', fontWeight: 800, marginBottom: '16px' }}>🎨 Identidade Visual</h2>

          <div style={{ display: 'grid', gap: '18px' }}>
            <label>
              <span className="field-label">Nome da Plataforma / Mentoria *</span>
              <input
                className="form-control"
                type="text"
                placeholder="Ex: Mentoria AprovaJá, Chega Junto Concurseiro..."
                value={form.nomePlataforma}
                onChange={(e) => handleChange('nomePlataforma', e.target.value)}
                required
              />
            </label>

            <label>
              <span className="field-label">Mensagem de Boas-Vindas ou Subtítulo</span>
              <input
                className="form-control"
                type="text"
                placeholder="Ex: Acompanhamento de elite para a sua aprovação"
                value={form.mensagemBoasVindas}
                onChange={(e) => handleChange('mensagemBoasVindas', e.target.value)}
              />
            </label>

            <div>
              <span className="field-label">Logo da Mentoria</span>
              <div style={{ display: 'flex', gap: '12px', alignItems: 'center', marginTop: '6px' }}>
                {form.logoUrl ? (
                  <img
                    src={form.logoUrl}
                    alt="Logo Marca"
                    style={{ maxHeight: '48px', maxWidth: '160px', objectFit: 'contain', borderRadius: '6px', border: '1px solid var(--border)' }}
                  />
                ) : (
                  <div style={{ width: '48px', height: '48px', borderRadius: '8px', background: 'var(--bg3)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '20px' }}>
                    🏆
                  </div>
                )}

                <div style={{ flex: 1 }}>
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleUploadLogo}
                    disabled={enviandoLogo}
                    style={{ fontSize: '12px' }}
                  />
                  {enviandoLogo && <span style={{ fontSize: '11px', color: 'var(--text3)', marginTop: '4px', display: 'block' }}>Enviando logo...</span>}
                  {form.logoUrl && !enviandoLogo && (
                    <button
                      type="button"
                      className="btn-secondary danger-text"
                      onClick={() => handleChange('logoUrl', '')}
                      style={{ padding: '2px 8px', fontSize: '11px', marginTop: '4px' }}
                    >
                      Remover Logo
                    </button>
                  )}
                </div>
              </div>
            </div>

            <div>
              <span className="field-label">Banner do dashboard</span>
              <p style={{ margin: '4px 0 8px', color: 'var(--text3)', fontSize: '12px' }}>Imagem opcional exibida no topo da página inicial do aluno. Recomendado: formato horizontal 3:1.</p>
              {form.bannerUrl && <img src={form.bannerUrl} alt="Prévia do banner" style={{ display: 'block', width: '100%', maxHeight: '150px', objectFit: 'cover', borderRadius: '8px', border: '1px solid var(--border)', marginBottom: '8px' }} />}
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
                <input type="file" accept="image/*" onChange={handleUploadBanner} disabled={enviandoBanner} style={{ fontSize: '12px' }} />
                {enviandoBanner && <span style={{ fontSize: '11px', color: 'var(--text3)' }}>Enviando banner...</span>}
                {form.bannerUrl && !enviandoBanner && <button type="button" className="btn-secondary danger-text" onClick={() => handleChange('bannerUrl', '')} style={{ padding: '4px 8px', fontSize: '11px' }}>Remover banner</button>}
              </div>
            </div>

            <div>
              <span className="field-label" style={{ marginBottom: '8px', display: 'block' }}>Paletas Recomendadas</span>
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                {PRESETS_CORES.map((preset) => (
                  <button
                    key={preset.nome}
                    type="button"
                    onClick={() => aplicarPreset(preset.primaria, preset.secundaria)}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: '6px',
                      padding: '6px 12px',
                      borderRadius: '20px',
                      border: '1px solid var(--border)',
                      background: 'var(--bg-card)',
                      fontSize: '12px',
                      cursor: 'pointer'
                    }}
                  >
                    <span style={{ width: '12px', height: '12px', borderRadius: '50%', background: preset.primaria, display: 'inline-block' }} />
                    {preset.nome}
                  </button>
                ))}
              </div>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
              <label>
                <span className="field-label">Cor Primária *</span>
                <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                  <input
                    type="color"
                    value={form.corPrimaria}
                    onChange={(e) => handleChange('corPrimaria', e.target.value)}
                    style={{ width: '40px', height: '38px', borderRadius: '6px', border: 'none', cursor: 'pointer', background: 'transparent' }}
                  />
                  <input
                    className="form-control"
                    type="text"
                    value={form.corPrimaria}
                    onChange={(e) => handleChange('corPrimaria', e.target.value)}
                  />
                </div>
              </label>

              <label>
                <span className="field-label">Cor Secundária / Destaques</span>
                <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                  <input
                    type="color"
                    value={form.corSecundaria}
                    onChange={(e) => handleChange('corSecundaria', e.target.value)}
                    style={{ width: '40px', height: '38px', borderRadius: '6px', border: 'none', cursor: 'pointer', background: 'transparent' }}
                  />
                  <input
                    className="form-control"
                    type="text"
                    value={form.corSecundaria}
                    onChange={(e) => handleChange('corSecundaria', e.target.value)}
                  />
                </div>
              </label>
            </div>

            <div style={{ marginTop: '12px', display: 'flex', justifyContent: 'flex-end' }}>
              <button type="submit" className="btn-primary" disabled={salvando} style={{ padding: '10px 24px', fontSize: '14px', background: form.corPrimaria, borderColor: form.corPrimaria }}>
                {salvando ? 'Salvando Marca…' : '💾 Salvar Identidade Visual'}
              </button>
            </div>
          </div>
        </form>

        {/* PRÉ-VISUALIZAÇÃO AO VIVO */}
        <div className="card-base" style={{ padding: '20px', position: 'sticky', top: '20px' }}>
          <h3 style={{ fontSize: '15px', fontWeight: 800, marginBottom: '14px', display: 'flex', alignItems: 'center', gap: '6px' }}>
            👁️ Pré-visualização do Aluno
          </h3>

          {/* MOCKUP DO HEADER */}
          <div style={{ background: 'var(--bg2)', padding: '12px 16px', borderRadius: '10px', marginBottom: '14px', border: '1px solid var(--border)' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
              {form.logoUrl ? (
                <img src={form.logoUrl} alt="Logo" style={{ maxHeight: '28px', maxWidth: '80px', objectFit: 'contain' }} />
              ) : (
                <span style={{ fontSize: '20px' }}>🏆</span>
              )}
              <div>
                <div style={{ fontSize: '13px', fontWeight: 800, color: 'var(--text)' }}>
                  {form.nomePlataforma || 'Chega Junto Concurseiro'}
                </div>
                <div style={{ fontSize: '10px', color: 'var(--text3)' }}>
                  {form.mensagemBoasVindas || 'Painel do Aluno'}
                </div>
              </div>
            </div>
          </div>

          {/* MOCKUP DE BOTÕES E BADGES COM A COR DO MENTOR */}
          <div style={{ display: 'grid', gap: '10px' }}>
            <button
              type="button"
              style={{
                background: form.corPrimaria,
                color: '#ffffff',
                border: 'none',
                padding: '10px 16px',
                borderRadius: '8px',
                fontWeight: 700,
                fontSize: '13px',
                cursor: 'pointer'
              }}
            >
              Exemplo de Botão Principal
            </button>

            <div style={{ display: 'flex', gap: '8px' }}>
              <span style={{ background: `color-mix(in srgb, ${form.corPrimaria} 20%, transparent)`, color: form.corPrimaria, padding: '4px 10px', borderRadius: '6px', fontSize: '11px', fontWeight: 700 }}>
                Tag da Matéria
              </span>
              <span style={{ background: `color-mix(in srgb, ${form.corSecundaria} 20%, transparent)`, color: form.corSecundaria, padding: '4px 10px', borderRadius: '6px', fontSize: '11px', fontWeight: 700 }}>
                Concurso Ativo
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
