package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

// шлет фронт
type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`

	TypeOfRepetition taskdomain.Repetition `json:"type_of_repetition"`
	Periodicity      int8                  `json:"periodicity"`
	ScheduledAt      string                `json:"scheduled_at"`
	CustomDates      []string              `json:"custom_dates"`
}

type taskUpdateMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`

	CreatedAt        string                `json:"created_at"`
	TypeOfRepetition taskdomain.Repetition `json:"type_of_repetition"`
	Periodicity      int8                  `json:"periodicity"`
	ScheduledAt      string                `json:"scheduled_at"`
	CustomDates      []string              `json:"custom_dates"`
}

// maybe later
// type repetitionDTO struct {
// 	TypeOfRepetition taskdomain.Repetition `json:"type_of_repetition"`
// 	Periodicity      int8                  `json:"periodicity"`
// 	ScheduledAt      string                `json:"scheduled_at"`
// 	CustomDates      []string              `json:"custom_dates"`
// }

// шлем на фронт
type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`

	TypeOfRepetition taskdomain.Repetition `json:"type_of_repetition"`
	Periodicity      int8                  `json:"periodicity"`
	ScheduledAt      time.Time             `json:"scheduled_at"`
	CustomDates      []string              `json:"custom_dates"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,

		TypeOfRepetition: task.TypeOfRepetition,
		Periodicity:      task.Periodicity,
		ScheduledAt:      task.ScheduledAt,
		CustomDates:      task.CustomDates,
	}
}
