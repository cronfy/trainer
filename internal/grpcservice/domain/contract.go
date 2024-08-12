package domain

import (
	"github.com/cronfy/trainer/internal/app/domain"
)

type MultiplyTaskUseCase interface {
	Get() domain.MultiplyTask
	Solve(task domain.MultiplyTask, solution int) bool
}
