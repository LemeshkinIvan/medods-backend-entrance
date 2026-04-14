package task

import (
	"context"
	"fmt"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	formatDate, err := parseDate(input.ScheduledAt)
	if err != nil {
		return nil, err
	}

	now := s.now()
	createdAt := now
	updatedAt := now

	if err := validateCustomDates(normalized.CustomDates, createdAt, formatDate); err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		Periodicity:      input.Periodicity,
		TypeOfRepetition: input.TypeOfRepetition,
		ScheduledAt:      formatDate,
		CustomDates:      normalized.CustomDates,
	}

	model.CreatedAt = createdAt
	model.UpdatedAt = updatedAt
	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	formatDate, err := parseDate(input.ScheduledAt)
	if err != nil {
		return nil, err
	}

	createdAt, err := parseDate(normalized.CreatedAt)
	if err != nil {
		return nil, err
	}

	now := s.now()

	if err := validateCustomDates(normalized.CustomDates, createdAt, formatDate); err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:               id,
		Title:            normalized.Title,
		Description:      normalized.Description,
		Status:           normalized.Status,
		Periodicity:      input.Periodicity,
		TypeOfRepetition: input.TypeOfRepetition,
		ScheduledAt:      formatDate,
		UpdatedAt:        now,
		CustomDates:      normalized.CustomDates,
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, date string) ([]taskdomain.Task, error) {
	formatDate, err := parseParamDate(date)
	if err != nil {
		return nil, err
	}

	// если параметр не пустой, то отдаем активные задачи
	if !formatDate.IsZero() {
		tasks, err := s.repo.ListByDate(ctx, formatDate)
		if err != nil {
			return nil, err
		}

		filtered := []taskdomain.Task{}
		for _, task := range tasks {
			if isActive(task, formatDate) {
				filtered = append(filtered, task)
			}
		}

		return filtered, nil
	}

	return s.repo.List(ctx)
}
