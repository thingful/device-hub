// Copyright © 2017 thingful

package config

import "fmt"

// Validate returns the first error it finds
func Validate(config *Configuration) error {

	checkers := []checker{isEmpty, noEmptyUIDs, uniqueUIDs, uniquePipeURIs, validatePipes}

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
			Severity: ERROR,
			Message:  "no pipes configured",
		}
	}

	if len(config.Listeners) == 0 {
		return false, &result{
			Severity: ERROR,
			Message:  "no listeners configured",
		}
	}

	if len(config.Endpoints) == 0 {
		return false, &result{
			Severity: ERROR,
			Message:  "no endpoints configured",
		}
	}

	return true, nil
}

func uniqueUIDs(config *Configuration) (bool, *result) {

	ok := true
	var ret *result
	uids := map[UID]bool{}

	f := func(e Endpoint) {

		alreadySeen, _ := uids[e.UID]
		if alreadySeen {
			ok = false
			ret = &result{
				Severity: ERROR,
				Message:  fmt.Sprintf("duplicate uid : %s", e.UID),
			}

			return
		}

		uids[e.UID] = true
	}

	config.Listeners.mapper(f)
	config.Endpoints.mapper(f)

	return ok, ret
}

func uniquePipeURIs(config *Configuration) (bool, *result) {

	ok := true
	var ret *result
	uris := map[string]bool{}

	f := func(p pipe) {

		alreadySeen, _ := uris[p.Uri]
		if alreadySeen {
			ok = false
			ret = &result{
				Severity: ERROR,
				Message:  fmt.Sprintf("duplicate pipe uri : %s", p.Uri),
			}

			return
		}

		uris[p.Uri] = true

	}

	config.Pipes.mapper(f)

	return ok, ret
}

func validatePipes(config *Configuration) (bool, *result) {

	// for each pipe there should be a profile, listener and endpoints
	// defined in the configuration
	ok := true
	var ret *result

	f := func(p pipe) {

		profileExists, _ := config.Profiles.FindByUID(p.Profile)

		if !profileExists {
			ok = false
			ret = &result{
				Severity: ERROR,
				Message:  fmt.Sprintf("missing profile (%s) for pipe with uri : %s", p.Profile, p.Uri),
			}
		}

		listenerExists, _ := config.Listeners.FindByUID(p.Listener)

		if !listenerExists {
			ok = false
			ret = &result{
				Severity: ERROR,
				Message:  fmt.Sprintf("missing listener (%s) for pipe with uri : %s", p.Listener, p.Uri),
			}
		}

		for _, e := range p.Endpoints {

			endpointExists, _ := config.Endpoints.FindByUID(e)

			if !endpointExists {
				ok = false
				ret = &result{
					Severity: ERROR,
					Message:  fmt.Sprintf("missing endpoint (%s) for pipe with uri : %s", e, p.Uri),
				}
			}
		}
	}

	config.Pipes.mapper(f)

	return ok, ret
}

func noEmptyUIDs(config *Configuration) (bool, *result) {

	ok := true
	var ret *result

	f := func(e Endpoint) {

		if e.UID == UID("") {
			ok = false
			ret = &result{
				Severity: WARNING,
				Message:  "empty uid in configuration",
			}

			return
		}
	}

	config.Listeners.mapper(f)
	config.Endpoints.mapper(f)

	return ok, ret
}
