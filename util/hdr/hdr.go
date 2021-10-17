package hdr

import (
	"net/http"
	"strings"
)

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := []string{"Content-Type", "Accept", "Authorization","Access-Control-Allow-Headers","X-Requested-With"}
		methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		}

		h.ServeHTTP(w, r)
	})
}
