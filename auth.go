package aviation

import (
	"context"

	"github.com/evergreen-ci/gimlet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/code"
	"google.golang.org/grpc/metadata"
)

// TODO: add interceptors to provide "role required" and "group
// member" using gimlet.Authenticator

func MakeAuthenticationRequiredUnaryInterceptor(um gimlet.UserManager, conf gimlet.UserMiddlewareConfiguration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return nil, grpc.Errorf(code.Unauthenticated, "missing metadata from context")
		}

		var (
			authDataAPIKey string
			authDataName   string
		)

		// Grab API auth details from header
		if len(meta[conf.HeaderKeyName]) > 0 {
			authDataAPIKey = meta[conf.HeaderKeyName][0]
		}
		if len(meta[conf.HeaderUserName]) > 0 {
			authDataName = meta[conf.HeaderUserName][0]
		}

		if len(authDataAPIKey) == 0 {
			return nil, grpc.Errorf(code.Unauthenticated, "user key not provided")
		}

		usr, err := um.GetUserByID(authDataName)
		if err != nil {
			return nil, grpc.Errorf(code.Unauthenticated, "problem finding user: %+v", err)
		}

		if usr == nil {
			return nil, grpc.Errorf(code.Unauthenticated, "user not found")
		}

		if usr.GetAPIKey() != authDataAPIKey {
			return nil, grpc.Errorf(code.Unauthenticated, "incorrect credentials")
		}

		ctx = SetRequestUser(ctx, usr)

		return handler(ctx, req)
	}
}

func MakeAuthenticationRequiredStreamingInterceptor(um gimlet.UserManager, conf gimlet.UserMiddlewareConfiguration) grpc.StreamingServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ctx := stream.Context()

		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return nil, grpc.Errorf(code.Unauthenticated, "missing metadata from context")
		}

		var (
			authDataAPIKey string
			authDataName   string
		)

		// Grab API auth details from header
		if len(meta[conf.HeaderKeyName]) > 0 {
			authDataAPIKey = meta[conf.HeaderKeyName][0]
		}
		if len(meta[conf.HeaderUserName]) > 0 {
			authDataName = meta[conf.HeaderUserName][0]
		}

		if len(authDataAPIKey) == 0 {
			return nil, grpc.Errorf(code.Unauthenticated, "user key not provided")
		}

		usr, err := um.GetUserByID(authDataName)
		if err != nil {
			return nil, grpc.Errorf(code.Unauthenticated, "problem finding user: %+v", err)
		}

		if usr == nil {
			return nil, grpc.Errorf(code.Unauthenticated, "user not found")
		}

		if usr.GetAPIKey() != authDataAPIKey {
			return nil, grpc.Errorf(code.Unauthenticated, "incorrect credentials")
		}

		return handler(ctx, req)
	}
}
