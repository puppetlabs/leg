package template

import (
	"context"

	"github.com/PaesslerAG/gval"
)

type DelimitedLanguage struct {
	Start, End string
	Language   gval.Language
	Extension  func(ctx context.Context, p *gval.Parser, a gval.Evaluable) (gval.Evaluable, error)
}

type DelimitedLanguageFactory interface {
	DelimitedLanguage(tpl gval.Language) *DelimitedLanguage
}

type DelimitedLanguageFactoryFunc func(tpl gval.Language) *DelimitedLanguage

var _ DelimitedLanguageFactory = DelimitedLanguageFactoryFunc(nil)

func (dlff DelimitedLanguageFactoryFunc) DelimitedLanguage(tpl gval.Language) *DelimitedLanguage {
	return dlff(tpl)
}

func DelimitedExpressionLanguageFactory(dl *DelimitedLanguage) DelimitedLanguageFactory {
	return DelimitedLanguageFactoryFunc(func(_ gval.Language) *DelimitedLanguage {
		return dl
	})
}
