// Package traefik_session_max_age is a plugin for the Traefik reverse proxy
// that sets cookie's max-age
package traefik_session_max_age

import (
	"context"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	CookieName string `json:"cookieName,omitempty"`
	MaxAge     int    `json:"maxAge,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type responseWriterWrapper struct {
	http.ResponseWriter
	cookieName string
	maxAge     int
}

func (rww responseWriterWrapper) Header() http.Header {
	return rww.ResponseWriter.Header()
}

func (rww responseWriterWrapper) WriteHeader(code int) {
	if rww.cookieName != "" {
		res := http.Response{Header: rww.ResponseWriter.Header()}
		cookies := res.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == rww.cookieName {
				cookie.MaxAge = rww.maxAge
				http.SetCookie(rww.ResponseWriter, cookie)
			}
		}
	}
	rww.ResponseWriter.WriteHeader(code)
}

func (rww responseWriterWrapper) Write(b []byte) (int, error) {
	return rww.ResponseWriter.Write(b)
}

// SessionMaxAge is a middleware for traefik middlware plugin to set cookie max age.
type SessionMaxAge struct {
	next       http.Handler
	cookieName string
	maxAge     int
	name       string
}

// New return a wrapped http.Handler.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &SessionMaxAge{
		next:       next,
		name:       name,
		maxAge:     config.MaxAge,
		cookieName: config.CookieName,
	}, nil
}

func (a *SessionMaxAge) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.next.ServeHTTP(responseWriterWrapper{ResponseWriter: rw, cookieName: a.cookieName, maxAge: a.maxAge}, req)
}
