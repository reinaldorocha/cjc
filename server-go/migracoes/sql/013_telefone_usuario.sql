-- Adiciona coluna de telefone no cadastro de usuários
ALTER TABLE usuarios
  ADD COLUMN telefone VARCHAR(40) NULL AFTER email;
