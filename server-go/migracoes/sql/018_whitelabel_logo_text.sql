-- Corrige colunas da tabela mentor_whitelabel que podem ter sido criadas com VARCHAR(255)
-- logo_url pode ser uma Data URL (base64) de imagem, precisando de espaço maior
ALTER TABLE mentor_whitelabel
  MODIFY COLUMN logo_url TEXT NULL,
  MODIFY COLUMN mensagem_boas_vindas TEXT NULL;
