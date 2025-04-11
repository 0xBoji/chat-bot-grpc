package middleware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorInterceptor is a gRPC interceptor that handles errors in a consistent way
func ErrorInterceptor(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Recover from panics
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("Panic recovered in %s: %v\n%s", info.FullMethod, r, debug.Stack())
				status.Errorf(codes.Internal, "Internal server error")
			}
		}()

		// Call the handler
		resp, err := handler(ctx, req)

		// Log the error
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				// Log based on error code
				switch st.Code() {
				case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.FailedPrecondition:
					// Client errors - log at info level
					logger.Printf("Client error in %s: %v", info.FullMethod, err)
				case codes.Internal, codes.DataLoss, codes.Unknown:
					// Server errors - log at error level with stack trace
					logger.Printf("Server error in %s: %v\n%s", info.FullMethod, err, debug.Stack())
				default:
					// Other errors - log at warning level
					logger.Printf("Error in %s: %v", info.FullMethod, err)
				}
			} else {
				// Unknown error type - log at error level with stack trace
				logger.Printf("Unknown error in %s: %v\n%s", info.FullMethod, err, debug.Stack())
			}
		}

		return resp, err
	}
}

// StreamErrorInterceptor is a gRPC interceptor for streaming RPCs
func StreamErrorInterceptor(logger *log.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Recover from panics
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("Panic recovered in streaming %s: %v\n%s", info.FullMethod, r, debug.Stack())
				status.Errorf(codes.Internal, "Internal server error")
			}
		}()

		// Call the handler
		err := handler(srv, ss)

		// Log the error
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				// Log based on error code
				switch st.Code() {
				case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.FailedPrecondition:
					// Client errors - log at info level
					logger.Printf("Client error in streaming %s: %v", info.FullMethod, err)
				case codes.Internal, codes.DataLoss, codes.Unknown:
					// Server errors - log at error level with stack trace
					logger.Printf("Server error in streaming %s: %v\n%s", info.FullMethod, err, debug.Stack())
				default:
					// Other errors - log at warning level
					logger.Printf("Error in streaming %s: %v", info.FullMethod, err)
				}
			} else {
				// Unknown error type - log at error level with stack trace
				logger.Printf("Unknown error in streaming %s: %v\n%s", info.FullMethod, err, debug.Stack())
			}
		}

		return err
	}
}
