package user

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func BenchmarkExtractFromGRPCRequest(b *testing.B) {
	// Simulate a realistic incoming context with multiple metadata keys
	md := metadata.MD{
		"content-type":         []string{"application/grpc"},
		"user-agent":          []string{"grpc-go/1.71.2"},
		"x-scope-orgid":      []string{"tenant-12345"},
		"x-forwarded-for":    []string{"10.0.0.1"},
		"x-request-id":       []string{"abc-def-ghi-jkl"},
		"traceparent":        []string{"00-abcdef1234567890abcdef1234567890-1234567890abcdef-01"},
		"grpc-accept-encoding": []string{"gzip"},
		"authorization":      []string{"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."},
		"x-amz-date":         []string{"20260507T183000Z"},
		"x-amz-security-token": []string{"FwoGZXIvYXdzEBYaDH..."},
	}
	ctx := metadata.NewIncomingContext(context.Background(), md)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ExtractFromGRPCRequest(ctx)
	}
}
