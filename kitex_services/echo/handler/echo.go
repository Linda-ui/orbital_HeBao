package handler

import (
	"context"

	echo "github.com/Linda-ui/orbital_HeBao/kitex_services/echo/kitex_gen/echo"
)

type EchoImpl struct{}

func (s *EchoImpl) EchoMethod(ctx context.Context, req *echo.EchoReq) (resp *echo.EchoResp, err error) {
	return &echo.EchoResp{Msg: req.Msg}, nil
}
