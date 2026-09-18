CREATE TABLE IF NOT EXISTS usuarios (
  id CHAR(36) PRIMARY KEY,
  nome VARCHAR(160) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  senha_hash VARCHAR(255) NOT NULL,
  papel ENUM('mestre','mentor','aluno') NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  ultimo_acesso_em DATETIME NULL,
  desativado_em DATETIME NULL,
  desativado_por CHAR(36) NULL,
  INDEX idx_usuarios_papel_ativo (papel, ativo)
);

CREATE TABLE IF NOT EXISTS sessoes_autenticacao (
  id CHAR(36) PRIMARY KEY,
  usuario_id CHAR(36) NOT NULL,
  token_hash CHAR(64) NOT NULL UNIQUE,
  csrf_hash CHAR(64) NOT NULL,
  expira_em DATETIME NOT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ultimo_uso_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  revogado_em DATETIME NULL,
  FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
  INDEX idx_sessoes_usuario (usuario_id, revogado_em, expira_em)
);

CREATE TABLE IF NOT EXISTS mentor_alunos (
  id CHAR(36) PRIMARY KEY,
  mentor_id CHAR(36) NOT NULL,
  aluno_id CHAR(36) NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  permite_cronograma_inteligente BOOLEAN NOT NULL DEFAULT FALSE,
  data_expiracao_plano DATE NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  desativado_em DATETIME NULL,
  desativado_por CHAR(36) NULL,
  aluno_ativo_chave CHAR(36) GENERATED ALWAYS AS (CASE WHEN ativo = TRUE THEN aluno_id ELSE NULL END) STORED,
  FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  UNIQUE KEY uk_aluno_mentor_ativo (aluno_ativo_chave),
  INDEX idx_mentor_alunos_mentor (mentor_id, ativo)
);

CREATE TABLE IF NOT EXISTS eventos_auditoria (
  id CHAR(36) PRIMARY KEY,
  executor_id CHAR(36) NOT NULL,
  aluno_id CHAR(36) NULL,
  acao VARCHAR(120) NOT NULL,
  entidade VARCHAR(120) NOT NULL,
  entidade_id CHAR(36) NULL,
  detalhes JSON NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_auditoria_aluno (aluno_id, criado_em),
  FOREIGN KEY (executor_id) REFERENCES usuarios(id)
);


CREATE TABLE IF NOT EXISTS concursos (
  id CHAR(36) PRIMARY KEY,
  nome VARCHAR(255) NOT NULL,
  banca VARCHAR(160) NOT NULL,
  cargo VARCHAR(255) NULL,
  salario DECIMAL(12,2) NULL,
  data_prova DATE NULL,
  pre_edital BOOLEAN NOT NULL DEFAULT FALSE,
  criado_por CHAR(36) NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  desativado_em DATETIME NULL,
  desativado_por CHAR(36) NULL,
  FOREIGN KEY (criado_por) REFERENCES usuarios(id)
);

CREATE TABLE IF NOT EXISTS aluno_concursos (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  concurso_id CHAR(36) NOT NULL,
  atribuido_por CHAR(36) NOT NULL,
  grupo ENUM('foco','mira','realizado') NOT NULL DEFAULT 'foco',
  ordem INT NOT NULL DEFAULT 0,
  contabiliza_estatisticas BOOLEAN NOT NULL DEFAULT TRUE,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  resultado ENUM('aguardando','aprovado','cadastro_reserva','reprovado','eliminado') NULL,
  classificacao INT NULL,
  nota_final DECIMAL(8,2) NULL,
  nomeado BOOLEAN NOT NULL DEFAULT FALSE,
  data_nomeacao DATE NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_aluno_concurso_ativo (aluno_id, concurso_id, ativo),
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  FOREIGN KEY (concurso_id) REFERENCES concursos(id),
  FOREIGN KEY (atribuido_por) REFERENCES usuarios(id)
);

CREATE TABLE IF NOT EXISTS editais (
  id CHAR(36) PRIMARY KEY,
  concurso_id CHAR(36) NOT NULL,
  nome VARCHAR(255) NOT NULL,
  versao VARCHAR(80) NULL,
  criado_por CHAR(36) NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (concurso_id) REFERENCES concursos(id),
  FOREIGN KEY (criado_por) REFERENCES usuarios(id)
);

CREATE TABLE IF NOT EXISTS edital_materias (
  id CHAR(36) PRIMARY KEY,
  edital_id CHAR(36) NOT NULL,
  nome VARCHAR(255) NOT NULL,
  ordem INT NOT NULL DEFAULT 0,
  peso DECIMAL(8,2) NULL,
  relevancia INT NULL,
  observacoes TEXT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (edital_id) REFERENCES editais(id),
  INDEX idx_edital_materias (edital_id, ordem, ativo)
);

CREATE TABLE IF NOT EXISTS edital_topicos (
  id CHAR(36) PRIMARY KEY,
  materia_id CHAR(36) NOT NULL,
  nome VARCHAR(500) NOT NULL,
  ordem INT NOT NULL DEFAULT 0,
  peso DECIMAL(8,2) NULL,
  relevancia INT NULL,
  observacoes TEXT NULL,
  materiais JSON NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (materia_id) REFERENCES edital_materias(id),
  INDEX idx_edital_topicos (materia_id, ordem, ativo)
);

CREATE TABLE IF NOT EXISTS edital_subtopicos (
  id CHAR(36) PRIMARY KEY,
  topico_id CHAR(36) NOT NULL,
  nome VARCHAR(500) NOT NULL,
  ordem INT NOT NULL DEFAULT 0,
  peso DECIMAL(8,2) NULL,
  relevancia INT NULL,
  observacoes TEXT NULL,
  materiais JSON NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (topico_id) REFERENCES edital_topicos(id),
  INDEX idx_edital_subtopicos (topico_id, ordem, ativo)
);

CREATE TABLE IF NOT EXISTS aluno_editais (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  edital_id CHAR(36) NOT NULL,
  atribuido_por CHAR(36) NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  atribuido_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_aluno_edital_ativo (aluno_id, edital_id, ativo),
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  FOREIGN KEY (edital_id) REFERENCES editais(id),
  FOREIGN KEY (atribuido_por) REFERENCES usuarios(id)
);

CREATE TABLE IF NOT EXISTS aluno_progresso_edital (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  tipo_item ENUM('topico','subtopico') NOT NULL,
  item_id CHAR(36) NOT NULL,
  estudado BOOLEAN NOT NULL DEFAULT FALSE,
  concluido_em DATETIME NULL,
  observacoes TEXT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_progresso_item (aluno_id, tipo_item, item_id),
  INDEX idx_progresso_aluno (aluno_id, estudado)
);

CREATE TABLE IF NOT EXISTS cronogramas (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  concurso_id CHAR(36) NULL,
  tipo ENUM('manual','agendado','ciclo_inteligente') NOT NULL,
  criado_por CHAR(36) NOT NULL,
  configuracao JSON NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  desativado_em DATETIME NULL,
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  FOREIGN KEY (concurso_id) REFERENCES concursos(id),
  FOREIGN KEY (criado_por) REFERENCES usuarios(id),
  INDEX idx_cronograma_aluno (aluno_id, ativo)
);

CREATE TABLE IF NOT EXISTS cronograma_itens (
  id CHAR(36) PRIMARY KEY,
  cronograma_id CHAR(36) NOT NULL,
  data_planejada DATE NULL,
  posicao_ciclo INT NULL,
  materia_id CHAR(36) NULL,
  topico_id CHAR(36) NULL,
  subtopico_id CHAR(36) NULL,
  duracao_minutos INT NULL,
  prioridade INT NULL,
  ordem INT NOT NULL DEFAULT 0,
  situacao ENUM('pendente','concluido','ignorado') NOT NULL DEFAULT 'pendente',
  concluido_em DATETIME NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (cronograma_id) REFERENCES cronogramas(id),
  INDEX idx_cronograma_itens (cronograma_id, data_planejada, ordem)
);

CREATE TABLE IF NOT EXISTS baralhos_cartoes (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NULL,
  concurso_id CHAR(36) NULL,
  edital_id CHAR(36) NULL,
  criado_por CHAR(36) NOT NULL,
  tipo_proprietario ENUM('mentor','aluno') NOT NULL,
  alcance ENUM('aluno','edital','global') NOT NULL,
  nome VARCHAR(255) NOT NULL,
  descricao TEXT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  FOREIGN KEY (concurso_id) REFERENCES concursos(id),
  FOREIGN KEY (edital_id) REFERENCES editais(id),
  FOREIGN KEY (criado_por) REFERENCES usuarios(id),
  INDEX idx_baralhos_alcance (alcance, aluno_id, edital_id, ativo)
);

CREATE TABLE IF NOT EXISTS cartoes_estudo (
  id CHAR(36) PRIMARY KEY,
  baralho_id CHAR(36) NOT NULL,
  tipo ENUM('basico','lacuna','multipla_escolha') NOT NULL,
  frente TEXT NOT NULL,
  verso TEXT NULL,
  alternativas JSON NULL,
  resposta_correta VARCHAR(255) NULL,
  topico_id CHAR(36) NULL,
  subtopico_id CHAR(36) NULL,
  criado_por CHAR(36) NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (baralho_id) REFERENCES baralhos_cartoes(id),
  FOREIGN KEY (criado_por) REFERENCES usuarios(id),
  INDEX idx_cartoes_baralho (baralho_id, ativo)
);

CREATE TABLE IF NOT EXISTS revisoes_cartoes (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  cartao_id CHAR(36) NOT NULL,
  qualidade INT NOT NULL,
  repeticoes INT NOT NULL DEFAULT 0,
  intervalo_dias INT NOT NULL DEFAULT 0,
  facilidade DECIMAL(5,2) NOT NULL DEFAULT 2.50,
  proxima_revisao DATE NOT NULL,
  revisado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  FOREIGN KEY (cartao_id) REFERENCES cartoes_estudo(id),
  INDEX idx_revisoes_cartoes (aluno_id, proxima_revisao)
);

CREATE TABLE IF NOT EXISTS sessoes_estudo (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  concurso_id CHAR(36) NULL,
  materia_id CHAR(36) NULL,
  topico_id CHAR(36) NULL,
  subtopico_id CHAR(36) NULL,
  segundos INT NOT NULL,
  modo VARCHAR(50) NULL,
  observacoes TEXT NULL,
  estudado_em DATETIME NOT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  INDEX idx_sessoes_estudo (aluno_id, estudado_em)
);

CREATE TABLE IF NOT EXISTS registros_questoes (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  concurso_id CHAR(36) NULL,
  materia_id CHAR(36) NULL,
  topico_id CHAR(36) NULL,
  subtopico_id CHAR(36) NULL,
  resolvidas INT NOT NULL,
  acertos INT NOT NULL,
  erros INT NOT NULL,
  registrado_em DATETIME NOT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  INDEX idx_registros_questoes (aluno_id, registrado_em)
);

CREATE TABLE IF NOT EXISTS revisoes_programadas (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  materia_id CHAR(36) NULL,
  topico_id CHAR(36) NULL,
  subtopico_id CHAR(36) NULL,
  ciclo_atual INT NOT NULL DEFAULT 0,
  proxima_data DATE NULL,
  percentual_anterior DECIMAL(5,2) NULL,
  concluida BOOLEAN NOT NULL DEFAULT FALSE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  INDEX idx_revisoes_programadas (aluno_id, concluida, proxima_data)
);

CREATE TABLE IF NOT EXISTS configuracoes_prova (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  concurso_id CHAR(36) NOT NULL,
  configuracao JSON NOT NULL,
  criado_por CHAR(36) NOT NULL,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_configuracao_prova (aluno_id, concurso_id),
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  FOREIGN KEY (concurso_id) REFERENCES concursos(id),
  FOREIGN KEY (criado_por) REFERENCES usuarios(id)
);

CREATE TABLE IF NOT EXISTS simulados (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  concurso_id CHAR(36) NULL,
  nome VARCHAR(255) NOT NULL,
  tipo ENUM('realizado','pendente') NOT NULL,
  realizado_em DATE NULL,
  link VARCHAR(2000) NULL,
  observacoes TEXT NULL,
  percentual DECIMAL(5,2) NULL,
  tempo_minutos INT NULL,
  questoes_feitas INT NULL,
  resultado JSON NULL,
  configuracao_usada JSON NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  INDEX idx_simulados_aluno (aluno_id, realizado_em, ativo)
);

CREATE TABLE IF NOT EXISTS resultados_materias_simulado (
  id CHAR(36) PRIMARY KEY,
  simulado_id CHAR(36) NOT NULL,
  materia_id CHAR(36) NULL,
  nome VARCHAR(255) NOT NULL,
  questoes INT NULL,
  acertos INT NULL,
  erros INT NULL,
  brancos INT NULL,
  percentual DECIMAL(5,2) NULL,
  pontos DECIMAL(10,2) NULL,
  FOREIGN KEY (simulado_id) REFERENCES simulados(id),
  INDEX idx_resultados_simulado (simulado_id)
);


ALTER TABLE cronogramas ADD COLUMN versao INT NOT NULL DEFAULT 1 AFTER criado_por;
ALTER TABLE cronogramas ADD COLUMN estado ENUM('rascunho','ativo','concluido') NOT NULL DEFAULT 'ativo' AFTER versao;

CREATE TABLE IF NOT EXISTS cronograma_resumos_diarios (
  id CHAR(36) PRIMARY KEY,
  cronograma_id CHAR(36) NOT NULL,
  aluno_id CHAR(36) NOT NULL,
  data_referencia DATE NOT NULL,
  conteudo JSON NOT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_resumo_diario (cronograma_id, data_referencia),
  FOREIGN KEY (cronograma_id) REFERENCES cronogramas(id),
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id),
  INDEX idx_resumos_aluno_data (aluno_id, data_referencia)
);


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


ALTER TABLE sessoes_estudo
  ADD COLUMN metricas JSON NULL AFTER observacoes,
  ADD COLUMN origem VARCHAR(50) NULL AFTER metricas,
  ADD COLUMN ativo BOOLEAN NOT NULL DEFAULT TRUE AFTER origem,
  ADD COLUMN atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER criado_em;

ALTER TABLE registros_questoes
  ADD COLUMN origem VARCHAR(50) NULL AFTER erros,
  ADD COLUMN ativo BOOLEAN NOT NULL DEFAULT TRUE AFTER origem,
  ADD COLUMN atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER criado_em;

ALTER TABLE revisoes_programadas
  ADD COLUMN concurso_id CHAR(36) NULL AFTER aluno_id,
  ADD COLUMN observacoes TEXT NULL AFTER percentual_anterior,
  ADD COLUMN concluida_em DATETIME NULL AFTER concluida,
  ADD COLUMN ativo BOOLEAN NOT NULL DEFAULT TRUE AFTER concluida_em,
  ADD INDEX idx_revisoes_concurso (aluno_id, concurso_id, ativo);

ALTER TABLE resultados_materias_simulado
  ADD COLUMN ativo BOOLEAN NOT NULL DEFAULT TRUE,
  ADD COLUMN desativado_em DATETIME NULL;


CREATE INDEX idx_sessoes_metricas ON sessoes_estudo (aluno_id, concurso_id, ativo, estudado_em, materia_id);
CREATE INDEX idx_questoes_metricas ON registros_questoes (aluno_id, concurso_id, ativo, registrado_em, materia_id);
CREATE INDEX idx_simulados_metricas ON simulados (aluno_id, concurso_id, ativo, tipo, realizado_em);
CREATE INDEX idx_resultados_materia_metricas ON resultados_materias_simulado (materia_id, ativo, simulado_id);
CREATE INDEX idx_progresso_item_metricas ON aluno_progresso_edital (aluno_id, tipo_item, item_id, estudado);


ALTER TABLE concursos
  ADD COLUMN logotipo LONGTEXT NULL AFTER cargo;


ALTER TABLE edital_materias ADD COLUMN materiais JSON NULL AFTER observacoes;

CREATE TABLE IF NOT EXISTS aluno_progresso_materiais (
  id CHAR(36) PRIMARY KEY,
  aluno_id CHAR(36) NOT NULL,
  tipo_item ENUM('materia','topico','subtopico') NOT NULL,
  item_id CHAR(36) NOT NULL,
  material_id VARCHAR(160) NOT NULL,
  concluido BOOLEAN NOT NULL DEFAULT FALSE,
  concluido_em DATETIME NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_progresso_material (aluno_id, tipo_item, item_id, material_id),
  INDEX idx_progresso_materiais_aluno (aluno_id, concluido),
  FOREIGN KEY (aluno_id) REFERENCES usuarios(id)
);


ALTER TABLE revisoes_cartoes
  ADD COLUMN concurso_id CHAR(36) NULL AFTER aluno_id,
  ADD FOREIGN KEY (concurso_id) REFERENCES concursos(id);


CREATE TABLE IF NOT EXISTS materiais_apoio (
  id CHAR(36) PRIMARY KEY,
  mentor_id CHAR(36) NOT NULL,
  titulo VARCHAR(255) NOT NULL,
  descricao TEXT NULL,
  tipo ENUM('arquivo','youtube','texto','link') NOT NULL,
  url VARCHAR(2048) NULL,
  texto TEXT NULL,
  escopo ENUM('global','edital') NOT NULL,
  edital_id CHAR(36) NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
  FOREIGN KEY (edital_id) REFERENCES editais(id),
  INDEX idx_materiais_apoio_mentor (mentor_id, ativo),
  INDEX idx_materiais_apoio_edital (edital_id, ativo)
);


ALTER TABLE materiais_apoio
  ADD COLUMN arquivo_nome VARCHAR(255) NULL,
  ADD COLUMN arquivo_caminho VARCHAR(1024) NULL,
  ADD COLUMN arquivo_mime VARCHAR(255) NULL;


-- Adiciona coluna de prazos de revisão configurados pelo mentor no concurso
ALTER TABLE concursos
  ADD COLUMN prazos_revisao VARCHAR(255) NOT NULL DEFAULT '1,7,30';


-- Adiciona coluna de telefone no cadastro de usuários
ALTER TABLE usuarios
  ADD COLUMN telefone VARCHAR(40) NULL AFTER email;


ALTER TABLE materiais_apoio ADD COLUMN pasta VARCHAR(255) DEFAULT '';


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


-- Corrige colunas da tabela mentor_whitelabel que podem ter sido criadas com VARCHAR(255)
-- logo_url pode ser uma Data URL (base64) de imagem, precisando de espaço maior
ALTER TABLE mentor_whitelabel
  MODIFY COLUMN logo_url TEXT NULL,
  MODIFY COLUMN mensagem_boas_vindas TEXT NULL;


CREATE TABLE IF NOT EXISTS cursos (
  id CHAR(36) PRIMARY KEY,
  mentor_id CHAR(36) NOT NULL,
  titulo VARCHAR(180) NOT NULL,
  descricao TEXT NOT NULL,
  categoria VARCHAR(80) NOT NULL DEFAULT 'Geral',
  capa_url TEXT NOT NULL,
  escopo ENUM('global','alunos','concursos') NOT NULL,
  destinatarios JSON NOT NULL,
  aulas JSON NOT NULL,
  ativo BOOLEAN NOT NULL DEFAULT TRUE,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
  INDEX idx_cursos_mentor (mentor_id, ativo)
);


CREATE TABLE IF NOT EXISTS cursos_pdfs (
  id CHAR(36) PRIMARY KEY,
  mentor_id CHAR(36) NOT NULL,
  nome VARCHAR(180) NOT NULL,
  conteudo MEDIUMBLOB NOT NULL,
  criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
  INDEX idx_cursos_pdfs_mentor (mentor_id)
);


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


CREATE TABLE IF NOT EXISTS cursos_capas (
 id CHAR(36) PRIMARY KEY,
 mentor_id CHAR(36) NOT NULL,
 tipo VARCHAR(30) NOT NULL,
 conteudo MEDIUMBLOB NOT NULL,
 criado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 FOREIGN KEY (mentor_id) REFERENCES usuarios(id),
 INDEX idx_cursos_capas_mentor (mentor_id)
);


ALTER TABLE cursos ADD COLUMN modo_exibicao VARCHAR(20) NOT NULL DEFAULT 'curso';


ALTER TABLE mentor_whitelabel
  ADD COLUMN banner_url TEXT NULL AFTER logo_url;
