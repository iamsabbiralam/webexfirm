package logging

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	correlationIDKey = "correlation-id"
)

type (
	correlationKey struct{}
	ctxLogger      struct{}
)

// WithLogger loads a logger with correlation id into the context.  This logger will
// be retrieved by any FromContext call using the returned context.
func WithLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	cid := CorrelationIDFromContext(ctx)
	if cid == "" {
		cid = generateCorrelationID()
		ctx = context.WithValue(ctx, correlationKey{}, cid)
	}
	return context.WithValue(ctx, ctxLogger{}, logger.WithField(correlationIDKey, cid))
}

// FromContext will check the context for a correlation id, if the correlation id
// isn't present then it will generate one, the correlation id will be a field
// in the logger.
func FromContext(ctx context.Context) *logrus.Entry {
	if lgr := ctx.Value(ctxLogger{}); lgr != nil {
		// Value can only be loaded by ContextLogger call
		return lgr.(*logrus.Entry)
	}
	return NewLogger().WithField(correlationIDKey, CorrelationIDFromContext(ctx))
}

// WithContext takes the context and appends the correlation id on the existing logger
func WithContext(ctx context.Context, logger logrus.FieldLogger) *logrus.Entry {
	cid := CorrelationIDFromContext(ctx)
	return logger.WithField(correlationIDKey, cid)
}

// CorrelationIdFromContext will return the correlation id from the context
// metadata if present, otherwise will return an empty string
func CorrelationIDFromContext(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		correlationIDs := md[correlationIDKey]
		if len(correlationIDs) > 0 {
			return correlationIDs[0]
		}
	}
	if cid := ctx.Value(correlationKey{}); cid != nil {
		return cid.(string)
	}
	return ""
}

// ContextWithCorrelationID will generate a correlation id and add it to the
// contexts metadata if one doesn't already exist
func ContextWithCorrelationID(ctx context.Context) context.Context {
	cid := CorrelationIDFromContext(ctx)
	if cid == "" {
		// If no correlation id is present on the context generate one and attach it
		// to the contexts metadata
		cid := generateCorrelationID()
		// Check to see if there is any metadata on the context
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md[correlationIDKey] = []string{cid}
		} else {
			md = metadata.Pairs(correlationIDKey, cid)
		}
		// Build a new context with the correlation id on the incoming metadata
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, correlationIDKey, cid)
	return ctx
}

const logLead = "x-log-"

func UnaryClientInterceptor(ignore ...string) grpc.UnaryClientInterceptor {
	ig := make(map[string]bool, len(ignore))
	for _, v := range ignore {
		ig[v] = true
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		lg := FromContext(ctx)
		md := metautils.ExtractOutgoing(ctx)
		for k, v := range lg.Data {
			if ig[k] {
				continue
			}
			md.Add(logLead+k, fmt.Sprint(v))
		}
		if md.Get(correlationIDKey) == "" {
			md.Set(correlationIDKey, CorrelationIDFromContext(ctx))
		}
		ctx = md.ToOutgoing(ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// UnaryServerInterceptor will make sure a correlation id exists on the context
func UnaryServerInterceptor(log *logrus.Entry, clear bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = ContextWithCorrelationID(ctx)
		ctx = WithLogger(ctx, log)
		ctx = forwardLogger(ctx, clear)
		return handler(ctx, req)
	}
}

// StreamServerInterceptor will make sure a correlation id exists on the context
func StreamServerInterceptor(log *logrus.Entry, clear bool) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := ContextWithCorrelationID(ss.Context())
		newCtx = WithLogger(newCtx, log)
		newCtx = forwardLogger(newCtx, clear)
		// Add the correlation id to the grpc logrus so that we get the correlation id
		// in middleware logs
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

func forwardLogger(ctx context.Context, clear bool) context.Context {
	md := metautils.ExtractIncoming(ctx)
	if md.Get(correlationIDKey) == "" {
		md.Set(correlationIDKey, CorrelationIDFromContext(ctx))
	}
	lg := FromContext(ctx)
	for k, v := range md {
		if strings.HasPrefix(k, logLead) {
			md.Del(k)
			// only load keys with values
			// clear flag removes metadata entries for external services to remove user input
			if len(v) > 0 && !clear {
				lg = lg.WithField(strings.TrimPrefix(k, logLead), fmt.Sprint(v[0]))
			}
		}
	}
	grpc_ctxtags.Extract(ctx).Set(correlationIDKey, md.Get(correlationIDKey))
	ctx = md.ToIncoming(ctx)
	lg = WithContext(ctx, lg)
	return WithLogger(ctx, lg)
}

// generateCorrelationId generates a random uuid to use as the correlation id
func generateCorrelationID() string {
	return uuid.New().String()
}

// LoggerMiddleware returns a middleware that loads the logger into the context.
func LoggerMiddleware(logger logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(WithLogger(r.Context(), logger)))
		})
	}
}

// LoggerMiddlewareFunc returns a middleware that loads the logger into the context.
func LoggerMiddlewareFunc(logger logrus.FieldLogger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r.WithContext(WithLogger(r.Context(), logger)))
		}
	}
}
