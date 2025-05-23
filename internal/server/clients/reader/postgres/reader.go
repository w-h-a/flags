package postgres

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/pkg/telemetry/log"
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
		log.Fatal(err)
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

	record := &record{}

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

	records := []*record{}

	for rows.Next() {
		record := &record{}

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

	readOne, err := c.conn.Prepare(`SELECT key, value FROM flags WHERE key = $1;`)
	if err != nil {
		log.Fatal(err)
	}
	c.readOne = readOne

	readMany, err := c.conn.Prepare(`SELECT key, value FROM flags;`)
	if err != nil {
		log.Fatal(err)
	}
	c.readMany = readMany

	return c
}
