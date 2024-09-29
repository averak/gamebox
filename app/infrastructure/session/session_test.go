package session

import (
	"fmt"
	"testing"
	"time"

	"github.com/averak/gamebox/testutils/faker"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestSession_IsExpired(t *testing.T) {
	now := time.Now()

	type fields struct {
		IssuedAt  time.Time
		ExpiresAt time.Time
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "有効期限 < 現在時刻 => true",
			fields: fields{
				ExpiresAt: now.Add(-1 * time.Second),
			},
			args: args{
				now: now,
			},
			want: true,
		},
		{
			name: "有効期限 = 現在時刻 => true",
			fields: fields{
				ExpiresAt: now,
			},
			args: args{
				now: now,
			},
			want: true,
		},
		{
			name: "有効期限 > 現在時刻 => false",
			fields: fields{
				ExpiresAt: now.Add(1 * time.Second),
			},
			args: args{
				now: now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				IssuedAt:  tt.fields.IssuedAt,
				ExpiresAt: tt.fields.ExpiresAt,
			}
			assert.Equalf(t, tt.want, s.IsExpired(tt.args.now), "IsExpired(%v)", tt.args.now)
		})
	}
}

func TestEncodeSessionToken(t *testing.T) {
	now := time.Now()

	type args struct {
		session   Session
		secretKey []byte
	}
	tests := []struct {
		name string
		args args
		then func(t *testing.T, got string, err error)
	}{
		{
			name: "セッショントークンを JWT で生成できる",
			args: args{
				session: Session{
					PrincipalID: faker.UUIDv5("u1"),
					IssuedAt:    now,
					ExpiresAt:   now.Add(1 * time.Hour),
				},
				secretKey: []byte("00000000000000000000000000000000"),
			},
			then: func(t *testing.T, got string, err error) {
				assert.NoError(t, err)

				decoded, _ := DecodeSessionToken(got, []byte("00000000000000000000000000000000"), now)
				assert.Equal(t, faker.UUIDv5("u1"), decoded.PrincipalID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeSessionToken(tt.args.session, tt.args.secretKey)
			tt.then(t, got, err)
		})
	}
}

func TestDecodeSessionToken(t *testing.T) {
	now := time.Now()

	type args struct {
		token     string
		secretKey []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Session
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "署名が正しい => セッショントークンをデコードできる",
			args: args{
				token: func() string {
					token, err := EncodeSessionToken(Session{PrincipalID: faker.UUIDv5("u1"), IssuedAt: now, ExpiresAt: now.Add(1 * time.Hour)}, []byte("00000000000000000000000000000000"))
					if err != nil {
						t.Fatal(err)
					}
					return token
				}(),
				secretKey: []byte("00000000000000000000000000000000"),
			},
			want: Session{
				PrincipalID: faker.UUIDv5("u1"),
				IssuedAt:    now,
				ExpiresAt:   now.Add(1 * time.Hour),
			},
			wantErr: assert.NoError,
		},
		{
			name: "署名が異なる => エラー",
			args: args{
				token: func() string {
					token, err := EncodeSessionToken(Session{PrincipalID: faker.UUIDv5("u1"), IssuedAt: now, ExpiresAt: now.Add(1 * time.Hour)}, []byte("00000000000000000000000000000000"))
					if err != nil {
						t.Fatal(err)
					}
					return token
				}(),
				secretKey: []byte("11111111111111111111111111111111"),
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.ErrorIs(t, err, jwt.ErrSignatureInvalid, msgAndArgs...)
			},
		},
		{
			name: "有効期限切れ => エラー",
			args: args{
				token: func() string {
					token, err := EncodeSessionToken(Session{PrincipalID: faker.UUIDv5("u1"), IssuedAt: now, ExpiresAt: now.Add(-1 * time.Hour)}, []byte("00000000000000000000000000000000"))
					if err != nil {
						t.Fatal(err)
					}
					return token
				}(),
				secretKey: []byte("00000000000000000000000000000000"),
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.ErrorIs(t, err, jwt.ErrTokenExpired, msgAndArgs...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeSessionToken(tt.args.token, tt.args.secretKey, now)
			if !tt.wantErr(t, err, fmt.Sprintf("DecodeSessionToken(%v, %v)", tt.args.token, tt.args.secretKey)) {
				return
			}

			// JWT は秒単位で有効期限を扱うため、時間の差分が 1 秒以内であれば同じとみなす。
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateApproxTime(time.Second)); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
