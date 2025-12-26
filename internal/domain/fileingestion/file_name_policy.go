package fileingestion

import (
	"errors"
	"regexp"
)

type RegexFileNamePolicy struct {
	regex *regexp.Regexp
}

func NewRegexFileNamePolicy(regex *regexp.Regexp) *RegexFileNamePolicy {
	return &RegexFileNamePolicy{regex: regex}
}

func (p *RegexFileNamePolicy) Match(name string) (bool, error) {
	if p.regex == nil {
		return false, errors.New("regex not initialized")
	}
	return p.regex.MatchString(name), nil
}
