package errors

import "github.com/Linda-ui/orbital_HeBao/hertz_gateway/entity"

// errSender implements the entity.ErrService interface.
type errSender struct{}

func (ejson *errSender) JSONEncode(e entity.Err) map[string]interface{} {
	return map[string]interface{}{
		"err_message": e.String(),
		"err_code":    int(e),
	}
}

func New() entity.ErrService {
	return &errSender{}
}
