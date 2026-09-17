ALTER TABLE edital_materias ADD COLUMN materiais JSON NULL AFTER observacoes;

CREATE TABLE IF NOT EXISTS aluno_progresso_materiais (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  tipo_item ENUM('materia','topico','subtopico') NOT NULL,
  item_id CHAR(36) NOT NULL,
  material_id VARCHAR(160) NOT NULL,
  concluido BOOLEAN NOT NULL DEFAULT FALSE,
  concluido_em DATETIME NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_progresso_material (aluno_id, tipo_item, item_id, material_id),
  INDEX idx_progresso_materiais_aluno (aluno_id, concluido),
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id)
);
