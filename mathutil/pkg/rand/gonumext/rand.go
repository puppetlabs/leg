// Package gonumext provides constructors for random factories using some of the
// algorithms provided by the Gonum PRNG library.
//
// It is provided as a separate package here so that users who do not need it do
// not bring in Gonum as a dependency when they use our rand package.
package gonumext

import (
	"github.com/puppetlabs/leg/mathutil/pkg/rand"
	exprand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mathext/prng"
)

// NewMersenneTwisterFactory creates a factory for producing RNGs using the
// common Mersenne Twister with the Mersenne prime 2^19937-1. Each RNG it
// produces will be seeded with the given seeder.
func NewMersenneTwisterFactory(seeder rand.Seeder) *rand.ExpRandFactory {
	return rand.NewExpRandFactory(func(seed uint64) exprand.Source {
		alg := prng.NewMT19937()
		alg.Seed(seed)
		return alg
	}, seeder)
}

// NewMersenneTwister64Factory creates a factory for producing RNGs using the
// common Mersenne Twister with the Mersenne prime 2^19937-1 and 64-bit word
// length. Each RNG it produces will be seeded with the given seeder.
func NewMersenneTwister64Factory(seeder rand.Seeder) *rand.ExpRandFactory {
	return rand.NewExpRandFactory(func(seed uint64) exprand.Source {
		alg := prng.NewMT19937_64()
		alg.Seed(seed)
		return alg
	}, seeder)
}
