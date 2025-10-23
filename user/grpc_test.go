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
