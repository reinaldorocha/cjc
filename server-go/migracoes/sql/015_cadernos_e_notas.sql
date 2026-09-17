CREATE TABLE IF NOT EXISTS cadernos (
    id VARCHAR(255) PRIMARY KEY,
    aluno_id VARCHAR(255) NOT NULL,
    titulo VARCHAR(255) NOT NULL,
    pasta VARCHAR(255) DEFAULT 'Geral',
    conteudo LONGTEXT,
    edital_id VARCHAR(255) NULL,
    materia_id VARCHAR(255) NULL,
    topico_id VARCHAR(255) NULL,
    cor VARCHAR(50) DEFAULT '#4f8ef7',
    ativo BOOLEAN DEFAULT TRUE,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_cadernos_aluno (aluno_id, ativo)
);
