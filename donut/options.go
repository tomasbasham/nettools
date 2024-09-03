package donut

import "net/http"

type option func(r *Resolver)

func WithDebug() option {
	return func(r *Resolver) {
		r.debug = true
	}
}

func WithClient(c *http.Client) option {
	return func(r *Resolver) {
		r.client = c
	}
}
