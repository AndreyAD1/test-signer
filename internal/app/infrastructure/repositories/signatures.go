package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	transaction, err := r.dbPool.Begin(ctx)
	if err != nil {
		log.Println("can not begin a transaction")
		return nil, err
	}
	defer func() {
		err := transaction.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf(
				"can not finish a transaction for a signature '%s'",
				signature.RequestID,
			)
		}
	}()
	insertQuery := `INSERT INTO signatures (id, request_id, user_id, created_at)
	VALUES ($1, $2, $3, $4) RETURNING id, request_id, user_id, created_at;`
	var savedSignature Signature
	err = transaction.QueryRow(
		ctx,
		insertQuery,
		signature.ID,
		signature.RequestID,
		signature.UserID,
		signature.CreatedAt,
	).Scan(
		&savedSignature.ID,
		&savedSignature.RequestID,
		&savedSignature.UserID,
		&savedSignature.CreatedAt,
	)
	if err != nil {
		var pgxError *pgconn.PgError
		if !errors.As(err, &pgxError) {
			log.Printf("unexpected DB error: %v", err)
			return nil, err
		}
		if pgxError.Code == pgerrcode.UniqueViolation {
			log.Printf("the signature already exists: %v", signature.RequestID)
			return nil, ErrDuplicate
		}
		return nil, err
	}
	insertAnswerQuery := `INSERT INTO test_details (signature_id, question, answer)
	VALUES ($1, $2, $3) RETURNING id, question, answer;`
	savedAnswers := []TestDetails{}
	for _, answer := range signature.Answers {
		var savedTestDetails TestDetails
		err := transaction.QueryRow(
			ctx, 
			insertAnswerQuery, 
			savedSignature.ID, 
			answer.Question, 
			answer.Answer,
		).Scan(
			&savedTestDetails.ID,
			&savedTestDetails.Question,
			&savedTestDetails.Answer,
		)
		if err != nil {
			log.Printf("unexpected DB error: %v", err)
			return nil, err
		}
		savedAnswers = append(savedAnswers, savedTestDetails)
	}
	savedSignature.Answers = savedAnswers
	if err := transaction.Commit(ctx); err != nil {
		log.Printf(
			"can not close a transaction for the signature %v - %v: %v",
			signature.RequestID,
			signature.UserID,
			err,
		)
		return &savedSignature, err
	}
	return &savedSignature, err
}

func (r *SignatureCollection) Query(ctx context.Context, spec Specification) ([]Signature, error) {
	return []Signature{}, ErrNotImplemented
}
