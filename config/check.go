// Copyright © 2017 thingful

package config

import "fmt"

// Validate returns the first error it finds
func Validate(config *Configuration) error {

	checkers := []checker{isEmpty, uniqueUIDs, uniquePipeURIs, validatePipes}

	for _, c := range checkers {

		ok, result := c(config)
		if !ok {
			return result
		}
	}
	return nil
}

type severity string

const (
	ERROR   = severity("Error")
	WARNING = severity("Warning")
)

type checker func(config *Configuration) (bool, *result)

type result struct {
	Severity severity
	Message  string
}

func (r *result) Error() string {
	return fmt.Sprintf("%s : %s", r.Severity, r.Message)
}

func isEmpty(config *Configuration) (bool, *result) {

	if config == nil {
		return false, &result{
			Severity: ERROR,
			Message:  "empty configuration",
		}
	}

	if len(config.Pipes) == 0 {
		return false, &result{
			Severity: WARNING,
			Message:  "no pipes configured",
		}
	}

	return true, nil
}

func uniqueUIDs(config *Configuration) (bool, *result) {
	return true, nil
}

func uniquePipeURIs(config *Configuration) (bool, *result) {
	return true, nil
}

func validatePipes(config *Configuration) (bool, *result) {
	return true, nil
}
