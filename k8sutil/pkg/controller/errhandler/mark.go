package errhandler

import (
	"github.com/puppetlabs/leg/errmap/pkg/errmark"
	"github.com/puppetlabs/leg/k8sutil/pkg/controller/obj/lifecycle"
	"k8s.io/apimachinery/pkg/api/errors"
)

var (
	RuleIsConflict = errmark.RuleAny(
		errmark.RuleFunc(errors.IsConflict),
		errmark.RuleFunc(errors.IsAlreadyExists),
	)

	RuleIsTimeout = errmark.RuleAny(
		errmark.RuleFunc(errors.IsTimeout),
		errmark.RuleFunc(errors.IsServerTimeout),
	)

	RuleIsForbidden = errmark.RuleFunc(errors.IsForbidden)

	RuleIsRequired = errmark.RuleType(&lifecycle.RequiredError{})
)
