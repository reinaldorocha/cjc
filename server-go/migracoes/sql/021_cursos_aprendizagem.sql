ALTER TABLE cursos ADD COLUMN publicado JSON NULL, ADD COLUMN revisao INT NOT NULL DEFAULT 1, ADD COLUMN revisao_publicada INT NOT NULL DEFAULT 0;
UPDATE cursos SET publicado=JSON_OBJECT('id',id,'titulo',titulo,'descricao',descricao,'categoria',categoria,'capaUrl',capa_url,'escopo',escopo,'destinatarios',destinatarios,'aulas',aulas), revisao_publicada=1;
CREATE TABLE IF NOT EXISTS cursos_progresso (
 aluno_id CHAR(36) NOT NULL, curso_id CHAR(36) NOT NULL, aula_id CHAR(36) NOT NULL,
 video_id VARCHAR(11) NOT NULL, posicao DOUBLE NOT NULL DEFAULT 0, duracao DOUBLE NOT NULL DEFAULT 0,
 concluida BOOLEAN NOT NULL DEFAULT FALSE, atualizado_em DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 PRIMARY KEY (aluno_id,curso_id,aula_id),
 FOREIGN KEY (aluno_id) REFERENCES usuarios(id), FOREIGN KEY (curso_id) REFERENCES cursos(id),
 INDEX idx_cursos_progresso (curso_id,aluno_id)
);
CREATE TABLE IF NOT EXISTS cursos_respostas (
 aluno_id CHAR(36) NOT NULL, curso_id CHAR(36) NOT NULL, aula_id CHAR(36) NOT NULL, questao_id CHAR(36) NOT NULL,
 resposta VARCHAR(255) NOT NULL, correta BOOLEAN NOT NULL, tentativas INT NOT NULL DEFAULT 1,
 atualizado_em DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 PRIMARY KEY (aluno_id,curso_id,aula_id,questao_id),
 FOREIGN KEY (aluno_id) REFERENCES usuarios(id), FOREIGN KEY (curso_id) REFERENCES cursos(id),
 FOREIGN KEY (questao_id) REFERENCES banco_questoes(id)
);
