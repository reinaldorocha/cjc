ALTER TABLE baralhos_cartoes
  ADD COLUMN baralho_pai_id CHAR(36) NULL AFTER edital_id,
  ADD COLUMN materia_id CHAR(36) NULL AFTER baralho_pai_id,
  ADD COLUMN topico_id CHAR(36) NULL AFTER materia_id,
  ADD COLUMN subtopico_id CHAR(36) NULL AFTER topico_id,
  ADD COLUMN icone VARCHAR(20) NULL AFTER descricao,
  ADD COLUMN ordem INT NOT NULL DEFAULT 0 AFTER icone,
  ADD INDEX idx_baralhos_hierarquia (baralho_pai_id, ordem),
  ADD CONSTRAINT fk_baralhos_pai FOREIGN KEY (baralho_pai_id) REFERENCES baralhos_cartoes(id);

ALTER TABLE cartoes_estudo
  MODIFY COLUMN tipo ENUM('basico','lacuna','multipla_escolha','certo_errado') NOT NULL,
  ADD COLUMN dica TEXT NULL AFTER verso,
  ADD COLUMN explicacao TEXT NULL AFTER resposta_correta,
  ADD COLUMN explicacoes_alternativas JSON NULL AFTER explicacao,
  ADD COLUMN etiquetas JSON NULL AFTER explicacoes_alternativas;
