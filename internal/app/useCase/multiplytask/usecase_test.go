package multiplytask_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	app "github.com/cronfy/trainer/internal/app/domain"
	"github.com/cronfy/trainer/internal/app/useCase/multiplytask"
	"github.com/cronfy/trainer/internal/app/useCase/multiplytask/mocks"
)

type useCase interface {
	Get() app.MultiplyTask
	Solve(task app.MultiplyTask, solution int) bool
}

type testObjects struct {
	randomTool *mocks.RandomTool
	useCase    useCase
}

func makeTestObjects(t *testing.T) *testObjects {
	randTool := mocks.NewRandomTool(t)

	return &testObjects{
		randomTool: randTool,
		useCase:    multiplytask.New(randTool),
	}
}

func TestUseCase_Get_Success(t *testing.T) {
	to := makeTestObjects(t)

	to.randomTool.EXPECT().RandomBetween(1, 20).Return(10)
	to.randomTool.EXPECT().RandomBetween(1, 11).Return(7)
	to.randomTool.EXPECT().Chance(int8(50)).Return(false)

	want := app.NewMultiplyTask([2]int{10, 7})
	got := to.useCase.Get()

	assert.Equal(t, want, got)
}

func TestUseCase_Get_SwitchesOperands(t *testing.T) {
	to := makeTestObjects(t)

	to.randomTool.EXPECT().RandomBetween(1, 20).Return(10)
	to.randomTool.EXPECT().RandomBetween(1, 11).Return(7)
	to.randomTool.EXPECT().Chance(int8(50)).Return(true)

	want := app.NewMultiplyTask([2]int{7, 10})
	got := to.useCase.Get()

	assert.Equal(t, want, got)
}

func TestUseCase_Solve_Success(t *testing.T) {
	to := makeTestObjects(t)

	got := to.useCase.Solve(app.NewMultiplyTask([2]int{5, 7}), 35)

	assert.Equal(t, true, got)
}

func TestUseCase_Solve_WrongAnswer(t *testing.T) {
	to := makeTestObjects(t)

	got := to.useCase.Solve(app.NewMultiplyTask([2]int{5, 7}), 25)

	assert.Equal(t, false, got)
}
