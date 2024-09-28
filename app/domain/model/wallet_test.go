package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWallet_Deposit(t *testing.T) {
	type fields struct {
		Balance Coins
	}
	type args struct {
		coins Coins
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Wallet
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "coins > 0 の場合 => 残高を加算できる",
			fields: fields{
				Balance: 0,
			},
			args: args{
				coins: 1,
			},
			want: &Wallet{
				Balance: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "coins = 0 の場合 => 何もしない",
			fields: fields{
				Balance: 0,
			},
			args: args{
				coins: 0,
			},
			want: &Wallet{
				Balance: 0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "coins < 0 の場合 => エラー",
			fields: fields{
				Balance: 0,
			},
			args: args{
				coins: -1,
			},
			want: &Wallet{
				Balance: 0,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wallet{
				Balance: tt.fields.Balance,
			}
			err := w.Deposit(tt.args.coins)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, w)
		})
	}
}

func TestWallet_Withdraw(t *testing.T) {
	type fields struct {
		Balance Coins
	}
	type args struct {
		coins Coins
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Wallet
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "coins > 0 && 残高 > coins の場合 => 残高を減算できる",
			fields: fields{
				Balance: 10,
			},
			args: args{
				coins: 1,
			},
			want: &Wallet{
				Balance: 9,
			},
			wantErr: assert.NoError,
		},
		{
			name: "coins > 0 && 残高 < coins の場合 => エラー",
			fields: fields{
				Balance: 10,
			},
			args: args{
				coins: 11,
			},
			want: &Wallet{
				Balance: 10,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrInsufficientCoins)
			},
		},
		{
			name: "coins = 0 の場合 => 何もしない",
			fields: fields{
				Balance: 10,
			},
			args: args{
				coins: 0,
			},
			want: &Wallet{
				Balance: 10,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wallet{
				Balance: tt.fields.Balance,
			}
			err := w.Withdraw(tt.args.coins)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, w)
		})
	}
}

func TestNewCoins(t *testing.T) {
	type args struct {
		v int
	}
	tests := []struct {
		name    string
		args    args
		want    Coins
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "v = 0 の場合 => 0",
			args: args{
				v: 0,
			},
			want:    0,
			wantErr: assert.NoError,
		},
		{
			name: "v > 0 の場合 => v",
			args: args{
				v: 1,
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "v < 0 の場合 => エラー",
			args: args{
				v: -1,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCoins(tt.args.v)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCoins_IsZero(t *testing.T) {
	tests := []struct {
		name string
		c    Coins
		want bool
	}{
		{
			name: "c = 0 の場合 => true",
			c:    0,
			want: true,
		},
		{
			name: "c > 0 の場合 => false",
			c:    1,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.IsZero())
		})
	}
}
