package mdval

import (
	"net/http"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncomingMD_Get(t *testing.T) {
	type fields struct {
		origin http.Header
	}
	type args struct {
		key incomingHeaderKey
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "キーが存在する場合 => 値を取得できる",
			fields: fields{
				origin: http.Header{
					textproto.CanonicalMIMEHeaderKey("K1"): []string{"v1"},
				},
			},
			args: args{
				key: "K1",
			},
			want:  "v1",
			want1: true,
		},
		{
			name: "キーが存在しない場合 => 空文字を返す",
			fields: fields{
				origin: http.Header{
					textproto.CanonicalMIMEHeaderKey("K1"): []string{"v1"},
				},
			},
			args: args{
				key: "K2",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IncomingMD{
				origin: tt.fields.origin,
			}
			got, got1 := i.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestIncomingMD_Set(t *testing.T) {
	type fields struct {
		origin http.Header
	}
	type args struct {
		key   incomingHeaderKey
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		then   func(t *testing.T, i IncomingMD)
	}{
		{
			name: "キーが存在しない場合 => 値が設定される",
			fields: fields{
				origin: http.Header{
					textproto.CanonicalMIMEHeaderKey("K1"): []string{"v1"},
				},
			},
			args: args{
				key:   "K2",
				value: "v2",
			},
			then: func(t *testing.T, i IncomingMD) {
				v, ok := i.Get("K2")
				assert.True(t, ok)
				assert.Equal(t, "v2", v)
			},
		},
		{
			name: "キーが存在する場合 => 値が上書きされる",
			fields: fields{
				origin: http.Header{
					textproto.CanonicalMIMEHeaderKey("K1"): []string{"v1"},
				},
			},
			args: args{
				key:   "K1",
				value: "v2",
			},
			then: func(t *testing.T, i IncomingMD) {
				v, ok := i.Get("K1")
				assert.True(t, ok)
				assert.Equal(t, "v2", v)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IncomingMD{
				origin: tt.fields.origin,
			}
			i.Set(tt.args.key, tt.args.value)
			tt.then(t, i)
		})
	}
}

func TestIncomingMD_ToMap(t *testing.T) {
	type fields struct {
		origin http.Header
	}
	tests := []struct {
		name   string
		fields fields
		want   map[incomingHeaderKey]string
	}{
		{
			name: "ヘッダーが空の場合 => 空のマップを返す",
			fields: fields{
				origin: http.Header{},
			},
			want: map[incomingHeaderKey]string{},
		},
		{
			name: "ヘッダーが存在する場合 => マップを返す",
			fields: fields{
				origin: http.Header{
					textproto.CanonicalMIMEHeaderKey("K1"): []string{"v1"},
					textproto.CanonicalMIMEHeaderKey("K2"): []string{"v2"},
				},
			},
			want: map[incomingHeaderKey]string{
				"K1": "v1",
				"K2": "v2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IncomingMD{
				origin: tt.fields.origin,
			}
			assert.Equalf(t, tt.want, i.ToMap(), "ToMap()")
		})
	}
}
