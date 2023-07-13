package errors

import (
	"reflect"
	"testing"

	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/entity"
)

func TestJSONEncode(t *testing.T) {
	type args struct {
		e entity.Err
	}

	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "10001 returns bad request error",
			args: args{10001},
			want: map[string]interface{}{
				"err_message": "bad request",
				"err_code":    10001,
			},
		},
		{
			name: "10006 returns response unable parse error",
			args: args{10006},
			want: map[string]interface{}{
				"err_message": "response unable parse",
				"err_code":    10006,
			},
		},
		{
			name: "unknown error",
			args: args{00000},
			want: map[string]interface{}{
				"err_message": "unknown error",
				"err_code":    00000,
			},
		},
	}

	testErrJSON := &errSender{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := testErrJSON.JSONEncode(test.args.e); !reflect.DeepEqual(got, test.want) {
				t.Errorf("return %v, want %v", got, test.want)
			}
		})
	}
}
