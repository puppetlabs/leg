package encoding

import (
	"fmt"
	"strings"

	errawr "github.com/puppetlabs/errawr-go"
	"github.com/puppetlabs/errawr-go/impl"
)

type ErrorDisplayEnvelope struct {
	Domain      string                  `json:"domain"`
	Section     string                  `json:"section"`
	Code        string                  `json:"code"`
	Title       string                  `json:"title"`
	Description *ErrorDescription       `json:"description,omitempty"`
	Arguments   map[string]interface{}  `json:"arguments,omitempty"`
	Formatted   *ErrorDescription       `json:"formatted,omitempty"`
	Causes      []*ErrorDisplayEnvelope `json:"causes,omitempty"`
}

func (ede ErrorDisplayEnvelope) AsError() errawr.Error {
	arguments := make(impl.ErrorArguments, len(ede.Arguments))
	for name, argument := range ede.Arguments {
		if argument == nil {
			continue
		}

		arguments[name] = &impl.ErrorArgument{
			Value: argument,
		}
	}

	prefix := fmt.Sprintf(`%s_%s_`, ede.Domain, ede.Section)
	code := strings.TrimPrefix(ede.Code, prefix)

	var e errawr.Error = &impl.Error{
		Version: errawr.Version,
		ErrorDomain: &impl.ErrorDomain{
			Key: ede.Domain,
		},
		ErrorSection: &impl.ErrorSection{
			Key: ede.Section,
		},
		ErrorCode:  code,
		ErrorTitle: ede.Title,
		ErrorDescription: &impl.ErrorDescription{
			Friendly:  ede.Description.Friendly,
			Technical: ede.Description.Technical,
		},
		ErrorArguments:   arguments,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSensitivity: errawr.ErrorSensitivityEdge,
	}

	for _, cause := range ede.Causes {
		if cause == nil {
			continue
		}

		e = e.WithCause(cause.AsError())
	}

	return e
}

func ForDisplay(e errawr.Error) *ErrorDisplayEnvelope {
	return ForDisplayWithSensitivity(e, errawr.ErrorSensitivityEdge)
}

func ForDisplayWithSensitivity(e errawr.Error, sensitivity errawr.ErrorSensitivity) *ErrorDisplayEnvelope {
	ede := &ErrorDisplayEnvelope{
		Domain:  e.Domain().Key(),
		Section: e.Section().Key(),
		Code:    e.ID(),
		Title:   e.Title(),
	}

	if e.Sensitivity() > sensitivity {
		return ede
	}

	causes := e.Causes()

	ede.Causes = make([]*ErrorDisplayEnvelope, len(causes))
	for i, cause := range causes {
		ede.Causes[i] = ForDisplay(cause)
	}

	ede.Description = &ErrorDescription{
		Friendly:  e.Description().Friendly(),
		Technical: e.Description().Technical(),
	}
	ede.Arguments = e.Arguments()
	ede.Formatted = &ErrorDescription{
		Friendly:  e.FormattedDescription().Friendly(),
		Technical: e.FormattedDescription().Technical(),
	}

	return ede
}
