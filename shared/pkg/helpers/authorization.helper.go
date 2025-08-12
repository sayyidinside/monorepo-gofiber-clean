package helpers

import (
	"context"
)

func SelfOrAdminOnly(ctx context.Context, user_id uint) bool {
	session_user_id := ctx.Value(CtxKeyUserID).(float64)
	is_admin := ctx.Value(CtxKeyIsAdmin).(bool)
	if user_id != uint(session_user_id) && !is_admin {
		return false
	}

	return true
}
