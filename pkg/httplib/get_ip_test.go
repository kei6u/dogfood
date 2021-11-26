package httplib

import (
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetIP(t *testing.T) {
	tests := []struct {
		name       string
		r          *http.Request
		setHeaders []func(r *http.Request)
		want       net.IP
	}{
		{
			name: "got from r.RemoteAddr",
			r:    httptest.NewRequest(http.MethodGet, "/test", nil),
			want: net.ParseIP("192.0.2.1"),
		},
		{
			name: "got from X-Forwarded-For",
			r:    httptest.NewRequest(http.MethodGet, "/test", nil),
			setHeaders: []func(r *http.Request){
				func(r *http.Request) {
					r.RemoteAddr = ""
					r.Header.Set("X-Forwarded-For", "192.0.2.1")
				},
			},
			want: net.ParseIP("192.0.2.1"),
		},
		{
			name: "got from x-forwarded-for",
			r:    httptest.NewRequest(http.MethodGet, "/test", nil),
			setHeaders: []func(r *http.Request){
				func(r *http.Request) {
					r.RemoteAddr = ""
					r.Header.Set("x-forwarded-for", "192.0.2.1")
				},
			},
			want: net.ParseIP("192.0.2.1"),
		},
		{
			name: "got from X-FORWARDED-FOR",
			r:    httptest.NewRequest(http.MethodGet, "/test", nil),
			setHeaders: []func(r *http.Request){
				func(r *http.Request) {
					r.RemoteAddr = ""
					r.Header.Set("X-FORWARDED-FOR", "192.0.2.1")
				},
			},
			want: net.ParseIP("192.0.2.1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, set := range tt.setHeaders {
				set(tt.r)
			}
			if got := GetIP(tt.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
