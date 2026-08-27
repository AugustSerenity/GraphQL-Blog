package auth

import "net/http"

const testUserID = "user-1"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithUserID(r.Context(), testUserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
