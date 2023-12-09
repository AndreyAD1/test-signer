package repositories

import "context"

type SignatureRepository interface {
	Add(context.Context, Signature) (*Signature, error)
	Query(context.Context, Specification) ([]Signature, error)
}

type Specification interface {
	ToSQL() (string, map[string]any)
}
