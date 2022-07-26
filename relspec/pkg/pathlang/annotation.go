package pathlang

import (
	"fmt"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
)

var (
	ExpressionAnnotation = evaluate.Annotation{
		Name:        "leg.relay.sh/expression",
		Description: "expression",
	}
)

func InfixOperationAnnotation(name string) evaluate.Annotation {
	return evaluate.Annotation{
		Name: "leg.relay.sh/operation",
		Attributes: map[string]any{
			"notation": "infix",
			"arity":    2,
			"name":     name,
		},
		Description: fmt.Sprintf("infix operation %s", name),
	}
}

func PrefixOperationAnnotation(name string) evaluate.Annotation {
	return evaluate.Annotation{
		Name: "leg.relay.sh/operation",
		Attributes: map[string]any{
			"notation": "prefix",
			"arity":    1,
			"name":     name,
		},
		Description: fmt.Sprintf("prefix operation %s", name),
	}
}

func OperandAnnotation(n int) evaluate.Annotation {
	return evaluate.Annotation{
		Name: "leg.relay.sh/operation.operand",
		Attributes: map[string]any{
			"number": n,
		},
		Description: fmt.Sprintf("operand %d", n),
	}
}
