package authorization

import (
	"fmt"
	"kubazulo/pkg/utils"
	"log"
	"net/http"
	"strings"
)

func startLocalListener(c utils.AuthorizationConfig, token *AuthorizationCode) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%s", c.RedirectPort)}

	http.HandleFunc(c.RedirectPath, func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatalf("Error while parsing form from response %s", err)
			return
		}
		for k, v := range r.Form {
			if k == "code" {
				token.Value = strings.Join(v, "")
			}
		}

		fmt.Fprintf(w, "%s", utils.SuccessMsg)
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			//log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}
