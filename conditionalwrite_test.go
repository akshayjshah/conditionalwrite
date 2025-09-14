package conditionalwrite

import (
	"fmt"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestConditionalWrite(t *testing.T) {
	// The MinIO testcontainers module requires a running (or socket-activated)
	// Docker daemon.
	const user, password = "admin", "password"
	mc, err := minio.Run(
		t.Context(),
		"minio/minio:RELEASE.2025-07-23T15-54-02Z",
		minio.WithUsername(user),
		minio.WithPassword(password),
	)
	if err != nil {
		t.Fatalf("start MinIO container: %v", err)
	}
	addr, err := mc.ConnectionString(t.Context())
	if err != nil {
		t.Fatalf("get MinIO connection string: %v", err)
	}

	c := NewClient(
		fmt.Sprintf("http://%s", addr), // endpoint
		user,
		password,
		"us-east-1", // region
		"test",      // bucket
		"text.txt",  // key
	)

	err = c.CreateBucket(t.Context())
	if err != nil {
		t.Fatalf("create bucket failed: %v", err)
	}

	etag, err := c.Set(t.Context(), strings.NewReader("one"), None)
	if err != nil {
		t.Fatalf("initial write failed: %v", err)
	}

	_, err = c.Set(t.Context(), strings.NewReader("two"), etag)
	if err != nil {
		t.Fatalf("overwrite with correct ETag failed: %v", err)
	}

	_, err = c.Set(t.Context(), strings.NewReader("three"), etag)
	if err == nil {
		t.Fatal("overwrite with incorrect ETag succeeded")
	}
	if !IsPreconditionFailed(err) {
		t.Fatalf("expected PreconditionFailed error, got %v", err)
	}
}
