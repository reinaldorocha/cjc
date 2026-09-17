CREATE TABLE IF NOT EXISTS banco_questoes (
  id CHAR(36) PRIMARY KEY,
  disciplina VARCHAR(150) NOT NULL,
  assunto VARCHAR(150) NOT NULL,
  tipo VARCHAR(30) NOT NULL DEFAULT 'multipla_escolha',
  enunciado TEXT NOT NULL,
  alternativas JSON NULL,
  resposta_correta VARCHAR(255) NOT NULL,
  explicacao TEXT NULL,
  alcance VARCHAR(20) NOT NULL DEFAULT 'global',
  concurso_id CHAR(36) NULL,
  criado_por CHAR(36) NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_banco_questoes_disciplina (disciplina),
  INDEX idx_banco_questoes_assunto (assunto),
  INDEX idx_banco_questoes_alcance (alcance, concurso_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS respostas_banco_questoes (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  questao_id CHAR(36) NOT NULL,
  concurso_id CHAR(36) NULL,
  resposta_aluno VARCHAR(255) NOT NULL,
  correto BOOLEAN NOT NULL,
  respondido_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_respostas_bq_aluno (aluno_id, questao_id),
  INDEX idx_respostas_bq_concurso (aluno_id, concurso_id),
  CONSTRAINT fk_respostas_bq_aluno FOREIGN KEY (aluno_id) REFERENCES usuarios(id) ON DELETE CASCADE,
  CONSTRAINT fk_respostas_bq_questao FOREIGN KEY (questao_id) REFERENCES banco_questoes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
