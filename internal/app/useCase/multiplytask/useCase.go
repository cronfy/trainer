package multiplytask

import (
	"math/rand"

	"github.com/cronfy/trainer/internal/app/domain"
)

type useCase struct{}

func New() *useCase {
	return &useCase{}
}

func (u *useCase) Get() domain.MultiplyTask {
	operands := make([]int, 2)

	operands[0] = randomBetween(1, 20)
	operands[1] = randomBetween(1, 11)

	if rand.Int31n(2) == 0 {
		operands[0], operands[1] = operands[1], operands[0]
	}

	return domain.MultiplyTask{
		Operands: []int{operands[0], operands[1]},
	}
}

func (u *useCase) Solve(task domain.MultiplyTask, solution int) bool {
	return task.Operands[0]*task.Operands[1] == solution
}

func randomBetween(min, max int) int {
	return int(rand.Int31n(int32(max-min))) + min
}
