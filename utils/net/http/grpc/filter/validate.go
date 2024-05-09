package filter

import (
	"context"

	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/utils/validation/validator"
	"google.golang.org/grpc"
)

func validate(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	if err := validator.Validator.Struct(req); err != nil {
		return nil, errorcode.InvalidArgument.Message(validator.Trans(err))
	}

	return handler(ctx, req)
}
