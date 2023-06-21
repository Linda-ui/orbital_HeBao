package kitex_handler

import (
	"context"
	"reflect"
	"testing"

	echo "github.com/Linda-ui/orbital_HeBao/kitex_services/kitex_gen/echo"
)

func TestEchoImpl_EchoMethod(t *testing.T) {
	type args struct {
		ctx context.Context
		req *echo.EchoReq
	}
	tests := []struct {
		name     string
		s        *EchoImpl
		args     args
		wantResp *echo.EchoResp
		wantErr  bool
	}{
		{
			name: "Test EchoMethod with valid request",
			s:    &EchoImpl{},
			args: args{
				ctx: context.Background(),
				req: &echo.EchoReq{
					Msg: "Hello, world!",
				},
			},
			wantResp: &echo.EchoResp{
				Msg: "Hello, world!",
			},
			wantErr: false,
		},
		{
			name: "Test EchoMethod with empty request",
			s:    &EchoImpl{},
			args: args{
				ctx: context.Background(),
				req: &echo.EchoReq{},
			},
			wantResp: &echo.EchoResp{},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EchoImpl{}
			gotResp, err := s.EchoMethod(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("EchoImpl.EchoMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("EchoImpl.EchoMethod() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
