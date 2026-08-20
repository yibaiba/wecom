package wecom

import (
	"errors"
	"fmt"
)

// Error is a WeCom API errcode response.
type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	if e.Msg == "" {
		return fmt.Sprintf("wecom api error %d", e.Code)
	}
	return fmt.Sprintf("wecom api %d: %s", e.Code, e.Msg)
}

func isTokenErr(err error) bool {
	var e Error
	if !errors.As(err, &e) {
		return false
	}
	switch e.Code {
	case 40014, 42001, 40001:
		return true
	default:
		return false
	}
}

// 81013: userid not in the app visible range. 60011: no privilege.
func isDirectoryScopeErr(err error) bool {
	var e Error
	if !errors.As(err, &e) {
		return false
	}
	return e.Code == 81013 || e.Code == 60011
}
