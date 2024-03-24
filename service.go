package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type PingService struct {
	DB *sql.DB
}

type Ping struct {
	Id     int
	IpAddr string
}

var ErrIpAlreadyExists error = errors.New("ip address already exists")

func (ps *PingService) Ping(ipAddr string) (*Ping, error) {
	row := ps.DB.QueryRow(
		`INSERT INTO ips (ip_addr)
		VALUES ($1) RETURNING id`, ipAddr)
	ping := &Ping{
		IpAddr: ipAddr,
	}
	err := row.Scan(&ping.Id)

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && pgerrcode.IsIntegrityConstraintViolation(e.Code) && e.Code == pgerrcode.UniqueViolation {
			return nil, ErrIpAlreadyExists
		}

		return nil, fmt.Errorf("create ping entry: %w", err)
	}

	return ping, nil
}
