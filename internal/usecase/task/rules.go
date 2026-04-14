package task

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func isActive(task taskdomain.Task, date time.Time) bool {
	switch task.TypeOfRepetition {

	case taskdomain.RepetitionDaily:
		if task.Periodicity <= 0 {
			return false
		}

		days := int8(date.Sub(task.CreatedAt).Hours() / 24)

		return days%task.Periodicity == 0

	case taskdomain.RepetitionMontly:
		return date.Day() == int(task.Periodicity)

	case taskdomain.RepetitionSpecific:
		return customDatesRuleStr(task.CustomDates, date)

	case taskdomain.RepetitionEven:
		return date.Day()%2 == 0

	case taskdomain.RepetitionOdd:
		return date.Day()%2 != 0
	}

	return false
}

// func between(t, start, end time.Time) bool {
// 	return !t.Before(start) && !t.After(end)
// }

func customDatesRuleStr(dates []string, date time.Time) bool {
	target := date.Format("2006-01-02")

	for _, d := range dates {
		if d == target {
			return true
		}
	}

	return false
}
