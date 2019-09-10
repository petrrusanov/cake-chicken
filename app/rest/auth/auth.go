package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/boilerplate/backend/app/rest"
	"github.com/boilerplate/backend/app/store/service"
)

// Authenticator auth config struct
type Authenticator struct {
	DataStore *service.DataStore
}

// Auth middleware populates user info
func (a *Authenticator) Auth(requireAuth bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeaderComponents := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

			if len(authHeaderComponents) != 2 || authHeaderComponents[0] != "Bearer" {
				log.Printf("[WARN] auth failed, incorrect auth header %s", r.Header.Get("Authorization"))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := authHeaderComponents[1]

			user, err := a.DataStore.FindUser(token)

			if err != nil {
				log.Printf("[WARN] auth failed, user not found, %s", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			r = rest.SetUserInfo(r, user)

			_ = rest.MustGetUserInfo(r)

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
