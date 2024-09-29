package advice

import (
	"fmt"
	"testing"

	"github.com/averak/gamebox/app/domain/model"
	"github.com/stretchr/testify/assert"
)

func Test_checkUserStatus(t *testing.T) {
	type args struct {
		user model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ユーザがアクティブな場合 => エラーなし",
			args: args{
				user: model.User{
					Status: model.UserStatusActive,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ユーザがアクティブでない場合 => エラーを返す",
			args: args{
				user: model.User{
					Status: model.UserStatusPending,
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, checkUserStatus(tt.args.user), fmt.Sprintf("checkUserStatus(%v)", tt.args.user))
		})
	}
}
