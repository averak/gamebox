package advice

import (
	"errors"
	"reflect"
	"testing"

	"github.com/averak/gamebox/app/domain/model"
	"github.com/averak/gamebox/protobuf/custom_option"
	"github.com/averak/gamebox/testutils/faker"
)

func TestMethodInfo_FindErrorDefinition(t *testing.T) {
	var (
		err1 = errors.New("error1")
		err2 = errors.New("error2")
	)

	type fields struct {
		errCauses map[error]*MethodErrDefinition
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *MethodErrDefinition
		want1  bool
	}{
		{
			name: "エラー定義が存在する => true",
			fields: fields{
				errCauses: map[error]*MethodErrDefinition{
					err1: {
						Code: 1,
					},
				},
			},
			args: args{
				err: err1,
			},
			want: &custom_option.MethodErrorDefinition{
				Code: 1,
			},
			want1: true,
		},
		{
			name: "エラー定義が存在しない => false",
			fields: fields{
				errCauses: map[error]*MethodErrDefinition{
					err1: {
						Code: 1,
					},
				},
			},
			args: args{
				err: err2,
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MethodInfo{
				errCauses: tt.fields.errCauses,
			}
			got, got1 := m.FindErrorDefinition(tt.args.err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindErrorDefinition() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindErrorDefinition() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRequest_Principal(t *testing.T) {
	type testCase[T any] struct {
		name  string
		r     Request[T]
		want  model.User
		want1 bool
	}
	tests := []testCase[any]{
		{
			name: "principal == nil の場合 => false",
			r: Request[any]{
				principal: nil,
			},
			want:  model.User{},
			want1: false,
		},
		{
			name: "principal != nil の場合 => true",
			r: Request[any]{
				principal: &model.User{
					ID: faker.UUIDv5("u1"),
				},
			},
			want: model.User{
				ID: faker.UUIDv5("u1"),
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.r.Principal()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Principal() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Principal() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
