package repository_redis

import (
	"testing"
	"time"

	modelsdb "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/database"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

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

	if err := repository.SavePendingUser(pendingUser, time.Minute); err != nil {
		t.Fatalf("SavePendingUser() error = %v", err)
	}
	hashCode, err := repository.GetPendingUserHashCode(pendingUser.ID)
	if err != nil || hashCode != pendingUser.HashCode {
		t.Fatalf("GetPendingUserHashCode() = %q, %v; want %q, nil", hashCode, err, pendingUser.HashCode)
	}
	got, err := repository.GetPendingUser(pendingUser.ID)
	if err != nil {
		t.Fatalf("GetPendingUser() error = %v", err)
	}
	if *got != *pendingUser {
		t.Fatalf("GetPendingUser() = %+v, want %+v", got, pendingUser)
	}

	if err := repository.DeletePendingUser(pendingUser.ID); err != nil {
		t.Fatalf("DeletePendingUser() error = %v", err)
	}
	if _, err := repository.GetPendingUserHashCode(pendingUser.ID); err == nil {
		t.Fatal("GetPendingUserHashCode() error = nil after deletion")
	}
}

func TestPendingUserRepositoryMissingUser(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	repository := NewRedisPendingUserRepository(client)

	if _, err := repository.GetPendingUser("missing"); err == nil {
		t.Fatal("GetPendingUser() error = nil for missing user")
	}
	if _, err := repository.GetPendingUserHashCode("missing"); err == nil {
		t.Fatal("GetPendingUserHashCode() error = nil for missing user")
	}
}
