package geo

import (
	"context"
	"testing"
	"time"
)

func TestLoaderRun(t *testing.T) {
	loader := NewLoader(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := loader.Run(ctx); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}
}
