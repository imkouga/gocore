package conf

import (
	"errors"
)

var (
	lineDelimiter = "\n"
	confDelimiter = "="
	confLineCount = 2
	defaultError  = errors.New(p_UNINITIALIZED)
)

const (
	p_SECTION_NOT_FOUND = "SECTION_MISMATCH"
	p_KEY_NOT_FOUND     = "KEY_MISMATCH"
	p_UNINITIALIZED     = "CONFIG_UNINITIALIZED"
)

type confError struct {
	reason string
	setion string
	key    string
}

func confErrorNew(section, key string) confError {
	ce := confError{}
	if 0 == len(section) {
		ce.key = key
		ce.reason = p_SECTION_NOT_FOUND
		return ce
	}

	ce.setion = section
	ce.reason = p_KEY_NOT_FOUND
	return ce
}

func (cf *confError) Error() string {
	if len(cf.setion) == 0 {
		return "[" + cf.setion + "] " + cf.reason
	}

	return "[" + cf.key + "] " + cf.reason
}
