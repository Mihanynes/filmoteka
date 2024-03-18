package auth

import (
	"context"
	"net/http"
	"strings"
)

func AdminAuthMiddleware(sm SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/admin/") {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)
		if err != nil {
			http.Error(w, "No admin auth", http.StatusUnauthorized)
			return
		}

		if !sess.IsAdmin {
			http.Error(w, "Not an admin", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
