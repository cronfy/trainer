package multiplytask

//go:generate mockery

type randomTool interface {
	RandomBetween(min, max int) int
	Chance(chance int8) bool
}
