package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
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
		detail := "failed to register postgres reader with otel"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	DRIVER = driver
}

type client struct {
	options  reader.Options
	conn     *sql.DB
	readOne  *sql.Stmt
	readMany *sql.Stmt
}

func (c *client) ReadByKey(ctx context.Context, key string) ([]byte, error) {
	row := c.readOne.QueryRowContext(ctx, key)

	record := &reader.Record{}

	if err := row.Scan(&record.Key, &record.Value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []byte{}, reader.ErrRecordNotFound
		}
		return nil, err
	}

	return record.Value, nil
}

func (c *client) Read(ctx context.Context) ([]byte, error) {
	rows, err := c.readMany.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	records := []*reader.Record{}

	for rows.Next() {
		record := &reader.Record{}

		if err := rows.Scan(&record.Key, &record.Value); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return []byte{}, nil
			}
			return nil, err
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := []byte{}

	for _, record := range records {
		result = append(result, []byte("\n")...)
		result = append(result, record.Value...)
	}

	return result, nil
}

func NewReader(opts ...reader.Option) reader.Reader {
	options := reader.NewOptions(opts...)

	if err := options.Validate(); err != nil {
		detail := "failed to validate postgres reader options"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	c := &client{
		options: options,
	}

	// postgres://user:password@host:port/db?sslmode=disable
	conn, err := sql.Open(DRIVER, c.options.Location)
	if err != nil {
		detail := "failed to connect with postgres reader"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	if err := conn.Ping(); err != nil {
		detail := "failed to ping with postgres reader"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	if err := otelsql.RecordStats(conn); err != nil {
		detail := "failed to initialize postgres instrumentation for postgres reader"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	c.conn = conn

	if _, err := c.conn.Exec(`CREATE TABLE IF NOT EXISTS flags (key text NOT NULL, value bytea, CONSTRAINT flags_pkey PRIMARY KEY (key));`); err != nil {
		detail := "failed to create table for postgres reader"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	if _, err := c.conn.Exec(`CREATE INDEX IF NOT EXISTS key_index_flags ON flags (key);`); err != nil {
		detail := "failed to create index for postgres reader"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}

	readOne, err := c.conn.Prepare(`SELECT key, value FROM flags WHERE key = $1;`)
	if err != nil {
		detail := "failed to prepare select statement for postgres reader"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}
	c.readOne = readOne

	readMany, err := c.conn.Prepare(`SELECT key, value FROM flags;`)
	if err != nil {
		detail := "failed to prepare select statement for postgres reader"
		slog.ErrorContext(context.Background(), detail, "error", err)
		panic(detail)
	}
	c.readMany = readMany

	return c
}
