ALTER TABLE cronogramas ADD COLUMN versao INT NOT NULL DEFAULT 1 AFTER criado_por;
ALTER TABLE cronogramas ADD COLUMN estado ENUM('rascunho','ativo','concluido') NOT NULL DEFAULT 'ativo' AFTER versao;

CREATE TABLE IF NOT EXISTS cronograma_resumos_diarios (
  id CHAR(36) PRIMARY KEY,
  cronograma_id CHAR(36) NOT NULL,
  aluno_id CHAR(36) NOT NULL,
  data_referencia DATE NOT NULL,
  conteudo JSON NOT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_resumo_diario (cronograma_id, data_referencia),
  FOREIGN KEY (cronograma_id) REFERENCES cronogramas(id),
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  INDEX idx_resumos_aluno_data (aluno_id, data_referencia)
);
