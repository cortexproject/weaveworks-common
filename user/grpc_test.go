package user

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestExtractFromGRPCRequest(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		expectedOrgID string
		expectedError error
		expectNewCtx  bool
	}{
		{
			name:          "successful extraction with single org ID",
			ctx:           metadata.NewIncomingContext(context.Background(), metadata.Pairs(lowerOrgIDHeaderName, "org-123")),
			expectedOrgID: "org-123",
			expectedError: nil,
			expectNewCtx:  true,
		},
		{
			name:          "no metadata in context",
			ctx:           context.Background(),
			expectedOrgID: "",
			expectedError: ErrNoOrgID,
			expectNewCtx:  false,
		},
		{
			name:          "missing org ID header",
			ctx:           metadata.NewIncomingContext(context.Background(), metadata.Pairs("other-header", "value")),
			expectedOrgID: "",
			expectedError: ErrNoOrgID,
			expectNewCtx:  false,
		},
		{
			name:          "multiple org IDs",
			ctx:           metadata.NewIncomingContext(context.Background(), metadata.Pairs(lowerOrgIDHeaderName, "org-123", lowerOrgIDHeaderName, "org-456")),
			expectedOrgID: "",
			expectedError: ErrNoOrgID,
			expectNewCtx:  false,
		},
		{
			name:          "empty org ID",
			ctx:           metadata.NewIncomingContext(context.Background(), metadata.Pairs(lowerOrgIDHeaderName, "")),
			expectedOrgID: "",
			expectedError: nil,
			expectNewCtx:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgID, newCtx, err := ExtractFromGRPCRequest(tt.ctx)

			// Check error
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err != tt.expectedError {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}

			// Check org ID
			if orgID != tt.expectedOrgID {
				t.Errorf("Expected org ID %q, got %q", tt.expectedOrgID, orgID)
			}

			// Check context
			if tt.expectNewCtx {
				if newCtx == tt.ctx {
					t.Error("Expected new context, got same context")
				}
				// Verify the org ID is properly injected into the new context
				extractedOrgID, extractErr := ExtractOrgID(newCtx)
				if extractErr != nil {
					t.Errorf("Failed to extract org ID from new context: %v", extractErr)
				} else if extractedOrgID != tt.expectedOrgID {
					t.Errorf("Expected extracted org ID %q, got %q", tt.expectedOrgID, extractedOrgID)
				}
			} else {
				if newCtx != tt.ctx {
					t.Error("Expected same context, got new context")
				}
			}
		})
	}
}

func BenchmarkExtractFromGRPCRequest(b *testing.B) {
	// Create a context with realistic metadata containing multiple headers
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		lowerOrgIDHeaderName, "org-123",
		"authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"user-agent", "grpc-go/1.50.1",
		"content-type", "application/grpc",
		"grpc-accept-encoding", "gzip",
		"grpc-encoding", "gzip",
		"x-forwarded-for", "192.168.1.100",
		"x-request-id", "req-12345",
		"x-trace-id", "trace-67890",
		"x-span-id", "span-abcdef",
		"x-custom-header", "custom-value",
		"x-api-version", "v1.2.3",
		"x-client-version", "1.0.0",
		"x-request-timeout", "30s",
		"x-retry-count", "0",
	))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _, err := ExtractFromGRPCRequest(ctx)
		if err != nil {
			b.Fatalf("ExtractFromGRPCRequest failed: %v", err)
		}
	}
}
