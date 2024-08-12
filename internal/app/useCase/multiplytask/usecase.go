package multiplytask

import (
	"github.com/cronfy/trainer/internal/app/domain"
)

type useCase struct {
	randomTool randomTool
}

func New(randomTool randomTool) *useCase {
	return &useCase{randomTool}
}

func (u *useCase) Get() domain.MultiplyTask {
	operands := make([]int, 2)

	operands[0] = u.randomTool.RandomBetween(1, 20)
	operands[1] = u.randomTool.RandomBetween(1, 11)

	if u.randomTool.Chance(50) {
		operands[0], operands[1] = operands[1], operands[0]
	}

	return domain.MultiplyTask{
		Operands: []int{operands[0], operands[1]},
	}
}

func (u *useCase) Solve(task domain.MultiplyTask, solution int) bool {
	return task.Operands[0]*task.Operands[1] == solution
}
