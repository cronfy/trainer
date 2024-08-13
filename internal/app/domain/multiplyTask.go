package domain

type MultiplyTask struct {
	operands [2]int
}

func NewMultiplyTask(operands [2]int) MultiplyTask {
	return MultiplyTask{operands: operands}
}

func (m *MultiplyTask) GetOperands() [2]int {
	return m.operands
}

func (m *MultiplyTask) Verify(solution int) bool {
	return m.operands[0]*m.operands[1] == solution
}
