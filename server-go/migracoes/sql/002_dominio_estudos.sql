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
