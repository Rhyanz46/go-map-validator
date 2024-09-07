package map_validator

import (
	"errors"
	"fmt"
	"regexp"
)

func doSafeRegexpMustCompile(data string) (rex *regexp.Regexp, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Assign error value during panic
			err = errors.New(fmt.Sprintf("Error when compiling regex: %v", r))
		}
	}()

	rex = regexp.MustCompile(data)
	return rex, nil
}
