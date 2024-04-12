package user

import (
	context "context"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type UserServiceResolvers struct{ Service UserServiceServer }

func (s *UserServiceResolvers) UserServiceSignup(ctx context.Context, in *SignupReq) (*wrapperspb.StringValue, error) {
	return s.Service.Signup(ctx, in)
}

type UserInput = User
type SignupReqInput = SignupReq
