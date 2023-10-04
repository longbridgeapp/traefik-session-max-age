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

type ResponseWriterWrapper struct {
	http.ResponseWriter
	cookieName string
	maxAge     int
}

func (rww ResponseWriterWrapper) Header() http.Header {
	return rww.ResponseWriter.Header()
}

func (rww ResponseWriterWrapper) WriteHeader(code int) {
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

func (rww ResponseWriterWrapper) Write(b []byte) (int, error) {
	return rww.ResponseWriter.Write(b)
}

type HeaderWrapper http.Header

func (h HeaderWrapper) Add(key, value string) {
	h.Add(key, value)
}

type SessionMaxAge struct {
	next       http.Handler
	cookieName string
	maxAge     int
	name       string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &SessionMaxAge{
		next:       next,
		name:       name,
		maxAge:     config.MaxAge,
		cookieName: config.CookieName,
	}, nil
}

func (a *SessionMaxAge) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.next.ServeHTTP(ResponseWriterWrapper{ResponseWriter: rw, cookieName: a.cookieName, maxAge: a.maxAge}, req)
}
