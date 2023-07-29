package handler

import (
	"context"
	"errors"
	"math"

	sum "github.com/Linda-ui/orbital_HeBao/kitex_services/sum/kitex_gen/sum"
)

type SumImpl struct{}

func (s *SumImpl) SumMethod(sctx context.Context, req *sum.SumReq) (resp *sum.SumResp, err error) {
	if (req.FirstNum > 0 && req.SecondNum > math.MaxInt64-req.FirstNum) ||
		(req.FirstNum < 0 && req.SecondNum < math.MinInt64-req.FirstNum) {
		return nil, errors.New("Overflow: sum is too large to be represented.")
	}
	return &sum.SumResp{Sum: req.FirstNum + req.SecondNum}, nil
}
