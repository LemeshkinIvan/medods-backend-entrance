package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (
			title, description, status, 
			created_at, updated_at, scheduled_at, 
			type_of_repetition, periodicity, dates
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING 
			id, title, description, 
			status, created_at, updated_at, 
			scheduled_at, type_of_repetition, periodicity, dates
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
		task.ScheduledAt,
		task.TypeOfRepetition,
		task.Periodicity,
		task.CustomDates,
	)
	created, err := scanTask(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT 
			id, title, description, 
			status, created_at, updated_at, 
			scheduled_at, type_of_repetition, periodicity, dates
		FROM tasks
		WHERE id = $1
		ORDER BY id DESC
	`

	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			updated_at = $4,
			scheduled_at = $5,
			type_of_repetition = $6,
			periodicity = $7,
			dates = $8
		WHERE id = $9
		RETURNING 
			id, title, description, status, 
			created_at,
			updated_at, 
			scheduled_at, 
			type_of_repetition, 
			periodicity, 
			dates
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
		task.ScheduledAt,
		task.TypeOfRepetition,
		task.Periodicity,
		task.CustomDates,
		task.ID,
	)
	updated, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT 
			id, 
			title, 
			description, 
			status, 
			created_at, 
			updated_at, 
			scheduled_at, 
			type_of_repetition, 
			periodicity,
			dates
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) ListByDate(ctx context.Context, date time.Time) ([]taskdomain.Task, error) {
	const query = `
		SELECT *
		FROM tasks
		WHERE created_at <= $1
		AND scheduled_at >= $1;
	`

	rows, err := r.pool.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task       taskdomain.Task
		status     string
		repetition string
	)

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.ScheduledAt,
		&repetition,
		&task.Periodicity,
		&task.CustomDates,
	); err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)
	task.TypeOfRepetition = taskdomain.Repetition(repetition)

	return &task, nil
}
