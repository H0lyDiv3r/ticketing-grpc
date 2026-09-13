package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	res, err := handler(ctx, req)
	duration := time.Since(start)
	code := status.Code(err)

	log.Printf("method: %s, time: %v, err: %v, code: %v", info.FullMethod, duration, err, code)

	return res, err
}
