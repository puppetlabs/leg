package algo

import (
	"testing"

	"github.com/reflect/gographt"
	"github.com/stretchr/testify/assert"
)

func TestTiernanSimpleCycles(t *testing.T) {
	g := gographt.NewDirectedPseudographWithFeatures(gographt.DeterministicIteration)
	for i := 0; i < 7; i++ {
		g.AddVertex(i)
	}

	g.Connect(0, 0)
	assert.Equal(t, [][]gographt.Vertex{{0}}, TiernanSimpleCyclesOf(g).Cycles())

	g.Connect(1, 1)
	assert.Equal(t, [][]gographt.Vertex{{0}, {1}}, TiernanSimpleCyclesOf(g).Cycles())

	g.Connect(0, 1)
	g.Connect(1, 0)
	assert.Equal(t, [][]gographt.Vertex{{0, 1}, {0}, {1}}, TiernanSimpleCyclesOf(g).Cycles())

	g.Connect(1, 2)
	g.Connect(2, 3)
	g.Connect(3, 0)
	assert.Equal(t, [][]gographt.Vertex{{0, 1, 2, 3}, {0, 1}, {0}, {1}}, TiernanSimpleCyclesOf(g).Cycles())

	g.Connect(6, 6)
	assert.Equal(t, [][]gographt.Vertex{{0, 1, 2, 3}, {0, 1}, {0}, {1}, {6}}, TiernanSimpleCyclesOf(g).Cycles())

	ns := map[int]int{
		1: 1,
		2: 3,
		3: 8,
		4: 24,
		5: 89,
		6: 415,
		7: 2372,
		8: 16072,
		9: 125673,
	}

	for size, count := range ns {
		g = gographt.NewDirectedPseudograph()
		for i := 0; i < size; i++ {
			g.AddVertex(i)
		}
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				g.Connect(i, j)
			}
		}

		assert.Len(t, TiernanSimpleCyclesOf(g).Cycles(), count)
	}
}
