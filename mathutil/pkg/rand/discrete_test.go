package rand_test

import (
	"log"
	"math"
	"math/cmplx"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/puppetlabs/leg/mathutil/pkg/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func numSamplesPerTest() int {
	if testing.Short() {
		// Note that you might encounter more spurious type I/II errors in short
		// mode. If in doubt, rerun the test in question without -short.
		return 250_000
	}

	return 2_500_000
}

// uniformCDF returns the cumulative value along a specified distribution.
func uniformCDF(min, max, x float64) float64 {
	if x < min {
		return 0
	}
	if x > max {
		return 1
	}
	return (x - min) / (max - min)
}

var (
	// Mapping of confidence levels (in basis points because they're easier to
	// work with numerically) to K(alpha) from the Kolmogorov distribution.
	KolmogorovKas = map[int]float64{
		10:   1.95,
		100:  1.6722,
		500:  1.3581,
		1000: 1.22385,
		2000: 1.07275,
	}
)

// kolmogorovSmirnovUniform finds the largest distance between the uniform
// distribution and the given samples.
func kolmogorovSmirnovUniform(samples []float64, distMin, distMax float64) float64 {
	sort.Float64s(samples)

	n := float64(len(samples))
	var d float64
	for i, sample := range samples {
		y := uniformCDF(distMin, distMax, sample)
		di := math.Max(y-float64(i)/n, float64(i+1)/n-y)
		if di > d {
			d = di
		}
	}

	return d
}

// 2D types for the minimum distance test.
type point complex128

func (p point) X() float64 { return real(p) }
func (p point) Y() float64 { return imag(p) }
func (p point) Dist(q point) float64 {
	return cmplx.Abs(complex128(p) - complex128(q))
}

type pointsX []point

func (px pointsX) Len() int           { return len(px) }
func (px pointsX) Swap(i, j int)      { px[i], px[j] = px[j], px[i] }
func (px pointsX) Less(i, j int) bool { return px[i].X() < px[j].X() }

// The parameter Q and the correction algorithm below are from M. Fischler
// (2002). https://www.osti.gov/biblio/794005
var (
	minimumDistanceWidth = 10_000.0
	minimumDistanceArea  = minimumDistanceWidth * minimumDistanceWidth
	minimumDistanceQ     = 3.0 * math.Sqrt(3.0) / (4.0 * math.Pi)
)

func testMinimumDistance(t *testing.T, rng rand.Rand, dists []float64, points []point, gen func() float64) {
	// This is the minimum distance test from the diehard tests by George
	// Marsaglia with adjustments from the paper by Fischler mentioned above.

	log.Printf("[%s] Number of squares: %d", t.Name(), len(dists))
	log.Printf("[%s] Number of samples per square: %d", t.Name(), len(points))

	n := float64(len(points))

	// Calculate the mean of the expected distribution.
	u := 2.0 * minimumDistanceArea / (n * (n - 1) * math.Pi)
	log.Printf("[%s] Expected μ: %f", t.Name(), u)
	log.Printf("[%s] Points to generate: %d", t.Name(), len(points))

	x := -1.0 / u

	for di := 0; di < len(dists); di++ {
		// Generate points.
		generationStart := time.Now()
		for i := 0; i < len(points); i++ {
			x, y := gen(), gen()
			points[i] = point(complex(float64(x), float64(y)))
		}
		log.Printf("[%s] [#%d] Generated points in %s", t.Name(), di, time.Since(generationStart))

		// Calculate minimum distance using brute force with a quick hack to
		// jump out of the loop early if possible. This can be further optimized
		// if takes too long.
		//
		// https://en.wikipedia.org/wiki/Closest_pair_of_points_problem
		minDistStart := time.Now()
		sort.Sort(pointsX(points))
		minDist := points[0].Dist(points[1])
		for pi := 0; pi < len(points)-1; pi++ {
			for pj := pi + 1; pj < len(points); pj++ {
				if points[pi].X()-points[pj].X() >= minDist {
					break
				}

				dist := points[pi].Dist(points[pj])
				if dist < minDist {
					minDist = dist
				}
			}
		}
		log.Printf("[%s] [#%d] Minimum distance: %f (in %s)", t.Name(), di, minDist, time.Since(minDistStart))

		a2 := minDist * minDist
		log.Printf("[%s] [#%d] a²: %f", t.Name(), di, a2)

		q := 1.0 + ((2.0+minimumDistanceQ)/6.0)*math.Pi*math.Pi*n*n*n*a2*a2/(minimumDistanceArea*minimumDistanceArea)
		log.Printf("[%s] [#%d] Overlap correction (q): %f", t.Name(), di, q)

		// P-value comes from an exponential distribution with mean u.
		dists[di] = 1 - math.Exp(x*a2)*q
		log.Printf("[%s] [#%d] P-value: %f", t.Name(), di, dists[di])
	}

	// Now we test whether our samples are uniformly distributed.
	d := kolmogorovSmirnovUniform(dists, 0, 1)
	log.Printf("[%s] Kolmogorov-Smirnov test statistic (D): %f", t.Name(), d)

	// At alpha = 0.001 (expressed here in bips, see above), we expect about 1
	// in every 1,000 tests to fail spuriously. Beyond that, we're in the
	// territory of diminishing returns.
	//
	// When we cause real problems with the RNG (like commenting out the modulo
	// bias avoidance code) our distribution breaks in a way so improbable no
	// reasonable significance level would not catch it.
	//
	// You can try other values of alpha here, for fun or deeper analysis, if
	// they exist in the table above.
	alpha := 10

	// For a goodness-of-fit test we reject the null hypothesis if D >
	// K_alpha / sqrt(len(dists)), where Pr(K <= K_alpha) = 1 - alpha.
	k := KolmogorovKas[alpha] / math.Sqrt(float64(len(dists)))
	log.Printf("[%s] Kolmogorov-Smirnov critical value: %f", t.Name(), k)

	assert.False(t,
		d > k,
		"Kolmogorov-Smirnov test rejects null hypothesis with D = %f, α = %f, K = %f",
		d, float64(alpha)/10000, k,
	)
}

func prepareMinimumDistance() (dists []float64, points []point) {
	// We want to produce 2500 sets of points for our uniformity check. Note
	// that lowering this number will make it much easier for false positives to
	// appear.
	dists = make([]float64, 2500)

	// Then we can determine how many points we should make per square.
	points = make([]point, numSamplesPerTest()/len(dists))

	return
}

func TestUint64N(t *testing.T) {
	rng, err := rand.DefaultFactory.New()
	require.NoError(t, err)

	dists, points := prepareMinimumDistance()

	tests := []uint64{
		100_000,
		1_000_000,
		// This test checks the binary-and fast path.
		1 << 32,
		// This test checks the modulo biasing of the function.
		math.MaxUint64 / 3 * 2,
	}
	for _, test := range tests {
		t.Run(strconv.FormatUint(test, 10), func(t *testing.T) {
			gen := func() float64 {
				rv, err := rand.Uint64N(rng, test)
				require.NoError(t, err)
				assert.Less(t, rv, test)

				return float64(rv) * minimumDistanceWidth / float64(test)
			}
			testMinimumDistance(t, rng, dists, points, gen)
		})
	}
}

func TestFloat64(t *testing.T) {
	rng, err := rand.DefaultFactory.New()
	require.NoError(t, err)

	dists, points := prepareMinimumDistance()

	gen := func() float64 {
		rv, err := rand.Float64(rng)
		require.NoError(t, err)
		assert.Less(t, rv, 1.0)

		return rv * minimumDistanceWidth
	}
	testMinimumDistance(t, rng, dists, points, gen)
}
