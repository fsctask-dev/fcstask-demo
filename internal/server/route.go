package server

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type RouterBuilder struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

type VersionBuilder struct {
	prefix      string
	mux         *http.ServeMux
	middlewares []Middleware
	parent      *RouterBuilder
}

func NewRouterBuilder() *RouterBuilder {
	return &RouterBuilder{
		mux: http.NewServeMux(),
	}
}

func (b *RouterBuilder) WithMiddleware(m Middleware) *RouterBuilder {
	b.middlewares = append(b.middlewares, m)
	return b
}

func (b *RouterBuilder) Version(prefix string) *VersionBuilder {
	return &VersionBuilder{
		prefix: prefix,
		mux:    http.NewServeMux(),
		parent: b,
	}
}

func (b *RouterBuilder) Build() http.Handler {
	var handler http.Handler = b.mux

	for i := len(b.middlewares) - 1; i >= 0; i-- {
		handler = b.middlewares[i](handler)
	}

	return handler
}

func (v *VersionBuilder) WithMiddleware(m Middleware) *VersionBuilder {
	v.middlewares = append(v.middlewares, m)
	return v
}

func (v *VersionBuilder) WithEcho(h http.HandlerFunc) *VersionBuilder {
	v.mux.HandleFunc("/echo", h)
	return v
}

func (v *VersionBuilder) Register() *RouterBuilder {
	var handler http.Handler = v.mux

	for i := len(v.middlewares) - 1; i >= 0; i-- {
		handler = v.middlewares[i](handler)
	}

	v.parent.mux.Handle(v.prefix+"/", http.StripPrefix(v.prefix, handler))
	return v.parent
}
