CREATE INDEX idx_sessoes_metricas ON sessoes_estudo (aluno_id, concurso_id, ativo, estudado_em, materia_id);
CREATE INDEX idx_questoes_metricas ON registros_questoes (aluno_id, concurso_id, ativo, registrado_em, materia_id);
CREATE INDEX idx_simulados_metricas ON simulados (aluno_id, concurso_id, ativo, tipo, realizado_em);
CREATE INDEX idx_resultados_materia_metricas ON resultados_materias_simulado (materia_id, ativo, simulado_id);
CREATE INDEX idx_progresso_item_metricas ON aluno_progresso_edital (aluno_id, tipo_item, item_id, estudado);
