package specifications

type SignatureSpecificationByID struct {
	ID string
}

func (s SignatureSpecificationByID) ToSQL() (string, map[string]any) {
	query := `SELECT id, request_id, user_id, created_at FROM signatures
	WHERE id = @id`
	return query, map[string]any{"id": s.ID}
}

func NewSignatureSpecificationByID(id string) SignatureSpecificationByID {
	return SignatureSpecificationByID{id}
}