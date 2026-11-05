package mix_test

import (
	"strconv"
	"testing"

	"github.com/hopeio/mix"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// detail carries the i18n key (reason) and the placeholder variables
// (metadata) in a standard google.rpc.ErrorInfo: status.code is the
// standard enum in the low two digits (below 100) and status.msg is the
// composite value.
// The assertion goes through st.Proto().GetDetails(), the wire level that
// actually goes on the network, to guard against the Any-inside-Any bug:
// WithDetails wraps an anypb.Any in another one.
func TestErrRespGRPCStatusDetail(t *testing.T) {
	e := mix.NewErrResp(mix.InvalidArgument, "auth.err.thirdLogin", map[string]string{"type": "Apple", "x": "y"})
	st := e.GRPCStatus()
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("code = %v", st.Code())
	}
	if st.Message() != "3" {
		t.Fatalf("message = %q", st.Message())
	}
	wire := st.Proto().GetDetails()
	if len(wire) != 1 {
		t.Fatalf("details = %d", len(wire))
	}
	if wire[0].TypeUrl != "type.googleapis.com/google.rpc.ErrorInfo" {
		t.Fatalf("typeUrl = %s", wire[0].TypeUrl)
	}
	var got errdetails.ErrorInfo
	if err := proto.Unmarshal(wire[0].Value, &got); err != nil {
		t.Fatal(err)
	}
	if got.Reason != "auth.err.thirdLogin" {
		t.Fatalf("reason = %q", got.Reason)
	}
	if got.Metadata["type"] != "Apple" || got.Metadata["x"] != "y" {
		t.Fatalf("metadata = %v", got.Metadata)
	}
	// grpc-go's Details() unwraps the standard types on its own.
	if d, ok := st.Details()[0].(*errdetails.ErrorInfo); !ok || d.Metadata["type"] != "Apple" {
		t.Fatalf("decoded detail = %T", st.Details()[0])
	}
}

func TestErrRespFromGRPCStatusDetails(t *testing.T) {
	src := mix.NewErrResp(mix.InvalidArgument, "auth.err.thirdLogin", map[string]string{"type": "Apple"})
	got := mix.ErrRespFrom(src.GRPCStatus().Err())
	if got.Code != mix.InvalidArgument || got.Msg != "auth.err.thirdLogin" || got.Data["type"] != "Apple" {
		t.Fatalf("got=%+v", got)
	}
}

// A composite business code travels over the gRPC channel: the code
// splits off the low two digits (below 100), the msg carries the whole
// value and the reason carries the i18n key.
func TestErrRespCompositeRoundtrip(t *testing.T) {
	code := mix.ErrCode(1001)*100 + mix.ErrCode(codes.PermissionDenied)
	e := mix.NewErrResp(code, "auth.err.notActivated", map[string]string{"a": "b"})
	st := e.GRPCStatus()
	if st.Code() != codes.PermissionDenied {
		t.Fatalf("code = %v", st.Code())
	}
	if st.Message() != strconv.Itoa(int(code)) {
		t.Fatalf("message = %q", st.Message())
	}
	got := mix.ErrRespFrom(st.Err())
	if got.Code != code || got.Msg != "auth.err.notActivated" || got.Data["a"] != "b" {
		t.Fatalf("got=%+v", got)
	}
}

// With both msg and data empty there are no details, to stay light.
func TestErrRespGRPCStatusNoDetail(t *testing.T) {
	st := mix.NewErrResp(mix.Internal, "", nil).GRPCStatus()
	if len(st.Details()) != 0 {
		t.Fatalf("details = %v", st.Details())
	}
}

// The data is empty but the msg is not, so the reason (the i18n key)
// still has to reach the client.
func TestErrRespGRPCStatusReasonWithoutData(t *testing.T) {
	st := mix.InvalidArgument.Msg("auth.err.notActivated", nil).GRPCStatus()
	details := st.Details()
	if len(details) != 1 {
		t.Fatalf("details = %v", details)
	}
	if d, ok := details[0].(*errdetails.ErrorInfo); !ok || d.Reason != "auth.err.notActivated" {
		t.Fatalf("detail = %+v", details[0])
	}
	if status.Code(mix.InvalidArgument.Msg("x", nil).GRPCStatus().Err()) != codes.InvalidArgument {
		t.Fatal("code")
	}
}
