package encoding

import (
	errawr "github.com/puppetlabs/errawr-go"
	"github.com/puppetlabs/errawr-go/impl"
)

type ErrorTransitEnvelope struct {
	Version     uint64                  `json:"version"`
	Domain      ErrorDomain             `json:"domain"`
	Section     ErrorSection            `json:"section"`
	Code        string                  `json:"code"`
	Title       string                  `json:"title"`
	Description ErrorDescription        `json:"description"`
	Arguments   ErrorArguments          `json:"arguments"`
	Metadata    ErrorMetadata           `json:"metadata,omitempty"`
	Causes      []*ErrorTransitEnvelope `json:"causes"`
	Buggy       bool                    `json:"buggy"`
	Sensitivity errawr.ErrorSensitivity `json:"sensitivity"`
}

func (ete ErrorTransitEnvelope) AsError() errawr.Error {
	sensitivity := ete.Sensitivity
	if ete.Buggy && sensitivity < errawr.ErrorSensitivityBug {
		sensitivity = errawr.ErrorSensitivityBug
	}

	var arguments impl.ErrorArguments
	if len(ete.Arguments) > 0 {
		arguments = make(impl.ErrorArguments, len(ete.Arguments))
		for name, argument := range ete.Arguments {
			if argument == nil {
				continue
			}

			arguments[name] = &impl.ErrorArgument{
				Value:       argument.Value,
				Description: argument.Description,
			}
		}
	}

	metadata := &impl.ErrorMetadata{}
	if hm := ete.Metadata.HTTPErrorMetadata; hm != nil {
		metadata.HTTPErrorMetadata = &impl.HTTPErrorMetadata{
			ErrorStatus:  hm.Status,
			ErrorHeaders: hm.Headers,
		}
	}

	var e errawr.Error = &impl.Error{
		Version: ete.Version,
		ErrorDomain: &impl.ErrorDomain{
			Key:   ete.Domain.Key,
			Title: ete.Domain.Title,
		},
		ErrorSection: &impl.ErrorSection{
			Key:   ete.Section.Key,
			Title: ete.Section.Title,
		},
		ErrorCode:  ete.Code,
		ErrorTitle: ete.Title,
		ErrorDescription: &impl.ErrorDescription{
			Friendly:  ete.Description.Friendly,
			Technical: ete.Description.Technical,
		},
		ErrorArguments:   arguments,
		ErrorMetadata:    metadata,
		ErrorSensitivity: sensitivity,
	}

	for _, cause := range ete.Causes {
		if cause == nil {
			continue
		}

		e = e.WithCause(cause.AsError())
	}

	if ete.Buggy {
		e = e.Bug()
	}

	return e
}

func ForTransit(e errawr.Error) *ErrorTransitEnvelope {
	ete := &ErrorTransitEnvelope{
		Version: errawr.Version,
		Domain: ErrorDomain{
			Key:   e.Domain().Key(),
			Title: e.Domain().Title(),
		},
		Section: ErrorSection{
			Key:   e.Section().Key(),
			Title: e.Section().Title(),
		},
		Code:  e.Code(),
		Title: e.Title(),
		Description: ErrorDescription{
			Friendly:  e.Description().Friendly(),
			Technical: e.Description().Technical(),
		},
		Arguments:   make(ErrorArguments, len(e.Arguments())),
		Sensitivity: e.Sensitivity(),
		Causes:      make([]*ErrorTransitEnvelope, len(e.Causes())),
		Buggy:       e.IsBug(),
	}

	for name, argument := range e.Arguments() {
		ete.Arguments[name] = &ErrorArgument{
			Value:       argument,
			Description: e.ArgumentDescription(name),
		}
	}

	if hm, ok := e.Metadata().HTTP(); ok {
		ete.Metadata.HTTPErrorMetadata = &HTTPErrorMetadata{
			Status:  hm.Status(),
			Headers: hm.Headers(),
		}
	}

	for i, cause := range e.Causes() {
		ete.Causes[i] = ForTransit(cause)
	}

	return ete
}
