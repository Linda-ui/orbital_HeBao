package main

import (
	echo "Orbital_Hebao/kitex_servers/kitex_gen/echo"
	"context"
)

type EchoImpl struct{}

func (s *EchoImpl) EchoMethod(ctx context.Context, req *echo.EchoReq) (resp *echo.EchoResp, err error) {
	return &echo.EchoResp{Msg: req.Msg}, nil
}
