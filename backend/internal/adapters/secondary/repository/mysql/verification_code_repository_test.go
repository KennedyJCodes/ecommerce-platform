package repository_mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

type verificationCodeTestDriver struct {
	connection *verificationCodeTestConnection
}

func (d *verificationCodeTestDriver) Open(string) (driver.Conn, error) {
	return d.connection, nil
}

type verificationCodeTestConnection struct {
	mu      sync.Mutex
	query   string
	args    []driver.NamedValue
	result  driver.Result
	execErr error
}

func (c *verificationCodeTestConnection) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepared statements are not supported by this test driver")
}

func (c *verificationCodeTestConnection) Close() error { return nil }

func (c *verificationCodeTestConnection) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions are not supported by this test driver")
}

func (c *verificationCodeTestConnection) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.query = query
	c.args = append([]driver.NamedValue(nil), args...)
	if c.execErr != nil {
		return nil, c.execErr
	}
	return c.result, nil
}

func (c *verificationCodeTestConnection) snapshot() (string, []driver.NamedValue) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.query, append([]driver.NamedValue(nil), c.args...)
}

func TestNewSQLVerificationCodeRepositoryRejectsNilDatabase(t *testing.T) {
	if _, err := NewSQLVerificationCodeRepository(nil); err == nil {
		t.Fatal("NewSQLVerificationCodeRepository(nil) error = nil, want error")
	}
}

func TestSaveCodeVerication(t *testing.T) {
	connection := &verificationCodeTestConnection{result: driver.RowsAffected(1)}
	driverName := "verification-code-test"
	sql.Register(driverName, &verificationCodeTestDriver{connection: connection})
	database, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer database.Close()

	repository, err := NewSQLVerificationCodeRepository(sqlx.NewDb(database, driverName))
	if err != nil {
		t.Fatalf("NewSQLVerificationCodeRepository() error = %v", err)
	}

	expiresAt := time.Date(2026, time.September, 25, 12, 0, 0, 0, time.UTC)
	if err := repository.SaveCodeVerication(context.Background(), 42, "hashed-code", expiresAt); err != nil {
		t.Fatalf("SaveCodeVerication() error = %v", err)
	}

	query, args := connection.snapshot()
	if !strings.Contains(query, "INSERT INTO user_verification_codes") {
		t.Errorf("query = %q, want an INSERT into user_verification_codes", query)
	}
	if len(args) != 3 {
		t.Fatalf("argument count = %d, want 3", len(args))
	}
	if args[0].Value != int64(42) || args[1].Value != "hashed-code" || args[2].Value != expiresAt {
		t.Errorf("unexpected arguments: %+v", args)
	}
}

func TestSaveCodeVericationReturnsDatabaseError(t *testing.T) {
	connection := &verificationCodeTestConnection{
		result:  driver.RowsAffected(0),
		execErr: errors.New("database unavailable"),
	}
	driverName := "verification-code-error-test"
	sql.Register(driverName, &verificationCodeTestDriver{connection: connection})
	database, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer database.Close()

	repository, err := NewSQLVerificationCodeRepository(sqlx.NewDb(database, driverName))
	if err != nil {
		t.Fatalf("NewSQLVerificationCodeRepository() error = %v", err)
	}
	if err := repository.SaveCodeVerication(context.Background(), 42, "hashed-code", time.Now()); err == nil {
		t.Fatal("SaveCodeVerication() error = nil, want database error")
	}
}
