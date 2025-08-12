package helpers

type contextKey string

const (
	CtxKeyIdentifier contextKey = "identifier"
	CtxKeyUsername   contextKey = "username"
	CtxKeyUserID     contextKey = "user_id"
	CtxKeyIsAdmin    contextKey = "is_admin"
	CtxKeyFunction   contextKey = "function"
)
