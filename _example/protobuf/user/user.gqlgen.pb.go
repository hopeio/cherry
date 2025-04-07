package user

import (
	context "context"
	fmt "fmt"
	graphql "github.com/99designs/gqlgen/graphql"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	io "io"
)

type UserServiceResolvers struct{ Service UserServiceServer }

func (s *UserServiceResolvers) UserServiceSignup(ctx context.Context, in *SignupReq) (*wrapperspb.StringValue, error) {
	return s.Service.Signup(ctx, in)
}
func (s *UserServiceResolvers) UserServiceGetUser(ctx context.Context, in *GetUserReq) (*User, error) {
	return s.Service.GetUser(ctx, in)
}

type UserInput = User
type SignupReqInput = SignupReq
type GetUserReqInput = GetUserReq

func MarshalGender(x Gender) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%q", x.String())
	})
}

func UnmarshalGender(v interface{}) (Gender, error) {
	code, ok := v.(string)
	if ok {
		return Gender(Gender_value[code]), nil
	}
	return 0, fmt.Errorf("cannot unmarshal Gender enum")
}

func MarshalRole(x Role) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%q", x.String())
	})
}

func UnmarshalRole(v interface{}) (Role, error) {
	code, ok := v.(string)
	if ok {
		return Role(Role_value[code]), nil
	}
	return 0, fmt.Errorf("cannot unmarshal Role enum")
}

func MarshalUserStatus(x UserStatus) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%q", x.String())
	})
}

func UnmarshalUserStatus(v interface{}) (UserStatus, error) {
	code, ok := v.(string)
	if ok {
		return UserStatus(UserStatus_value[code]), nil
	}
	return 0, fmt.Errorf("cannot unmarshal UserStatus enum")
}

func MarshalUserErr(x UserErr) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%q", x.String())
	})
}

func UnmarshalUserErr(v interface{}) (UserErr, error) {
	code, ok := v.(string)
	if ok {
		return UserErr(UserErr_value[code]), nil
	}
	return 0, fmt.Errorf("cannot unmarshal UserErr enum")
}
