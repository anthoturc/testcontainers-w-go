package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB

// createPgContainer will standup a Postgres testcontainer for use across all tests
func createPgContainer(pgConf *PgConf) (*postgres.PostgresContainer, error) {
	container, err := postgres.RunContainer(context.Background(),
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithInitScripts(filepath.Join("db", "ips.sql")),
		postgres.WithDatabase(pgConf.Database),
		postgres.WithPassword(pgConf.Password),
		postgres.WithUsername(pgConf.UserName),
		postgres.WithDatabase(pgConf.Database),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	return container, err
}

// TestMain will setup the dependencies needed for all tests
func TestMain(m *testing.M) {
	pgConf := DefaultPgConf()
	container, err := createPgContainer(pgConf)
	if err != nil {
		log.Fatalf("unable to start container: %+v", err)
	}

	defer func() {
		if err := container.Terminate(context.Background()); err != nil {
			log.Fatalf("unable to stop container: %+v", err)
		}
	}()

	connStr, err := container.ConnectionString(context.Background(), fmt.Sprintf("sslmode=%s", pgConf.SSLMode), fmt.Sprintf("dbname=%s", pgConf.Database))
	if err != nil {
		log.Fatalf("Failed to connect to Pg %s", connStr)
	}

	db, err = Open(connStr)
	if err != nil {
		log.Fatal("couldn't open with connstr")
	}
	defer db.Close()

	code := m.Run()

	os.Exit(code)
}

func TestPing(t *testing.T) {
	ps := &PingService{
		DB: db,
	}

	ipAddr := "12.123.1.1"
	ping, err := ps.Ping(ipAddr)
	if err != nil {
		t.Errorf("unexpected failure: %+v", err)
		t.FailNow()
	}

	assert.NotNil(t, ping)
	assert.NotNil(t, ping.IpAddr)
	assert.Equal(t, ipAddr, ping.IpAddr)
	assert.True(t, ping.Id > -1)
}

func TestPingWithIpAlreadyExists(t *testing.T) {
	ps := &PingService{
		DB: db,
	}

	ipAddr := "123.123.123.1"
	_, err := ps.Ping(ipAddr)
	assert.NoError(t, err)

	_, err = ps.Ping(ipAddr)
	assert.Error(t, err, "There m be an error!")
	assert.ErrorAs(t, err, &ErrIpAlreadyExists, "The error should be an ErrIpAlreadyExists")
}
