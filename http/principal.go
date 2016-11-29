package gorgonzola

import "net/http"

type Principal struct {
	next ContextHandler
}

func (p Principal) ServeCtxHTTP(w http.ResponseWriter, r *http.Request, c *Context) {
	if pr := r.Header.Get("X-Principal"); pr != "" {
		p.next(w, r, c)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
