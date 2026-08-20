package wecom

import "fmt"

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
	e, ok := err.(Error)
	if !ok {
		return false
	}
	switch e.Code {
	case 40014, 42001, 40001:
		return true
	default:
		return false
	}
}
