package task

import (
	"encoding/json"
	"testing"
)

func TestFetch(t *testing.T) {
	type args struct {
		t        *Client
		tType    string
		consumer string
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    *ApiResponse[T]
		wantErr bool
	}
	tests := []testCase[any]{
		{
			name: "Test 1",
			args: args{
				t: &Client{
					Url:   "https://tasks.k8test.bc2c.cn/",
					Token: "f2e76228-18e7-40b8-bf02-95c9c7679ef5",
				},
				tType:    "v_man",
				consumer: "consumer",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e := Fetch[map[string]any](tt.args.t, tt.args.tType, tt.args.consumer)
			if e != nil {
				println(e.Error())
			}
			b, _ := json.Marshal(got)
			println(string(b))

		})
	}
}
