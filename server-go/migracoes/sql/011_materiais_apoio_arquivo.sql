ALTER TABLE materiais_apoio
  ADD COLUMN arquivo_nome VARCHAR(255) NULL,
  ADD COLUMN arquivo_caminho VARCHAR(1024) NULL,
  ADD COLUMN arquivo_mime VARCHAR(255) NULL;
