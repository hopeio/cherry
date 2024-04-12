package user

import (
	fmt "fmt"
	graphql "github.com/99designs/gqlgen/graphql"
	io "io"
)

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
