package kitex_handler

import (
	"context"

	sum "github.com/Linda-ui/orbital_HeBao/kitex_services/kitex_gen/sum"
)

type SumImpl struct{}

func (s *SumImpl) SumMethod(ctx context.Context, req *sum.SumReq) (resp *sum.SumResp, err error) {
	return &sum.SumResp{Sum: req.FirstNum + req.SecondNum}, nil
}
