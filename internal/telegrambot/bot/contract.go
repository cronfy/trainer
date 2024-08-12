package bot

import (
	"oap/trainer/internal/app/domain"
)

type multiplyTaskUseCase interface {
	Get() domain.MultiplyTask
	Solve(task domain.MultiplyTask, solution int) bool
}
