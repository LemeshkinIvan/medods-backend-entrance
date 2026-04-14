package task

import "errors"

var ErrInvalidInput = errors.New("invalid task input")
var ErrCustomDateBefore = errors.New("Custom time cannot be earlier than the task creation")
var ErrCustomDateAfter = errors.New("Custom time cannot be later than scheduleAt")
