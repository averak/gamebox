package mdval

import (
	"net/http"

	"connectrpc.com/connect"
)

type (
	incomingHeaderKey  string
	outgoingHeaderKey  string
	outgoingTrailerKey string
	OutgoingHeaderMD   map[outgoingHeaderKey]string
	OutgoingTrailerMD  map[outgoingTrailerKey]string
)

const (
	IdempotencyKey         incomingHeaderKey = "x-idempotency-key"
	SessionTokenKey        incomingHeaderKey = "x-session-token"
	DebugAdjustedTimeKey   incomingHeaderKey = "x-debug-adjustment-timestamp"
	DebugSpoofingUserIDKey incomingHeaderKey = "x-debug-spoofing-uid"

	RespondTimestampKey outgoingHeaderKey = "x-respond-timestamp"
	ServerVersionKey    outgoingHeaderKey = "x-server-version"
)

type IncomingMD struct {
	origin http.Header
}

func NewIncomingMD(header http.Header) IncomingMD {
	return IncomingMD{
		origin: header,
	}
}

func (i IncomingMD) Get(key incomingHeaderKey) (string, bool) {
	v := i.origin.Get(string(key))
	return v, v != ""
}

func (i IncomingMD) Set(key incomingHeaderKey, value string) {
	i.origin.Set(string(key), value)
}

func (i IncomingMD) ToMap() map[incomingHeaderKey]string {
	res := make(map[incomingHeaderKey]string, len(i.origin))
	for k, v := range i.origin {
		if len(v) > 0 {
			res[incomingHeaderKey(k)] = v[0]
		}
	}
	return res
}

func SetOutgoingHeader(response connect.AnyResponse, md OutgoingHeaderMD) {
	for key, value := range md {
		response.Header().Set(string(key), value)
	}
}

func SetOutgoingTrailer(response connect.AnyResponse, md OutgoingTrailerMD) {
	for key, value := range md {
		response.Trailer().Set(string(key), value)
	}
}
