package supervisor

import "errors"

// errorsAsType provides typed error unwrapping on Go versions before 1.26.
func errorsAsType[T error](err error) (T, bool) {
	var target T
	ok := errors.As(err, &target)
	return target, ok
}
