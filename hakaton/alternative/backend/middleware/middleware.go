package middleware

import (
	"backend/apiutils"
	"backend/store"
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func WithUserID(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, UserIDKey, id)
}

func GetUserID(ctx context.Context) (uint64, bool) {
	value := ctx.Value(UserIDKey)
	if value == nil {
		return 0, false
	}
	id, ok := value.(uint64)
	return id, ok
}

var allowed = map[string]bool{
	"http://localhost:8080":      true,
	"http://127.0.0.1:8080":      true,
	"http://localhost:8030":      true,
	"http://127.0.0.1:8030":      true,
	"http://89.208.210.115:8030": true,
	"http://localhost:3001":      true,
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(redisStore *store.RedisStore) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := r.Cookie("session_id")
			if errors.Is(err, http.ErrNoCookie) {
				log.Info().Msg("no session cookie found in auth middleware")
				apiutils.WriteError(w, http.StatusUnauthorized, "no session cookie")
				return
			}
			if err != nil {
				log.Error().Err(err).Msg("error getting session cookie in auth middleware")
				apiutils.WriteError(w, http.StatusInternalServerError, "internal server error")
				return
			}

			user, ok := redisStore.GetUserBySession(session.Value)
			if !ok {
				apiutils.WriteError(w, http.StatusUnauthorized, "invalid session")
				return
			}

			ctx := WithUserID(r.Context(), user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
