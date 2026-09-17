-- Migration 009: Add concurso_id to revisoes_cartoes if not exists and backfill NULL records
SET @dbname = DATABASE();
SET @tablename = 'revisoes_cartoes';
SET @columnname = 'concurso_id';
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = @dbname
    AND TABLE_NAME = @tablename
    AND COLUMN_NAME = @columnname
  ) > 0,
  'SELECT 1',
  'ALTER TABLE revisoes_cartoes ADD COLUMN concurso_id CHAR(36) NULL AFTER aluno_id, ADD FOREIGN KEY (concurso_id) REFERENCES concursos(id);'
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- Backfill NULL concurso_id in revisoes_cartoes using active student contest
UPDATE revisoes_cartoes r
JOIN aluno_concursos ac ON ac.aluno_id = r.aluno_id AND ac.ativo = TRUE
SET r.concurso_id = ac.concurso_id
WHERE r.concurso_id IS NULL;
