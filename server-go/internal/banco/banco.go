package banco

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"chega-junto-concurseiro-web/internal/configuracao"

	_ "github.com/go-sql-driver/mysql"
)

func Abrir(ctx context.Context, cfg configuracao.Configuracao) (*sql.DB, error) {
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(cfg.DBNome) {
		return nil, fmt.Errorf("DB_NOME contém caracteres inválidos")
	}
	if cfg.DBCriarBanco {
		servidorDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true&loc=UTC", cfg.DBUsuario, cfg.DBSenha, cfg.DBHost, cfg.DBPorta)
		servidor, err := sql.Open("mysql", servidorDSN)
		if err != nil {
			return nil, err
		}
		if err := servidor.PingContext(ctx); err != nil {
			servidor.Close()
			return nil, err
		}
		if _, err := servidor.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS `"+cfg.DBNome+"` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
			servidor.Close()
			return nil, err
		}
		servidor.Close()
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=UTC&multiStatements=true", cfg.DBUsuario, cfg.DBSenha, cfg.DBHost, cfg.DBPorta, cfg.DBNome)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
