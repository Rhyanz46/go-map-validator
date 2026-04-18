package map_validator

import "net/http"

// ValidateJSON is a one-shot helper that runs the full Load→Validate→Bind
// pipeline against a JSON HTTP request and returns the bound value of type T.
//
// It is functionally equivalent to:
//
//	op, err := NewValidateBuilder().SetRules(rules).LoadJsonHttp(r)
//	...
//	extra, err := op.RunValidate()
//	...
//	var out T
//	err = extra.Bind(&out)
//
// Errors are forwarded verbatim from the underlying layers (ErrNoRules,
// ErrInvalidJsonFormat, validation errors) so callers keep full fidelity for
// logging and response messages.
//
// A nil request returns an error — it never panics. Empty rules return
// ErrNoRules. The returned T is the zero value on any error path.
func ValidateJSON[T any](r *http.Request, rules RulesWrapper) (T, error) {
	var zero T
	op, err := NewValidateBuilder().SetRules(rules).LoadJsonHttp(r)
	if err != nil {
		return zero, err
	}
	extra, err := op.RunValidate()
	if err != nil {
		return zero, err
	}
	var out T
	if err := extra.Bind(&out); err != nil {
		return zero, err
	}
	return out, nil
}
