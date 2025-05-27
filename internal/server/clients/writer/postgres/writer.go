package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/w-h-a/flags/internal/server/clients/writer"
	"github.com/w-h-a/pkg/telemetry/log"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var DRIVER string

func init() {
	driver, err := otelsql.Register(
		"postgres",
		otelsql.TraceQueryWithoutArgs(),
		otelsql.TraceRowsClose(),
		otelsql.TraceRowsAffected(),
		otelsql.WithSystem(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		log.Fatal(err)
	}

	DRIVER = driver
}

type client struct {
	options writer.Options
	conn    *sql.DB
	write   *sql.Stmt
}

func (c *client) Write(ctx context.Context, key string, bs []byte) error {
	if _, err := c.write.ExecContext(ctx, key, bs); err != nil {
		return err
	}

	return nil
}

func NewWriter(opts ...writer.Option) writer.Writer {
	options := writer.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		log.Fatal(err)
	}

	c := &client{
		options: options,
	}

	// postgres://user:password@host:port/db?sslmode=disable
	conn, err := sql.Open(DRIVER, c.options.Location)
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := otelsql.RecordStats(conn); err != nil {
		log.Fatal(err)
	}

	c.conn = conn

	if _, err := c.conn.Exec(`CREATE TABLE IF NOT EXISTS flags (key text NOT NULL, value bytea, CONSTRAINT flags_pkey PRIMARY KEY (key));`); err != nil {
		log.Fatal(err)
	}

	if _, err := c.conn.Exec(`CREATE INDEX IF NOT EXISTS key_index_flags ON flags (key);`); err != nil {
		log.Fatal(err)
	}

	write, err := c.conn.Prepare(`INSERT INTO flags (key, value) VALUES ($1, $2::bytea) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`)
	if err != nil {
		log.Fatal(err)
	}
	c.write = write

	return c
}
