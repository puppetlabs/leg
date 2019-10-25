package testutil

import (
	"fmt"

	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
	"github.com/puppetlabs/errawr-go/v2/pkg/impl"
)

func NewStubError(code string) errawr.Error {
	return &impl.Error{
		ErrorArguments: impl.ErrorArguments{},
		ErrorCode:      code,
		ErrorDescription: &impl.ErrorDescription{
			Friendly:  "A test error occurred.",
			Technical: "A test error occurred.",
		},
		ErrorDomain: &impl.ErrorDomain{
			Key:   "!test",
			Title: "Test",
		},
		ErrorSection: &impl.ErrorSection{
			Key:   "stub",
			Title: "Stub errors",
		},
		ErrorSensitivity: errawr.ErrorSensitivityNone,
		ErrorTitle:       "Stub error",
		Version:          errawr.Version,
	}
}

func IsStubError(code string, err errawr.Error) bool {
	return err != nil && err.Is(fmt.Sprintf("!test_stub_%s", code))
}
