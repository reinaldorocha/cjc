CREATE TABLE IF NOT EXISTS cursos (
  id CHAR(36) PRIMARY KEY,
  mentor_id CHAR(36) NOT NULL,
  titulo VARCHAR(180) NOT NULL,
  descricao TEXT NOT NULL,
  categoria VARCHAR(80) NOT NULL DEFAULT 'Geral',
  capa_url TEXT NOT NULL,
  escopo ENUM('global','alunos','concursos') NOT NULL,
  destinatarios JSON NOT NULL,
  aulas JSON NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
  INDEX idx_cursos_mentor (mentor_id, ativo)
);
