CREATE TABLE IF NOT EXISTS cursos_capas (
 id CHAR(36) PRIMARY KEY,
 mentor_id CHAR(36) NOT NULL,
 tipo VARCHAR(30) NOT NULL,
 conteudo MEDIUMBLOB NOT NULL,
 criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
 INDEX idx_cursos_capas_mentor (mentor_id)
);
