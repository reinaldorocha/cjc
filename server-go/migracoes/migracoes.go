// Package migracoes aplica, em ordem, as alterações versionadas do banco.
package migracoes

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed sql/*.sql
var arquivos embed.FS

const migracaoInicial = "001_schema_inicial.sql"

func Aplicar(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS historico_migracoes (arquivo VARCHAR(255) PRIMARY KEY, aplicado_em DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		return err
	}
	entradas, err := arquivos.ReadDir("sql")
	if err != nil {
		return err
	}
	nomes := make([]string, 0, len(entradas))
	for _, e := range entradas {
		if !e.IsDir() {
			nomes = append(nomes, e.Name())
		}
	}
	sort.Strings(nomes)
	for _, nome := range nomes {
		if nome == migracaoInicial {
			var migracoesLegadas int
			if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM historico_migracoes WHERE arquivo <> ?`, nome).Scan(&migracoesLegadas); err != nil {
				return err
			}
			if migracoesLegadas > 0 {
				if _, err := db.ExecContext(ctx, `INSERT IGNORE INTO historico_migracoes (arquivo) VALUES (?)`, nome); err != nil {
					return err
				}
				continue
			}
		}
		var existe int
		if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM historico_migracoes WHERE arquivo = ?`, nome).Scan(&existe); err != nil {
			return err
		}
		if existe > 0 {
			continue
		}
		conteudo, err := arquivos.ReadFile("sql/" + nome)
		if err != nil {
			return err
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, comando := range strings.Split(string(conteudo), ";") {
			if strings.TrimSpace(comando) == "" {
				continue
			}
			if _, err = tx.ExecContext(ctx, comando); err != nil {
				tx.Rollback()
				return fmt.Errorf("migração %s: %w", nome, err)
			}
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO historico_migracoes (arquivo) VALUES (?)`, nome); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
