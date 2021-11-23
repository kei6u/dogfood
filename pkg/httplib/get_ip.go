package httplib

import (
	"net"
	"net/http"
)

// GetIP gets ip address from http.Request.
func GetIP(r *http.Request) net.IP {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		for _, h := range []string{"X-Forwarded-For", "x-forwarded-for", "X-FORWARDED-FOR"} {
			if ip = r.Header.Get(h); ip != "" {
				break
			}
		}
	}
	return net.ParseIP(ip)
}
