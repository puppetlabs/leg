package template

type WithStringFormatter struct {
	StringFormatter
}

var _ StringJoinerOption = WithStringFormatter{}

func (wsf WithStringFormatter) ApplyToStringJoinerOptions(target *StringJoinerOptions) {
	target.Formatter = wsf.StringFormatter
}

type WithEmptyStringsEliminated bool

var _ StringJoinerOption = WithEmptyStringsEliminated(false)

func (wese WithEmptyStringsEliminated) ApplyToStringJoinerOptions(target *StringJoinerOptions) {
	target.EliminateEmptyStrings = bool(wese)
}

type WithJoiner struct {
	Joiner
}

var _ Option = WithJoiner{}

func (wj WithJoiner) ApplyToOptions(target *Options) {
	target.Joiner = wj.Joiner
}

type WithDelimitedLanguage struct {
	*DelimitedLanguage
}

var _ Option = WithDelimitedLanguage{}

func (wdl WithDelimitedLanguage) ApplyToOptions(target *Options) {
	target.DelimitedLanguageFactories = append(target.DelimitedLanguageFactories, DelimitedExpressionLanguageFactory(wdl.DelimitedLanguage))
}

type WithDelimitedLanguageFactory struct {
	DelimitedLanguageFactory
}

var _ Option = WithDelimitedLanguageFactory{}

func (wdlf WithDelimitedLanguageFactory) ApplyToOptions(target *Options) {
	target.DelimitedLanguageFactories = append(target.DelimitedLanguageFactories, wdlf.DelimitedLanguageFactory)
}
