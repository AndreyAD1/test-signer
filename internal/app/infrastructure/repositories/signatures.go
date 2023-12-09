package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SignatureCollection struct {
	dbPool *pgxpool.Pool
}

func NewSignatureCollection(ctx context.Context, dbURL string) (*SignatureCollection, error) {
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to create a connection pool: DB URL '%s': %w",
			dbURL,
			err,
		)
	}
	if err := dbPool.Ping(ctx); err != nil {
		log.Printf("unable to connect to the DB '%v'", dbURL)
		return nil, err
	}
	return &SignatureCollection{dbPool}, nil
}

func (r *SignatureCollection) Add(ctx context.Context, signature Signature) (*Signature, error) {
	return &signature, nil
}

func (r *SignatureCollection) Query(ctx context.Context, spec Specification) ([]Signature, error) {
	return []Signature{}, ErrNotImplemented
}