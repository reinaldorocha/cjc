-- Adiciona coluna de prazos de revisão configurados pelo mentor no concurso
ALTER TABLE concursos
  ADD COLUMN prazos_revisao VARCHAR(255) NOT NULL DEFAULT '1,7,30';
