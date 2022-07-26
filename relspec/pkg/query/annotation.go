package query

import (
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
)

var (
	QueryAnnotation = evaluate.Annotation{
		Name:        "leg.relay.sh/query",
		Description: "query",
	}
	QueryResultAnnotation = evaluate.Annotation{
		Name:        "leg.relay.sh/query.expansion",
		Description: "query result",
	}
	TraversalAnnotation = evaluate.Annotation{
		Name:        "leg.relay.sh/traversal",
		Description: "traversal",
	}
	TraversalKeyAnnotation = evaluate.Annotation{
		Name:        "leg.relay.sh/traversal.key",
		Description: "traversal key",
	}
	TraversalDataAnnotation = evaluate.Annotation{
		Name:        "leg.relay.sh/traversal.data",
		Description: "traversal data",
	}
)
