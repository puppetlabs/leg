package fnlib_test

import (
	"context"

	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/puppetlabs/leg/relspec/pkg/relspec"
)

type testID struct {
	Name string
}

func (ti testID) Less(other testID) bool {
	return ti.Name < other.Name
}

type testReferences = *ref.Log[testID]

type testMappingTypeResolver struct{}

var _ relspec.MappingTypeResolver[testReferences] = &testMappingTypeResolver{}

func (*testMappingTypeResolver) ResolveMappingType(ctx context.Context, tm map[string]any) (*evaluate.Result[testReferences], error) {
	return evaluate.ContextualizedResult(
		evaluate.NewMetadata(ref.InitialLog(ref.Observed(testID{Name: "nope"}))),
	), nil
}
