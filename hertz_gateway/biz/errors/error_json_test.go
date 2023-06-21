package errors

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		e Err
	}
	tests := []struct {
		name string
		args args
		want ErrJSON
	}{
		{
			name: "10001 returns bad request error",
			args: args{10001},
			want: ErrJSON{
				ErrCode:    10001,
				ErrMessage: "bad request",
			},
		},
		{
			name: "10002 returns server not found error",
			args: args{10002},
			want: ErrJSON{
				ErrCode:    10002,
				ErrMessage: "server not found",
			},
		},
		{
			name: "10003 returns server method not found error",
			args: args{10003},
			want: ErrJSON{
				ErrCode:    10003,
				ErrMessage: "server method not found",
			},
		},
		{
			name: "10004 returns request server fail error",
			args: args{10004},
			want: ErrJSON{
				ErrCode:    10004,
				ErrMessage: "request server fail",
			},
		},
		{
			name: "10005 returns server handle fail error",
			args: args{10005},
			want: ErrJSON{
				ErrCode:    10005,
				ErrMessage: "server handle fail",
			},
		},
		{
			name: "10006 returns response unable parse error",
			args: args{10006},
			want: ErrJSON{
				ErrCode:    10006,
				ErrMessage: "response unable parse",
			},
		},
		{
			name: "000 returns unknown error",
			args: args{000},
			want: ErrJSON{
				ErrCode:    000,
				ErrMessage: "unknown error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
