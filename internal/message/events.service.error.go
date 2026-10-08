package msgerr

import "errors"

var (
	CapacityFulled = errors.New("event capacity fulled")
)
