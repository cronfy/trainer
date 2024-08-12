package random

import (
	"math/rand"
)

type tool struct {
}

func New() *tool {
	return &tool{}
}

func (*tool) RandomBetween(min, max int) int {
	return int(rand.Int31n(int32(max-min))) + min
}

func (*tool) Chance(chance int8) bool {
	return chance > int8(rand.Int31n(100))
}
