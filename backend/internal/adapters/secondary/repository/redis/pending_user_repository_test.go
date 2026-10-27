package repository_redis

import (
	"context"
	"testing"
	"time"

	modelsdb "github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/database"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newPendingUserRepositoryForTest(t *testing.T) (*RedisPendingUserRepository, *redis.Client) {
	t.Helper()
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	return NewRedisPendingUserRepository(client).(*RedisPendingUserRepository), client
}

func TestPendingUserRepositorySaveGetAndDelete(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	repository := NewRedisPendingUserRepository(client)
	pendingUser := &modelsdb.PendingUser{
		ID:           "pending-1",
		Username:     "alice",
		PasswordHash: "password-hash",
		Email:        "alice@example.com",
		HashCode:     "code-hash",
		Attempts:     2,
	}

	ctx := context.Background()
	if err := repository.SavePendingUser(ctx, pendingUser, time.Minute); err != nil {
		t.Fatalf("SavePendingUser() error = %v", err)
	}
	hashCode, err := repository.GetPendingUserHashCode(ctx, pendingUser.ID)
	if err != nil || hashCode != pendingUser.HashCode {
		t.Fatalf("GetPendingUserHashCode() = %q, %v; want %q, nil", hashCode, err, pendingUser.HashCode)
	}
	got, err := repository.GetPendingUser(ctx, pendingUser.ID)
	if err != nil {
		t.Fatalf("GetPendingUser() error = %v", err)
	}
	if *got != *pendingUser {
		t.Fatalf("GetPendingUser() = %+v, want %+v", got, pendingUser)
	}

	if err := repository.DeletePendingUser(ctx, pendingUser.ID); err != nil {
		t.Fatalf("DeletePendingUser() error = %v", err)
	}
	if _, err := repository.GetPendingUserHashCode(ctx, pendingUser.ID); err == nil {
		t.Fatal("GetPendingUserHashCode() error = nil after deletion")
	}
}

func TestPendingUserRepositoryMissingUser(t *testing.T) {
	repository, _ := newPendingUserRepositoryForTest(t)

	ctx := context.Background()
	if _, err := repository.GetPendingUser(ctx, "missing"); err == nil {
		t.Fatal("GetPendingUser() error = nil for missing user")
	}
	if _, err := repository.GetPendingUserHashCode(ctx, "missing"); err == nil {
		t.Fatal("GetPendingUserHashCode() error = nil for missing user")
	}
	if _, err := repository.RegisterVerificationAttempt(ctx, "missing", 5); err == nil {
		t.Fatal("RegisterVerificationAttempt() error = nil for missing user")
	}
}

func TestPendingUserRepositoryRegisterVerificationAttempt(t *testing.T) {
	repository, client := newPendingUserRepositoryForTest(t)
	pendingUser := &modelsdb.PendingUser{
		ID:       "pending-attempts",
		HashCode: "code-hash",
		Attempts: 0,
	}
	ctx := context.Background()
	if err := repository.SavePendingUser(ctx, pendingUser, time.Minute); err != nil {
		t.Fatalf("SavePendingUser() error = %v", err)
	}

	for want, i := 1, 1; i <= 5; want, i = want+1, i+1 {
		got, err := repository.RegisterVerificationAttempt(ctx, pendingUser.ID, 5)
		if err != nil || got != want {
			t.Fatalf("RegisterVerificationAttempt() = %d, %v; want %d, nil", got, err, want)
		}
	}

	if _, err := repository.RegisterVerificationAttempt(ctx, pendingUser.ID, 5); err == nil {
		t.Fatal("RegisterVerificationAttempt() error = nil after max attempts")
	}
	if exists, err := client.Exists(context.Background(), "pending_user:"+pendingUser.ID).Result(); err != nil || exists != 0 {
		t.Fatalf("pending user exists after exceeding attempts = %d, %v; want 0, nil", exists, err)
	}
}
