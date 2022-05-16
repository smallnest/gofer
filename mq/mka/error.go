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
	if len(err) == 0 {
		return fmt.Sprintf("kafka write errors (%d/%d)", err.Count(), len(err))
	}

	if len(err) == 1 {
		return fmt.Sprintf("kafka write errors (%d/%d), error: %v", err.Count(), len(err), err[0])
	}

	return fmt.Sprintf("kafka write errors (%d/%d), the first two errors: %v; %v", err.Count(), len(err), err[0], err[1])

}
