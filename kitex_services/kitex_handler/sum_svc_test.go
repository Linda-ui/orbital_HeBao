package kitex_handler

import (
	"context"
	"reflect"
	"testing"

	sum "github.com/Linda-ui/orbital_HeBao/kitex_services/kitex_gen/sum"
)

func TestSumImpl_SumMethod(t *testing.T) {
	type args struct {
		ctx context.Context
		req *sum.SumReq
	}
	tests := []struct {
		name     string
		s        *SumImpl
		args     args
		wantResp *sum.SumResp
		wantErr  bool
	}{
		{
			name: "Test SumMethod with valid request",
			s:    &SumImpl{},
			args: args{
				ctx: context.Background(),
				req: &sum.SumReq{
					FirstNum:  2,
					SecondNum: 3,
				},
			},
			wantResp: &sum.SumResp{
				Sum: 5,
			},
			wantErr: false,
		},
		{
			name: "Test SumMethod with negative numbers",
			s:    &SumImpl{},
			args: args{
				ctx: context.Background(),
				req: &sum.SumReq{
					FirstNum:  -5,
					SecondNum: 3,
				},
			},
			wantResp: &sum.SumResp{
				Sum: -2,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SumImpl{}
			gotResp, err := s.SumMethod(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SumImpl.SumMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("SumImpl.SumMethod() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
