package fn

import (
	"fmt"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
)

func KeywordInvocationAnnotation(name string) evaluate.Annotation {
	return evaluate.Annotation{
		Name: "leg.relay.sh/fn",
		Attributes: map[string]any{
			"name":  name,
			"style": "keyword",
		},
		Description: fmt.Sprintf("invocation of function %s", name),
	}
}

func PositionalInvocationAnnotation(name string) evaluate.Annotation {
	return evaluate.Annotation{
		Name: "leg.relay.sh/fn",
		Attributes: map[string]any{
			"name":  name,
			"style": "keyword",
		},
		Description: fmt.Sprintf("invocation of function %s", name),
	}
}

func ArgAnnotation(id any) evaluate.Annotation {
	return evaluate.Annotation{
		Name: "leg.relay.sh/fn.arg",
		Attributes: map[string]any{
			"arg": id,
		},
		Description: fmt.Sprintf("arg %s", id),
	}
}
