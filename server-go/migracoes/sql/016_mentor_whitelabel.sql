CREATE TABLE IF NOT EXISTS mentor_whitelabel (
  mentor_id CHAR(36) NOT NULL PRIMARY KEY,
  nome_plataforma VARCHAR(255) NOT NULL DEFAULT 'Chega Junto Concurseiro',
  logo_url TEXT NULL,
  cor_primaria VARCHAR(50) NOT NULL DEFAULT '#4f8ef7',
  cor_secundaria VARCHAR(50) NOT NULL DEFAULT '#7c5cfc',
  mensagem_boas_vindas VARCHAR(255) NULL DEFAULT 'Análise completa da sua preparação',
  atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (mentor_id) REFERENCES usuarios(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
