package mka

import "fmt"

type WriteErrors []error

// Count counts the number of non-nil errors in err.
func (err WriteErrors) Count() int {
	n := 0

	for _, e := range err {
		if e != nil {
			n++
		}
	}

	return n
}

func (err WriteErrors) Error() string {
	var errors []error

	n := 0
	for _, e := range err {
		e := e
		if e != nil {
			n++
			if n < 3 {
				errors = append(errors, e)
			}
		}
	}

	if n == 0 {
		return fmt.Sprintf("kafka write errors (%d/%d)", n, len(err))
	}

	if n == 1 {
		return fmt.Sprintf("kafka write errors (%d/%d), error: %v", n, len(err), errors[0])
	}

	return fmt.Sprintf("kafka write errors (%d/%d), the first two errors: %v; %v", n, len(err), errors[0], errors[1])

}
