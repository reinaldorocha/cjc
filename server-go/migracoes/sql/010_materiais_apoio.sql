CREATE TABLE IF NOT EXISTS materiais_apoio (
  id CHAR(36) PRIMARY KEY,
  mentor_id CHAR(36) NOT NULL,
  titulo VARCHAR(255) NOT NULL,
  descricao TEXT NULL,
  tipo ENUM('arquivo','youtube','texto','link') NOT NULL,
  url VARCHAR(2048) NULL,
  texto TEXT NULL,
  escopo ENUM('global','edital') NOT NULL,
  edital_id CHAR(36) NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
  FOREIGN KEY (edital_id) REFERENCES editais(id),
  INDEX idx_materiais_apoio_mentor (mentor_id, ativo),
  INDEX idx_materiais_apoio_edital (edital_id, ativo)
);
