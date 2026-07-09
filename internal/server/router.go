package server

import (
	"fmt"
	"net/http"
)

type ApiVers string

var (
	ApiVers1 = ApiVers("v1")
	ApiVers2 = ApiVers("v2")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVers
}

func NewApiVersionRouter(apiversion ApiVers) *APIVersionRouter {

	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiversion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.Handler)
	}

}
