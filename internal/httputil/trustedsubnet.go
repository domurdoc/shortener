package httputil

import (
	"net"
	"net/http"
)

// CheckSubnet is a middleware function that checks if the request comes from a trusted subnet.
func CheckSubnet(fn http.HandlerFunc, trustedSubnet string) http.HandlerFunc {
	if trustedSubnet == "" {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}
	}
	_, ipNET, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ipStr := r.Header.Get("X-Real-IP")
		if ipStr == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if !ipNET.Contains(ip) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		fn(w, r)
	}
}
