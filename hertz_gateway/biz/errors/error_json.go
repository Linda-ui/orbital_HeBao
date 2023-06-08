package errors

type ErrJSON struct {
	ErrCode    int64  `json:"err_code"`
	ErrMessage string `json:"err_message"`
}

func New(e Err) ErrJSON {
	return ErrJSON{
		ErrCode:    int64(e),
		ErrMessage: e.String(),
	}
}
