package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// режимы повторений
type Repetition string

const (
	RepetitionDaily    Repetition = "daily"
	RepetitionMontly   Repetition = "monthly"
	RepetitionSpecific Repetition = "specific_dates"
	RepetitionEven     Repetition = "even"
	RepetitionOdd      Repetition = "odd"
)

// refactor later. maybe map...
func (s Repetition) Valid() bool {
	switch s {
	case RepetitionDaily, RepetitionMontly, RepetitionSpecific, RepetitionEven, RepetitionOdd:
		return true
	default:
		return false
	}
}

// модель таблицы из БД
type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`

	TypeOfRepetition Repetition `json:"type_of_repetition"`
	Periodicity      int8       `json:"periodicity"`
	ScheduledAt      time.Time  `json:"scheduled_at"`
	CustomDates      []string   `json:"custom_dates"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
