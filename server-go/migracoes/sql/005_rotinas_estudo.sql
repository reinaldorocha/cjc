ALTER TABLE sessoes_estudo
  ADD COLUMN metricas JSON NULL AFTER observacoes,
  ADD COLUMN origem VARCHAR(50) NULL AFTER metricas,
  ADD COLUMN ativo BOOLEAN NOT NULL DEFAULT TRUE AFTER origem,
  ADD COLUMN atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER criado_em;

ALTER TABLE registros_questoes
  ADD COLUMN origem VARCHAR(50) NULL AFTER erros,
  ADD COLUMN ativo BOOLEAN NOT NULL DEFAULT TRUE AFTER origem,
  ADD COLUMN atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER criado_em;

ALTER TABLE revisoes_programadas
  ADD COLUMN concurso_id CHAR(36) NULL AFTER aluno_id,
  ADD COLUMN observacoes TEXT NULL AFTER percentual_anterior,
  ADD COLUMN concluida_em DATETIME NULL AFTER concluida,
  ADD COLUMN ativo BOOLEAN NOT NULL DEFAULT TRUE AFTER concluida_em,
  ADD INDEX idx_revisoes_concurso (aluno_id, concurso_id, ativo);

ALTER TABLE resultados_materias_simulado
  ADD COLUMN ativo BOOLEAN NOT NULL DEFAULT TRUE,
  ADD COLUMN desativado_em DATETIME NULL;
