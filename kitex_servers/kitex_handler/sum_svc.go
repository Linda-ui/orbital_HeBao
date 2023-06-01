package main

import (
	sum "Orbital_Hebao/kitex_servers/kitex_gen/sum"
	"context"
)

type SumImpl struct{}

func (s *SumImpl) SumMethod(ctx context.Context, req *sum.SumReq) (resp *sum.SumResp, err error) {
	return &sum.SumResp{Sum: req.FirstNum + req.SecondNum}, nil
}
