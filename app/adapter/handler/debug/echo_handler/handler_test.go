package echo_handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/averak/gamebox/app/registry"
	"github.com/averak/gamebox/protobuf/api/debug"
	"github.com/averak/gamebox/protobuf/api/debug/debugconnect"
	"github.com/averak/gamebox/testutils"
	"github.com/averak/gamebox/testutils/bdd"
	"github.com/averak/gamebox/testutils/testconnect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_handler_EchoV1(t *testing.T) {
	mux, err := registry.InitializeAPIServerMux(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	now := time.Now()

	type given struct{}
	type when struct {
		req *debug.EchoServiceEchoV1Request
	}
	type then = func(*testing.T, *connect.Response[debug.EchoServiceEchoV1Response], error)
	tests := []bdd.Testcase[given, when, then]{
		{
			Behaviors: []bdd.Behavior[when, then]{
				{
					When: when{
						req: &debug.EchoServiceEchoV1Request{
							Message: "echo",
						},
					},
					Then: func(t *testing.T, got *connect.Response[debug.EchoServiceEchoV1Response], err error) {
						require.NoError(t, err)

						want := &debug.EchoServiceEchoV1Response{
							Message:   "echo",
							Timestamp: timestamppb.New(now),
						}
						assert.EqualExportedValues(t, want, got.Msg)
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.Run(t, func(t *testing.T, given given, when when, then then) {
			defer testutils.Teardown(t)

			got, err := testconnect.MethodInvoke(
				debugconnect.NewEchoServiceClient(http.DefaultClient, server.URL).EchoV1,
				when.req,
				testconnect.WithAdjustedTime(now),
			)
			then(t, got, err)
		})
	}
}
