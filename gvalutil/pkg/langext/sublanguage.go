package langext

import (
	"context"

	"github.com/PaesslerAG/gval"
)

// PrefixSublanguage is a helper that simply calls the given sublanguage as an extension.
func PrefixSublanguage(r rune, sub gval.Language, ext func(context.Context, *gval.Parser, gval.Evaluable) (gval.Evaluable, error)) gval.Language {
	return gval.PrefixExtension(r, func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
		eval, err := p.ParseSublanguage(ctx, sub)
		if err != nil {
			return nil, err
		}

		return ext(ctx, p, eval)
	})
}
