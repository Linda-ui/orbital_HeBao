package entity

type Err int

const (
	Err_BadRequest           Err = 10001
	Err_ServerNotFound       Err = 10002
	Err_ServerMethodNotFound Err = 10003
	Err_RequestServerFail    Err = 10004
	Err_ServerHandleFail     Err = 10005
	Err_ResponseUnableParse  Err = 10006
)

type ErrService interface {
	JSONEncode(e Err) map[string]interface{}
}

func (e Err) String() string {
	switch e {
	case Err_BadRequest:
		return "bad request"
	case Err_ServerNotFound:
		return "server not found"
	case Err_ServerMethodNotFound:
		return "server method not found"
	case Err_RequestServerFail:
		return "request server fail"
	case Err_ServerHandleFail:
		return "server handle fail"
	case Err_ResponseUnableParse:
		return "response unable parse"
	default:
		return "unknown error"
	}
}
