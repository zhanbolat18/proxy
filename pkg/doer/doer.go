package doer

import "net/http"

type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}
