package advice

import (
	"net/http"
	"net/textproto"
	"testing"
	"time"

	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/infrastructure/connect/mdval"
	"github.com/averak/gamebox/testutils/bdd"
	"github.com/averak/gamebox/testutils/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_fixIdempotencyKey(t *testing.T) {
	type args struct {
		incomingMD mdval.IncomingMD
	}
	tests := []struct {
		name string
		args args
		then func(*testing.T, *preconditionParamsBuilder)
	}{
		{
			name: "ヘッダが未指定の場合 => UUID を生成する",
			args: args{
				incomingMD: mdval.NewIncomingMD(http.Header{}),
			},
			then: func(t *testing.T, builder *preconditionParamsBuilder) {
				assert.NotEmpty(t, builder.raw.IdempotencyKey)
			},
		},
		{
			name: "ヘッダが指定されている場合 => ヘッダの値を採用する",
			args: args{
				incomingMD: mdval.NewIncomingMD(http.Header{
					textproto.CanonicalMIMEHeaderKey(string(mdval.IdempotencyKey)): {faker.UUIDv5("i1").String()},
				}),
			},
			then: func(t *testing.T, builder *preconditionParamsBuilder) {
				assert.Equal(t, faker.UUIDv5("i1"), builder.raw.IdempotencyKey)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := newPreconditionParamsBuilder()
			fixIdempotencyKey(tt.args.incomingMD, builder)
			tt.then(t, builder)
		})
	}
}

func Test_fixTimestamp(t *testing.T) {
	now := time.Now()

	type given struct {
		conf *config.Config
	}
	type when struct {
		incomingMD mdval.IncomingMD
		now        time.Time
	}
	type then func(*testing.T, *preconditionParamsBuilder, error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Name: "デバッグモードでない場合",
			Given: given{
				conf: &config.Config{Debug: false},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "現在時刻をそのまま採用する",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							// 無視される
							textproto.CanonicalMIMEHeaderKey(string(mdval.DebugAdjustedTimeKey)): {"2000-01-01T00:00:00Z"},
						}),
					},
					Then: func(t *testing.T, builder *preconditionParamsBuilder, err error) {
						require.NoError(t, err)
						assert.Equal(t, now, builder.raw.RequestedTime)
						assert.Equal(t, now, builder.raw.Now)
					},
				},
			},
		},
		{
			Name: "デバッグモードの場合",
			Given: given{
				conf: &config.Config{Debug: true},
			},
			Behaviors: []bdd.Behavior[when, then]{
				{
					Name: "デバッグ用ヘッダが指定されていない場合 => 現在時刻を採用する",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{}),
					},
					Then: func(t *testing.T, builder *preconditionParamsBuilder, err error) {
						require.NoError(t, err)
						assert.Equal(t, now, builder.raw.RequestedTime)
						assert.Equal(t, now, builder.raw.Now)
					},
				},
				{
					Name: "デバッグ用ヘッダが指定された場合 => デバッグ用ヘッダの値を採用する",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							textproto.CanonicalMIMEHeaderKey(string(mdval.DebugAdjustedTimeKey)): {"2000-01-01T00:00:00Z"},
						}),
					},
					Then: func(t *testing.T, builder *preconditionParamsBuilder, err error) {
						require.NoError(t, err)
						assert.Equal(t, now, builder.raw.RequestedTime)
						assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), builder.raw.Now)
					},
				},
				{
					Name: "デバッグ用ヘッダが不正なフォーマットの場合 => エラー",
					When: when{
						incomingMD: mdval.NewIncomingMD(http.Header{
							textproto.CanonicalMIMEHeaderKey(string(mdval.DebugAdjustedTimeKey)): {"invalid format"},
						}),
					},
					Then: func(t *testing.T, builder *preconditionParamsBuilder, err error) {
						require.Error(t, err)
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			builder := newPreconditionParamsBuilder()
			err := fixTimestamp(given.conf, when.incomingMD, builder, now)
			then(t, builder, err)
		})
	}
}

func Test_preconditionParamsBuilder_build(t *testing.T) {
	now := time.Now()

	type fields struct {
		raw preconditionParams
	}
	type args struct {
		requireMasterVersion bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    preconditionParams
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "全てのフィールドが設定されている => 構築に成功する",
			fields: fields{
				raw: preconditionParams{
					RequestedTime:  now,
					Now:            now,
					RequestID:      faker.UUIDv5("RequestID"),
					IdempotencyKey: faker.UUIDv5("IdempotencyKey"),
				},
			},
			args: args{
				requireMasterVersion: true,
			},
			want: preconditionParams{
				RequestID:      faker.UUIDv5("RequestID"),
				IdempotencyKey: faker.UUIDv5("IdempotencyKey"),
				RequestedTime:  now,
				Now:            now,
			},
			wantErr: assert.NoError,
		},
		{
			name: "RequestedTime が未設定 => エラー",
			fields: fields{
				raw: preconditionParams{
					RequestID:      faker.UUIDv5("RequestID"),
					IdempotencyKey: faker.UUIDv5("IdempotencyKey"),
					Now:            now,
				},
			},
			args: args{
				requireMasterVersion: true,
			},
			want:    preconditionParams{},
			wantErr: assert.Error,
		},
		{
			name: "Now が未設定 => エラー",
			fields: fields{
				raw: preconditionParams{
					RequestID:      faker.UUIDv5("RequestID"),
					IdempotencyKey: faker.UUIDv5("IdempotencyKey"),
					RequestedTime:  now,
				},
			},
			args: args{
				requireMasterVersion: true,
			},
			want:    preconditionParams{},
			wantErr: assert.Error,
		},
		{
			name: "RequestID が未設定 => エラー",
			fields: fields{
				raw: preconditionParams{
					IdempotencyKey: faker.UUIDv5("IdempotencyKey"),
					RequestedTime:  now,
					Now:            now,
				},
			},
			args: args{
				requireMasterVersion: true,
			},
			want:    preconditionParams{},
			wantErr: assert.Error,
		},
		{
			name: "IdempotencyKey が未設定 => エラー",
			fields: fields{
				raw: preconditionParams{
					RequestID:     faker.UUIDv5("RequestID"),
					RequestedTime: now,
					Now:           now,
				},
			},
			args: args{
				requireMasterVersion: true,
			},
			want:    preconditionParams{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &preconditionParamsBuilder{
				raw: tt.fields.raw,
			}
			got, err := b.build()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
