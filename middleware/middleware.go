package middleware

import (
	"net/http"
)

type Middleware interface {
	Open(name string) (http.File, error)
	SetNext(middleware http.FileSystem)
}

type FileSystemMiddleware struct {
	FileSystem  http.FileSystem
	Middlewares []Middleware
}

func CreateChain(middlewares []Middleware, final http.FileSystem) http.FileSystem {
	m := middlewares[0]
	for i := 1; i < len(middlewares); i++ {
		m.SetNext(middlewares[i])
		m = middlewares[i]
	}
	m.SetNext(final)
	return middlewares[0]
}
