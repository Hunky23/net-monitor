package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(ctx context.Context, db *sql.DB) error {
	dirs, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("база данных: директория миграций: %w", err)
	}

	for _, d := range dirs {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".sql") {
			continue
		}

		script, err := fs.ReadFile(migrationsFS, "migrations/"+d.Name())
		if err != nil {
			return fmt.Errorf("база данных: чтение миграции %s: %w", d.Name(), err)
		}

		if _, err := db.ExecContext(ctx, string(script)); err != nil {
			return fmt.Errorf("база данных: применение миграции %s: %w", d.Name(), err)
		}
	}

	return nil
}
