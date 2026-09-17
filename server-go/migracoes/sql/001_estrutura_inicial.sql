CREATE TABLE IF NOT EXISTS usuarios (
  id CHAR(36) PRIMARY KEY,
  nome VARCHAR(160) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  senha_hash VARCHAR(255) NOT NULL,
  papel ENUM('mestre','mentor','aluno') NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  ultimo_acesso_em DATETIME NULL,
  desativado_em DATETIME NULL,
  desativado_por CHAR(36) NULL,
  INDEX idx_usuarios_papel_ativo (papel, ativo)
);

CREATE TABLE IF NOT EXISTS sessoes_autenticacao (
  id CHAR(36) PRIMARY KEY,
  usuario_id CHAR(36) NOT NULL,
  token_hash CHAR(64) NOT NULL UNIQUE,
  csrf_hash CHAR(64) NOT NULL,
  expira_em DATETIME NOT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ultimo_uso_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  revogado_em DATETIME NULL,
  FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
  INDEX idx_sessoes_usuario (usuario_id, revogado_em, expira_em)
);

CREATE TABLE IF NOT EXISTS mentor_alunos (
  id CHAR(36) PRIMARY KEY,
  mentor_id CHAR(36) NOT NULL,
  aluno_id CHAR(36) NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  permite_cronograma_inteligente BOOLEAN NOT NULL DEFAULT FALSE,
  data_expiracao_plano DATE NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  desativado_em DATETIME NULL,
  desativado_por CHAR(36) NULL,
  aluno_ativo_chave CHAR(36) GENERATED ALWAYS AS (CASE WHEN ativo = TRUE THEN aluno_id ELSE NULL END) STORED,
  FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  UNIQUE KEY uk_aluno_mentor_ativo (aluno_ativo_chave),
  INDEX idx_mentor_alunos_mentor (mentor_id, ativo)
);

CREATE TABLE IF NOT EXISTS eventos_auditoria (
  id CHAR(36) PRIMARY KEY,
  executor_id CHAR(36) NOT NULL,
  aluno_id CHAR(36) NULL,
  acao VARCHAR(120) NOT NULL,
  entidade VARCHAR(120) NOT NULL,
  entidade_id CHAR(36) NULL,
  detalhes JSON NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_auditoria_aluno (aluno_id, criado_em),
  FOREIGN KEY (executor_id) REFERENCES usuarios(id)
);
