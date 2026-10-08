package msgerr

import "errors"

var (
	AlredyJoinedEvents = errors.New("already joined event")
	AlredySavedEvent   = errors.New("already saved event")
	NotJoinedEvents    = errors.New("not joined this event")
	NotSavedEvent      = errors.New("not saved this event")
)
