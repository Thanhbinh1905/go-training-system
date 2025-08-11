package contextkey

type ctxKey string

const (
	ctxUserIDKey   ctxKey = "userID"
	ctxUserRoleKey ctxKey = "userRole"
)

func CtxUserIDKey() ctxKey {
	return ctxUserIDKey
}

func CtxUserRoleKey() ctxKey {
	return ctxUserRoleKey
}
