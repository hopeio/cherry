package mix

import (
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// SysErr ErrCode = -1
	Success            ErrCode = 0
	Canceled           ErrCode = 1
	Unknown            ErrCode = 2
	InvalidArgument    ErrCode = 3
	DeadlineExceeded   ErrCode = 4
	NotFound           ErrCode = 5
	AlreadyExists      ErrCode = 6
	PermissionDenied   ErrCode = 7
	ResourceExhausted  ErrCode = 8
	FailedPrecondition ErrCode = 9
	Aborted            ErrCode = 10
	OutOfRange         ErrCode = 11
	Unimplemented      ErrCode = 12
	Internal           ErrCode = 13
	Unavailable        ErrCode = 14
	DataLoss           ErrCode = 15
	Unauthenticated    ErrCode = 16
)

var codeMsgMap = map[ErrCode]string{
	Success:            "Success",
	Canceled:           "Canceled",
	Unknown:            "Unknown",
	InvalidArgument:    "InvalidArgument",
	DeadlineExceeded:   "DeadlineExceeded",
	NotFound:           "NotFound",
	AlreadyExists:      "AlreadyExists",
	PermissionDenied:   "PermissionDenied",
	ResourceExhausted:  "ResourceExhausted",
	FailedPrecondition: "FailedPrecondition",
	Aborted:            "Aborted",
	OutOfRange:         "OutOfRange",
	Unimplemented:      "Unimplemented",
	Internal:           "Internal",
	Unavailable:        "Unavailable",
	DataLoss:           "DataLoss",
	Unauthenticated:    "Unauthenticated",
}

// Not safe for concurrent use; call it during initialization.
func RegisterErrCode(code ErrCode, msg string) {
	codeMsgMap[code] = msg
}

func RegisterErrCodeMap(enum map[int32]string) {
	for code, msg := range enum {
		codeMsgMap[ErrCode(code)] = msg
	}
}

type ErrCode int32

// ErrCode layout: the low two digits hold the standard gRPC code, the high
// part holds the business code (bizcode*100+grpcCode, the same idea as the
// 2xx/3xx/4xx grouping of HTTP status codes, readable straight in decimal:
// 100107 = UserErrLogin+InvalidArgument).
// Define an error code with plain arithmetic:
// `ErrCode(bizCode)*100 + ErrCode(codes.X)`, and split it apart again with
// GRPCCode()/BizCode(). A pure gRPC error (business code 0) is its own value,
// so the constants above and the existing registration behaviour stay
// compatible.
// On the gRPC channel status.code only carries the standard enum below 100;
// the business code travels in status.msg (the composite value as a string)
// and in the ErrorInfo detail (reason/metadata).
const grpcCodeBase int32 = 100

// GRPCCode returns the standard gRPC code from the low two digits; for a
// pure gRPC error it is the error itself.
func (x ErrCode) GRPCCode() codes.Code {
	return codes.Code(int32(x) % grpcCodeBase)
}

// BizCode returns the business code from the high part; 0 for a pure gRPC
// error.
func (x ErrCode) BizCode() int32 {
	return int32(x) / grpcCodeBase
}

func (x ErrCode) String() string {
	if value, ok := codeMsgMap[x]; ok {
		return value
	}
	// composite: the business code was registered as a plain value
	// (RegisterErrCodeMap registers a proto enum)
	if biz := x.BizCode(); biz != 0 {
		if value, ok := codeMsgMap[ErrCode(biz)]; ok {
			return value
		}
	}
	return "Unknown Error, Code:" + strconv.Itoa(int(x))
}

func (x ErrCode) ErrResp() *ErrResp {
	return &ErrResp{Code: x, Msg: x.String()}
}

func (x ErrCode) Msg(msg string, data map[string]string) *ErrResp {
	return &ErrResp{Code: x, Msg: msg, Data: data}
}

func (x ErrCode) Wrap(err error) *ErrResp {
	if err == nil {
		return x.ErrResp()
	}
	errResp := ErrRespFrom(err)
	errResp.Code = x
	return errResp
}

func (x ErrCode) Error() string {
	return x.String()
}

// GRPCStatus splits the error for gRPC: the code uses the standard enum
// below 100 (so it stays valid) and the msg carries the composite value,
// which the client parses back into the business code.
func (x ErrCode) GRPCStatus() *status.Status {
	return status.New(x.GRPCCode(), strconv.Itoa(int(x)))
}

