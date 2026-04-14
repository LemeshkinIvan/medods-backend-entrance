package task

import (
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func validateCustomDates(dates []string, createdAt time.Time, deadline time.Time) error {
	for _, i := range dates {
		date, err := parseDate(i)
		if err != nil {
			return err
		}

		if date.Before(createdAt) {
			return ErrCustomDateBefore
		}

		if date.After(deadline) {
			return ErrCustomDateAfter
		}
	}

	return nil
}

func parseParamDate(input string) (time.Time, error) {
	if input == "" {
		return time.Time{}, nil
	}

	// формат 18.11.2026
	if t, err := time.Parse("02.01.2006", input); err == nil {
		return t, nil
	}

	// формат 18.11.2026 15:04
	// лучше без точек, а через - в url
	if t, err := time.Parse("02.01.2006 15:04", input); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("invalid date format: use DD.MM.YYYY or DD.MM.YYYY HH:MM")
}

func parseDate(input string) (time.Time, error) {
	if input == "" {
		return time.Time{}, fmt.Errorf("%w: date is empty", ErrInvalidInput)
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	t, err := time.ParseInLocation("02.01.2006 15:04", input, loc)
	if err != nil {
		t, err = time.ParseInLocation("02.01.2006", input, loc)
		if err != nil {
			return time.Time{}, fmt.Errorf("%w: invalid date format", ErrInvalidInput)
		}
	}

	return t, nil
}

// здесь можно было бы сделать generic-валидацию, но для текущего кейса
// проще и читабельнее оставить вот так. мало ли логика изменится в одном из методов
func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.TypeOfRepetition == "" {
		return CreateInput{}, fmt.Errorf("%w: empty type of repetition", ErrInvalidInput)
	}

	if !input.TypeOfRepetition.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid type of repetition", ErrInvalidInput)
	}

	if err := validateType(input.TypeOfRepetition, input.CustomDates); err != nil {
		return CreateInput{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	if input.Periodicity < 0 && input.Periodicity > 31 {
		return CreateInput{}, fmt.Errorf("%w: invalid periodicity", ErrInvalidInput)
	}

	if input.ScheduledAt == "" {
		return CreateInput{}, fmt.Errorf("%w: invalid scheduled_at", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.TypeOfRepetition == "" {
		return UpdateInput{}, fmt.Errorf("%w: empty type of repetition", ErrInvalidInput)
	}

	if !input.TypeOfRepetition.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid type of repetition", ErrInvalidInput)
	}

	if err := validateType(input.TypeOfRepetition, input.CustomDates); err != nil {
		return UpdateInput{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	if input.Periodicity < 0 && input.Periodicity > 31 {
		return UpdateInput{}, fmt.Errorf("%w: invalid periodicity", ErrInvalidInput)
	}

	if input.ScheduledAt == "" {
		return UpdateInput{}, fmt.Errorf("%w: invalid scheduled_at", ErrInvalidInput)
	}

	return input, nil
}

func validateType(repetition taskdomain.Repetition, dates []string) error {
	var hasDate = len(dates) > 0

	if repetition != taskdomain.RepetitionSpecific && hasDate {
		return fmt.Errorf("custom_dates allowed only for specific_dates")
	}

	if repetition == taskdomain.RepetitionSpecific && !hasDate {
		return fmt.Errorf("custom_dates required for specific_dates")
	}

	return nil
}
