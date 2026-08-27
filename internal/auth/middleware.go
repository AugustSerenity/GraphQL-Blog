package auth

import "net/http"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")

		if userID == "" {
			http.Error(w, "user is not authenticated", http.StatusUnauthorized)
			return
		}

		ctx := WithUserID(r.Context(), userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
