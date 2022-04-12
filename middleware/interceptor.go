package middleware

import (
	"context"
	"net/http"

	db "github.com/sanches1984/gopkg-pg-orm"

	"google.golang.org/grpc"
)

// NewDBServerInterceptor wrap endpoint with middleware mixing in db connection
func NewDBServerInterceptor(dbClient db.IClient, option ...db.Option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(db.NewContext(ctx, dbClient.WrapWithContext(ctx), option...), req)
	}
}

// NewDBServerMiddleware wrap endpoint with middleware mixing in db connection
func NewDBServerMiddleware(dbClient db.IClient, option ...db.Option) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			r = r.WithContext(db.NewContext(ctx, dbClient.WrapWithContext(ctx), option...))
			next.ServeHTTP(w, r)
		})
	}
}
