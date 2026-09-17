CREATE TABLE IF NOT EXISTS cursos_pdfs (
  id CHAR(36) PRIMARY KEY,
  mentor_id CHAR(36) NOT NULL,
  nome VARCHAR(180) NOT NULL,
  conteudo MEDIUMBLOB NOT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
  INDEX idx_cursos_pdfs_mentor (mentor_id)
);
