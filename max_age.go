// Package traefik_session_max_age is a plugin for the Traefik reverse proxy
// that sets cookie's max-age
package traefik_session_max_age

import (
	"context"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
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
	cookies := rww.ResponseWriter.Header()["Set-Cookie"]
	if len(cookies) == 0 {
		rww.ResponseWriter.WriteHeader(code)
		return
	}
	for i, line := range cookies {
		parts := strings.Split(textproto.TrimString(line), ";")
		if len(parts) == 1 && parts[0] == "" {
			continue
		}
		parts[0] = textproto.TrimString(parts[0])
		name, _, ok := strings.Cut(parts[0], "=")
		if !ok {
			continue
		}
		buf := make([]byte, 0, 19)
		name = textproto.TrimString(name)
		if name == rww.cookieName {
			var b strings.Builder
			b.WriteString(line)
			if rww.maxAge > 0 {
				b.WriteString("; Max-Age=")
				b.Write(strconv.AppendInt(buf, int64(rww.maxAge), 10))
			} else if rww.maxAge < 0 {
				b.WriteString("; Max-Age=0")
			}
			cookies[i] = b.String()
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
