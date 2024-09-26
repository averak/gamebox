package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWallet_deposit(t *testing.T) {
	type fields struct {
		Balance int
	}
	type args struct {
		amount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Wallet
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "amount = 0 の場合 => 何もしない",
			fields: fields{
				Balance: 10,
			},
			args: args{
				amount: 0,
			},
			want: &Wallet{
				Balance: 10,
			},
			wantErr: assert.NoError,
		},
		{
			name: "amount > 0 の場合 => 残高を加算できる",
			fields: fields{
				Balance: 10,
			},
			args: args{
				amount: 1,
			},
			want: &Wallet{
				Balance: 11,
			},
			wantErr: assert.NoError,
		},
		{
			name: "amount < 0 の場合 => エラー",
			fields: fields{
				Balance: 10,
			},
			args: args{
				amount: -1,
			},
			want: &Wallet{
				Balance: 10,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wallet{
				Balance: tt.fields.Balance,
			}
			err := w.deposit(tt.args.amount)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, w)
		})
	}
}

func TestWallet_withdraw(t *testing.T) {
	type fields struct {
		Balance int
	}
	type args struct {
		amount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Wallet
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "amount = 0 の場合 => 何もしない",
			fields: fields{
				Balance: 10,
			},
			args: args{
				amount: 0,
			},
			want: &Wallet{
				Balance: 10,
			},
			wantErr: assert.NoError,
		},
		{
			name: "amount > 0 && 残高 > amount の場合 => 残高を減算できる",
			fields: fields{
				Balance: 10,
			},
			args: args{
				amount: 1,
			},
			want: &Wallet{
				Balance: 9,
			},
			wantErr: assert.NoError,
		},
		{
			name: "amount > 0 && 残高 < amount の場合 => エラー",
			fields: fields{
				Balance: 10,
			},
			args: args{
				amount: 11,
			},
			want: &Wallet{
				Balance: 10,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrInsufficientBalance)
			},
		},
		{
			name: "amount < 0 の場合 => エラー",
			fields: fields{
				Balance: 10,
			},
			args: args{
				amount: -1,
			},
			want: &Wallet{
				Balance: 10,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wallet{
				Balance: tt.fields.Balance,
			}
			if !tt.wantErr(t, w.withdraw(tt.args.amount)) {
				return
			}
			assert.Equal(t, tt.want, w)
		})
	}
}
