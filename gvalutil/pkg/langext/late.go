package langext

import (
	"context"

	"github.com/PaesslerAG/gval"
)

// Late is a helper to provide late-binding for a language; that is, the
// completed language is passed back in as an argument.
func Late(fn func(complete gval.Language) gval.Language) (bound gval.Language) {
	complete := gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
		return p.ParseSublanguage(ctx, bound)
	})
	bound = fn(complete)
	return
}
