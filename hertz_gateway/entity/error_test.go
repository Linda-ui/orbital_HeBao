package entity

import "testing"

func TestErr_String(t *testing.T) {
	tests := []struct {
		name string
		e    Err
		want string
	}{
		{
			name: "bad request",
			e:    Err_BadRequest,
			want: "bad request",
		},
		{
			name: "default case",
			e:    10007,
			want: "unknown error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Err.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
