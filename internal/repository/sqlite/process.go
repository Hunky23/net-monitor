package sqlite

import (
	"changeme/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

type ProcessRepository struct {
	db *sql.DB
}

func NewProcessRepository(db *sql.DB) *ProcessRepository {
	return &ProcessRepository{db: db}
}

func (r ProcessRepository) Save(ctx context.Context, process *domain.Process) error {
	if process.ID == 0 {
		return r.insert(ctx, process)
	}

	return r.update(ctx, process)
}

func (r ProcessRepository) FindByFilters(ctx context.Context, filters any) ([]domain.Process, error) {
	return nil, nil
}

func (r ProcessRepository) insert(ctx context.Context, process *domain.Process) error {
	const query = `
		INSERT INTO processes (pid, process_name, display_name, path, status, type, icon, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		process.PID,
		process.ProcessName,
		process.DisplayName,
		process.Path,
		process.Status,
		process.Type,
		process.Icon,
		process.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("репозиторий: создание процесса: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("репозиторий: создание процесса: получение id: %w", err)
	}

	process.ID = id

	return nil
}

func (r ProcessRepository) update(ctx context.Context, process *domain.Process) error {
	const query = `
		UPDATE processes
		SET pid = ?, process_name = ?, display_name = ?, path = ?, status = ?, type = ?, icon = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		process.PID,
		process.ProcessName,
		process.DisplayName,
		process.Path,
		process.Status,
		process.Type,
		process.Icon,
		process.ID,
	)
	if err != nil {
		return fmt.Errorf("репозиторий: обновление процесса: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("репозиторий: обновление процесса: чтение числа обновлённых строк: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("репозиторий: обновление процесса: не найден процесс с id=%d", process.ID)
	}

	return nil
}
