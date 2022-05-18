package gossip

import (
	"math"
	"math/rand"
)

func randomSample[T comparable](k int, pop []T) []T {
	n := len(pop)
	if k >= n {
		panic("Population size must be greater than sample size")
	}
	sample := make([]T, k)
	idx := rand.Perm(n)[:k]
	for i, j := range idx {
		sample[i] = pop[j]
	}
	return sample
}

func retransmitLimit(multiplier, fanout, nodes int) int {
	conv := math.Ceil(logBase(float64(fanout+1), float64(nodes)))
	return multiplier * int(conv)

}

func logBase(base, x float64) float64 {
	return math.Log2(x) / math.Log2(base)
}
